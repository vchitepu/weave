package server

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/vchitepu/weave/internal/htmlrenderer"
	"github.com/vchitepu/weave/internal/theme"
	"github.com/vchitepu/weave/internal/watcher"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	goldrenderer "github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// Input represents a Markdown source. Either Path or Data is set.
type Input struct {
	Path string // file path (re-read on each request)
	Data []byte // pre-read bytes (stdin)
}

// Start launches the web server and blocks until interrupted.
func Start(inputs []Input, th theme.Theme, port int) error {
	h := newHandler(inputs, th)

	// Collect file paths for watching
	var paths []string
	for _, in := range inputs {
		if in.Path != "" {
			paths = append(paths, in.Path)
		}
	}

	// Start file watcher
	stop, err := watcher.Watch(paths, func() {
		h.broadcast()
	})
	if err != nil {
		return fmt.Errorf("failed to start file watcher: %w", err)
	}
	defer stop()

	// Start HTTP server
	addr := fmt.Sprintf(":%d", port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %w", port, err)
	}

	fmt.Fprintf(os.Stderr, "Weave web viewer running at http://localhost:%d\nPress Ctrl+C to stop.\n", port)

	srv := &http.Server{Handler: h}

	go func() {
		if err := srv.Serve(ln); err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		}
	}()

	// Block until signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	fmt.Fprintln(os.Stderr)

	return srv.Close()
}

type handler struct {
	inputs []Input
	th     theme.Theme

	mu      sync.Mutex
	clients map[chan struct{}]struct{}
}

func newHandler(inputs []Input, th theme.Theme) *handler {
	return &handler{
		inputs:  inputs,
		th:      th,
		clients: make(map[chan struct{}]struct{}),
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/", "":
		h.handleRoot(w, r)
	case "/events":
		h.handleEvents(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *handler) handleRoot(w http.ResponseWriter, r *http.Request) {
	htmlR := htmlrenderer.New(h.th)
	md := goldmark.New(
		goldmark.WithExtensions(extension.Table, extension.Strikethrough, extension.TaskList),
		goldmark.WithRenderer(
			goldrenderer.NewRenderer(
				goldrenderer.WithNodeRenderers(
					util.Prioritized(htmlR, htmlrenderer.Priority),
				),
			),
		),
	)

	var body bytes.Buffer
	for i, input := range h.inputs {
		if i > 0 {
			name := filepath.Base(input.Path)
			if name == "" || name == "." {
				name = fmt.Sprintf("input-%d", i+1)
			}
			body.WriteString(fmt.Sprintf("<hr class=\"file-separator\">\n<p class=\"file-label\">%s</p>\n",
				htmlrenderer.HtmlEscapeString(name)))
		}

		data, err := h.readInput(input)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading input: %v", err), http.StatusInternalServerError)
			return
		}

		var rendered bytes.Buffer
		if err := md.Convert(data, &rendered); err != nil {
			http.Error(w, fmt.Sprintf("Error rendering markdown: %v", err), http.StatusInternalServerError)
			return
		}
		body.Write(rendered.Bytes())
	}

	page := buildPage(h.th, body.String())

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(page))
}

func (h *handler) readInput(input Input) ([]byte, error) {
	if input.Path != "" {
		return os.ReadFile(input.Path)
	}
	return input.Data, nil
}

func (h *handler) handleEvents(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	flusher.Flush()

	ch := make(chan struct{}, 1)
	h.mu.Lock()
	h.clients[ch] = struct{}{}
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.clients, ch)
		h.mu.Unlock()
	}()

	for {
		select {
		case <-ch:
			fmt.Fprintf(w, "data: reload\n\n")
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

func (h *handler) broadcast() {
	h.mu.Lock()
	defer h.mu.Unlock()
	for ch := range h.clients {
		select {
		case ch <- struct{}{}:
		default:
			// Don't block if client is slow
		}
	}
}

func buildPage(th theme.Theme, body string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>weave</title>
<style>
%s
%s
</style>
</head>
<body>
%s
<script>
const es = new EventSource("/events");
es.onmessage = function(e) {
  if (e.data === "reload") {
    location.reload();
  }
};
</script>
</body>
</html>`, htmlrenderer.ThemeCSS(th), htmlrenderer.BaseCSS, body)
}

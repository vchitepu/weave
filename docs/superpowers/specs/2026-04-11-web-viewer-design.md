# Web Viewer Design (`--web`)

**Date:** 2026-04-11
**Status:** Approved

## Overview

Add a `--web` flag to `weave` that launches a local HTTP server and renders Markdown in a browser instead of the terminal. The server watches source files for changes and pushes live reloads to the browser via server-sent events (SSE). The feature supports all existing input modes: single file, multiple files, and stdin.

## Architecture

Three new internal packages are added, plus CLI changes:

```
internal/
  htmlrenderer/   — goldmark NodeRenderer emitting HTML+CSS
  server/         — HTTP server, SSE broadcaster, page template
  watcher/        — fsnotify wrapper for file change detection
cmd/weave/
  main.go         — --web and --port flags, run() dispatch
```

### Data flow

```
CLI args / stdin
      ↓
  run() in main.go
      ↓  (--web set)
  server.Start(inputs, theme, port)
      ↓
  GET /  ──► htmlrenderer.Render(inputs, theme) ──► full HTML page
  GET /events ──► SSE stream (holds open)
      ↑
  watcher ──► file change ──► broadcaster.Notify() ──► "data: reload\n\n"
```

## Section 1: CLI Integration

### New flags

| Flag | Type | Default | Description |
|---|---|---|---|
| `--web` | bool | `false` | Launch web viewer instead of terminal output |
| `--port` | int | `7331` | HTTP port; silently ignored if `--web` is not set |

### `run()` changes

After resolving theme and reading inputs, if `--web` is set, `run()` calls `server.Start(inputs, th, port)` and returns its result. The existing terminal render path is completely unchanged.

### Startup output

Printed to stderr (not stdout, to stay pipe-friendly):

```
Weave web viewer running at http://localhost:7331
Press Ctrl+C to stop.
```

### Flag interactions

- `--web` + `--width`: width is silently ignored in web mode (browser handles layout)
- `--web` + `--theme`: honoured — passed through to htmlrenderer and CSS custom properties
- `--web` + stdin: works — bytes captured once at startup, no live reload (nothing to watch)
- `--web` + multiple files: all files rendered in sequence with an `<hr>` + filename label separator

### Error cases

| Situation | Behaviour |
|---|---|
| Port already in use | Print error to stderr, exit non-zero |
| File not found at startup | Print error to stderr, exit non-zero |
| File disappears after startup | Return HTTP 500 with an error page on next request |

## Section 2: HTML Renderer (`internal/htmlrenderer`)

Implements `goldmark.NodeRenderer` via `RegisterFuncs()` — same interface pattern as `internal/renderer`.

### Node-to-HTML mapping

| Markdown element | HTML output |
|---|---|
| H1–H6 | `<h1>`–`<h6>` styled via CSS class |
| Paragraph | `<p>` |
| Fenced code block | `<div class="code-container lang-go"><span class="lang-badge">go</span><pre><code>…</code></pre></div>` |
| Indented code block | Same container, no badge |
| Container types (`tree`, `diagram`, `shell`) | Additional CSS class on container div for distinct border/badge color |
| Blockquote | `<blockquote>` with left-border CSS rule |
| Ordered list | `<ol><li>` |
| Unordered list | `<ul><li>` |
| Task list item | `<li class="task-item"><span class="checkbox checked|unchecked">` |
| Table | `<table><thead><tbody><tr><th><td>` |
| Emphasis | `<em>` |
| Strong | `<strong>` |
| Strikethrough | `<del>` |
| Inline code | `<code class="inline-code">` |
| Link | `<a href="…">` |
| Image | `<img src="…" alt="…">` |
| Thematic break | `<hr>` |

### CSS delivery

A single embedded CSS string (Go string constant or `embed.FS`) is injected into the page `<head>`. Colors are encoded as CSS custom properties (`--color-heading1`, `--color-code-bg`, etc.) set on `:root`. Dark and light themes set different values. No external fonts or CDN dependencies — the page works fully offline.

### Syntax highlighting

Chroma is called with the `html` formatter (instead of `terminal256` used by the terminal renderer), producing `<span class="chroma-…">` elements. The chroma CSS for the selected style (`github-dark` for dark theme, `xcode` for light) is also embedded in the page `<head>`.

## Section 3: HTTP Server and Live Reload (`internal/server`)

### Entry point

```go
func Start(inputs []Input, th theme.Theme, port int) error
```

`Input` is a struct holding either a file path (`Path string`) or pre-read bytes (`Data []byte`, used for stdin).

### Endpoints

**`GET /`**
- Re-reads all file-backed inputs from disk on every request (supports both live reload and manual browser refresh)
- Renders via `htmlrenderer`
- Returns a complete HTML page: fragment wrapped in `<html><head>…</head><body>…</body></html>`
- The `<head>` contains: embedded CSS, embedded chroma CSS, and a `<script>` block that opens `EventSource("/events")` and calls `location.reload()` on any message event

**`GET /events`**
- Sets `Content-Type: text/event-stream`, `Cache-Control: no-cache`, flushes headers
- Registers the `http.ResponseWriter` (as an `http.Flusher`) in the broadcaster
- Blocks until client disconnects via `r.Context().Done()`
- Unregisters on disconnect

### Broadcaster

A mutex-protected map of channels. When `watcher` fires a change event, `broadcaster.Notify()` writes `data: reload\n\n` to every registered channel and flushes each.

### Port selection

Default `7331`. Overridden by `--port`. If the port is already in use, `Start()` returns an error immediately with a descriptive message.

### Stdin mode

Input bytes captured once at startup. Watcher not started. SSE endpoint still exists but never fires. Manual browser refresh re-renders from the captured bytes.

### Shutdown

`Start()` runs `http.ListenAndServe` in a goroutine and blocks on `os.Signal` (`SIGINT`, `SIGTERM`). On signal, prints a newline to stderr and returns `nil`.

## Section 4: File Watcher (`internal/watcher`)

### Entry point

```go
func Watch(paths []string, onChange func()) (stop func(), err error)
```

- Wraps `fsnotify.NewWatcher`
- Watches all provided paths for `Write` and `Create` events
- Calls `onChange` (which triggers `broadcaster.Notify()`) on any event, with a short debounce (~100ms) to avoid multiple rapid reloads on a single save
- Returns a `stop` function that closes the watcher goroutine cleanly
- If `paths` is empty (stdin mode), returns a no-op stop function and `nil` error — no watcher is started

## Section 5: Multi-file Rendering

When multiple file inputs are provided, the HTML page renders them in sequence. Between each file a separator is inserted:

```html
<hr class="file-separator">
<p class="file-label">filename.md</p>
```

This mirrors the terminal mode `fileSeparator()` behaviour. Each file is re-read from disk on every `GET /` request.

## Section 6: Testing

### `internal/htmlrenderer`
- Unit tests per node type: assert HTML output contains expected tags, classes, and content
- Tests for both dark and light themes (verify CSS custom properties differ)
- Container type tests: `tree`, `diagram`, `shell` get correct CSS class
- Integration test using `testdata/full.md`: assert broad structural presence of `<table>`, `<blockquote>`, task checkboxes, code containers

### `internal/watcher`
- Write a temp file, start watcher, modify file, assert callback fires within timeout
- Assert no-op watcher (empty paths) never calls callback

### `internal/server` (using `httptest`)
- `GET /` returns 200, `Content-Type: text/html`, rendered content present
- `GET /` re-renders on each request: modify temp file between two fetches, assert content differs
- `GET /events` returns `Content-Type: text/event-stream`
- Multi-file: response contains both rendered files and separator markup

### `cmd/weave`
- `--web --port 0` starts without error (port 0 = OS-assigned free port)
- `--web` with a missing file exits non-zero

No end-to-end browser automation. SSE live reload verified at unit level via broadcaster tests.

## Dependencies

One new direct dependency:

| Package | Purpose |
|---|---|
| `github.com/fsnotify/fsnotify` | Cross-platform file system change notifications |

All other functionality uses the Go standard library (`net/http`, `os/signal`, `sync`, `embed`) and existing project dependencies (goldmark, chroma, lipgloss, cobra).

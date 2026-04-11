package htmlrenderer

import (
	"bytes"
	"strings"

	"github.com/alecthomas/chroma/v2"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/vchitepu/weave/internal/theme"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/util"
)

type containerType int

const (
	containerCode containerType = iota
	containerTree
	containerDiagram
	containerShell
)

func detectContainer(lang string) containerType {
	switch strings.ToLower(lang) {
	case "tree":
		return containerTree
	case "ascii", "diagram", "art", "mermaid":
		return containerDiagram
	case "bash", "sh", "shell", "console", "terminal":
		return containerShell
	default:
		return containerCode
	}
}

func (r *Renderer) renderFencedCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	n := node.(*ast.FencedCodeBlock)
	lang := ""
	if n.Language(source) != nil {
		lang = string(n.Language(source))
	}

	code := collectLines(node, source)
	writeCodeContainer(w, code, lang, r.theme)

	return ast.WalkSkipChildren, nil
}

func (r *Renderer) renderCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	code := collectLines(node, source)
	writeCodeContainer(w, code, "", r.theme)

	return ast.WalkSkipChildren, nil
}

// collectLines extracts the text content from a code block's lines.
func collectLines(node ast.Node, source []byte) string {
	var codeBuf bytes.Buffer
	lines := node.Lines()
	for i := 0; i < lines.Len(); i++ {
		line := lines.At(i)
		codeBuf.Write(line.Value(source))
	}
	return strings.TrimRight(codeBuf.String(), "\n")
}

// writeCodeContainer writes an HTML code container to the writer.
func writeCodeContainer(w util.BufWriter, code, lang string, th theme.Theme) {
	ct := detectContainer(lang)

	// Determine CSS class variant and badge label
	var variant string
	var badge string

	switch ct {
	case containerTree:
		variant = " container-tree"
		badge = "tree"
	case containerDiagram:
		variant = " container-diagram"
		badge = "diagram"
	case containerShell:
		variant = " container-shell"
		badge = "$"
	default:
		variant = ""
		badge = lang
	}

	_, _ = w.WriteString(`<div class="code-container` + variant + `">`)
	_, _ = w.WriteString("\n")

	if badge != "" {
		_, _ = w.WriteString(`<span class="lang-badge">` + htmlEscapeString(badge) + `</span>`)
		_, _ = w.WriteString("\n")
	}

	// Syntax highlight for code containers with a non-empty language
	var content string
	if ct == containerCode && lang != "" {
		content = highlightCodeHTML(code, lang, th)
	} else {
		content = htmlEscapeString(code)
	}

	_, _ = w.WriteString("<pre>" + content + "</pre>\n")
	_, _ = w.WriteString("</div>\n")
}

// highlightCodeHTML uses chroma to syntax-highlight code with inline styles.
func highlightCodeHTML(code, lang string, th theme.Theme) string {
	lexer := lexers.Get(lang)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	style := styles.Get(th.ChromaStyle)
	if style == nil {
		style = styles.Fallback
	}

	formatter := chromahtml.New(
		chromahtml.WithClasses(false),
		chromahtml.InlineCode(true),
		chromahtml.PreventSurroundingPre(true),
	)

	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		return htmlEscapeString(code)
	}

	var buf bytes.Buffer
	err = formatter.Format(&buf, style, iterator)
	if err != nil {
		return htmlEscapeString(code)
	}

	return buf.String()
}

// ChromaCSS returns chroma CSS class definitions for the given theme.
// Provided for future class-based highlighting use.
func ChromaCSS(th theme.Theme) string {
	style := styles.Get(th.ChromaStyle)
	if style == nil {
		style = styles.Fallback
	}
	formatter := chromahtml.New(chromahtml.WithClasses(true))
	var buf bytes.Buffer
	_ = formatter.WriteCSS(&buf, style)
	return buf.String()
}

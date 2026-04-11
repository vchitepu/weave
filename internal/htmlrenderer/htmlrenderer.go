package htmlrenderer

import (
	"bytes"

	"github.com/vchitepu/weave/internal/theme"
	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
	goldrenderer "github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// Priority is the goldmark renderer priority for the HTML renderer.
// Lower values = higher priority. We use 100 to take precedence over
// goldmark's default HTML renderer (which uses 1000).
const Priority = 100

// Renderer implements goldmark's NodeRenderer interface for HTML output.
type Renderer struct {
	theme theme.Theme
}

// New creates a new HTML Renderer.
func New(th theme.Theme) *Renderer {
	return &Renderer{theme: th}
}

// RegisterFuncs registers AST node render functions.
func (r *Renderer) RegisterFuncs(reg goldrenderer.NodeRendererFuncRegisterer) {
	// Block nodes
	reg.Register(ast.KindDocument, r.renderDocument)
	reg.Register(ast.KindParagraph, r.renderParagraph)
	reg.Register(ast.KindHeading, r.renderHeading)
	reg.Register(ast.KindThematicBreak, r.renderThematicBreak)
	reg.Register(ast.KindFencedCodeBlock, r.renderFencedCodeBlock)
	reg.Register(ast.KindCodeBlock, r.renderCodeBlock)
	reg.Register(ast.KindBlockquote, r.renderBlockquote)
	reg.Register(ast.KindList, r.renderList)
	reg.Register(ast.KindListItem, r.renderListItem)

	// Table nodes (from goldmark extension)
	reg.Register(east.KindTable, r.renderTable)
	reg.Register(east.KindTableHeader, r.renderTableHeader)
	reg.Register(east.KindTableRow, r.renderTableRow)
	reg.Register(east.KindTableCell, r.renderTableCell)

	// Strikethrough (from goldmark extension)
	reg.Register(east.KindStrikethrough, r.renderStrikethrough)

	// Task list checkboxes (from goldmark extension)
	reg.Register(east.KindTaskCheckBox, r.renderTaskCheckBox)

	// Inline nodes
	reg.Register(ast.KindText, r.renderText)
	reg.Register(ast.KindString, r.renderString)
	reg.Register(ast.KindEmphasis, r.renderEmphasis)
	reg.Register(ast.KindCodeSpan, r.renderCodeSpan)
	reg.Register(ast.KindLink, r.renderLink)
	reg.Register(ast.KindImage, r.renderImage)
}

func (r *Renderer) renderDocument(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *Renderer) renderParagraph(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<p>")
	} else {
		_, _ = w.WriteString("</p>\n")
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderText(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*ast.Text)
	_, _ = w.Write(htmlEscape(n.Segment.Value(source)))
	if n.HardLineBreak() {
		_, _ = w.WriteString("<br>\n")
	} else if n.SoftLineBreak() {
		_, _ = w.WriteString("\n")
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderString(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*ast.String)
	_, _ = w.Write(htmlEscape(n.Value))
	return ast.WalkContinue, nil
}

// htmlEscape escapes HTML special characters in a byte slice.
func htmlEscape(b []byte) []byte {
	var buf bytes.Buffer
	for _, c := range b {
		switch c {
		case '&':
			buf.WriteString("&amp;")
		case '<':
			buf.WriteString("&lt;")
		case '>':
			buf.WriteString("&gt;")
		case '"':
			buf.WriteString("&quot;")
		default:
			buf.WriteByte(c)
		}
	}
	return buf.Bytes()
}

// htmlEscapeString escapes HTML special characters in a string.
func htmlEscapeString(s string) string {
	return string(htmlEscape([]byte(s)))
}

// HtmlEscapeString is the exported version for use by the server package.
func HtmlEscapeString(s string) string {
	return htmlEscapeString(s)
}

// collectText recursively collects text content from an AST node tree.
func collectText(node ast.Node, source []byte) string {
	var buf bytes.Buffer
	collectTextInto(&buf, node, source)
	return buf.String()
}

func collectTextInto(buf *bytes.Buffer, node ast.Node, source []byte) {
	switch n := node.(type) {
	case *ast.Text:
		buf.Write(n.Segment.Value(source))
		if n.SoftLineBreak() || n.HardLineBreak() {
			buf.WriteString("\n")
		}
	case *ast.String:
		buf.Write(n.Value)
	default:
		for child := node.FirstChild(); child != nil; child = child.NextSibling() {
			collectTextInto(buf, child, source)
		}
	}
}

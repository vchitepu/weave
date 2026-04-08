package renderer

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/vchitepu/weave/internal/theme"
	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
	goldrenderer "github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// Priority is the goldmark renderer priority for the weave renderer.
// Lower values = higher priority. We use 100 to take precedence over
// goldmark's default HTML renderer (which uses 1000).
const Priority = 100

// rightMargin is the number of columns left empty on the right side of all
// block-level elements so they don't butt against the terminal edge.
const rightMargin = 2

// leftPad is the number of spaces prepended to every block-level text line
// to give breathing room from the left terminal edge.
const leftPad = 2

// pad is the string form of leftPad, used for prefixing lines.
const pad = "  "

// Renderer implements goldmark's NodeRenderer interface.
type Renderer struct {
	theme     theme.Theme
	width     int
	tableData *tableState
	// listPrefixWidths tracks current list-item marker widths for wrapped text.
	listPrefixWidths []int
	// tightListItemActive indicates we're buffering inline text for a tight list item.
	tightListItemActive bool
	// paraBuf accumulates inline content while inside a paragraph so the full
	// text can be word-wrapped as a unit before being written.
	paraBuf *strings.Builder
}

// New creates a new weave Renderer.
func New(th theme.Theme, width int) *Renderer {
	if width < 20 {
		width = 20
	}
	return &Renderer{
		theme: th,
		width: width,
	}
}

// contentWidth returns the usable text width after subtracting both margins.
func (r *Renderer) contentWidth() int {
	w := r.width - rightMargin - leftPad
	if w < 20 {
		w = 20
	}
	return w
}

// wrapText word-wraps s to contentWidth, preserving existing hard newlines,
// and prepends pad to every output line.
func (r *Renderer) wrapText(s string) string {
	cw := r.contentWidth()
	// lipgloss.NewStyle().Width() hard-wraps; we use it purely for wrapping.
	wrapped := lipgloss.NewStyle().Width(cw).Render(strings.TrimRight(s, "\n"))
	lines := strings.Split(wrapped, "\n")
	var out strings.Builder
	for _, line := range lines {
		out.WriteString(pad + line + "\n")
	}
	return out.String()
}

func (r *Renderer) writeWrappedListParagraph(w util.BufWriter, s string, prefixWidth int) {
	wrapWidth := r.contentWidth() - prefixWidth
	if wrapWidth < 10 {
		wrapWidth = 10
	}
	wrapped := lipgloss.NewStyle().Width(wrapWidth).Render(strings.TrimRight(s, "\n"))
	lines := strings.Split(wrapped, "\n")
	if len(lines) == 0 {
		return
	}
	_, _ = w.WriteString(lines[0])
	cont := strings.Repeat(" ", prefixWidth)
	for _, line := range lines[1:] {
		_, _ = w.WriteString("\n" + cont + line)
	}
}

// RegisterFuncs registers AST node render functions.
func (r *Renderer) RegisterFuncs(reg goldrenderer.NodeRendererFuncRegisterer) {
	// Block nodes
	reg.Register(ast.KindDocument, r.renderDocument)
	reg.Register(ast.KindParagraph, r.renderParagraph)
	reg.Register(ast.KindHeading, r.renderHeadingEntering)
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

// renderParagraph buffers inline content on entering and flushes it word-wrapped on exit.
func (r *Renderer) renderParagraph(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		r.paraBuf = &strings.Builder{}
		return ast.WalkContinue, nil
	}

	// Flush buffered paragraph text, word-wrapped and left-padded.
	if r.paraBuf != nil {
		if node.Parent() != nil && node.Parent().Kind() == ast.KindListItem && len(r.listPrefixWidths) > 0 {
			r.writeWrappedListParagraph(w, r.paraBuf.String(), r.listPrefixWidths[len(r.listPrefixWidths)-1])
		} else {
			_, _ = w.WriteString(r.wrapText(r.paraBuf.String()))
		}
		r.paraBuf = nil
	}

	// Inside a list item, single blank line; otherwise double.
	if node.Parent() != nil && node.Parent().Kind() == ast.KindListItem {
		_, _ = w.WriteString("\n")
	} else {
		_, _ = w.WriteString("\n")
	}
	return ast.WalkContinue, nil
}

// renderText writes to paraBuf when inside a paragraph, otherwise directly.
func (r *Renderer) renderText(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*ast.Text)
	text := string(n.Segment.Value(source))

	if r.paraBuf == nil && len(r.listPrefixWidths) > 0 {
		for p := node.Parent(); p != nil; p = p.Parent() {
			if p.Kind() == ast.KindListItem {
				r.paraBuf = &strings.Builder{}
				r.tightListItemActive = true
				break
			}
		}
	}

	if r.paraBuf != nil {
		r.paraBuf.WriteString(text)
		if n.SoftLineBreak() {
			r.paraBuf.WriteString(" ") // turn soft breaks into spaces for wrapping
		}
		if n.HardLineBreak() {
			r.paraBuf.WriteString("\n")
		}
	} else {
		_, _ = w.Write([]byte(text))
		if n.SoftLineBreak() {
			_, _ = w.WriteString("\n")
		}
		if n.HardLineBreak() {
			_, _ = w.WriteString("\n")
		}
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderString(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*ast.String)
	if r.paraBuf != nil {
		r.paraBuf.Write(n.Value)
	} else {
		_, _ = w.Write(n.Value)
	}
	return ast.WalkContinue, nil
}

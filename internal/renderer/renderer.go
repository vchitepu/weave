package renderer

import (
	"github.com/vinaychitepu/shine/internal/theme"
	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
	goldrenderer "github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// Priority is the goldmark renderer priority for the shine renderer.
// Lower values = higher priority. We use 100 to take precedence over
// goldmark's default HTML renderer (which uses 1000).
const Priority = 100

// rightMargin is the number of columns left empty on the right side of
// all block-level elements (rules, code containers) so they don't butt
// against the terminal edge.
const rightMargin = 2

// Renderer implements goldmark's NodeRenderer interface.
type Renderer struct {
	theme     theme.Theme
	width     int
	tableData *tableState
}

// New creates a new shine Renderer.
func New(th theme.Theme, width int) *Renderer {
	if width < 20 {
		width = 20
	}
	return &Renderer{
		theme: th,
		width: width,
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

// contentWidth returns the usable line width after subtracting the right margin.
func (r *Renderer) contentWidth() int {
	w := r.width - rightMargin
	if w < 20 {
		w = 20
	}
	return w
}

func (r *Renderer) renderDocument(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *Renderer) renderParagraph(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		// Inside a list item, use single newline (the list handles spacing)
		if node.Parent() != nil && node.Parent().Kind() == ast.KindListItem {
			_, _ = w.WriteString("\n")
		} else {
			_, _ = w.WriteString("\n\n")
		}
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderText(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*ast.Text)
	segment := n.Segment
	text := segment.Value(source)
	_, _ = w.Write(text)
	if n.SoftLineBreak() {
		_, _ = w.WriteString("\n")
	}
	if n.HardLineBreak() {
		_, _ = w.WriteString("\n")
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderString(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*ast.String)
	_, _ = w.Write(n.Value)
	return ast.WalkContinue, nil
}

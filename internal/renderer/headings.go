package renderer

import (
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/util"
)

func (r *Renderer) renderHeadingEntering(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	n := node.(*ast.Heading)

	text := collectText(node, source)

	var styled string
	switch {
	case n.Level == 1:
		styled = r.theme.H1.Render(text)
	case n.Level == 2:
		styled = r.theme.H2.Render(text)
	default:
		styled = r.theme.H3.Render(text)
	}

	_, _ = w.WriteString(pad + styled)
	_, _ = w.WriteString("\n")

	// H1 gets a full-width rule
	if n.Level == 1 {
		rule := r.theme.HeadingRule.Render(strings.Repeat("─", r.contentWidth()))
		_, _ = w.WriteString(pad + rule)
		_, _ = w.WriteString("\n")
	}

	_, _ = w.WriteString("\n")

	// Skip children — we already rendered the text
	return ast.WalkSkipChildren, nil
}

func (r *Renderer) renderThematicBreak(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	rule := r.theme.HorizontalRule.Render(strings.Repeat("─", r.contentWidth()))
	_, _ = w.WriteString(pad + rule)
	_, _ = w.WriteString("\n\n")
	return ast.WalkContinue, nil
}

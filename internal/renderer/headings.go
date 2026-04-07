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

	var textBuf strings.Builder
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		if t, ok := c.(*ast.Text); ok {
			textBuf.Write(t.Segment.Value(source))
		}
	}
	text := textBuf.String()

	var styled string
	switch {
	case n.Level == 1:
		styled = r.theme.H1.Render(text)
	case n.Level == 2:
		styled = r.theme.H2.Render(text)
	default:
		styled = r.theme.H3.Render(text)
	}

	_, _ = w.WriteString(styled)
	_, _ = w.WriteString("\n")

	// H1 gets a full-width rule
	if n.Level == 1 {
		rule := r.theme.HeadingRule.Render(strings.Repeat("─", r.width))
		_, _ = w.WriteString(rule)
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
	rule := r.theme.HorizontalRule.Render(strings.Repeat("─", r.width))
	_, _ = w.WriteString(rule)
	_, _ = w.WriteString("\n\n")
	return ast.WalkContinue, nil
}

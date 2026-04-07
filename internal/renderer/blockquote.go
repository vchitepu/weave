package renderer

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/util"
)

func (r *Renderer) renderBlockquote(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		bar := r.theme.BlockquoteBar.Render("│")
		_, _ = w.WriteString(bar + " ")
	} else {
		_, _ = w.WriteString("\n")
	}
	return ast.WalkContinue, nil
}

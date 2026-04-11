package htmlrenderer

import (
	"fmt"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/util"
)

func (r *Renderer) renderHeading(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Heading)
	level := n.Level
	if level < 1 {
		level = 1
	}
	if level > 6 {
		level = 6
	}
	tag := fmt.Sprintf("h%d", level)
	if entering {
		_, _ = w.WriteString("<" + tag + ">")
	} else {
		_, _ = w.WriteString("</" + tag + ">\n")
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderThematicBreak(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	_, _ = w.WriteString("<hr>\n")
	return ast.WalkContinue, nil
}

package renderer

import (
	"fmt"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/util"
)

func (r *Renderer) renderList(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		_, _ = w.WriteString("\n")
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderListItem(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	// Calculate indent depth
	depth := 0
	parent := node.Parent()
	for parent != nil {
		if parent.Kind() == ast.KindList {
			depth++
		}
		parent = parent.Parent()
	}
	indent := strings.Repeat("  ", depth-1)

	// Determine bullet or number
	list := node.Parent().(*ast.List)
	if list.IsOrdered() {
		pos := 1
		for sib := node.PreviousSibling(); sib != nil; sib = sib.PreviousSibling() {
			pos++
		}
		start := list.Start
		if start > 0 {
			pos = start + pos - 1
		}
		_, _ = w.WriteString(fmt.Sprintf("%s%d. ", indent, pos))
	} else {
		_, _ = w.WriteString(fmt.Sprintf("%s• ", indent))
	}

	return ast.WalkContinue, nil
}

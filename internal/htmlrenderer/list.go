package htmlrenderer

import (
	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/util"
)

func (r *Renderer) renderList(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.List)
	if n.IsOrdered() {
		if entering {
			_, _ = w.WriteString("<ol>\n")
		} else {
			_, _ = w.WriteString("</ol>\n")
		}
	} else {
		if entering {
			_, _ = w.WriteString("<ul>\n")
		} else {
			_, _ = w.WriteString("</ul>\n")
		}
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderListItem(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		// Check if this list item contains a TaskCheckBox
		isTask := false
		for child := node.FirstChild(); child != nil; child = child.NextSibling() {
			for inline := child.FirstChild(); inline != nil; inline = inline.NextSibling() {
				if inline.Kind() == east.KindTaskCheckBox {
					isTask = true
					break
				}
			}
			if isTask {
				break
			}
		}
		if isTask {
			_, _ = w.WriteString(`<li class="task-item">`)
		} else {
			_, _ = w.WriteString("<li>")
		}
	} else {
		_, _ = w.WriteString("</li>\n")
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderTaskCheckBox(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*east.TaskCheckBox)
	if n.IsChecked {
		_, _ = w.WriteString(`<span class="checkbox checked">&#x2713;</span> `)
	} else {
		_, _ = w.WriteString(`<span class="checkbox unchecked">&#x25CB;</span> `)
	}
	return ast.WalkContinue, nil
}

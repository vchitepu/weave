package htmlrenderer

import (
	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/util"
)

func (r *Renderer) renderTable(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<table>\n")
	} else {
		_, _ = w.WriteString("</tbody>\n</table>\n")
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderTableHeader(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<thead>\n<tr>\n")
	} else {
		_, _ = w.WriteString("</tr>\n</thead>\n<tbody>\n")
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderTableRow(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<tr>\n")
	} else {
		_, _ = w.WriteString("</tr>\n")
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderTableCell(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	// Determine if parent is a TableHeader
	isHeader := false
	if node.Parent() != nil && node.Parent().Kind() == east.KindTableHeader {
		isHeader = true
	}

	if isHeader {
		if entering {
			_, _ = w.WriteString("<th>")
		} else {
			_, _ = w.WriteString("</th>\n")
		}
	} else {
		if entering {
			_, _ = w.WriteString("<td>")
		} else {
			_, _ = w.WriteString("</td>\n")
		}
	}
	return ast.WalkContinue, nil
}

package htmlrenderer

import (
	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/util"
)

func (r *Renderer) renderEmphasis(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Emphasis)
	if n.Level == 2 {
		if entering {
			_, _ = w.WriteString("<strong>")
		} else {
			_, _ = w.WriteString("</strong>")
		}
	} else {
		if entering {
			_, _ = w.WriteString("<em>")
		} else {
			_, _ = w.WriteString("</em>")
		}
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderCodeSpan(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		text := collectText(node, source)
		_, _ = w.WriteString(`<code class="inline-code">`)
		_, _ = w.WriteString(htmlEscapeString(text))
		_, _ = w.WriteString("</code>")
		return ast.WalkSkipChildren, nil
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Link)
	if entering {
		_, _ = w.WriteString(`<a href="` + htmlEscapeString(string(n.Destination)) + `">`)
	} else {
		_, _ = w.WriteString("</a>")
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderImage(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*ast.Image)
	alt := string(n.Text(source))
	_, _ = w.WriteString(`<img src="` + htmlEscapeString(string(n.Destination)) + `" alt="` + htmlEscapeString(alt) + `">`)
	return ast.WalkSkipChildren, nil
}

func (r *Renderer) renderStrikethrough(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	_ = node.(*east.Strikethrough)
	if entering {
		_, _ = w.WriteString("<del>")
	} else {
		_, _ = w.WriteString("</del>")
	}
	return ast.WalkContinue, nil
}

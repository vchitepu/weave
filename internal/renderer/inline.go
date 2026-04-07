package renderer

import (
	"fmt"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/util"
)

// inlineWrite writes s to paraBuf if inside a paragraph, otherwise directly to w.
func (r *Renderer) inlineWrite(w util.BufWriter, s string) {
	if r.paraBuf != nil {
		r.paraBuf.WriteString(s)
	} else {
		_, _ = w.WriteString(s)
	}
}

func (r *Renderer) renderEmphasis(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Emphasis)
	if n.Level == 2 {
		if entering {
			r.inlineWrite(w, "\033[1m")
		} else {
			r.inlineWrite(w, "\033[22m")
		}
	} else {
		if entering {
			r.inlineWrite(w, "\033[3m")
		} else {
			r.inlineWrite(w, "\033[23m")
		}
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderCodeSpan(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		text := collectText(node, source)
		styled := r.theme.InlineCode.Render(" " + text + " ")
		r.inlineWrite(w, styled)
		return ast.WalkSkipChildren, nil
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Link)
	if entering {
		text := collectText(node, source)
		linkText := r.theme.LinkText.Render(text)
		linkURL := r.theme.LinkURL.Render(fmt.Sprintf("(%s)", string(n.Destination)))
		r.inlineWrite(w, linkText+" "+linkURL)
		return ast.WalkSkipChildren, nil
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderImage(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Image)
	if entering {
		alt := string(n.Text(source))
		styled := r.theme.ImageAlt.Render(fmt.Sprintf("[image: %s]", alt))
		r.inlineWrite(w, styled)
		return ast.WalkSkipChildren, nil
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderStrikethrough(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		r.inlineWrite(w, "\033[9m")
	} else {
		r.inlineWrite(w, "\033[29m")
	}
	return ast.WalkContinue, nil
}

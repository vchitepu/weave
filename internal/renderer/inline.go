package renderer

import (
	"fmt"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/util"
)

func (r *Renderer) renderEmphasis(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Emphasis)
	if n.Level == 2 {
		// Bold
		if entering {
			_, _ = w.WriteString("\033[1m")
		} else {
			_, _ = w.WriteString("\033[22m") // bold off
		}
	} else {
		// Italic
		if entering {
			_, _ = w.WriteString("\033[3m")
		} else {
			_, _ = w.WriteString("\033[23m")
		}
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderCodeSpan(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		// Collect all child text
		var text string
		for c := node.FirstChild(); c != nil; c = c.NextSibling() {
			if t, ok := c.(*ast.Text); ok {
				text += string(t.Segment.Value(source))
			}
		}
		styled := r.theme.InlineCode.Render(" " + text + " ")
		_, _ = w.WriteString(styled)
		return ast.WalkSkipChildren, nil
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Link)
	if entering {
		// Collect link text from children
		var text string
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			if t, ok := c.(*ast.Text); ok {
				text += string(t.Segment.Value(source))
			}
		}
		linkText := r.theme.LinkText.Render(text)
		linkURL := r.theme.LinkURL.Render(fmt.Sprintf("(%s)", string(n.Destination)))
		_, _ = w.WriteString(linkText + " " + linkURL)
		return ast.WalkSkipChildren, nil
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderImage(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Image)
	if entering {
		alt := string(n.Text(source))
		styled := r.theme.ImageAlt.Render(fmt.Sprintf("[image: %s]", alt))
		_, _ = w.WriteString(styled)
		return ast.WalkSkipChildren, nil
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderStrikethrough(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("\033[9m") // strikethrough on
	} else {
		_, _ = w.WriteString("\033[29m") // strikethrough off
	}
	return ast.WalkContinue, nil
}

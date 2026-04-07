package renderer

import (
	"bytes"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/util"
)

func (r *Renderer) renderBlockquote(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	// Collect all text content from children
	text := collectText(node, source)

	bar := r.theme.BlockquoteBar.Render("│")
	text = strings.TrimRight(text, "\n")
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		_, _ = w.WriteString(pad + bar + " " + line + "\n")
	}
	_, _ = w.WriteString("\n")

	return ast.WalkSkipChildren, nil
}

// collectText recursively collects text content from an AST node tree.
func collectText(node ast.Node, source []byte) string {
	var buf bytes.Buffer
	collectTextInto(&buf, node, source)
	return buf.String()
}

func collectTextInto(buf *bytes.Buffer, node ast.Node, source []byte) {
	switch n := node.(type) {
	case *ast.Text:
		buf.Write(n.Segment.Value(source))
		if n.SoftLineBreak() || n.HardLineBreak() {
			buf.WriteString("\n")
		}
	case *ast.String:
		buf.Write(n.Value)
	default:
		for child := node.FirstChild(); child != nil; child = child.NextSibling() {
			collectTextInto(buf, child, source)
		}
	}
}

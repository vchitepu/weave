package renderer

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
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
		if r.tightListItemActive && r.paraBuf != nil && len(r.listPrefixWidths) > 0 {
			r.writeWrappedListParagraph(w, r.paraBuf.String(), r.listPrefixWidths[len(r.listPrefixWidths)-1])
			r.paraBuf = nil
			r.tightListItemActive = false
		}

		// Ensure each list item ends with a newline.
		// For loose lists (with paragraphs), renderParagraph handles the newline.
		// For tight lists (no paragraphs), we need to add it here.
		// Check if the last child is NOT a paragraph — if so, add a newline.
		lastChild := node.LastChild()
		if lastChild == nil || lastChild.Kind() != ast.KindParagraph {
			_, _ = w.WriteString("\n")
		}
		if len(r.listPrefixWidths) > 0 {
			r.listPrefixWidths = r.listPrefixWidths[:len(r.listPrefixWidths)-1]
		}
		return ast.WalkContinue, nil
	}

	// If a tight list item is currently buffering text, flush it before
	// starting a nested or sibling list item so prefixes don't concatenate.
	if r.tightListItemActive && r.paraBuf != nil && len(r.listPrefixWidths) > 0 {
		r.writeWrappedListParagraph(w, r.paraBuf.String(), r.listPrefixWidths[len(r.listPrefixWidths)-1])
		_, _ = w.WriteString("\n")
		r.paraBuf = nil
		r.tightListItemActive = false
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

	// Check if this list item has a task checkbox — if so, suppress the bullet
	// and write only the indent prefix. The checkbox symbol is rendered by
	// renderTaskCheckBox.
	isTask := false
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		// TaskCheckBox is an inline node inside a TextBlock/Paragraph child
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
		prefix := fmt.Sprintf("%s%s", pad, indent)
		// +2 accounts for the checkbox symbol + space ("✓ " or "○ ") written by
		// renderTaskCheckBox, so continuation lines of wrapped text align correctly.
		r.listPrefixWidths = append(r.listPrefixWidths, lipgloss.Width(prefix)+2)
		_, _ = w.WriteString(prefix)
	} else if list.IsOrdered() {
		pos := 1
		for sib := node.PreviousSibling(); sib != nil; sib = sib.PreviousSibling() {
			pos++
		}
		start := list.Start
		if start > 0 {
			pos = start + pos - 1
		}
		prefix := fmt.Sprintf("%s%s%d. ", pad, indent, pos)
		r.listPrefixWidths = append(r.listPrefixWidths, lipgloss.Width(prefix))
		_, _ = w.WriteString(prefix)
	} else {
		prefix := fmt.Sprintf("%s%s• ", pad, indent)
		r.listPrefixWidths = append(r.listPrefixWidths, lipgloss.Width(prefix))
		_, _ = w.WriteString(prefix)
	}

	if node.FirstChild() != nil && node.FirstChild().Kind() == ast.KindText {
		r.paraBuf = &strings.Builder{}
		r.tightListItemActive = true
	}

	return ast.WalkContinue, nil
}

func (r *Renderer) renderTaskCheckBox(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*east.TaskCheckBox)
	if n.IsChecked {
		_, _ = w.WriteString(r.theme.TaskChecked.Render("✓") + " ")
	} else {
		_, _ = w.WriteString(r.theme.TaskUnchecked.Render("○") + " ")
	}
	return ast.WalkContinue, nil
}

package renderer

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/util"
)

// tableState accumulates rows during table rendering.
type tableState struct {
	headers []string
	rows    [][]string
	current []string
	inHead  bool
}

func (r *Renderer) renderTable(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		r.tableData = &tableState{}
		return ast.WalkContinue, nil
	}

	// Exiting table: render the accumulated data
	td := r.tableData
	if td == nil {
		return ast.WalkContinue, nil
	}

	// Calculate column widths
	numCols := len(td.headers)
	colWidths := make([]int, numCols)
	for i, h := range td.headers {
		if len(h) > colWidths[i] {
			colWidths[i] = len(h)
		}
	}
	for _, row := range td.rows {
		for i, cell := range row {
			if i < numCols && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	borderColor := r.theme.TableBorder
	hSep := r.tableHSep(colWidths, borderColor, "├", "┼", "┤")
	topBorder := r.tableHSep(colWidths, borderColor, "┌", "┬", "┐")
	bottomBorder := r.tableHSep(colWidths, borderColor, "└", "┴", "┘")

	bStyle := lipgloss.NewStyle().Foreground(borderColor)

	var buf strings.Builder

	// Top border
	buf.WriteString(topBorder + "\n")

	// Header row
	buf.WriteString(bStyle.Render("│"))
	for i, h := range td.headers {
		padded := r.padCell(h, colWidths[i])
		styled := r.theme.TableHeader.Render(padded)
		buf.WriteString(" " + styled + " " + bStyle.Render("│"))
	}
	buf.WriteString("\n")

	// Header separator
	buf.WriteString(hSep + "\n")

	// Data rows
	for _, row := range td.rows {
		buf.WriteString(bStyle.Render("│"))
		for i := 0; i < numCols; i++ {
			cell := ""
			if i < len(row) {
				cell = row[i]
			}
			padded := r.padCell(cell, colWidths[i])
			buf.WriteString(" " + padded + " " + bStyle.Render("│"))
		}
		buf.WriteString("\n")
	}

	// Bottom border
	buf.WriteString(bottomBorder + "\n")

	// Prefix every line with leftPad spaces.
	for _, line := range strings.Split(strings.TrimRight(buf.String(), "\n"), "\n") {
		_, _ = w.WriteString(pad + line + "\n")
	}
	_, _ = w.WriteString("\n")

	r.tableData = nil
	return ast.WalkContinue, nil
}

func (r *Renderer) tableHSep(colWidths []int, borderColor lipgloss.Color, left, mid, right string) string {
	bStyle := lipgloss.NewStyle().Foreground(borderColor)
	var parts []string
	for _, w := range colWidths {
		parts = append(parts, strings.Repeat("─", w+2))
	}
	return bStyle.Render(left + strings.Join(parts, mid) + right)
}

func (r *Renderer) padCell(text string, width int) string {
	padding := width - len(text)
	if padding < 0 {
		padding = 0
	}
	return text + strings.Repeat(" ", padding)
}

func (r *Renderer) renderTableHeader(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if r.tableData == nil {
		return ast.WalkContinue, nil
	}
	if entering {
		r.tableData.inHead = true
		r.tableData.current = nil
	} else {
		r.tableData.headers = r.tableData.current
		r.tableData.inHead = false
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderTableRow(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		if r.tableData != nil {
			r.tableData.current = nil
		}
	} else {
		if r.tableData != nil {
			if r.tableData.inHead {
				r.tableData.headers = r.tableData.current
			} else {
				r.tableData.rows = append(r.tableData.rows, r.tableData.current)
			}
		}
	}
	return ast.WalkContinue, nil
}

func (r *Renderer) renderTableCell(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		text := collectText(node, source)
		if r.tableData != nil {
			r.tableData.current = append(r.tableData.current, strings.TrimSpace(text))
		}
		return ast.WalkSkipChildren, nil
	}
	return ast.WalkContinue, nil
}

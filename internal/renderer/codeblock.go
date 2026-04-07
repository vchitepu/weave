package renderer

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/charmbracelet/lipgloss"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/util"
)

type containerType int

const (
	containerCode containerType = iota
	containerTree
	containerDiagram
	containerShell
)

func detectContainer(lang string) containerType {
	switch strings.ToLower(lang) {
	case "tree":
		return containerTree
	case "ascii", "diagram", "art", "mermaid":
		return containerDiagram
	case "bash", "sh", "shell", "console", "terminal":
		return containerShell
	default:
		return containerCode
	}
}

func (r *Renderer) renderFencedCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	n := node.(*ast.FencedCodeBlock)
	lang := ""
	if n.Language(source) != nil {
		lang = string(n.Language(source))
	}

	code := collectLines(node, source)
	rendered := r.renderCodeContainer(w, code, lang)
	_, _ = w.WriteString(rendered)

	return ast.WalkSkipChildren, nil
}

func (r *Renderer) renderCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	code := collectLines(node, source)
	rendered := r.renderCodeContainer(w, code, "")
	_, _ = w.WriteString(rendered)

	return ast.WalkSkipChildren, nil
}

// collectLines extracts the text content from a code block's lines.
func collectLines(node ast.Node, source []byte) string {
	var codeBuf bytes.Buffer
	lines := node.Lines()
	for i := 0; i < lines.Len(); i++ {
		line := lines.At(i)
		codeBuf.Write(line.Value(source))
	}
	return strings.TrimRight(codeBuf.String(), "\n")
}

// renderCodeContainer builds a bordered lipgloss container for code content.
func (r *Renderer) renderCodeContainer(w util.BufWriter, code, lang string) string {
	ct := detectContainer(lang)
	isMermaid := strings.ToLower(lang) == "mermaid"

	// Determine border color and header label
	var borderColor lipgloss.Color
	var headerLabel string

	switch ct {
	case containerTree:
		borderColor = r.theme.TreeBorder
		headerLabel = "tree"
	case containerDiagram:
		borderColor = r.theme.DiagramBorder
		headerLabel = "diagram"
	case containerShell:
		borderColor = r.theme.ShellBorder
		headerLabel = "$"
	default:
		borderColor = r.theme.CodeBorder
		headerLabel = lang
	}

	// Syntax highlight the code (only for code containers with a language)
	highlighted := code
	if ct == containerCode && lang != "" {
		highlighted = r.highlightCode(code, lang)
	}

	// For mermaid, prepend a "diagram not rendered" note
	if isMermaid {
		highlighted = "[diagram not rendered]\n\n" + code
	}

	// Build the container using lipgloss
	// Subtract right margin so the box doesn't run to the terminal edge
	const rightMargin = 2
	boxWidth := r.width - rightMargin
	if boxWidth < 22 {
		boxWidth = 22
	}
	contentWidth := boxWidth - 4 // account for border + padding
	if contentWidth < 20 {
		contentWidth = 20
	}

	// Pad each line to fill the container width
	codeLines := strings.Split(highlighted, "\n")
	var paddedLines []string
	for _, line := range codeLines {
		visLen := lipgloss.Width(line)
		padding := contentWidth - visLen
		if padding < 0 {
			padding = 0
		}
		paddedLines = append(paddedLines, line+strings.Repeat(" ", padding))
	}
	paddedContent := strings.Join(paddedLines, "\n")

	// Build the full box with lipgloss (includes its own top border).
	border := lipgloss.RoundedBorder()
	box := lipgloss.NewStyle().
		BorderStyle(border).
		BorderForeground(borderColor).
		Padding(0, 1).
		Width(boxWidth).
		Render(paddedContent)

	// Replace the top border line with a hand-crafted one that embeds the
	// label. We can't slice the lipgloss-rendered line by rune index because
	// it contains ANSI escape codes; instead we drop lipgloss's top line and
	// build our own plain-text top border, then color it.
	bStyle := lipgloss.NewStyle().Foreground(borderColor)
	// lipgloss Width() sets content width; the outer border adds 2 more columns.
	// So the visual outer width of the box is boxWidth+2, meaning our hand-crafted
	// top border needs boxWidth+2 visible rune columns total.
	outerWidth := boxWidth + 2
	innerWidth := outerWidth - 2 // subtract the two corner runes
	var topLine string
	if headerLabel != "" {
		label := fmt.Sprintf(" %s ", headerLabel)
		labelWidth := len([]rune(label))
		dashCount := innerWidth - labelWidth - 1 // -1 for the leading "─"
		if dashCount < 0 {
			dashCount = 0
		}
		topLine = bStyle.Render("╭─" + label + strings.Repeat("─", dashCount) + "╮")
	} else {
		topLine = bStyle.Render("╭" + strings.Repeat("─", innerWidth) + "╮")
	}

	// Drop lipgloss's first line (its top border) and prepend ours.
	boxLines := strings.SplitN(box, "\n", 2)
	rest := ""
	if len(boxLines) > 1 {
		rest = boxLines[1]
	}
	return topLine + "\n" + rest + "\n\n"
}

// highlightCode applies chroma syntax highlighting to the given code.
func (r *Renderer) highlightCode(code, lang string) string {
	lexer := lexers.Get(lang)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	style := styles.Get(r.theme.ChromaStyle)
	if style == nil {
		style = styles.Fallback
	}

	formatter := formatters.Get("terminal256")
	if formatter == nil {
		formatter = formatters.Fallback
	}

	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		return code
	}

	var buf bytes.Buffer
	err = formatter.Format(&buf, style, iterator)
	if err != nil {
		return code
	}

	return strings.TrimRight(buf.String(), "\n\r")
}

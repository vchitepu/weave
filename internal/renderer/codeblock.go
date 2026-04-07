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
	contentWidth := r.width - 4 // account for border + padding
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

	// Build the box with lipgloss
	border := lipgloss.RoundedBorder()
	box := lipgloss.NewStyle().
		BorderStyle(border).
		BorderForeground(borderColor).
		Padding(0, 1).
		Width(r.width).
		Render(paddedContent)

	// Replace the top border line to include the header label
	boxLines := strings.Split(box, "\n")
	if len(boxLines) > 0 && headerLabel != "" {
		topBorder := boxLines[0]
		labelStr := fmt.Sprintf(" %s ", headerLabel)
		runes := []rune(topBorder)
		if len(runes) > 3 {
			labelRunes := []rune(labelStr)
			insertEnd := 2 + len(labelRunes)
			if insertEnd < len(runes) {
				newTop := make([]rune, 0, len(runes))
				newTop = append(newTop, runes[:2]...)
				newTop = append(newTop, labelRunes...)
				newTop = append(newTop, runes[insertEnd:]...)
				boxLines[0] = string(newTop)
			}
		}
	}

	return strings.Join(boxLines, "\n") + "\n\n"
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

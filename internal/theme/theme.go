package theme

import "github.com/charmbracelet/lipgloss"

// Theme holds all styles used by the renderer.
type Theme struct {
	// Text
	Normal        lipgloss.Style
	Bold          lipgloss.Style
	Italic        lipgloss.Style
	Strikethrough lipgloss.Style
	Dim           lipgloss.Style

	// Headings
	H1          lipgloss.Style
	H2          lipgloss.Style
	H3          lipgloss.Style
	HeadingRule lipgloss.Style

	// Code
	CodeBorder  lipgloss.Color
	CodeHeader  lipgloss.Style
	InlineCode  lipgloss.Style
	ChromaStyle string // chroma style name

	// Container variants
	TreeBorder    lipgloss.Color
	DiagramBorder lipgloss.Color
	ShellBorder   lipgloss.Color

	// Blockquote
	BlockquoteBar  lipgloss.Style
	BlockquoteText lipgloss.Style

	// Table
	TableHeader lipgloss.Style
	TableBorder lipgloss.Color

	// Links / Images
	LinkText lipgloss.Style
	LinkURL  lipgloss.Style
	ImageAlt lipgloss.Style

	// Horizontal rule
	HorizontalRule lipgloss.Style
}

// DarkTheme returns the built-in dark theme.
func DarkTheme() Theme {
	return Theme{
		Normal:        lipgloss.NewStyle(),
		Bold:          lipgloss.NewStyle().Bold(true),
		Italic:        lipgloss.NewStyle().Italic(true),
		Strikethrough: lipgloss.NewStyle().Strikethrough(true),
		Dim:           lipgloss.NewStyle().Faint(true),

		H1:          lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#61AFEF")),
		H2:          lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#C678DD")),
		H3:          lipgloss.NewStyle().Foreground(lipgloss.Color("#E5C07B")),
		HeadingRule: lipgloss.NewStyle().Foreground(lipgloss.Color("#3E4451")),

		CodeBorder:  lipgloss.Color("#3E4451"),
		CodeHeader:  lipgloss.NewStyle().Foreground(lipgloss.Color("#ABB2BF")).Faint(true),
		InlineCode:  lipgloss.NewStyle().Background(lipgloss.Color("#2C313C")),
		ChromaStyle: "dracula",

		TreeBorder:    lipgloss.Color("#98C379"),
		DiagramBorder: lipgloss.Color("#61AFEF"),
		ShellBorder:   lipgloss.Color("#E5C07B"),

		BlockquoteBar:  lipgloss.NewStyle().Foreground(lipgloss.Color("#3E4451")),
		BlockquoteText: lipgloss.NewStyle().Foreground(lipgloss.Color("#5C6370")),

		TableHeader: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#61AFEF")),
		TableBorder: lipgloss.Color("#3E4451"),

		LinkText: lipgloss.NewStyle().Foreground(lipgloss.Color("#61AFEF")),
		LinkURL:  lipgloss.NewStyle().Faint(true),
		ImageAlt: lipgloss.NewStyle().Italic(true).Faint(true),

		HorizontalRule: lipgloss.NewStyle().Foreground(lipgloss.Color("#3E4451")),
	}
}

// LightTheme returns the built-in light theme.
func LightTheme() Theme {
	return Theme{
		Normal:        lipgloss.NewStyle(),
		Bold:          lipgloss.NewStyle().Bold(true),
		Italic:        lipgloss.NewStyle().Italic(true),
		Strikethrough: lipgloss.NewStyle().Strikethrough(true),
		Dim:           lipgloss.NewStyle().Faint(true),

		H1:          lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#0366D6")),
		H2:          lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#6F42C1")),
		H3:          lipgloss.NewStyle().Foreground(lipgloss.Color("#B08800")),
		HeadingRule: lipgloss.NewStyle().Foreground(lipgloss.Color("#D0D7DE")),

		CodeBorder:  lipgloss.Color("#D0D7DE"),
		CodeHeader:  lipgloss.NewStyle().Foreground(lipgloss.Color("#57606A")).Faint(true),
		InlineCode:  lipgloss.NewStyle().Background(lipgloss.Color("#F6F8FA")),
		ChromaStyle: "github",

		TreeBorder:    lipgloss.Color("#1A7F37"),
		DiagramBorder: lipgloss.Color("#0366D6"),
		ShellBorder:   lipgloss.Color("#B08800"),

		BlockquoteBar:  lipgloss.NewStyle().Foreground(lipgloss.Color("#D0D7DE")),
		BlockquoteText: lipgloss.NewStyle().Foreground(lipgloss.Color("#57606A")),

		TableHeader: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#0366D6")),
		TableBorder: lipgloss.Color("#D0D7DE"),

		LinkText: lipgloss.NewStyle().Foreground(lipgloss.Color("#0366D6")),
		LinkURL:  lipgloss.NewStyle().Faint(true),
		ImageAlt: lipgloss.NewStyle().Italic(true).Faint(true),

		HorizontalRule: lipgloss.NewStyle().Foreground(lipgloss.Color("#D0D7DE")),
	}
}

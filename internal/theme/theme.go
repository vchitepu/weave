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

		H1:          lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8BA4D4")),
		H2:          lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#A8A0D6")),
		H3:          lipgloss.NewStyle().Foreground(lipgloss.Color("#C9A86A")),
		HeadingRule: lipgloss.NewStyle().Foreground(lipgloss.Color("#3A3F4B")),

		CodeBorder:  lipgloss.Color("#3A3F4B"),
		CodeHeader:  lipgloss.NewStyle().Foreground(lipgloss.Color("#B7C0D0")).Faint(true),
		InlineCode:  lipgloss.NewStyle().Background(lipgloss.Color("#262B36")),
		ChromaStyle: "dracula",

		TreeBorder:    lipgloss.Color("#5D7290"),
		DiagramBorder: lipgloss.Color("#7B74A6"),
		ShellBorder:   lipgloss.Color("#9F8656"),

		BlockquoteBar:  lipgloss.NewStyle().Foreground(lipgloss.Color("#3A3F4B")),
		BlockquoteText: lipgloss.NewStyle().Foreground(lipgloss.Color("#9BA3B2")),

		TableHeader: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8BA4D4")),
		TableBorder: lipgloss.Color("#3A3F4B"),

		LinkText: lipgloss.NewStyle().Foreground(lipgloss.Color("#7FA3C8")),
		LinkURL:  lipgloss.NewStyle().Faint(true),
		ImageAlt: lipgloss.NewStyle().Italic(true).Faint(true),

		HorizontalRule: lipgloss.NewStyle().Foreground(lipgloss.Color("#3A3F4B")),
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

		H1:          lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#3F5F8A")),
		H2:          lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#665E95")),
		H3:          lipgloss.NewStyle().Foreground(lipgloss.Color("#8D6B3F")),
		HeadingRule: lipgloss.NewStyle().Foreground(lipgloss.Color("#C2C7D0")),

		CodeBorder:  lipgloss.Color("#C2C7D0"),
		CodeHeader:  lipgloss.NewStyle().Foreground(lipgloss.Color("#5F6B7A")).Faint(true),
		InlineCode:  lipgloss.NewStyle().Background(lipgloss.Color("#EEF1F5")),
		ChromaStyle: "github",

		TreeBorder:    lipgloss.Color("#5E7699"),
		DiagramBorder: lipgloss.Color("#7D73A3"),
		ShellBorder:   lipgloss.Color("#9A7B52"),

		BlockquoteBar:  lipgloss.NewStyle().Foreground(lipgloss.Color("#C2C7D0")),
		BlockquoteText: lipgloss.NewStyle().Foreground(lipgloss.Color("#6C7483")),

		TableHeader: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#3F5F8A")),
		TableBorder: lipgloss.Color("#C2C7D0"),

		LinkText: lipgloss.NewStyle().Foreground(lipgloss.Color("#496A92")),
		LinkURL:  lipgloss.NewStyle().Faint(true),
		ImageAlt: lipgloss.NewStyle().Italic(true).Faint(true),

		HorizontalRule: lipgloss.NewStyle().Foreground(lipgloss.Color("#C2C7D0")),
	}
}

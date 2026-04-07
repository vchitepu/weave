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

		H1:          lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#9CCFD8")),
		H2:          lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#C4A7E7")),
		H3:          lipgloss.NewStyle().Foreground(lipgloss.Color("#F6C177")),
		HeadingRule: lipgloss.NewStyle().Foreground(lipgloss.Color("#393552")),

		CodeBorder:  lipgloss.Color("#393552"),
		CodeHeader:  lipgloss.NewStyle().Foreground(lipgloss.Color("#E0DEF4")).Faint(true),
		InlineCode:  lipgloss.NewStyle().Background(lipgloss.Color("#2A273F")),
		ChromaStyle: "dracula",

		TreeBorder:    lipgloss.Color("#9CCFD8"),
		DiagramBorder: lipgloss.Color("#C4A7E7"),
		ShellBorder:   lipgloss.Color("#F6C177"),

		BlockquoteBar:  lipgloss.NewStyle().Foreground(lipgloss.Color("#393552")),
		BlockquoteText: lipgloss.NewStyle().Foreground(lipgloss.Color("#908CAA")),

		TableHeader: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#9CCFD8")),
		TableBorder: lipgloss.Color("#393552"),

		LinkText: lipgloss.NewStyle().Foreground(lipgloss.Color("#9CCFD8")),
		LinkURL:  lipgloss.NewStyle().Faint(true),
		ImageAlt: lipgloss.NewStyle().Italic(true).Faint(true),

		HorizontalRule: lipgloss.NewStyle().Foreground(lipgloss.Color("#393552")),
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

		H1:          lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#286983")),
		H2:          lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#907AA9")),
		H3:          lipgloss.NewStyle().Foreground(lipgloss.Color("#EA9D34")),
		HeadingRule: lipgloss.NewStyle().Foreground(lipgloss.Color("#9893A5")),

		CodeBorder:  lipgloss.Color("#9893A5"),
		CodeHeader:  lipgloss.NewStyle().Foreground(lipgloss.Color("#575279")).Faint(true),
		InlineCode:  lipgloss.NewStyle().Background(lipgloss.Color("#F2E9E1")),
		ChromaStyle: "github",

		TreeBorder:    lipgloss.Color("#286983"),
		DiagramBorder: lipgloss.Color("#907AA9"),
		ShellBorder:   lipgloss.Color("#EA9D34"),

		BlockquoteBar:  lipgloss.NewStyle().Foreground(lipgloss.Color("#9893A5")),
		BlockquoteText: lipgloss.NewStyle().Foreground(lipgloss.Color("#797593")),

		TableHeader: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#286983")),
		TableBorder: lipgloss.Color("#9893A5"),

		LinkText: lipgloss.NewStyle().Foreground(lipgloss.Color("#286983")),
		LinkURL:  lipgloss.NewStyle().Faint(true),
		ImageAlt: lipgloss.NewStyle().Italic(true).Faint(true),

		HorizontalRule: lipgloss.NewStyle().Foreground(lipgloss.Color("#9893A5")),
	}
}

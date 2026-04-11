package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vchitepu/weave/internal/pager"
	"github.com/vchitepu/weave/internal/renderer"
	"github.com/vchitepu/weave/internal/theme"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	goldrenderer "github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
	"golang.org/x/term"
)

var (
	version   = "dev"
	themeFlag string
	widthFlag int
)

const (
	separatorLeftPad     = 2
	separatorRightMargin = 2
	separatorPad         = "  "
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "weave [file...]",
		Short:   "A terminal Markdown viewer with rich visual containers",
		Version: version,
		Args:    cobra.ArbitraryArgs,
		RunE:    run,
	}

	rootCmd.Flags().StringVar(&themeFlag, "theme", "", "Override theme (dark|light)")
	rootCmd.Flags().IntVar(&widthFlag, "width", 0, "Override terminal width")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	// Multi-file mode: render each file with separators.
	if len(args) >= 2 {
		// Detect terminal width
		width := widthFlag
		autoWidth := widthFlag == 0
		if width == 0 {
			if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
				width = w
			} else {
				width = 80
			}
		}
		width = normalizeWidth(width, autoWidth)

		// Validate theme flag
		if themeFlag != "" && themeFlag != "dark" && themeFlag != "light" {
			return fmt.Errorf("weave: invalid theme %q (use 'dark' or 'light')", themeFlag)
		}
		th := theme.Detect(themeFlag)
		md := buildMarkdown(th, width)

		var combined strings.Builder
		for i, path := range args {
			if i > 0 {
				combined.WriteString(fileSeparator(path, width, th))
			}
			rendered, err := renderFile(path, md)
			if err != nil {
				return err
			}
			combined.WriteString(rendered)
		}

		output := combined.String()

		// Determine if we should page
		isTTY := term.IsTerminal(int(os.Stdout.Fd()))
		if isTTY {
			_, termHeight, err := term.GetSize(int(os.Stdout.Fd()))
			if err != nil {
				termHeight = 24
			}
			lineCount := strings.Count(output, "\n")
			if pager.ShouldPage(lineCount, termHeight) {
				return pager.Run(output)
			}
		}

		_, err := fmt.Fprint(os.Stdout, output)
		return err
	}

	// Read input
	var input []byte
	var err error

	if len(args) == 1 {
		input, err = os.ReadFile(args[0])
		if err != nil {
			return fmt.Errorf("weave: no such file: %s", args[0])
		}
	} else {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			return cmd.Help()
		}
		input, err = io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("weave: failed to read stdin: %w", err)
		}
	}

	// Detect terminal width
	width := widthFlag
	autoWidth := widthFlag == 0
	if width == 0 {
		if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
			width = w
		} else {
			width = 80
		}
	}
	width = normalizeWidth(width, autoWidth)

	// Validate theme flag
	if themeFlag != "" && themeFlag != "dark" && themeFlag != "light" {
		return fmt.Errorf("weave: invalid theme %q (use 'dark' or 'light')", themeFlag)
	}

	// Detect theme
	th := theme.Detect(themeFlag)

	md := buildMarkdown(th, width)

	// Render
	var buf bytes.Buffer
	if err := md.Convert(input, &buf); err != nil {
		return fmt.Errorf("weave: render error: %w", err)
	}

	output := buf.String()

	// Determine if we should page
	isTTY := term.IsTerminal(int(os.Stdout.Fd()))
	if isTTY {
		_, termHeight, err := term.GetSize(int(os.Stdout.Fd()))
		if err != nil {
			termHeight = 24
		}
		lineCount := strings.Count(output, "\n")
		if pager.ShouldPage(lineCount, termHeight) {
			return pager.Run(output)
		}
	}

	// Write directly to stdout
	_, err = fmt.Fprint(os.Stdout, output)
	return err
}

func normalizeWidth(width int, auto bool) int {
	if width < 20 {
		return 20
	}
	if auto && width > 120 {
		return 120
	}
	return width
}

func fileSeparator(filename string, width int, th theme.Theme) string {
	contentWidth := width - separatorRightMargin - separatorLeftPad
	if contentWidth < 1 {
		contentWidth = 1
	}

	rule := th.HorizontalRule.Render(strings.Repeat("─", contentWidth))
	label := th.Dim.Render(filename)

	return "\n" + separatorPad + rule + "\n" + separatorPad + label + "\n\n"
}

func buildMarkdown(th theme.Theme, width int) goldmark.Markdown {
	r := renderer.New(th, width)

	return goldmark.New(
		goldmark.WithExtensions(extension.Table, extension.Strikethrough, extension.TaskList),
		goldmark.WithRenderer(
			goldrenderer.NewRenderer(
				goldrenderer.WithNodeRenderers(
					util.Prioritized(r, renderer.Priority),
				),
			),
		),
	)
}

func renderFile(path string, md goldmark.Markdown) (string, error) {
	input, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("weave: no such file: %s", path)
	}

	var buf bytes.Buffer
	if err := md.Convert(input, &buf); err != nil {
		return "", fmt.Errorf("weave: render error for %s: %w", path, err)
	}

	return buf.String(), nil
}

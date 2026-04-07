package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vinaychitepu/shine/internal/pager"
	"github.com/vinaychitepu/shine/internal/renderer"
	"github.com/vinaychitepu/shine/internal/theme"
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

func main() {
	rootCmd := &cobra.Command{
		Use:     "shine [file]",
		Short:   "A terminal Markdown viewer with rich visual containers",
		Version: version,
		Args:    cobra.MaximumNArgs(1),
		RunE:    run,
	}

	rootCmd.Flags().StringVar(&themeFlag, "theme", "", "Override theme (dark|light)")
	rootCmd.Flags().IntVar(&widthFlag, "width", 0, "Override terminal width")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	// Read input
	var input []byte
	var err error

	if len(args) == 1 {
		input, err = os.ReadFile(args[0])
		if err != nil {
			return fmt.Errorf("shine: no such file: %s", args[0])
		}
	} else {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			return cmd.Help()
		}
		input, err = io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("shine: failed to read stdin: %w", err)
		}
	}

	// Detect terminal width
	width := widthFlag
	if width == 0 {
		if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
			width = w
		} else {
			width = 80
		}
	}

	// Detect theme
	th := theme.Detect(themeFlag)

	// Build goldmark with our renderer
	r := renderer.New(th, width)
	md := goldmark.New(
		goldmark.WithExtensions(extension.Table, extension.Strikethrough),
		goldmark.WithRenderer(
			goldrenderer.NewRenderer(
				goldrenderer.WithNodeRenderers(
					util.Prioritized(r, 100),
				),
			),
		),
	)

	// Render
	var buf bytes.Buffer
	if err := md.Convert(input, &buf); err != nil {
		return fmt.Errorf("shine: render error: %w", err)
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

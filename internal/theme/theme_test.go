package theme

import (
	"fmt"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func colorString(c lipgloss.TerminalColor) string {
	if c == nil {
		return ""
	}
	return fmt.Sprint(c)
}

func TestDarkThemeNotNil(t *testing.T) {
	th := DarkTheme()
	if !th.H1.GetBold() {
		t.Fatal("DarkTheme H1 style should be bold")
	}
	if th.H1.GetForeground() == nil {
		t.Fatal("DarkTheme H1 style should have a foreground color")
	}
	if th.ChromaStyle != "github-dark" {
		t.Fatalf("DarkTheme ChromaStyle = %q, want %q", th.ChromaStyle, "github-dark")
	}
	if got := colorString(th.H1.GetForeground()); got != "#8BA4D4" {
		t.Fatalf("DarkTheme H1 color = %q, want %q", got, "#8BA4D4")
	}
	if got := colorString(th.CodeBorder); got != "#3A3F4B" {
		t.Fatalf("DarkTheme CodeBorder = %q, want %q", got, "#3A3F4B")
	}
}

func TestLightThemeNotNil(t *testing.T) {
	th := LightTheme()
	if !th.H1.GetBold() {
		t.Fatal("LightTheme H1 style should be bold")
	}
	if th.H1.GetForeground() == nil {
		t.Fatal("LightTheme H1 style should have a foreground color")
	}
	if th.ChromaStyle != "xcode" {
		t.Fatalf("LightTheme ChromaStyle = %q, want %q", th.ChromaStyle, "xcode")
	}
	if got := colorString(th.H1.GetForeground()); got != "#3F5F8A" {
		t.Fatalf("LightTheme H1 color = %q, want %q", got, "#3F5F8A")
	}
	if got := colorString(th.CodeBorder); got != "#C2C7D0" {
		t.Fatalf("LightTheme CodeBorder = %q, want %q", got, "#C2C7D0")
	}
}

func TestDetectThemeFallbackDark(t *testing.T) {
	// Clear all env vars that could influence detection.
	t.Setenv("SHINE_THEME", "")
	t.Setenv("COLORFGBG", "")
	t.Setenv("TERM_PROGRAM", "")

	th := Detect("")
	if th.ChromaStyle != "github-dark" {
		t.Fatalf("Detect fallback: ChromaStyle = %q, want %q (dark)", th.ChromaStyle, "github-dark")
	}
}

func TestDetectThemeExplicitOverride(t *testing.T) {
	t.Setenv("SHINE_THEME", "light")
	t.Setenv("COLORFGBG", "")
	t.Setenv("TERM_PROGRAM", "")

	th := Detect("")
	if th.ChromaStyle != "xcode" {
		t.Fatalf("Detect with SHINE_THEME=light: ChromaStyle = %q, want %q", th.ChromaStyle, "xcode")
	}
}

func TestDetectThemeFlagOverride(t *testing.T) {
	// Flag should take priority over env vars.
	t.Setenv("SHINE_THEME", "dark")
	t.Setenv("COLORFGBG", "")
	t.Setenv("TERM_PROGRAM", "")

	th := Detect("light")
	if th.ChromaStyle != "xcode" {
		t.Fatalf("Detect with flag=light: ChromaStyle = %q, want %q", th.ChromaStyle, "xcode")
	}
}

func TestDetectCOLORFGBG_Dark(t *testing.T) {
	t.Setenv("SHINE_THEME", "")
	t.Setenv("COLORFGBG", "15;0")
	t.Setenv("TERM_PROGRAM", "")

	th := Detect("")
	if th.ChromaStyle != "github-dark" {
		t.Fatalf("COLORFGBG=15;0: ChromaStyle = %q, want %q (dark)", th.ChromaStyle, "github-dark")
	}
}

func TestDetectCOLORFGBG_LightMultiSegment(t *testing.T) {
	t.Setenv("SHINE_THEME", "")
	t.Setenv("COLORFGBG", "0;15;255")
	t.Setenv("TERM_PROGRAM", "")

	th := Detect("")
	if th.ChromaStyle != "xcode" {
		t.Fatalf("COLORFGBG=0;15;255: ChromaStyle = %q, want %q (light)", th.ChromaStyle, "xcode")
	}
}

func TestDetectCOLORFGBG_InvalidFallsThrough(t *testing.T) {
	t.Setenv("SHINE_THEME", "")
	t.Setenv("COLORFGBG", "invalid")
	t.Setenv("TERM_PROGRAM", "")

	// "invalid" has no semicolon, so COLORFGBG parsing fails.
	// No TERM_PROGRAM set, so falls through to default dark.
	th := Detect("")
	if th.ChromaStyle != "github-dark" {
		t.Fatalf("COLORFGBG=invalid: ChromaStyle = %q, want %q (dark fallback)", th.ChromaStyle, "github-dark")
	}
}

func TestDetectTERM_PROGRAM_AppleTerminal(t *testing.T) {
	t.Setenv("SHINE_THEME", "")
	t.Setenv("COLORFGBG", "")
	t.Setenv("TERM_PROGRAM", "Apple_Terminal")

	th := Detect("")
	if th.ChromaStyle != "xcode" {
		t.Fatalf("TERM_PROGRAM=Apple_Terminal: ChromaStyle = %q, want %q (light)", th.ChromaStyle, "xcode")
	}
}

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
	if th.ChromaStyle != "dracula" {
		t.Fatalf("DarkTheme ChromaStyle = %q, want %q", th.ChromaStyle, "dracula")
	}
	if got := colorString(th.H1.GetForeground()); got != "#9CCFD8" {
		t.Fatalf("DarkTheme H1 color = %q, want %q", got, "#9CCFD8")
	}
	if got := colorString(th.CodeBorder); got != "#393552" {
		t.Fatalf("DarkTheme CodeBorder = %q, want %q", got, "#393552")
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
	if th.ChromaStyle != "github" {
		t.Fatalf("LightTheme ChromaStyle = %q, want %q", th.ChromaStyle, "github")
	}
	if got := colorString(th.H1.GetForeground()); got != "#286983" {
		t.Fatalf("LightTheme H1 color = %q, want %q", got, "#286983")
	}
	if got := colorString(th.CodeBorder); got != "#9893A5" {
		t.Fatalf("LightTheme CodeBorder = %q, want %q", got, "#9893A5")
	}
}

func TestDetectThemeFallbackDark(t *testing.T) {
	// Clear all env vars that could influence detection.
	t.Setenv("SHINE_THEME", "")
	t.Setenv("COLORFGBG", "")
	t.Setenv("TERM_PROGRAM", "")

	th := Detect("")
	if th.ChromaStyle != "dracula" {
		t.Fatalf("Detect fallback: ChromaStyle = %q, want %q (dark)", th.ChromaStyle, "dracula")
	}
}

func TestDetectThemeExplicitOverride(t *testing.T) {
	t.Setenv("SHINE_THEME", "light")
	t.Setenv("COLORFGBG", "")
	t.Setenv("TERM_PROGRAM", "")

	th := Detect("")
	if th.ChromaStyle != "github" {
		t.Fatalf("Detect with SHINE_THEME=light: ChromaStyle = %q, want %q", th.ChromaStyle, "github")
	}
}

func TestDetectThemeFlagOverride(t *testing.T) {
	// Flag should take priority over env vars.
	t.Setenv("SHINE_THEME", "dark")
	t.Setenv("COLORFGBG", "")
	t.Setenv("TERM_PROGRAM", "")

	th := Detect("light")
	if th.ChromaStyle != "github" {
		t.Fatalf("Detect with flag=light: ChromaStyle = %q, want %q", th.ChromaStyle, "github")
	}
}

func TestDetectCOLORFGBG_Dark(t *testing.T) {
	t.Setenv("SHINE_THEME", "")
	t.Setenv("COLORFGBG", "15;0")
	t.Setenv("TERM_PROGRAM", "")

	th := Detect("")
	if th.ChromaStyle != "dracula" {
		t.Fatalf("COLORFGBG=15;0: ChromaStyle = %q, want %q (dark)", th.ChromaStyle, "dracula")
	}
}

func TestDetectCOLORFGBG_LightMultiSegment(t *testing.T) {
	t.Setenv("SHINE_THEME", "")
	t.Setenv("COLORFGBG", "0;15;255")
	t.Setenv("TERM_PROGRAM", "")

	th := Detect("")
	if th.ChromaStyle != "github" {
		t.Fatalf("COLORFGBG=0;15;255: ChromaStyle = %q, want %q (light)", th.ChromaStyle, "github")
	}
}

func TestDetectCOLORFGBG_InvalidFallsThrough(t *testing.T) {
	t.Setenv("SHINE_THEME", "")
	t.Setenv("COLORFGBG", "invalid")
	t.Setenv("TERM_PROGRAM", "")

	// "invalid" has no semicolon, so COLORFGBG parsing fails.
	// No TERM_PROGRAM set, so falls through to default dark.
	th := Detect("")
	if th.ChromaStyle != "dracula" {
		t.Fatalf("COLORFGBG=invalid: ChromaStyle = %q, want %q (dark fallback)", th.ChromaStyle, "dracula")
	}
}

func TestDetectTERM_PROGRAM_AppleTerminal(t *testing.T) {
	t.Setenv("SHINE_THEME", "")
	t.Setenv("COLORFGBG", "")
	t.Setenv("TERM_PROGRAM", "Apple_Terminal")

	th := Detect("")
	if th.ChromaStyle != "github" {
		t.Fatalf("TERM_PROGRAM=Apple_Terminal: ChromaStyle = %q, want %q (light)", th.ChromaStyle, "github")
	}
}

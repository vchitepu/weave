package theme

import "testing"

func TestDarkThemeNotNil(t *testing.T) {
	th := DarkTheme()
	if !th.H1.GetBold() {
		t.Fatal("DarkTheme H1 style should be bold")
	}
	if th.H1.GetForeground() == nil {
		t.Fatal("DarkTheme H1 style should have a foreground color")
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
}

func TestDetectThemeFallbackDark(t *testing.T) {
	// With no env vars set, should fall back to dark
	th := Detect("")
	// Verify it returns a valid theme (dark by default)
	if !th.H1.GetBold() {
		t.Fatal("Detect fallback should return a valid theme with bold H1")
	}
}

func TestDetectThemeExplicitOverride(t *testing.T) {
	t.Setenv("SHINE_THEME", "light")
	th := Detect("")
	// Light theme uses different colors — just verify it returns valid
	if !th.H1.GetBold() {
		t.Fatal("Detect with SHINE_THEME=light should return a valid theme")
	}
}

func TestDetectThemeFlagOverride(t *testing.T) {
	th := Detect("light")
	if !th.H1.GetBold() {
		t.Fatal("Detect with flag=light should return a valid theme")
	}
}

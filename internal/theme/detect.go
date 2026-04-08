package theme

import (
	"os"
	"strconv"
	"strings"
)

// Detect returns a Theme based on the flag override, env vars, or heuristics.
// flagValue is the --theme flag value ("dark", "light", or "" for auto).
func Detect(flagValue string) Theme {
	mode := detectMode(flagValue)
	if mode == "light" {
		return LightTheme()
	}
	return DarkTheme()
}

func detectMode(flagValue string) string {
	// 1. Flag override
	if flagValue == "dark" || flagValue == "light" {
		return flagValue
	}

	// 2. WEAVE_THEME env var
	if env := os.Getenv("WEAVE_THEME"); env == "dark" || env == "light" {
		return env
	}

	// 3. COLORFGBG — format "fg;bg", bg < 128 means dark
	if colorfgbg := os.Getenv("COLORFGBG"); colorfgbg != "" {
		parts := strings.Split(colorfgbg, ";")
		if len(parts) >= 2 {
			if bg, err := strconv.Atoi(parts[len(parts)-1]); err == nil {
				if bg < 128 {
					return "dark"
				}
				return "light"
			}
		}
	}

	// 4. TERM_PROGRAM heuristics
	termProgram := os.Getenv("TERM_PROGRAM")
	switch termProgram {
	case "Apple_Terminal":
		return "light"
	}

	// 5. Fallback
	return "dark"
}

package pager

import (
	"os"
	"os/exec"
	"strings"
)

// ShouldPage returns true if content lines exceed terminal height.
func ShouldPage(contentLines, termHeight int) bool {
	return contentLines > termHeight
}

// PagerCmd returns the pager command and its arguments.
func PagerCmd() (string, []string) {
	pager := os.Getenv("PAGER")
	if pager == "" {
		return "less", []string{"-R"}
	}
	parts := strings.Fields(pager)
	if len(parts) == 1 {
		return parts[0], nil
	}
	return parts[0], parts[1:]
}

// Run pipes the given content through the system pager.
func Run(content string) error {
	cmd, args := PagerCmd()
	pagerCmd := exec.Command(cmd, args...)
	pagerCmd.Stdin = strings.NewReader(content)
	pagerCmd.Stdout = os.Stdout
	pagerCmd.Stderr = os.Stderr
	return pagerCmd.Run()
}

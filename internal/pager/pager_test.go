package pager

import "testing"

func TestShouldPageTrue(t *testing.T) {
	if !ShouldPage(100, 24) {
		t.Fatal("expected ShouldPage=true for 100 lines in 24-line terminal")
	}
}

func TestShouldPageFalse(t *testing.T) {
	if ShouldPage(10, 24) {
		t.Fatal("expected ShouldPage=false for 10 lines in 24-line terminal")
	}
}

func TestPagerCommand(t *testing.T) {
	cmd, args := PagerCmd()
	if cmd == "" {
		t.Fatal("expected non-empty pager command")
	}
	_ = args
}

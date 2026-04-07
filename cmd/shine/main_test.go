package main

import "testing"

func TestNormalizeWidth_AutoWidthCappedAt120(t *testing.T) {
	if got := normalizeWidth(221, true); got != 120 {
		t.Fatalf("normalizeWidth(221, auto=true) = %d, want 120", got)
	}
}

func TestNormalizeWidth_ExplicitWidthNotCapped(t *testing.T) {
	if got := normalizeWidth(400, false); got != 400 {
		t.Fatalf("normalizeWidth(400, auto=false) = %d, want 400", got)
	}
}

func TestNormalizeWidth_Minimum20(t *testing.T) {
	if got := normalizeWidth(10, false); got != 20 {
		t.Fatalf("normalizeWidth(10, auto=false) = %d, want 20", got)
	}
}

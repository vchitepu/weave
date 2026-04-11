package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vchitepu/weave/internal/theme"
)

func TestHandler_RootReturnsHTML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.md")
	if err := os.WriteFile(path, []byte("# Hello World"), 0644); err != nil {
		t.Fatal(err)
	}

	inputs := []Input{{Path: path}}
	h := newHandler(inputs, theme.DarkTheme())

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}

	ct := res.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "text/html") {
		t.Fatalf("expected text/html, got %s", ct)
	}

	body, _ := io.ReadAll(res.Body)
	bodyStr := string(body)
	if !strings.Contains(bodyStr, "Hello World") {
		t.Fatalf("expected rendered content, got: %s", bodyStr[:200])
	}
	if !strings.Contains(bodyStr, "<html>") {
		t.Fatalf("expected full HTML page")
	}
	if !strings.Contains(bodyStr, "EventSource") {
		t.Fatalf("expected SSE script in page")
	}
}

func TestHandler_RootReRendersOnEachRequest(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.md")
	if err := os.WriteFile(path, []byte("# Version 1"), 0644); err != nil {
		t.Fatal(err)
	}

	inputs := []Input{{Path: path}}
	h := newHandler(inputs, theme.DarkTheme())

	req1 := httptest.NewRequest("GET", "/", nil)
	rec1 := httptest.NewRecorder()
	h.ServeHTTP(rec1, req1)
	body1, _ := io.ReadAll(rec1.Result().Body)

	if err := os.WriteFile(path, []byte("# Version 2"), 0644); err != nil {
		t.Fatal(err)
	}

	req2 := httptest.NewRequest("GET", "/", nil)
	rec2 := httptest.NewRecorder()
	h.ServeHTTP(rec2, req2)
	body2, _ := io.ReadAll(rec2.Result().Body)

	if !strings.Contains(string(body1), "Version 1") {
		t.Error("first request should contain Version 1")
	}
	if !strings.Contains(string(body2), "Version 2") {
		t.Error("second request should contain Version 2")
	}
}

func TestHandler_StdinInput(t *testing.T) {
	inputs := []Input{{Data: []byte("# From Stdin")}}
	h := newHandler(inputs, theme.DarkTheme())

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	body, _ := io.ReadAll(rec.Result().Body)
	if !strings.Contains(string(body), "From Stdin") {
		t.Fatalf("expected stdin content, got: %s", string(body)[:200])
	}
}

func TestHandler_MultiFile(t *testing.T) {
	dir := t.TempDir()
	path1 := filepath.Join(dir, "file1.md")
	path2 := filepath.Join(dir, "file2.md")
	os.WriteFile(path1, []byte("# First"), 0644)
	os.WriteFile(path2, []byte("# Second"), 0644)

	inputs := []Input{{Path: path1}, {Path: path2}}
	h := newHandler(inputs, theme.DarkTheme())

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	body, _ := io.ReadAll(rec.Result().Body)
	bodyStr := string(body)
	if !strings.Contains(bodyStr, "First") {
		t.Error("expected First file content")
	}
	if !strings.Contains(bodyStr, "Second") {
		t.Error("expected Second file content")
	}
	if !strings.Contains(bodyStr, "file-separator") {
		t.Error("expected file separator between files")
	}
}

func TestHandler_NotFound(t *testing.T) {
	inputs := []Input{{Data: []byte("# Test")}}
	h := newHandler(inputs, theme.DarkTheme())

	req := httptest.NewRequest("GET", "/nonexistent", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Result().StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Result().StatusCode)
	}
}

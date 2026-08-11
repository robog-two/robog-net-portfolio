package main

import (
	"strings"
	"testing"
)

// TestViewFitsTerminal is the regression test for the original jankiness: the
// old view drew roughly eleven rows of chrome while reserving eight, so it was
// taller than the screen on every frame.
func TestViewFitsTerminal(t *testing.T) {
	sizes := []struct{ w, h int }{
		{80, 24}, {120, 40}, {40, 12}, {minWidth, minHeight}, {200, 60},
	}
	for _, size := range sizes {
		m := NewModel(Config{OutputDir: t.TempDir()})
		m.editor.SetValue(strings.Repeat("a line of text\n", 200))

		updated, _ := m.resize(size.w, size.h)
		got := updated.(Model).View().Content

		if lines := strings.Count(got, "\n") + 1; lines > size.h {
			t.Errorf("%dx%d: view is %d rows, terminal has %d",
				size.w, size.h, lines, size.h)
		}
	}
}

func TestViewFitsTerminalInPreview(t *testing.T) {
	m := NewModel(Config{OutputDir: t.TempDir()})
	m.editor.SetValue("# Heading\n\n" + strings.Repeat("body text ", 400))

	updated, _ := m.resize(80, 24)
	m = updated.(Model)
	updated, _ = m.enterPreview()
	m = updated.(Model)

	if m.mode != modePreview {
		t.Fatalf("mode = %v, want preview (status: %q)", m.mode, m.status)
	}
	got := m.View().Content
	if lines := strings.Count(got, "\n") + 1; lines > 24 {
		t.Errorf("preview is %d rows, terminal has 24", lines)
	}
}

func TestLayoutReservesExactlyTheChromeItDraws(t *testing.T) {
	// If someone adds a row to View without adding it to layout.go, this
	// fails.
	m := NewModel(Config{OutputDir: t.TempDir()})
	updated, _ := m.resize(80, 24)
	m = updated.(Model)

	if got, want := m.layout.ContentHeight, 24-chromeRows; got != want {
		t.Errorf("ContentHeight = %d, want %d", got, want)
	}
	rows := strings.Split(m.View().Content, "\n")
	if got := len(rows) - m.layout.ContentHeight; got != chromeRows {
		t.Errorf("view drew %d rows of chrome, layout reserves %d", got, chromeRows)
	}
}

func TestTooSmallTerminalDoesNotOverflow(t *testing.T) {
	m := NewModel(Config{OutputDir: t.TempDir()})
	updated, _ := m.resize(20, 5)
	m = updated.(Model)

	if !m.layout.TooSmall {
		t.Fatal("20x5 should be flagged too small")
	}
	got := m.View().Content
	if lines := strings.Count(got, "\n") + 1; lines > 5 {
		t.Errorf("too-small notice is %d rows, terminal has 5", lines)
	}
}

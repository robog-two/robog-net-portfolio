package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// ctrl builds a ctrl+<r> key press, and tab/esc build those keys, so the tests
// exercise the same key.Matches path the running program does.
func ctrl(r rune) tea.KeyPressMsg {
	return tea.KeyPressMsg{Code: r, Mod: tea.ModCtrl}
}

func plain(r rune) tea.KeyPressMsg {
	return tea.KeyPressMsg{Code: r, Text: string(r)}
}

func special(code rune) tea.KeyPressMsg {
	return tea.KeyPressMsg{Code: code}
}

func newTestModel(t *testing.T) Model {
	t.Helper()
	m := NewModel(Config{OutputDir: t.TempDir()})
	updated, _ := m.resize(80, 24)
	return updated.(Model)
}

func send(m Model, msgs ...tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	for _, msg := range msgs {
		var updated tea.Model
		updated, cmd = m.Update(msg)
		m = updated.(Model)
	}
	return m, cmd
}

// TestSaveIsOneKeystroke covers the central UX change: saving no longer means
// tabbing to a button row and arrowing across it.
func TestSaveIsOneKeystroke(t *testing.T) {
	m := newTestModel(t)
	m.editor.SetValue("Business cards printed!")

	m, _ = send(m, ctrl('s'))

	if m.savedPath == "" {
		t.Fatalf("ctrl+s did not save (status: %q)", m.status)
	}
	if !strings.HasSuffix(m.savedPath, "-business-cards-printed.md") {
		t.Errorf("saved to %s", m.savedPath)
	}
	data, err := os.ReadFile(m.savedPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "snippet: yes") {
		t.Errorf("front matter missing from %q", data)
	}
}

// TestSaveDoesNotQuit: saving used to return tea.Quit, so an accidental save
// ended the session and you had to relaunch to keep writing.
func TestSaveDoesNotQuit(t *testing.T) {
	m := newTestModel(t)
	m.editor.SetValue("still writing")

	m, _ = send(m, ctrl('s'))
	m.editor.InsertText(" — and more")
	m, _ = send(m, ctrl('s'))

	if got := m.editor.Value(); got != "still writing — and more" {
		t.Errorf("buffer after re-save = %q", got)
	}
	entries, _ := os.ReadDir(filepath.Dir(m.savedPath))
	if len(entries) != 1 {
		t.Errorf("re-saving produced %d files, want 1", len(entries))
	}
}

// TestTabReachesTheEditor: tab used to be stolen for focus switching, so an
// indented code block could not be typed.
func TestTabReachesTheEditor(t *testing.T) {
	m := newTestModel(t)
	m.editor.SetValue("- item\n")
	m, _ = send(m, special(tea.KeyTab))
	m.editor.InsertText("- nested")

	if got, want := m.editor.Value(), "- item\n"+indent+"- nested"; got != want {
		t.Errorf("editor = %q, want %q", got, want)
	}
}

// TestQuitConfirmsWhenDirty guards unsaved work behind a second press.
func TestQuitConfirmsWhenDirty(t *testing.T) {
	m := newTestModel(t)
	m.editor.SetValue("unsaved draft")

	m, cmd := send(m, ctrl('c'))
	if cmd != nil {
		t.Fatal("first ctrl+c quit with unsaved changes")
	}
	if !strings.Contains(m.status, "unsaved") {
		t.Errorf("status = %q, want an unsaved-changes warning", m.status)
	}

	_, cmd = send(m, ctrl('c'))
	if cmd == nil {
		t.Error("second ctrl+c did not quit")
	}
}

func TestQuitIsImmediateWhenClean(t *testing.T) {
	m := newTestModel(t)
	if _, cmd := send(m, ctrl('c')); cmd == nil {
		t.Error("ctrl+c on an empty buffer should quit at once")
	}

	m = newTestModel(t)
	m.editor.SetValue("saved work")
	m, _ = send(m, ctrl('s'))
	if _, cmd := send(m, ctrl('c')); cmd == nil {
		t.Error("ctrl+c after saving should quit at once")
	}
}

func TestTypingCancelsPendingQuit(t *testing.T) {
	m := newTestModel(t)
	m.editor.SetValue("draft")

	m, _ = send(m, ctrl('c')) // arms the confirmation
	m, _ = send(m, plain('x'))
	if m.quitting {
		t.Error("typing should cancel a pending quit")
	}
	if _, cmd := send(m, ctrl('c')); cmd != nil {
		t.Error("ctrl+c after typing should re-arm, not quit")
	}
}

// TestErrorsDoNotDestroyTheDraft: a failed save used to replace the entire
// view with an error and an instruction to quit.
func TestErrorsDoNotDestroyTheDraft(t *testing.T) {
	dir := t.TempDir()
	blocked := filepath.Join(dir, "blocked")
	if err := os.WriteFile(blocked, []byte("not a directory"), 0o644); err != nil {
		t.Fatal(err)
	}

	m := NewModel(Config{OutputDir: blocked})
	updated, _ := m.resize(80, 24)
	m = updated.(Model)
	m.editor.SetValue("precious draft")

	m, _ = send(m, ctrl('s'))

	if m.savedPath != "" {
		t.Fatal("save should have failed")
	}
	if m.statusLevel != statusError || m.status == "" {
		t.Errorf("statusLevel = %v, status = %q; want an error message", m.statusLevel, m.status)
	}
	if got := m.editor.Value(); got != "precious draft" {
		t.Errorf("draft was lost: %q", got)
	}
	if !strings.Contains(m.View().Content, "precious draft") {
		t.Error("editor is no longer rendered after an error")
	}
}

func TestPreviewRoundTrip(t *testing.T) {
	m := newTestModel(t)
	m.editor.SetValue("# Hello\n\nsome *body* text")

	m, _ = send(m, ctrl('p'))
	if m.mode != modePreview {
		t.Fatalf("ctrl+p did not enter preview (status: %q)", m.status)
	}

	m, _ = send(m, special(tea.KeyEsc))
	if m.mode != modeEdit {
		t.Error("esc did not return to the editor")
	}
	if got := m.editor.Value(); got != "# Hello\n\nsome *body* text" {
		t.Errorf("buffer changed across preview: %q", got)
	}
}

func TestPreviewRefusesEmptyBuffer(t *testing.T) {
	m := newTestModel(t)
	m, _ = send(m, ctrl('p'))
	if m.mode != modeEdit {
		t.Error("previewing an empty buffer should stay in the editor")
	}
	if m.statusLevel != statusWarn {
		t.Errorf("statusLevel = %v, want a warning", m.statusLevel)
	}
}

func TestSaveRefusesEmptyBuffer(t *testing.T) {
	m := newTestModel(t)
	m, _ = send(m, ctrl('s'))
	if m.savedPath != "" {
		t.Error("an empty buffer should not be saved")
	}
	entries, _ := os.ReadDir(m.cfg.Dir())
	if len(entries) != 0 {
		t.Errorf("empty save created %d files", len(entries))
	}
}

// TestPasteGoesThroughThePipeline checks the preprocessing hook is actually
// wired to the editor.
func TestPasteGoesThroughThePipeline(t *testing.T) {
	m := newTestModel(t)
	m, _ = send(m, tea.PasteMsg{Content: "pasted\r\ntext"})

	if got := m.editor.Value(); got != "pasted\ntext" {
		t.Errorf("editor = %q, want the normalised paste", got)
	}
}

func TestFileDropReportsThroughTheStatusLine(t *testing.T) {
	m := newTestModel(t)
	dropped := filepath.Join(t.TempDir(), "toot.mp3")
	if err := os.WriteFile(dropped, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	m, _ = send(m, tea.PasteMsg{Content: dropped})

	if got := m.editor.Value(); got != dropped {
		t.Errorf("editor = %q, want the dropped path", got)
	}
	if !strings.Contains(m.status, "not enabled") {
		t.Errorf("status = %q, want the attachments explanation", m.status)
	}
}

func TestStatusExpires(t *testing.T) {
	m := newTestModel(t)
	m = m.setStatus(statusInfo, "hello")
	gen := m.statusGen

	m, _ = send(m, statusExpiredMsg{gen: gen})
	if m.status != "" {
		t.Errorf("status = %q, want cleared", m.status)
	}

	// A stale timer must not wipe a newer message.
	m = m.setStatus(statusInfo, "newer")
	m, _ = send(m, statusExpiredMsg{gen: gen})
	if m.status != "newer" {
		t.Errorf("a stale expiry cleared %q", "newer")
	}
}

func TestBackgroundColorRestylesEverything(t *testing.T) {
	m := newTestModel(t)
	m, _ = send(m, tea.BackgroundColorMsg{Color: lipgloss.Color("#ffffff")})
	if m.theme.Glamour != "light" {
		t.Errorf("Glamour style = %q, want light for a light background", m.theme.Glamour)
	}
	if m.renderer != nil {
		t.Error("renderer should be invalidated when the theme changes")
	}
}

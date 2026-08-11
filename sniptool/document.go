package main

import (
	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
)

// Editor is the seam between the model and whatever is actually doing the
// editing.
//
// There is exactly one implementation today — textareaEditor, a thin wrapper
// over bubbles' textarea that edits Markdown source as plain text. The
// interface exists so a richer, more WYSIWYG-ish editor (one that styles
// headings and emphasis in place, draws a link as a chip rather than raw
// [text](url), shows an attachment as a tile) can arrive as a second
// implementation without model.go, view.go or input.go changing. Nothing
// outside this file may reach past the interface to the textarea.
type Editor interface {
	Focus() tea.Cmd
	Blur()
	Focused() bool

	// SetSize is given the exact cell rectangle the editor may draw into.
	// Callers guarantee width and height are at least 1.
	SetSize(width, height int)

	// Value is the Markdown source, and is what gets written to disk. Any
	// implementation must round-trip it: after SetValue(v), Value() == v.
	// A WYSIWYG editor may render something quite different from this; it
	// still has to be able to hand back canonical source.
	Value() string
	SetValue(string)

	// InsertText inserts at the cursor. The input pipeline routes every paste
	// and drop through here, so a future editor inherits attachment insertion
	// without doing anything.
	InsertText(string)

	Update(tea.Msg) (Editor, tea.Cmd)
	View() string

	// Cursor is the real terminal cursor, positioned relative to the editor's
	// own top-left cell. view.go offsets it by the editor's screen position.
	// Returning nil hides the cursor.
	Cursor() *tea.Cursor
}

// textareaEditor is the plain-source editor.
type textareaEditor struct {
	ta textarea.Model
}

func newTextareaEditor(t Theme) *textareaEditor {
	ta := textarea.New()
	ta.Placeholder = "Write your post here. Markdown + HTML supported."
	ta.Prompt = "│ "
	ta.ShowLineNumbers = false
	ta.CharLimit = 0

	// Use the real terminal cursor rather than the drawn one: it blinks
	// correctly without us pumping blink messages, and it survives focus
	// changes. See view.go, which places it via the returned tea.Cursor.
	ta.SetVirtualCursor(false)
	ta.SetStyles(t.Textarea)

	return &textareaEditor{ta: ta}
}

func (e *textareaEditor) Focus() tea.Cmd { return e.ta.Focus() }
func (e *textareaEditor) Blur()          { e.ta.Blur() }
func (e *textareaEditor) Focused() bool  { return e.ta.Focused() }

func (e *textareaEditor) SetSize(w, h int) {
	// Prompt and ShowLineNumbers must be set before SetWidth; they are, in
	// newTextareaEditor.
	e.ta.SetWidth(w)
	e.ta.SetHeight(h)
}

func (e *textareaEditor) Value() string       { return e.ta.Value() }
func (e *textareaEditor) SetValue(s string)   { e.ta.SetValue(s) }
func (e *textareaEditor) InsertText(s string) { e.ta.InsertString(s) }
func (e *textareaEditor) View() string        { return e.ta.View() }
func (e *textareaEditor) Cursor() *tea.Cursor { return e.ta.Cursor() }

// indent is what the tab key inserts. Spaces rather than a literal tab because
// Markdown's rules for tabs differ between code blocks and list nesting, and
// two of them is one level of list nesting.
const indent = "  "

func (e *textareaEditor) Update(msg tea.Msg) (Editor, tea.Cmd) {
	// bubbles' textarea has no tab binding — a tab press carries no printable
	// text, so it would otherwise do nothing at all. The old tool hid this by
	// stealing tab for focus switching; now that tab reaches the editor, the
	// editor has to act on it.
	if key, ok := msg.(tea.KeyPressMsg); ok && key.Code == tea.KeyTab && key.Mod == 0 {
		e.ta.InsertString(indent)
		return e, nil
	}

	var cmd tea.Cmd
	e.ta, cmd = e.ta.Update(msg)
	return e, cmd
}

// applyTheme lets the model restyle the editor when the terminal reports its
// background colour. Kept off the Editor interface because it is a detail of
// how this implementation is styled, not of editing.
func (e *textareaEditor) applyTheme(t Theme) { e.ta.SetStyles(t.Textarea) }

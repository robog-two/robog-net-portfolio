package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// View assembles exactly the rows layout.go accounts for — no more, no fewer —
// and positions the real terminal cursor over the editor.
func (m Model) View() tea.View {
	v := tea.NewView("")
	v.AltScreen = true
	v.WindowTitle = "snip"

	if m.layout.TooSmall {
		v.SetContent(m.theme.TooSmall.Render(fmt.Sprintf(
			"Terminal too small.\nNeed at least %d×%d, have %d×%d.",
			minWidth, minHeight, m.layout.Width, m.layout.Height,
		)))
		return v
	}

	var content string
	if m.mode == modePreview {
		content = m.preview.View()
	} else {
		content = m.editor.View()
	}

	rows := []string{
		m.titleBar(),
		"",
		content,
		"",
		m.statusLine(),
		m.helpLine(),
	}
	v.SetContent(strings.Join(rows, "\n"))

	// The cursor belongs to the editor, offset by where the editor sits on
	// screen. In v2 this is declarative, so it stays correct across renders
	// without any show/hide commands.
	if m.mode == modeEdit {
		if c := m.editor.Cursor(); c != nil {
			c.Position.Y += contentTop
			v.Cursor = c
		}
	}
	return v
}

// titleBar shows the mode on the left and the pending filename on the right.
// The filename comes from the same Post value that save writes, so what is
// shown is what lands on disk.
func (m Model) titleBar() string {
	left := m.theme.Title.Render("New post")
	if m.mode == modePreview {
		left = m.theme.Title.Render("Preview")
		if pct := m.preview.ScrollPercent(); m.preview.TotalLineCount() > m.preview.Height() {
			left += m.theme.Muted.Render(fmt.Sprintf("  %3.0f%%", pct*100))
		}
	}

	right := m.theme.Filename.Render(m.displayPath())
	gap := m.layout.Width - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 1 {
		// Too narrow for both; the mode label is the more useful half.
		return truncate(left, m.layout.Width)
	}
	return left + strings.Repeat(" ", gap) + right
}

func (m Model) statusLine() string {
	if m.status == "" {
		return ""
	}
	style := m.theme.Status
	switch m.statusLevel {
	case statusWarn:
		style = m.theme.Warning
	case statusError:
		style = m.theme.Error
	}
	return truncate(style.Render(m.status), m.layout.Width)
}

// helpLine is generated from the key bindings, so it cannot describe keys that
// do not exist. The old version was a hand-written string per focus state.
func (m Model) helpLine() string {
	if m.mode == modePreview {
		return m.help.View(previewHelp(m.keys))
	}
	return m.help.View(editHelp(m.keys))
}

// displayPath is the post's path, shortened relative to the working directory
// when that is shorter, so the title bar is not dominated by /Users/....
func (m Model) displayPath() string {
	path := m.savedPath
	if path == "" {
		path = filepath.Join(m.cfg.Dir(), m.post().Filename())
	}
	if cwd, err := os.Getwd(); err == nil {
		if rel, err := filepath.Rel(cwd, path); err == nil && !strings.HasPrefix(rel, "..") {
			return rel
		}
	}
	return path
}

func truncate(s string, width int) string {
	if width <= 0 || lipgloss.Width(s) <= width {
		return s
	}
	return lipgloss.NewStyle().MaxWidth(width).Render(s)
}

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress Ctrl+C to quit.\n", m.err)
	}

	if m.appState == stateMainMenu {
		return m.renderMainMenu()
	}
	if m.appState == statePickThread {
		return m.renderPickThread()
	}
	if m.appState == stateNewThread {
		return m.renderNewThread()
	}

	// stateNewSnippet or stateAppendEntry: preview/editor
	if m.mode == modePreview {
		header := headerStyle.Render("── Preview ──")
		help := helpStyle.Render("tab/esc: back to editor • ctrl+c: quit")
		return fmt.Sprintf("%s\n\n%s\n%s\n", header, previewBoxStyle.Render(m.preview), help)
	}

	// Edit mode for snippet or append entry
	now := time.Now()
	preview := m.textarea.Value()
	if strings.TrimSpace(preview) == "" {
		preview = "(empty)"
	}

	var info string
	if m.appState == stateNewSnippet {
		fileName := fmt.Sprintf("%s-%s.md", now.Format("2006-01-02"), slugify(preview))
		info = fmt.Sprintf("File: %s   Time: %s",
			outputDir+"/"+fileName,
			now.Format("Mon Jan 2 3:04 PM"),
		)
	} else if m.appState == stateAppendEntry {
		fileName := fmt.Sprintf("%s-%s.md", now.Format("2006-01-02"), slugify(preview))
		info = fmt.Sprintf("File: %s/%s   Time: %s",
			outputDir+"/"+m.selectedSlug,
			fileName,
			now.Format("Mon Jan 2 3:04 PM"),
		)
	}

	var buttons []string
	for id := buttonID(0); id < numButtons; id++ {
		label := buttonLabels[id]
		if m.focus == focusButtons && id == m.activeButton {
			buttons = append(buttons, buttonFocusedStyle.Render(label))
		} else {
			buttons = append(buttons, buttonStyle.Render(label))
		}
	}
	buttonRow := lipgloss.JoinHorizontal(lipgloss.Top, buttons...)

	var help string
	if m.focus == focusEditor {
		help = "tab: focus buttons • ctrl+c: quit"
	} else {
		help = "←/→: pick button • enter: activate • tab: back to editor • ctrl+c: quit"
	}

	status := ""
	if m.statusMsg != "" {
		status = "\n" + helpStyle.Render(m.statusMsg)
	}

	title := "New post"
	if m.appState == stateAppendEntry {
		title = "Add to thread"
	}

	return fmt.Sprintf(
		"%s\n%s\n\n%s\n\n%s%s\n%s\n",
		headerStyle.Render(title),
		helpStyle.Render(info),
		m.textarea.View(),
		buttonRow,
		status,
		helpStyle.Render(help),
	)
}

func (m model) renderMainMenu() string {
	var items []string
	for id := menuNewSnippet; id < numMenuItems; id++ {
		label := menuLabels[id]
		if int(id) == m.menuCursor {
			items = append(items, menuItemFocusedStyle.Render("> "+label))
		} else {
			items = append(items, menuItemStyle.Render("  "+label))
		}
	}

	help := helpStyle.Render("↑/↓: navigate • enter: select • ctrl+c: quit")
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s\n",
		headerStyle.Render("Sniptool"),
		strings.Join(items, "\n"),
		help,
	)
}

func (m model) renderPickThread() string {
	header := headerStyle.Render("Pick a thread to append to:")
	help := helpStyle.Render("↑/↓/j/k: navigate • enter: select • esc: back")
	return fmt.Sprintf("%s\n%s\n%s\n", header, m.threadList.View(), help)
}

func (m model) renderNewThread() string {
	titleLabel := "Title"
	bodyLabel := "Body"

	switch m.focus {
	case focusTitleInput:
		titleLabel = headerStyle.Render(titleLabel + " (focused)")
	case focusThreadBody:
		bodyLabel = headerStyle.Render(bodyLabel + " (focused)")
	}

	var buttons []string
	for id := buttonID(0); id < numButtons; id++ {
		label := buttonLabels[id]
		if m.focus == focusButtons && id == m.activeButton {
			buttons = append(buttons, buttonFocusedStyle.Render(label))
		} else {
			buttons = append(buttons, buttonStyle.Render(label))
		}
	}
	buttonRow := lipgloss.JoinHorizontal(lipgloss.Top, buttons...)

	var help string
	if m.focus == focusButtons {
		help = "←/→: pick button • enter: save • tab: back to fields • esc: cancel"
	} else {
		help = "tab: next field • esc: cancel"
	}

	return fmt.Sprintf(
		"%s\n%s\n%s\n\n%s\n%s\n\n%s\n\n%s\n%s\n",
		headerStyle.Render("New Thread Post"),
		titleLabel,
		m.titleInput.View(),
		bodyLabel,
		m.threadBody.View(),
		buttonRow,
		helpStyle.Render(help),
		helpStyle.Render("ctrl+c: quit"),
	)
}

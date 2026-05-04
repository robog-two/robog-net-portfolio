package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	if !m.ready {
		return ""
	}

	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress Ctrl+C to quit.\n", m.err)
	}

	var content string
	switch m.appState {
	case stateHome:
		content = m.renderHome()
	case stateEditor:
		content = m.renderEditor()
	case statePreview:
		content = m.renderPreview()
	}

	return lipgloss.NewStyle().Width(m.width).Render(content)
}

// --- home screen ---------------------------------------------------------

func (m model) renderHome() string {
	header := homeHeaderStyle.Render("Sniptool")

	var items []string
	for i, item := range m.homeItems {
		label := item.label
		if i == m.homeCursor {
			if item.kind == homeAction {
				items = append(items, actionItemFocusedStyle.Render(label))
			} else {
				items = append(items, threadItemFocusedStyle.Render(label))
			}
		} else {
			if item.kind == homeAction {
				items = append(items, actionItemStyle.Render(label))
			} else {
				items = append(items, threadItemStyle.Render(label))
			}
		}
	}

	var banner string
	if m.savedMsg != "" {
		banner = savedBannerStyle.Render("✓ " + m.savedMsg)
	}

	help := homeHelpStyle.Render("↑/↓: navigate • enter: select • ctrl+c: quit")

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		"",
		lipgloss.JoinVertical(lipgloss.Left, items...),
		"",
		banner,
		"",
		help,
	)
}

// --- editor screen -------------------------------------------------------

func (m model) renderEditor() string {
	// Build breadcrumb header
	var breadcrumb string
	switch m.editKind {
	case editSnippet:
		breadcrumb = "Sniptool › New Snippet Post"
	case editThreadPost:
		breadcrumb = "Sniptool › New Thread Post"
	case editThreadReply:
		// Find thread title
		threadTitle := "Unknown"
		for _, t := range m.threads {
			if t.slug == m.editSlug {
				threadTitle = t.title
				break
			}
		}
		breadcrumb = fmt.Sprintf("Sniptool › Reply to: %s", threadTitle)
	}
	header := editorHeaderStyle.Render(breadcrumb)

	// Build file path preview
	var filePath string
	preview := strings.TrimSpace(m.textarea.Value())
	if preview == "" {
		preview = "(empty)"
	}

	switch m.editKind {
	case editSnippet:
		filePath = fmt.Sprintf("File: %s/YYYY-MM-DD-%s.md", outputDir, slugify(preview))
	case editThreadPost:
		if m.editTitle == "" {
			filePath = fmt.Sprintf("File: %s/(untitled).md", blogDir)
		} else {
			filePath = fmt.Sprintf("File: %s/%s.md", blogDir, slugify(m.editTitle))
		}
	case editThreadReply:
		filePath = fmt.Sprintf("File: %s/%s/YYYY-MM-DD-%s.md", outputDir, m.editSlug, slugify(preview))
	}
	info := infoStyle.Render(filePath)

	// Title field (only for thread posts)
	var titleField string
	if m.editKind == editThreadPost {
		titleLabel := "Title"
		if m.focus == focusTitle {
			titleLabel = titleLabelStyle.Render(titleLabel + " (editing)")
		}
		titleField = lipgloss.JoinVertical(lipgloss.Left, titleLabel, m.titleInput.View(), "")
	}

	// Body field
	bodyLabel := "Content"
	if m.focus == focusEditor {
		bodyLabel = bodyLabelStyle.Render(bodyLabel + " (editing)")
	}
	body := lipgloss.JoinVertical(lipgloss.Left, bodyLabel, editorBoxStyle.Render(m.textarea.View()), "")

	// Buttons
	var buttons []string
	btnLabels := []string{"Preview", "Save", "Back"}
	for i := 0; i < 3; i++ {
		id := buttonID(i)
		label := btnLabels[i]
		if m.focus == focusButtons && id == m.activeButton {
			buttons = append(buttons, buttonFocusedStyle.Render(label))
		} else {
			buttons = append(buttons, buttonStyle.Render(label))
		}
	}
	buttonRow := lipgloss.JoinHorizontal(lipgloss.Top, buttons...)

	// Help text
	var help string
	if m.focus == focusButtons {
		help = "←/→: pick button • enter/space: select • esc: back"
	} else if m.editKind == editThreadPost && m.focus == focusTitle {
		help = "tab: next field • ctrl+c: quit"
	} else {
		help = "tab: next field • ctrl+c: quit"
	}
	helpText := editorHelpStyle.Render(help)

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		info,
		titleField,
		body,
		buttonRow,
		"",
		helpText,
	)
}

// --- preview screen ------------------------------------------------------

func (m model) renderPreview() string {
	header := previewHeaderStyle.Render("── Preview ──")

	// Buttons: Back and Save
	var buttons []string
	btnLabels := []string{"← Back to Editor", "Save & Return Home"}
	for i := 0; i < 2; i++ {
		id := buttonID(i)
		if m.focus == focusButtons && id == m.activeButton {
			buttons = append(buttons, previewButtonFocusedStyle.Render(btnLabels[i]))
		} else {
			buttons = append(buttons, previewButtonStyle.Render(btnLabels[i]))
		}
	}
	buttonRow := lipgloss.JoinHorizontal(lipgloss.Top, buttons...)

	help := previewHelpStyle.Render("↑/↓: scroll • ←/→: pick button • enter/space: select • esc: back to editor")

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		"",
		m.viewport.View(),
		"",
		buttonRow,
		"",
		help,
	)
}

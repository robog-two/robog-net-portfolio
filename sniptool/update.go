package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		m.resizeComponents()
		return m, nil

	case tea.KeyMsg:
		// Global quit
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		// State-specific handling
		switch m.appState {
		case stateHome:
			return m.updateHome(msg)
		case stateEditor:
			return m.updateEditor(msg)
		case statePreview:
			return m.updatePreview(msg)
		}
	}

	return m, nil
}

// --- home screen updates -------------------------------------------------

func (m model) updateHome(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.homeCursor > 0 {
			m.homeCursor--
		}
	case "down", "j":
		if m.homeCursor < len(m.homeItems)-1 {
			m.homeCursor++
		}
	case "enter", " ":
		item := m.homeItems[m.homeCursor]
		if item.kind == homeAction {
			// First action is snippet, second is thread post
			if m.homeCursor == 0 {
				m.editKind = editSnippet
			} else {
				m.editKind = editThreadPost
			}
		} else {
			m.editKind = editThreadReply
			m.editSlug = item.slug
		}

		// Enter editor
		m.appState = stateEditor
		m.focus = focusEditor
		if m.editKind == editThreadPost {
			m.focus = focusTitle
			m.titleInput.Focus()
		} else {
			m.textarea.Focus()
		}
		m.textarea.Reset()
		m.titleInput.Reset()
		m.editTitle = ""
		m.savedMsg = ""
		m.resizeComponents()
	}

	return m, nil
}

// --- editor screen updates -----------------------------------------------

func (m model) updateEditor(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab":
		switch m.focus {
		case focusTitle:
			m.titleInput.Blur()
			m.focus = focusEditor
			m.textarea.Focus()
		case focusEditor:
			m.textarea.Blur()
			m.focus = focusButtons
		case focusButtons:
			m.focus = focusTitle
			m.titleInput.Focus()
		}
		return m, nil

	case "esc":
		// Go back to home
		m.appState = stateHome
		m.loadHomeItems()
		m.focus = focusTitle
		return m, nil

	case "left", "h":
		if m.focus == focusButtons && m.activeButton > 0 {
			m.activeButton--
		}
		return m, nil

	case "right", "l":
		if m.focus == focusButtons && m.activeButton < numButtons-1 {
			m.activeButton++
		}
		return m, nil

	case "enter", " ":
		if m.focus == focusButtons {
			return m.activateEditorButton()
		}
	}

	// Forward keystrokes to focused widget
	if m.focus == focusTitle {
		var cmd tea.Cmd
		m.titleInput, cmd = m.titleInput.Update(msg)
		m.editTitle = m.titleInput.Value()
		return m, cmd
	}

	if m.focus == focusEditor {
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) activateEditorButton() (tea.Model, tea.Cmd) {
	switch m.activeButton {
	case btnPreview:
		body := strings.TrimSpace(m.textarea.Value())
		if body == "" {
			return m, nil // silently do nothing
		}
		rendered, err := renderMarkdown(body)
		if err != nil {
			m.err = err
			return m, nil
		}
		m.preview = rendered
		m.viewport.SetContent(rendered)
		m.viewport.GotoTop()
		m.appState = statePreview
		m.focus = focusButtons
		m.activeButton = btnSave // default to Save button in preview
		return m, nil

	case btnSave:
		return m.saveContent()

	case btnQuit:
		m.appState = stateHome
		m.loadHomeItems()
		m.focus = focusTitle
		return m, nil
	}

	return m, nil
}

func (m model) saveContent() (tea.Model, tea.Cmd) {
	var path string
	var err error

	switch m.editKind {
	case editSnippet:
		body := strings.TrimSpace(m.textarea.Value())
		if body == "" {
			return m, nil
		}
		path, err = saveSnippet(body)

	case editThreadPost:
		title := strings.TrimSpace(m.editTitle)
		body := strings.TrimSpace(m.textarea.Value())
		if title == "" || body == "" {
			return m, nil
		}
		path, err = saveThreadPost(title, body)

	case editThreadReply:
		body := strings.TrimSpace(m.textarea.Value())
		if body == "" {
			return m, nil
		}
		path, err = appendThreadEntry(m.editSlug, body)
	}

	if err != nil {
		m.err = err
		return m, nil
	}

	// Success: show banner and return to home
	m.lastPath = path
	m.savedMsg = "Post saved"
	m.appState = stateHome
	m.loadHomeItems()
	m.focus = focusTitle
	return m, nil
}

// --- preview screen updates ----------------------------------------------

func (m model) updatePreview(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Back to editor
			m.appState = stateEditor
			m.focus = focusEditor
			if m.editKind == editThreadPost {
				m.focus = focusTitle
			}
			return m, nil

		case "left", "h":
			if m.activeButton > 0 {
				m.activeButton--
			}
			return m, nil

		case "right", "l":
			if m.activeButton < 1 {
				m.activeButton++
			}
			return m, nil

		case "enter", " ":
			switch m.activeButton {
			case 0: // Back button
				m.appState = stateEditor
				m.focus = focusEditor
				if m.editKind == editThreadPost {
					m.focus = focusTitle
				}
			case 1: // Save button
				return m.saveContent()
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

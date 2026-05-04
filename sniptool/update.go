package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		taWidth := msg.Width - uiPadding
		if taWidth < minTextareaWidth {
			taWidth = minTextareaWidth
		}
		taHeight := msg.Height - uiMargin
		if taHeight < minTextareaHeight {
			taHeight = minTextareaHeight
		}
		m.textarea.SetWidth(taWidth)
		m.textarea.SetHeight(taHeight)
		m.threadBody.SetWidth(taWidth)
		m.threadBody.SetHeight(taHeight - 3)
		m.titleInput.Width = taWidth
		m.threadList.SetWidth(taWidth)
		m.threadList.SetHeight(taHeight)
		return m, nil

	case tea.KeyMsg:
		// Global quit
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		// Main menu
		if m.appState == stateMainMenu {
			switch msg.String() {
			case "up", "k":
				if m.menuCursor > 0 {
					m.menuCursor--
				}
			case "down", "j":
				if m.menuCursor < int(numMenuItems)-1 {
					m.menuCursor++
				}
			case "enter", " ":
				switch menuItem(m.menuCursor) {
				case menuNewSnippet:
					m.appState = stateNewSnippet
					m.focus = focusEditor
					m.textarea.Focus()
					m.statusMsg = ""
				case menuNewThread:
					m.appState = stateNewThread
					m.focus = focusTitleInput
					m.titleFocused = true
					m.titleInput.Focus()
					m.statusMsg = ""
				case menuAppendThread:
					threads, _ := scanThreadPosts()
					m.threadPosts = threads
					items := make([]list.Item, len(threads))
					for i, t := range threads {
						items[i] = t
					}
					m.threadList.SetItems(items)
					m.appState = statePickThread
					m.focus = focusList
					m.statusMsg = ""
				}
			}
			return m, nil
		}

		// Pick thread list
		if m.appState == statePickThread {
			if msg.String() == "enter" {
				if sel, ok := m.threadList.SelectedItem().(threadPostMeta); ok {
					m.selectedSlug = sel.slug
					m.appState = stateAppendEntry
					m.focus = focusEditor
					m.textarea.Reset()
					m.textarea.Focus()
					m.statusMsg = ""
				}
			} else if msg.String() == "esc" {
				m.appState = stateMainMenu
				m.menuCursor = 0
				m.statusMsg = ""
				return m, nil
			}
			var cmd tea.Cmd
			m.threadList, cmd = m.threadList.Update(msg)
			return m, cmd
		}

		// New thread: title input, body, or buttons
		if m.appState == stateNewThread {
			if msg.String() == "tab" {
				switch m.focus {
				case focusTitleInput:
					m.titleInput.Blur()
					m.focus = focusThreadBody
					m.threadBody.Focus()
					m.titleFocused = false
				case focusThreadBody:
					m.threadBody.Blur()
					m.focus = focusButtons
					m.titleFocused = false
				case focusButtons:
					m.focus = focusTitleInput
					m.titleInput.Focus()
					m.titleFocused = true
				}
				return m, nil
			}
			if msg.String() == "esc" {
				m.appState = stateMainMenu
				m.menuCursor = 0
				m.titleInput.Reset()
				m.threadBody.Reset()
				m.focus = focusTitleInput
				m.titleFocused = true
				m.statusMsg = ""
				return m, nil
			}

			if m.focus == focusTitleInput {
				var cmd tea.Cmd
				m.titleInput, cmd = m.titleInput.Update(msg)
				return m, cmd
			}

			if m.focus == focusThreadBody {
				var cmd tea.Cmd
				m.threadBody, cmd = m.threadBody.Update(msg)
				return m, cmd
			}

			if m.focus == focusButtons {
				switch msg.String() {
				case "left", "h":
					if m.activeButton > 0 {
						m.activeButton--
					}
				case "right", "l":
					if m.activeButton < numButtons-1 {
						m.activeButton++
					}
				case "enter", " ":
					return m.activateThreadButton()
				}
				return m, nil
			}
		}

		// New snippet or append entry: editor with buttons
		if m.appState == stateNewSnippet || m.appState == stateAppendEntry {
			if msg.String() == "tab" {
				if m.focus == focusEditor {
					m.focus = focusButtons
					m.textarea.Blur()
				} else {
					m.focus = focusEditor
					m.textarea.Focus()
				}
				return m, nil
			}
			if msg.String() == "esc" && m.appState == stateAppendEntry {
				m.appState = statePickThread
				m.focus = focusList
				m.selectedSlug = ""
				m.statusMsg = ""
				return m, nil
			}

			if m.focus == focusEditor {
				var cmd tea.Cmd
				m.textarea, cmd = m.textarea.Update(msg)
				return m, cmd
			}

			if m.focus == focusButtons {
				switch msg.String() {
				case "left", "h":
					if m.activeButton > 0 {
						m.activeButton--
					}
				case "right", "l":
					if m.activeButton < numButtons-1 {
						m.activeButton++
					}
				case "enter", " ":
					return m.activateButton()
				}
				return m, nil
			}
		}
	}

	return m, nil
}

func (m model) activateThreadButton() (tea.Model, tea.Cmd) {
	switch m.activeButton {
	case btnPreview:
		// Preview not yet supported for thread posts
		m.statusMsg = "Preview not supported for thread posts"
		return m, nil

	case btnSave:
		title := strings.TrimSpace(m.titleInput.Value())
		body := strings.TrimSpace(m.threadBody.Value())

		if title == "" {
			m.statusMsg = "Title cannot be empty."
			return m, nil
		}
		if body == "" {
			m.statusMsg = "Body cannot be empty."
			return m, nil
		}

		path, err := saveThreadPost(title, body)
		if err != nil {
			m.err = err
			return m, nil
		}
		m.saved = true
		m.savedPath = path
		return m, tea.Quit

	case btnQuit:
		m.appState = stateMainMenu
		m.menuCursor = 0
		m.titleInput.Reset()
		m.threadBody.Reset()
		m.focus = focusTitleInput
		m.titleFocused = true
		m.statusMsg = ""
		return m, nil
	}
	return m, nil
}

func (m model) activateButton() (tea.Model, tea.Cmd) {
	switch m.activeButton {
	case btnPreview:
		post := m.textarea.Value()
		if strings.TrimSpace(post) == "" {
			m.statusMsg = "Nothing to preview yet."
			return m, nil
		}
		rendered, err := renderMarkdown(post)
		if err != nil {
			m.err = err
			return m, nil
		}
		m.preview = rendered
		m.mode = modePreview
		m.statusMsg = ""
		return m, nil

	case btnSave:
		post := m.textarea.Value()
		if strings.TrimSpace(post) == "" {
			m.statusMsg = "Can't save an empty post."
			return m, nil
		}

		if m.appState == stateNewSnippet {
			path, err := saveSnippet(post)
			if err != nil {
				m.err = err
				return m, nil
			}
			m.saved = true
			m.savedPath = path
			return m, tea.Quit

		} else if m.appState == stateAppendEntry {
			path, err := appendThreadEntry(m.selectedSlug, post)
			if err != nil {
				m.err = err
				return m, nil
			}
			m.saved = true
			m.savedPath = path
			return m, tea.Quit
		}

	case btnQuit:
		if m.appState == stateNewSnippet || m.appState == stateAppendEntry {
			m.appState = stateMainMenu
			m.menuCursor = 0
			m.textarea.Reset()
			m.selectedSlug = ""
			m.statusMsg = ""
			return m, nil
		}
	}
	return m, nil
}

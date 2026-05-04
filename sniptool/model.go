package main

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// UI dimensions
const (
	minEditorWidth  = 20
	minEditorHeight = 5
	uiPadding       = 4
	uiMargin        = 8
)

// --- application state ---------------------------------------------------

type model struct {
	// Control state
	appState appState
	focus    focusRegion
	ready    bool

	// Home screen
	homeItems  []homeItem
	homeCursor int
	threads    []threadPostMeta
	savedMsg   string // last save message (cleared after render)

	// Editor
	editKind   editKind
	editSlug   string // parent slug for replies
	editTitle  string // title content for thread posts
	titleInput textinput.Model
	textarea   textarea.Model

	// Preview
	preview      string
	viewport     viewport.Model
	activeButton buttonID

	// Status
	width    int
	height   int
	err      error
	lastPath string // path of last saved file
}

func newModel() model {
	ta := textarea.New()
	ta.Placeholder = "Write your content here. Markdown + HTML supported."
	ta.Prompt = "│ "
	ta.ShowLineNumbers = false

	ti := textinput.New()
	ti.Placeholder = "Post title..."
	ti.CharLimit = 100

	m := model{
		appState:     stateHome,
		focus:        focusTitle,
		textarea:     ta,
		titleInput:   ti,
		homeCursor:   0,
		activeButton: btnSave,
		viewport:     viewport.New(0, 0),
	}

	m.loadHomeItems()
	return m
}

func (m *model) resizeComponents() {
	if !m.ready {
		return
	}

	// measure chrome heights (placeholder renders)
	headerH := lipgloss.Height(editorHeaderStyle.Render("x"))
	infoH := lipgloss.Height(infoStyle.Render("x"))
	labelH := 1
	btnH := lipgloss.Height(buttonStyle.Render("x"))
	helpH := lipgloss.Height(editorHelpStyle.Render("x"))
	chrome := headerH + infoH + labelH + btnH + helpH + 2 // 2 spacing

	if m.editKind == editThreadPost {
		chrome += 3 // title label + title input + gap
	}

	boxPad := 4 // border (1 each side) + padding (1 each side)
	w := m.width - boxPad
	h := m.height - chrome - 2 // 2 for box border
	if w < 20 {
		w = 20
	}
	if h < 5 {
		h = 5
	}

	m.textarea.SetWidth(w)
	m.textarea.SetHeight(h)
	m.titleInput.Width = m.width - 2

	m.viewport.Width = m.width - 4
	m.viewport.Height = m.height - 8 // header + buttons + help
}

func (m *model) loadHomeItems() {
	m.homeItems = []homeItem{
		{kind: homeAction, label: "[+] New Snippet Post"},
		{kind: homeAction, label: "[+] New Thread Post"},
	}

	threads, _ := scanThreadPosts()
	m.threads = threads
	for _, t := range threads {
		m.homeItems = append(m.homeItems, homeItem{
			kind:  homeThread,
			label: t.title,
			slug:  t.slug,
		})
	}

	m.homeCursor = 0
}

func (m model) Init() tea.Cmd {
	return nil
}

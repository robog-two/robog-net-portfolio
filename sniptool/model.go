package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

// UI dimensions
const (
	defaultTextareaWidth  = 80
	defaultTextareaHeight = 15
	defaultThreadHeight   = 12
	minTextareaWidth      = 20
	minTextareaHeight     = 5
	uiPadding             = 4
	uiMargin              = 8
)

// --- application state ---------------------------------------------------

type model struct {
	// Control state
	appState appState
	mode     mode
	focus    focusRegion

	// Main menu
	menuCursor int

	// Editor (snippet & thread entry)
	textarea     textarea.Model
	activeButton buttonID
	preview      string

	// Thread creation
	titleInput   textinput.Model
	threadBody   textarea.Model
	titleFocused bool

	// Thread selection
	threadList  list.Model
	threadPosts []threadPostMeta

	// Status
	selectedSlug string
	width        int
	height       int
	statusMsg    string
	err          error
	saved        bool
	savedPath    string
}

func newModel() model {
	ta := textarea.New()
	ta.Placeholder = "Write your post here. Markdown + HTML supported."
	ta.Prompt = "│ "
	ta.ShowLineNumbers = false
	ta.SetWidth(defaultTextareaWidth)
	ta.SetHeight(defaultTextareaHeight)

	ti := textinput.New()
	ti.Placeholder = "Post title..."
	ti.CharLimit = 100
	ti.Width = defaultTextareaWidth

	threadBody := textarea.New()
	threadBody.Placeholder = "Write thread post content here."
	threadBody.Prompt = "│ "
	threadBody.ShowLineNumbers = false
	threadBody.SetWidth(defaultTextareaWidth)
	threadBody.SetHeight(defaultThreadHeight)

	threadList := list.New([]list.Item{}, list.NewDefaultDelegate(), defaultTextareaWidth, defaultTextareaHeight)
	threadList.Title = "Select a thread to append to:"

	return model{
		appState:    stateMainMenu,
		mode:        modeEdit,
		focus:       focusMenu,
		textarea:    ta,
		titleInput:  ti,
		threadBody:  threadBody,
		threadList:  threadList,
		activeButton: btnPreview,
		menuCursor:  0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

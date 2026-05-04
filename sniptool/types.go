package main

// --- state machine -------------------------------------------------------

type appState int

const (
	stateMainMenu appState = iota
	stateNewSnippet
	stateNewThread
	statePickThread
	stateAppendEntry
)

type focusRegion int

const (
	focusEditor focusRegion = iota
	focusButtons
	focusMenu
	focusList
	focusTitleInput
	focusThreadBody
)

type mode int

const (
	modeEdit mode = iota
	modePreview
)

// --- ui controls ---------------------------------------------------------

type buttonID int

const (
	btnPreview buttonID = iota
	btnSave
	btnQuit
	numButtons
)

var buttonLabels = map[buttonID]string{
	btnPreview: "Preview",
	btnSave:    "Save",
	btnQuit:    "Quit",
}

type menuItem int

const (
	menuNewSnippet menuItem = iota
	menuNewThread
	menuAppendThread
	numMenuItems
)

var menuLabels = map[menuItem]string{
	menuNewSnippet:   "New Snippet",
	menuNewThread:    "New Thread",
	menuAppendThread: "Append to Thread",
}

// --- data models ---------------------------------------------------------

type threadPostMeta struct {
	slug  string
	title string
	path  string
}

func (t threadPostMeta) Title() string       { return t.title }
func (t threadPostMeta) Description() string { return "/blog/" + t.slug + "/" }
func (t threadPostMeta) FilterValue() string { return t.title }

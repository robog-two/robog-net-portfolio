package main

// --- state machine -------------------------------------------------------

type appState int

const (
	stateHome appState = iota
	stateEditor
	statePreview
)

type focusRegion int

const (
	focusTitle focusRegion = iota
	focusEditor
	focusButtons
)

// --- buttons (editor and preview screens) --------------------------------

type buttonID int

const (
	btnPreview buttonID = iota
	btnSave
	btnQuit
	numButtons
)

// --- editor context -----------------------------------------------------

type editKind int

const (
	editSnippet editKind = iota
	editThreadPost
	editThreadReply
)

// --- home screen items --------------------------------------------------

type homeItemKind int

const (
	homeAction homeItemKind = iota
	homeThread
)

type homeItem struct {
	kind  homeItemKind
	label string
	slug  string // empty for actions
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

package main

import (
	"strings"
	"time"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/glamour/v2"
)

type mode int

const (
	modeEdit mode = iota
	modePreview
)

// statusLevel colours the status line.
type statusLevel int

const (
	statusInfo statusLevel = iota
	statusWarn
	statusError
)

// statusTTL is how long a status message lingers before clearing itself.
const statusTTL = 4 * time.Second

// statusExpiredMsg clears the status line. The generation guards against a
// stale timer wiping a newer message.
type statusExpiredMsg struct{ gen int }

// Model is the application state.
//
// Errors live in the status line rather than in an `err` field. Previously any
// error — a failed write, a glamour hiccup — made View render nothing but the
// error and an instruction to quit, discarding the unsaved draft. Nothing in
// this program is worth losing a draft over.
type Model struct {
	cfg Config
	// createdAt is captured once at startup and drives both the front matter
	// and the filename, so the name shown in the title bar is the name that
	// gets written.
	createdAt time.Time
	keys      KeyMap
	theme     Theme

	editor   Editor
	preview  viewport.Model
	help     help.Model
	pipeline Pipeline
	store    Store

	mode   mode
	layout Layout

	// savedPath is the file this post owns, empty until the first save. Once
	// claimed, every later save rewrites the same file instead of creating a
	// new one.
	savedPath string
	// savedBody is the body as last written, for dirty tracking.
	savedBody string
	// quitting is set by the first Quit press when there are unsaved changes;
	// a second press confirms.
	quitting bool

	status      string
	statusLevel statusLevel
	statusGen   int

	// renderer is rebuilt when the width or theme changes. Preview word wrap
	// used to be hardcoded at 80 columns regardless of terminal size.
	renderer      *glamour.TermRenderer
	rendererWidth int
	rendererStyle string
}

func NewModel(cfg Config) Model {
	theme := NewTheme(true) // provisional; replaced by BackgroundColorMsg
	return Model{
		cfg:       cfg,
		createdAt: time.Now(),
		keys:      DefaultKeyMap(),
		theme:     theme,
		editor:    newTextareaEditor(theme),
		preview:   viewport.New(),
		help:      help.New(),
		pipeline:  NewPipeline(),
		store:     NewStore(cfg),
		mode:      modeEdit,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.RequestBackgroundColor,
		m.editor.Focus(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		return m.resize(msg.Width, msg.Height)

	case tea.BackgroundColorMsg:
		m.theme = NewTheme(msg.IsDark())
		if ed, ok := m.editor.(*textareaEditor); ok {
			ed.applyTheme(m.theme)
		}
		m.help.Styles = m.theme.Help
		m.renderer = nil // force a rebuild at the new style
		return m, nil

	case statusExpiredMsg:
		if msg.gen == m.statusGen {
			m.status = ""
		}
		return m, nil

	case tea.PasteMsg:
		// Every paste and every drag-and-drop arrives here. This is the whole
		// preprocessing hook: classify, transform, insert.
		return m.handlePaste(string(msg.Content))

	case tea.KeyPressMsg:
		return m.handleKey(msg)
	}

	// Anything else (cursor blink, focus events) goes to the component that
	// owns the current mode. The old code gated this on focus, which starved
	// the editor of its own messages whenever the button row had focus.
	return m.updateActive(msg)
}

func (m Model) resize(w, h int) (tea.Model, tea.Cmd) {
	m.layout = newLayout(w, h)
	if m.layout.TooSmall {
		return m, nil
	}
	m.editor.SetSize(m.layout.ContentWidth, m.layout.ContentHeight)
	m.preview.SetWidth(m.layout.ContentWidth)
	m.preview.SetHeight(m.layout.ContentHeight)

	// Re-render the preview at the new width so a resize does not leave stale
	// wrapping on screen.
	if m.mode == modePreview {
		return m.renderPreview()
	}
	return m, nil
}

func (m Model) handleKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	// Quit is global, and is the only key that needs confirmation.
	if key.Matches(msg, m.keys.Quit) {
		if m.dirty() && !m.quitting {
			m.quitting = true
			return m.setStatus(statusWarn, "unsaved changes — ctrl+c again to discard"), nil
		}
		return m, tea.Quit
	}
	// Any other key cancels a pending quit.
	m.quitting = false

	if key.Matches(msg, m.keys.Save) {
		return m.save()
	}

	switch m.mode {
	case modePreview:
		if key.Matches(msg, m.keys.Back) {
			m.mode = modeEdit
			m.status = ""
			return m, m.editor.Focus()
		}
		var cmd tea.Cmd
		m.preview, cmd = m.preview.Update(msg)
		return m, cmd

	default:
		if key.Matches(msg, m.keys.Preview) {
			return m.enterPreview()
		}
		// Everything else, tab included, belongs to the editor.
		var cmd tea.Cmd
		m.editor, cmd = m.editor.Update(msg)
		return m, cmd
	}
}

func (m Model) updateActive(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if m.mode == modePreview {
		m.preview, cmd = m.preview.Update(msg)
		return m, cmd
	}
	m.editor, cmd = m.editor.Update(msg)
	return m, cmd
}

func (m Model) handlePaste(raw string) (tea.Model, tea.Cmd) {
	if m.mode != modeEdit {
		return m, nil
	}
	res := m.pipeline.Run(raw, Context{Store: m.store, Cfg: m.cfg})
	m.editor.InsertText(res.Text)
	if res.Note != "" {
		return m.setStatus(statusInfo, res.Note), m.expireStatus()
	}
	return m, nil
}

func (m Model) enterPreview() (tea.Model, tea.Cmd) {
	if strings.TrimSpace(m.editor.Value()) == "" {
		return m.setStatus(statusWarn, "nothing to preview yet"), m.expireStatus()
	}
	m.mode = modePreview
	m.editor.Blur()
	m.status = ""
	return m.renderPreview()
}

func (m Model) renderPreview() (tea.Model, tea.Cmd) {
	r, err := m.termRenderer()
	if err != nil {
		m.mode = modeEdit
		return m.setStatus(statusError, "preview unavailable: "+err.Error()), tea.Batch(m.expireStatus(), m.editor.Focus())
	}
	out, err := r.Render(m.editor.Value())
	if err != nil {
		m.mode = modeEdit
		return m.setStatus(statusError, "preview failed: "+err.Error()), tea.Batch(m.expireStatus(), m.editor.Focus())
	}
	m.preview.SetContent(out)
	m.preview.GotoTop()
	return m, nil
}

// termRenderer lazily builds a glamour renderer for the current width and
// theme, reusing it until one of those changes.
func (m *Model) termRenderer() (*glamour.TermRenderer, error) {
	width := m.layout.ContentWidth
	if width <= 0 {
		width = minWidth
	}
	if m.renderer != nil && m.rendererWidth == width && m.rendererStyle == m.theme.Glamour {
		return m.renderer, nil
	}
	r, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle(m.theme.Glamour),
		glamour.WithWordWrap(width),
		glamour.WithEmoji(),
	)
	if err != nil {
		return nil, err
	}
	m.renderer, m.rendererWidth, m.rendererStyle = r, width, m.theme.Glamour
	return r, nil
}

func (m Model) save() (tea.Model, tea.Cmd) {
	body := m.editor.Value()
	if strings.TrimSpace(body) == "" {
		return m.setStatus(statusWarn, "nothing to save"), m.expireStatus()
	}

	post := m.post()
	path, err := Write(m.cfg.Dir(), m.savedPath, post)
	if err != nil {
		// The draft is untouched; the user can fix the problem and retry.
		return m.setStatus(statusError, err.Error()), m.expireStatus()
	}

	m.savedPath = path
	m.savedBody = body
	// Saving deliberately does not quit. It used to, which made an accidental
	// save unrecoverable and forced a relaunch to keep writing.
	return m.setStatus(statusInfo, "saved "+m.displayPath()), m.expireStatus()
}

// post builds the Post for the current buffer. Used both to save and to show
// the pending filename in the title bar, which is why they cannot disagree.
func (m Model) post() Post {
	return Post{
		CreatedAt:   m.createdAt,
		FrontMatter: NewFrontMatter(m.createdAt),
		Body:        m.editor.Value(),
	}
}

func (m Model) dirty() bool {
	return m.editor.Value() != m.savedBody && strings.TrimSpace(m.editor.Value()) != ""
}

func (m Model) setStatus(level statusLevel, text string) Model {
	m.status = text
	m.statusLevel = level
	m.statusGen++
	return m
}

func (m Model) expireStatus() tea.Cmd {
	gen := m.statusGen
	return tea.Tick(statusTTL, func(time.Time) tea.Msg {
		return statusExpiredMsg{gen: gen}
	})
}

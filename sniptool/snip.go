package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

const (
	outputDir = "./src/s"
	blogDir   = "./src/blog"

	// UI dimensions
	defaultTextareaWidth  = 80
	defaultTextareaHeight = 15
	defaultThreadHeight   = 12
	minTextareaWidth      = 20
	minTextareaHeight     = 5
	uiPadding             = 4
	uiMargin              = 8
)

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

// --- styles ---------------------------------------------------------------

var (
	buttonStyle = lipgloss.NewStyle().
			Padding(0, 3).
			Margin(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240"))

	buttonFocusedStyle = buttonStyle.
				BorderForeground(lipgloss.Color("212")).
				Foreground(lipgloss.Color("212")).
				Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1)

	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("99")).
			Bold(true)

	previewBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1)

	menuItemStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Margin(0, 0, 1, 0)

	menuItemFocusedStyle = menuItemStyle.
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("212")).
				Foreground(lipgloss.Color("212")).
				Bold(true)
)

// --- data models ---------------------------------------------------------

type threadPostMeta struct {
	slug  string
	title string
	path  string
}

func (t threadPostMeta) Title() string       { return t.title }
func (t threadPostMeta) Description() string { return "/blog/" + t.slug + "/" }
func (t threadPostMeta) FilterValue() string { return t.title }

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

		// New thread: title input vs body
		if m.appState == stateNewThread {
			if msg.String() == "tab" {
				if m.titleFocused {
					m.titleInput.Blur()
					m.focus = focusThreadBody
					m.threadBody.Focus()
				} else {
					m.threadBody.Blur()
					m.focus = focusTitleInput
					m.titleInput.Focus()
				}
				return m, nil
			}
			if msg.String() == "esc" {
				m.appState = stateMainMenu
				m.menuCursor = 0
				m.titleInput.Reset()
				m.threadBody.Reset()
				m.statusMsg = ""
				return m, nil
			}
			if m.titleFocused {
				var cmd tea.Cmd
				m.titleInput, cmd = m.titleInput.Update(msg)
				return m, cmd
			}
			var cmd tea.Cmd
			m.threadBody, cmd = m.threadBody.Update(msg)
			return m, cmd
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

	if m.titleFocused {
		titleLabel = headerStyle.Render(titleLabel + " (focused)")
		bodyLabel = bodyLabel
	} else {
		titleLabel = titleLabel
		bodyLabel = headerStyle.Render(bodyLabel + " (focused)")
	}

	help := helpStyle.Render("tab: switch fields • esc: cancel • enter in buttons: save")
	var buttons []string
	for id := buttonID(0); id < numButtons; id++ {
		label := buttonLabels[id]
		buttons = append(buttons, buttonStyle.Render(label))
	}
	buttonRow := lipgloss.JoinHorizontal(lipgloss.Top, buttons...)

	return fmt.Sprintf(
		"%s\n%s\n%s\n\n%s\n%s\n\n%s\n\n%s\n%s\n",
		headerStyle.Render("New Thread Post"),
		titleLabel,
		m.titleInput.View(),
		bodyLabel,
		m.threadBody.View(),
		buttonRow,
		help,
		helpStyle.Render("ctrl+c: quit"),
	)
}

// --- main -----------------------------------------------------------------

func main() {
	p := tea.NewProgram(newModel(), tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
	if m, ok := finalModel.(model); ok && m.saved {
		fmt.Printf("Wrote %s\n", m.savedPath)
	}
}

// --- utilities -----------------------------------------------------------

// Markdown rendering
func renderMarkdown(md string) (string, error) {
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)
	if err != nil {
		return "", err
	}
	return r.Render(md)
}

// String formatting
func slugify(s string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9]`)
	stripped := re.ReplaceAllString(s, "")
	if stripped == "" {
		return "post"
	}
	if len(stripped) > 10 {
		stripped = stripped[:10]
	}
	return strings.ToLower(stripped)
}

func formatTime(t time.Time) string {
	hour := t.Hour() % 12
	if hour == 0 {
		hour = 12
	}
	ampm := "AM"
	if t.Hour() >= 12 {
		ampm = "PM"
	}
	return fmt.Sprintf("%d:%02d%s", hour, t.Minute(), ampm)
}

// --- file operations ----------------------------------------------------

func saveSnippet(body string) (string, error) {
	now := time.Now()
	fileName := fmt.Sprintf("%s-%s.md", now.Format("2006-01-02"), slugify(body))
	header := buildSnippetHeader(now)
	path := outputDir + "/" + fileName

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", err
	}
	if err := os.WriteFile(path, []byte(header+body), 0644); err != nil {
		return "", err
	}
	return path, nil
}

func buildSnippetHeader(t time.Time) string {
	return fmt.Sprintf(`---
date: %04d-%02d-%02d
time: %s
tags: post
snippet: yes
layout: blog.njk
description: yes
---
`, t.Year(), int(t.Month()), t.Day(), formatTime(t))
}

// --- thread management --------------------------------------------------

func appendThreadEntry(parentSlug, body string) (string, error) {
	now := time.Now()
	fileName := fmt.Sprintf("%s-%s.md", now.Format("2006-01-02"), slugify(body))
	dirPath := filepath.Join(outputDir, parentSlug)
	path := filepath.Join(dirPath, fileName)

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return "", err
	}

	header := fmt.Sprintf(`---
date: %04d-%02d-%02d
time: %s
thread: %s
layout: thread-entry.njk
---
`, now.Year(), int(now.Month()), now.Day(), formatTime(now), parentSlug)

	if err := os.WriteFile(path, []byte(header+"\n"+body), 0644); err != nil {
		return "", err
	}
	return path, nil
}

func scanThreadPosts() ([]threadPostMeta, error) {
	var posts []threadPostMeta

	entries, err := os.ReadDir(blogDir)
	if err != nil {
		return nil, err
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}

		data, err := os.ReadFile(filepath.Join(blogDir, e.Name()))
		if err != nil {
			continue
		}

		if !strings.Contains(string(data), "thread: true") {
			continue
		}

		slug := strings.TrimSuffix(e.Name(), ".md")
		title := extractFrontmatterTitle(string(data))

		posts = append(posts, threadPostMeta{
			slug:  slug,
			title: title,
			path:  filepath.Join(blogDir, e.Name()),
		})
	}

	return posts, nil
}

func extractFrontmatterTitle(content string) string {
	// Find the frontmatter block
	parts := strings.SplitN(content, "---", 3)
	if len(parts) < 2 {
		return "Unknown"
	}

	// Extract title from frontmatter
	re := regexp.MustCompile(`(?m)^title:\s*(.+)$`)
	matches := re.FindStringSubmatch(parts[1])
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return "Unknown"
}

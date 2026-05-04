package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

const outputDir = "./src/s"

// --- view modes -----------------------------------------------------------

type mode int

const (
	modeEdit mode = iota
	modePreview
)

// --- focus regions --------------------------------------------------------

type focusRegion int

const (
	focusEditor focusRegion = iota
	focusButtons
)

// --- buttons --------------------------------------------------------------

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
)

// --- model ----------------------------------------------------------------

type model struct {
	textarea     textarea.Model
	mode         mode
	focus        focusRegion
	activeButton buttonID
	preview      string
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
	ta.SetWidth(80)
	ta.SetHeight(15)
	ta.Focus()

	return model{
		textarea:     ta,
		mode:         modeEdit,
		focus:        focusEditor,
		activeButton: btnPreview,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Leave room for header + button row + help line.
		taWidth := msg.Width - 4
		if taWidth < 20 {
			taWidth = 20
		}
		taHeight := msg.Height - 8
		if taHeight < 5 {
			taHeight = 5
		}
		m.textarea.SetWidth(taWidth)
		m.textarea.SetHeight(taHeight)
		return m, nil

	case tea.KeyMsg:
		// Global keys.
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "tab":
			// Swap focus between editor and button row.
			if m.mode == modePreview {
				// In preview mode, tab returns to edit mode.
				m.mode = modeEdit
				m.statusMsg = ""
				return m, m.textarea.Focus()
			}
			if m.focus == focusEditor {
				m.focus = focusButtons
				m.textarea.Blur()
			} else {
				m.focus = focusEditor
				return m, m.textarea.Focus()
			}
			return m, nil
		case "esc":
			// In preview mode, esc returns to edit mode.
			if m.mode == modePreview {
				m.mode = modeEdit
				m.statusMsg = ""
				return m, m.textarea.Focus()
			}
		}

		if m.mode == modePreview {
			// Any other key in preview mode: ignore (or could scroll).
			return m, nil
		}

		// Mode is modeEdit. Route by focus.
		switch m.focus {
		case focusEditor:
			var cmd tea.Cmd
			m.textarea, cmd = m.textarea.Update(msg)
			return m, cmd

		case focusButtons:
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

	// Pass everything else to the textarea if it's focused.
	if m.focus == focusEditor && m.mode == modeEdit {
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		return m, cmd
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
		now := time.Now()
		fileName := fmt.Sprintf("%s-%s.md", now.Format("2006-01-02"), slugify(post))
		header := buildHeader(now)
		path := outputDir + "/" + fileName

		if err := os.MkdirAll(outputDir, 0755); err != nil {
			m.err = err
			return m, nil
		}
		if err := os.WriteFile(path, []byte(header+post), 0644); err != nil {
			m.err = err
			return m, nil
		}
		m.saved = true
		m.savedPath = path
		return m, tea.Quit

	case btnQuit:
		return m, tea.Quit
	}
	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress Ctrl+C to quit.\n", m.err)
	}

	if m.mode == modePreview {
		header := headerStyle.Render("── Preview ──")
		help := helpStyle.Render("tab/esc: back to editor • ctrl+c: quit")
		return fmt.Sprintf("%s\n\n%s\n%s\n", header, previewBoxStyle.Render(m.preview), help)
	}

	// Edit mode.
	now := time.Now()
	preview := m.textarea.Value()
	if strings.TrimSpace(preview) == "" {
		preview = "(empty)"
	}
	fileName := fmt.Sprintf("%s-%s.md", now.Format("2006-01-02"), slugify(preview))
	info := fmt.Sprintf("File: %s   Time: %s",
		outputDir+"/"+fileName,
		now.Format("Mon Jan 2 3:04 PM"),
	)

	// Render button row.
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

	return fmt.Sprintf(
		"%s\n%s\n\n%s\n\n%s%s\n%s\n",
		headerStyle.Render("New post"),
		helpStyle.Render(info),
		m.textarea.View(),
		buttonRow,
		status,
		helpStyle.Render(help),
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

// --- helpers (unchanged from previous version) ---------------------------

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

func buildHeader(t time.Time) string {
	hour := t.Hour() % 12
	if hour == 0 {
		hour = 12
	}
	ampm := "AM"
	if t.Hour() >= 12 {
		ampm = "PM"
	}
	timestr := fmt.Sprintf("%d:%02d%s", hour, t.Minute(), ampm)

	return fmt.Sprintf(`---
date: %04d-%02d-%02d
time: %s
tags: post
snippet: yes
layout: blog.njk
description: yes
---
`, t.Year(), int(t.Month()), t.Day(), timestr)
}

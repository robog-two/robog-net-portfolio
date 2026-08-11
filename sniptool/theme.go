package main

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/textarea"
	"charm.land/lipgloss/v2"
)

// Theme holds every style the program uses.
//
// Lip Gloss v2 removed automatic light/dark detection, so the terminal's
// background arrives as a tea.BackgroundColorMsg and the whole theme is
// rebuilt from it. Keeping the styles in one value — rather than as package
// level vars, as before — is what makes that rebuild possible.
//
// No style here sets a margin. Margins on the old help style silently added
// rows the layout arithmetic did not account for, which is what pushed the
// view past the bottom of the screen. Spacing is now the layout's job alone.
type Theme struct {
	Title    lipgloss.Style
	Filename lipgloss.Style
	Muted    lipgloss.Style
	Status   lipgloss.Style
	Warning  lipgloss.Style
	Error    lipgloss.Style
	TooSmall lipgloss.Style

	Textarea textarea.Styles
	Help     help.Styles

	// Glamour is the name of the built-in glamour style used for previews.
	// v2 dropped WithAutoStyle, so the choice is explicit.
	Glamour string
}

func NewTheme(isDark bool) Theme {
	lightDark := lipgloss.LightDark(isDark)

	var (
		accent = lightDark(lipgloss.Color("#8B5CF6"), lipgloss.Color("#C4B5FD"))
		muted  = lightDark(lipgloss.Color("#6B7280"), lipgloss.Color("#9CA3AF"))
		warn   = lightDark(lipgloss.Color("#B45309"), lipgloss.Color("#FBBF24"))
		fail   = lightDark(lipgloss.Color("#B91C1C"), lipgloss.Color("#FCA5A5"))
	)

	glamourStyle := "light"
	if isDark {
		glamourStyle = "dark"
	}

	return Theme{
		Title:    lipgloss.NewStyle().Foreground(accent).Bold(true),
		Filename: lipgloss.NewStyle().Foreground(muted),
		Muted:    lipgloss.NewStyle().Foreground(muted),
		Status:   lipgloss.NewStyle().Foreground(accent),
		Warning:  lipgloss.NewStyle().Foreground(warn),
		Error:    lipgloss.NewStyle().Foreground(fail),
		TooSmall: lipgloss.NewStyle().Foreground(warn),
		Textarea: textarea.DefaultStyles(isDark),
		Help:     help.DefaultStyles(isDark),
		Glamour:  glamourStyle,
	}
}

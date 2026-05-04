package main

import "github.com/charmbracelet/lipgloss"

// --- home screen styles --------------------------------------------------

var (
	homeHeaderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("33")).
			Bold(true)

	actionItemStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Foreground(lipgloss.Color("241"))

	actionItemFocusedStyle = actionItemStyle.
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("33")).
				Foreground(lipgloss.Color("33")).
				Bold(true)

	threadItemStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Foreground(lipgloss.Color("250"))

	threadItemFocusedStyle = threadItemStyle.
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("212")).
				Foreground(lipgloss.Color("212")).
				Bold(true)

	homeHelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1)
)

// --- editor screen styles ------------------------------------------------

var (
	editorHeaderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("33")).
			Bold(true)

	breadcrumbStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("243"))

	editorBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("33")).
			Padding(0, 1)

	titleLabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("33")).
			Bold(true)

	bodyLabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("33")).
			Bold(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1).
			MarginBottom(1)

	buttonStyle = lipgloss.NewStyle().
			Padding(0, 3).
			Margin(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240"))

	buttonFocusedStyle = buttonStyle.
				BorderForeground(lipgloss.Color("33")).
				Foreground(lipgloss.Color("33")).
				Bold(true)

	editorHelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1)
)

// --- preview screen styles -----------------------------------------------

var (
	previewHeaderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("82")).
			Bold(true)

	previewButtonStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Margin(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240"))

	previewButtonFocusedStyle = previewButtonStyle.
					BorderForeground(lipgloss.Color("82")).
					Foreground(lipgloss.Color("82")).
					Bold(true)

	previewHelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1)

	savedBannerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("82")).
			Bold(true).
			Padding(0, 1)
)

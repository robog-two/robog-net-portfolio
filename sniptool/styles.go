package main

import "github.com/charmbracelet/lipgloss"

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

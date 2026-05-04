package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

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

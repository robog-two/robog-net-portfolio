package main

import (
	"flag"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

func main() {
	dir := flag.String("dir", "", "directory to write snippets into (default: <repo root>/"+snippetDir+")")
	flag.Parse()

	if err := run(*dir); err != nil {
		fmt.Fprintln(os.Stderr, "snip:", err)
		os.Exit(1)
	}
}

func run(dir string) error {
	cfg, err := ResolveConfig(dir)
	if err != nil {
		return err
	}

	// AltScreen is no longer a program option in v2; it is a field on the
	// View. See view.go.
	final, err := tea.NewProgram(NewModel(cfg)).Run()
	if err != nil {
		return err
	}

	if m, ok := final.(Model); ok && m.savedPath != "" {
		fmt.Println("Wrote", m.savedPath)
	}
	return nil
}

package main

import "charm.land/bubbles/v2/key"

// KeyMap is the whole command surface of the editor.
//
// This replaces the old focusable button row. Reaching Save used to cost four
// keystrokes — tab out of the editor, right, right, enter — and the focus
// machinery that made it work also stole tab from the editor, so an indented
// code block or a nested list was impossible to type. Direct bindings cost one
// keystroke and give tab back.
//
// Bindings carry their own help text, so the help line is generated from this
// value and cannot drift out of sync with what the keys actually do.
type KeyMap struct {
	Save    key.Binding
	Preview key.Binding
	Back    key.Binding
	Quit    key.Binding

	// Preview-only scrolling, delegated to the viewport.
	Up       key.Binding
	Down     key.Binding
	PageUp   key.Binding
	PageDown key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Save: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "save"),
		),
		Preview: key.NewBinding(
			key.WithKeys("ctrl+p"),
			key.WithHelp("ctrl+p", "preview"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc", "ctrl+p"),
			key.WithHelp("esc", "back"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup", "b"),
			key.WithHelp("pgup", "page up"),
		),
		PageDown: key.NewBinding(
			// v2 stringifies the space bar as "space", not " ".
			key.WithKeys("pgdown", "f", "space"),
			key.WithHelp("pgdn", "page down"),
		),
	}
}

// editHelp and previewHelp implement help.KeyMap for each mode, so the help
// line shows only the keys that currently do something.
type editHelp KeyMap

func (k editHelp) ShortHelp() []key.Binding {
	return []key.Binding{k.Save, k.Preview, k.Quit}
}

func (k editHelp) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

type previewHelp KeyMap

func (k previewHelp) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.PageDown, k.Back, k.Save, k.Quit}
}

func (k previewHelp) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

package main

// The old layout reserved height-8 for the editor while actually drawing
// around eleven rows of chrome — the help style carried a top margin and the
// bordered button row was three rows tall — so the view was taller than the
// screen on every frame. In the alternate screen buffer that means tearing and
// content scrolling off the top. That was the jank.
//
// Rather than a magic number, the rows of chrome are named here and summed.
// Adding a row to the view means adding it to this list; the two cannot drift
// apart without someone noticing.
const (
	rowTitle    = 1 // "New post — 2026-08-06-....md"
	rowGapTop   = 1 // blank line under the title
	rowGapLower = 1 // blank line above the status line
	rowStatus   = 1 // status line, always reserved so nothing reflows
	rowHelp     = 1 // generated help line

	chromeRows = rowTitle + rowGapTop + rowGapLower + rowStatus + rowHelp

	// contentTop is the row the editor and preview start on. It is also the
	// cursor's Y offset.
	contentTop = rowTitle + rowGapTop

	// minWidth and minHeight are the smallest terminal the UI is legible in.
	// Below this the program says so instead of drawing a broken frame.
	minWidth  = 34
	minHeight = chromeRows + 3
)

// Layout is the resolved geometry for one frame.
type Layout struct {
	Width  int
	Height int

	// ContentWidth and ContentHeight are the editor's and preview's exact
	// cell rectangle.
	ContentWidth  int
	ContentHeight int

	// TooSmall reports that the terminal cannot fit the UI.
	TooSmall bool
}

func newLayout(width, height int) Layout {
	l := Layout{Width: width, Height: height}
	if width < minWidth || height < minHeight {
		l.TooSmall = true
		return l
	}
	l.ContentWidth = width
	l.ContentHeight = height - chromeRows
	return l
}

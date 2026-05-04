package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/glamour"
)

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

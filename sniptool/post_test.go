package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestSlugify(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
	}{
		// The cases the old implementation got wrong: it stripped spaces
		// rather than converting them, then cut at 10 bytes.
		{"words become hyphens", "Business cards printed!", "business-cards-printed"},
		{"old bug case", "Made new ringtones for my phone", "made-new-ringtones-for-my-phone"},
		{"leading punctuation", "## A heading", "a-heading"},
		{"collapses runs", "a---b   c", "a-b-c"},
		{"trims edges", "!!hello!!", "hello"},
		{"lowercases", "SHOUTING", "shouting"},
		{"skips blank lines", "\n\n\nfirst real line", "first-real-line"},
		{"only first line", "title here\nbody text", "title-here"},
		{"empty", "", "post"},
		{"whitespace only", "   \n\t ", "post"},
		{"punctuation only", "!!! ???", "post"},
		{"unicode letters kept", "café society", "café-society"},
		{"markdown link", "[phonecall.mp3](/blog/media/x.mp3)", "phonecall-mp3-blog-media-x-mp3"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := Slugify(tc.body); got != tc.want {
				t.Errorf("Slugify(%q) = %q, want %q", tc.body, got, tc.want)
			}
		})
	}
}

func TestSlugifyTruncatesOnWordBoundary(t *testing.T) {
	body := strings.Repeat("alpha ", 20)
	got := Slugify(body)
	if len(got) > maxSlugLen {
		t.Errorf("slug %q is %d bytes, want <= %d", got, len(got), maxSlugLen)
	}
	if strings.HasSuffix(got, "-") || strings.Contains(got, "--") {
		t.Errorf("slug %q has a ragged edge", got)
	}
	// Truncation must not split a word.
	for _, w := range strings.Split(got, "-") {
		if w != "alpha" {
			t.Errorf("slug %q contains partial word %q", got, w)
		}
	}
}

func TestSlugifySingleLongWordIsNotEmpty(t *testing.T) {
	got := Slugify(strings.Repeat("x", 200))
	if got == "" || got == "post" {
		t.Fatalf("got %q, want a hard-truncated slug", got)
	}
	if len(got) > maxSlugLen {
		t.Errorf("slug is %d bytes, want <= %d", len(got), maxSlugLen)
	}
}

func TestFrontMatterMatchesExistingPosts(t *testing.T) {
	// This is the exact header of src/s/2026-05-01-workingonn.md. The site
	// selects snippets on these keys, so drift here breaks the blog.
	when := time.Date(2026, 5, 1, 20, 29, 0, 0, time.UTC)
	want := "---\ndate: 2026-05-01\ntime: 8:29PM\ntags: post\nsnippet: yes\nlayout: blog.liquid\ndescription: yes\n---\n"
	if got := NewFrontMatter(when).Render(); got != want {
		t.Errorf("front matter mismatch\ngot:\n%s\nwant:\n%s", got, want)
	}
}

func TestFrontMatterTwelveHourEdges(t *testing.T) {
	tests := []struct {
		hour, min int
		want      string
	}{
		{0, 0, "12:00AM"},  // midnight, not 0:00AM
		{0, 5, "12:05AM"},  // minute keeps its leading zero
		{12, 0, "12:00PM"}, // noon, not 0:00PM
		{13, 7, "1:07PM"},  // hour drops its leading zero
		{23, 59, "11:59PM"},
	}
	for _, tc := range tests {
		when := time.Date(2026, 1, 2, tc.hour, tc.min, 0, 0, time.UTC)
		if got := NewFrontMatter(when).Time; got != tc.want {
			t.Errorf("%02d:%02d -> %q, want %q", tc.hour, tc.min, got, tc.want)
		}
	}
}

func TestPostRenderEndsWithExactlyOneNewline(t *testing.T) {
	for _, body := range []string{"hello", "hello\n", "hello\n\n\n"} {
		got := Post{FrontMatter: NewFrontMatter(time.Now()), Body: body}.Render()
		if !strings.HasSuffix(got, "hello\n") {
			t.Errorf("body %q rendered to %q", body, got)
		}
		if strings.HasSuffix(got, "\n\n") {
			t.Errorf("body %q left a trailing blank line", body)
		}
	}
}

func TestPostFilename(t *testing.T) {
	p := Post{
		CreatedAt: time.Date(2026, 8, 6, 9, 0, 0, 0, time.UTC),
		Body:      "Business cards printed!",
	}
	if got, want := p.Filename(), "2026-08-06-business-cards-printed.md"; got != want {
		t.Errorf("Filename() = %q, want %q", got, want)
	}
}

func TestWriteAvoidsCollisions(t *testing.T) {
	dir := t.TempDir()
	p := Post{
		CreatedAt:   time.Date(2026, 8, 6, 9, 0, 0, 0, time.UTC),
		FrontMatter: NewFrontMatter(time.Date(2026, 8, 6, 9, 0, 0, 0, time.UTC)),
		Body:        "Same opening words",
	}

	first, err := Write(dir, "", p)
	if err != nil {
		t.Fatal(err)
	}
	second, err := Write(dir, "", p)
	if err != nil {
		t.Fatal(err)
	}
	if first == second {
		t.Fatalf("second save reused %s; it would have overwritten the first post", first)
	}
	if want := "2026-08-06-same-opening-words-2.md"; filepath.Base(second) != want {
		t.Errorf("second path = %s, want basename %s", second, want)
	}
}

func TestWriteReusesClaimedPath(t *testing.T) {
	dir := t.TempDir()
	p := Post{CreatedAt: time.Now(), FrontMatter: NewFrontMatter(time.Now()), Body: "draft"}

	path, err := Write(dir, "", p)
	if err != nil {
		t.Fatal(err)
	}
	p.Body = "draft, revised"
	again, err := Write(dir, path, p)
	if err != nil {
		t.Fatal(err)
	}
	if again != path {
		t.Fatalf("re-save went to %s, want %s", again, path)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("re-saving created %d files, want 1", len(entries))
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(string(data), "draft, revised\n") {
		t.Errorf("file was not updated: %q", data)
	}
}

func TestWriteCreatesMissingDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "s")
	p := Post{CreatedAt: time.Now(), FrontMatter: NewFrontMatter(time.Now()), Body: "hi"}
	if _, err := Write(dir, "", p); err != nil {
		t.Fatal(err)
	}
}

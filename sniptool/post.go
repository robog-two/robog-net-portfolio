package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

// maxSlugLen bounds the slug portion of a filename. Truncation happens on a
// word boundary, so this is an upper bound rather than an exact length.
const maxSlugLen = 40

// FrontMatter is the Eleventy front matter written above every snippet. The
// field set and its order match the posts already in src/s exactly — the site
// selects snippets with `where: 'data.snippet', 'yes'`, so these keys are load
// bearing. Adding a title or real tags later means adding a field here and a
// line in Render, and nothing else.
type FrontMatter struct {
	Date        string
	Time        string
	Tags        string
	Snippet     string
	Layout      string
	Description string
}

func NewFrontMatter(t time.Time) FrontMatter {
	return FrontMatter{
		Date: t.Format("2006-01-02"),
		// "3:04PM" is Go's reference layout for 12-hour time with no leading
		// zero on the hour, which is the format the existing posts use.
		Time:        t.Format("3:04PM"),
		Tags:        "post",
		Snippet:     "yes",
		Layout:      "blog.liquid",
		Description: "yes",
	}
}

func (f FrontMatter) Render() string {
	var b strings.Builder
	b.WriteString("---\n")
	fmt.Fprintf(&b, "date: %s\n", f.Date)
	fmt.Fprintf(&b, "time: %s\n", f.Time)
	fmt.Fprintf(&b, "tags: %s\n", f.Tags)
	fmt.Fprintf(&b, "snippet: %s\n", f.Snippet)
	fmt.Fprintf(&b, "layout: %s\n", f.Layout)
	fmt.Fprintf(&b, "description: %s\n", f.Description)
	b.WriteString("---\n")
	return b.String()
}

// Post is a snippet ready to be written.
type Post struct {
	// CreatedAt is captured once, when the editor session starts, and drives
	// both the front matter and the filename. The old code read the clock
	// separately in View and in the save handler, so the filename shown while
	// typing was not provably the filename written.
	CreatedAt   time.Time
	FrontMatter FrontMatter
	Body        string
}

// Filename is the YYYY-MM-DD-slug.md name for this post.
func (p Post) Filename() string {
	return fmt.Sprintf("%s-%s.md", p.CreatedAt.Format("2006-01-02"), Slugify(p.Body))
}

// Render is the exact bytes that hit disk. The body always ends in exactly one
// newline: POSIX wants the trailing newline, and more than one is noise in the
// rendered excerpt.
func (p Post) Render() string {
	return p.FrontMatter.Render() + strings.TrimRight(p.Body, "\n") + "\n"
}

// Slugify derives a URL-ish slug from the first non-empty line of a post.
//
// The old implementation deleted every non-alphanumeric rune — spaces
// included — and then cut at 10 characters, which is why the existing posts
// are named "businessca" and "madenewrin". This keeps word boundaries as
// hyphens and truncates between words.
func Slugify(body string) string {
	line := firstNonEmptyLine(body)

	var b strings.Builder
	lastHyphen := true // suppresses a leading hyphen
	for _, r := range line {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(unicode.ToLower(r))
			lastHyphen = false
		case !lastHyphen:
			b.WriteByte('-')
			lastHyphen = true
		}
	}
	slug := strings.Trim(b.String(), "-")
	if slug == "" {
		return "post"
	}
	return truncateOnWord(slug, maxSlugLen)
}

func firstNonEmptyLine(body string) string {
	for _, line := range strings.Split(body, "\n") {
		if strings.TrimSpace(line) != "" {
			return line
		}
	}
	return ""
}

// truncateOnWord cuts a hyphenated slug at or before max, preferring a hyphen
// boundary. A single word longer than max is cut hard rather than dropped, so
// the result is never empty.
func truncateOnWord(slug string, max int) string {
	if len(slug) <= max {
		return slug
	}
	cut := slug[:max]
	if i := strings.LastIndexByte(cut, '-'); i > 0 {
		return cut[:i]
	}
	return strings.TrimRight(cut, "-")
}

// Write saves a post, returning the path written.
//
// If path is empty a fresh, non-colliding name is chosen inside dir; otherwise
// path is overwritten. That split is what lets the editor save repeatedly
// without littering the directory: the first save claims a name, and every
// save after that rewrites the same file.
func Write(dir, path string, p Post) (string, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create %s: %w", dir, err)
	}
	if path == "" {
		var err error
		path, err = availablePath(dir, p.Filename())
		if err != nil {
			return "", err
		}
	}
	if err := os.WriteFile(path, []byte(p.Render()), 0o644); err != nil {
		return "", fmt.Errorf("write %s: %w", path, err)
	}
	return path, nil
}

// availablePath finds an unused name, appending -2, -3, … before the
// extension. Two posts on the same day with the same opening words used to
// silently overwrite each other.
func availablePath(dir, name string) (string, error) {
	ext := filepath.Ext(name)
	stem := strings.TrimSuffix(name, ext)

	for i := 1; i < 1000; i++ {
		candidate := filepath.Join(dir, stem+ext)
		if i > 1 {
			candidate = filepath.Join(dir, fmt.Sprintf("%s-%d%s", stem, i, ext))
		}
		_, err := os.Stat(candidate)
		if os.IsNotExist(err) {
			return candidate, nil
		}
		if err != nil {
			return "", fmt.Errorf("stat %s: %w", candidate, err)
		}
	}
	return "", fmt.Errorf("no available filename for %q in %s", name, dir)
}

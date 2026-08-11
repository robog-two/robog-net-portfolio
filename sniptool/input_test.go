package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestSplitShellWords(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []string
	}{
		{"plain", "/a/b.png", []string{"/a/b.png"}},
		{"two files", "/a/b.png /a/c.png", []string{"/a/b.png", "/a/c.png"}},
		// Terminal.app backslash-escapes spaces in dropped paths.
		{"backslash escape", `/a/my\ file.png`, []string{"/a/my file.png"}},
		// iTerm2 and Ghostty single-quote them instead.
		{"single quoted", `'/a/my file.png'`, []string{"/a/my file.png"}},
		{"double quoted", `"/a/my file.png"`, []string{"/a/my file.png"}},
		{"mixed quoting", `'/a/one two.png' /a/three.png`, []string{"/a/one two.png", "/a/three.png"}},
		{"newline separated", "/a/b.png\n/a/c.png", []string{"/a/b.png", "/a/c.png"}},
		{"collapses whitespace", "  /a/b.png   /a/c.png  ", []string{"/a/b.png", "/a/c.png"}},
		{"escaped quote in name", `/a/it\'s.png`, []string{"/a/it's.png"}},
		{"empty", "", nil},
		// An unterminated quote is not well-formed shell quoting, so it is
		// prose, not a drop.
		{"unterminated quote", `'/a/b.png`, nil},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := splitShellWords(tc.in); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("splitShellWords(%q) = %#v, want %#v", tc.in, got, tc.want)
			}
		})
	}
}

// fakeStat reports the given paths as existing and nothing else.
func fakeStat(existing ...string) func(string) (fs.FileInfo, error) {
	set := make(map[string]bool, len(existing))
	for _, p := range existing {
		set[p] = true
	}
	return func(p string) (fs.FileInfo, error) {
		if set[p] {
			return fakeFileInfo{}, nil
		}
		return nil, os.ErrNotExist
	}
}

type fakeFileInfo struct{ fs.FileInfo }

func TestClassifyFileDrop(t *testing.T) {
	c := Classifier{Stat: fakeStat("/tmp/a.png", "/tmp/b c.mp3"), Home: "/home/u"}

	p := c.Classify("/tmp/a.png '/tmp/b c.mp3'")
	if p.Kind != PasteFilePaths {
		t.Fatalf("Kind = %v, want %v", p.Kind, PasteFilePaths)
	}
	if want := []string{"/tmp/a.png", "/tmp/b c.mp3"}; !reflect.DeepEqual(p.Paths, want) {
		t.Errorf("Paths = %#v, want %#v", p.Paths, want)
	}
}

func TestClassifyRejectsProse(t *testing.T) {
	// The stat is deliberately generous: even so, prose must not classify as a
	// drop, because a false positive would mangle ordinary pasted text.
	c := Classifier{Stat: fakeStat("/tmp/a.png"), Home: "/home/u"}
	prose := []string{
		"Business cards printed! So exciting.",
		"README.md",                  // relative, and terminals send absolute
		"/tmp/a.png is a nice image", // path plus words
		"/tmp/does-not-exist.png",    // absolute but missing
		"see /tmp/a.png",
	}
	for _, in := range prose {
		if got := c.Classify(in); got.Kind == PasteFilePaths {
			t.Errorf("Classify(%q) = file paths, want not a drop", in)
		}
	}
}

func TestClassifyFileURIList(t *testing.T) {
	// The XDND text/uri-list convention several terminals use for multi-file
	// drags: newline-separated file:// URIs, with an optional '#' comment.
	c := Classifier{Stat: fakeStat("/tmp/a.jpg", "/tmp/b c.jpg"), Home: "/home/u"}
	raw := "#urilist\nfile:///tmp/a.jpg\nfile:///tmp/b%20c.jpg\n"

	p := c.Classify(raw)
	if p.Kind != PasteFilePaths {
		t.Fatalf("Kind = %v, want %v", p.Kind, PasteFilePaths)
	}
	want := []string{"/tmp/a.jpg", "/tmp/b c.jpg"}
	if !reflect.DeepEqual(p.Paths, want) {
		t.Errorf("Paths = %#v, want %#v", p.Paths, want)
	}
}

func TestClassifyExpandsHome(t *testing.T) {
	c := Classifier{Stat: fakeStat("/home/u/pic.png"), Home: "/home/u"}
	p := c.Classify("~/pic.png")
	if p.Kind != PasteFilePaths {
		t.Fatalf("Kind = %v, want file paths", p.Kind)
	}
	if want := []string{"/home/u/pic.png"}; !reflect.DeepEqual(p.Paths, want) {
		t.Errorf("Paths = %#v, want %#v", p.Paths, want)
	}
}

func TestClassifyURLs(t *testing.T) {
	c := Classifier{Stat: fakeStat(), Home: "/home/u"}

	p := c.Classify("https://localify.org/ http://example.com/x")
	if p.Kind != PasteURLs {
		t.Fatalf("Kind = %v, want %v", p.Kind, PasteURLs)
	}
	if len(p.URLs) != 2 {
		t.Errorf("URLs = %#v, want 2", p.URLs)
	}

	for _, in := range []string{
		"check out https://localify.org/", // URL plus words
		"mailto:someone@example.com",      // wrong scheme
		"just some text",
	} {
		if got := c.Classify(in); got.Kind == PasteURLs {
			t.Errorf("Classify(%q) = urls, want plain", in)
		}
	}
}

func TestClassifyNormalisesLineEndings(t *testing.T) {
	c := Classifier{Stat: fakeStat(), Home: ""}
	p := c.Classify("one\r\ntwo\rthree")
	if p.Raw != "one\ntwo\nthree" {
		t.Errorf("Raw = %q, want normalised newlines", p.Raw)
	}
	if p.Kind != PastePlain {
		t.Errorf("Kind = %v, want plain", p.Kind)
	}
}

func TestPipelinePlainPasteIsVerbatim(t *testing.T) {
	p := NewPipeline()
	const text = "Business cards printed!\n\nFor those viewing remotely:"
	res := p.Run(text, Context{Store: disabledStore{}})
	if res.Text != text {
		t.Errorf("plain paste was altered:\ngot  %q\nwant %q", res.Text, text)
	}
	if res.Note != "" {
		t.Errorf("plain paste produced note %q, want none", res.Note)
	}
}

func TestPipelineFileDropWithDisabledStore(t *testing.T) {
	// With the default Store, a drop still inserts something useful and
	// explains itself rather than failing silently.
	dir := t.TempDir()
	a := filepath.Join(dir, "toot.mp3")
	if err := os.WriteFile(a, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	p := NewPipeline()
	res := p.Run(a, Context{Store: disabledStore{}})
	if res.Text != a {
		t.Errorf("Text = %q, want the dropped path %q", res.Text, a)
	}
	if !strings.Contains(res.Note, "not enabled") {
		t.Errorf("Note = %q, want an explanation", res.Note)
	}
}

// stubStore stands in for a real media store, to prove the transformer
// produces correct Markdown once a Store exists.
type stubStore struct{}

func (stubStore) Put(src string) (Attachment, error) {
	return Attachment{Source: src, Ref: "/blog/media/" + filepath.Base(src)}, nil
}

func TestPipelineFileDropWithWorkingStore(t *testing.T) {
	dir := t.TempDir()
	audio := filepath.Join(dir, "toot.mp3")
	image := filepath.Join(dir, "shot.png")
	for _, f := range []string{audio, image} {
		if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	p := NewPipeline()
	res := p.Run(audio+" "+image, Context{Store: stubStore{}})

	want := "[toot.mp3](/blog/media/toot.mp3)\n![shot.png](/blog/media/shot.png)"
	if res.Text != want {
		t.Errorf("Text =\n%q\nwant\n%q", res.Text, want)
	}
	if res.Note != "attached 2 files" {
		t.Errorf("Note = %q", res.Note)
	}
}

func TestPipelineURLPasteBecomesLink(t *testing.T) {
	p := NewPipeline()
	res := p.Run("https://example.com/post", Context{Store: disabledStore{}})
	want := "[https://example.com/post](https://example.com/post)"
	if res.Text != want {
		t.Errorf("Text = %q, want %q", res.Text, want)
	}
}

func TestPipelineImageURLPasteBecomesEmbed(t *testing.T) {
	p := NewPipeline()
	res := p.Run("https://example.com/pic.png", Context{Store: disabledStore{}})
	want := "![https://example.com/pic.png](https://example.com/pic.png)"
	if res.Text != want {
		t.Errorf("Text = %q, want %q", res.Text, want)
	}
}

func TestPipelineMultipleURLsOnePerLine(t *testing.T) {
	p := NewPipeline()
	res := p.Run("https://example.com/a https://example.com/b.jpg", Context{Store: disabledStore{}})
	want := "[https://example.com/a](https://example.com/a)\n" +
		"![https://example.com/b.jpg](https://example.com/b.jpg)"
	if res.Text != want {
		t.Errorf("Text =\n%q\nwant\n%q", res.Text, want)
	}
}

func TestPipelineFileDropWithRealFileStore(t *testing.T) {
	root := t.TempDir()
	cfg := Config{RepoRoot: root, MediaDir: filepath.Join(root, mediaDir)}
	store := FileStore{cfg: cfg}

	dir := t.TempDir()
	image := filepath.Join(dir, "shot.png")
	if err := os.WriteFile(image, []byte("pixels"), 0o644); err != nil {
		t.Fatal(err)
	}

	p := NewPipeline()
	res := p.Run(image, Context{Store: store, Cfg: cfg})
	if want := "![shot.png](/blog/media/shot.png)"; res.Text != want {
		t.Errorf("Text = %q, want %q", res.Text, want)
	}
	if _, err := os.Stat(filepath.Join(cfg.MediaDir, "shot.png")); err != nil {
		t.Errorf("file was not copied into MediaDir: %v", err)
	}
}

// TestPipelineMultiFileDropAsURIList reproduces a multi-file Finder drag
// delivered as a text/uri-list — the shape that used to be misclassified as
// prose and inserted verbatim instead of becoming attachments.
func TestPipelineMultiFileDropAsURIList(t *testing.T) {
	root := t.TempDir()
	cfg := Config{RepoRoot: root, MediaDir: filepath.Join(root, mediaDir)}
	store := FileStore{cfg: cfg}

	dir := t.TempDir()
	a := filepath.Join(dir, "one.jpg")
	b := filepath.Join(dir, "two.jpg")
	for _, f := range []string{a, b} {
		if err := os.WriteFile(f, []byte(f), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	raw := "file://" + a + "\nfile://" + b + "\n"
	p := NewPipeline()
	res := p.Run(raw, Context{Store: store, Cfg: cfg})

	want := "![one.jpg](/blog/media/one.jpg)\n![two.jpg](/blog/media/two.jpg)"
	if res.Text != want {
		t.Errorf("Text =\n%q\nwant\n%q", res.Text, want)
	}
	for _, name := range []string{"one.jpg", "two.jpg"} {
		if _, err := os.Stat(filepath.Join(cfg.MediaDir, name)); err != nil {
			t.Errorf("%s was not copied into MediaDir: %v", name, err)
		}
	}
}

func TestKindForPath(t *testing.T) {
	tests := map[string]AttachmentKind{
		"/a/b.PNG":  AttachmentImage,
		"/a/b.jpeg": AttachmentImage,
		"/a/b.mp3":  AttachmentAudio,
		"/a/b.pdf":  AttachmentFile,
		"/a/b":      AttachmentFile,
	}
	for path, want := range tests {
		if got := kindForPath(path); got != want {
			t.Errorf("kindForPath(%q) = %v, want %v", path, got, want)
		}
	}
}

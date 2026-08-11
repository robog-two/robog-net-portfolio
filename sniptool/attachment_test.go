package main

import (
	"os"
	"path/filepath"
	"testing"
)

func testConfig(t *testing.T) Config {
	t.Helper()
	root := t.TempDir()
	return Config{RepoRoot: root, MediaDir: filepath.Join(root, mediaDir)}
}

func TestFileStorePutCopiesAndRefsFile(t *testing.T) {
	cfg := testConfig(t)
	srcDir := t.TempDir()
	src := filepath.Join(srcDir, "shot.png")
	if err := os.WriteFile(src, []byte("pixels"), 0o644); err != nil {
		t.Fatal(err)
	}

	att, err := FileStore{cfg: cfg}.Put(src)
	if err != nil {
		t.Fatal(err)
	}
	if want := "/blog/media/shot.png"; att.Ref != want {
		t.Errorf("Ref = %q, want %q", att.Ref, want)
	}
	dest := filepath.Join(cfg.MediaDir, "shot.png")
	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "pixels" {
		t.Errorf("copied content = %q, want %q", got, "pixels")
	}
}

func TestFileStorePutSanitizesName(t *testing.T) {
	cfg := testConfig(t)
	srcDir := t.TempDir()
	src := filepath.Join(srcDir, "My Screenshot (1).PNG")
	if err := os.WriteFile(src, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	att, err := FileStore{cfg: cfg}.Put(src)
	if err != nil {
		t.Fatal(err)
	}
	if want := "/blog/media/my-screenshot-1.png"; att.Ref != want {
		t.Errorf("Ref = %q, want %q", att.Ref, want)
	}
}

func TestFileStorePutDedupesByContent(t *testing.T) {
	cfg := testConfig(t)
	srcDir := t.TempDir()
	a := filepath.Join(srcDir, "a.png")
	b := filepath.Join(srcDir, "b.png") // different name, same bytes
	for _, p := range []string{a, b} {
		if err := os.WriteFile(p, []byte("same bytes"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	store := FileStore{cfg: cfg}
	first, err := store.Put(a)
	if err != nil {
		t.Fatal(err)
	}
	second, err := store.Put(b)
	if err != nil {
		t.Fatal(err)
	}
	if first.Ref != second.Ref {
		t.Errorf("dropping identical content twice produced two refs: %q and %q", first.Ref, second.Ref)
	}

	entries, err := os.ReadDir(cfg.MediaDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Errorf("MediaDir has %d entries, want 1 (deduped)", len(entries))
	}
}

func TestFileStorePutAvoidsNameCollisionForDifferentContent(t *testing.T) {
	cfg := testConfig(t)
	srcDir := t.TempDir()
	a := filepath.Join(srcDir, "sub", "shot.png")
	b := filepath.Join(srcDir, "shot.png")
	if err := os.MkdirAll(filepath.Dir(a), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(a, []byte("one"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(b, []byte("two"), 0o644); err != nil {
		t.Fatal(err)
	}

	store := FileStore{cfg: cfg}
	first, err := store.Put(a)
	if err != nil {
		t.Fatal(err)
	}
	second, err := store.Put(b)
	if err != nil {
		t.Fatal(err)
	}
	if first.Ref == second.Ref {
		t.Errorf("two different files collapsed to one ref %q", first.Ref)
	}
}

func TestNewStoreFallsBackWhenNoRepoRoot(t *testing.T) {
	if _, ok := NewStore(Config{OutputDir: "/tmp/wherever"}).(disabledStore); !ok {
		t.Error("NewStore with no MediaDir should return disabledStore")
	}
	cfg := testConfig(t)
	if _, ok := NewStore(cfg).(FileStore); !ok {
		t.Error("NewStore with a resolved repo root should return FileStore")
	}
}

func TestConfigRefRejectsPathOutsideSrc(t *testing.T) {
	cfg := testConfig(t)
	if _, err := cfg.Ref(filepath.Join(cfg.RepoRoot, "..", "elsewhere.png")); err == nil {
		t.Error("expected an error for a path outside RepoRoot/src")
	}
}

func TestAttachmentMarkdown(t *testing.T) {
	img := Attachment{Ref: "/blog/media/shot.png", Kind: AttachmentImage}
	if got, want := img.Markdown("a shot"), "![a shot](/blog/media/shot.png)"; got != want {
		t.Errorf("Markdown() = %q, want %q", got, want)
	}
	file := Attachment{Ref: "/blog/media/notes.pdf", Kind: AttachmentFile}
	if got, want := file.Markdown(""), "[/blog/media/notes.pdf](/blog/media/notes.pdf)"; got != want {
		t.Errorf("Markdown() with empty label = %q, want %q", got, want)
	}
}

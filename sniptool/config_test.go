package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveConfigFindsRepoRoot(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "deno.json"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	deep := filepath.Join(root, "src", "blog", "media")
	if err := os.MkdirAll(deep, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Chdir(deep)

	cfg, err := ResolveConfig("")
	if err != nil {
		t.Fatal(err)
	}
	// t.TempDir can sit under a symlink (/var -> /private/var on macOS), so
	// compare resolved paths.
	want, _ := filepath.EvalSymlinks(filepath.Join(root, snippetDir))
	got, _ := filepath.EvalSymlinks(cfg.OutputDir)
	if got != want {
		t.Errorf("OutputDir = %s, want %s", got, want)
	}

	wantMedia, _ := filepath.EvalSymlinks(filepath.Join(root, mediaDir))
	gotMedia, _ := filepath.EvalSymlinks(cfg.MediaDir)
	if gotMedia != wantMedia {
		t.Errorf("MediaDir = %s, want %s", gotMedia, wantMedia)
	}
}

// TestResolveConfigHonoursFlagStillFindsMediaDir: -dir only redirects where
// snippets are written; attachments still belong under the repo's
// src/blog/media, so MediaDir must still resolve.
func TestResolveConfigHonoursFlagStillFindsMediaDir(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "deno.json"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Chdir(root)

	cfg, err := ResolveConfig("elsewhere")
	if err != nil {
		t.Fatal(err)
	}
	want, _ := filepath.EvalSymlinks(filepath.Join(root, mediaDir))
	got, _ := filepath.EvalSymlinks(cfg.MediaDir)
	if got != want {
		t.Errorf("MediaDir = %s, want %s", got, want)
	}
}

func TestConfigRefOmitsSrc(t *testing.T) {
	cfg := Config{RepoRoot: "/repo"}
	got, err := cfg.Ref("/repo/src/blog/media/toot.mp3")
	if err != nil {
		t.Fatal(err)
	}
	if want := "/blog/media/toot.mp3"; got != want {
		t.Errorf("Ref() = %q, want %q", got, want)
	}
}

// TestResolveConfigFailsOutsideRepo: the old tool silently created a stray
// ./src/s wherever it happened to be run.
func TestResolveConfigFailsOutsideRepo(t *testing.T) {
	t.Chdir(t.TempDir())
	if _, err := ResolveConfig(""); err == nil {
		t.Error("expected an error outside the repository, got none")
	}
}

func TestResolveConfigHonoursFlag(t *testing.T) {
	t.Chdir(t.TempDir())
	cfg, err := ResolveConfig("some/where")
	if err != nil {
		t.Fatal(err)
	}
	if !filepath.IsAbs(cfg.OutputDir) {
		t.Errorf("OutputDir = %s, want an absolute path", cfg.OutputDir)
	}
}

func TestConfigSubdirNestsPosts(t *testing.T) {
	// Threads already live on disk as src/s/<thread>/; this is the hook.
	cfg := Config{OutputDir: "/repo/src/s", Subdir: "digital-minimalism"}
	if got, want := cfg.Dir(), "/repo/src/s/digital-minimalism"; got != want {
		t.Errorf("Dir() = %s, want %s", got, want)
	}
}

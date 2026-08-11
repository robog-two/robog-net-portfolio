package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// snippetDir is where snippets live, relative to the repository root.
const snippetDir = "src/s"

// mediaDir is where dropped attachments are copied to, relative to the
// repository root. Every post-facing Ref is this directory's contents with
// "src" trimmed off the front — see Config.Ref.
const mediaDir = "src/blog/media"

// repoMarkers identify the repository root when walking up from the working
// directory. deno.json is checked first because it is what actually defines
// this site's build.
var repoMarkers = []string{"deno.json", ".git"}

// Config is everything the program needs from the outside world.
type Config struct {
	// OutputDir is the resolved, absolute directory snippets are written to.
	OutputDir string

	// RepoRoot and MediaDir are absolute paths, set only when a repository
	// root was actually found. -dir lets OutputDir be pointed anywhere, so a
	// caller outside a repo can still write snippets; MediaDir being empty is
	// how NewStore knows attachments have nowhere well-defined to land, and
	// falls back to disabledStore.
	RepoRoot string
	MediaDir string

	// Subdir, when set, nests the post inside OutputDir — the shape threads
	// already take on disk (src/s/digital-minimalism/). Nothing sets this
	// yet; it is here so thread support is a flag and a path join rather than
	// a change to how posts are written.
	Subdir string
}

// ResolveConfig works out where to write.
//
// The old tool hardcoded a relative "./src/s" while build.sh installed the
// binary a directory above the repo root, so running it from anywhere else
// silently created a stray src/s in the wrong place. This walks up to the
// repository root instead, and fails loudly rather than guessing.
func ResolveConfig(dirFlag string) (Config, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return Config{}, fmt.Errorf("get working directory: %w", err)
	}
	// The repo root drives MediaDir regardless of -dir, because attachments
	// always belong under the site's src/blog/media, not wherever -dir points
	// snippets at. Only OutputDir resolution is allowed to fail over to the
	// flag; a missing repo root otherwise still fails loudly, as before.
	root, rootErr := findRepoRoot(cwd)

	var cfg Config
	if dirFlag != "" {
		abs, err := filepath.Abs(dirFlag)
		if err != nil {
			return Config{}, fmt.Errorf("resolve -dir: %w", err)
		}
		cfg.OutputDir = abs
	} else {
		if rootErr != nil {
			return Config{}, rootErr
		}
		cfg.OutputDir = filepath.Join(root, snippetDir)
	}

	if rootErr == nil {
		cfg.RepoRoot = root
		cfg.MediaDir = filepath.Join(root, mediaDir)
	}
	return cfg, nil
}

// Dir is the directory a post is written to.
func (c Config) Dir() string {
	if c.Subdir == "" {
		return c.OutputDir
	}
	return filepath.Join(c.OutputDir, c.Subdir)
}

// Ref maps an absolute path under RepoRoot/src to the site-relative URL a
// post should link to, e.g. .../src/blog/media/toot.mp3 becomes
// "/blog/media/toot.mp3". The site serves everything under src/ from "/", so
// omitting that one path segment is the entire mapping.
func (c Config) Ref(absPath string) (string, error) {
	rel, err := filepath.Rel(filepath.Join(c.RepoRoot, "src"), absPath)
	if err != nil {
		return "", fmt.Errorf("map %s to a site path: %w", absPath, err)
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("%s is outside %s", absPath, filepath.Join(c.RepoRoot, "src"))
	}
	return "/" + filepath.ToSlash(rel), nil
}

func findRepoRoot(start string) (string, error) {
	dir := start
	for {
		for _, marker := range repoMarkers {
			if _, err := os.Stat(filepath.Join(dir, marker)); err == nil {
				return dir, nil
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New(
				"not inside the site repository (no deno.json or .git found above " +
					start + "); pass -dir to write somewhere explicit")
		}
		dir = parent
	}
}

package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

// ErrAttachmentsDisabled is returned by the default Store. Attachments are
// detected (see input.go) but not yet ingested; this is the one thing standing
// between detection and a working feature.
var ErrAttachmentsDisabled = errors.New("attachments are not enabled")

// Attachment is a file that has been pulled into the site alongside a post.
type Attachment struct {
	// Source is where the file came from on the local disk.
	Source string
	// Ref is the site-relative URL the post should link to, e.g.
	// "/blog/media/ringtones/toot.mp3".
	Ref string
	// Kind drives how the attachment is written into the Markdown: an image
	// gets ![](), everything else gets a plain link.
	Kind AttachmentKind
}

type AttachmentKind int

const (
	AttachmentFile AttachmentKind = iota
	AttachmentImage
	AttachmentAudio
)

// Store ingests a local file into the site and hands back the reference a post
// should use.
type Store interface {
	Put(sourcePath string) (Attachment, error)
}

// disabledStore is the Store used when no repository root was found (see
// Config.MediaDir). It refuses everything, which makes the file-drop
// transformer fall through to inserting the raw paths with a note explaining
// why.
type disabledStore struct{}

func (disabledStore) Put(string) (Attachment, error) {
	return Attachment{}, ErrAttachmentsDisabled
}

// NewStore picks the Store a Config can actually support: a working FileStore
// when a repository root was resolved, disabledStore otherwise.
func NewStore(cfg Config) Store {
	if cfg.MediaDir == "" || cfg.RepoRoot == "" {
		return disabledStore{}
	}
	return FileStore{cfg: cfg}
}

// FileStore copies files into Config.MediaDir and maps them to the site URL
// the published post should use.
//
// Dedupe is by content hash rather than by name: dropping the same
// screenshot into two different posts should reuse one file on disk, not
// create screenshot.png and screenshot-2.png that happen to be identical.
// Hashing is done by walking MediaDir on every Put rather than keeping an
// index — this tool's media directories are small enough (dozens of files)
// that the simplicity is worth more than the saved stat calls, and a
// persistent index would need its own invalidation story for files added
// outside the tool.
type FileStore struct {
	cfg Config
}

func (s FileStore) Put(sourcePath string) (Attachment, error) {
	data, err := os.ReadFile(sourcePath)
	if err != nil {
		return Attachment{}, fmt.Errorf("read %s: %w", sourcePath, err)
	}
	hash := sha256.Sum256(data)

	if existing, err := findByHash(s.cfg.MediaDir, hash); err != nil {
		return Attachment{}, err
	} else if existing != "" {
		ref, err := s.cfg.Ref(existing)
		if err != nil {
			return Attachment{}, err
		}
		return Attachment{Source: sourcePath, Ref: ref}, nil
	}

	if err := os.MkdirAll(s.cfg.MediaDir, 0o755); err != nil {
		return Attachment{}, fmt.Errorf("create %s: %w", s.cfg.MediaDir, err)
	}
	dest, err := availablePath(s.cfg.MediaDir, sanitizeFilename(filepath.Base(sourcePath)))
	if err != nil {
		return Attachment{}, err
	}
	if err := os.WriteFile(dest, data, 0o644); err != nil {
		return Attachment{}, fmt.Errorf("write %s: %w", dest, err)
	}

	ref, err := s.cfg.Ref(dest)
	if err != nil {
		return Attachment{}, err
	}
	return Attachment{Source: sourcePath, Ref: ref}, nil
}

// findByHash walks dir looking for a file whose content already matches
// hash, returning its path or "" if dir doesn't exist yet or nothing
// matches.
func findByHash(dir string, hash [sha256.Size]byte) (string, error) {
	if _, err := os.Stat(dir); errors.Is(err, fs.ErrNotExist) {
		return "", nil
	}

	var found string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || found != "" || d.IsDir() {
			return err
		}
		h, err := hashFile(path)
		if err != nil {
			return err
		}
		if h == hash {
			found = path
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("scan %s for duplicates: %w", dir, err)
	}
	return found, nil
}

func hashFile(path string) ([sha256.Size]byte, error) {
	var zero [sha256.Size]byte
	f, err := os.Open(path)
	if err != nil {
		return zero, err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return zero, err
	}
	var out [sha256.Size]byte
	copy(out[:], h.Sum(nil))
	return out, nil
}

// sanitizeFilename makes a dropped file's basename safe to publish as a URL
// path segment: lowercased, whitespace and anything else odd collapsed to a
// single hyphen. The extension is kept verbatim (lowercased) so kindForPath
// keeps working on the result.
func sanitizeFilename(name string) string {
	ext := strings.ToLower(filepath.Ext(name))
	base := strings.TrimSuffix(name, filepath.Ext(name))

	var b strings.Builder
	lastHyphen := true
	for _, r := range base {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(unicode.ToLower(r))
			lastHyphen = false
		case r == '-' || r == '_':
			b.WriteRune(r)
			lastHyphen = false
		case !lastHyphen:
			b.WriteByte('-')
			lastHyphen = true
		}
	}
	sanitized := strings.Trim(b.String(), "-")
	if sanitized == "" {
		sanitized = "file"
	}
	return sanitized + ext
}

// Markdown renders an attachment as the snippet source that should be
// inserted. Written now because it is pure, trivial, and the exact thing a
// Store implementation will need on day one.
func (a Attachment) Markdown(label string) string {
	if label == "" {
		label = a.Ref
	}
	if a.Kind == AttachmentImage {
		return "![" + label + "](" + a.Ref + ")"
	}
	return "[" + label + "](" + a.Ref + ")"
}

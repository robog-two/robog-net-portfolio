package main

import (
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Everything the user types arrives as a key press, but everything the user
// *pastes or drops* arrives as one tea.PasteMsg holding a whole blob of text.
// That blob is the only hook a terminal gives us: a file dragged onto the
// window is delivered as bracketed paste containing its shell-quoted path.
//
// So this file is the preprocessing layer. Classify decides what a blob
// appears to be, and a Pipeline of Transformers decides what to insert. The
// model calls Pipeline.Run and inserts the result; it knows nothing about
// paths, URLs or shell quoting.

// PasteKind is what a pasted blob appears to be.
type PasteKind int

const (
	// PastePlain is ordinary text. The overwhelmingly common case.
	PastePlain PasteKind = iota
	// PasteFilePaths is a drag-and-drop: one or more shell-quoted paths that
	// all exist on disk.
	PasteFilePaths
	// PasteURLs is one or more http(s) URLs and nothing else.
	PasteURLs
)

func (k PasteKind) String() string {
	switch k {
	case PasteFilePaths:
		return "file paths"
	case PasteURLs:
		return "urls"
	default:
		return "plain"
	}
}

// Paste is a classified blob of pasted text.
type Paste struct {
	// Raw is the blob with line endings normalised, and nothing else changed.
	Raw   string
	Kind  PasteKind
	Paths []string // set when Kind is PasteFilePaths
	URLs  []string // set when Kind is PasteURLs
}

// Classifier turns raw paste blobs into Pastes. Stat is injectable so the
// classification rules can be tested without touching the filesystem.
type Classifier struct {
	Stat func(string) (fs.FileInfo, error)
	// Home expands a leading ~. Empty means no expansion.
	Home string
}

func NewClassifier() Classifier {
	home, _ := os.UserHomeDir()
	return Classifier{Stat: os.Stat, Home: home}
}

func (c Classifier) Classify(raw string) Paste {
	raw = strings.ReplaceAll(raw, "\r\n", "\n")
	raw = strings.ReplaceAll(raw, "\r", "\n")
	p := Paste{Raw: raw, Kind: PastePlain}

	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return p
	}

	if paths, ok := c.asFilePaths(trimmed); ok {
		p.Kind, p.Paths = PasteFilePaths, paths
		return p
	}
	if urls, ok := asURLs(trimmed); ok {
		p.Kind, p.URLs = PasteURLs, urls
		return p
	}
	return p
}

// asFilePaths reports whether the blob is a file drop. The test is
// deliberately strict — every shell-word must resolve to something that
// actually exists — because misclassifying prose as a drop would be far worse
// than missing a drop. Existence is the strongest signal available and costs
// one stat per word.
//
// Single-file drags and multi-file drags aren't guaranteed to arrive in the
// same shape: several terminals switch from plain shell-quoted paths to the
// XDND/text-uri-list convention (newline-separated file:// URIs, optionally
// with '#' comment lines) once more than one file is selected. Both shapes
// are handled here so a two-file drop isn't silently treated as prose.
func (c Classifier) asFilePaths(s string) ([]string, bool) {
	words := splitShellWords(s)
	if len(words) == 0 {
		return nil, false
	}
	stat := c.Stat
	if stat == nil {
		stat = os.Stat
	}
	out := make([]string, 0, len(words))
	for _, w := range words {
		if strings.HasPrefix(w, "#") {
			continue // text/uri-list comment line
		}
		raw := w
		if p, ok := fileURIToPath(w); ok {
			raw = p
		}
		abs := c.expand(raw)
		// A relative word is never a drop: terminals always deliver absolute
		// paths, and treating relative words as paths would misfire on
		// ordinary pasted prose like "README.md".
		if !filepath.IsAbs(abs) {
			return nil, false
		}
		if _, err := stat(abs); err != nil {
			return nil, false
		}
		out = append(out, abs)
	}
	if len(out) == 0 {
		return nil, false
	}
	return out, true
}

func (c Classifier) expand(p string) string {
	if c.Home != "" && (p == "~" || strings.HasPrefix(p, "~/")) {
		return filepath.Join(c.Home, strings.TrimPrefix(strings.TrimPrefix(p, "~"), "/"))
	}
	return p
}

// fileURIToPath decodes a file:// URI (the text/uri-list form multi-file
// drags commonly arrive as) into a local filesystem path. It reports false
// for anything that isn't a file:// URI, including ordinary http(s) URLs and
// plain paths, so callers can try it opportunistically without disturbing
// either.
func fileURIToPath(s string) (string, bool) {
	u, err := url.Parse(s)
	if err != nil || u.Scheme != "file" || u.Path == "" {
		return "", false
	}
	path := u.Path
	if u.Host != "" && u.Host != "localhost" {
		// A remote host segment (file://host/path) — keep it well-formed
		// rather than silently dropping which machine it names.
		path = "/" + u.Host + path
	}
	return path, true
}

func asURLs(s string) ([]string, bool) {
	fields := strings.Fields(s)
	if len(fields) == 0 {
		return nil, false
	}
	out := make([]string, 0, len(fields))
	for _, f := range fields {
		u, err := url.Parse(f)
		if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") {
			return nil, false
		}
		out = append(out, f)
	}
	return out, true
}

// splitShellWords splits a drop payload the way the terminal wrote it.
// Terminals differ: Terminal.app backslash-escapes spaces, iTerm2 and Ghostty
// single-quote paths containing them, and multiple files may be separated by
// spaces or newlines. Handling all three is what makes multi-file drops of
// awkwardly named files work.
func splitShellWords(s string) []string {
	var (
		words []string
		cur   strings.Builder
		quote rune // 0, '\'' or '"'
		open  bool // cur holds something, even if empty (e.g. "")
	)
	flush := func() {
		if open {
			words = append(words, cur.String())
			cur.Reset()
			open = false
		}
	}

	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		switch {
		case quote == 0 && r == '\\':
			// Backslash escapes the next rune, including a newline.
			if i+1 < len(runes) {
				i++
				cur.WriteRune(runes[i])
				open = true
			}
		case quote == 0 && (r == '\'' || r == '"'):
			quote = r
			open = true
		case quote != 0 && r == quote:
			quote = 0
		case quote == 0 && (r == ' ' || r == '\t' || r == '\n'):
			flush()
		default:
			cur.WriteRune(r)
			open = true
		}
	}
	// An unterminated quote means this was never well-formed shell quoting,
	// so it is not a drop.
	if quote != 0 {
		return nil
	}
	flush()
	return words
}

// --- pipeline -------------------------------------------------------------

// Result is what a paste turns into.
type Result struct {
	// Text is inserted at the cursor.
	Text string
	// Note, when non-empty, is surfaced to the user as a status line.
	Note string
}

// Context is what transformers are allowed to reach.
type Context struct {
	Store Store
	Cfg   Config
}

// Transformer converts a classified paste into insertable text. It returns
// false to decline, in which case the next transformer is tried.
type Transformer interface {
	Name() string
	Transform(Paste, Context) (Result, bool)
}

// Pipeline runs transformers in order; the first to accept wins. The last
// entry always accepts, so Run always produces something.
type Pipeline struct {
	classifier   Classifier
	transformers []Transformer
}

func NewPipeline() Pipeline {
	return Pipeline{
		classifier: NewClassifier(),
		transformers: []Transformer{
			fileDropTransformer{},
			urlTransformer{},
			// Additional transformers go here, ahead of plainTransformer. A
			// transformer that fetches a page's <title> for the link text —
			// or prompts for label text — only has to be appended to this
			// slice. Nothing else changes.
			plainTransformer{},
		},
	}
}

func (p Pipeline) Run(raw string, ctx Context) Result {
	paste := p.classifier.Classify(raw)
	for _, t := range p.transformers {
		if res, ok := t.Transform(paste, ctx); ok {
			return res
		}
	}
	return Result{Text: paste.Raw}
}

// fileDropTransformer handles drag-and-drop. It hands each path to the Store
// and inserts the resulting Markdown. With the default disabledStore every
// Put fails, so it inserts the paths verbatim and explains why — which is
// strictly better than today's silent raw insert, and becomes a working
// attachment feature the moment a real Store is supplied.
type fileDropTransformer struct{}

func (fileDropTransformer) Name() string { return "file-drop" }

func (fileDropTransformer) Transform(p Paste, ctx Context) (Result, bool) {
	if p.Kind != PasteFilePaths {
		return Result{}, false
	}

	store := ctx.Store
	if store == nil {
		store = disabledStore{}
	}

	var (
		lines  []string
		failed int
		reason error
	)
	for _, path := range p.Paths {
		att, err := store.Put(path)
		if err != nil {
			failed++
			reason = err
			lines = append(lines, path)
			continue
		}
		// The Store decides where a file lands; presentation is decided here,
		// from the extension, so a Store never has to think about Markdown.
		att.Kind = kindForPath(path)
		lines = append(lines, att.Markdown(filepath.Base(path)))
	}

	res := Result{Text: strings.Join(lines, "\n")}
	switch {
	case failed == len(p.Paths):
		res.Note = fmt.Sprintf("dropped %s — %v; inserted %s instead",
			plural(len(p.Paths), "file"), reason, plural(len(p.Paths), "path"))
	case failed > 0:
		res.Note = fmt.Sprintf("attached %d of %s (%v)",
			len(p.Paths)-failed, plural(len(p.Paths), "file"), reason)
	default:
		res.Note = "attached " + plural(len(p.Paths), "file")
	}
	return res, true
}

// urlTransformer handles a paste of one or more bare URLs. Left alone they'd
// fall through to plainTransformer and render as inert text in Markdown; this
// wraps each one as a link, or — when the URL's extension says it's an image
// — an embed, so a pasted image URL shows up in the preview the same way a
// dropped image file does.
type urlTransformer struct{}

func (urlTransformer) Name() string { return "url" }

func (urlTransformer) Transform(p Paste, _ Context) (Result, bool) {
	if p.Kind != PasteURLs {
		return Result{}, false
	}
	lines := make([]string, len(p.URLs))
	for i, u := range p.URLs {
		lines[i] = Attachment{Ref: u, Kind: kindForURL(u)}.Markdown("")
	}
	return Result{Text: strings.Join(lines, "\n")}, true
}

// kindForURL reads the URL's path the same way kindForPath reads a filename,
// ignoring any query string or fragment.
func kindForURL(rawURL string) AttachmentKind {
	u, err := url.Parse(rawURL)
	if err != nil {
		return AttachmentFile
	}
	return kindForPath(u.Path)
}

// plainTransformer is the terminal case: insert the text as pasted. It always
// accepts.
type plainTransformer struct{}

func (plainTransformer) Name() string { return "plain" }

func (plainTransformer) Transform(p Paste, _ Context) (Result, bool) {
	return Result{Text: p.Raw}, true
}

func kindForPath(path string) AttachmentKind {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".avif", ".svg":
		return AttachmentImage
	case ".mp3", ".wav", ".ogg", ".m4a", ".flac":
		return AttachmentAudio
	default:
		return AttachmentFile
	}
}

func plural(n int, word string) string {
	if n == 1 {
		return fmt.Sprintf("1 %s", word)
	}
	return fmt.Sprintf("%d %ss", n, word)
}

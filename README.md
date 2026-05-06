# Portfolio Website

*[Generative AI Disclosure](https://robog.net/docs/generative-ai-usage): The initial prototype for this website was created with the assistance of Claude Code, using a Large Language Model. The original design, written content, and build system were created originally by me. Commits authored by Claude Code are marked as such in the commit description. No images or art on this website were created by a generative AI model. If you suspect a portion of code may be written by AI and undisclosed, please contact me personally at robog.net/about*

My portfolio website built with Deno, Sass, and 11ty and hosted by Cloudflare at [robog.net](https://robog.net/).

## Project Structure

```
src/
├── _includes/                 # Page templates and layouts
├── _scripts/                  # Client-side JavaScript
├── _styles/                   # Sass/SCSS stylesheets
├── gallery/                   # Static assets for art pieces, some interactive some not
├── s/                         # Snippet posts (collected in live blogs and in the main blog feed)
├── blog/                      # Blog posts and assets
│   ├── attachments/           # Blog post images and media
│   └── (...).md               # Individual blog posts (Markdown)
├── index.html                 # Homepage file
├── about.md                   # About page content
├── blog.md                    # Blog listing page
└── gallery.md                 # Art portfolio listing page
```


```bash
deno task dev
```

### Build static files

```bash
deno task build
```

### License details
CC0 License in LICENSE.txt

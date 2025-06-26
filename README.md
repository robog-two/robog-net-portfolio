# Portfolio Website

A modern, responsive portfolio website built with Eleventy (11ty) and Sass,
hosted on the edge by Cloudflare at [robog.net](https://robog.net/).

## Tech Stack

### Core Technologies

- **[Eleventy (11ty)](https://www.11ty.dev/)** - Static site generator for
  building fast, modern websites
- **[Sass (SCSS)](https://sass-lang.com/)** - CSS preprocessor for maintainable
  stylesheets with variables, mixins, and modular architecture
- **[Deno](https://deno.land/)** - JavaScript/TypeScript runtime used by
  Eleventy for build processes
- **[Nunjucks](https://mozilla.github.io/nunjucks/)** - Templating engine for
  dynamic HTML generation

## Project Structure

```
src/
├── _fonts/                    # Custom web fonts (TeX Gyre Heros family)
│
├── _includes/                 # Nunjucks templates and layouts
│   ├── base.njk              # Base HTML template with head/body structure
│   ├── home.njk              # Homepage layout with special styling
│   ├── page.njk              # Standard page layout with header/nav
│   └── blog.njk              # Blog post layout with image zoom functionality
│
├── _scripts/                  # Client-side JavaScript
│   ├── header-shrink.js       # Sticky header animation on scroll
│   ├── home-blur.js           # Homepage navigation hover effects
│   └── image-zoom.js          # Modal image zoom for blog posts
│
├── _styles/                   # Sass/SCSS stylesheets
│   ├── _globals.scss          # Design system tokens and variables
│   ├── _fonts.scss            # Font-face declarations
│   ├── _navigation.scss       # Header and navigation components
│   ├── _shared-components.scss # Reusable mixins and placeholders
│   ├── base.scss             # CSS reset and foundation styles
│   ├── home.scss             # Homepage-specific styles (corner navigation)
│   ├── about.scss            # About page styles (pink theme)
│   ├── apps.scss             # Apps portfolio page (green theme)
│   ├── blog.scss             # Blog listing and post styles (yellow theme)
│   └── gallery.scss          # Gallery page styles (cyan theme)
│
├── blog/                      # Blog posts and assets
│   ├── attachments/           # Blog post images and media
│   └── (...).md             # Individual blog posts (Markdown)
│
├── index.html                 # Homepage with corner navigation
├── about.md                   # About page content
├── apps.md                    # Apps portfolio content
├── blog.md                    # Blog listing page
└── gallery.md                 # Gallery showcase content
```

## Design System

### Color Palette & Theming

The site uses a consistent 4-color brand palette with page-specific themes:

- **Pink** (`#ff6fab`) - About & CV page theme
- **Yellow** (`#ffc533`) - Blog page theme
- **Green** (`#aaf751`) - Apps portfolio theme
- **Cyan** (`#85f7ff`) - Gallery page theme
- **Primary Text** (`#0d0748`) - Main text color

### Typography

- **Headers**: TeX Gyre Heros (custom web font)
- **Body Text**: Lato (Google Fonts fallback)
- **Responsive Scale**: Fluid typography using rem units

### Spacing System

Consistent spacing scale using design tokens:

```scss
$spacing-xs: 0.25rem; // 4px
$spacing-sm: 0.5rem; // 8px
$spacing-md: 1rem; // 16px
$spacing-lg: 1.5rem; // 24px
$spacing-xl: 2rem; // 32px
```

## Architecture Patterns

### Sass Organization

- **Modular Architecture**: Separate files for different concerns
- **Design Tokens**: Centralized variables in `_globals.scss`
- **Mixins & Placeholders**: Reusable patterns in `_shared-components.scss`
- **Component-Based**: Each page has its own stylesheet
- **BEM-like Naming**: Consistent CSS class naming conventions

### Template Hierarchy

```
base.njk (HTML foundation)
├── home.njk (homepage layout)
├── page.njk (standard pages)
└── blog.njk (blog posts)
```

### Content Management

- **Markdown**: All content written in Markdown for easy editing
- **Front Matter**: YAML metadata for titles, layouts, and styling
- **Collections**: Eleventy automatically generates blog post collections
- **Asset Pipeline**: Automatic processing of styles, scripts, and fonts

## Key Features

### Interactive Homepage

- **Corner Navigation**: Triangular navigation elements in each corner
- **Hover Effects**: Animated triangles that expand on hover
- **Blur Animation**: Background content blurs when navigation is hovered
- **Business Card**: Centered layout with color blocks and typography

### Responsive Grid System

- **CSS Grid**: Flexible, responsive layouts for portfolio items
- **Auto-fit**: Automatically adjusts columns based on screen size
- **Consistent Cards**: Reusable card components across pages

### Blog System

- **Image Zoom**: Click images to view in full-screen modal
- **Responsive Images**: Breakout images on larger screens
- **Automatic Listing**: Generated blog index with excerpts and dates
- **Markdown Processing**: Full Markdown support with syntax highlighting

### Performance Optimizations

- **Font Display Swap**: Optimized web font loading
- **CSS Optimization**: Modular Sass compilation
- **Minimal JavaScript**: Lightweight interactions with vanilla JS
- **Static Generation**: Pre-built HTML for fast loading

## Development

### Local Development

```bash
deno task dev
```

### Build for Production

```bash
deno task build
```

### File Watching

Eleventy automatically watches for changes to:

- Markdown content files
- Nunjucks templates
- Sass stylesheets
- JavaScript files

## Browser Support

- Modern browsers with CSS Grid support
- Progressive enhancement for older browsers (plain old HTML)
- Responsive design for mobile, tablet, and desktop
- Accessible markup with proper ARIA labels & semantic HTML

## License

MIT

# Accessibility Audit: robog.net

**Date:** 2026-11-07
**Standard:** WCAG 2.2, Level AA
**Scope:** Every page Eleventy builds from `src/` (37 HTML pages), including blog posts, snippets, docs, the gallery, and the interactive pieces under `gallery/piece/` and `spheres/`.

## How this was tested

1. Built the site locally with `deno task build` and served `public/`.
2. Ran **axe-core 4.x** on every page in headless Chromium with the `wcag2a`, `wcag2aa`, `wcag21a`, `wcag21aa`, `wcag22aa` and `best-practice` rule sets.
3. Tested by hand with Playwright:
   - Tab order and focus indicators.
   - The image-zoom modal, using the keyboard.
   - Reflow at 320px wide (WCAG 1.4.10).
   - Text at 200% size (WCAG 1.4.4).
   - Whether a sticky header hides focused items (WCAG 2.4.11).
   - `prefers-reduced-motion`.
   - Link styling, heading outlines and alt text across all pages.
4. Read the templates, SCSS and scripts to find where each issue comes from.
5. Checked the résumé and recipe PDFs for tags.

Automated tools find roughly a third of accessibility problems. A short test with a real screen reader (VoiceOver on macOS/iOS, Orca on Linux) is still worth doing after the fixes below.

---

## Summary

The design has a good base: the body text has strong contrast (`#0d0748` on `#fffdf7`), there's a global `:focus-visible` outline, links in body text are bold and underlined, most images have descriptive alt text, the homepage nav has `aria-label`s, and the PDFs are tagged. Most of the problems come from a few shared places, so fixing about six files would resolve most of this report.

| Severity | Count | Examples |
|---|---|---|
| 🔴 Critical / serious | 6 | Pages lose `<!DOCTYPE>` and `lang`; nested links on the blog index; zoom modal can't be used with a keyboard; animation that never stops |
| 🟠 Moderate | 9 | Low-contrast text, filename alt text, empty "Apps" link, multiple `<h1>`s, untitled iframe, horizontal scrolling at 320px |
| 🟡 Minor / best practice | 8 | Interactive art pieces with no title or text alternative, pinch-zoom disabled, `<a>` with no `href` |

## Fix status (updated 2026-11-07)

All 23 items below have been fixed, except where noted. After the fixes:

- **axe-core reports 0 violations on all 37 pages** (WCAG 2.0/2.1/2.2 A and AA plus best-practice rules). It found 10 rule failures before.
- **The pages look the same.** Every page was screenshotted at 1280px and 390px before and after, then compared pixel by pixel. Every content heading keeps its exact font, size, margins and position. The only intended visual changes are the darker text from #7 and the removal of sideways scrolling from #12.

How the headings were fixed without touching the markdown: the build step (`demoteContentHeadings` in `.eleventy.deno.ts`) shifts every heading inside `<main>` down one level and records the original level in `data-md-level`. The Sass helper `h()` in `_styles/_globals.scss` then matches headings by that original level, so a markdown `#` still looks like it always did, but it's now an `<h2>` under the page's single `<h1>`.

Known differences and limits:
- Pages now render in standards mode instead of quirks mode. CSS adjustments reproduce the old spacing everywhere except two spots:
  - On the AI-usage doc, at phone width, a footnote marker wraps onto its own line, and that line is about 10px taller.
  - On the slm docs page, one line holding only inline code sits 1px lower.
- The decorative green "Apps" corner no longer shows a pointer cursor on hover, because it isn't a link any more.
- #20 (a typed note option for Passing Notes) was not added, because it would change the page's look. The email link on the About page remains the accessible way to get in touch.
- New behaviors that only show up for the people who need them:
  - A "Skip to content" link and a homepage "Pause animations" button, both visible only on keyboard focus.
  - The homepage corner triangles now also grow when focused with the keyboard.
  - Press M to mute the art pieces.
  - The homepage cubes and sky start paused for anyone with "reduce motion" turned on.
- The build also had a race: file-size labels such as "103kb" sometimes disappeared or showed "0b", because sizes were read from `public/` before the files were copied there. Sizes are now read from the source files, so the labels always appear.

---

## 🔴 Critical / serious

### 1. The build strips `<!DOCTYPE>`, `<html lang="en">`, `<head>` and `<body>` from every page
**WCAG 3.1.1 Language of Page (A) · affects all 28 Eleventy-rendered pages**

The `pdf-filesize` transform in `.eleventy.deno.ts` wraps the whole rendered page in `<!doctype html><html><body>…`, parses it, and returns `dom.body.innerHTML`. The page already contains its own `<html>`/`<head>`/`<body>`, so the parser merges them and only the inner body markup comes back. The result:

- No `lang` attribute. Screen readers fall back to the user's default voice, so a French or Spanish VoiceOver user will hear the site read with the wrong pronunciation.
- No doctype, so every page renders in **quirks mode** (`document.compatMode === "BackCompat"`, confirmed in testing). This changes box sizing and line height in subtle ways.

**Fix:** return the whole serialized document, not just the body, and don't wrap a full page in a second document:

```ts
const dom = new DOMParser().parseFromString(content, "text/html");
// …append filesize spans…
return "<!DOCTYPE html>\n" + dom.documentElement.outerHTML;
```

(There's also a stray ` }}` in `base.liquid` line 2: `<html lang="en" {% … %} }}>`. Remove it.)

### 2. Nested links on the blog listing
**WCAG 4.1.2 Name, Role, Value (A), 2.4.4 Link Purpose · `src/blog.liquid`**

Each `.blog-post-card` wraps the header, excerpt **and** a "Read More" `<a>` inside another `<a>`. HTML doesn't allow links inside links, so the parser closes the outer link early and leaves an orphan `</a>`. Screen readers then hear one very long link (title + date + the whole excerpt), followed by a second "Read More" link to the same page. That's two tab stops per post, and "Read More" doesn't say which post it opens. The outer link also has `title=""` because `post.title` is undefined; it should be `post.data.title`.

**Fix:** make the `<h2>` the only link, and if the whole card should be clickable, stretch that link with a pseudo-element:

```liquid
<article class="blog-post-card">
  <header class="post-header">
    <h2><a href="{{ post.url }}" class="card-link">{{ post.data.title }}</a></h2>
    …
  </header>
  <p class="post-excerpt">…</p>
</article>
```
```scss
.blog-post-card { position: relative; }
.card-link::after { content: ""; position: absolute; inset: 0; }
```

### 3. The image-zoom modal can't be used with a keyboard
**WCAG 2.1.1 Keyboard (A), 2.4.3 Focus Order (A), 4.1.2 (A) · `src/_scripts/image-zoom.js`**

- Images only respond to `click`. They have no `tabindex`, role or key handler, so keyboard users can't open the zoomed view.
- When the modal opens, focus stays on `<body>`. It doesn't move to the close button, isn't kept inside the modal, and doesn't return to the image afterwards.
- The modal has no `role="dialog"`, `aria-modal` or label.
- A new `keydown` listener is added each time the modal opens and is never removed. Pressing Escape a second time throws `Failed to execute 'removeChild'` (reproduced).
- The page sets `cursor: pointer` and a hover scale on every post image, but nothing tells screen-reader users that an image can be zoomed.

**Fix:** use a native `<dialog>`, which handles focus trapping, Escape and focus return for you:

```js
blogImages.forEach((img) => {
  const btn = document.createElement("button");
  btn.className = "zoom-trigger";
  btn.setAttribute("aria-label", `Enlarge image: ${img.alt}`);
  img.replaceWith(btn); btn.append(img);
  btn.addEventListener("click", () => open(img, btn));
});
function open(img, opener) {
  const d = document.createElement("dialog");
  d.className = "image-modal";
  d.setAttribute("aria-label", img.alt || "Enlarged image");
  d.innerHTML = `<button class="modal-close" aria-label="Close">×</button><img class="modal-image">`;
  d.querySelector("img").src = img.src; d.querySelector("img").alt = img.alt;
  d.querySelector(".modal-close").onclick = () => d.close();
  d.addEventListener("close", () => { d.remove(); opener.focus(); });
  document.body.append(d); d.showModal();
}
```

### 4. Animation that runs forever and ignores `prefers-reduced-motion`
**WCAG 2.2.2 Pause, Stop, Hide (A), 2.3.3 Animation from Interactions (AAA, best practice)**

- Homepage: from 15 seconds after load, the three p5 cubes flip and spin every 5 seconds, with no way to stop them. Anything that moves automatically for more than 5 seconds alongside other content needs a pause control or must stop on its own.
- The homepage sky background scrolls endlessly (`animation: scroll-left 480s infinite`).
- The site has no `prefers-reduced-motion` rule anywhere (`grep` finds nothing).
- Smaller motion: card hover lifts, image hover scale, and the corner triangles growing to fill the whole screen on hover.

**Fix:** in `base.scss`:

```scss
@media (prefers-reduced-motion: reduce) {
  *, *::before, *::after {
    animation: none !important;
    transition-duration: 0.01ms !important;
  }
}
```

In `home-blocks.js`, check `matchMedia('(prefers-reduced-motion: reduce)').matches` and skip the random animations when it's true. Either stop the cubes after a few cycles or let a click/Enter on them toggle the animation.

### 5. Homepage link names don't match the visible text, and one link is empty
**WCAG 2.5.3 Label in Name (A), 2.4.4 (A) · `src/index.html`**

- The visible text "Thoughts" has `aria-label="Blog page"`, so a speech-control user who says "click Thoughts" gets nothing (axe: `label-content-name-mismatch`). "About" → "About and CV page" and "Gallery" → "Gallery page" technically pass, but they're wordier than needed.
- The green bottom-right corner is a real, focusable link with `href="#"`, no visible text, and `aria-label="Apps page"`. It goes nowhere: screen-reader users hear "Apps page, link" and keyboard users land on a dead end.

**Fix:** remove the `aria-label`s and let the visible text be the name. Turn the Apps corner into a non-interactive `<span aria-hidden="true">` until the page exists.

### 6. Back-link names include "Back Arrow"
**WCAG 2.4.4 / 1.1.1 · `src/_includes/page.liquid`**

The arrow image has `alt="Back Arrow"`, so the link is read as "Back Arrow Home" / "Back Arrow All Posts" on all 22 pages. The arrow is decorative. Also, the `href` has about 24 leading spaces (`href="                        /"`). Browsers trim them, but they're messy.

**Fix:** `<img src="/_favicon/arrow.svg" alt="">` and `href="{% if navOverrideLink %}…{% endif %}"` with no spaces. Below 768px, `.back-text` is `display:none`, so the link's only name would be empty. Use a visually-hidden class instead of `display:none` there.

---

## 🟠 Moderate

### 7. Text contrast
**WCAG 1.4.3 Contrast (Minimum) (AA)**

| Where | Colors | Ratio | Needed |
|---|---|---|---|
| Snippet dates (`.snippet-card .post-header`, blog.scss) | `#8f8f8f` on `#fffdf7`, 12.8px | **3.17:1** | 4.5:1 |
| Gallery "*Medium:*" / "*Created:*" labels (`section.gallery-item em`) | `#00c2ea` on white, 17.6px bold | **2.12:1** | 4.5:1 |
| Blog card dates (`rgba(13,7,72,.7)` on white) | ≈ 7:1 | ✅ passes | |
| Link hover states that use the theme color as text on `#0d0748` | yellow/pink/green/cyan on navy | ✅ pass (7–12:1) | |

**Fix:** use `#6b6b6b` (5.2:1) or darker for snippet dates. For gallery labels, keep the navy text and show the cyan as an underline or left border, or use a darker cyan such as `#00728a` (5.6:1).

Also: the colored page headers (`#00c2ea`, `#a0e80a`, `#ffc533`, `#ff6fab`) all pass with navy text. Don't switch them to white text.

### 8. Filenames used as alt text
**WCAG 1.1.1 Non-text Content (A) · `src/s/2026-08-11-this-game-is-so-cute.md`**

The six Tiny Glade screenshots have alt text such as `20260806002923_1.jpg` (this is what markdown-it outputs when the alt is empty). Screen readers read out the filename. Write a short description of each scene, or use `![](…)` with an empty alt only if they are purely decorative (they aren't, since they're the subject of the post).

### 9. Slide images and charts need a fuller text alternative
**WCAG 1.1.1 (A) · `src/blog/human-app-ai-world.md`, `src/gallery/piece/drawings/index.html`**

- The keynote slides use the slide title as alt text ("Artificial Intelligence Opinion Survey", "User research word cloud", "Examples of data sources"). The prose around them covers most of it, but data-heavy slides (pie chart, survey, word cloud) should state the key numbers or words in the alt text or the next paragraph.
- On desktop, `slide-viewer.js` copies each image into `#slideViewer`, so screen readers can hear each slide twice. Add `aria-hidden="true"` to `#slideViewer`. The original images stay in the post at `height: 0`, so they're still read in order.
- The drawings gallery uses alt text such as "Folioslides 0–5" and "Homework 1". Describe what each drawing shows, as the main gallery page already does well.

### 10. Heading structure
**WCAG 1.3.1 Info and Relationships (A), 2.4.6 Headings and Labels (AA)**

The page title is the `<h1>` in the header, but markdown `#` headings in posts also become `<h1>`:

- `/about/`: 4 × h1
- `/docs/generative-ai-usage/`: 5 × h1
- `/blog/podman-…/` and `/blog/i-wrote-a-fan-controller/`: 3 × h1

Screen-reader users who move through the page by heading see a flat list with no hierarchy. Change in-post `#` to `##` (and `##` → `###`), then adjust `post-content` in `_shared-components.scss` so `h2` gets the current `h1` styling. Also, on the homepage, "Writes Joyful, Useful Software." is an `<h2>` but it's a tagline, not a section heading. Use a `<p>`.

### 11. Untitled iframe
**WCAG 4.1.2 (A) · `src/blog/birthday-cake.md`**

The embedded recipe PDF `<iframe>` has no `title`. Add `title="Baffling blueberry birthday cake recipe (PDF)"`. Embedded PDF viewers are hard to use with screen readers, so also add a plain download link next to it.

### 12. Horizontal scrolling at 320px
**WCAG 1.4.10 Reflow (AA)**

- `/gallery/` scrolls to 352px because `%grid-layout` uses `minmax(350px, 1fr)`. Use `minmax(min(350px, 100%), 1fr)`.
- `/blog/podman-systemd-rootless-containers/` scrolls to 503px because of `<pre><code>` blocks. Code is allowed to scroll, but only inside its own box: add `pre { overflow-x: auto; }` (and `tabindex="0"` on scrollable `<pre>` so keyboard users can scroll it).
- `footer { width: 100vw; position: absolute; height: 100px; }`: `100vw` includes the scrollbar on desktop browsers and can create a small horizontal scrollbar. A fixed height can also cut off footer text when users enlarge text. Use `width: 100%` and `min-height` instead of `height`.

### 13. Homepage keyboard focus doesn't match hover
**WCAG 2.4.7 Focus Visible (AA), 1.4.13**

The focus outline appears, but as a plain square around a triangle (see `home.scss`). Keyboard users never see the triangle grow, which is the main feedback mouse users get. Add `:focus-visible` next to each `:hover` rule in `home.scss` (and in `triangle-gradient`'s `@at-root`). `home-blur.js` also only listens for `mouseenter`/`mouseleave`, so add `focus`/`blur` too. (It also leaves a `console.log("yep")` in place.)

### 14. No skip link
**WCAG 2.4.1 Bypass Blocks (A)**

Pages have only one nav link before `<main>`, so this is low impact, but a skip link is cheap. Add `<a class="skip-link" href="#main">Skip to content</a>` in `page.liquid`, plus `id="main"` on `<main>`.

### 15. Links that look the same as body text
**WCAG 1.4.1 Use of Color (A)**

Most links are bold and underlined, which works. But file links (`.pdf`, `.mp3`, `.txt`) and snippet links have `text-decoration: none` and rely on a border or icon. That's fine as long as the border stays. The blog card link (`style="text-decoration: none"`) relies on its heading style alone. After fix #2, add an underline on hover/focus at least.

---

## 🟡 Minor / best practice

16. **Interactive pieces have no title** (`concerto-in-petri-dish`, `herb-journey`, `our-world-planes-overhead`, `maze`, `ozymandias`, `stack-game`, `bulletin-board-letters`): add a `<title>` to each. Seven pages show only the URL in tabs, history and screen readers (axe: `document-title`, WCAG 2.4.2).
17. **Canvas-only art has no text alternative** (WCAG 1.1.1): give each p5 piece a visible or visually-hidden `<h1>` and a short description, e.g. reuse the descriptions already on `/gallery/`. p5's `describe()` adds a screen-reader description to the canvas in one line.
18. **No way to stop audio** in Herb Journey, Our World Planes Overhead, Concerto and Stack Game (WCAG 1.4.2). Playback starts on click, so this is technically user-started, but sound that loops forever should have a mute/stop key or button.
19. **Pinch-zoom disabled** on `/spheres/passing-notes/` (`user-scalable=no, maximum-scale=1`, WCAG 1.4.4). Remove those parts, or call `preventDefault` only on the canvas's touch events.
20. **Passing Notes is drawing-only.** It's offered as a contact method, but there's no keyboard or screen-reader way to write a note. The email link next to it covers this; consider adding a text field to the note page as well.
21. **`<a onclick="redrawMaze()">` with no `href`** in `/spheres/maze/` can't be focused. Use a `<button>`.
22. **`aria-label` on a plain `<div>`** in `bulletin-board-letters` (axe: `aria-prohibited-attr`). Add `role="img"`.
23. **Misc. markup:** `<figcaption>` used outside a `<figure>` on `/about/` (wrap the image and caption in `<figure>`); `<br/><br/><br/>` used for spacing; the Mastodon `rel="me"` `<a>` sits inside `<head>`, which is invalid. Use `<link rel="me" href="…">` instead. The snippet `<time>` element includes the text "Snippet Post •", so move that label outside `<time>`.

---

## What's already working well ✅

- Body text contrast is very high (18:1).
- `:focus-visible` outline is defined globally, and every link tested showed a visible focus ring.
- Tab order follows visual order on every page tested.
- Most images on the About page and blog posts have good, descriptive alt text (the About photo is excellent).
- `font-display: swap` and system fallback fonts mean text never stays invisible while fonts load.
- Body text is 1.1rem with line-height 1.6–1.7, a comfortable size.
- The résumé and recipe PDFs are tagged (`StructTreeRoot`) and the résumé declares `/Lang (en)`.
- The homepage nav is wrapped in `<nav aria-label="Main navigation">`, and decorative triangles are `aria-hidden`.
- The print stylesheet is thoughtful.
- The sticky header didn't hide focused elements at 200% text size.

## Suggested order of work

1. Fix the build transform (#1). One file, fixes 28 pages.
2. Blog listing links (#2) and image-zoom modal (#3).
3. Reduced-motion block plus homepage animation control (#4).
4. Homepage link names and dead Apps link (#5), back-arrow alt (#6).
5. Contrast tweaks (#7), alt text (#8, #9), heading levels (#10).
6. Everything else as time allows.

Re-run the check after changes:

```bash
deno task build
```
```bash
npx @axe-core/cli http://localhost:8080 --tags wcag2a,wcag2aa,wcag21aa,wcag22aa
```

## References

- W3C, *Web Content Accessibility Guidelines (WCAG) 2.2*: https://www.w3.org/TR/WCAG22/
- W3C WAI, *Understanding SC 2.5.3 Label in Name*: https://www.w3.org/WAI/WCAG22/Understanding/label-in-name
- W3C WAI ARIA Authoring Practices, *Dialog (Modal) Pattern*: https://www.w3.org/WAI/ARIA/apg/patterns/dialog-modal/
- MDN, *`<dialog>` element*: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/dialog
- MDN, *`prefers-reduced-motion`*: https://developer.mozilla.org/en-US/docs/Web/CSS/@media/prefers-reduced-motion
- Deque, *axe-core rule descriptions*: https://dequeuniversity.com/rules/axe/
- p5.js reference, *`describe()`*: https://p5js.org/reference/p5/describe/

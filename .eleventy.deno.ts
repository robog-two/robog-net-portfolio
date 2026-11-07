import eleventySass from "@11tyrocks/eleventy-plugin-sass-lightningcss";
import { eleventyImageTransformPlugin } from "@11ty/eleventy-img";
import markdownItCheckbox from "markdown-it-task-checkbox";
import markdownItFootnote from "markdown-it-footnote";
import syntaxHighlight from "@11ty/eleventy-plugin-syntaxhighlight";
import markdownIt from "markdown-it";
import { DOMParser } from "@b-fuze/deno-dom";
import type { EleventyConfig, EleventyPage } from "@11ty/eleventy";

// Format file size in human-readable format
function formatFileSize(bytes: number): string {
  const units = ["b", "kb", "mb", "gb"];
  let size = bytes;
  let unitIndex = 0;

  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024;
    unitIndex++;
  }

  return `${Math.ceil(parseFloat(size.toFixed(size < 10 ? 2 : 1)))}${
    units[unitIndex]
  }`;
}

export default (eleventyConfig: EleventyConfig) => {
  eleventyConfig.addPlugin(eleventySass);
  eleventyConfig.addPlugin(syntaxHighlight);

  // Add image optimization plugin
  eleventyConfig.addPlugin(eleventyImageTransformPlugin, {
    formats: ["webp", "svg"],
    widths: ["1280", null],
    defaultAttributes: {
      loading: "lazy",
      decoding: "async",
    },
    svgShortCircuit: true,
  });

  // Filters which are used like {{ data.thing | filter }}
  eleventyConfig.addFilter("md", function (content = "") {
    return markdownIt({ html: true }).render(content);
  });
  eleventyConfig.addFilter("excerpt", function (content = "") {
    return htmlToText(markdownIt({ html: true }).render(content))
      .replace(/\s+/g, " ") /* bug fix: convert line breaks to preserve spacing */
      .trim().slice(0, 200)
      .replace(/\s+\w+$/, "") /* remove last cut-off word chunk */
      .replace(/[\.\,\-\—\;\:]+\s*$/, "") /* remove last punctuation/spacing after cutoff */
      + "…";
  });
  eleventyConfig.addFilter("keys", function (content = {}) {
    return JSON.stringify(Object.keys(content));
  });

  eleventyConfig.amendLibrary("md", (mdLib: any) => {
    mdLib
      .use(markdownItCheckbox)
      .use(markdownItFootnote);
  });

  eleventyConfig.addPassthroughCopy("src/cv");

  eleventyConfig.addPassthroughCopy({ "src/_favicon": "/" });

  // Gallery pieces are completely static and should not be processed
  eleventyConfig.addPassthroughCopy("src/gallery/piece");
  eleventyConfig.addPassthroughCopy("src/spheres");

  // Copy font files to output
  eleventyConfig.addPassthroughCopy("src/_fonts");

  // Copy script files to output
  eleventyConfig.addPassthroughCopy("src/_scripts");

  // Copy blog attachments to output (still needed for source images)
  eleventyConfig.addPassthroughCopy("src/blog/media");

  // Copy license file
  eleventyConfig.addPassthroughCopy("LICENSE.txt");

  // SVG Badges
  eleventyConfig.addPassthroughCopy("src/badges");

  // Post-process every HTML page: add file sizes to download links,
  // fix heading levels inside the page content, and make code blocks
  // reachable by keyboard.
  eleventyConfig.addTransform("html-postprocess", async function (
      this: { page: EleventyPage },
      content: string
  ): Promise<string> {
    if (!this.page.outputPath || !this.page.outputPath.endsWith(".html")) {
      return content;
    }

    // Full documents are parsed as-is so the doctype, <html lang>, <head>
    // and <body> survive. Fragments (no <html>) are wrapped and unwrapped.
    const isFullDocument = /<html[\s>]/i.test(content);

    try {
      const dom = new DOMParser().parseFromString(
        isFullDocument
          ? content
          : "<!doctype html><html><body>" + content + "</body></html>",
        "text/html",
      );

      await addFileSizes(dom, this.page);
      demoteContentHeadings(dom);

      // Scrollable code blocks need to be focusable so keyboard users can
      // scroll them (WCAG 2.1.1).
      for (const pre of dom.querySelectorAll("main pre")) {
        if (!pre.hasAttribute("tabindex")) pre.setAttribute("tabindex", "0");
      }

      if (!isFullDocument) return dom.body.innerHTML;
      const doctype = /^\s*<!doctype[^>]*>/i.exec(content)?.[0] ?? "<!DOCTYPE html>";
      return doctype + "\n" + dom.documentElement!.outerHTML;
    } catch {
      // If DOM parsing fails, return content as-is
      return content;
    }
  });

  return {
    dir: {
      input: "src",
      output: "public",
    },
  };
};

// Append a file size label to links that point at downloads.
async function addFileSizes(dom: any, page: EleventyPage) {
  const links = [
    ...dom.querySelectorAll("a[href$='.pdf']"),
    ...dom.querySelectorAll("a[href$='.txt']"),
    ...dom.querySelectorAll("a[href$='.mp3']"),
  ];

  for (const link of links) {
    const href = link.getAttribute("href");
    if (!href) continue;

    try {
      // Read the size from the source file. The copy in public/ may not
      // exist (or may be half-written) when this transform runs.
      const candidates = href.startsWith("/")
        // Under src/, or at the project root (e.g. /LICENSE.txt)
        ? ["src" + href, "." + href]
        // Relative path - resolve from page directory
        : [page.inputPath.replace(/[^/]*$/, "") + href];

      let stats: Deno.FileInfo | undefined;
      for (const candidate of candidates) {
        try {
          stats = await Deno.stat(decodeURIComponent(candidate));
          break;
        } catch {
          // try the next location
        }
      }
      if (!stats) continue;
      const sizeStr = formatFileSize(stats.size);

      // Create and append span
      const span = dom.createElement("span");
      span.className = "filesize";
      span.textContent = sizeStr;
      link.appendChild(span);
    } catch {
      // File not found or error reading file, skip this link
      continue;
    }
  }
}

// The page title in the header is the page's only <h1>. Markdown "#"
// headings inside <main> would otherwise be extra <h1>s, so shift every
// heading in <main> down one level (h1 -> h2, h2 -> h3, ...).
// `data-md-level` keeps the original level so the styles stay the same;
// see the `h()` helper in _styles/_globals.scss.
function demoteContentHeadings(dom: any) {
  const main = dom.querySelector("main");
  if (!main || !dom.querySelector(".page-header h1") || !main.querySelector("h1")) {
    return;
  }

  for (const heading of main.querySelectorAll("h1, h2, h3, h4, h5, h6")) {
    const level = Number(heading.tagName[1]);
    const newLevel = Math.min(level + 1, 6);
    const replacement = dom.createElement("h" + newLevel);
    for (const name of heading.getAttributeNames()) {
      replacement.setAttribute(name, heading.getAttribute(name));
    }
    replacement.setAttribute("data-md-level", String(level));
    if (level + 1 > 6) replacement.setAttribute("aria-level", String(level + 1));
    while (heading.firstChild) replacement.appendChild(heading.firstChild);
    heading.parentNode.replaceChild(replacement, heading);
  }
}

// Source - https://stackoverflow.com/a/50822488
// Posted by chrmcpn, modified by community. See post 'Timeline' for change history
// Retrieved 2026-05-01, License - CC BY-SA 4.0

function htmlToText(html: string) {
  //remove code brakes and tabs
  // soft line breaks separate words, so they become spaces
  html = html.replace(/\n/g, " ");
  html = html.replace(/\t/g, " ");

  //keep html brakes and tabs
  html = html.replace(/<\/td>/g, " ");
  html = html.replace(/<\/table>/g, " ");
  html = html.replace(/<\/tr>/g, " ");
  html = html.replace(/<\/p>/g, " ");
  html = html.replace(/<\/div>/g, " ");
  html = html.replace(/<\/h>/g, " ");
  html = html.replace(/<br>/g, " ");
  html = html.replace(/<br( )*\/>/g, " ");
  html = html.replace(/>/g, "> ");

  //parse html into text
  const dom = (new DOMParser()).parseFromString(
    "<!doctype html><html><body>" + html + "</body></html>",
    "text/html",
  );
  return dom.body.textContent;
}

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
    return htmlToText(markdownIt({ html: true }).render(content)).slice(0, 200)
      .replace(/\.+$/, "").replace(/\s\w+$/, "") + "…";
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

  // Add file sizes to PDF links
  eleventyConfig.addTransform("pdf-filesize", async function (
      this: { page: EleventyPage },
      content: string
  ): Promise<string> {
    if (!this.page.outputPath.endsWith(".html")) {
      return content;
    }

    try {
      const dom = new DOMParser().parseFromString(
        "<!doctype html><html><body>" + content + "</body></html>",
        "text/html",
      );

      const links = [
        ...dom.querySelectorAll("a[href$='.pdf']"),
        ...dom.querySelectorAll("a[href$='.txt']"),
        ...dom.querySelectorAll("a[href$='.mp3']"),
      ];

      for (const link of links) {
        const href = link.getAttribute("href");
        if (!href) continue;

        try {
          // Resolve file path relative to project root
          let filePath = href;
          if (href.startsWith("/")) {
            filePath = "public" + href;
          } else {
            // Relative path - resolve from page directory
            const pageDir = this.page.inputPath.replace(/[^/]*$/, "");
            filePath = pageDir.replace("src/", "") + href;
          }

          const stats = await Deno.stat(filePath);
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

      return dom.body.innerHTML;
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

// Source - https://stackoverflow.com/a/50822488
// Posted by chrmcpn, modified by community. See post 'Timeline' for change history
// Retrieved 2026-05-01, License - CC BY-SA 4.0

function htmlToText(html: string) {
  //remove code brakes and tabs
  html = html.replace(/\n/g, "");
  html = html.replace(/\t/g, "");

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

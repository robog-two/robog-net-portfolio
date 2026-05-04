import eleventySass from "@11tyrocks/eleventy-plugin-sass-lightningcss";
import { eleventyImageTransformPlugin } from "@11ty/eleventy-img";
import markdownItCheckbox from "markdown-it-task-checkbox";
import markdownItFootnote from "markdown-it-footnote";
import syntaxHighlight from "@11ty/eleventy-plugin-syntaxhighlight";
import markdownIt from "markdown-it";
import { DOMParser } from "@b-fuze/deno-dom";

export default (eleventyConfig) => {
  eleventyConfig.addPlugin(eleventySass);
  eleventyConfig.addPlugin(syntaxHighlight);

  // Add image optimization plugin
  eleventyConfig.addPlugin(eleventyImageTransformPlugin, {
    formats: ["webp", "avif", "svg"],
    widths: ["240", "480", "760", "1280", null],
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
  eleventyConfig.addCollection("threadEntries", function(collection: any) {
    // Returns all pages with a 'thread' field set (thread entry pages)
    return collection
      .getAll()
      .filter((item: any) => item.data.thread);
  });

  eleventyConfig.amendLibrary("md", (mdLib) => {
    mdLib
      .use(markdownItCheckbox)
      .use(markdownItFootnote);
  });

  eleventyConfig.addPassthroughCopy("src/cv");

  //eleventyConfig.setServerPassthroughCopyBehavior("passthrough");
  eleventyConfig.addPassthroughCopy({ "src/_favicon": "/" });

  // Gallery pieces are completely static and should not be processed
  eleventyConfig.addPassthroughCopy("src/gallery/piece");

  // Copy font files to output
  eleventyConfig.addPassthroughCopy("src/_fonts");

  // Copy script files to output
  eleventyConfig.addPassthroughCopy("src/_scripts");

  // Copy blog attachments to output (still needed for source images)
  eleventyConfig.addPassthroughCopy("src/blog/media");

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
    "<!doctype html><html><body>" + html +  "</body></html>",
    "text/html",
  );
  return dom.body.textContent;
}

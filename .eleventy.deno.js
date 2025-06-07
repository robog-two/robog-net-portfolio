import eleventySass from "npm:@11tyrocks/eleventy-plugin-sass-lightningcss";
import { eleventyImageTransformPlugin } from "npm:@11ty/eleventy-img";
import markdownItCheckbox from "npm:markdown-it-task-checkbox";
import markdownItFootnote from "npm:markdown-it-footnote";

export default (eleventyConfig) => {
  eleventyConfig.addPlugin(eleventySass);

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

  eleventyConfig.amendLibrary("md", (mdLib) => {
    mdLib
      .use(markdownItCheckbox)
      .use(markdownItFootnote);
  });

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

const eleventySass = require("@11tyrocks/eleventy-plugin-sass-lightningcss");
const { eleventyImageTransformPlugin } = require("@11ty/eleventy-img");
const markdownItCheckbox = require("markdown-it-task-checkbox");
const markdownItFootnote = require("markdown-it-footnote");

module.exports = (eleventyConfig) => {
  eleventyConfig.addPlugin(eleventySass);

  // Add image optimization plugin
  eleventyConfig.addPlugin(eleventyImageTransformPlugin, {
    formats: ["webp", "avif", "jpeg", "svg"],
    widths: ["auto"],
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

  eleventyConfig.setServerPassthroughCopyBehavior("passthrough");
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

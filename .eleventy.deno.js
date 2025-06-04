import eleventySass from "npm:@11tyrocks/eleventy-plugin-sass-lightningcss";
import { eleventyImageTransformPlugin } from "npm:@11ty/eleventy-img";
import markdownItCheckbox from "npm:markdown-it-task-checkbox"

export default (eleventyConfig) => {
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

  eleventyConfig.amendLibrary("md", (mdLib) => mdLib.use(markdownItCheckbox));

  eleventyConfig.setServerPassthroughCopyBehavior("passthrough");
  eleventyConfig.addPassthroughCopy({ "src/_favicon": "/" });

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

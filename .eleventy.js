const eleventySass = require("@11tyrocks/eleventy-plugin-sass-lightningcss");
const { eleventyImageTransformPlugin } = require("@11ty/eleventy-img");

module.exports = (eleventyConfig) => {
  eleventyConfig.addPlugin(eleventySass);

  // Add image optimization plugin
  eleventyConfig.addPlugin(eleventyImageTransformPlugin, {
    formats: ["webp", "jpeg"],
    widths: ["auto"],
    defaultAttributes: {
      loading: "lazy",
      decoding: "async",
    },
  });

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

import eleventySass from "npm:@11tyrocks/eleventy-plugin-sass-lightningcss";
import { eleventyImageTransformPlugin } from "npm:@11ty/eleventy-img";

export default (eleventyConfig) => {
  eleventyConfig.addPlugin(eleventySass);
  
  // Add image optimization plugin
  eleventyConfig.addPlugin(eleventyImageTransformPlugin, {
    formats: ["webp", "jpeg"],
    widths: ["auto"],
    defaultAttributes: {
      loading: "lazy",
      decoding: "async"
    }
  });
  
  // Copy font files to output
  eleventyConfig.addPassthroughCopy("src/_fonts");
  
  // Copy script files to output
  eleventyConfig.addPassthroughCopy("src/_scripts");
  
  // Copy blog attachments to output (still needed for source images)
  eleventyConfig.addPassthroughCopy("src/blog/attachments");

  return {
    dir: {
      input: "src",
      output: "public",
    },
  };
};
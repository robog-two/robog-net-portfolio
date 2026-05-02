#!deno

console.log("Type your post (markdown + html supported) and press Ctrl+D to write to file.")

// thanks to Claude 4.7 sonnet
const post = await new Response(Deno.stdin.readable).text();

const fileName = post.replace(/[^a-zA-Z0-9]/g, "").slice(0, 10).toLowerCase() +
  ".md";

const now = new Date();
const minutes = now.getMinutes() < 10
  ? "0" + now.getMinutes()
  : now.getMinutes();
const hour = now.getHours() % 12 == 0 ? 12 : now.getHours() % 12;
const ampm = now.getHours() >= 12 ? "PM" : "AM";
const timestr = `${hour}:${minutes}${ampm}`;

const header = `
---
date: ${now.getFullYear()}-${now.getMonth()}-${now.getDay()}
time: ${timestr}
tags: post
snippet: yes
---
`;

Deno.writeFileSync(
  "./src/snippets/" + fileName,
  new TextEncoder().encode(header + post),
);

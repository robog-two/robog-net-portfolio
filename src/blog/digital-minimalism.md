---
title: Digital Minimalism Journey
layout: blog.njk
date: 2026-05-04
description: yes
live: yes
tags: post
---

I am on a continuing journey towards digital minimalism. This is a live-blog of little tips, thoughts, and updates related to reducing the prevalence and impact that digital devices have on my life.

<!-- TODO: Move this stuff into a template file separate from this -->
<div style="display: block; width: 100%; height: 2px; background-color: lightgray; margin-top: 2em; margin-bottom: 4em;"></div>


{%- for snippet in collections.post-digital-minimalism -%}

<article class="snippet-card">
<header class="post-header">
<time datetime="{{ snippet.data.date | date: '%Y-%m-%d' }}">{{ snippet.data.date |
date: '%B %d, %Y' }}  {{ snippet.data.time }}</time>
</header>
{{ snippet.data.page.rawInput | md }}
</article>

{%- endfor -%}

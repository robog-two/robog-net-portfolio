---
title: Blog
layout: page.njk
style: blog
---

Welcome to my blog where I share thoughts on technology, computer science, and my projects.

{%- for post in collections.post reversed -%}
<article class="blog-post-card">
    <header class="post-header">
        <h2><a href="{{ post.url }}">{{ post.data.title }}</a></h2>
        <time datetime="{{ post.data.date | date: '%Y-%m-%d' }}">{{ post.data.date | date: '%B %d, %Y' }}</time>
    </header>
    {%- if post.data.excerpt -%}
    <p class="post-excerpt">{{ post.data.excerpt }}</p>
    {%- endif -%}
    <a href="{{ post.url }}" class="read-more">Read More →</a>
</article>
{%- endfor -%}
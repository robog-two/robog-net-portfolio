---
title: Sam's Blog.
layout: page.njk
style: blog
---

{%- for post in collections.post reversed -%}

<article class="blog-post-card">
<a href="{{ post.url }}" style="text-decoration: none" title="{{ post.title }}">
<header class="post-header">
<h2>{{ post.data.title }}</h2>
<time datetime="{{ post.data.date | date: '%Y-%m-%d' }}">{{ post.data.date |
date: '%B %d, %Y' }}</time>
</header> {%- if post.data.excerpt -%}
<p class="post-excerpt">{{ post.data.excerpt }}</p> {%- endif -%}
<a href="{{ post.url }}" class="read-more">Read More</a>

</a>
</article>

{%- endfor -%}

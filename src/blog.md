---
title: Sam's Thoughts.
layout: page.njk
primaryColor: '#ffc533'
style: blog
---

{%- for post in collections.post reversed -%}

{%- if post.data.snippet == "yes" -%}
<article class="snippet-card">
<header class="post-header">
<span>Sam's Snippets ✂️</span>
<time datetime="{{ post.data.date | date: '%Y-%m-%d' }}">{{ post.data.date |
date: '%m/%d/%y' }}  {{ post.data.time }}</time>
</header>
{{ post.data.page.rawInput | md }}
</article>
{%- else -%}
<article class="blog-post-card">
<a href="{{ post.url }}" style="text-decoration: none" title="{{ post.title }}">
<header class="post-header">
<h2>{{ post.data.title }}</h2>
<time datetime="{{ post.data.date | date: '%Y-%m-%d' }}">{{ post.data.date |
date: '%B %d, %Y' }}</time>
</header>
<p class="post-excerpt">{{ post.data.page.rawInput | excerpt }}</p>
<a href="{{ post.url }}" class="read-more">Read More</a>

</a>
</article>
{%- endif -%}

{%- endfor -%}

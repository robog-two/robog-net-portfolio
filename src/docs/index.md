---
title: Sam's Documentation.
layout: about.liquid
primaryColor: '#a0e80a'
description: Various documentation for my public APIs, libraries, and projects.
hideNav: true
---

A collection of information and guides pertaining to my various projects.

<ul>
{%- for docroot in collections.docroot -%}

<li><a href="{{ docroot.url }}">{{ docroot.data.title }}</a></li>

{%- endfor -%}

</ul>

---
title: Bomb&nbsp;Squad&colon; How&nbsp;I&nbsp;defused Carnegie&nbsp;Mellon's Binary&nbsp;Bomb
layout: blog.njk
date: 2025-06-04
tags: post
excerpt: "... I *really* wanted to have a perfect defusal, so I hatched a plan: I would put the bomb in a sealed docker container where no radio signals could escape ..."
---

## Document To-dos

 - [ ] Find source code and upload to GitHub
 - [ ] Post Content
    - [ ] Intro, what is the Binary Bomb
    - [ ] Journey into docker
    - [ ] Journey into brute forcing with deno
    - [ ] Decoding the bomb with GDB
    - [ ] Further potential for exploration

## Level 3 Jump Table

In assembly, the jump table for level_3() looks like this:

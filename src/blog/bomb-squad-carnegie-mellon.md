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
  - [ ] Journey into podman/flawless defusal
    - [ ] Path 1: MITM
    - [ ] Path 2: Cut all contact, gdb jumping
    - [ ] Path 3: Nops with docker network
  - [ ] Journey into brute forcing with Deno
  - [ ] Decoding the bomb with GDB
  - [ ] Further potential for exploration

## What is the Binary Bomb?

Carnegie Mellon University's Randal Bryant and David O'Hallaron[^1] created a
challenge for students learning C that's similar to a CTF. "Dr. Evil" (of the
Austin Powers franchise) has planted a bomb which can only be defused through
solving a series of puzzles based on math, recursion, and C data types.

In a classroom setting, each bomb networks to a central server, where failed
solutions or "explosions" subtract points on a leaderboard and successful
solutions add points, with the date of the completed solve being used as a
tiebreaker.

For my bomb, I endeavored to have a perfect solution, solving the bomb
completely with zero explosions on the leaderboard.

## My Journey to a Flawless Defusal

### Path #1: No SSL, no problem!

When I heard that the bombs networked, my first plan was to run the bomb inside
a container, and network it to a fake server where only I would receive the
information about when it explodes and what the score is. I knew that the
central server used http (not https), so it would be relatively trivial to
intercept and redirect traffic to an artificial server as there is no
verification of authenticity or encryption present.

One small problem was finding out how exactly requests to the server were
formatted. I had a general idea of the domains and ports being used to contact
the server from where the leaderboard page was located, so I started by adding
these to the hosts file inside of my podman container and redirecting them to
localhost, where I used `npx http-server` to log any incoming requests.

Unfortunately, it seemed like

### Path #2: Jmp of Faith

### Path #3: Nop nop nop <div style="display: inline-block;transform: rotate(45deg);">◔</div> • • •

## Potential for Further Exploration

[^1]: You can find the original lab
    [here](https://csapp.cs.cmu.edu/3e/labs.html), or use an archive.
    [Internet Archive](https://web.archive.org/web/20250604211438/https://csapp.cs.cmu.edu/3e/labs.html)
    / [Archive Today](https://archive.ph/rIfZt)

[^2]: It would be interesting to revisit this project with Siemen's
    (Edgeshark)[https://github.com/siemens/edgeshark] to see if I could have
    made more headway using this method. I think this would also be a great tool
    for security research in any containerized environment, so I'll surely find
    use for it later.

---
title: Organizing&nbspMy&nbspFiles (The&nbspHard&nbspWay)
layout: blog.njk
date: 2025-06-04
tags: # post TODO: uncomment upon publishing
description: My hard drive needs a spring cleaning. Instead of actually doing the work, I decided to write a cutting-edge agentic multi-modal local-llm file organizing app in Python. (Because this is obviously easier than organizing them myself!)
---

My Documents folder is a mess. Recently, while I was organizing my photos, code
repositories, screenshots, PDF's etc. etc. etc., I realized that a lot of the
places I put my files are already named after the files that live in them. So,
in the classic developer fashion of wasting hours to save minutes, I decided to
make a large language model scan the folders, scan the files, and then put each
file in its related folder based on the name. Couldn't be _that_ complicated,
right?

## Choosing Tech Stack

I chose the tech stack carefully to meet my needs, there are definitely other
ways to write this app, but I decided that this stack would be best for my use
cases.

- **Python** — Python is a good choice here for two reasons, the most obvious
  being the number of supported ML libraries, in the event that I need to do
  additional machine learning work aside from the LLM itself. The second great
  reason is that Python is present on most Linux distros by default, which means
  that if I want to run this app on a server in the cloud or on my variety of
  IoT devices, it's fairly straightforward plug-and-play (though it might
  require hosting the LLM somewhere else)

- **Ollama** — The gold standard for running LLMs locally. I especially
  appreciate the API parity with OpenAI, which makes it easy for me to swap
  between running models locally for cost and running them in the cloud for
  accuracy and offloading the heavy resource usage of LLMs

- **Textual** - I really wanted a
  [Bubble Tea–style](https://github.com/charmbracelet/bubbletea) CLI that could
  take full advantage of modern terminal features to convey information about
  file types and display file tree graphs with less code complexity. Plus, my
  app would need to process input and display output in parallel (say, the LLM
  is organizing files, but names a folder poorly, you should be able to correct
  it with a new prompt without waiting for it to finish organizing) I chose
  Textual as it's pretty much the only Python CLI library that fits the bill.

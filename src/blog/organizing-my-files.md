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

# Getting Started <!-- TODO: Rename this heading to something more descriptive -->

## Choosing my Tech Stack

I chose the tech stack carefully to meet my needs, there are definitely other
ways to write this app, but I decided that this stack would be best for my use
cases.

- **Python** - Python is a good choice here for two reasons, the most obvious
  being the number of supported ML libraries, in the event that I need to do
  additional machine learning work aside from the LLM itself. The second great
  reason is that Python is present on most Linux distros by default, which means
  that if I want to run this app on a server in the cloud or on my variety of
  IoT devices, it's fairly straightforward plug-and-play (though it might
  require hosting the LLM somewhere else)

- **Ollama** - The gold standard for running LLMs locally. I especially
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

## Getting Started Quickly with Claude[^1]

Something new I've started trying in the age of AI is using an LLM to build out
the start of my project. It's really useful for getting started with frameworks
quickly, and since I am new to Textual, I decided to ask Claude to:

```llm prompt
Write a basic main file skeleton using the textual library for python that creates a CLI window with an input at the bottom and a space for messages to be placed starting above the input and moving upwards. Don't implement any functionality when messages are sent except for placing them into the message window. Use the internet to consult documentation before drafting your code.
```

This gave me a great starting point. There was a small issue with the styling of
the input, and the theme was really strange. I _could_ have used some more
prompting with Claude to fix this, but I find that it's better to fix Claude's
errors yourself with documentation, as they're usually areas where your
assumptions would be challenged by the actual behavior of the framework. Plus,
Claude and pretty much all LLMs start to make a mess out of code as they
continue to edit it iteratively, since they need to be prompted to clean up or
refactor code.

I fixed the input box, added the solarized-light theme, and I ended up with a
basic chat interface that looked like this:

![A visual studio code terminal displaying the following text in a beige terminal UI window: Sage: Hi, I'm Sage. I can help you sort and organize your files. Just let me know how I can assist.](/blog/media/sage-cli-1.png)

## Connecting Ollama

Ollama with Python is very straightforward, I created a new LLM.py to store my
simple wrapper, and then I only need to call the message() function when the
user submits the input box in my main Textual UI.

Viola! We have an LLM chat app.

![A visual studio code terminal displaying the following text in a beige terminal UI window: Sage: Hi, I'm Sage. I can help you sort and organize your files. Just let me know how I can assist. You: Hi. Sage: Hello! How can I assist you today? You: What is 2+2. Sage: 2 plus 2 equals 4.](/blog/media/sage-cli-2.png)

Well, it's not _that_ simple. While this is an LLM chat app, we have more
choices to make.

[^1]: [Claude](https://claude.ai/) is [Anthropic](https://www.anthropic.com)'s
    Large Language Model chatbot/research assistant/programming agent. It does
    more than just that, but for the purposes of this article, that's how I'm
    using it.

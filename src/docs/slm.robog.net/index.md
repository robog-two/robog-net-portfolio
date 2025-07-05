---
title: Small Language Model API
layout: docs.njk
tags:
  - docroot
  - docs-slm
description: Public API for using various small language models.
---

The Small Language Model API allows you to integrate Ollama-based language
models into your project for free. I created it because I wanted to be able to
write webapps that use these models without requiring users to have lots of
available processing power. Please use it kindly as I am hosting this service on
my own hardware.

# Synchronous Completion

- HTTPS Only (insecure requests return a 521 error)
- The request will take ~30-60s for the LLM to process
- This method does not support streaming

```bash
curl https://slm.robog.net/respond \
  --request POST \
  --header "Content-Type: application/json" \
  --data @- << EOF
    [
      {
        "role": "user",
        "content": "Tell me a joke!"
      }
    ]
EOF
```

# Streamed Completion

`// TODO`

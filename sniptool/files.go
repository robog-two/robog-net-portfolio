package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	outputDir = "./src/s"
	blogDir   = "./src/blog"
)

// --- file operations ----------------------------------------------------

func saveSnippet(body string) (string, error) {
	now := time.Now()
	fileName := fmt.Sprintf("%s-%s.md", now.Format("2006-01-02"), slugify(body))
	header := buildSnippetHeader(now)
	path := outputDir + "/" + fileName

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", err
	}
	if err := os.WriteFile(path, []byte(header+body), 0644); err != nil {
		return "", err
	}
	return path, nil
}

func buildSnippetHeader(t time.Time) string {
	return fmt.Sprintf(`---
date: %04d-%02d-%02d
time: %s
tags: post
snippet: yes
layout: blog.njk
description: yes
---
`, t.Year(), int(t.Month()), t.Day(), formatTime(t))
}

// --- thread management --------------------------------------------------

func saveThreadPost(title, body string) (string, error) {
	now := time.Now()
	fileName := fmt.Sprintf("%s.md", slugify(title))
	path := filepath.Join(blogDir, fileName)

	header := fmt.Sprintf(`---
title: %s
layout: thread.njk
date: %04d-%02d-%02d
tags:
  - post
description: yes
thread: true
---
`, title, now.Year(), int(now.Month()), now.Day())

	if err := os.WriteFile(path, []byte(header+"\n"+body), 0644); err != nil {
		return "", err
	}
	return path, nil
}

func appendThreadEntry(parentSlug, body string) (string, error) {
	now := time.Now()
	fileName := fmt.Sprintf("%s-%s.md", now.Format("2006-01-02"), slugify(body))
	dirPath := filepath.Join(outputDir, parentSlug)
	path := filepath.Join(dirPath, fileName)

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return "", err
	}

	header := fmt.Sprintf(`---
date: %04d-%02d-%02d
time: %s
thread: %s
layout: thread-entry.njk
---
`, now.Year(), int(now.Month()), now.Day(), formatTime(now), parentSlug)

	if err := os.WriteFile(path, []byte(header+"\n"+body), 0644); err != nil {
		return "", err
	}
	return path, nil
}

func scanThreadPosts() ([]threadPostMeta, error) {
	var posts []threadPostMeta

	entries, err := os.ReadDir(blogDir)
	if err != nil {
		return nil, err
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}

		data, err := os.ReadFile(filepath.Join(blogDir, e.Name()))
		if err != nil {
			continue
		}

		if !strings.Contains(string(data), "thread: true") {
			continue
		}

		slug := strings.TrimSuffix(e.Name(), ".md")
		title := extractFrontmatterTitle(string(data))

		posts = append(posts, threadPostMeta{
			slug:  slug,
			title: title,
			path:  filepath.Join(blogDir, e.Name()),
		})
	}

	return posts, nil
}

func extractFrontmatterTitle(content string) string {
	parts := strings.SplitN(content, "---", 3)
	if len(parts) < 2 {
		return "Unknown"
	}

	re := regexp.MustCompile(`(?m)^title:\s*(.+)$`)
	matches := re.FindStringSubmatch(parts[1])
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return "Unknown"
}

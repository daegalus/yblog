---
author: Yulian Kuncheff
date: 2025-03-14T10:00:00Z
draft: false
slug: go-embed-gotcha
title: Go embed.FS requires a file to exist
tags:
  - go
  - til
type: til
---

`embed.FS` in Go will fail at compile time if the directory you're embedding is completely empty. You need at least one file (even a `.gitkeep`) for the embed directive to work. This caught me off guard when setting up new content directories.

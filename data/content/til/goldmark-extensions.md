---
author: Yulian Kuncheff
date: 2025-03-13T15:00:00Z
draft: false
slug: goldmark-extensions
title: Goldmark supports custom inline parsers
tags:
  - go
  - goldmark
  - til
type: til
---

You can add custom inline parsers to Goldmark by implementing the `parser.InlineParser` interface with a `Trigger()` byte and a `Parse()` function. I used this to add `[[wiki-link]]` syntax to my blog's knowledge base. The key is setting the priority higher than the default link parser so `[[` gets matched before `[`.

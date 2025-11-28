# {{ .Title }}

**Author**: {{ .Author }}

Welcome to your new book project! This repository is designed to let you focus on writing while DraftForge handles the formatting and publishing details.

## 🚀 Quick Start

1.  **Write**: Create new chapters in the `chapters/` directory.
    -   Example: `chapters/01-the-beginning.md`
2.  **Preview**: Run `draftforge serve` to see a live preview of your book.
3.  **Publish**: Push your changes to GitHub to automatically generate EPUB and PDF files.

## 📂 Project Structure

-   `chapters/`: **Your Manuscript**. Numbered markdown files (e.g., `01.md`, `02.md`).
-   `docs/`: **World Bible**. Keep track of characters, locations, and research here.
-   `editorial/`: **Editing History**. Archived editorial goals and feedback.
-   `draftforge.yaml`: **Configuration**. Settings for your book's title, author, and build options.
-   `CREATIVE.md`: **Story Bible**. High-level overview of your story's core elements.
-   `EDITORIAL.md`: **Current Focus**. What you are currently working on (e.g., "First Draft", "Fixing Plot Holes").
-   `AGENTS.md`: **AI Context**. Instructions for AI collaborators.

## ✍️ Writing Guide

### Frontmatter

Each chapter file should start with a small metadata block (frontmatter) so DraftForge knows how to handle it:

```yaml
---
title: The Beginning
part: Part 1
pov: Alice
---
```

### Editorial Workflow

1.  **Open a Round**: Define your goal in `EDITORIAL.md`.
2.  **Execute**: Write or edit your chapters.
3.  **Close a Round**: Move `EDITORIAL.md` to `editorial/YYYY-MM-DD-goal.md` and start a new one.

## 🤖 AI Collaboration

This project is "AI-Native". The `docs/` folder and `AGENTS.md` file are designed to provide context to AI tools, helping them understand your story's world and characters.

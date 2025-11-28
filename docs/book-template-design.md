# Book Repository Template Design (DraftForge Edition)

This template represents a "thin client" approach. The repository contains only **content** and **configuration**. All logic (compilation, stats, AI orchestration) is handled by the external `draftforge` CLI or GitHub Actions.

**Delivery Mechanism**: These files will be embedded in the DraftForge binary (using Go `embed`) and generated dynamically via `df init`.

## 1. Directory Structure

The template initializes with a clean, content-focused structure.

```text
/
├── .github/
│   └── workflows/
│       └── draftforge.yaml     # CI/CD pipeline using DraftForge Actions
├── chapters/ *                # Source of truth. Numbered markdown files.
│   └── .gitkeep
├── stats/                     # JSON data (optional, can be git-ignored or committed for history)
│   └── .gitkeep
├── build/                     # Git-ignored output folder (epub, pdf, html)
├── docs/                      # World bible and research
│   ├── characters/
│   └── world/
├── .gitignore                 # Standard ignores
├── draftforge.yaml *          # Project configuration (replaces Taskfile logic)
├── README.md *                # Project root & quickstart
├── REPO-STRUCTURE.md          # Technical manual for the repo
├── CREATIVE.md                # The "Story Bible" stub
├── EDITORIAL.md               # Current editorial focus
├── editorial/                 # Archived editorial rounds
│   └── .gitkeep
└── AGENTS.md                  # Context for AI collaborators

## 2. Essential Boilerplate Files

### `draftforge.yaml`
The central configuration file. Generated dynamically with user input.
```yaml
project:
  title: "{{ .Title }}"   # Injected during generation
  author: "{{ .Author }}" # Injected during generation

structure:
  chapters_dir: "chapters"
  stats_dir: "stats"
  output_dir: "build"

workflows:
  stats: true
  epub: true
  pdf: false # Optional
```

### `.github/workflows/draftforge.yml`
Standard workflow that delegates to DraftForge's reusable actions.
```yaml
name: DraftForge Pipeline
on: [push]
jobs:
  build:
    uses: draftforge/actions/.github/workflows/pipeline.yml@v1
    with:
      config: draftforge.yaml
```

### `README.md`
Focused on writing, not tooling.
- **Sections**:
    - Title & Logline
    - "How to Write" (pointing to `chapters/`)
    - "How to Preview": "Run `df serve` or push to GitHub."

### `CREATIVE.md` & `AGENTS.md`
Remains the same. These are critical for AI context and human collaboration.

### `EDITORIAL.md`
Defines the *current* focus of the project.
- **Content**:
    - "Current Round: First Draft / Pacing Pass / Character Consistency"
    - "Goals: Fix plot holes in Ch 1-5"
    - "Anti-Goals: Don't worry about line editing yet."

## 3. Core Patterns & Guidance

### The "No-Code" Author Experience
- **Pattern**: Users never see Go code or shell scripts.
- **Workflow**:
    1.  Write in `chapters/`.
    2.  Push to GitHub -> Actions auto-generate stats/ebooks.
    3.  (Optional) Run `draftforge start` locally for a live preview.

### The Editorial Migration Pattern
Treat editing like database migrations. Never edit aimlessly; edit with a specific schema change in mind.
- **Workflow**:
    1.  **Open Round**: Create `EDITORIAL.md` defining the focus (e.g., "Fixing Plot Holes").
    2.  **Execute**: Make changes across chapters.
    3.  **Close Round**: Move `EDITORIAL.md` to `editorial/YYYY-MM-DD-fixing-plot-holes.md`.
    4.  **Repeat**: Start a new `EDITORIAL.md` for the next pass.
- **Benefit**: Preserves the history of *why* changes were made, not just *what* changed.

### The Chapter Frontmatter Spec
Strict adherence to frontmatter is still required for the platform to understand the content.
```yaml
---
title: Chapter Title
part: I
pov: character_name
---
```

### AI Context
- **Pattern**: The `draftforge` CLI/Action will automatically aggregate `docs/` and `chapters/` when prompting AI agents, so the file structure *is* the context window.

## 4. Implementation Steps for DraftForge API

1.  **Create `templates/` Package**: Inside the DraftForge Go codebase.
2.  **Embed Files**: Use `//go:embed *` to bundle the template files.
3.  **Implement `df init`**:
    - Prompt user for `Title` and `Author`.
    - Render `draftforge.yaml` template with these values.
    - Write all other static files (`.gitignore`, `README.md`, etc.) to the target directory.
    - Initialize `git init`.

## 5. Proposed "First Run" Experience

1.  User runs `df init my-new-book`.
2.  CLI asks: "What is the title?" -> "The Great Novel".
3.  CLI asks: "Author name?" -> "Jane Doe".
4.  CLI generates folder `my-new-book/` with all files pre-configured.
5.  User `cd my-new-book` and starts writing.

## 6. Benefits & Author Coaching

### Why This Process? (The "Pitch")
1.  **Zero-Config Start**: Authors don't need to learn Git, Makefiles, or Pandoc. They just type `df init` and start writing Markdown.
2.  **Ownership**: The author owns the text files on their local machine. No vendor lock-in to a SaaS platform.
3.  **Invisible DevOps**: CI/CD happens magically. They push code, and a PDF appears. It feels like a service, but it's just a repo.
4.  **AI-Native Foundation**: By structuring data correctly from Day 1 (frontmatter, `docs/` folder), the project is "pre-indexed" for any future AI agent they want to hire.

### Coaching Strategies

#### Phase 1: The "Just Write" Phase
*   **Goal**: Get them comfortable with Markdown and the file structure.
*   **Coach**: "Don't worry about the `build` folder. Just create `chapters/01.md` and write. The system will handle the rest."
*   **Habit**: Encourage them to commit/push daily, framing it as "saving to the cloud" rather than "version control."

#### Phase 2: The "Feedback Loop" Phase
*   **Goal**: Introduce the stats and build artifacts.
*   **Coach**: "Check the `stats/` folder (or the dashboard). See how your word count is trending? That's your velocity."
*   **Gamification**: Use the automated stats to celebrate milestones (e.g., "You hit 10k words!").

#### Phase 3: The "Gardener" Phase
*   **Goal**: expand `docs/` and `CREATIVE.md`.
*   **Coach**: "As you invent characters, give them a file in `docs/characters/`. This isn't just for you; it's teaching the AI to act as your continuity editor."
*   **Metaphor**: The repo is a garden. The story is the flower, but the `docs/` are the soil that supports it.

#### Phase 4: The "Director" Phase (Editorial Migrations)
*   **Goal**: Structured editing using `EDITORIAL.md`.
*   **Coach**: "Don't try to fix everything at once. Open an 'Editorial Round' for *just* dialogue. Finish that, archive it, then start a round for *pacing*."
*   **Shift**: The author moves from just a writer to a creative director, managing the project's evolution in distinct, committable stages.

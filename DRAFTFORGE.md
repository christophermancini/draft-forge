# DraftForge Platform Blueprint

This document captures the concept, requirements, and recommended architecture for DraftForge—the platform that walks authors through the end-to-end book-production process, while transparently storing work in GitHub and orchestrating automations through GitHub Actions (GHA).

The goal is a UX where non-technical writers, editors, agents, and production partners collaborate without touching Git/YAML, yet all artifacts, workflows, and rights remain under the author’s control.

---

## 1. Core Principles

1. **User-owned repos** – Each project lives in the author’s GitHub repo. DraftForge is merely a client/service layer.
2. **Transparent automation** – All builds, exports, AI runs happen via GitHub Actions. DraftForge provides recipes and dashboards.
3. **Composable workflow** – Chapters, metadata, legal documents, AI prompts, and exports follow conventions (e.g., Markdown + metadata blocks) for easy tooling.
4. **AI-augmented, not AI-owned** – AI assists in drafting, analysis, exports, but every action is versioned and attributable.
5. **Open exit** – Users can leave at any time with a plain Git repo containing all content and workflow definitions.

---

## 2. User Roles & Journeys

### Authors
- Create/import manuscripts (Markdown-first).
- Request stats, exports (HTML/EPUB/PDF/DOCX), AI suggestions, or legal packets via one-click UI.
- Track tasks (e.g., “Chapter 12 revision,” “Legal docs prepared”).

### Editors/Agents
- Comment/review directly in DraftForge UI while committing suggestions to branches or PRs.
- Approve/reject workflow runs (e.g., release packaging).

### Tech/Production Partners
- Define new workflows (e.g., SSML generation, audiobook scripts, timeline builds).
- Monitor resource usage / GHA credits.

---

## 3. Features & Modules

### 3.1 Content Layer
- Markdown chapters with frontmatter (title, POV, entry/exit hooks, tags).
- Ancillary docs (legal evidence, timelines, AI transcripts) stored under `/docs`.
- Metadata index (JSON/YAML) describing book structure for quick reference.

### 3.2 Workflow Engine (GitHub Actions)
- **Stats** – run `cmd/generate-stats` (as seen in Immortal) after chapter changes.
- **Lint/Format** – dprint or prettier for Markdown + code.
- **Compile** – Pandoc-based pipeline to create `build/book.md`, EPUB, PDF, HTML (already prototyped), optional DOCX/MOBI.
- **Publish** – optionally auto-upload artifacts to S3/Netlify/BookFunnel when a release tag is pushed.
- **AI Tasks** – templated prompts for editing suggestions, continuity checks, query-letter drafts, etc.

### 3.3 UX Layer (DraftForge Web App)
- **Project Dashboard** – status of workflows, build artifacts, outstanding edits.
- **Chapter Editor** – WYSIWYM (what-you-see-is-what-you-mean) Markdown editor with AI assist.
- **Review Mode** – in-browser diff + comment threads tied to Git branches/PRs.
- **Automation Catalog** – toggle workflows (stats, compile, SSML export); config gets committed as GitHub Action files.
- **AI Request Console** – submit structured prompts with references (chapters, characters), store outputs in repo.

### 3.4 Collaboration
- Role-based invites (author, editor, agent, legal), mapped to GitHub org permissions.
- Activity log showing commits, AI runs, workflow statuses.
- Notifications (email/SMS/slack) for major events (legal builds ready, AI review done, exporter finished).

### 3.5 Legal & Release Toolkit
- Prebuilt templates for legal filing packets (like Rachel/Eleanor’s case).
- Evidence bundler: collects stats, metadata, neural scans, statements into exportable PDF/HTML.
- Query/pitch builder with AI assistance plus tracked revisions.

### 3.6 Audio & Advanced Exports (Phase 2)
- Audiobook script generator (plain text or SSML).
- TTS integration via SSML exports; per-chapter audio output artifact.
- Timeline visualizer (Scrivener-style corkboard) generated from metadata.

---

## 4. Architecture Overview

```
User UI (DraftForge Web App)
   ├── Auth (GitHub OAuth) + role management
   ├── Markdown editor + AI assistant (calls backend)
   ├── Workflow dashboard (reads GitHub status APIs)
   ├── Artifact viewer/downloader (pull from GitHub releases/artifacts)

DraftForge Backend
   ├── Stores project configuration (which workflows enabled, AI prompt templates)
   ├── Manages OAuth tokens to interact with GitHub repos & Actions
   ├── Queues AI requests, writes results back to repo (PR or branch)
   └── Provides public URLs for artifacts if needed (optional CDN)

GitHub Repo (per project)
   ├── /chapters/*.md + metadata
   ├── /docs (legal, research, AI transcripts)
   ├── /cmd (stats, compile scripts as in Immortal)
   ├── /Taskfile or /gha workflows (managed via DraftForge)
   └── /build (outputs) – ignored except when uploading artifacts
```

---

## 5. Prompt / Context Document for AI Agents

When engaging another AI (coding agent, UX designer, etc.), provide:

### Prompt Header
```
System: You are DraftForge’s engineering/design agent. You extend a GitHub-based novel-writing platform that orchestrates workflows via GitHub Actions. All output must be repo-ready (docs, code, workflows) and respect user ownership.
```

### Context Block (example)
```
Context:
- Repos contain Markdown chapters with YAML frontmatter (see /chapters)
- Automation pipeline uses Taskfile + cmd/ scripts (stats, compile, etc.)
- Users interact via DraftForge UI; GitHub remains source of truth
- Workflows must be transparent: commit YAML to .github/workflows/
- New features should be toggled via DraftForge config but ultimately reside in repo
- AI outputs (edits, summaries) are saved as Markdown under /ai or inline comments
- Access control maps back to GitHub teams/roles
- Future phases include audio exports (SSML) and AI editorial suggestions

Current Sprint Objectives:
1. Add HTML exporter (Pandoc) – ✅ done (example commit)
2. Define audiobook script pipeline (text + SSML)
3. Build AI comment bot that reviews PRs for continuity issues
```

### Ask to Agent
```
Tasks:
- Create `audiobook` Taskfile target that outputs plain-text + SSML per chapter.
- Add `.github/workflows/audiobook.yml` calling the new task on push to main.
- Document usage in README (new section “Audiobook Exports”).
Deliverables:
- Modified Taskfile
- New workflow file
- README update
```

---

## 6. Phase Plan

1. **MVP (Weeks 0-6)**
   - GitHub OAuth integration
   - Repo import/creation flow
   - Chapter editor with AI assist (basic prompt templates)
   - Workflow dashboard reading existing Actions (stats, compile)

2. **Expansion (Weeks 6-12)**
   - Automation catalog (enable/disable export formats, AI jobs)
   - Collaboration (comments, review assignments)
   - Legal toolkit (evidence bundler template)

3. **Advanced (Weeks 12-20)**
   - Audio/SSML pipeline
   - Timeline visualizer
   - Marketplace for custom workflows (community-contributed Actions)

4. **General Release**
   - Pricing tied to GHA usage (pass-through or prepaid credits)
   - Support & documentation portal
   - API for third-party tools (e.g., Scrivener importer, Vellum export)

---

## 7. Key Open Questions
- How to sandbox drafts for AI without exposing sensitive data? (e.g., optional “air-gapped” mode using local actions)
- Licensing for AI models integrated into platform (OpenAI, Anthropic, self-hosted LLM)
- Monetization beyond GHA credits (premium templates, dedicated support)?
- Handling large media assets (audio, video) – store in Git LFS, S3, or user-provided storage?
- Compliance needs for legal docs (chain of custody, signatures)

---

This blueprint should be enough for another AI or dev team to spin up the system without additional context, while keeping the user-centric, GitHub-backed philosophy intact.

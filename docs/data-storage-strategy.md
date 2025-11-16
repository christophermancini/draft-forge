# Data Storage Strategy

This document outlines where different types of data are stored in DraftForge and the reasoning behind these decisions.

---

## 🎯 Core Principle: Authors Own Everything

**All content and project-specific data lives in the author's GitHub repository.**
**DraftForge's database serves as a performance cache and coordination layer.**

---

## 📊 Storage Locations

### 1. **Author's Repository** (Source of Truth)

#### Manuscript Content
**Location:** `chapters/`, `manuscript/`
**Format:** Markdown files
**Ownership:** Author

```
chapters/
├── 01-the-beginning.md
├── 02-the-journey.md
└── ...
```

**Why in repo:**
- ✅ Version controlled
- ✅ Portable and exportable
- ✅ Works offline
- ✅ No vendor lock-in
- ✅ Standard Markdown

---

#### Project Statistics
**Location:** `.draftforge/stats.json`
**Format:** JSON
**Ownership:** Author
**Update frequency:** On every push

```json
{
  "version": "1.0.0",
  "project": {
    "name": "My Novel",
    "type": "novel",
    "target_word_count": 80000
  },
  "current": {
    "word_count": 45678,
    "chapter_count": 12,
    "character_count": 234567,
    "last_updated": "2025-10-29T14:30:00Z",
    "last_commit_sha": "abc123def456"
  },
  "breakdown": {
    "chapters": [
      {
        "file": "chapters/01-the-beginning.md",
        "title": "The Beginning",
        "word_count": 3850,
        "character_count": 19234,
        "last_modified": "2025-10-28T10:15:00Z"
      }
    ],
    "manuscript": [
      {
        "file": "manuscript/outline.md",
        "word_count": 1250
      }
    ]
  },
  "history": [
    {
      "date": "2025-10-29",
      "word_count": 45678,
      "chapter_count": 12,
      "commits": 3,
      "words_written": 1250,
      "commit_sha": "abc123def456"
    },
    {
      "date": "2025-10-28",
      "word_count": 44428,
      "chapter_count": 12,
      "commits": 2,
      "words_written": 890,
      "commit_sha": "def456abc123"
    }
  ],
  "streaks": {
    "current_days": 7,
    "longest_days": 21,
    "last_write_date": "2025-10-29",
    "total_writing_days": 87
  },
  "generated_by": "draftforge-stats-action",
  "generated_at": "2025-10-29T14:30:00Z"
}
```

**Why in repo:**
- ✅ Portable with project
- ✅ Version controlled (git log shows progress)
- ✅ Can view without DraftForge
- ✅ Enables offline analysis
- ✅ Authors own their metrics

**Update workflow:**
1. Author pushes commit
2. GitHub Action counts words
3. Action updates `.draftforge/stats.json`
4. Action commits stats file
5. Webhook notifies DraftForge API
6. API syncs to database

---

#### Project Configuration
**Location:** `.draftforge/config.yml`
**Format:** YAML
**Ownership:** Author

```yaml
version: "1.0"

project:
  name: "My Novel"
  type: novel
  genre: fantasy
  language: en
  target_word_count: 80000

structure:
  chapters_directory: chapters
  manuscript_directory: manuscript
  assets_directory: assets

agents:
  continuity:
    enabled: true
    model: claude-3-5-sonnet
    trigger: pr
    context_window: 5
    auto_comment: true

  style:
    enabled: true
    model: gpt-4o-mini
    trigger: commit
    checks:
      - tense_consistency
      - pov_shifts
      - readability

  timeline:
    enabled: false

  fact:
    enabled: false

compilation:
  pandoc_template: default
  include_cover: true
  auto_build_on_tag: true
  formats:
    - epub
    - pdf

github_actions:
  stats_enabled: true
  compile_enabled: true
  ai_review_enabled: true
```

**Why in repo:**
- ✅ Version controlled (can revert settings)
- ✅ Portable with project
- ✅ Human-readable and editable
- ✅ Can be customized per branch

---

#### Agent Results
**Location:** `.draftforge/agent-runs/`
**Format:** JSON files (one per run)
**Ownership:** Author

```
.draftforge/agent-runs/
├── 2025-10-29_14-30-00_continuity_pr-123.json
├── 2025-10-29_10-15-00_style_commit-abc123.json
└── ...
```

**Example file:**
```json
{
  "run_id": "df-run-12345",
  "agent_type": "continuity",
  "trigger": "pr",
  "trigger_ref": "123",
  "timestamp": "2025-10-29T14:30:00Z",
  "model": "claude-3-5-sonnet",
  "status": "completed",
  "summary": "Found 2 continuity issues",
  "issues": [
    {
      "id": "cont-001",
      "severity": "error",
      "category": "character_consistency",
      "file": "chapters/05-chapter-five.md",
      "line": 45,
      "message": "Character 'Sarah' has blue eyes here, but brown eyes in Chapter 3.",
      "suggestion": "Update to match earlier description or explain the change."
    }
  ],
  "stats": {
    "files_analyzed": 5,
    "total_words": 12500,
    "tokens_used": 12500,
    "cost_cents": 25,
    "duration_seconds": 8
  }
}
```

**Why in repo:**
- ✅ Complete audit trail
- ✅ Can review past agent feedback
- ✅ Portable with project
- ✅ Version controlled

---

#### Metadata
**Location:** `manuscript/metadata.yml`
**Format:** YAML
**Ownership:** Author

```yaml
title: "My Epic Novel"
subtitle: "A Tale of Adventure"
author:
  name: "Jane Author"
  email: "jane@example.com"
  website: "https://janeauthor.com"

series:
  name: "The Great Saga"
  number: 1
  total: 3

publishing:
  isbn: "978-0-123456-78-9"
  edition: "First Edition"
  copyright_year: 2025
  publisher: "Self-Published"

categories:
  - Fiction
  - Fantasy
  - Adventure

keywords:
  - magic
  - dragons
  - quest

description: |
  A thrilling adventure story about...

dedication: "For my family"

cover_image: "../assets/cover.png"
```

**Why in repo:**
- ✅ Used for EPUB/PDF generation
- ✅ Version controlled
- ✅ Portable
- ✅ Standard format

---

### 2. **DraftForge Database** (Performance Cache & Coordination)

#### User Accounts
**Table:** `users`
**Contains:** GitHub OAuth data, preferences, profile

**Why in database:**
- ✅ Authentication requires central storage
- ✅ Fast user lookup
- ✅ Not project-specific

---

#### Project Registry
**Table:** `projects`
**Contains:** Project metadata, GitHub repo links, settings cache

**Why in database:**
- ✅ Fast project listing for dashboard
- ✅ Enables search and filtering
- ✅ Links users to their repos

**Sync strategy:**
- Database is synced FROM `.draftforge/config.yml`
- Repository is source of truth
- Database caches for performance

---

#### Stats Cache
**Table:** `project_stats`
**Contains:** Cached word counts, streaks, recent history

**Why in database:**
- ✅ Fast dashboard queries
- ✅ Aggregations and trends
- ✅ Real-time updates

**Sync strategy:**
1. GitHub Action updates `.draftforge/stats.json`
2. Webhook triggers sync to database
3. Dashboard reads from database
4. If database unavailable, fall back to reading from repo

---

#### AI Credits & Billing
**Table:** `ai_credits`, `ai_usage_log`
**Contains:** Token balances, subscription tiers, usage history

**Why in database:**
- ✅ Real-time credit checking
- ✅ Billing and invoicing
- ✅ Usage analytics
- ✅ Rate limiting
- ✅ Not project-specific

**NOT in repository:** Billing data is service-level, not project-specific

---

#### Agent Run Queue
**Table:** `agent_runs`
**Contains:** Queued/running agent jobs, execution status

**Why in database:**
- ✅ Job queue management
- ✅ Real-time status updates
- ✅ Priority and retry logic

**Sync strategy:**
- Results written to `.draftforge/agent-runs/` after completion
- Database tracks execution state
- Repository stores final results

---

#### Collaborators
**Table:** `collaborators`
**Contains:** Multi-user access, invitations, permissions

**Why in database:**
- ✅ Access control (security-sensitive)
- ✅ Real-time permission checks
- ✅ Invitation system

**Could be in repo:** Future consideration for `.draftforge/collaborators.yml`

---

#### Export Jobs
**Table:** `export_jobs`
**Contains:** EPUB/PDF generation queue and status

**Why in database:**
- ✅ Job queue management
- ✅ Output storage links (Cloudflare R2)
- ✅ Temporary (auto-expire after 30 days)

---

#### Notifications
**Table:** `notifications`
**Contains:** In-app alerts, unread status

**Why in database:**
- ✅ Real-time delivery
- ✅ Read/unread tracking
- ✅ Temporary (not historical)

---

#### Webhooks
**Table:** `webhooks`, `webhook_deliveries`
**Contains:** Integration configs, delivery logs

**Why in database:**
- ✅ Service-level configuration
- ✅ Delivery tracking and retry logic
- ✅ Not project-specific

---

## 🔄 Sync Patterns

### **Repository → Database (Pull)**
Triggered by GitHub webhooks:

1. **Push event:**
   - Sync `.draftforge/config.yml` → `projects.settings`
   - Sync `.draftforge/stats.json` → `project_stats`

2. **File change detection:**
   - Only sync if `.draftforge/` files changed
   - Validate JSON/YAML before syncing

### **Database → Repository (Push)**
Triggered by user actions:

1. **Settings changed in UI:**
   - Update `projects.settings` in database
   - Create commit to update `.draftforge/config.yml`
   - Use GitHub API to commit

2. **Agent run completed:**
   - Store summary in database
   - Create commit with `.draftforge/agent-runs/*.json`

---

## 🎯 Decision Matrix

| Data Type | Repository | Database | Reason |
|-----------|-----------|----------|--------|
| **Manuscript content** | ✅ Primary | ❌ | Author ownership, version control |
| **Project stats** | ✅ Primary | ✅ Cache | Portability + performance |
| **Project config** | ✅ Primary | ✅ Cache | Portability + performance |
| **Agent results** | ✅ Archive | ✅ Current | Audit trail + performance |
| **User accounts** | ❌ | ✅ Only | Authentication requirement |
| **AI credits** | ❌ | ✅ Only | Billing, not project-specific |
| **Collaborators** | ❌* | ✅ Only | Security, access control |
| **Export outputs** | ❌** | ✅ Only | Temporary, large files |
| **Notifications** | ❌ | ✅ Only | Ephemeral, real-time |
| **Webhooks** | ❌ | ✅ Only | Service-level config |

\* Could add `.draftforge/collaborators.yml` in future
\*\* Outputs stored in Cloudflare R2, not GitHub

---

## 🔐 Data Ownership Philosophy

### **Author Owns:**
- All manuscript content
- All project-specific data (stats, config, agent results)
- All metadata (title, author, ISBN, etc.)

### **DraftForge Owns:**
- User authentication records
- Service usage (AI credits, billing)
- Platform-level data (notifications, webhooks)

### **Shared/Synced:**
- Project registry (repo is source of truth)
- Project stats (repo is source of truth, DB caches)
- Project settings (repo is source of truth, DB caches)

---

## 💡 Benefits of This Approach

1. **Portability:** Authors can export their repo and continue elsewhere
2. **Transparency:** All project data visible in git history
3. **Offline Work:** Stats and config available locally
4. **Version Control:** Settings and stats are versioned
5. **Performance:** Database caching for fast queries
6. **No Lock-in:** DraftForge is optional after initial setup
7. **Data Ownership:** Authors control their content and metrics

---

## 🚀 Implementation Priority

### Phase 1 (MVP)
- ✅ `.draftforge/config.yml` (project settings)
- ✅ `.draftforge/stats.json` (word counts, streaks)
- ✅ Database sync via webhooks
- ✅ GitHub Action for stats generation

### Phase 2
- ✅ `.draftforge/agent-runs/` (agent results archive)
- ✅ Bi-directional sync (UI changes → repo)
- ✅ Conflict resolution (repo wins)

### Phase 3
- ✅ `.draftforge/collaborators.yml` (optional)
- ✅ `.draftforge/exports/` (build artifacts)
- ✅ Offline mode (work without API)

---

## 📝 File Format Standards

### JSON Schema Validation
All `.draftforge/*.json` files validated against JSON schemas:
- `stats.schema.json`
- `agent-run.schema.json`

### YAML Schema Validation
All `.draftforge/*.yml` files validated against YAML schemas:
- `config.schema.yml`
- `metadata.schema.yml`

Schemas stored in `scaffolds/common/schemas/` and deployed with projects.

---

**Last Updated:** October 29, 2025
**Version:** 1.0.0
**Status:** Approved

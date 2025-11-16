# CLAUDE.md

## Project: DraftForge
**Tagline:** Forge your draft. Keep your voice.

**Purpose:** An AI-assisted authoring platform for writers who want control over their craft. DraftForge helps authors organize, refine, and publish long-form creative works (novels, screenplays, non-fiction) using Markdown, Git, and AI editorial agents—without sacrificing their authentic voice.

---

## 🎯 Core Philosophy

**Authors Own Everything**
- Every project lives in the author's GitHub account
- No vendor lock-in—all files are standard Markdown
- Authors can export and continue elsewhere at any time

**AI Assists, Never Replaces**
- AI acts as editorial assistant: critique, suggestions, consistency checks
- Authors make all creative decisions
- No "generate chapter" features—we support authentic creativity

**Git-Native Architecture**
- Version control is built-in, not bolted-on
- GitHub Actions automate the boring parts (stats, builds, linting)
- Branching for drafts, PRs for editorial review

**Multi-Model Flexibility**
- Authors choose AI models based on task and budget
- Support for Claude, GPT-4, Gemini, and local models
- BYOK (Bring Your Own Key) option for advanced users

---

## 🏗️ Architecture Overview

### Tech Stack
| Layer | Technology | Rationale |
|-------|-----------|-----------|
| **Frontend** | SvelteKit + Tailwind | Fast, modern, great DX |
| **Backend** | Go (Fiber) + PostgreSQL | Performance, type safety, easy deployment |
| **Storage** | Cloudflare R2 | Cost-effective object storage for exports |
| **Auth** | GitHub OAuth + GitHub App | Seamless GitHub integration |
| **AI Routing** | OpenRouter API | Multi-model support without vendor lock-in |
| **CI/CD** | GitHub Actions | Already in author's workflow |
| **Editor** | Monaco (VS Code editor) | Familiar, powerful, extensible |

### System Components

```
┌─────────────────────────────────────────────────────┐
│                 DraftForge Platform                 │
│                                                     │
│  ┌─────────────┐         ┌──────────────┐         │
│  │   Web UI    │ ←─────→ │   Go API     │         │
│  │ (SvelteKit) │         │   (Fiber)    │         │
│  └─────────────┘         └──────────────┘         │
│         │                        │                  │
│         │                        ↓                  │
│         │                 ┌──────────────┐         │
│         │                 │  PostgreSQL  │         │
│         │                 └──────────────┘         │
│         │                                           │
│         ↓                                           │
│  ┌─────────────────────────────────────┐          │
│  │      Author's GitHub Repo           │          │
│  │  ┌─────────────────────────────┐   │          │
│  │  │ chapters/                   │   │          │
│  │  │ manuscript/                 │   │          │
│  │  │ .github/workflows/          │   │          │
│  │  │ .draftforge/                │   │          │
│  │  └─────────────────────────────┘   │          │
│  └─────────────────────────────────────┘          │
│         │                                           │
│         ↓                                           │
│  ┌──────────────┐         ┌──────────────┐        │
│  │ GitHub       │         │  OpenRouter  │        │
│  │ Actions      │ ←─────→ │  (AI Models) │        │
│  └──────────────┘         └──────────────┘        │
└─────────────────────────────────────────────────────┘
```

---

## 📁 Repository Structure

### Platform Repository (This Repo)
```bash
draftforge/
├── cmd/
│   ├── api/             # Main API server
│   └── cli/             # CLI tools for local use
├── internal/
│   ├── auth/            # GitHub OAuth handling
│   ├── projects/        # Project CRUD
│   ├── ai/              # AI model routing
│   ├── github/          # GitHub API interactions
│   ├── db/              # Database utilities & migrations
│   └── scaffold/        # Project scaffolding logic
├── frontend/            # SvelteKit web app
│   ├── src/
│   │   ├── routes/      # Pages
│   │   ├── lib/         # Components & utilities
│   │   └── app.html
│   └── package.json
├── agents/              # AI agent definitions
│   ├── continuity/
│   ├── style/
│   ├── timeline/
│   └── README.md
├── scaffolds/           # Project templates
│   ├── novel/
│   ├── screenplay/
│   ├── technical-book/
│   └── non-fiction/
├── infra/               # Infrastructure (Terraform, Docker, etc.)
├── .github/
│   └── workflows/       # Platform CI/CD
├── docs/
│   ├── architecture.md
│   ├── getting-started.md
│   └── ai-agents.md
├── go.mod               # Go dependencies
├── Taskfile.yaml        # Development tasks
├── CLAUDE.md            # This file
└── README.md
```

### User Project Structure (What Gets Created)
```bash
my-novel/                          # Author's GitHub repo
├── chapters/
│   ├── 01-the-beginning.md
│   ├── 02-the-journey.md
│   └── ...
├── manuscript/
│   ├── metadata.yml               # Title, author, ISBN, etc.
│   ├── outline.md
│   └── character-bible.md
├── assets/
│   ├── cover.png
│   └── images/
├── .draftforge/
│   ├── config.yml                 # Project settings
│   ├── agents.yml                 # Which agents to run, when
│   └── templates/                 # Custom export templates
├── .github/
│   └── workflows/
│       ├── stats.yml              # Word count, chapter tracking
│       ├── compile.yml            # Build EPUB/PDF
│       └── ai-review.yml          # Run AI agents on PRs
├── outputs/                       # Git-ignored build artifacts
│   ├── my-novel.epub
│   └── my-novel.pdf
└── README.md                      # Auto-generated project info
```

---

## 🤖 AI Agent System

### Agent Philosophy
- **Domain-specific:** Each agent has one clear responsibility
- **Cost-conscious:** Only analyze changed content, not entire manuscript
- **Configurable:** Authors enable/disable agents per project
- **Transparent:** All suggestions are logged and reviewable

### Core Agents (MVP)

#### 1. ContinuityBot
**Model:** Claude 3.5 Sonnet (best reasoning for complex context)
**Triggers:** On PR creation, on-demand
**Checks:**
- Character consistency (appearance, traits, backstory)
- World rules (magic systems, technology, geography)
- Plot threads (unresolved storylines, contradictions)

**Example Output:**
```markdown
## Continuity Issues Found

### Chapter 12, Line 45
❌ Character "Sarah" has blue eyes here, but brown eyes in Chapter 3.

### Chapter 15, Line 102
⚠️  The magic system established in Chapter 2 states spells require
verbal components, but this spell is cast silently.
```

#### 2. StyleBot
**Model:** GPT-4o-mini (fast, good for stylistic patterns)
**Triggers:** On every commit to draft branch
**Checks:**
- Voice consistency (tense shifts, POV breaks)
- Sentence variety and rhythm
- Overused words or phrases
- Readability metrics (Flesch-Kincaid, etc.)

#### 3. TimelineBot
**Model:** Gemini 1.5 Pro (excellent at structured analysis)
**Triggers:** On-demand or weekly
**Checks:**
- Chronological consistency
- Scene pacing (too fast/slow)
- Chapter length balance
- Time elapsed between events

#### 4. FactBot
**Model:** GPT-4.1 Turbo (best for factual accuracy)
**Triggers:** On-demand (for non-fiction/historical)
**Checks:**
- Historical accuracy
- Geographic correctness
- Technical plausibility
- Citation verification (for non-fiction)

### Agent Configuration

Authors configure agents in `.draftforge/agents.yml`:

```yaml
agents:
  continuity:
    enabled: true
    model: claude-3-5-sonnet
    trigger: pr
    context_window: 5  # chapters before/after

  style:
    enabled: true
    model: gpt-4o-mini
    trigger: commit
    checks:
      - tense_consistency
      - pov_shifts
      - readability

  timeline:
    enabled: false  # Disabled by author

  fact:
    enabled: true
    model: gpt-4-turbo
    trigger: manual
```

---

## 🚀 Development Roadmap

### Phase 1: Core Platform (Months 1-6)
**Goal:** Prove concept with technical writers

**Features:**
- [ ] GitHub OAuth + App integration
- [ ] Project scaffolding (novel, technical book templates)
- [ ] Web-based Markdown editor with preview
- [ ] Basic AI review (ContinuityBot + StyleBot)
- [ ] GitHub Actions for word count stats
- [ ] EPUB/PDF compilation via Pandoc
- [ ] User dashboard (projects, stats, credits)

**Success Metrics:**
- 100-500 active users
- 15% free-to-paid conversion
- <$0.10 cost per AI review
- 90%+ uptime

### Phase 2: Enhanced Workflow (Months 7-12)
**Goal:** Expand to power users and fiction authors

**Features:**
- [ ] Multi-model AI routing with cost comparison
- [ ] Timeline visualization (interactive scene graph)
- [ ] Character database with relationship mapping
- [ ] Collaborative editing (invite co-authors/editors)
- [ ] Version comparison UI (visual diff)
- [ ] Export theme marketplace
- [ ] Mobile-optimized editor
- [ ] Offline mode (PWA)

**Success Metrics:**
- 1,000-5,000 active users
- 20% conversion rate
- 85%+ gross margin
- 12+ month median retention

### Phase 3: Advanced Features (Year 2+)
**Goal:** Category leadership and vertical expansion

**Features:**
- [ ] Series management (multi-book projects)
- [ ] Screenplay and stage play templates
- [ ] Advanced citation and sourcing for non-fiction
- [ ] Direct KDP/IngramSpark export
- [ ] Custom agent training (user-specific style guides)
- [ ] Audiobook script generation (TTS integration)
- [ ] Translation pipeline
- [ ] "Verified human-written" certification

**Success Metrics:**
- 10,000+ active users
- Category leadership position
- Profitability
- Strong community and word-of-mouth

---

## 💰 Business Model

### Pricing Tiers

**Free (Hobbyist)**
- 1 active project
- 10,000 AI tokens/month (~5 reviews)
- Basic agents (Style, Continuity)
- Community support
- EPUB/PDF exports

**Creator ($19/month or $15/month annual)**
- 5 active projects
- 200,000 AI tokens/month (~100 reviews)
- All agents
- Priority support
- Advanced exports (mobi, custom templates)
- Version history (90 days)

**Professional ($49/month or $39/month annual)**
- Unlimited projects
- 750,000 AI tokens/month (unlimited for some models)
- All features
- Collaboration (5 team members)
- Custom agent training
- BYOK option
- Priority support (24h response)
- Version history (unlimited)

**Enterprise (Custom)**
- White-label options
- Custom AI models
- On-premise deployment
- SSO integration
- Dedicated support

### Credit System
- Additional tokens: $10 per 100,000
- Credits don't expire for paid users
- Transparent cost breakdown per model

---

## 📊 Key Metrics & Goals

### Technical Metrics
- **API Response Time:** p95 < 200ms
- **Uptime:** 99.9%+
- **AI Review Cost:** <$0.10 per review on average
- **Build Success Rate:** 98%+ (EPUB/PDF compilation)

### Business Metrics
- **Gross Margin:** >85% target
- **CAC:** <$30 (organic + paid)
- **LTV:CAC:** >3:1
- **Free-to-Paid Conversion:** 15-20%
- **Monthly Churn:** <5%
- **NPS:** >50

### Product Metrics
- **Time to First Project:** <5 minutes
- **Chapters Written per Month:** Average 10+ (engaged users)
- **AI Reviews per User:** 5-15/month
- **Weekly Active Users:** 60%+ of registered users

---

## 🎨 Design Principles

### User Experience
1. **Invisible Complexity:** Git should be transparent to non-technical users
2. **Progressive Disclosure:** Show advanced features only when needed
3. **Fast Feedback:** AI reviews should feel near-instant (<10 seconds)
4. **Delight in Details:** Small animations, helpful tooltips, smart defaults

### Code Quality
1. **Type Safety:** Use Go's type system; TypeScript on frontend
2. **Testing:** >80% coverage on critical paths
3. **Documentation:** Every public API has godoc/JSDoc
4. **Performance:** Profile regularly; optimize hot paths

### AI Integration
1. **Cost Transparency:** Show users what each review costs
2. **Model Choice:** Let users pick models for different tasks
3. **Graceful Degradation:** Fall back to cheaper models if quota exceeded
4. **Privacy First:** Never train on user content without explicit opt-in

---

## 🧪 Development Guidelines

### When Working on This Project

**DO:**
- Write clean, idiomatic Go and modern JavaScript/Svelte
- Follow established patterns in the codebase
- Add tests for new features
- Update documentation as you code
- Ask clarifying questions when requirements are unclear
- Consider cost implications of AI features
- Think about non-technical user experience

**DON'T:**
- Generate user-facing creative content (novels, chapters, etc.)
- Make architectural changes without discussion
- Add dependencies without justification
- Compromise on security or privacy
- Assume all users are developers
- Over-engineer solutions

### Code Style
- **Go:** Use `gofumpt` for formatting, `golangci-lint` for linting
- **JavaScript/Svelte:** Use `prettier` for formatting, `eslint` for linting
- **Markdown:** Keep clean and minimal; use standard formatting
- **YAML:** Validate all workflow files before committing

### Git Workflow
- **Branch Naming:** `feature/description`, `fix/description`, `docs/description`
- **Commits:** Conventional Commits format (`feat:`, `fix:`, `docs:`, etc.)
- **PRs:** Include description, screenshots (if UI), test coverage
- **Reviews:** At least one approval before merge to main

---

## 🔮 Future Ideas (Parking Lot)

These are ideas for later, not current priorities:

- **"Ask Editor" Chat:** Real-time AI assistant in the editor
- **AI-Aided Outlining:** Structured brainstorming tool
- **Blockchain Verification:** "Authentic human writing" badges
- **Voice Input:** Dictation with AI cleanup
- **Reading Analytics:** Track how beta readers engage with drafts
- **Publishing Marketplace:** Connect authors with editors, cover designers
- **Learning Mode:** Tutorial system for new writers
- **Genre Templates:** Specialized scaffolds (mystery, romance, sci-fi)
- **Auto-Translation:** Multi-language export pipeline

---

## 🤝 Contributing

This project is currently in private development. Once we reach MVP, we'll open-source the agent system and accept community contributions for:
- New agent types
- Export templates
- Language support
- Bug fixes and improvements

---

## 📞 Contact

**Maintainer:** Christopher Mancini
**AI Collaborator:** Claude (Anthropic)
**Last Updated:** October 29, 2025

---

## 🧠 Claude's Role in This Repository

When you (Claude) are working in this repository, you should:

1. **Focus on Infrastructure & Tooling**
   - Design and implement the platform, not user content
   - Write specs, APIs, and documentation
   - Build and refine AI agent systems
   - Optimize costs and performance

2. **Maintain Code Quality**
   - Follow the style guidelines above
   - Write tests for new features
   - Keep documentation in sync with code
   - Consider edge cases and error handling

3. **Think About Users**
   - Remember many users won't be developers
   - Make git operations invisible when possible
   - Optimize for writing experience, not technical complexity
   - Balance power features with simplicity

4. **Be Cost-Conscious**
   - Every AI call costs money
   - Optimize prompts and context windows
   - Cache aggressively where appropriate
   - Provide cost estimates for new features

5. **Ask Questions**
   - When requirements are ambiguous, ask for clarification
   - When multiple approaches exist, present trade-offs
   - When scope is unclear, help define it
   - Never assume—verify understanding

6. **Respect Boundaries**
   - Never generate user-facing creative content
   - Don't make major architectural decisions unilaterally
   - Flag security or privacy concerns immediately
   - Escalate when unsure

---

**Remember:** DraftForge exists to amplify authors' voices, not replace them. Every feature should serve that mission.

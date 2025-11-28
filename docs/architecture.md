# Architecture Documentation

## System Architecture

DraftForge follows a modern, cloud-native architecture designed for scalability, maintainability, and cost-effectiveness.

---

## High-Level Architecture

```
┌─────────────────────────────────────────────────────┐
│                      User Layer                      │
│  ┌────────────┐           ┌────────────┐           │
│  │  Browser   │           │  GitHub    │           │
│  │  (Svelte)  │           │  (OAuth)   │           │
│  └────────────┘           └────────────┘           │
└─────────────┬───────────────────┬───────────────────┘
              │                   │
              │ HTTPS             │ OAuth
              ↓                   ↓
┌─────────────────────────────────────────────────────┐
│              Application Layer (Go)                  │
│  ┌──────────────────────────────────────────────┐  │
│  │         Fiber Web Framework                   │  │
│  │  ┌────────┐  ┌────────┐  ┌────────────┐     │  │
│  │  │  Auth  │  │Projects│  │  AI Agents │     │  │
│  │  │Handler │  │Handler │  │  Handler   │     │  │
│  │  └────────┘  └────────┘  └────────────┘     │  │
│  └──────────────────────────────────────────────┘  │
│                        ↓                            │
│  ┌──────────────────────────────────────────────┐  │
│  │           Business Logic Layer                │  │
│  │  ┌────────┐  ┌─────────┐  ┌──────────┐      │  │
│  │  │  Auth  │  │Projects │  │ AI Router│      │  │
│  │  │Service │  │Service  │  │ Service  │      │  │
│  │  └────────┘  └─────────┘  └──────────┘      │  │
│  └──────────────────────────────────────────────┘  │
└────────────┬────────────┬────────────┬─────────────┘
             │            │            │
             ↓            ↓            ↓
┌─────────────────────────────────────────────────────┐
│                  Data Layer                          │
│  ┌────────────┐  ┌─────────────┐  ┌─────────────┐  │
│  │PostgreSQL  │  │  GitHub API │  │ OpenRouter  │  │
│  │  Database  │  │             │  │  (AI APIs)  │  │
│  └────────────┘  └─────────────┘  └─────────────┘  │
└─────────────────────────────────────────────────────┘
```

---

## Component Details

### Frontend (SvelteKit)

**Technology:** SvelteKit 2.5 + TypeScript + Tailwind CSS + DaisyUI

**Responsibilities:**

- User interface and interaction
- Client-side routing
- Form validation
- Monaco editor integration
- Real-time preview

**Key Features:**

- Server-side rendering (SSR)
- Progressive enhancement
- Optimistic UI updates
- WebSocket support (future)

**Hosting:** Cloudflare Pages

- Global CDN distribution
- Zero-config deployment
- Automatic HTTPS

---

### Backend API (Go + Fiber)

**Technology:** Go 1.24 + Fiber v2

**Architecture Pattern:** Clean Architecture / Hexagonal Architecture

```
backend/
├── cmd/              # Application entry points
├── internal/         # Private application logic
│   ├── auth/        # Authentication domain
│   ├── projects/    # Project management domain
│   ├── ai/          # AI orchestration domain
│   └── github/      # GitHub integration domain
└── pkg/             # Public, reusable packages
```

**Domain Structure:**

```
internal/projects/
├── models.go        # Data models
├── repository.go    # Database operations
├── service.go       # Business logic
├── handlers.go      # HTTP handlers
└── routes.go        # Route definitions
```

**Hosting:** DigitalOcean App Platform

- Auto-scaling based on load
- Built-in monitoring
- Zero-downtime deployments

---

### Database (PostgreSQL)

**Version:** PostgreSQL 17

**Schema Design Principles:**

- Normalized data structure
- JSONB for flexible metadata
- Foreign keys with CASCADE
- Indexes on frequently queried columns
- Timestamp tracking with triggers

**Migration Strategy:**

- Version-controlled SQL migrations
- Forward and backward migrations
- Safe, idempotent operations

**Hosting:** DigitalOcean Managed PostgreSQL

- Automated backups
- Point-in-time recovery
- Connection pooling
- High availability

---

## Authentication Flow

```
User               Frontend              Backend              GitHub
 │                    │                     │                    │
 │  1. Click Login    │                     │                    │
 ├───────────────────>│                     │                    │
 │                    │  2. GET /auth/github│                    │
 │                    ├────────────────────>│                    │
 │                    │                     │  3. Redirect OAuth │
 │                    │                     ├───────────────────>│
 │                    │                     │                    │
 │  4. Authorize                            │                    │
 ├─────────────────────────────────────────────────────────────>│
 │                    │                     │  5. Callback+code  │
 │                    │                     │<───────────────────┤
 │                    │                     │  6. Exchange token │
 │                    │                     ├───────────────────>│
 │                    │                     │  7. Access token   │
 │                    │                     │<───────────────────┤
 │                    │                     │  8. Get user info  │
 │                    │                     ├───────────────────>│
 │                    │  9. JWT + Refresh   │                    │
 │                    │<────────────────────┤                    │
 │  10. Store tokens  │                     │                    │
 │<───────────────────┤                     │                    │
```

**Token Management:**

- **Access Token (JWT):** 15-minute expiry, contains user ID and permissions
- **Refresh Token:** 30-day expiry, stored in database, can be revoked
- **GitHub Token:** Stored encrypted in database for API calls

---

## Project Creation Flow

```
User             Frontend          Backend           GitHub API
 │                  │                 │                  │
 │  1. Create       │                 │                  │
 │  Project Form    │                 │                  │
 ├────────────────>│                 │                  │
 │                  │ 2. POST         │                  │
 │                  │ /projects       │                  │
 │                  ├────────────────>│                  │
 │                  │                 │  3. Create Repo  │
 │                  │                 ├─────────────────>│
 │                  │                 │  4. Repo created │
 │                  │                 │<─────────────────┤
 │                  │                 │  5. Create files │
 │                  │                 │  (scaffold)      │
 │                  │                 ├─────────────────>│
 │                  │                 │  6. Setup Actions│
 │                  │                 ├─────────────────>│
 │                  │                 │  7. Initial      │
 │                  │                 │  commit          │
 │                  │                 ├─────────────────>│
 │                  │  8. Project     │                  │
 │                  │  created        │                  │
 │                  │<────────────────┤                  │
 │  9. Redirect to  │                 │                  │
 │  project page    │                 │                  │
 │<─────────────────┤                 │                  │
```

**Scaffolding Process:**

1. Select template based on `project_type`
2. Populate template variables (name, author, etc.)
3. Generate file tree
4. Create GitHub repository
5. Push initial commit with all files
6. Configure webhooks and Actions
7. Store project metadata in database

---

## AI Agent System

### Agent Execution Flow

```
Trigger            Queue              Worker             AI API
 │                  │                   │                  │
 │  1. PR created   │                   │                  │
 │  (GitHub Action) │                   │                  │
 ├─────────────────>│                   │                  │
 │                  │  2. Enqueue       │                  │
 │                  │  agent run        │                  │
 │                  ├──────────────────>│                  │
 │                  │                   │  3. Fetch context│
 │                  │                   │  (changed files) │
 │                  │                   │                  │
 │                  │                   │  4. Build prompt │
 │                  │                   │                  │
 │                  │                   │  5. Call AI      │
 │                  │                   ├─────────────────>│
 │                  │                   │  6. Response     │
 │                  │                   │<─────────────────┤
 │                  │                   │  7. Parse results│
 │                  │                   │                  │
 │                  │                   │  8. Post comment │
 │                  │                   │  to PR           │
 │                  │  9. Update status │                  │
 │                  │<──────────────────┤                  │
 │  10. Notify user │                   │                  │
 │<─────────────────┤                   │                  │
```

### Agent Architecture

**Agent Types:**

1. **ContinuityBot** - Character/plot consistency
2. **StyleBot** - Writing style analysis
3. **TimelineBot** - Chronological validation
4. **FactBot** - Factual accuracy checking

**Agent Structure:**

```go
type Agent interface {
    Name() string
    Trigger() TriggerType
    Model() string
    Execute(ctx Context) (Result, error)
}

type Context struct {
    ProjectID    int
    Files        []File
    PreviousRuns []Run
    Settings     AgentSettings
}

type Result struct {
    Issues      []Issue
    Suggestions []Suggestion
    TokensUsed  int
}
```

**Queue Implementation:**

- PostgreSQL-based job queue (simple, reliable)
- Retry logic with exponential backoff
- Dead letter queue for failed jobs
- Priority levels based on trigger type

---

## GitHub Integration

### Repository Structure

DraftForge creates this structure in user repos:

```
user-project/
├── chapters/                   # Manuscript chapters
│   ├── 01-chapter-one.md
│   └── 02-chapter-two.md
├── manuscript/
│   ├── metadata.yml           # Book metadata
│   ├── outline.md
│   └── character-bible.md
├── .draftforge/
│   ├── config.yml             # Project settings
│   └── agents.yml             # Agent configuration
├── .github/
│   └── workflows/
│       ├── stats.yml          # Word count tracking
│       ├── compile.yml        # EPUB/PDF build
│       └── ai-review.yml      # AI agent triggers
└── README.md
```

### GitHub Actions Integration

**Stats Workflow** (`stats.yml`):

- Triggers on every push
- Counts words per chapter
- Updates project stats in DraftForge
- Generates progress chart

**Compile Workflow** (`compile.yml`):

- Triggers on tag creation
- Uses Pandoc to build EPUB/PDF
- Uploads artifacts to Cloudflare R2
- Creates GitHub release

**AI Review Workflow** (`ai-review.yml`):

- Triggers on PR creation
- Calls DraftForge API to queue agent runs
- Posts results as PR comments
- Blocks merge if critical issues found (optional)

### Automation Catalog UX

- In the DraftForge app, authors toggle workflows (stats, compile, AI review, audio/SSML, etc.) from a catalog UI.
- Enabling/disabling writes the corresponding `.github/workflows/*.yml` (and `.draftforge/config.yml` if needed) via commit/PR so automation remains transparent and user-owned.
- Workflow defaults live in the platform, but the repo is always the source of truth.

---

## Data Flow Patterns

### Read Path (Get Project)

```
User Request → Fiber Handler → Service Layer → Repository → PostgreSQL
                    ↓              ↓              ↓
                  JWT Auth     Business Logic   Query
```

### Write Path (Update Project)

```
User Request → Fiber Handler → Service Layer → Repository → PostgreSQL
                    ↓              ↓              ↓              ↓
                  Validate     Transform      Transaction    Commit
                                   ↓
                              Side Effects
                             (GitHub API)
```

### Agent Execution Path

```
GitHub Webhook → Queue → Worker → OpenRouter → Result Parser → Database
                                      ↓                            ↓
                                  AI Model                    Comment on PR
```

---

## Security Architecture

### Authentication

- GitHub OAuth 2.0 for user login
- JWT for API authentication
- Refresh tokens for session management
- PKCE for OAuth flow (future enhancement)

### Authorization

- Row-level security via user_id checks
- Project ownership validation
- GitHub token scoping (minimal permissions)

### Data Security

- Passwords never stored (OAuth only)
- GitHub tokens encrypted at rest (AES-256)
- HTTPS everywhere
- CORS configured for frontend origin only

### API Security

- Rate limiting per user tier
- Request size limits
- SQL injection prevention (parameterized queries)
- XSS prevention (Content Security Policy)

---

## Scalability Considerations

### Horizontal Scaling

- Stateless API servers (can run multiple instances)
- Load balancing via DigitalOcean/Cloudflare
- Database connection pooling

### Vertical Scaling

- Database can scale up as needed
- Separate read replicas (future)
- Caching layer (Redis) for hot data (future)

### Cost Optimization

- Efficient AI token usage (context pruning)
- Cloudflare caching for static assets
- Lazy loading of large files
- Archive old projects to cold storage

---

## Monitoring & Observability

### Metrics (Future)

- API response times (p50, p95, p99)
- Database query performance
- AI token consumption
- Error rates by endpoint
- Active users and sessions

### Logging

- Structured JSON logs
- Request/response logging
- Error stack traces
- Audit trail for sensitive operations

### Alerting

- Database connection failures
- High error rates
- AI API failures
- Low credit balance warnings

---

## Disaster Recovery

### Backup Strategy

- PostgreSQL: Daily automated backups (7-day retention)
- Point-in-time recovery available
- GitHub repos: User owns, naturally backed up

### Failure Scenarios

**Database Failure:**

- DigitalOcean auto-failover (managed service)
- Restore from backup (RPO: 24 hours)

**API Server Failure:**

- Auto-restart via App Platform
- Multiple instances for redundancy

**GitHub API Outage:**

- Queue operations for retry
- Graceful degradation (read-only mode)

**AI API Failure:**

- Fallback to alternative models
- Retry with exponential backoff
- User notification of delays

---

## Future Enhancements

### Phase 2 Additions

- Redis cache layer
- WebSocket for real-time updates
- Elasticsearch for full-text search
- Background job dashboard

### Phase 3 Additions

- Multi-region deployment
- Read replicas
- GraphQL API
- Real-time collaboration (CRDT)

---

## Technology Choices Rationale

**Go + Fiber:**

- High performance, low memory footprint
- Excellent concurrency primitives
- Type safety
- Fast compilation
- Simple deployment (single binary)

**SvelteKit:**

- Minimal bundle size
- Great DX with reactive primitives
- SSR out of the box
- Cloudflare Pages deployment

**PostgreSQL:**

- Mature, battle-tested
- Excellent JSONB support
- Strong consistency guarantees
- Rich extension ecosystem

**DigitalOcean:**

- Competitive pricing
- Simple, predictable billing
- Good developer experience
- Sufficient scale for MVP

**OpenRouter:**

- Multi-model support
- Transparent pricing
- Single API for all providers
- No vendor lock-in

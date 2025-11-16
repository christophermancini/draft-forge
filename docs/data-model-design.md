# Data Model Design

This document provides a comprehensive design of all data models for the DraftForge platform, including database schema, relationships, and business logic considerations.

---

## 🎯 Design Principles

1. **Normalized but pragmatic** - Balance normalization with query performance
2. **Audit trails** - Track creation and updates for all major entities
3. **Soft deletes where appropriate** - Allow recovery of accidentally deleted data
4. **JSONB for flexibility** - Use structured JSON for settings and metadata
5. **Indexes on query patterns** - Index foreign keys and frequently filtered columns
6. **Timestamps with timezone** - Always use `TIMESTAMP WITH TIME ZONE`

---

## 📊 Entity Relationship Diagram

```
┌─────────────┐
│    Users    │
└──────┬──────┘
       │ 1
       │
       │ N
┌──────▼───────────┐        ┌─────────────────┐
│    Projects      │───────▶│  Project Stats  │
└──────┬───────────┘   1:1  └─────────────────┘
       │ 1
       │
       ├─────────────────────┐
       │ N                   │ N
┌──────▼──────┐      ┌───────▼────────┐
│ Agent Runs  │      │  Collaborators │
└─────────────┘      └────────────────┘

┌──────────────┐
│  AI Credits  │◀───── 1:1 ───── Users
└──────────────┘

┌───────────────┐
│ AI Usage Log  │◀───── N:1 ───── Users
└───────┬───────┘
        │ N:1
        ▼
    Projects

┌─────────────────┐
│  Export Jobs    │◀───── N:1 ───── Projects
└─────────────────┘

┌─────────────────┐
│  Subscriptions  │◀───── 1:1 ───── Users
└─────────────────┘

┌─────────────────┐
│  Notifications  │◀───── N:1 ───── Users
└─────────────────┘

┌─────────────────┐
│  Webhooks       │◀───── N:1 ───── Projects
└─────────────────┘
```

---

## 📋 Core Entities

### 1. Users

**Purpose:** Store user account information from GitHub OAuth

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,

    -- GitHub OAuth data
    github_id INTEGER UNIQUE NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255),
    avatar_url TEXT,
    github_access_token TEXT, -- Encrypted at application layer
    github_refresh_token TEXT, -- Encrypted at application layer
    github_token_expires_at TIMESTAMP WITH TIME ZONE,

    -- Profile
    display_name VARCHAR(255),
    bio TEXT,
    website_url TEXT,

    -- Settings
    preferences JSONB DEFAULT '{}', -- UI preferences, notifications, etc.

    -- Status
    is_active BOOLEAN DEFAULT TRUE,
    email_verified BOOLEAN DEFAULT FALSE,
    last_login_at TIMESTAMP WITH TIME ZONE,

    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft delete
);

CREATE INDEX idx_users_github_id ON users(github_id);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_is_active ON users(is_active) WHERE deleted_at IS NULL;
```

**Preferences JSONB structure:**
```json
{
  "theme": "dark",
  "editor": {
    "fontSize": 14,
    "lineHeight": 1.6,
    "vim_mode": false
  },
  "notifications": {
    "email_on_agent_complete": true,
    "email_on_low_credits": true
  },
  "default_agent_settings": {
    "continuity_enabled": true,
    "style_enabled": true
  }
}
```

---

### 2. Projects

**Purpose:** User's writing projects (novels, screenplays, etc.)

```sql
CREATE TABLE projects (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Basic info
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL, -- URL-safe identifier
    description TEXT,
    project_type VARCHAR(50) NOT NULL, -- 'novel', 'screenplay', 'technical-book', 'non-fiction', 'poetry', 'short-story'

    -- GitHub integration
    github_repo_id BIGINT,
    github_repo_name VARCHAR(255), -- owner/repo
    github_repo_url TEXT,
    github_default_branch VARCHAR(100) DEFAULT 'main',
    github_webhook_id BIGINT,
    github_webhook_secret VARCHAR(255), -- Encrypted

    -- Project structure
    chapters_directory VARCHAR(255) DEFAULT 'chapters',
    manuscript_directory VARCHAR(255) DEFAULT 'manuscript',

    -- Settings
    settings JSONB DEFAULT '{}', -- Agent configs, compile settings, etc.

    -- Metadata
    genre VARCHAR(100),
    target_word_count INTEGER,
    language VARCHAR(10) DEFAULT 'en',

    -- Status
    status VARCHAR(50) DEFAULT 'active', -- 'active', 'archived', 'completed'
    visibility VARCHAR(20) DEFAULT 'private', -- 'private', 'unlisted', 'public'

    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    archived_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE, -- Soft delete

    UNIQUE(user_id, slug),
    CONSTRAINT valid_project_type CHECK (project_type IN ('novel', 'screenplay', 'technical-book', 'non-fiction', 'poetry', 'short-story')),
    CONSTRAINT valid_status CHECK (status IN ('active', 'archived', 'completed'))
);

CREATE INDEX idx_projects_user_id ON projects(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_projects_github_repo_id ON projects(github_repo_id);
CREATE INDEX idx_projects_status ON projects(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_projects_created_at ON projects(created_at DESC);
```

**Settings JSONB structure:**
```json
{
  "agents": {
    "continuity": {
      "enabled": true,
      "model": "claude-3-5-sonnet",
      "trigger": "pr",
      "context_window": 5,
      "auto_comment": true
    },
    "style": {
      "enabled": true,
      "model": "gpt-4o-mini",
      "trigger": "commit"
    },
    "timeline": {
      "enabled": false
    }
  },
  "compilation": {
    "pandoc_template": "default",
    "include_cover": true,
    "auto_build_on_tag": true
  },
  "github_actions": {
    "stats_enabled": true,
    "compile_enabled": true,
    "ai_review_enabled": true
  }
}
```

---

### 3. Project Stats

**Purpose:** Track word counts, chapters, and writing progress

```sql
CREATE TABLE project_stats (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,

    -- Current stats
    word_count INTEGER DEFAULT 0,
    chapter_count INTEGER DEFAULT 0,
    character_count INTEGER DEFAULT 0,

    -- Git tracking
    last_commit_sha VARCHAR(255),
    last_commit_at TIMESTAMP WITH TIME ZONE,
    total_commits INTEGER DEFAULT 0,

    -- Writing streaks
    current_streak_days INTEGER DEFAULT 0,
    longest_streak_days INTEGER DEFAULT 0,
    last_write_date DATE,

    -- Historical data (stored as time series)
    stats_history JSONB DEFAULT '[]', -- Array of daily snapshots

    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(project_id)
);

CREATE INDEX idx_project_stats_project_id ON project_stats(project_id);
CREATE INDEX idx_project_stats_updated_at ON project_stats(updated_at DESC);
```

**Stats History JSONB structure:**
```json
[
  {
    "date": "2025-10-29",
    "word_count": 45678,
    "chapter_count": 12,
    "words_written_today": 1250,
    "commits": 3
  },
  {
    "date": "2025-10-28",
    "word_count": 44428,
    "chapter_count": 12,
    "words_written_today": 890,
    "commits": 2
  }
]
```

---

### 4. AI Credits

**Purpose:** Track AI token usage and subscription limits

```sql
CREATE TABLE ai_credits (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Credit balance
    credits_remaining INTEGER NOT NULL DEFAULT 0, -- In tokens
    credits_total INTEGER NOT NULL DEFAULT 0,
    credits_bonus INTEGER DEFAULT 0, -- From promotions, referrals

    -- Subscription
    subscription_tier VARCHAR(50) DEFAULT 'free', -- 'free', 'creator', 'professional', 'enterprise'
    subscription_status VARCHAR(50) DEFAULT 'active', -- 'active', 'canceled', 'expired', 'paused'
    subscription_period VARCHAR(20) DEFAULT 'monthly', -- 'monthly', 'annual'

    -- Billing cycle
    current_period_start TIMESTAMP WITH TIME ZONE,
    current_period_end TIMESTAMP WITH TIME ZONE,
    next_billing_date TIMESTAMP WITH TIME ZONE,

    -- Usage this period
    tokens_used_this_period INTEGER DEFAULT 0,
    tokens_limit_per_period INTEGER,

    -- Stripe integration
    stripe_customer_id VARCHAR(255),
    stripe_subscription_id VARCHAR(255),

    -- BYOK (Bring Your Own Key)
    byok_enabled BOOLEAN DEFAULT FALSE,
    byok_openrouter_key TEXT, -- Encrypted

    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(user_id),
    CONSTRAINT valid_subscription_tier CHECK (subscription_tier IN ('free', 'creator', 'professional', 'enterprise')),
    CONSTRAINT valid_subscription_status CHECK (subscription_status IN ('active', 'canceled', 'expired', 'paused'))
);

CREATE INDEX idx_ai_credits_user_id ON ai_credits(user_id);
CREATE INDEX idx_ai_credits_subscription_tier ON ai_credits(subscription_tier);
CREATE INDEX idx_ai_credits_next_billing_date ON ai_credits(next_billing_date);
```

**Tier Limits:**
```go
const (
    FreeTier         = 10000    // tokens/month
    CreatorTier      = 200000   // tokens/month
    ProfessionalTier = 750000   // tokens/month
    EnterpriseTier   = -1       // unlimited (or custom)
)
```

---

### 5. AI Usage Log

**Purpose:** Track every AI API call for billing and analytics

```sql
CREATE TABLE ai_usage_log (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    project_id INTEGER REFERENCES projects(id) ON DELETE SET NULL,
    agent_run_id INTEGER REFERENCES agent_runs(id) ON DELETE SET NULL,

    -- API call details
    agent_type VARCHAR(50) NOT NULL, -- 'continuity', 'style', 'timeline', 'fact'
    model_name VARCHAR(100) NOT NULL, -- 'claude-3-5-sonnet', 'gpt-4o-mini', etc.
    provider VARCHAR(50) NOT NULL, -- 'openrouter', 'direct-openai', 'byok'

    -- Token usage
    tokens_prompt INTEGER NOT NULL,
    tokens_completion INTEGER NOT NULL,
    tokens_total INTEGER NOT NULL,

    -- Cost
    cost_cents INTEGER NOT NULL, -- Total cost in cents
    cost_per_token_cents DECIMAL(10, 8), -- For analytics

    -- Context
    context_type VARCHAR(50), -- 'pr_review', 'commit_hook', 'manual', 'scheduled'
    files_analyzed INTEGER,

    -- Result
    status VARCHAR(20) NOT NULL, -- 'success', 'failed', 'timeout', 'rate_limited'
    error_code VARCHAR(50),
    error_message TEXT,
    response_time_ms INTEGER,

    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT valid_status CHECK (status IN ('success', 'failed', 'timeout', 'rate_limited'))
);

CREATE INDEX idx_ai_usage_user_id ON ai_usage_log(user_id);
CREATE INDEX idx_ai_usage_project_id ON ai_usage_log(project_id);
CREATE INDEX idx_ai_usage_created_at ON ai_usage_log(created_at DESC);
CREATE INDEX idx_ai_usage_agent_type ON ai_usage_log(agent_type);
CREATE INDEX idx_ai_usage_status ON ai_usage_log(status);

-- Partitioning strategy (future optimization)
-- Partition by month for better query performance
```

---

### 6. Agent Runs

**Purpose:** Track AI agent execution history and results

```sql
CREATE TABLE agent_runs (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Agent details
    agent_type VARCHAR(50) NOT NULL, -- 'continuity', 'style', 'timeline', 'fact'
    trigger VARCHAR(50) NOT NULL, -- 'manual', 'commit', 'pr', 'scheduled', 'webhook'
    trigger_ref VARCHAR(255), -- commit SHA, PR number, etc.

    -- Context
    context JSONB DEFAULT '{}', -- Files analyzed, settings used, etc.

    -- Execution
    status VARCHAR(20) NOT NULL DEFAULT 'queued', -- 'queued', 'running', 'completed', 'failed', 'canceled'
    priority INTEGER DEFAULT 5, -- 1-10, higher = more urgent

    -- Results
    results JSONB, -- Issues found, suggestions, etc.
    summary TEXT, -- Human-readable summary
    issues_found INTEGER DEFAULT 0,
    severity_breakdown JSONB, -- Count by severity level

    -- Error handling
    error_message TEXT,
    retry_count INTEGER DEFAULT 0,
    max_retries INTEGER DEFAULT 3,

    -- GitHub integration
    github_comment_id BIGINT, -- Comment posted to PR
    github_check_run_id BIGINT, -- GitHub Checks API

    -- Performance
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    duration_seconds INTEGER, -- Calculated

    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT valid_status CHECK (status IN ('queued', 'running', 'completed', 'failed', 'canceled')),
    CONSTRAINT valid_priority CHECK (priority BETWEEN 1 AND 10)
);

CREATE INDEX idx_agent_runs_project_id ON agent_runs(project_id);
CREATE INDEX idx_agent_runs_user_id ON agent_runs(user_id);
CREATE INDEX idx_agent_runs_status ON agent_runs(status);
CREATE INDEX idx_agent_runs_created_at ON agent_runs(created_at DESC);
CREATE INDEX idx_agent_runs_agent_type ON agent_runs(agent_type);
CREATE INDEX idx_agent_runs_trigger ON agent_runs(trigger);
```

**Results JSONB structure:**
```json
{
  "issues": [
    {
      "id": "cont-001",
      "severity": "error",
      "category": "character_consistency",
      "file": "chapters/05-chapter-five.md",
      "line": 45,
      "message": "Character 'Sarah' has blue eyes here, but brown eyes in Chapter 3.",
      "suggestion": "Update to match earlier description or explain the change.",
      "confidence": 0.95
    },
    {
      "id": "cont-002",
      "severity": "warning",
      "category": "world_rules",
      "file": "chapters/12-chapter-twelve.md",
      "line": 102,
      "message": "Magic system requires verbal components, but spell cast silently.",
      "suggestion": "Add explanation or revise scene.",
      "confidence": 0.88
    }
  ],
  "stats": {
    "files_analyzed": 5,
    "total_words": 12500,
    "issues_by_severity": {
      "error": 1,
      "warning": 3,
      "info": 5
    }
  },
  "model_info": {
    "model": "claude-3-5-sonnet",
    "tokens_used": 12500,
    "cost_cents": 25
  }
}
```

---

### 7. Collaborators

**Purpose:** Multi-user access to projects

```sql
CREATE TABLE collaborators (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    invited_by_user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,

    -- Permissions
    role VARCHAR(50) NOT NULL DEFAULT 'viewer', -- 'owner', 'editor', 'reviewer', 'viewer'
    permissions JSONB DEFAULT '{}', -- Granular permissions

    -- Invitation
    invitation_status VARCHAR(50) DEFAULT 'pending', -- 'pending', 'accepted', 'declined', 'revoked'
    invitation_token VARCHAR(255) UNIQUE,
    invitation_expires_at TIMESTAMP WITH TIME ZONE,

    -- Activity
    last_accessed_at TIMESTAMP WITH TIME ZONE,

    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    removed_at TIMESTAMP WITH TIME ZONE,

    UNIQUE(project_id, user_id),
    CONSTRAINT valid_role CHECK (role IN ('owner', 'editor', 'reviewer', 'viewer')),
    CONSTRAINT valid_invitation_status CHECK (invitation_status IN ('pending', 'accepted', 'declined', 'revoked'))
);

CREATE INDEX idx_collaborators_project_id ON collaborators(project_id);
CREATE INDEX idx_collaborators_user_id ON collaborators(user_id);
CREATE INDEX idx_collaborators_role ON collaborators(role);
CREATE INDEX idx_collaborators_invitation_status ON collaborators(invitation_status);
```

**Permissions JSONB structure:**
```json
{
  "can_read": true,
  "can_write": false,
  "can_delete": false,
  "can_invite": false,
  "can_run_agents": true,
  "can_view_stats": true,
  "can_export": false,
  "can_manage_settings": false
}
```

---

### 8. Export Jobs

**Purpose:** Track EPUB/PDF export jobs

```sql
CREATE TABLE export_jobs (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Export details
    export_type VARCHAR(50) NOT NULL, -- 'epub', 'pdf', 'mobi', 'docx', 'html'
    export_format_version VARCHAR(20), -- '3.0' for EPUB, etc.

    -- Source
    source_ref VARCHAR(255), -- commit SHA or tag
    source_branch VARCHAR(100),

    -- Settings
    settings JSONB DEFAULT '{}', -- Template, fonts, includes, etc.

    -- Status
    status VARCHAR(50) DEFAULT 'queued', -- 'queued', 'processing', 'completed', 'failed'
    progress_percent INTEGER DEFAULT 0,

    -- Output
    output_url TEXT, -- Cloudflare R2 URL
    output_size_bytes BIGINT,
    output_checksum VARCHAR(64),
    expires_at TIMESTAMP WITH TIME ZONE, -- Auto-delete after 30 days

    -- Error handling
    error_message TEXT,
    build_log TEXT,

    -- Performance
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    duration_seconds INTEGER,

    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT valid_status CHECK (status IN ('queued', 'processing', 'completed', 'failed')),
    CONSTRAINT valid_export_type CHECK (export_type IN ('epub', 'pdf', 'mobi', 'docx', 'html'))
);

CREATE INDEX idx_export_jobs_project_id ON export_jobs(project_id);
CREATE INDEX idx_export_jobs_user_id ON export_jobs(user_id);
CREATE INDEX idx_export_jobs_status ON export_jobs(status);
CREATE INDEX idx_export_jobs_created_at ON export_jobs(created_at DESC);
CREATE INDEX idx_export_jobs_expires_at ON export_jobs(expires_at) WHERE status = 'completed';
```

---

### 9. Notifications

**Purpose:** In-app and email notifications

```sql
CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Notification details
    type VARCHAR(50) NOT NULL, -- 'agent_complete', 'low_credits', 'export_ready', 'invitation', 'comment'
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,

    -- Related entities
    related_entity_type VARCHAR(50), -- 'project', 'agent_run', 'export_job'
    related_entity_id INTEGER,

    -- Delivery
    is_read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMP WITH TIME ZONE,

    -- Actions
    action_url TEXT,
    action_label VARCHAR(100),

    -- Priority
    priority VARCHAR(20) DEFAULT 'normal', -- 'low', 'normal', 'high', 'urgent'

    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE, -- Auto-delete old notifications

    CONSTRAINT valid_priority CHECK (priority IN ('low', 'normal', 'high', 'urgent'))
);

CREATE INDEX idx_notifications_user_id ON notifications(user_id) WHERE is_read = FALSE;
CREATE INDEX idx_notifications_created_at ON notifications(created_at DESC);
CREATE INDEX idx_notifications_type ON notifications(type);
CREATE INDEX idx_notifications_expires_at ON notifications(expires_at) WHERE expires_at IS NOT NULL;
```

---

### 10. Webhooks

**Purpose:** Outbound webhooks for integrations

```sql
CREATE TABLE webhooks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,

    -- Webhook config
    url TEXT NOT NULL,
    secret VARCHAR(255), -- For HMAC signature

    -- Events to subscribe to
    events TEXT[] NOT NULL, -- ['agent.run.completed', 'export.completed', 'stats.updated']

    -- Status
    is_active BOOLEAN DEFAULT TRUE,

    -- Delivery stats
    last_delivery_at TIMESTAMP WITH TIME ZONE,
    last_status_code INTEGER,
    total_deliveries INTEGER DEFAULT 0,
    failed_deliveries INTEGER DEFAULT 0,
    consecutive_failures INTEGER DEFAULT 0,

    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    disabled_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_webhooks_user_id ON webhooks(user_id);
CREATE INDEX idx_webhooks_project_id ON webhooks(project_id);
CREATE INDEX idx_webhooks_is_active ON webhooks(is_active);
```

---

### 11. Webhook Deliveries

**Purpose:** Track webhook delivery attempts

```sql
CREATE TABLE webhook_deliveries (
    id SERIAL PRIMARY KEY,
    webhook_id INTEGER NOT NULL REFERENCES webhooks(id) ON DELETE CASCADE,

    -- Event
    event_type VARCHAR(50) NOT NULL,
    event_payload JSONB NOT NULL,

    -- Delivery
    status VARCHAR(20) NOT NULL, -- 'pending', 'delivered', 'failed'
    status_code INTEGER,
    response_body TEXT,
    response_time_ms INTEGER,

    -- Retry
    attempt_number INTEGER DEFAULT 1,
    next_retry_at TIMESTAMP WITH TIME ZONE,

    -- Error
    error_message TEXT,

    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    delivered_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT valid_status CHECK (status IN ('pending', 'delivered', 'failed'))
);

CREATE INDEX idx_webhook_deliveries_webhook_id ON webhook_deliveries(webhook_id);
CREATE INDEX idx_webhook_deliveries_status ON webhook_deliveries(status);
CREATE INDEX idx_webhook_deliveries_created_at ON webhook_deliveries(created_at DESC);
CREATE INDEX idx_webhook_deliveries_next_retry_at ON webhook_deliveries(next_retry_at) WHERE status = 'pending';
```

---

## 🔧 Utility Tables

### 12. Audit Log

**Purpose:** Track important user actions for security and compliance

```sql
CREATE TABLE audit_log (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,

    -- Action
    action VARCHAR(100) NOT NULL, -- 'login', 'project.create', 'agent.run', 'settings.update'
    entity_type VARCHAR(50),
    entity_id INTEGER,

    -- Context
    ip_address INET,
    user_agent TEXT,

    -- Changes (for update operations)
    changes JSONB,

    -- Metadata
    metadata JSONB,

    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_audit_log_user_id ON audit_log(user_id);
CREATE INDEX idx_audit_log_action ON audit_log(action);
CREATE INDEX idx_audit_log_created_at ON audit_log(created_at DESC);
CREATE INDEX idx_audit_log_entity ON audit_log(entity_type, entity_id);

-- Partition by month for performance
```

---

### 13. Feature Flags

**Purpose:** Control feature rollouts and A/B testing

```sql
CREATE TABLE feature_flags (
    id SERIAL PRIMARY KEY,

    -- Flag
    flag_key VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,

    -- Status
    is_enabled BOOLEAN DEFAULT FALSE,

    -- Rollout strategy
    rollout_percentage INTEGER DEFAULT 0, -- 0-100
    rollout_user_ids INTEGER[], -- Specific users
    rollout_tiers TEXT[], -- Subscription tiers

    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT valid_rollout_percentage CHECK (rollout_percentage BETWEEN 0 AND 100)
);

CREATE INDEX idx_feature_flags_is_enabled ON feature_flags(is_enabled);
```

---

## 🔐 Security Considerations

### Encryption at Rest
- `github_access_token`, `github_refresh_token` - Encrypted with AES-256
- `github_webhook_secret` - Encrypted
- `byok_openrouter_key` - Encrypted
- Encryption keys managed via environment variables (future: AWS KMS, HashiCorp Vault)

### Sensitive Data Access
- Use application-level encryption/decryption
- Never log tokens or secrets
- Rotate encryption keys periodically

### Data Retention
- Soft delete for users and projects (30-day recovery window)
- Auto-expire export jobs after 30 days
- Auto-delete old notifications after 90 days
- Partition and archive AI usage logs older than 1 year

---

## 📈 Performance Optimizations

### Indexes
- All foreign keys indexed
- Frequently filtered columns indexed (status, created_at, etc.)
- Composite indexes where appropriate

### Partitioning Strategy
- `ai_usage_log` - Partition by month
- `audit_log` - Partition by month
- `webhook_deliveries` - Partition by month

### Archival
- Move old data to separate archive tables
- Use TimescaleDB for time-series data (future)

### Caching Strategy
- User profiles and preferences - Redis cache (5 min TTL)
- Project settings - Redis cache (1 min TTL)
- AI credit balance - Redis cache (30 sec TTL)
- Stats history - Aggregate in background job

---

## 🧪 Data Validation Rules

### Business Rules
1. User cannot have > 100 active projects (free tier)
2. Project slug must be URL-safe (lowercase, hyphen, underscore only)
3. Agent runs auto-cancel after 5 minutes
4. Export jobs auto-expire after 30 days
5. Webhook disabled after 10 consecutive failures
6. Collaborators limited to 5 per project (free tier)

### Constraints
- Email format validation (application layer)
- Username: 3-39 characters, alphanumeric + hyphen
- Project name: 1-100 characters
- Slug: 1-100 characters
- Token amounts: Non-negative

---

## 🚀 Migration Strategy

### Phase 1 (Current)
- Users, Projects, AI Credits, AI Usage Log, Agent Runs, Project Stats

### Phase 2 (Months 2-3)
- Collaborators, Export Jobs, Notifications

### Phase 3 (Months 4-6)
- Webhooks, Webhook Deliveries, Feature Flags, Audit Log

### Rollback Plan
- Each migration has corresponding `.down.sql`
- Test rollback in staging before production deployment
- Keep backups before major schema changes

---

## 📝 Notes

### Future Enhancements
- Full-text search (PostgreSQL `tsvector` or Elasticsearch)
- Character/world database tables
- Timeline event tracking
- Reading analytics (if users share drafts)
- Version comparison metadata
- Social features (following, likes, comments)

### Alternatives Considered
- NoSQL for flexibility → Rejected (relational data, ACID compliance critical)
- Separate analytics DB → Deferred (premature optimization)
- Graph DB for relationships → Deferred (not needed at MVP scale)

---

**Last Updated:** October 29, 2025
**Schema Version:** 1.0.0
**Status:** Draft - Ready for Implementation

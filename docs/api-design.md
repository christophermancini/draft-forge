# API Design Documentation

## Overview

The DraftForge API is a RESTful service built with Go Fiber that provides endpoints for project management, GitHub integration, and AI-assisted writing features.

**Base URL:** `http://localhost:8080/api/v1`

**Authentication:** JWT-based with GitHub OAuth

---

## Authentication Flow

### 1. GitHub OAuth

```
GET /api/v1/auth/github
```

Redirects user to GitHub OAuth consent screen.

**Response:** HTTP 302 redirect to GitHub

---

### 2. GitHub Callback

```
GET /api/v1/auth/github/callback?code={code}&state={state}
```

Handles OAuth callback from GitHub, exchanges code for access token.

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "github_id": 12345,
    "username": "authorname",
    "email": "author@example.com",
    "avatar_url": "https://avatars.githubusercontent.com/u/12345"
  }
}
```

---

### 3. Refresh Token

```
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

## User Endpoints

### Get Current User

```
GET /api/v1/users/me
Authorization: Bearer {access_token}
```

**Response:**
```json
{
  "id": 1,
  "github_id": 12345,
  "username": "authorname",
  "email": "author@example.com",
  "avatar_url": "https://avatars.githubusercontent.com/u/12345",
  "credits": {
    "remaining": 150000,
    "total": 200000,
    "subscription_tier": "creator"
  }
}
```

---

## Project Endpoints

### List Projects

```
GET /api/v1/projects
Authorization: Bearer {access_token}
```

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Results per page (default: 20)

**Response:**
```json
{
  "projects": [
    {
      "id": 1,
      "name": "My Novel",
      "slug": "my-novel",
      "description": "A thrilling adventure story",
      "project_type": "novel",
      "github_repo_url": "https://github.com/authorname/my-novel",
      "stats": {
        "word_count": 45678,
        "chapter_count": 12
      },
      "created_at": "2025-10-01T12:00:00Z",
      "updated_at": "2025-10-29T14:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 3
  }
}
```

---

### Get Project

```
GET /api/v1/projects/{id}
Authorization: Bearer {access_token}
```

**Response:**
```json
{
  "id": 1,
  "name": "My Novel",
  "slug": "my-novel",
  "description": "A thrilling adventure story",
  "project_type": "novel",
  "github_repo_id": 67890,
  "github_repo_name": "authorname/my-novel",
  "github_repo_url": "https://github.com/authorname/my-novel",
  "settings": {
    "agents": {
      "continuity": {
        "enabled": true,
        "model": "claude-3-5-sonnet",
        "trigger": "pr"
      },
      "style": {
        "enabled": true,
        "model": "gpt-4o-mini",
        "trigger": "commit"
      }
    }
  },
  "stats": {
    "word_count": 45678,
    "chapter_count": 12,
    "last_commit_sha": "abc123def456",
    "last_commit_at": "2025-10-29T10:15:00Z"
  },
  "created_at": "2025-10-01T12:00:00Z",
  "updated_at": "2025-10-29T14:30:00Z"
}
```

---

### Create Project

```
POST /api/v1/projects
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "name": "My Novel",
  "description": "A thrilling adventure story",
  "project_type": "novel",
  "github_repo_name": "my-novel"
}
```

**Project Types:** `novel`, `screenplay`, `technical-book`, `non-fiction`

**Response:** Same as Get Project (201 Created)

**Side Effects:**
1. Creates a new GitHub repository under user's account
2. Scaffolds project structure based on `project_type`
3. Initializes `.draftforge/` configuration
4. Sets up GitHub Actions workflows
5. Creates initial commit

---

### Update Project

```
PATCH /api/v1/projects/{id}
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "name": "My Updated Novel Title",
  "description": "New description",
  "settings": {
    "agents": {
      "continuity": {
        "enabled": false
      }
    }
  }
}
```

**Response:** Same as Get Project

---

### Delete Project

```
DELETE /api/v1/projects/{id}
Authorization: Bearer {access_token}
```

**Response:** 204 No Content

**Note:** This does NOT delete the GitHub repository, only removes it from DraftForge.

---

## AI Agent Endpoints

### List Agent Runs

```
GET /api/v1/projects/{project_id}/agents/runs
Authorization: Bearer {access_token}
```

**Query Parameters:**
- `agent_type` (optional): Filter by agent type
- `status` (optional): Filter by status

**Response:**
```json
{
  "runs": [
    {
      "id": 1,
      "agent_type": "continuity",
      "trigger": "pr",
      "status": "completed",
      "results": {
        "issues_found": 2,
        "suggestions": [
          {
            "severity": "error",
            "chapter": "chapter-05",
            "line": 45,
            "message": "Character 'Sarah' has blue eyes here, but brown eyes in Chapter 3."
          }
        ]
      },
      "tokens_used": 12500,
      "cost_cents": 25,
      "started_at": "2025-10-29T10:00:00Z",
      "completed_at": "2025-10-29T10:01:15Z"
    }
  ]
}
```

---

### Trigger Agent Run

```
POST /api/v1/projects/{project_id}/agents/run
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "agent_type": "continuity",
  "context": {
    "commit_sha": "abc123def456",
    "files_changed": ["chapters/05-chapter-five.md"]
  }
}
```

**Response:**
```json
{
  "run_id": 1,
  "status": "queued",
  "message": "Agent run queued successfully"
}
```

**Note:** Returns immediately. Agent runs asynchronously in background.

---

### Get Agent Run Status

```
GET /api/v1/projects/{project_id}/agents/runs/{run_id}
Authorization: Bearer {access_token}
```

**Response:** Same as individual run from List Agent Runs

---

## GitHub Integration Endpoints

### List Repositories

```
GET /api/v1/github/repos
Authorization: Bearer {access_token}
```

**Response:**
```json
{
  "repositories": [
    {
      "id": 67890,
      "name": "my-novel",
      "full_name": "authorname/my-novel",
      "url": "https://github.com/authorname/my-novel",
      "private": true,
      "is_draftforge_project": true
    }
  ]
}
```

---

### Sync Repository

```
POST /api/v1/projects/{project_id}/sync
Authorization: Bearer {access_token}
```

Syncs project stats from GitHub (word count, commit history, etc.)

**Response:**
```json
{
  "synced_at": "2025-10-29T14:30:00Z",
  "stats": {
    "word_count": 45678,
    "chapter_count": 12,
    "last_commit_sha": "abc123def456"
  }
}
```

---

## Analytics Endpoints

### Get Project Stats

```
GET /api/v1/projects/{project_id}/stats
Authorization: Bearer {access_token}
```

**Query Parameters:**
- `period` (optional): `day`, `week`, `month`, `all` (default: `month`)

**Response:**
```json
{
  "word_count": {
    "current": 45678,
    "history": [
      { "date": "2025-10-01", "count": 35000 },
      { "date": "2025-10-08", "count": 38500 },
      { "date": "2025-10-15", "count": 42000 },
      { "date": "2025-10-29", "count": 45678 }
    ]
  },
  "chapter_count": 12,
  "writing_streak": {
    "current": 7,
    "longest": 21
  },
  "ai_usage": {
    "tokens_used": 125000,
    "cost_cents": 250,
    "reviews_run": 15
  }
}
```

---

## Credits & Billing Endpoints

### Get Credit Balance

```
GET /api/v1/credits
Authorization: Bearer {access_token}
```

**Response:**
```json
{
  "remaining": 150000,
  "total": 200000,
  "subscription_tier": "creator",
  "subscription_expires_at": "2025-11-01T00:00:00Z"
}
```

---

### Get Usage Log

```
GET /api/v1/credits/usage
Authorization: Bearer {access_token}
```

**Query Parameters:**
- `start_date` (optional): ISO 8601 date
- `end_date` (optional): ISO 8601 date

**Response:**
```json
{
  "usage": [
    {
      "id": 1,
      "agent_type": "continuity",
      "model_name": "claude-3-5-sonnet",
      "tokens_used": 12500,
      "cost_cents": 25,
      "project_name": "My Novel",
      "created_at": "2025-10-29T10:00:00Z"
    }
  ],
  "total_tokens": 125000,
  "total_cost_cents": 250
}
```

---

## Error Responses

All errors follow this format:

```json
{
  "error": true,
  "message": "Human-readable error message",
  "code": "ERROR_CODE",
  "details": {
    "field": "validation error details"
  }
}
```

**Common Error Codes:**
- `UNAUTHORIZED` (401): Invalid or missing JWT token
- `FORBIDDEN` (403): User doesn't have access to resource
- `NOT_FOUND` (404): Resource not found
- `VALIDATION_ERROR` (422): Invalid request body
- `RATE_LIMIT_EXCEEDED` (429): Too many requests
- `INSUFFICIENT_CREDITS` (402): Not enough AI credits
- `INTERNAL_ERROR` (500): Server error

---

## Rate Limiting

- **Free tier:** 100 requests/hour
- **Creator tier:** 500 requests/hour
- **Professional tier:** 2000 requests/hour

Rate limit headers included in all responses:
```
X-RateLimit-Limit: 500
X-RateLimit-Remaining: 487
X-RateLimit-Reset: 1698764400
```

---

## Webhooks (Future)

DraftForge will send webhooks for events like:
- `agent.run.completed`
- `project.stats.updated`
- `credits.low`

Webhook payload format:
```json
{
  "event": "agent.run.completed",
  "timestamp": "2025-10-29T10:01:15Z",
  "data": {
    "project_id": 1,
    "run_id": 1,
    "agent_type": "continuity",
    "status": "completed"
  }
}
```

---

## SDK & Client Libraries

### Go Client (Official)

```go
import "github.com/yourusername/draftforge-go"

client := draftforge.NewClient("your-api-key")
projects, err := client.Projects.List()
```

### TypeScript Client (Official)

```typescript
import { DraftForgeClient } from '@draftforge/client';

const client = new DraftForgeClient({ apiKey: 'your-api-key' });
const projects = await client.projects.list();
```

---

## API Versioning

- Current version: **v1**
- Version is specified in URL path: `/api/v1/...`
- Breaking changes will increment major version (v2, v3, etc.)
- Old versions supported for 12 months after deprecation notice

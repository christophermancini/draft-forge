# DraftForge

**Forge your draft. Keep your voice.**

An AI-assisted authoring platform for writers who want control over their craft. DraftForge helps authors organize, refine, and publish long-form creative works using Markdown, Git, and AI editorial agents—without sacrificing their authentic voice.

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.24-blue.svg)](https://golang.org)
[![SvelteKit](https://img.shields.io/badge/svelte-5.1-orange.svg)](https://kit.svelte.dev)

---

## 🎯 Philosophy

- **Authors Own Everything** - Projects live in your GitHub account, no vendor lock-in
- **AI Assists, Never Replaces** - AI provides editorial feedback, you make creative decisions
- **Git-Native** - Version control is built-in, not bolted-on
- **Multi-Model Flexibility** - Choose AI models based on task and budget

---

## 🚀 Quick Start

### Prerequisites

- **Go 1.24+** ([install](https://golang.org/dl/))
- **Node.js 20+** ([install](https://nodejs.org/))
- **Docker** ([install](https://docs.docker.com/get-docker/))
- **Task** ([install](https://taskfile.dev/installation/))
- **GitHub Account** (for OAuth and repository integration)

### Installation

1. **Clone the repository**

```bash
git clone https://github.com/yourusername/draftforge.git
cd draftforge
```

2. **Set up environment variables**

```bash
cp .env.local.example .env.local
```

Edit `.env.local` and fill in your configuration:

```env
# Database
DATABASE_URL=postgres://postgres:password@localhost:5432/draftforge?sslmode=disable

# GitHub OAuth (create app at https://github.com/settings/developers)
GITHUB_CLIENT_ID=your_client_id_here
GITHUB_CLIENT_SECRET=your_client_secret_here

# JWT Secrets (generate with: openssl rand -base64 32)
JWT_SECRET=your_jwt_secret_here
REFRESH_TOKEN_SECRET=your_refresh_secret_here

# AI Services
OPENROUTER_API_KEY=your_openrouter_key_here
```

3. **Initialize the project**

```bash
task setup
```

This will:

- Install Go dependencies
- Start PostgreSQL in Docker
- Run database migrations
- Install frontend dependencies

4. **Start development environment**

```bash
task dev
```

This starts:

- PostgreSQL database on port 5432
- Go API server on port 8080
- SvelteKit frontend on port 5173

5. **Open the application**

Navigate to [http://localhost:5173](http://localhost:5173)

---

## 📁 Project Structure

```
draftforge/
├── cmd/
│   ├── api/                # Main API application
│   └── cli/                # CLI tools (migrations, etc.)
├── internal/               # Private application code
│   ├── auth/               # Authentication & JWT handling
│   ├── projects/           # Project management
│   ├── ai/                 # AI agent orchestration
│   ├── github/             # GitHub API integration
│   ├── db/                 # Database utilities & migrations
│   └── scaffold/           # Project scaffolding
├── frontend/               # SvelteKit web application
│   ├── src/
│   │   ├── routes/         # Pages & endpoints
│   │   └── lib/            # Components & utilities
│   └── package.json
├── docs/                   # Documentation
│   ├── api-design.md       # API specification
│   ├── architecture.md     # System architecture
│   └── market-analysis.md  # Market research
├── agents/                 # AI agent definitions
├── scaffolds/              # Project templates
├── infra/                  # Infrastructure (Terraform, etc.)
├── go.mod                  # Go dependencies
├── Taskfile.yaml           # Task runner configuration
├── CLAUDE.md               # AI collaboration guidelines
└── README.md               # This file
```

---

## 🛠️ Development

### Available Tasks

View all available tasks:

```bash
task --list
```

**Common tasks:**

```bash
# Setup and initialization
task setup              # Initialize project (deps, DB, migrations)

# Development
task dev                # Start full dev environment
task api:dev            # Start only the API server
task frontend:dev       # Start only the frontend

# Database
task db:up              # Start PostgreSQL
task db:down            # Stop PostgreSQL
task db:migrate         # Run pending migrations
task db:migrate-down    # Rollback last migration
task db:reset           # Drop, recreate, and migrate DB
task db:console         # Open psql console

# Building
task go:build           # Build Go binaries
task frontend:build     # Build frontend for production

# Testing
task go:test            # Run Go tests
task frontend:test      # Run frontend tests
task test               # Run all tests

# Code quality
task go:lint            # Run golangci-lint
task go:fmt             # Format Go code
```

### Database Migrations

Create a new migration:

```bash
# Create migration files
touch backend/pkg/db/migrations/002_add_feature.up.sql
touch backend/pkg/db/migrations/002_add_feature.down.sql
```

Run migrations:

```bash
task db:migrate
```

Rollback:

```bash
task db:migrate-down
```

### Adding a New API Endpoint

1. Define handler in `internal/<domain>/handlers.go`
2. Add route in `cmd/api/main.go`
3. Update `docs/api-design.md`
4. Write tests in `internal/<domain>/handlers_test.go`

Example:

```go
// internal/projects/handlers.go
func GetProject(c *fiber.Ctx) error {
    id := c.Params("id")
    // ... implementation
    return c.JSON(project)
}

// cmd/api/main.go
api.Get("/projects/:id", projects.GetProject)
```

---

## 🧪 Testing

### Backend Tests

```bash
go test ./... -v
```

### Frontend Tests

```bash
cd frontend
npm run test
```

### Integration Tests

```bash
task test
```

---

## 🔐 GitHub Setup

### 1. Create GitHub OAuth App

1. Go to [GitHub Developer Settings](https://github.com/settings/developers)
2. Click "New OAuth App"
3. Fill in:
   - **Application name:** DraftForge Local
   - **Homepage URL:** `http://localhost:5173`
   - **Authorization callback URL:** `http://localhost:8080/api/v1/auth/github/callback`
4. Copy Client ID and Client Secret to `.env.local`

### 2. Create GitHub App (Optional, for repo management)

1. Go to [GitHub Apps](https://github.com/settings/apps)
2. Click "New GitHub App"
3. Fill in:
   - **GitHub App name:** DraftForge
   - **Homepage URL:** `http://localhost:5173`
   - **Webhook URL:** `http://localhost:8080/api/v1/webhooks/github`
   - **Permissions:**
     - Repository contents: Read & write
     - Repository metadata: Read-only
     - Workflows: Read & write
4. Generate and download private key
5. Save as `backend/github-app-private-key.pem`
6. Copy App ID to `.env.local`

---

## 📊 Database Schema

```sql
users
├── id (PK)
├── github_id (unique)
├── username
├── email
├── avatar_url
└── created_at

projects
├── id (PK)
├── user_id (FK)
├── name
├── slug
├── project_type (novel, screenplay, etc.)
├── github_repo_url
└── settings (JSONB)

ai_credits
├── id (PK)
├── user_id (FK)
├── credits_remaining
├── subscription_tier
└── subscription_expires_at

ai_usage_log
├── id (PK)
├── user_id (FK)
├── project_id (FK)
├── agent_type
├── model_name
├── tokens_used
└── cost_cents

agent_runs
├── id (PK)
├── project_id (FK)
├── agent_type
├── status
└── results (JSONB)
```

See [`internal/db/migrations/001_initial_schema.up.sql`](internal/db/migrations/001_initial_schema.up.sql) for full schema.

---

## 🤖 AI Agents

DraftForge includes several specialized AI agents:

- **ContinuityBot** - Checks for character, world, and plot consistency
- **StyleBot** - Analyzes voice, tense, and readability
- **TimelineBot** - Validates chronological consistency
- **FactBot** - Verifies factual accuracy (for non-fiction)

See [docs/ai-agents.md](docs/ai-agents.md) for detailed documentation.

---

## 🌐 API Documentation

Full API documentation: [docs/api-design.md](docs/api-design.md)

**Base URL:** `http://localhost:8080/api/v1`

**Key endpoints:**

- `POST /auth/github` - GitHub OAuth login
- `GET /projects` - List user's projects
- `POST /projects` - Create new project
- `POST /projects/:id/agents/run` - Trigger AI agent

---

## 🗺️ Roadmap

See [docs/roadmap.md](docs/roadmap.md) for the path from PoC to MVP with actionable tasks and references.

---

## 🚢 Deployment

### Backend (DigitalOcean App Platform)

```bash
# Build binary
task api:build

# Deploy (configure DO App Platform to use ./bin/api)
```

### Frontend (Cloudflare Pages)

```bash
# Build frontend
task frontend:build

# Deploy
wrangler pages deploy frontend/build
```

### Database (DigitalOcean Managed PostgreSQL)

Set production `DATABASE_URL` in environment.

---

## 🤝 Contributing

This project is currently in private development. See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

---

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

---

## 📞 Support

- **Documentation:** [docs/](docs/)
- **Issues:** [GitHub Issues](https://github.com/yourusername/draftforge/issues)
- **Email:** support@draftforge.io

---

## 🙏 Acknowledgments

- Built with [Go Fiber](https://gofiber.io/)
- Frontend powered by [SvelteKit](https://kit.svelte.dev/)
- AI routing via [OpenRouter](https://openrouter.ai/)
- Designed with help from Claude (Anthropic)

---

**Made with ❤️ for writers who code (or coders who write)**

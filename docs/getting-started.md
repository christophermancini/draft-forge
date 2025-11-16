# Getting Started Guide

Welcome to DraftForge! This guide will help you set up your development environment and understand the project structure.

---

## Prerequisites Checklist

Before starting, ensure you have:

- [ ] **Go 1.24+** installed and in PATH
- [ ] **Node.js 20+** installed
- [ ] **Docker Desktop** running
- [ ] **Task** CLI tool installed
- [ ] **GitHub account** with OAuth app created
- [ ] **Code editor** (VS Code recommended)

---

## Initial Setup (5-10 minutes)

### 1. Clone and Configure

```bash
# Clone repository
git clone https://github.com/yourusername/draftforge.git
cd draftforge

# Create environment file
cp .env.local.example .env.local
```

### 2. Set Up GitHub OAuth

1. Visit [GitHub Developer Settings](https://github.com/settings/developers)
2. Click "New OAuth App"
3. Fill in details:
   ```
   Application name: DraftForge Local
   Homepage URL: http://localhost:5173
   Authorization callback URL: http://localhost:8080/api/v1/auth/github/callback
   ```
4. Copy Client ID and Client Secret

### 3. Configure Environment

Edit `.env.local`:

```env
# Required - Update these values
GITHUB_CLIENT_ID=your_github_client_id_here
GITHUB_CLIENT_SECRET=your_github_client_secret_here

# Generate secrets (run: openssl rand -base64 32)
JWT_SECRET=generate_a_random_secret
REFRESH_TOKEN_SECRET=generate_another_random_secret

# Optional - For AI features
OPENROUTER_API_KEY=sk-or-v1-...

# Leave defaults for local development
DATABASE_URL=postgres://postgres:password@localhost:5432/draftforge?sslmode=disable
API_PORT=8080
CORS_ORIGINS=http://localhost:5173
```

### 4. Initialize Project

```bash
task setup
```

This command will:
- ✅ Install Go dependencies
- ✅ Start PostgreSQL in Docker
- ✅ Run database migrations
- ✅ Install frontend dependencies

**Expected output:**
```
task: [go:mod] go mod download
task: [db:up] docker run -d --name draftforge-db ...
task: [db:migrate] Running migrations...
Migrations completed successfully
task: [frontend:install] npm install
```

### 5. Start Development Server

```bash
task dev
```

This starts three services:
- 🟢 **PostgreSQL** on `localhost:5432`
- 🟢 **Go API** on `http://localhost:8080`
- 🟢 **SvelteKit Frontend** on `http://localhost:5173`

### 6. Verify Setup

Open [http://localhost:5173](http://localhost:5173) in your browser.

You should see the DraftForge homepage with a "Login with GitHub" button.

---

## Project Structure Explained

```
```
draftforge/
├── cmd/
│   ├── api/               # Main HTTP server
│   │   └── main.go        # Entry point
│   └── cli/               # CLI tools (migrations)
├── internal/              # Private packages
│   ├── auth/              # JWT, OAuth handlers
│   ├── projects/          # Project CRUD
│   ├── ai/                # AI agent system
│   ├── github/            # GitHub API client
│   ├── db/                # Database utilities & migrations
│   └── scaffold/          # Project templates
├── go.mod
└── go.sum
```
├── frontend/                   # SvelteKit app
│   ├── src/
│   │   ├── routes/            # Pages (file-based routing)
│   │   │   ├── +page.svelte   # Homepage
│   │   │   └── +layout.svelte # Root layout
│   │   ├── lib/               # Shared code
│   │   │   └── components/    # Reusable components
│   │   ├── app.css            # Global styles
│   │   └── app.html           # HTML template
│   ├── package.json
│   ├── svelte.config.js
│   └── vite.config.ts
│
├── docs/                       # Documentation
│   ├── api-design.md          # API specification
│   ├── architecture.md        # System design
│   ├── getting-started.md     # This file
│   └── market-analysis.md     # Market research
│
├── Taskfile.yaml              # Task definitions
├── .env.local.example         # Environment template
├── .gitignore
├── README.md
└── CLAUDE.md                  # AI collaboration guide
```

---

## Common Development Tasks

### Database Operations

```bash
# Start PostgreSQL
task db:up

# Stop PostgreSQL
task db:down

# Run migrations
task db:migrate

# Rollback last migration
task db:migrate-down

# Reset database (drops all data!)
task db:reset

# Open PostgreSQL console
task db:console
```

### Backend Development

```bash
# Start API with hot reload
task api:dev

# Build binary
task go:build

# Run tests
task go:test

# Lint code
task go:lint

# Format code
task go:fmt
```

### Frontend Development

```bash
# Start dev server
task frontend:dev

# Build for production
task frontend:build

# Preview production build
task frontend:preview

# Run tests
task frontend:test
```

### Full Stack

```bash
# Start everything
task dev

# Run all tests
task test

# Clean build artifacts
task clean
```

---

## Making Your First API Call

### 1. Start the server

```bash
task dev
```

### 2. Health check

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "ok",
  "service": "draftforge-api"
}
```

### 3. API version

```bash
curl http://localhost:8080/api/v1/
```

Expected response:
```json
{
  "message": "DraftForge API v1"
}
```

---

## Testing Authentication Flow

### 1. Start the application

```bash
task dev
```

### 2. Open browser

Navigate to [http://localhost:5173](http://localhost:5173)

### 3. Click "Login with GitHub"

You'll be redirected to GitHub for authorization.

### 4. Authorize the app

After authorizing, you'll be redirected back to DraftForge with:
- JWT access token (stored in localStorage)
- User information

### 5. Verify authentication

Check browser DevTools → Application → Local Storage:
- `access_token`: Should contain JWT
- `user`: Should contain your GitHub info

---

## Creating Your First Database Migration

### 1. Create migration files

```bash
# Create new migration
touch internal/db/migrations/002_add_feature.up.sql
touch internal/db/migrations/002_add_feature.down.sql
```

### 2. Write migration

**002_add_feature.up.sql:**
```sql
-- Add new column to users table
ALTER TABLE users ADD COLUMN bio TEXT;
```

**002_add_feature.down.sql:**
```sql
-- Remove column (rollback)
ALTER TABLE users DROP COLUMN bio;
```

### 3. Run migration

```bash
task db:migrate
```

### 4. Verify

```bash
task db:console
# Then in psql:
\d users
```

---

## Adding a New API Endpoint

### Example: Get user profile

**1. Define handler** (`internal/auth/handlers.go`):

```go
package auth

import "github.com/gofiber/fiber/v2"

func GetProfile(c *fiber.Ctx) error {
    // TODO: Get user from JWT
    return c.JSON(fiber.Map{
        "message": "Profile endpoint",
    })
}
```

**2. Register route** (`cmd/api/main.go`):

```go
import "github.com/yourusername/draftforge/internal/auth"

// In main():
api := app.Group("/api/v1")
api.Get("/profile", auth.GetProfile)
```

**3. Test**

```bash
curl http://localhost:8080/api/v1/profile
```

**4. Document** in `docs/api-design.md`

---

## Troubleshooting

### Database won't start

**Error:** `port 5432 already in use`

**Solution:**
```bash
# Check if PostgreSQL is already running
docker ps

# Stop existing container
docker stop draftforge-db
docker rm draftforge-db

# Try again
task db:up
```

### Go dependencies not found

**Error:** `could not import github.com/...`

**Solution:**
```bash
cd backend
go mod tidy
go mod download
```

### Frontend won't start

**Error:** `Cannot find module '@sveltejs/kit'`

**Solution:**
```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
```

### Port conflicts

**Error:** `address already in use`

**Solution:**
```bash
# Check what's using the port
# Windows:
netstat -ano | findstr :8080

# Kill the process or change port in .env.local
API_PORT=8081
```

### Migration fails

**Error:** `migration failed`

**Solution:**
```bash
# Check migration syntax
cat internal/db/migrations/XXX_name.up.sql

# Reset database (WARNING: deletes all data)
task db:reset
```

---

## Next Steps

Now that you have DraftForge running locally:

1. **Read the documentation:**
   - [API Design](./api-design.md) - API endpoints
   - [Architecture](./architecture.md) - System design
   - [CLAUDE.md](../CLAUDE.md) - Project philosophy

2. **Explore the codebase:**
   - Start with `cmd/api/main.go`
   - Check out the database schema in migrations
   - Look at frontend routes in `frontend/src/routes/`

3. **Make your first change:**
   - Add a new API endpoint
   - Create a new database migration
   - Build a new frontend component

4. **Join the development:**
   - Check open issues
   - Read [CONTRIBUTING.md](../CONTRIBUTING.md)
   - Submit your first PR

---

## Getting Help

- **Documentation:** Check `docs/` directory
- **GitHub Issues:** [Report bugs or request features](https://github.com/yourusername/draftforge/issues)
- **Code Examples:** See `docs/api-design.md` for request/response examples

---

**Welcome to DraftForge! Happy coding! 🚀**

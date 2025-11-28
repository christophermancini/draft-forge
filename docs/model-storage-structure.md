# Model & Storage Structure Proposal

Goal: separate domain models (used across services/handlers) from DB-facing structs, keep storage-specific details contained, and reduce scanning boilerplate.

## Principles

- Domain models live in `internal/models` with JSON tags only (no DB hints or nullable wrappers).
- DB structs live in per-table packages under `internal/db` (e.g., `internal/db/project`), with SQL/x tags and null handling contained there.
- Converters translate `db.Project` ↔ `models.Project`. Storage-specific logic (NullString, IDs, repo info) stays in the DB package.
- Stores expose/accept `models.*`, hiding DB details from callers.
- Prefer `sqlx` for simpler scanning/named args; still raw SQL, but less boilerplate.

## Package Layout (proposed)

```
internal/models/
  project.go      // Project model (JSON tags)
  user.go         // User model
  agent_run.go    // AgentRun model, etc.

internal/db/project/
  dto.go          // DBProject struct with db tags, nullable fields
  convert.go      // ToModel/FromModel helpers
  store.go        // Uses sqlx, returns models.Project

internal/db/user/
  ...

internal/db/agent/
  ...
```

## Usage Pattern

- Service layer depends on stores that accept/return `models.*`.
- Handler → Service → Store: only models cross boundaries; storage details don’t leak.
- Adding caches (e.g., Redis) is easy because models are DB-agnostic.

## API Payloads & Metadata (proposed)

- Use JSON:API-like envelopes for consistency and SDK friendliness:
  - Single: `{ "data": <model>, "meta": { ...optional... } }`
  - List: `{ "data": [<model>...], "meta": { "total": n, "page": p, "per_page": k, "next_cursor": "...?" } }`
- Surface response-only fields (e.g., `repo_url`, `scaffold_path`) in `meta` or as siblings inside `data` when they are part of the model’s response contract, but do not bake storage-only fields into domain models.
- Errors: HTTP status + `{ "error": { "message": "...", "code": "...", "details": ...? } }`.
- Keeps models reusable for SDKs/caches while allowing pagination/telemetry without breaking clients.

## sqlx Decision

- Adopt `sqlx` to reduce manual scanning/null handling. Still raw SQL (migration-friendly), minimal dependency cost.
- Keep queries in stores; converters handle nullable fields.

## Migration Steps

1. Add `internal/models` for Project/User/AgentRun with JSON tags.
2. Create `internal/db/project` (dto + convert + store) using sqlx; replace current project store with the new one.
3. Update services/handlers to use models.Project (no DB structs).
4. Repeat for User and AgentRun stores.
5. Remove old store structs once migrated.

## Notes

- Repo metadata stays in models.Project.GitHubRepo; DB dto handles nullable columns.
- Keep migrations unchanged; only Go layering changes.

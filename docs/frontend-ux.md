# Frontend UI/UX Plan (PoC)

Goal: deliver a usable, consistent UI for auth, project creation, and agent runs, aligned with the data/meta API envelopes. Keep components small, typed, and testable; prefer clarity over abstraction until MVP.

## Architecture

- Framework: SvelteKit with TypeScript.
- State: local component state + lightweight stores (`svelte/store`) for auth tokens/user/profile; avoid heavy global state.
- HTTP: fetch helpers that unwrap `{ data, meta }` envelopes and handle errors centrally.
- Layout: shared shell with top nav (brand, user menu), content area for views.

## API Conventions

- Responses: `{ data: ..., meta: ... }`; errors: `{ error: { message, code?, details? } }`.
- Auth: `/auth/github/start` (redirect URL from `data.auth_url`), `/auth/github/callback`, `/me`.
- Projects: `GET /projects` -> `data: Project[]`, `POST /projects` -> `data: Project`, `meta.repo_url`, `meta.scaffold_path`.
- Agents: `POST /projects/:id/agents/run` -> `data: AgentRun`, `meta.message`; `GET /projects/:id/agents/runs/:run_id` -> `data: AgentRun` (add list endpoint when ready).

## Key Screens (PoC)

1. **Landing/Auth**
   - GitHub login button (calls `/auth/github/start` and redirects to callback).
   - After callback, store JWT, fetch `/me` (shows username/avatar).

2. **Projects List**
   - Table/cards showing `name`, `slug`, `project_type`, `github_repo?.url`.
   - “New Project” form modal: fields `name`, `project_type`, `template` (default `novel`), `use_github` toggle, optional `github_owner`.
   - Show `meta.repo_url` and `meta.scaffold_path` on create success.

3. **Project Detail (lightweight)**
   - Show project info; button to trigger agent run (Continuity by default).
   - List recent agent runs (once list endpoint exists); show status, agent_type, link to artifact if available.

## Component Structure

- `src/lib/api/client.ts`: fetch wrapper to handle envelopes, inject JWT, parse errors.
- `src/lib/stores/auth.ts`: writable store for token + user (from `/me`).
- `src/lib/components`: buttons, forms, modals; project list item; agent run card.
- `src/routes/+layout.svelte`: shell with nav; handles redirect to login if no token (except public routes).
- `src/routes/+page.svelte`: landing/login.
- `src/routes/projects/+page.svelte`: list + create form.
- `src/routes/projects/[slug]/+page.svelte`: detail + agent trigger (stub until list endpoint).

## Styling & Quality

- Keep it clean and minimal: consistent spacing, neutral palette, clear typography; avoid unstyled defaults.
- Form validation on required fields; inline error messages on failed API calls.
- Type all API responses; centralize models in `src/lib/types.ts`.
- Tests: prefer component tests for forms and API helper tests; ensure critical flows (login, project create) have coverage.

## Next Implementation Steps

1. Add API client wrapper to unwrap `{ data, meta }` and raise errors.
2. Build auth flow (login button → callback handler → store JWT → fetch `/me`).
3. Implement projects list/create UI using new fields (`template`, `use_github`, `github_owner`); display repo/scaffold meta.
4. Add project detail stub with agent run trigger (list endpoint to follow).

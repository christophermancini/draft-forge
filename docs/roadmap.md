# Roadmap: PoC → MVP

Actionable phases to reach a working platform. Each task references an existing design doc for context.

## Phase 0: Proof of Concept (PoC)

Goal: Single-project happy path with GitHub OAuth, repo scaffold, and one AI agent run on PR.

1. Auth + Session

- Implement GitHub OAuth login (see docs/architecture.md#project-creation-flow).
- Persist user/session tokens; stub minimal user model in DB (docs/data-model-design.md#users).

2. Project Scaffold + Repo Create

- API to create project and seed repo from template; commit base tree and Actions (docs/architecture.md#repository-structure).
- Support at least one scaffold (e.g., novel) with chapters/manuscript/.draftforge config.
- Task entrypoint: Taskfile target to scaffold locally for e2e testing (docs/getting-started.md).

3. Agent Run (single agent)

- Implement queue + worker for one agent (ContinuityBot) with stubbed AI call (can be mocked) (docs/architecture.md#ai-agent-system).
- Expose `POST /projects/{id}/agents/run` and `GET /projects/{id}/agents/runs/{run_id}` per docs/api-design.md#ai-agent-endpoints.
- Persist run in `agent_runs` table (docs/data-model-design.md#6-agent-runs); write JSON artifact to `.draftforge/agent-runs/` (docs/data-storage-strategy.md#agent-results).

4. GitHub Action Trigger (PR)

- Add `.github/workflows/ai-review.yml` that calls the API to queue the agent on PR open (docs/architecture.md#github-actions-integration).
- PoC can post results to logs or commit artifact; PR comment optional.

5. Minimal Frontend UI

- Simple page to create project and list agent runs (status/results) (docs/architecture.md#ai-agent-system).
- No editor needed; just auth + list.

Exit criteria: Login works, project scaffolded to GitHub, PR triggers a ContinuityBot queue/run, result JSON written, and UI shows run status.

## Phase 1: MVP

Goal: Multi-project support, two agents, basic dashboard, and GitHub PR comments.

1. Data Model + Credits

- Implement core tables: users, projects, agent_runs, ai_usage_log (docs/data-model-design.md#phase-1-current).
- Track tokens/cost per run; surface in API responses.

2. Agent System Hardening

- Add StyleBot alongside ContinuityBot; model routing config in `.draftforge/agents.yml` (docs/architecture.md#ai-agent-system).
- Retry/backoff, dead-letter handling, max runtime cancel (docs/data-model-design.md#6-agent-runs).

3. PR Comment + Checks

- Post summarized results as PR comment and/or GitHub Check Run (docs/api-design.md#ai-agent-endpoints).
- Include link to `.draftforge/agent-runs/<run>.json` in comment.

4. Automation Catalog Toggles

- Frontend switchboard to enable/disable workflows (stats, compile, ai-review) and commit YAML to `.github/workflows/` (docs/architecture.md#automation-catalog-ux).
- Backend writes/updates workflow files and records settings in `.draftforge/config.yml`.

5. Dashboard + Project List

- UI for multiple projects, showing last runs, credits, and workflow status (docs/architecture.md#project-creation-flow).
- Filter/sort agent runs by status/type.

6. CI/Lint + Packaging

- Taskfile targets for backend lint/test (`golangci-lint`, `go test ./...`) and frontend lint/build (`npm run lint && npm run check`) (docs/getting-started.md).
- Add GH Action for API/Frontend lint/test on PR (docs/architecture.md#github-actions-integration).

Exit criteria: Users can manage multiple projects, toggle workflows, run ContinuityBot + StyleBot via PR or manual trigger, see comments on PRs, and monitor runs/credits in UI.

## Phase 2: Nice-to-Haves (post-MVP)

- Audio/SSML export workflow toggle (docs/ai-agents.md#future-extensions).
- TimelineBot scheduling (weekly) and result visualization.
- FactBot for non-fiction; adds citation hints.
- Evidence/Legal export templates (DRAFTFORGE.md#35-legal--release-toolkit).
- Pricing/usage surfaced in UI; webhook/email notifications for `agent.run.completed`.

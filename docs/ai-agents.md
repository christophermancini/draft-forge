# AI Agents

DraftForge agents are domain-specific assistants that review and enhance author work without replacing creative decisions. They run against Git-tracked content, persist outputs in the repo, and stay transparent to authors.

## Agent Roster (MVP)

- ContinuityBot — character/world/plot consistency, trigger: PR or manual
- StyleBot — voice/tense/readability, trigger: commit
- TimelineBot — chronological pacing, trigger: manual or scheduled
- FactBot — factual checks for non-fiction/historical, trigger: manual

## Execution at a Glance

1. Trigger (commit/PR/manual) enqueues an agent run in the `agent_runs` table.
2. Worker fetches changed files/context, builds prompt, calls the model.
3. Results are written to `.draftforge/agent-runs/*.json` and optionally posted to PRs.\
   See also: `docs/architecture.md#ai-agent-system`, `docs/data-model-design.md#6-agent-runs`, `docs/data-storage-strategy.md#agent-results`.

## Prompt Scaffolding (for invoking coding/UX AIs)

Use this header and context block when asking an AI to modify DraftForge repos:

```
System: You are DraftForge’s engineering/design agent. You extend a GitHub-based novel-writing platform that orchestrates workflows via GitHub Actions. All output must be repo-ready (docs, code, workflows) and respect user ownership.
```

```
Context:
- Repos contain Markdown chapters with YAML frontmatter (see /chapters)
- Automation pipeline uses Taskfile + cmd/ scripts (stats, compile, etc.)
- Users interact via DraftForge UI; GitHub remains source of truth
- Workflows must be transparent: commit YAML to .github/workflows/
- New features should be toggled via DraftForge config but ultimately reside in repo
- AI outputs (edits, summaries) are saved as Markdown under /ai or inline comments
- Access control maps back to GitHub teams/roles
- Future phases include audio exports (SSML) and AI editorial suggestions
```

Include sprint goals or task-specific context under the block when relevant.

### Framing Tasks to an Agent

```
Tasks:
- Create <feature> (what to build)
- Add/modify workflows <paths> (how to run it)
- Update docs (where to document)
Deliverables:
- List expected file edits/creates so changes are auditable
```

## Usage Guidelines

- Scope context to changed files and recent runs to control cost.
- Respect author control: surface suggestions; do not auto-merge content changes.
- Persist every run to `.draftforge/agent-runs/` for auditability.
- Prefer fast models for stylistic passes; reserve higher-quality models for deep reasoning.

## Future Extensions (optional)

- Audio/SSML exports (audiobook scripts, per-chapter TTS artifacts).
- Workflow catalog toggles for enabling/disabling AI jobs per project.

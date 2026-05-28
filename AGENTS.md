# AGENTS.md

This repository builds **Headlamp Agent**, a local-first desktop operator tool.

Headlamp is the tool a user uses to operate their own machine through a simple chat interface, local LLMs, governed local tools, and explicit approval flows.

## Non-negotiable product direction

Do not couple this project to `headlamp.dev` or Kubernetes. That is a separate project with the same word in its name.

Do not turn this into a cloud-first SaaS platform.

Do not make the agent the operator. The human remains the operator. Headlamp is the tool/interface, and agents are capabilities inside the tool.

## Product principles

1. Local-first by default.
2. Desktop-native, lightweight, and useful quickly.
3. Simple chat UI first, advanced dashboards later.
4. Local LLM support first through OpenAI-compatible endpoints.
5. Read/observe before act.
6. Ask approval before sensitive actions.
7. Every proposed and executed action must be logged.
8. Keep implementation modular so tools, skills, providers, and policies can be added.

## MVP priorities

Build toward this path first:

```txt
User opens desktop app
  → configures local LLM endpoint
  → chats with Headlamp
  → Headlamp can inspect local context
  → Headlamp proposes actions
  → user approves or rejects actions
  → approved actions execute locally
  → outcomes are logged
```

## Preferred technical direction

- Desktop shell: Tauri
- UI: React + TypeScript
- Styling: Tailwind CSS
- Local persistence: SQLite or app-local file store
- Model abstraction: OpenAI-compatible chat completions first
- Tool execution: capability registry with policy checks
- Package manager: pnpm
- Monorepo shape: `apps/` and `packages/`

## Safety requirements

Sensitive actions require explicit approval:

- shell commands
- file writes or deletes
- git commits, pushes, branch changes
- network calls outside configured allowlist
- reading secrets or credential stores
- process termination
- destructive operations

The default posture is read-only until a capability is granted.

## Documentation expectations

When adding a feature, update the relevant docs in `docs/` and keep README aligned with actual code.

Avoid exaggerated claims. Be honest about what is implemented versus planned.

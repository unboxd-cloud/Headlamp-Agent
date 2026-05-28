# Headlamp Agent

**Headlamp Agent** is a local-first desktop operator tool with a simple chat UI and local LLM support.

Headlamp is the tool the user uses to operate their machine: it helps them see the current context, understand what is happening, choose the next action, and execute approved actions safely.

It is designed to run on a user-controlled desktop, work with local models by default, and assist across local files, tools, repositories, scripts, apps, and workflows.

This project is intentionally **not coupled to headlamp.dev** and is not Kubernetes-first. The name is literal and functional: a headlamp is a tool for operating in the dark by lighting up what is directly in front of the user.

## Product intent

Headlamp Agent should be:

- **Local-first**: local model execution by default, with user-owned data and local storage.
- **Desktop-native**: a tool for operating from the user's own computer.
- **Simple**: a clean chat UI with focused operator controls.
- **Action-oriented**: optimized for understanding, deciding, and doing useful work.
- **Governed**: actions are reviewed, logged, and permissioned.
- **Composable**: local tools, model providers, skills, and workflows can be added over time.

## MVP scope

The first version focuses on:

1. Desktop app shell
2. Simple chat interface
3. Local LLM connection
4. Local tool execution boundary
5. Conversation and action history
6. Human approval before sensitive actions
7. Lightweight configuration

## Reference architecture

```txt
Desktop Operator Tool
  ├─ Chat UI
  ├─ Settings
  ├─ Conversation History
  └─ Action Review Panel

Operator Core
  ├─ Message Orchestrator
  ├─ Local LLM Adapter
  ├─ Tool Registry
  ├─ Permission Guard
  ├─ Memory Store
  └─ Audit Log

Local Integrations
  ├─ File system tools
  ├─ Shell command tools
  ├─ Git/repository tools
  ├─ Browser/search tools
  └─ Future MCP tools
```

## Local LLM strategy

The default architecture should support local OpenAI-compatible endpoints so users can bring their own local model runner, such as:

- LM Studio
- Ollama
- llama.cpp server
- Docker Model Runner
- LiteLLM gateway over local models

The app should not assume cloud LLM usage.

## Safety model

Headlamp Agent is read-first and approval-first.

```txt
Observe → Explain → Recommend → Ask Approval → Act → Log Outcome
```

Sensitive actions must require explicit user approval before execution, especially:

- shell commands
- file writes/deletes
- git commits/pushes
- network calls
- credential access
- destructive operations

## Planned structure

```txt
apps/
  desktop/        Desktop operator UI
  operator/       Local operator service

packages/
  core/           Agent orchestration primitives
  llm/            Local model adapters
  tools/          Tool registry and built-in tools
  governance/     Permission and audit controls
  ui/             Shared UI components

docs/
  architecture.md
  local-llm.md
  safety-model.md
  roadmap.md
```

## Development status

This repository has just been initialized. The first milestone is a working local desktop chat app connected to a local LLM endpoint.

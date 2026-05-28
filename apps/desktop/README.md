# Headlamp Desktop

Headlamp Desktop is the local operator dashboard for Headlamp Agent.

This app should be built with Wails so the Go runtime remains the backend and the UI stays lightweight.

## MVP screens

- Setup Plan
- Setup Verify
- Agent Health
- Host Inventory

## Backend methods

The desktop app should expose these Go methods first:

```go
GetSetupPlan()
VerifySetup()
GetAgentHealth()
GetInventory()
```

## Current status

The current DMG packages CLI binaries. The next DMG should package this Wails desktop app once implemented.

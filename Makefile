PODMAN_COMPOSE ?= podman compose
COMPOSE_FILE := deploy/podman/compose.yaml
DIST_DIR := dist
HEADLAMP_BIN := $(DIST_DIR)/headlamp
NODE_AGENT_BIN := $(DIST_DIR)/headlamp-node-agent
NODE_AGENT_CONFIG := /tmp/headlamp-agent/config.json

.PHONY: dev
dev: podman-up run-node-agent

.PHONY: build
build: build-headlamp build-node-agent

.PHONY: build-headlamp
build-headlamp:
	mkdir -p $(DIST_DIR)
	go build -o $(HEADLAMP_BIN) ./cmd/headlamp

.PHONY: build-node-agent
build-node-agent:
	mkdir -p $(DIST_DIR)
	go build -o $(NODE_AGENT_BIN) ./cmd/headlamp-node-agent

.PHONY: setup-plan
setup-plan: build-headlamp
	$(HEADLAMP_BIN) setup plan

.PHONY: setup-verify
setup-verify: build-headlamp
	$(HEADLAMP_BIN) setup verify || true

.PHONY: run-node-agent
run-node-agent: build-node-agent
	mkdir -p /tmp/headlamp-agent
	printf '{\n  "agentId": "agent:local",\n  "hostId": "host:local",\n  "mode": "supervised",\n  "listen": "127.0.0.1:9080"\n}\n' > $(NODE_AGENT_CONFIG)
	$(NODE_AGENT_BIN) --config $(NODE_AGENT_CONFIG)

.PHONY: test-node-agent
test-node-agent:
	curl -s http://127.0.0.1:9080/healthz && echo
	curl -s http://127.0.0.1:9080/inventory && echo

.PHONY: podman-up
podman-up:
	$(PODMAN_COMPOSE) -f $(COMPOSE_FILE) up -d

.PHONY: podman-down
podman-down:
	$(PODMAN_COMPOSE) -f $(COMPOSE_FILE) down

.PHONY: podman-logs
podman-logs:
	$(PODMAN_COMPOSE) -f $(COMPOSE_FILE) logs -f

.PHONY: podman-reset
podman-reset:
	$(PODMAN_COMPOSE) -f $(COMPOSE_FILE) down -v

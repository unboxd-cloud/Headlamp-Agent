PODMAN_COMPOSE ?= podman compose
COMPOSE_FILE := deploy/podman/compose.yaml

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

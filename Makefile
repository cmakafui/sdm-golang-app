# Variables
DOCKER_COMPOSE = docker compose
DOCKER_COMPOSE_FILE = compose.yaml

# Targets
.PHONY: build up down stop restart logs shell

# Build the Docker images
build:
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) build

# Start the Docker containers
up: build
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) up -d

# Stop the Docker containers
stop:
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) stop

# Restart the Docker containers
restart: stop up

# Remove the Docker containers
down:
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) down

# View logs from the Docker containers
logs:
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) logs -f

# Access the shell inside the running container
shell:
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) exec sdm-app /bin/sh

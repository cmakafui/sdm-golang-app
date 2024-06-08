.PHONY: dev-up dev-down dev-logs dev-rebuild dev-restart
.PHONY: up down logs rebuild restart

# Targets for development environment
dev-up:
	docker compose -f compose-dev.yaml up --build -d

dev-down:
	docker compose -f compose-dev.yaml down

dev-logs:
	docker compose -f compose-dev.yaml logs -f

dev-rebuild:
	docker compose -f compose-dev.yaml up --build -d --force-recreate

dev-restart: dev-down dev-up

# Targets for production environment
up:
	docker compose -f compose-prod.yaml up --build -d

down:
	docker compose -f compose-prod.yaml down

logs:
	docker compose -f compose-prod.yaml logs -f

rebuild:
	docker compose -f compose-prod.yaml up --build -d --force-recreate

restart: down up

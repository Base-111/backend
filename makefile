# По умолчанию показываем помощь
.DEFAULT_GOAL := help

# Команда для compose (можно переопределить: make COMPOSE="docker-compose")
COMPOSE ?= docker compose

# Подключаем .env, если он есть
ifneq (,$(wildcard .env))
include .env
export $(shell sed -ne 's/^\([^#[:space:]][^=[:space:]]*\)=.*/\1/p' .env)
endif

# Собираем URL БД из переменных окружения
DB_URL := postgres://${DB_USERNAME}:${DB_PASSWORD}@localhost:${DB_PORT}/${DB_NAME}?sslmode=disable

.PHONY: help create-app restart-app create-migrations delete-migrations \
        ollama-webui ollama-pull

help: ## Show available commands
	@echo "Available commands:"
	@grep -E '^[a-zA-Z0-9_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}'

create-app: ## Запустить сервисы в фоне
	$(COMPOSE) up -d

restart-app: ## Rebuild and reboot services
	$(COMPOSE) up -d --build

create-migrations: ## Commit migrations (up)
	@migrate -database "$(DB_URL)" -path ./migrations up

delete-migrations: ## Rollback migrations (down)
	@migrate -database "$(DB_URL)" -path ./migrations down

ollama-webui: ## Run ollama и webui
	$(COMPOSE) up -d ollama webui

# Примеры использования:
#   make ollama-pull model=llama2
#   make ollama-pull model=llama3
#   make ollama-pull model=phi3
# Актуальные модели: https://github.com/ollama/ollama?tab=readme-ov-file#model-library
ollama-pull: ## Download model ollama: make ollama-pull model=<name>
ifndef model
	$(error Choose model: make ollama-pull model=llama3)
endif
	$(COMPOSE) exec ollama ollama pull $(model)

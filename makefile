.PHONY: ollama-webui
include .$(PWD)/.env

create-app:
	docker-compose up -d

restart-app:
	docker-compose up -d --build

create-migrations:
	migrate -database postgres://${DB_USERNAME}:${DB_PASSWORD}@localhost:${DB_PORT}/${DB_NAME}?sslmode=disable -path ./migrations up

delete-migrations:
	migrate -database postgres://${DB_USERNAME}:${DB_PASSWORD}@localhost:${DB_PORT}/${DB_NAME}?sslmode=disable -path ./migrations down

ollama-webui:
	docker-compose up -d ollama webui

.PHONY: help deps generate build run test clean dev

BINARY_NAME = auth-service
PROTO_DIR = proto
PROTO_FILE = $(PROTO_DIR)/auth.proto
PB_DIR = pkg/pb
DOCKER_IMAGE= auth-service
DB_URL = postgres://authuser:authpass@localhost:5433/authdb?sslmode=disable

help:
	@echo "=== Auth Service Makefile ==="
	@echo "  make deps     		- Установить зависимости"
	@echo "  make generate 		- Сгенерировать код из proto"
	@echo "  make grpc-clean	- Очистка gRPC кода"
	@echo "  make test     		- Запустить тесты"
	@echo "  make dev      		- Полный цикл: deps -> generate -> run"

# ===========================================================================
# Local dev
# ===========================================================================

deps:
	@echo "Устанавливаем зависимости..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Зависимости установлены!"

generate:
	@echo "Генерируем gRPC код..."
	mkdir -p $(PB_DIR)
	protoc --go_out=./$(PB_DIR) --go-grpc_out=./$(PB_DIR) $(PROTO_FILE)
	@echo "Генерация завершена!"

grpc-clean:
	@echo "Очищаем gRPC код..."
	rm -rf $(PB_DIR)/
	@echo "gRPC код очищен!"

test:
	@echo "Запускаем тесты..."
	go test ./... -v

dev: deps generate run

# ===========================================================================
# Docker
# ===========================================================================

compose-up:
	@echo "Запускаем все сервисы через Docker Compose..."
	docker compose up -d
	@echo "Сервисы запущены!"

compose-down:
	@echo "Останавливаем все сервисы..."
	docker compose down

compose restart: compose-down compose-up

docker-dev: deps generate compose-up
		@echo "Dev среда с Docker запущена"

# ===========================================================================
# Database
# ===========================================================================

db-connect:
	@echo "Подключаемся к PostgreSQL..."
	docker exec -it $$(docker compose ps -q postgres) psql -U authuser -d authdb

# ===========================================================================
# Migrations
# ===========================================================================

migrate-up:
	@echo "Применяются миграции..."
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	@echo "Откатываем миграции..."
	migrate -path migrations -database "$(DB_URL)" down

migrate-create:
	@read -p "Введите название миграции: " name; \
	migrate create -seq -ext sql -dir migrations "$$name"
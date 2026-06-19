# Путь к golangci-lint (убедись, что он установлен локально)
LINT_BIN := $(HOME)/go/bin/golangci-lint

## build: Build an application
.PHONY: build_signaling
build_signaling:
	go build --ldflags="-s -w" -v cmd/signaling/main.go

## build: Run an application
.PHONY: run_signaling
run_signaling:
	go run cmd/signaling/main.go

## Как пользоваться:
## make lint — проверить весь проект.
## make lint-pkg PKG=./internal/room — проверить только пакет комнаты.
## make lint-file FILE=internal/chat/client.go — проверить один файл.

.PHONY: lint
lint: ## Запуск линтера на весь проект
	$(LINT_BIN) run ./...

.PHONY: lint-pkg
lint-pkg: ## Запуск линтера на конкретный пакет (используй: make lint-pkg PKG=./internal/chat)
	$(LINT_BIN) run $(PKG)/...

.PHONY: lint-file
lint-file: ## Запуск линтера на файл (используй: make lint-file FILE=main.go)
	$(LINT_BIN) run $(FILE)
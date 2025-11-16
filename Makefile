GO      ?= go
APP     ?= prservice
PKG     ?= ./...

DC        ?= docker-compose
DC_FLAGS  ?= -d --build
DC_SERVICE ?= app

.PHONY: build
build:
	$(GO) build -o $(APP) ./cmd/main.go

.PHONY: test
test:
	$(GO) test $(PKG)

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: tidy
tidy:
	$(GO) mod tidy

.PHONY: fmt
fmt:
	$(GO) fmt $(PKG)

##### DOCKER COMPOSE #####
.PHONY: up
up:
	$(DC) up $(DC_FLAGS)

.PHONY: down
down:
	$(DC) down

.PHONY: logs
logs:
	$(DC) logs -f $(DC_SERVICE)

.PHONY: loadtest
loadtest:
	$(DC) run --rm k6

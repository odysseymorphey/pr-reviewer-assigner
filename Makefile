NAME ?= prservice
DST ?= ./

DC ?= docker-compose
DCFLAGS ?= -d --build
DCSERVICE ?= app-1

.PHONY: build
build:
	go build -o $(NAME) $(DST)

.PHONY: test
test:
	go test ./

.PHONY: lint
lint:
	golangcli-lint run ./

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: up
up:
	$(DC) up $(DCFLACGS)

.PHONY: down
down:
	$(DC) down

.PHONY: logs
logs:
	$(DC) logs -f $(DCSERVICE)
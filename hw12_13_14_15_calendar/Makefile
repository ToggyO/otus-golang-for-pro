BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

ifndef config
	config=./configs/config.development.toml
endif

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run-local: build
	$(BIN) -c $(config)

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/... ./pkg/...

install-lint-deps:
	#(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.49.0

lint: install-lint-deps
	/home/otogushakov/go/bin/golangci-lint run ./...

run-integration-tests:
	@docker-compose -f deployments/docker-compose.yml -f deployments/docker-compose.stage.yml -f deployments/docker-compose.local.yml \
		--project-directory ./configs/ up \
		 --build --abort-on-container-exit
	@docker-compose -f deployments/docker-compose.yml -f deployments/docker-compose.stage.yml -f deployments/docker-compose.local.yml \
		--project-directory ./configs/ down \
        --rmi local \
		--volumes \
		--timeout 60; \

migrate-up:
	$(BIN) migrate -up -c $(config)

migrate-down:
	$(BIN) migrate -down -c $(config)


.PHONY: build run build-img run-img version test lint

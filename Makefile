APP_NAME = feed-api
APP_EXECUTABLE = "./build/$(APP_COMMIT)/$(APP_NAME)"
APP_COMMIT ?= $(shell git rev-parse HEAD)
APP_VERSION ?= v1.0.0

build:
	mkdir -p build/
	go build -ldflags "-X main.commit=$(APP_COMMIT)" -o $(APP_EXECUTABLE) cmd/app/*.go

run:
	go run cmd/app/*.go

test:
	./scripts/unit-test.sh

integration-test:
	./scripts/integration-test.sh

lint:
	./scripts/lint.sh

.PHONY: build run test integration-test lint
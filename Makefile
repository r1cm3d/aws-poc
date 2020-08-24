.PHONY: clean build test

include scripts/env/.env
export

all: unit-test

lint:
	@echo "\nApplying golint\n"
	@golint ./...

fmt:
	@echo "\nFormatting scripts\n"
	@shfmt -w scripts/*sh
	@echo "\nFormatting terraform files"
	@terraform fmt terraform/
	@echo "\nFormatting go files\n"
	@go fmt ./...

clean:
	@echo "\nRemoving localstack container\n"
	@(@docker rm -f aws || \
	  rm terraform/*tfstate*) 2>/dev/null | true

build: clean fmt
	@echo "\nBuilding application\n"
	@go build cmd/main.go

unit-test: build
	@echo "\nRunning unit tests\n"
	@go test -cover -v -short ./...

run-dep: clean
	@echo "\nStarting localstack container and creating AWS local resources\n"
	@docker-compose up -d --force-recreate
	@echo "\nWaiting until localstack be ready"
	@sleep 6
	@echo "\nApplying terraform scripts"
	cd terraform && \
	terraform init && \
	terraform plan && \
	terraform apply -auto-approve

integration-test: run-dep build
	@echo "\nRunning integration tests\n"
	@go test -cover -v -run Integration ./...

test: unit-test integration-test
	@echo "\nRunning tests\n"

codecov: run-dep
	@echo "\nRunning Codecov\n"

run: run-dep build run-local
	@echo "\nRunning locally"
	@go test -race -coverprofile=coverage.txt -covermode=atomic -cover ./...

run-local:
	@echo "\nRunning without building it"
	@go run cmd/main.go

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
	@terraform fmt deployments/terraform/
	@echo "\nFormatting go files\n"
	@goimports -w $(shell find . -type f -name '*.go')

clean:
	@echo "\nRemoving localstack container\n"
	@(@docker rm -f aws || \
	  rm -rf deployments/terraform/*tfstate* deployments/terraform/.terraform) 2>/dev/null | true

build: clean fmt lint
	@echo "\nBuilding application\n"
	@go build cmd/main.go

unit-test: build
	@echo "\nRunning unit tests\n"
	@go test -cover -short ./...

up: clean
	@echo "\nStarting localstack container and creating AWS local resources\n"
	@docker-compose up -d --force-recreate
	@echo "\nWaiting until localstack be ready"
	@until docker inspect --format='{{json .State.Health}}' aws | grep -o healthy; do sleep 1; done
	@echo "\nCleaning AWS resources"
	-@cd scripts && bash cleanup-aws-rs.sh
	@echo "\nApplying terraform scripts"
	-cd deployments/terraform/ && \
	terraform init && \
	terraform destroy  -auto-approve && \
	terraform plan && \
	terraform apply -auto-approve

integration-test: up build
	@echo "\nRunning integration tests\n"
	@go test -cover -run Integration ./...

test: fmt unit-test up integration-test
	@echo "\nRunning tests\n"

codecov: up
	@echo "\nRunning Codecov\n"

run: up build run-local
	@echo "\nRunning locally"
	@go test -race -coverprofile=coverage.txt -covermode=atomic -cover ./...

run-local:
	@echo "\nRunning without building it"
	@go run cmd/main.go

stop:
	@echo "\nStopping containers"
	@docker-compose stop

.PHONY: clean build test

include scripts/env/.env
export

all: run-dep build run
	@echo "\nBuilding application from scratch and running it\n"

lint:
	@echo "\nApplying golint\n"
	@golint ./...

fmt:
	@echo "\nFormatting scripts\n"
	@shfmt -w scripts/*sh
	@echo "\nFormattin go files\n"
	@go fmt ./... 

clean:
	@echo "\nRemoving localstack container\n"
	@(docker rm -f aws && \
 	  rm -rf .localstack) 2>/dev/null | true

build: clean fmt
	@echo "\nBuilding application\n"

unit-test: build
	@echo "\nRunning unit tests\n"

run-dep: clean
	@echo "\nStarting localstack container and creating AWS local resources\n"
	@docker-compose up -d && \
	cd scripts && bash init-aws-rs.sh

integration-test: build
	@echo "\nRunning integration tests\n"

test: unit-test integration-test
	@echo "\nRunning tests\n"

run:
	@echo "\nRunning without building it"

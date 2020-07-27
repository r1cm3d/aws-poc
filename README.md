# aws-poc 

**TL;DR:**
```console
make run
```

## Prerequisites
[![Docker](https://img.shields.io/badge/Docker-19.03.9-blue)](https://www.docker.com/)
[![Docker-compose](https://img.shields.io/badge/Docker--compose-1.25.5-blue)](https://github.com/docker/compose/releases)
[![GNU Make](https://img.shields.io/badge/GNU%20Make-4.2.1-lightgrey)](https://www.gnu.org/software/make/)
[![GNU Bash](https://img.shields.io/badge/GNU%20Bash-4.2.1-lightgrey)](https://www.gnu.org/software/bash/)
[![aws-cli](https://img.shields.io/badge/aws--cli-1.18.95-yellow)](https://github.com/aws/aws-cli)
[![shfmt](https://img.shields.io/badge/shfmt-v3.1.0-lightgrey)](https://github.com/mvdan/sh)

## Table of Contents
* [TL;DR](#aws-poc)
* [Prerequisites](#prerequisites)
* [About the Project](#about-the-project)
* [Getting Started](#getting-started)
* [Testing](#testing)
* [Run](#run)
* [Team](#team)

## About The Project

The goal of this project is consume AWS resources with Go programming language. It consumes two SQS queues, persists in DynamoDB tables according some business rules, it downloads and uploads some files in S3 and posts it again in others SQS queues.

## Getting Started

To run this project locally you must have the technologies as the [prerequisites section](#prerequisites)

### Testing
#### Unit tests
```sh
make unit-test
```

#### Integration tests
```sh
make integration-test
```

#### All tests
```sh
make test
```

### Run
#### Build all dependencies and run 
```sh
make run
```

#### Just run without any build
```sh
make run-local
```

It will run `go run` without build AWS infrastructure locally.   

## Contributing 

This project follows this [style guide](https://github.com/golang/go/wiki/CodeReviewComments#error-strings) and this [package structure](https://github.com/golang-standards/project-layout). To contribute you must follow these standards.

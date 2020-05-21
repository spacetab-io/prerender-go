IMAGE_NAME = spacetabio/prerender-go
IMAGE_VERSION = 0.1.5

deps:
	go mod vendor

build:
	go build -o ./bin/prerender ./cmd/prerender/main.go
.PHONY: build

build_vendor:
	go build -mod=vendor -o ./bin/prerender ./cmd/prerender/main.go
.PHONY: build_vendor

build_for_docker:
	GOOS=linux GOARCH=amd64 go build -o ./bin/prerender ./cmd/prerender/main.go
.PHONY: build_for_docker

build_vendor_for_docker:
	GOOS=linux GOARCH=amd64 go build -mod=vendor -o ./bin/prerender ./cmd/prerender/main.go
.PHONY: build_for_docker

run:
	./bin/prerender
.PHONY: run

run_in_docker:
	docker run \
	--rm \
	--name headless-shell \
	-v $$(pwd)/configuration/:/app/configuration/ \
	-v $$(pwd)/bin/prerender:/app/bin/prerender \
	-v $$(pwd)/pages:/app/pages \
	-w /app \
	alpeware/chrome-headless-trunk \
	./bin/prerender

all: build_vendor
.PHONY: all

all_ubuntu: build_vendor_for_docker
.PHONY: all_ubuntu

go: build run
.PHONY: go

go_docker: deps build_for_docker run_in_docker
.PHONY: go_docker

run_headless_shell:
	docker run -it --rm -p 9222:9222 \
     --name=chrome-headless \
     -v /tmp/chromedata/:/data alpeware/chrome-headless-trunk bash

## -------------------

## lint and test stuff

get_lint_config:
	@[ -f ./.golangci.yml ] && echo ".golangci.yml exists" || ( echo "getting .golangci.yml" && curl -O https://raw.githubusercontent.com/spacetab-io/docker-images-golang/master/linter/.golangci.yml )
.PHONY: get_lint_config

lint: get_lint_config
	golangci-lint run
.PHONY: lint

test-unit:
	go test ./... --race --cover -count=1 -timeout 1s -coverprofile=c.out -v
.PHONY: test-unit

coverage-html:
	go tool cover -html=c.out -o coverage.html
.PHONE: coverage-html

test: deps test-unit coverage-html
.PHONY: test

## -------------------

## image stuff

image_build:
	docker build -t ${IMAGE_NAME}:${IMAGE_VERSION} .

image_push:
	docker push ${IMAGE_NAME}:${IMAGE_VERSION}
IMAGE_NAME = spacetabio/prerender-go
IMAGE_VERSION = v1.2.0

deps:
	go mod vendor

build:
	go build -o ./bin/prerender .
.PHONY: build

build_vendor:
	go build -mod=vendor -o ./bin/prerender .
.PHONY: build_vendor

build_for_docker:
	GOOS=linux GOARCH=amd64 go build -o ./bin/prerender .
.PHONY: build_for_docker

build_vendor_for_docker:
	GOOS=linux GOARCH=amd64 go build -mod=vendor -o ./bin/prerender .
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

# ----
## LINTER stuff start

linter_include_check:
	@[ -f linter.mk ] && echo "linter.mk include exists" || (echo "getting linter.mk from github.com" && curl -sO https://raw.githubusercontent.com/spacetab-io/makefiles/master/golang/linter.mk)

.PHONY: lint
lint: linter_include_check
	@make -f linter.mk go_lint

## LINTER stuff end
# ----

# ----
## TESTS stuff start

tests_include_check:
	@[ -f tests.mk ] && echo "tests.mk include exists" || (echo "getting tests.mk from github.com" && curl -sO https://raw.githubusercontent.com/spacetab-io/makefiles/master/golang/tests.mk)

tests: tests_include_check
	@make -f tests.mk go_tests
.PHONY: tests

tests_html: tests_include_check
	@make -f tests.mk go_tests_html
.PHONY: tests_html

## TESTS stuff end
# ----

## image stuff

image_build:
	docker build -t ${IMAGE_NAME}:${IMAGE_VERSION} .

image_push:
	docker push ${IMAGE_NAME}:${IMAGE_VERSION}

image: image_build image_push
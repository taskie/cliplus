.PHONY: build install test coverage dep-ensure dep-graph pre-commit

build:
	go build -v -ldflags "-s -w -X cliplus.Revision=$(shell git rev-parse --short HEAD)"

install:
	go install -v -ldflags "-s -w -X cliplus.Revision=$(shell git rev-parse --short HEAD)"

test:
	go test

coverage:
	mkdir -p test/coverage
	go test -coverprofile=test/coverage/cover.out
	go tool cover -html=test/coverage/cover.out -o test/coverage/cover.html

pre-commit:
	$(MAKE) build
	$(MAKE) coverage

version=0.1.0

.PHONY: all

all:
	@echo "make <cmd>"
	@echo ""
	@echo "commands:"
	@echo "  build         - build the source code"
	@echo "  lint          - lint the source code"
	@echo "  fmt           - format the code with gofmt"
	@echo "  install       - install dependencies"

lint:
	@go vet ./...
	@golint ./...

fmt:
	@go fmt ./...

build: lint
	@go build ./...

install:
	@go get -u golang.org/x/lint/golint

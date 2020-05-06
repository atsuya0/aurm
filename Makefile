.PHONY: build install

build:
	@goimports -w .
	@go build

install:
	@goimports -w .
	@go install

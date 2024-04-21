PHONY: build
build:
	go build -C ./cmd -o web


PHONY: int
int: 
	go test -C integration -v


PHONY: unit
unit:
	go test -C services -v


PHONY: tests
tests: unit int
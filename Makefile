.PHONY: build
build:
	go build -C ./cmd -o web


.PHONY: int
int: 
	go test -C integration -v


.PHONY: unit
unit:
	go test -C services -v


.PHONY: staticcheck
staticcheck:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	golangci-lint run 


.PHONY: tests
tests: staticcheck unit int


.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v


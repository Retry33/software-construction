.PHONY: run run-clean test cover lint build

run:
	go run ./cmd/vyshkat

run-clean:
	go run ./cmd/vyshkat --no-demo

build:
	go build -o bin/vyshkat ./cmd/vyshkat

test:
	go test ./...

cover:
	go test -coverprofile=coverage.out ./internal/...
	go tool cover -func=coverage.out | tail -1

lint:
	gofmt -l .
	go vet ./...
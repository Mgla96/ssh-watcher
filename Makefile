.PHONY: lint test build

lint:
	golangci-lint run

test:
	GOWORK=off go test ./...

coverage:
	GOWORK=off go test ./... -coverprofile=coverage.out

build:
	go build -o ssh-watcher ./cmd/ssh-watcher
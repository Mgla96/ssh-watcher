.PHONY: lint test

lint:
	golangci-lint run

test:
	GOWORK=off go test ./...

coverage:
	GOWORK=off go test ./... -coverprofile=coverage.out

.PHONY: lint test build gen-docs

lint:
	golangci-lint run

test:
	GOWORK=off go test ./...

coverage:
	GOWORK=off go test ./... -coverprofile=coverage.out

build:
	go build -o ssh-watcher ./cmd/ssh-watcher

gen-docs:
	for d in $(shell find $(CURDIR)/internal -type f -name '*.go' | xargs -n 1 dirname | sort -u); \
	do \
	  cd $$d; \
  	  echo generating $$d/README.md; \
  	  gomarkdoc > README.md; \
  	  cd $(CURDIR); \
	done

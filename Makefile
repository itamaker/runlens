BINARY := runlens

.PHONY: build test example release-check snapshot

build:
	mkdir -p dist
	go build -o dist/$(BINARY) .

test:
	go test ./...

example:
	go run . summary -input examples/run.jsonl

release-check:
	goreleaser check

snapshot:
	goreleaser release --snapshot --clean


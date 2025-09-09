./bin/viki: $(shell find . -name '*.go')
	go build -o ./bin/viki ./cmd/viki

.PHONY: test
test:
	@go test ./...
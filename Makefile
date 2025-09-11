./bin/viki: $(shell find . -name '*.go') $(shell find . -name '*.tpl')
	go build -o ./bin/viki ./cmd/viki

.PHONY: test
test:
	@go test ./...
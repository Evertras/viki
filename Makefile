./bin/viki: $(shell find . -name '*.go') ./lib/viki/assets.go
	go build -o ./bin/viki ./cmd/viki

./lib/viki/assets.go: $(shell find ./lib/viki/static -type f) $(shell find lib/viki/templates -type f) ./lib/viki/gen-assets.sh
	cd lib/viki && ./gen-assets.sh

.PHONY: test
test: ./lib/viki/assets.go
	@go test ./...
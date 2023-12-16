export GO111MODULE=on

GO_FILES:=$(shell find . -type f -name '*.go' -print)

bin/credstash: $(GO_FILES)
	@go build -o $@ github.com/iyuuya/go-credstash/cmd/credstash

.PHONY: run
run:
	@go run github.com/iyuuya/go-credstash/cmd/credstash

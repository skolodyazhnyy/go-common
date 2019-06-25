.PHONY: dep lint test cover emojify

GO := go
PACKAGE_DIRS := $(shell $(GO) list ./... | grep -v /vendor/)
COVER_MSG := coverage.msg

dep:
	go get ./...

lint:
	golangci-lint run

test:
	go test ./...

cover:
	$(GO) test $(PACKAGE_DIRS) --cover > $(COVER_MSG)

emojify:
	./.drone/coverage_emojify $(COVER_MSG)

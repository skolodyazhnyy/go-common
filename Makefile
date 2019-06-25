.PHONY: dep lint test coverage

GO := go
NAME := go-common
OS := $(shell uname)
PACKAGE_DIRS := $(shell $(GO) list ./... | grep -v /vendor/)
PKGS := $(shell go list ./... | grep -v /vendor)
PKGS := $(subst  :,_,$(PKGS))
COVER_RAW := coverage.raw
COVER_MSG := coverage.msg

dep:
	go get ./...

lint:
	golangci-lint run

test:
	go test ./...

cover:
	$(GO) test $(PACKAGE_DIRS) --cover >> $(COVER_RAW)

emojify:
	echo '```' > $(COVER_MSG)
	cat $(COVER_RAW) >> $(COVER_MSG)
	echo '```' >> $(COVER_MSG)

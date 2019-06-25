.PHONY: dep lint test coverage

GO := go
NAME := go-common
OS := $(shell uname)
PACKAGE_DIRS := $(shell $(GO) list ./... | grep -v /vendor/)
PKGS := $(shell go list ./... | grep -v /vendor)
PKGS := $(subst  :,_,$(PKGS))
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
	sed -i 's/ok/:white_check_mark:/g' $(COVER_MSG)
	sed -i 's/?/:broken_heart:/g' $(COVER_MSG)
	sed -i 's/FAIL/FAIL :tomato:/g' $(COVER_MSG)

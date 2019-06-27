.PHONY: dep lint test cover emojify

dep:
	go get ./...

lint:
	golangci-lint run

test:
	go test ./...

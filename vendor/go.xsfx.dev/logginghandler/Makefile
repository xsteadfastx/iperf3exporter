.PHONY: test lint

test:
	go test -v ./...

lint:
	golangci-lint run --enable-all --disable=godox --timeout 10m

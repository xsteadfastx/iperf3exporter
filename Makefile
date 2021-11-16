.PHONY: build
build:
	goreleaser build --rm-dist --snapshot

.PHONY: release
release:
	goreleaser release --rm-dist --snapshot --skip-publish

.PHONY: generate
generate:
	go generate

.PHONY: lint
lint:
	golangci-lint run \
		--enable-all \
		--disable=exhaustivestruct,godox,varnamelen

.PHONY: test
test:
	go test -v -race -cover -coverprofile=coverage.out

.PHONY: coverage
coverage: test
	go tool cover -html=coverage.out

.PHONY: tidy
tidy:
	go mod tidy
	go mod vendor

.PHONY: install-tools
install-tools:
	go list -f '{{range .Imports}}{{.}} {{end}}' tools/tools.go | xargs go install -v

.PHONY: install-goreleaser
install-goreleaser:
	go install -v github.com/goreleaser/goreleaser

.PHONY: install-golangci-lint
install-golangci-lint:
	go install -v github.com/golangci/golangci-lint/cmd/golangci-lint

.PHONY: run-prometheus
run-prometheus:
	docker run --rm -ti \
		-v $(PWD)/test/prometheus.yml:/etc/prometheus/prometheus.yml \
		-p 9090:9090 \
		prom/prometheus

.PHONY: build
build:
	goreleaser build --clean --snapshot

.PHONY: test
test:
	go test -v ./...

.PHONY: dependencies
dependencies:
	go mod tidy

.PHONY: install-tools
install-tools: dependencies
	cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

.PHONY: lint
lint: lint-go lint-proto

.PHONY: lint-go
lint-go: dependencies
	golangci-lint run ./... --fix

.PHONY: lint-proto
lint-proto: dependencies
	buf lint
	buf format -w
	buf breaking --against ".git#branch=main"

.PHONY: gen
gen:
	buf generate

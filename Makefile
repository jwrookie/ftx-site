.PHONY: install-tools
install-tools: ## Install necessary tools
	@bash -c 'go install github.com/golang/mock/mockgen@v1.6.0'
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.46.2
	@go install golang.org/x/tools/cmd/goimports@latest

.PHONY: codegen
codegen: ## Run code generation
	./scripts/mockgen.sh

.PHONY: gofmt
gofmt: ## Format the source code
	@find . -type f -name "*.go" | xargs gofmt -w -s

.PHONY: goimports
goimports: ## Format the source code
	@find . -type f -name "*.go" | xargs goimports -w -v

.PHONY: lint
lint: ## Apply go lint check
	@golangci-lint run --timeout 10m ./...

.PHONY: test
test: ## Run the unit tests
	@UNIT_TEST=true go test -v ./... -count 1 -failfast -short

.PHONY: precommit
precommit: gofmt goimports lint test
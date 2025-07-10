ALLDOCS=$(shell find . -name '*.md' -type f | sort)
GOFILES=$(shell find . -type f -name '*.go' \
    -not -path "./vendor/*" \
    -not -path "./third_party/*" \
    -not -path "./.idea/*" \
    -not -name '*.pb.go' \
    -not -name '*mock*.go')

.PHONY:	fmt
fmt:
	@gofumpt -l -w $(GOFILES)
	@goimports -l -w $(GOFILES)

.PHONY: check-license
check-license:
	@echo "Ruining check-license"
	@missing=$$(find . -type f -name "*.go" \
	  -not -path "./vendor/*" \
	  -not -path "./third_party/*" \
	  -not -path "./.idea/*" \
	  -not -name '*.pb.go' \
	  -not -name '*mock*.go' \
	  | xargs grep -L "Licensed under the Apache License"); \
	if [ -n "$$missing" ]; then \
	  echo "The following files are missing the license header:"; \
	  echo "$$missing"; \
	  exit 1; \
	else \
	  echo "All Go files contain the license header."; \
	fi

.PHONY:	lint
lint:
	@echo "lint code..."
	@golangci-lint run -c .golangci.yaml
	@echo "finished the lint"

.PHONY: test
test:
	@echo "Running tests..."
	@go test -v -race -timeout=30s -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out
	@rm -f coverage.out

.PHONY: check
check:
	@make check-license
	@make fmt
	@make lint
	@make test
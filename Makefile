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

.PHONY:	lint
lint:
	@golangci-lint run -c .golangci.yaml

.PHONY: test
test:
	@go test -v -race -timeout=30s -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out
	@rm -f coverage.out

.PHONY: check
check:
	@make fmt
	@make lint
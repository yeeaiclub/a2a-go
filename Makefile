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

.PHONY: fmt-fix
fmt-fix:

.PHONY:	lint
lint:
	@golangci-lint run -c .golangci.yaml
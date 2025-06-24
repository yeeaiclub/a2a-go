.PHONY:	fmt
fmt:
	@goimports -l -w $$(find . -type f -name '*.go'  -not -path "./.idea/*" -not -name '*.pb.go' -not -name '*mock*.go')
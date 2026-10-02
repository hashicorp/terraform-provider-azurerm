TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PKG_NAME=tfe

fmt:
	gofmt -w $(GOFMT_FILES)

lint:
	@golangci-lint run ; if [ $$? -ne 0 ]; then \
		echo ""; \
		echo "golangci-lint found some code style issues."; \
		exit 1; \
	fi

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

.PHONY: fmt lint fmtcheck


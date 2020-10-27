TEST?=$$(go list ./... |grep -v 'vendor'|grep -v 'examples')
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=azurerm
TESTTIMEOUT=180m

.EXPORT_ALL_VARIABLES:
  TF_SCHEMA_PANIC_ON_ERROR=1
  GO111MODULE=on
  GOFLAGS=-mod=vendor

default: build

tools:
	@echo "==> installing required tooling..."
	@sh "$(CURDIR)/scripts/gogetcookie.sh"
	GO111MODULE=off go get -u github.com/client9/misspell/cmd/misspell
	GO111MODULE=off go get -u github.com/bflad/tfproviderlint/cmd/tfproviderlint
	GO111MODULE=off go get -u github.com/bflad/tfproviderdocs
	GO111MODULE=off go get -u github.com/katbyte/terrafmt
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH || $$GOPATH)/bin v1.24.0

build: fmtcheck generate
	go install

build-docker:
	mkdir -p bin
	docker run --rm -v $$(pwd)/bin:/go/bin -v $$(pwd):/go/src/github.com/terraform-providers/terraform-provider-azurerm -w /go/src/github.com/terraform-providers/terraform-provider-azurerm -e GOOS golang:1.13 make build

fmt:
	@echo "==> Fixing source code with gofmt..."
	# This logic should match the search logic in scripts/gofmtcheck.sh
	find . -name '*.go' | grep -v vendor | xargs gofmt -s -w

# Currently required by tf-deploy compile, duplicated by linters
fmtcheck:
	@sh "$(CURDIR)/scripts/gofmtcheck.sh"
	@sh "$(CURDIR)/scripts/timeouts.sh"
	@sh "$(CURDIR)/scripts/check-test-package.sh"

terrafmt:
	@echo "==> Fixing acceptance test terraform blocks code with terrafmt..."
	@find azurerm | egrep "_test.go" | sort | while read f; do terrafmt fmt -f $$f; done
	@echo "==> Fixing website terraform blocks code with terrafmt..."
	@find . | egrep html.markdown | sort | while read f; do terrafmt fmt $$f; done

generate:
	go generate ./azurerm/internal/provider/

goimports:
	@echo "==> Fixing imports code with goimports..."
	goimports -w $(PKG_NAME)/

lint:
	./scripts/run-lint.sh

# we have split off static check because it causes travis to fail with an OOM error
lintunused:
	@echo "==> Checking source code against static check linters..."
	(while true; do sleep 300; echo "(I'm still alive and linting!)"; done) & PID=$$!; echo $$PID; \
	golangci-lint run ./... -v --no-config --concurrency 1 --deadline=30m10s --disable-all --enable=unused; ES=$$?; kill -9 $$PID; exit $$ES

lintrest:
	./scripts/run-lint-rest.sh

depscheck:
	@echo "==> Checking source code with go mod tidy..."
	@go mod tidy
	@git diff --exit-code -- go.mod go.sum || \
		(echo; echo "Unexpected difference in go.mod/go.sum files. Run 'go mod tidy' command or revert any go.mod/go.sum changes and commit."; exit 1)
	@echo "==> Checking source code with go mod vendor..."
	@go mod vendor
	@git diff --compact-summary --exit-code -- vendor || \
		(echo; echo "Unexpected difference in vendor/ directory. Run 'go mod vendor' command or revert any go.mod/go.sum/vendor changes and commit."; exit 1)

tflint:
	./scripts/run-tflint.sh

whitespace:
	@echo "==> Fixing source code with whitespace linter..."
	golangci-lint run ./... --no-config --disable-all --enable=whitespace --fix

test-docker:
	docker run --rm -v $$(pwd):/go/src/github.com/terraform-providers/terraform-provider-azurerm -w /go/src/github.com/terraform-providers/terraform-provider-azurerm golang:1.13 make test

test: fmtcheck
	@TEST=$(TEST) ./scripts/run-gradually-deprecated.sh
	@TEST=$(TEST) ./scripts/run-test.sh

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout $(TESTTIMEOUT) -ldflags="-X=github.com/terraform-providers/terraform-provider-azurerm/version.ProviderVersion=acc"

acctests: fmtcheck
	TF_ACC=1 go test -v ./azurerm/internal/services/$(SERVICE)/tests/ $(TESTARGS) -timeout $(TESTTIMEOUT) -ldflags="-X=github.com/terraform-providers/terraform-provider-azurerm/version.ProviderVersion=acc"

debugacc: fmtcheck
	TF_ACC=1 dlv test $(TEST) --headless --listen=:2345 --api-version=2 -- -test.v $(TESTARGS)

website-lint:
	@echo "==> Checking documentation spelling..."
	@misspell -error -source=text -i hdinsight,exportfs website/
	@echo "==> Checking documentation for errors..."
	@tfproviderdocs check -provider-name=azurerm -require-resource-subcategory \
		-allowed-resource-subcategories-file website/allowed-subcategories
	@sh -c "'$(CURDIR)/scripts/terrafmt-website.sh'"

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

scaffold-website:
	./scripts/scaffold-website.sh

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.PHONY: build build-docker test test-docker testacc vet fmt fmtcheck errcheck scaffold-website test-compile website website-test

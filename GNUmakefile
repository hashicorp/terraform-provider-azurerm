TEST?=$$(go list ./... |grep -v 'vendor'|grep -v 'examples')
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=azurerm

#make sure we catch schema errors during testing
TF_SCHEMA_PANIC_ON_ERROR=1
GO111MODULE=on
GOFLAGS=-mod=vendor

default: build

tools:
	@echo "==> installing required tooling..."
	@sh "$(CURDIR)/scripts/gogetcookie.sh"
	GO111MODULE=off go get -u github.com/client9/misspell/cmd/misspell
	GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	GO111MODULE=off go get -u github.com/bflad/tfproviderlint/cmd/tfproviderlint

build: fmtcheck
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

goimports:
	@echo "==> Fixing imports code with goimports..."
	goimports -w $(PKG_NAME)/

lint:
	@echo "==> Checking source code against linters..."
	golangci-lint run ./...

# we have split off static check because it causes travis to fail with an OOM error
lintstatic:
	@echo "==> Checking source code against static check linters..."
	(while true; do sleep 300; echo "(I'm still alive and linting!)"; done) & PID=$$!; echo $$PID; \
	golangci-lint run ./... -v --no-config --concurrency 1 --deadline=30m10s --disable-all --enable=staticcheck; ES=$$?; kill -9 $$PID; exit $$ES

lintrest:
	@echo "==> Checking source code against linters..."
	(while true; do sleep 300; echo "(I'm still alive and linting!)"; done) & PID=$$!; echo $$PID; \
	golangci-lint run ./... -v --concurrency 1 --deadline=30m10s --config .golangci-travisrest.yml ; ES=$$?; kill -9 $$PID; exit $$ES

tflint:
	@echo "==> Checking source code against terraform provider linters..."
	@tfproviderlint \
        -R001 -R002 -R003 -R004\
        -S001 -S002 -S003 -S004 -S005 -S006 -S007 -S008 -S009 -S010 -S011 -S012 -S013 -S014 -S015 -S016 -S017 -S018 -S019\
        ./$(PKG_NAME)

whitespace:
	@echo "==> Fixing source code with whitespace linter..."
	golangci-lint run ./... --no-config --disable-all --enable=whitespace --fix

test-docker:
	docker run --rm -v $$(pwd):/go/src/github.com/terraform-providers/terraform-provider-azurerm -w /go/src/github.com/terraform-providers/terraform-provider-azurerm golang:1.13 make test

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 180m -ldflags="-X=github.com/terraform-providers/terraform-provider-azurerm/version.ProviderVersion=acc"

debugacc: fmtcheck
	TF_ACC=1 dlv test $(TEST) --headless --listen=:2345 --api-version=2 -- -test.v $(TESTARGS)

website-lint:
	@echo "==> Checking website against linters..."
	@misspell -error -source=text -i hdinsight website/

website-registrycheck:
	@sh "$(CURDIR)/scripts/website-registrycheck.sh"

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.PHONY: build build-docker test test-docker testacc vet fmt fmtcheck errcheck test-compile website website-test

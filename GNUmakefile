TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

ifndef TF_PARL_NUM
	TF_PARL_NUM=20
endif

default: build

build: fmtcheck
	go install

test-install:
	go test -i $(TEST) || exit 1

test: fmtcheck test-install
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

# example usages:
#  1. Enable parallel run for all acceptance tests with default number(20) of
#  threads:
#    TF_PARL=1 make testacc
#  2. Enable parallel run for acceptance tests with matching pattern.
#    TF_PARL=1 TESTARGS="-run TestAccAzureRMStorage" make testacc
testacc: fmtcheck test-install
	echo $(TEST) | \
	TF_ACC=1 xargs -t -n1 go test -v $(TESTARGS) -timeout 300m -parallel=$(TF_PARL_NUM) 2>&1

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

vendor-status:
	@govendor status

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./aws"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

.PHONY: build test testacc vet fmt fmtcheck errcheck vendor-status test-compile


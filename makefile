TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
NAME?=$$(basename "$$(pwd)")
VERSION:=$(if $(VERSION),$(VERSION),$$(git describe --abbrev=0 --tags)+$$(git rev-parse --short=8 HEAD))

default: build

build: fmtcheck
	gox -verbose \
		-ldflags "-X main.version=$(VERSION)" \
		-os "linux darwin" \
		-arch "amd64" \
		-output "dist/$(NAME)-$(VERSION)-{{.OS}}_{{.Arch}}/$(NAME)_$(VERSION)"

dist:
	cd dist &&\
	find * -type d -exec tar --mtime 1970-01-01 -zcf {}.tar.gz {} \;

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 180m

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
	@sh "$(CURDIR)/scripts/gofmtcheck.sh"

errcheck:
	@sh "$(CURDIR)/scripts/errcheck.sh"

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

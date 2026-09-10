TEST?=$$(go list ./... |grep -v 'vendor'|grep -v 'examples')
TESTTIMEOUT=180m
TF_SCHEMA_PANIC_ON_ERROR=1

# The single source of truth for the golangci-lint version is the 'version:' field in
# scripts/.custom-gcl.yml (it is required to live there for the plugin build); everything
# else, including the CI workflows, derives it from that file.
GOLANGCI_LINT_VERSION := $(shell sed -n 's/^version: *//p' scripts/.custom-gcl.yml)

# The single source of truth for the actionlint version is the go install pin
# in .github/workflows/workflow-actionlint.yml.
ACTIONLINT_VERSION := $(shell sed -n 's/.*actionlint\/cmd\/actionlint@//p' .github/workflows/workflow-actionlint.yml)


.EXPORT_ALL_VARIABLES:

default: build

help: ## Show this help
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-18s\033[0m%s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) }' $(MAKEFILE_LIST)

##@ Deprecated (remove at the end of 2026)
fmtcheck: ## renamed to quick-checks
	@echo "NOTE: 'make fmtcheck' has been renamed to 'make quick-checks' to reflect what it actually runs and will be removed in the future."
	@$(MAKE) quick-checks

tflint: ## renamed to tfproviderlint
	@echo "NOTE: 'make tflint' has been renamed to 'make tfproviderlint' to reflect what it actually runs and will be removed in the future."
	@$(MAKE) tfproviderlint

golangci-fix: ## renamed to lint-fix
	@echo "NOTE: 'make golangci-fix' has been renamed to 'make lint-fix' and will be removed in the future."
	@$(MAKE) lint-fix

##@ Build & Generate
tools: ## Install the tools required to develop the provider
	@echo "==> installing required tooling..."
	go install github.com/client9/misspell/cmd/misspell@latest
	go install github.com/YakDriver/tfproviderdocs@latest
	go install github.com/katbyte/terrafmt@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install mvdan.cc/gofumpt@latest
	go install github.com/rhysd/actionlint/cmd/actionlint@$(ACTIONLINT_VERSION)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $$(go env GOPATH || $$GOPATH)/bin $(GOLANGCI_LINT_VERSION)
	@$(MAKE) golangci-with-modules

build: quick-checks generate ## Run the quick checks, generate code, and compile the provider
	go install

generate: tools ## Regenerate auto-generated code
	go generate ./internal/services/...
	go generate ./internal/provider/

gencheck: generate ## Check that generated code matches what is committed
	@echo "==> Comparing generated code to committed code..."
	@git diff --compact-summary --exit-code -- ./ || \
		(echo; echo "Unexpected difference in generated code. Run 'make generate' to update the generated code and commit."; echo "If you added or modified a resource, ensure 'go generate' directives are up to date."; exit 1)

##@ Formatting & Quick Checks
# All top-level locations containing Go source, excluding vendor.
GOPATHS=main.go helpers internal version

# The fixers here (plus goimports below) should match the checks in scripts/checks/fmt-check.sh
fmt: ## Fix Go formatting (gofmt, gofumpt, whitespace)
	@echo "==> Fixing source code with gofmt..."
	@gofmt -s -w $(GOPATHS)
	@echo "==> Fixing source code with gofumpt..."
	@gofumpt -w $(GOPATHS)
	@echo "==> Fixing source code with whitespace linter..."
	@golangci-lint run ./... --no-config --enable-only=whitespace --fix

# goimports runs via `golangci-lint fmt` as the standalone binary is single-threaded and far slower
goimports: ## Fix Go import ordering/grouping (slower than fmt, so kept separate)
	@echo "==> Fixing imports with goimports and gci..."
	@golangci-lint fmt -E goimports,gci

quick-checks: ## Run the quick CI checks (formatting)
	@echo "==> Running the set of quick CI checks (formatting)..."
	@sh "$(CURDIR)/scripts/checks/fmt-check.sh"
	@sh "$(CURDIR)/scripts/checks/terrafmt-acctests.sh"

terrafmt: ## Fix terraform blocks in acceptance tests and website docs
	@echo "==> Fixing acceptance test terraform blocks code with terrafmt..."
	@terrafmt fmt -f -p "*_test.go" ./internal
	@echo "==> Fixing website terraform blocks code with terrafmt..."
	@terrafmt fmt -p "*.html.markdown" .

##@ Linting & Dependencies
# golangci-lint module plugins (azproviderlint) only exist in a custom-built binary, so lint
# targets use scripts/golangci-with-modules, rebuilt automatically whenever the config (which
# pins the golangci-lint version) changes. The pinned version is installed before building so
# the host binary running 'golangci-lint custom' always matches the pin. The config filename
# and its living in the build cwd are both fixed by golangci-lint, hence the cd into scripts/.
scripts/golangci-with-modules: scripts/.custom-gcl.yml
	@echo "==> Building golangci-lint with plugins (scripts/golangci-with-modules)..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $$(go env GOPATH)/bin $(GOLANGCI_LINT_VERSION)
	@cd scripts && $$(go env GOPATH)/bin/golangci-lint custom

golangci-with-modules: ## Build golangci-lint with plugins (automatic when the config or pinned version changes)
	@if [ -x scripts/golangci-with-modules ] && ! ./scripts/golangci-with-modules version 2>/dev/null | grep -qF -- "$(GOLANGCI_LINT_VERSION:v%=%)"; then \
		echo "==> scripts/golangci-with-modules is not $(GOLANGCI_LINT_VERSION), rebuilding..."; \
		rm -f scripts/golangci-with-modules; \
	fi
	@$(MAKE) scripts/golangci-with-modules

lint: golangci-with-modules ## Check source code with the golangci linters
	@echo "==> Checking source code with golangci-lint..."
	@./scripts/golangci-with-modules run -v ./...

lint-fix: golangci-with-modules ## Fix source code with all golangci linters
	@echo "==> Fixing source code with all golangci linters..."
	@./scripts/golangci-with-modules run ./... --fix

# tfproviderlint and azproviderlint run as part of lint; these targets run just their checks
tfproviderlint: golangci-with-modules ## Check terraform schema definitions with only the tfproviderlint checks
	@echo "==> Checking terraform schemas with tfproviderlint (via golangci-lint)..."
	@./scripts/golangci-with-modules run -v --enable-only tfproviderlint ./...

azproviderlint: golangci-with-modules ## Check source code with only the azproviderlint checks
	@echo "==> Checking source code with azproviderlint (via golangci-lint)..."
	@./scripts/golangci-with-modules run -v --enable-only azproviderlint ./...

yamllint: ## Check YAML files with yamllint (config in .yamllint.yml)
	@command -v yamllint >/dev/null || (echo "yamllint not installed. Install via: brew install yamllint (macOS) or pip install yamllint" && exit 1)
	@echo "==> Checking YAML files with yamllint..."
	@yamllint -s .

actionlint: ## Check GitHub workflows with actionlint (incl. shellcheck on run blocks)
	@command -v actionlint >/dev/null || (echo "actionlint not installed. Install via 'make tools' or: go install github.com/rhysd/actionlint/cmd/actionlint@$(ACTIONLINT_VERSION)" && exit 1)
	@echo "==> Checking workflows with actionlint..."
	@actionlint

shellcheck: ## Check shell scripts with shellcheck
	@command -v shellcheck >/dev/null || (echo "shellcheck not installed. Install via: brew install shellcheck (macOS) or apt install shellcheck (Linux)" && exit 1)
	@echo "==> Checking shell scripts with shellcheck..."
	@shellcheck scripts/*.sh scripts/checks/*.sh scripts/automation/*.sh || \
		(echo; echo "ShellCheck found issues in shell scripts."; echo "Review the errors above and fix them. See https://www.shellcheck.net/ for detailed explanations of each rule."; exit 1)

depscheck: ## Check that go.mod/go.sum and vendor/ are in sync
	@echo "==> Checking dependencies.."
	@./scripts/checks/track2-check.sh
	@echo "==> Checking source code with go mod tidy..."
	@go mod tidy
	@git diff --exit-code -- go.mod go.sum || \
		(echo; echo "Unexpected difference in go.mod/go.sum files. Run 'go mod tidy' command or revert any go.mod/go.sum changes and commit."; echo "Do not modify files in the vendor/ directory directly."; exit 1)
	@echo "==> Checking source code with go mod vendor..."
	@go mod vendor
	@git diff --compact-summary --exit-code -- vendor || \
		(echo; echo "Unexpected difference in vendor/ directory. Run 'go mod vendor' command or revert any go.mod/go.sum/vendor changes and commit."; echo "Do not modify files in the vendor/ directory directly."; exit 1)

##@ Testing
test: ## Run the unit tests
	@TEST=$(TEST) ./scripts/checks/gradually-deprecated.sh
	@TEST=$(TEST) ./scripts/checks/test.sh

testacc: ## Run acceptance tests for a package (TEST=./internal/services/<service>)
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout $(TESTTIMEOUT) -ldflags="-X=github.com/hashicorp/terraform-provider-azurerm/version.ProviderVersion=acc"

acctests: ## Run acceptance tests for a service (SERVICE=<service>)
	TF_ACC=1 go test -v ./internal/services/$(SERVICE) $(TESTARGS) -timeout $(TESTTIMEOUT) -ldflags="-X=github.com/hashicorp/terraform-provider-azurerm/version.ProviderVersion=acc"

debugacc: ## Run acceptance tests under the delve debugger (TEST=./internal/services/<service>)
	TF_ACC=1 dlv test $(TEST) --headless --listen=:2345 --api-version=2 -- -test.v $(TESTARGS)

##@ Generated Code
prepare: ## Remove all generated files ahead of a full regeneration
	@echo "==> Preparing the repository (removing all '*_gen.go' files)..."
	@find . -iname \*_gen.go -type f -delete
	@echo "==> Preparing the repository (removing all '*_gen_test.go' files)..."
	@find . -iname \*_gen_test.go -type f -delete

##@ Website & Documentation
# markdown checked by markdownlint: website docs, README, contributing docs, and the
# .github markdown (PR/issue templates etc). The resource/data-source docs are exempt
# (leading \# ignore glob, \# escapes the hash from make) as they will soon be generated.
MARKDOWN_INPUTS='website/docs/**/*.markdown' README.md 'contributing/**/*.md' '.github/**/*.md' '\#website/docs/r' '\#website/docs/d'

markdownlint: ## Check repo markdown with markdownlint (config in .markdownlint.yml)
	@command -v markdownlint-cli2 >/dev/null || (echo "markdownlint-cli2 not installed. Install via: brew install markdownlint-cli2 (macOS) or npm install -g markdownlint-cli2" && exit 1)
	@echo "==> Checking markdown with markdownlint..."
	@markdownlint-cli2 $(MARKDOWN_INPUTS)

website-lint: ## Check website documentation for issues
	@echo "==> Checking documentation for .html.markdown extension present"
	@if ! find website/docs -type f -not -name "*.html.markdown" -print -exec false {} +; then \
		echo "ERROR: file extension should be .html.markdown"; \
		echo "All documentation files must use the .html.markdown extension."; \
		exit 1; \
	fi
	@echo "==> Checking documentation spelling..."
	@misspell -error -source=text -i hdinsight,exportfs website/ || \
		(echo; echo "Spelling errors found in documentation. Install misspell: go install github.com/client9/misspell/cmd/misspell@latest"; exit 1)
	@echo "==> Checking documentation for errors..."
	@tfproviderdocs check -provider-name=azurerm -require-resource-subcategory \
		-allowed-resource-subcategories-file website/allowed-subcategories || \
		(echo; echo "Documentation validation failed. Check that your docs follow the provider documentation format."; echo "See: contributing/topics/guide-new-resource.md for documentation requirements."; exit 1)
	@sh -c "'$(CURDIR)/scripts/checks/terrafmt-website.sh'"

document-validate: ## Check website documentation against resource schemas
	@./scripts/checks/document-validate.sh

document-fix: ## Fix website documentation issues against resource schemas
	@./scripts/checks/document-fix.sh

document-lint: ## Check website documentation with document-lint
	@echo "==> Checking documentation with document-lint..."
	@go run $(CURDIR)/internal/tools/document-lint/main.go check

scaffold-website: ## Scaffold website documentation for a new resource or data source
	@./scripts/website-scaffold.sh

##@ Other
teamcity-test: ## Test the TeamCity configuration
	@$(MAKE) -C .teamcity tools
	@$(MAKE) -C .teamcity test

validate-examples: build ## Check that the terraform examples are valid
	@echo "==> Validating examples..."
	@./scripts/checks/examples-validate.sh

schemagen: ## Generate a schema snapshot (RESOURCE_TYPE=<resource>)
	@go run ./internal/tools/generator-schema-snapshot $(RESOURCE_TYPE)

resource-counts: ## Print the number of resources and data sources in the provider
	@go test -v ./internal/provider -run=TestProvider_counts

pr-check: generate build test lint website-lint ## Run the same set of checks CI runs against a PR

.PHONY: default help tools build fmt goimports quick-checks fmtcheck terrafmt generate lint actionlint yamllint markdownlint shellcheck depscheck gencheck tfproviderlint tflint azproviderlint lint-fix golangci-fix test testacc acctests debugacc prepare website-lint document-validate document-fix document-lint scaffold-website teamcity-test validate-examples schemagen resource-counts pr-check

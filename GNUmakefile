TEST?=$$(go list ./... |grep -v 'vendor'|grep -v 'examples')
TESTTIMEOUT=180m
TF_SCHEMA_PANIC_ON_ERROR=1

# Go dev tools are built into .tools/bin (gitignored) from the versions pinned in .tools/go.mod,
# the single source of truth for make and CI (dependabot keeps them bumped). Each target lists
# the tool binaries it needs as prerequisites, so the first run builds them and nothing has to
# be go installed by hand. .tools/bin is also prepended to PATH so the check scripts and
# go generate (which shell out to goimports/gofumpt/terrafmt by name) use the pinned builds.
TOOLS_BIN=.tools/bin
ACTIONLINT=$(TOOLS_BIN)/actionlint
GOFUMPT=$(TOOLS_BIN)/gofumpt
GOIMPORTS=$(TOOLS_BIN)/goimports
GOLANGCI_LINT=$(TOOLS_BIN)/golangci-lint
GOTESTSUM=$(TOOLS_BIN)/gotestsum
MISSPELL=$(TOOLS_BIN)/misspell
TCTEST=$(TOOLS_BIN)/tctest
TERRAFMT=$(TOOLS_BIN)/terrafmt
TFPROVIDERDOCS=$(TOOLS_BIN)/tfproviderdocs
PATH := $(CURDIR)/$(TOOLS_BIN):$(PATH)

# non-Go tools also live in .tools/bin at pinned versions, but the pins are here (dependabot
# cannot bump them): shellcheck is a static binary downloaded from its github releases, yamllint
# is pip installed into a repo-local venv and markdownlint-cli2 is npm installed into a repo-local
# prefix. all rebuild when this makefile changes.
MARKDOWNLINT_CLI2_VERSION=0.23.2
SHELLCHECK_VERSION=v0.11.0
YAMLLINT_VERSION=1.38.0
MARKDOWNLINT=$(TOOLS_BIN)/markdownlint-cli2
SHELLCHECK=$(TOOLS_BIN)/shellcheck
YAMLLINT=$(TOOLS_BIN)/yamllint

# golangci-lint with the azproviderlint/tfproviderlint module plugins compiled in
# (.tools/.custom-gcl.yml); the lint targets use this binary, the plain one bootstraps
# `golangci-lint custom` and runs the formatters
GOLANGCI_LINT_MODULES=$(TOOLS_BIN)/golangci-with-modules

.EXPORT_ALL_VARIABLES:

# one rule builds any Go tool: the import path comes from the tool directives in .tools/go.mod
# (via go list tool), so the makefile never repeats it - add a tool there and a variable above
$(TOOLS_BIN)/%: .tools/go.mod .tools/go.sum
	@echo "==> Building $* (version pinned in .tools/go.mod)..."
	@cd .tools && go build -o bin/$* $$(go list tool | grep "/$*$$")

# explicit rules take precedence over the pattern rule above; golangci-lint custom must run from
# the directory holding .custom-gcl.yml, hence the cd
$(GOLANGCI_LINT_MODULES): .tools/.custom-gcl.yml $(GOLANGCI_LINT)
	@echo "==> Building golangci-lint with plugins (versions pinned in .tools/.custom-gcl.yml)..."
	@cd .tools && bin/golangci-lint custom

$(SHELLCHECK): GNUmakefile
	@echo "==> Downloading shellcheck $(SHELLCHECK_VERSION)..."
	@mkdir -p $(TOOLS_BIN)
	@os=$$(uname | tr 'A-Z' 'a-z'); arch=$$(uname -m); [ "$$arch" = "arm64" ] && arch=aarch64; \
		curl -sSfL "https://github.com/koalaman/shellcheck/releases/download/$(SHELLCHECK_VERSION)/shellcheck-$(SHELLCHECK_VERSION).$$os.$$arch.tar.xz" \
		| tar -xJ -O shellcheck-$(SHELLCHECK_VERSION)/shellcheck > $@ && chmod +x $@

$(YAMLLINT): GNUmakefile
	@command -v python3 >/dev/null || (echo "python3 is required to install yamllint (macOS: xcode CLT; Debian/Ubuntu: apt install python3-venv)" && exit 1)
	@echo "==> Installing yamllint $(YAMLLINT_VERSION) into .tools/venv..."
	@mkdir -p $(TOOLS_BIN)
	@python3 -m venv .tools/venv && .tools/venv/bin/pip install -q yamllint==$(YAMLLINT_VERSION) && ln -sf ../venv/bin/yamllint $@

$(MARKDOWNLINT): GNUmakefile
	@command -v npm >/dev/null || (echo "npm is required to install markdownlint-cli2 (macOS: brew install node; Debian/Ubuntu: apt install npm)" && exit 1)
	@echo "==> Installing markdownlint-cli2 $(MARKDOWNLINT_CLI2_VERSION) into .tools/npm..."
	@mkdir -p $(TOOLS_BIN) .tools/npm
	@npm install --silent --no-audit --no-fund --prefix .tools/npm markdownlint-cli2@$(MARKDOWNLINT_CLI2_VERSION) && ln -sf ../npm/node_modules/.bin/markdownlint-cli2 $@

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
tools: $(ACTIONLINT) $(GOFUMPT) $(GOIMPORTS) $(GOLANGCI_LINT) $(GOLANGCI_LINT_MODULES) $(GOTESTSUM) $(MISSPELL) $(TCTEST) $(TERRAFMT) $(TFPROVIDERDOCS) $(MARKDOWNLINT) $(SHELLCHECK) $(YAMLLINT) ## Install all pinned dev tools into .tools/bin (targets install what they need on demand)

build: quick-checks generate ## Run the quick checks, generate code, and compile the provider
	go install

generate: $(GOFUMPT) $(GOIMPORTS) $(TERRAFMT) ## Regenerate auto-generated code
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
fmt: $(GOFUMPT) $(GOLANGCI_LINT) ## Fix Go formatting (gofmt, gofumpt, whitespace)
	@echo "==> Fixing source code with gofmt..."
	@gofmt -s -w $(GOPATHS)
	@echo "==> Fixing source code with gofumpt..."
	@$(GOFUMPT) -w $(GOPATHS)
	@echo "==> Fixing source code with whitespace linter..."
	@$(GOLANGCI_LINT) run ./... --no-config --enable-only=whitespace --fix

# goimports runs via `golangci-lint fmt` as the standalone binary is single-threaded and far slower
goimports: $(GOLANGCI_LINT) ## Fix Go import ordering/grouping (slower than fmt, so kept separate)
	@echo "==> Fixing imports with goimports and gci..."
	@$(GOLANGCI_LINT) fmt -E goimports,gci

quick-checks: $(GOLANGCI_LINT) $(TERRAFMT) ## Run the quick CI checks (formatting)
	@echo "==> Running the set of quick CI checks (formatting)..."
	@sh "$(CURDIR)/scripts/checks/fmt-check.sh"
	@sh "$(CURDIR)/scripts/checks/terrafmt-acctests.sh"

terrafmt: $(TERRAFMT) ## Fix terraform blocks in acceptance tests and website docs
	@echo "==> Fixing acceptance test terraform blocks code with terrafmt..."
	@$(TERRAFMT) fmt -f -p "*_test.go" ./internal
	@echo "==> Fixing website terraform blocks code with terrafmt..."
	@$(TERRAFMT) fmt -p "*.html.markdown" .

##@ Linting & Dependencies
# golangci-lint module plugins (azproviderlint, tfproviderlint) only exist in a custom-built
# binary, so the lint targets use .tools/bin/golangci-with-modules, rebuilt automatically
# whenever .tools/.custom-gcl.yml or the golangci-lint pin in .tools/go.mod changes
golangci-with-modules: $(GOLANGCI_LINT_MODULES) ## Build golangci-lint with plugins into .tools/bin (automatic when the pins change)

lint: $(GOLANGCI_LINT_MODULES) ## Check source code with the golangci linters
	@echo "==> Checking source code with golangci-lint..."
	@$(GOLANGCI_LINT_MODULES) run -v ./...

lint-fix: $(GOLANGCI_LINT_MODULES) ## Fix source code with all golangci linters
	@echo "==> Fixing source code with all golangci linters..."
	@$(GOLANGCI_LINT_MODULES) run ./... --fix

# tfproviderlint and azproviderlint run as part of lint; these targets run just their checks
tfproviderlint: $(GOLANGCI_LINT_MODULES) ## Check terraform schema definitions with only the tfproviderlint checks
	@echo "==> Checking terraform schemas with tfproviderlint (via golangci-lint)..."
	@$(GOLANGCI_LINT_MODULES) run -v --enable-only tfproviderlint ./...

azproviderlint: $(GOLANGCI_LINT_MODULES) ## Check source code with only the azproviderlint checks
	@echo "==> Checking source code with azproviderlint (via golangci-lint)..."
	@$(GOLANGCI_LINT_MODULES) run -v --enable-only azproviderlint ./...

yamllint: $(YAMLLINT) ## Check YAML files with yamllint (config in .yamllint.yml)
	@echo "==> Checking YAML files with yamllint..."
	@$(YAMLLINT) -s .

actionlint: $(ACTIONLINT) $(SHELLCHECK) ## Check GitHub workflows with actionlint (incl. shellcheck on run blocks)
	@echo "==> Checking workflows with actionlint..."
	@$(ACTIONLINT) -shellcheck=$(SHELLCHECK)

shellcheck: $(SHELLCHECK) ## Check shell scripts with shellcheck
	@echo "==> Checking shell scripts with shellcheck..."
	@$(SHELLCHECK) scripts/*.sh scripts/checks/*.sh scripts/automation/*.sh || \
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
	@echo "==> Checking .tools/go.mod with go mod tidy..."
	@cd .tools && go mod tidy
	@git diff --exit-code -- .tools/go.mod .tools/go.sum || \
		(echo; echo "Unexpected difference in .tools/go.mod/go.sum. Run 'cd .tools && go mod tidy' and commit."; exit 1)
	@echo "==> Checking .tools/.custom-gcl.yml golangci-lint version matches .tools/go.mod..."
	@modv=$$(cd .tools && go list -m -f '{{.Version}}' github.com/golangci/golangci-lint/v2); \
		gclv=$$(sed -n 's/^version: *//p' .tools/.custom-gcl.yml); \
		[ "$$modv" = "$$gclv" ] || \
		(echo; echo "golangci-lint version mismatch: .tools/go.mod has $$modv but .tools/.custom-gcl.yml has $$gclv - update .custom-gcl.yml to match."; exit 1)

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

markdownlint: $(MARKDOWNLINT) ## Check repo markdown with markdownlint (config in .markdownlint.yml)
	@echo "==> Checking markdown with markdownlint..."
	@$(MARKDOWNLINT) $(MARKDOWN_INPUTS)

website-lint: $(MISSPELL) $(TFPROVIDERDOCS) $(TERRAFMT) ## Check website documentation for issues
	@echo "==> Checking documentation for .html.markdown extension present"
	@if ! find website/docs -type f -not -name "*.html.markdown" -print -exec false {} +; then \
		echo "ERROR: file extension should be .html.markdown"; \
		echo "All documentation files must use the .html.markdown extension."; \
		exit 1; \
	fi
	@echo "==> Checking documentation spelling..."
	@$(MISSPELL) -error -source=text -i hdinsight,exportfs website/ || \
		(echo; echo "Spelling errors found in documentation."; exit 1)
	@echo "==> Checking for locale-specific Microsoft Learn links..."
	@! grep -rnE '(learn|docs)\.microsoft\.com/[a-z]{2}-[a-z]{2}/' website/ || \
		(echo; echo "Remove the locale segment (e.g. /en-us/) from Microsoft Learn links so readers are served their own language."; exit 1)
	@echo "==> Checking documentation for errors..."
	@$(TFPROVIDERDOCS) check -provider-name=azurerm -require-resource-subcategory \
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

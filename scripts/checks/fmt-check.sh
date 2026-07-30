#!/usr/bin/env bash
# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0


if ! command -v golangci-lint >/dev/null; then
    echo "ERROR: golangci-lint is not installed. Run 'make tools' to install required tooling."
    exit 1
fi

# The checks here should match the fixers in the GNUmakefile `fmt` target.

# Check gofmt, gofumpt and goimports via `golangci-lint fmt`, configured in ../../.golangci.yml
echo "==> Checking that code complies with formatting requirements (gofmt, gofumpt, goimports)..."

if ! golangci-lint fmt --diff; then
    echo "You can use the command: \`make fmt\` to reformat code."
    exit 1
fi

# Check whitespace
echo "==> Checking that code complies with the whitespace linter..."

if ! golangci-lint run ./... --no-config --enable-only=whitespace; then
    echo "You can use the command: \`make fmt\` to reformat code."
    exit 1
fi

exit 0

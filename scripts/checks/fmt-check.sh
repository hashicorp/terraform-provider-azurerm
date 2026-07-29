#!/usr/bin/env bash
# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0


for tool in gofmt gofumpt goimports golangci-lint; do
    if ! command -v "$tool" >/dev/null; then
        echo "ERROR: $tool is not installed. Run 'make tools' to install required tooling."
        exit 1
    fi
done

# All top-level locations containing Go source, excluding vendor.
# This list should match the one in ../../GNUmakefile.
gopaths="main.go helpers internal utils version"

# Check gofmt
echo "==> Checking that code complies with gofmt requirements..."

gofmt_files=$(gofmt -s -l $gopaths)
if [ -n "${gofmt_files}" ]; then
    echo 'gofmt needs running on the following files:'
    echo "${gofmt_files}"
    echo "You can use the command: \`make fmt\` to reformat code."
    exit 1
fi

# Check gofumpt
echo "==> Checking that code complies with gofumpt requirements..."

gofumpt_files=$(gofumpt -l $gopaths)
if [ -n "${gofumpt_files}" ]; then
    echo 'gofumpt needs running on the following files:'
    echo "${gofumpt_files}"
    echo "You can use the command: \`make fmt\` to reformat code."
    exit 1
fi

# Check whitespace
echo "==> Checking that code complies with the whitespace linter..."

if ! golangci-lint run ./... --no-config --enable-only=whitespace; then
    echo "You can use the command: \`make fmt\` to reformat code."
    exit 1
fi

# Check goimports
echo "==> Checking that imports comply with goimports requirements..."

goimports_files=$(goimports -l $gopaths)
if [ -n "${goimports_files}" ]; then
    echo 'goimports needs running on the following files:'
    echo "${goimports_files}"
    echo "You can use the command: \`make fmt\` to reformat code."
    exit 1
fi

exit 0
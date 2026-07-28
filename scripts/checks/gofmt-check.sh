#!/usr/bin/env bash
# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0


# Check gofmt
echo "==> Checking that code complies with gofmt requirements..."

# This filter should match the search filter in ../../GNUmakefile
gofmt_files=$(find . -name '*.go' | grep -v vendor | xargs gofmt -s -l)
if [ -n "${gofmt_files}" ]; then
    echo 'gofmt needs running on the following files:'
    echo "${gofmt_files}"
    echo "You can use the command: \`make fmt\` to reformat code."
    exit 1
fi

# Check gofumpt
echo "==> Checking that code complies with gofumpt requirements..."

# This filter should match the search filter in ../../GNUmakefile
gofumpt_files=$(find . -name '*.go' | grep -v vendor | xargs gofumpt -l)
if [ -n "${gofumpt_files}" ]; then
    echo 'gofumpt needs running on the following files:'
    echo "${gofumpt_files}"
    echo "You can use the command: \`make fmt\` to reformat code."
    exit 1
fi

exit 0
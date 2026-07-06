#!/usr/bin/env bash
# Copyright IBM Corp. 2023, 2026
# SPDX-License-Identifier: MPL-2.0


function on_failure {
  echo ""
  echo "==> golangci-lint failed!"
  echo "    To auto-fix some issues run: make golangci-fix"
  echo "    Common issues: unused variables, formatting, error handling, naming conventions"
  echo "    Docs: https://golangci-lint.run/"
  echo ""
}

function runLinters {
  echo "==> Checking source code against linters..."
  golangci-lint run -v ./...
}

function main {
  runLinters || { on_failure; exit 1; }
}

main

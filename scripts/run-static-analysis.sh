#!/usr/bin/env bash
# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

on_failure() {
  echo ""
  echo "Static analysis failed. Run 'make static-analysis' locally to reproduce."
  echo "This check validates project-specific code conventions that standard linters cannot enforce."
  echo "See the error output above for details on which rules failed and how to fix them."
}

function runStaticAnalysis {
# This tool checks for code conformity within the provider e.g. are the correct Go types used in TypedSDK structs.
  go run internal/tools/static-analysis/main.go || { on_failure; exit 1; }
}

function main {
  runStaticAnalysis
}

main
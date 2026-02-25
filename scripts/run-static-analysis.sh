#!/usr/bin/env bash
# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0


function on_failure {
  echo ""
  echo "==> Static analysis failed!"
  echo "    Check the output above for specific violations."
  echo "    Common issues: incorrect Go types in TypedSDK structs, missing required fields."
  echo "    Run locally: go run internal/tools/static-analysis/main.go -fail-on-error=false"
  echo ""
}

function runStaticAnalysis {
# This tool checks for code conformity within the provider e.g. are the correct Go types used in TypedSDK structs.
# Currently will not fail GHA's etc as we have existing violations in `main`. -fail-on-error=false can be removed when
# these are resolved to prevent PRs introducing this in future.
  go run internal/tools/static-analysis/main.go -fail-on-error=false
}

function main {
  runStaticAnalysis || { on_failure; exit 1; }
}

main
#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


function runStaticAnalysis {
# This tool checks for code conformity within the provider e.g. are the correct Go types used in TypedSDK structs.
# Currently will not fail GHA's etc as we have existing violations in `main`. -fail-on-error=false can be removed when
# these are resolved to prevent PRs introducing this in future.
  go run internal/tools/static-analysis/main.go -fail-on-error=false
}

function main {
  runStaticAnalysis
}

main
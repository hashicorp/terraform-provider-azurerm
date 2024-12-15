#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


function runTests {
  echo "==> Running Unit Tests..."
  go test -i $TEST || exit 1
  go test -v $TEST "$TESTARGS" -timeout=30s -parallel=20
}

function main {
  runTests
}

main

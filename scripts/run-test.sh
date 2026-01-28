#!/usr/bin/env bash
# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0


function runTests {
  echo "==> Running Unit Tests..."
  go test -v $TEST "$TESTARGS" -timeout=30s -parallel=20
}

function main {
  runTests
}

main

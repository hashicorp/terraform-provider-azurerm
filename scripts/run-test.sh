#!/usr/bin/env bash

function runTests {
  echo "==> Running Unit Tests..."
  go test -i $TEST || exit 1
  go test -v $TEST "$TESTARGS" -timeout=30s -parallel=20
}

function main {
  runTests
}

main

#!/usr/bin/env bash

function runLinters {
  echo "==> Checking source code against linters..."
  golangci-lint run -v ./...
}

function main {
  runLinters
}

main

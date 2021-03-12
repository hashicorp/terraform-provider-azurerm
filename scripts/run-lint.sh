#!/usr/bin/env bash

function runLinters {
  echo "==> Checking source code against linters..."
  golangci-lint run ./...
}

function main {
  runLinters
}

main

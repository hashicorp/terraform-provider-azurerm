#!/usr/bin/env bash
# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0


function runLinters {
  echo "==> Checking source code against linters..."
  golangci-lint run -v ./...
}

function main {
  runLinters
}

main

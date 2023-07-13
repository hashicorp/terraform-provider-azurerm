#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


function runLinters {
  echo "==> Checking source code against linters..."
  golangci-lint run -v ./...
}

function main {
  runLinters
}

main

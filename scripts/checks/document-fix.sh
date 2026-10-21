#!/usr/bin/env bash
# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0


echo "==> Validating and fixing documentation..."
if ! go run ./internal/tools/document-fmt/main.go fix; then
  echo "==> Checking for items that are not automatically fixed..."
  if ! go run ./internal/tools/document-fmt/main.go validate; then
    echo "==> Fixing documentation failed, remaining errors require manual fixes..."
  fi
fi

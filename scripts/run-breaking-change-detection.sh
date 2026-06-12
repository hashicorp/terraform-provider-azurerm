#!/usr/bin/env bash
# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0


function on_failure {
  echo ""
  echo "==> Breaking change detection failed!"
  echo "    Your changes modify the provider schema in a backwards-incompatible way."
  echo "    If this is intentional, the change must be gated behind features.FivePointOh()"
  echo "    and documented in the 5.0 upgrade guide."
  echo "    See: contributing/topics/guide-breaking-changes.md"
  echo ""
}

function runDetect {
  go run internal/tools/schema-api/main.go -detect .release/provider-schema.json
}

function main {
  runDetect || { on_failure; exit 1; }
}

main
#!/usr/bin/env bash
# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0


function runDetect {
  go run internal/tools/schema-api/main.go -detect .release/provider-schema.json
}

function main {
  runDetect
}

main
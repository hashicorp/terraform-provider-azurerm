#!/usr/bin/env bash

function runDetect {
  go run internal/tools/schema-api/main.go -detect .release/provider-schema.json
}

function main {
  runDetect
}

main
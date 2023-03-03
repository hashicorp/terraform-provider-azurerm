#!/usr/bin/env bash

function runDetect {
  go run internal/tools/schema-api/main.go -detect azurermProviderSchema.json
}

function main {
  runDetect
}

main
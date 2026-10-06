#!/bin/bash -ex
cd "$(dirname "$0")"
go install tool
mage build coverage
cat report.out

if type -p golangci-lint >/dev/null; then
  golangci-lint run ./...
fi
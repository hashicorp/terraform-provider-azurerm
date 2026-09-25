#!/bin/bash -ex
cd "$(dirname "$0")"
go install tool
mage build coverage
cat report.out
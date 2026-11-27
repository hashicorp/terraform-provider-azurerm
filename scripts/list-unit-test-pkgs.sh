#!/usr/bin/env bash
# Copyright IBM Corp. 2014, 2026
# SPDX-License-Identifier: MPL-2.0

# Emits Go package paths containing at least one unit test, benchmark, example, or fuzz target,
# skipping ~830+ acceptance-only (TestAcc*) packages so unit test runs do not link hundreds of
# ~300-500 MB full-provider test binaries.

set -euo pipefail

exec go run ./internal/tools/test-pkg-filter -mode=unit -path=.

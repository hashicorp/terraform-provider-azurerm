#!/usr/bin/env bash
# Copyright IBM Corp. 2014, 2026
# SPDX-License-Identifier: MPL-2.0

# Ensures internal/sdk and production resource packages under internal/services do not
# re-introduce direct imports of the monolithic internal/clients package (which pulls all
# 135 service clients and >1,500 Azure SDK packages into every subpackage compile unit).

set -euo pipefail

echo "==> Checking that internal/sdk and internal/services production code do not import internal/clients..."

violations="$(
  grep -rnE '"github\.com/hashicorp/terraform-provider-azurerm/internal/clients"' \
    internal/sdk \
    internal/preflight \
    internal/services \
    --include='*.go' \
    --exclude='*_test.go' || true
)"

if [ -n "$violations" ]; then
  echo "ERROR: Direct import of monolithic 'internal/clients' found in production SDK/service files:" >&2
  echo "$violations" >&2
  echo "" >&2
  echo "Use the per-service '<service>/client.From(...)', 'account.From(...)', or 'features.From(...)' helpers instead." >&2
  exit 1
fi

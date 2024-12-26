#!/bin/bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

if [[ "$(uname)" == "Darwin" ]]; then
  echo "(Using BSD sed)"
  SED="sed -E"
else
  echo "(Using GNU sed)"
  SED="sed -r"
fi

DATE="$(date '+%B %d, %Y')"
PROVIDER_URL="https:\/\/github.com\/hashicorp\/terraform-provider-azurerm\/issues"

echo "Preparing changelog for release..."

if [[ ! -f CHANGELOG.md ]]; then
  echo "Error: CHANGELOG.md not found."
  exit 2
fi

# Get the next release
RELEASE="$($SED -n 's/^## v?([0-9.]+) \(Unreleased\)/\1/p' CHANGELOG.md)"
if [[ "${RELEASE}" == "" ]]; then
  echo "Error: could not determine next release in CHANGELOG.md" >&2
  exit 3
fi

# Replace [GH-nnnn] references with issue links
( set -x; ${debug}$SED -i.bak "s/\[GH-([0-9]+)\]/\(\[#\1\]\(${PROVIDER_URL}\/\1\)\)/g" CHANGELOG.md )

# Set the date for the latest release
( set -x; ${debug}$SED -i.bak "s/^(## v?[0-9.]+) \(Unreleased\)/\1 (${DATE})/i" CHANGELOG.md )

${debug}rm CHANGELOG.md.bak

echo "exporting Provider Schema JSON"
(
  set -x
  ${debug}go run internal/tools/schema-api/main.go -export .release/provider-schema.json
)
#!/bin/bash
# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

# Set to "echo " to print commands instead of running them
debug="${debug:-}"

echo "Preparing changelog for release..."

echo "Generating changelog..."
output="$(${debug}changeloggy generate)"
echo "${output}"

RELEASE="$(echo "${output}" | sed -E -n 's/^Generated v?([0-9][0-9.]*) .*/\1/p')"
if [[ "${RELEASE}" == "" ]]; then
  echo "Error: could not determine release version from changeloggy output" >&2
  exit 3
fi

echo "exporting Provider Schema JSON"
(
  set -x
  # shellcheck disable=SC2086 # debug is intentionally unquoted for command prefix pattern
  ${debug}go run internal/tools/schema-api/main.go -export .release/provider-schema.json
)

# Update the version file with this new version
printf "%s" "${RELEASE}" > version/VERSION
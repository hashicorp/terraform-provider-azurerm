#!/usr/bin/env bash
# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

# Runs the schema linter, but only against resources and data sources that are
# NEWLY ADDED in this pull request. It is a no-op unless the PR adds a
# *_resource.go or *_data_source.go file under internal/services/.
#
# The linter is scoped to the new resource/data source type name(s) so that
# pre-existing findings elsewhere in the provider do not affect the result.

set -euo pipefail

echo "==> Detecting newly added resource / data source files..."

# Only files newly ADDED in this PR (not modified/renamed), relative to the
# merge-base with main.
new_files=$(git diff --diff-filter=A origin/main --name-only --merge-base 2>/dev/null || true)

if [ -z "${new_files}" ]; then
  echo "    No new files detected. ✓"
  exit 0
fi

# Keep only newly added resource / data source implementation files (no tests).
schema_files=()
while IFS= read -r f; do
  [ -n "${f}" ] || continue
  case "${f}" in
    internal/services/*_resource.go | internal/services/*_data_source.go | internal/services/*_datasource.go) ;;
    *) continue ;;
  esac
  case "${f}" in
    *_test.go) continue ;;
  esac
  schema_files+=("${f}")
done <<< "${new_files}"

if [ ${#schema_files[@]} -eq 0 ]; then
  echo "    No new resource or data source files detected. ✓"
  exit 0
fi

echo "    New resource / data source files:"
printf '      • %s\n' "${schema_files[@]}"

# Determine the "azurerm_*" type name(s) declared in the new files so the lint
# can be scoped to just them. Prefer the typed-SDK ResourceType() declaration,
# falling back to any azurerm_* string literal in the file.
#
# A while-read loop is used instead of mapfile for portability (macOS ships an
# old bash without mapfile).
resource_names=()
while IFS= read -r name; do
  [ -n "${name}" ] || continue
  resource_names+=("${name}")
done < <(
  grep -hEA2 'ResourceType\(\) string' "${schema_files[@]}" 2>/dev/null |
    grep -oE '"azurerm_[a-z0-9_]+"' | tr -d '"' | sort -u
)

if [ ${#resource_names[@]} -eq 0 ]; then
  while IFS= read -r name; do
    [ -n "${name}" ] || continue
    resource_names+=("${name}")
  done < <(
    grep -hoE '"azurerm_[a-z0-9_]+"' "${schema_files[@]}" 2>/dev/null | tr -d '"' | sort -u
  )
fi

if [ ${#resource_names[@]} -eq 0 ]; then
  echo ""
  echo "::warning::Could not determine the azurerm_* type name from the new file(s); skipping schema-lint."
  echo "New resources should be typed and declare their type via ResourceType()."
  exit 0
fi

names_csv=$(
  IFS=,
  echo "${resource_names[*]}"
)

echo ""
echo "==> Running schema-lint for: ${names_csv}"
echo ""

# Fails (non-zero) on any error-severity finding for the new resource(s);
# warnings are printed but do not fail the build. Pass -fix to surface the
# suggested remediation for fixable findings.
go run ./internal/tools/schema-lint check -resource="${names_csv}" -fix

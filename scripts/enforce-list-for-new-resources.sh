#!/usr/bin/env bash
# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

# This script enforces that all NEW resources added to the provider include a
# corresponding list resource implementation (*_resource_list.go).
#
# This check can be skipped by applying the "allow-without-list" or
# "list-not-supported" label to the pull request.

echo "==> Checking that new resources include a list implementation..."

# Only look at files newly added in this PR (not modified/renamed)
new_files=$(git diff --diff-filter=A origin/main --name-only --merge-base 2>/dev/null || true)

if [ -z "$new_files" ]; then
  echo "    No new files detected. ✓"
  exit 0
fi

missing=()

while IFS= read -r f; do
  # Only consider resource implementation files under internal/services/
  case "$f" in
    internal/services/*_resource.go) ;;
    *) continue ;;
  esac

  # Skip test files, list files
  case "$f" in
    *_test.go)          continue ;;
    *_resource_list.go) continue ;;
  esac

  # Derive the expected list file: foo_resource.go -> foo_resource_list.go
  expected_list="${f%_resource.go}_resource_list.go"

  # Check if the list file exists in the repo (already present or added in this PR)
  if [ ! -f "$expected_list" ]; then
    missing+=("$f")
  fi
done <<< "$new_files"

if [ ${#missing[@]} -eq 0 ]; then
  echo "    All new resources have a corresponding list implementation. ✓"
  exit 0
fi

echo ""
echo "ERROR: The following new resource(s) are missing a list implementation:"
echo ""
for f in "${missing[@]}"; do
  expected_list="${f%_resource.go}_resource_list.go"
  echo "  • $f"
  echo "    expected: $expected_list"
done
echo ""
echo "Every new resource must include a list resource so that it is compatible"
echo "with Terraform's 'list' block (Terraform >= 1.14)."
echo ""
echo "To add a list implementation, create the *_resource_list.go file and"
echo "register it in the service's ListResources() method in registration.go."
echo ""
echo "For detailed instructions, see the contributor guide:"
echo "  https://github.com/hashicorp/terraform-provider-azurerm/blob/main/contributing/topics/guide-list-resource.md"
echo ""
echo "If this resource genuinely cannot support listing (e.g. it has no List API),"
echo "please explain why in the PR description or with a comment and a maintainer will apply the"
echo "'allow-without-list' or 'list-not-supported' label to skip this check."
exit 1

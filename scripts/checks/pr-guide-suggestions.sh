#!/usr/bin/env bash
# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

# Runs the AST-based schema linter, but only against schema properties that this
# pull request ADDS relative to the branch it is merging into.
#
# Renaming or changing an existing property is a breaking change, so those are
# intentionally out of scope: the linter only reports findings on lines the
# change set adds (newly added properties and entirely new resources / data
# sources).
#
# Mechanics:
#   1. Resolve the base branch.
#   2. Build the linter from internal/tools/pr-guide-suggestions.
#   3. Run `pr-guide-suggestions check -diff-base <base>` against the current
#      checkout, which parses the changed Go source directly (no provider build
#      or JSON schema export) and reports only findings on properties added since
#      base.

set -euo pipefail

repo_root="$(git rev-parse --show-toplevel)"
cd "${repo_root}"

# GITHUB_BASE_REF is set by GitHub Actions for pull_request events; default to
# main for local runs.
base_ref="${GITHUB_BASE_REF:-main}"

echo "==> Resolving base branch '${base_ref}'..."
git fetch --no-tags --quiet origin "${base_ref}" 2>/dev/null || true
if git rev-parse --verify --quiet "origin/${base_ref}" >/dev/null; then
  base="origin/${base_ref}"
elif git rev-parse --verify --quiet "${base_ref}" >/dev/null; then
  base="${base_ref}"
else
  echo "::warning::Could not resolve base branch '${base_ref}'; skipping pr-guide-suggestions diff."
  exit 0
fi
echo "    base = ${base} ($(git rev-parse "${base}"))"

# mktemp -d is used (rather than a fixed path) so concurrent runs do not collide.
bin_dir="$(mktemp -d)"
cleanup() {
  rm -rf "${bin_dir}" >/dev/null 2>&1 || true
}
trap cleanup EXIT

echo "==> Building the schema linter..."
go build -o "${bin_dir}/pr-guide-suggestions" ./internal/tools/pr-guide-suggestions

echo ""
echo "==> Linting schema properties added since '${base_ref}'..."
echo ""

"${bin_dir}/pr-guide-suggestions" check -disable SL001 -diff-base "${base}" -C "${repo_root}"

#!/usr/bin/env bash
# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

# Runs the schema linter, but only against schema properties that this pull
# request ADDS relative to the branch it is merging into.
#
# Renaming or changing an existing property is a breaking change, so those are
# intentionally out of scope: the linter only holds NEWLY ADDED properties (and
# entirely new resources / data sources) to the rules.
#
# Mechanics:
#   1. Export the base branch's provider schema from a temporary git worktree.
#   2. Run `schema-lint check -diff <base>` against the current checkout, which
#      reports findings only on properties absent from the base schema.

set -euo pipefail

# GITHUB_BASE_REF is set by GitHub Actions for pull_request events; default to
# main for local runs.
base_ref="${GITHUB_BASE_REF:-main}"

echo "==> Resolving base branch '${base_ref}'..."
git fetch --no-tags --quiet origin "${base_ref}" 2>/dev/null || true
if git rev-parse --verify --quiet "origin/${base_ref}" >/dev/null; then
  base_sha="$(git rev-parse "origin/${base_ref}")"
elif git rev-parse --verify --quiet "${base_ref}" >/dev/null; then
  base_sha="$(git rev-parse "${base_ref}")"
else
  echo "::warning::Could not resolve base branch '${base_ref}'; skipping schema-lint diff."
  exit 0
fi
echo "    base = ${base_sha}"

# mktemp is used (rather than a fixed path) so concurrent runs do not collide.
# The worktree lives in a subdirectory that must NOT pre-exist, so it is created
# under an mktemp -d parent.
base_schema="$(mktemp)"
tmp_parent="$(mktemp -d)"
worktree_dir="${tmp_parent}/base"

cleanup() {
  git worktree remove --force "${worktree_dir}" >/dev/null 2>&1 || true
  rm -rf "${tmp_parent}" "${base_schema}" >/dev/null 2>&1 || true
}
trap cleanup EXIT

echo "==> Exporting the base schema from a worktree at ${base_sha}..."
git worktree add --detach --quiet "${worktree_dir}" "${base_sha}"
(
  cd "${worktree_dir}"
  go run internal/tools/schema-api/main.go -export "${base_schema}"
)

echo ""
echo "==> Linting schema properties added since '${base_ref}'..."
echo ""

# Fails (non-zero) on any error-severity finding for a newly added property;
# warnings are printed but do not fail the build. Pass -fix to surface the
# suggested remediation for fixable findings.
go run ./internal/tools/schema-lint check -diff "${base_schema}" -fix

#!/usr/bin/env bash
# Copyright IBM Corp. 2014, 2026
# SPDX-License-Identifier: MPL-2.0

# A best-effort reminder to update the upgrade guide. It fails when a pull
# request adds any of the following to non-test provider code without changing
# an upgrade guide (website/docs/<version>-upgrade-guide.html.markdown):
#
#   - a major version feature flag check, e.g. features.SixPointOh()
#   - a deprecated property (`Deprecated:`)
#   - a deprecated resource or data source (`DeprecationMessage`,
#     `DeprecatedInFavourOfResource()`, `DeprecatedInFavourOfDataSource()`)
#
# Only added lines are looked at, so moved code is reported as new. To skip a
# line, put `// lintignore: upgrade-guide` on it or on the line above.

set -euo pipefail

# Exit codes: 0 = pass, 1 = the upgrade guide was not updated, 2 = the check
# itself could not run (e.g. a git error).
repo_root="$(git rev-parse --show-toplevel)" || exit 2
cd "${repo_root}" || exit 2

# GITHUB_BASE_REF is set by GitHub Actions for pull requests; local runs use
# main. If origin is a fork, set UPGRADE_GUIDE_BASE (e.g. upstream/main) to
# compare against the right branch.
base_ref="${GITHUB_BASE_REF:-main}"

if [ -n "${UPGRADE_GUIDE_BASE:-}" ]; then
  echo "==> Using base '${UPGRADE_GUIDE_BASE}' from UPGRADE_GUIDE_BASE..."
  if ! git rev-parse --verify --quiet "${UPGRADE_GUIDE_BASE}^{commit}" >/dev/null; then
    echo "::error::UPGRADE_GUIDE_BASE '${UPGRADE_GUIDE_BASE}' is not a commit; the upgrade guide check could not run."
    exit 2
  fi
  base="${UPGRADE_GUIDE_BASE}"
else
  echo "==> Resolving base branch '${base_ref}'..."
  # Only CI fetches: a local run should not move the contributor's
  # remote-tracking branches as a side effect.
  if [ -n "${GITHUB_ACTIONS:-}" ]; then
    git fetch --no-tags --quiet origin "${base_ref}" 2>/dev/null || true
  fi
  if git rev-parse --verify --quiet "origin/${base_ref}" >/dev/null; then
    base="origin/${base_ref}"
  elif git rev-parse --verify --quiet "${base_ref}" >/dev/null; then
    base="${base_ref}"
  else
    # In CI a pass here would hide the check never having run, so fail with the
    # "could not run" exit code; locally there may simply be no base branch.
    if [ -n "${GITHUB_ACTIONS:-}" ]; then
      echo "::error::Could not resolve base branch '${base_ref}'; the upgrade guide check could not run."
      exit 2
    fi
    echo "Could not resolve base branch '${base_ref}'; skipping the upgrade guide check."
    exit 0
  fi
fi
echo "    base = ${base} ($(git rev-parse "${base}"))"

merge_base="$(git merge-base "${base}" HEAD)" || {
  echo "::error::'${base}' and HEAD have no merge base; the upgrade guide check could not run."
  exit 2
}

# Explicit prefixes and no external diff, so a contributor's git config cannot
# change the output. Tests and the tools under internal/tools are not provider
# code.
diff="$(git -c core.quotePath=false diff --no-color --no-ext-diff --src-prefix=a/ --dst-prefix=b/ --unified=0 -M \
  "${merge_base}" -- 'internal/*.go' ':(exclude)internal/*_test.go' ':(exclude)internal/tools/')" || exit 2

# Prints "file<TAB>line<TAB>text" for each added line that is a flag check or a
# deprecation. Whole-line comments are skipped.
candidates="$(printf '%s\n' "${diff}" | awk '
  /^diff --git / { file = ""; in_hunk = 0; next }
  !in_hunk && /^\+\+\+ / { file = ($0 == "+++ /dev/null") ? "" : substr($0, 7); next }
  /^@@/ { split($3, start, ","); line = substr(start[1], 2) + 0; in_hunk = 1; next }
  in_hunk && /^\+/ {
    text = substr($0, 2)
    trimmed = text
    sub(/^[ \t]+/, "", trimmed)
    if (file != "" && trimmed !~ /^\/\// &&
        (text ~ /features\.(DeprecatedIn)?[A-Z][a-z]+PointOh\(/ ||
         text ~ /^[ \t]*Deprecated:/ ||
         text ~ /^[ \t]*DeprecationMessage:/ ||
         text ~ /^func \(.*\) (DeprecationMessage|DeprecatedInFavourOfResource|DeprecatedInFavourOfDataSource)\(\)/)) {
      printf "%s\t%d\t%s\n", file, line, trimmed
    }
    line++
  }
')" || exit 2

ignore_re='lintignore:[[:space:]]*upgrade-guide'
found=()
while IFS=$'\t' read -r file line text; do
  [ -n "${file}" ] || continue
  if [[ "${text}" =~ ${ignore_re} ]]; then
    continue
  fi
  # the directive on the line above must be a comment of its own
  if [ "${line}" -gt 1 ] && [ -f "${file}" ]; then
    above="$(sed -n "$((line - 1))p" "${file}")"
    if [[ "${above}" =~ ^[[:space:]]*// ]] && [[ "${above}" =~ ${ignore_re} ]]; then
      continue
    fi
  fi
  found+=("${file}:${line}: ${text}")
done <<< "${candidates}"

if [ "${#found[@]}" -eq 0 ]; then
  echo "No new major version feature flag checks or deprecations found, nothing to check."
  exit 0
fi

echo ""
echo "==> New major version feature flag checks and deprecations:"
printf '    %s\n' "${found[@]}"
echo ""

# Any upgrade guide in website/docs counts; released guides live in
# website/docs/guides and are not matched.
guide_spec=':(glob)website/docs/*-upgrade-guide.html.markdown'
guide_lines_added="$(git diff --numstat "${merge_base}" -- "${guide_spec}" | awk '{ n += $1 } END { print n + 0 }')" || exit 2

if [ "${guide_lines_added}" -gt 0 ]; then
  echo "OK: this pull request adds ${guide_lines_added} line(s) to the upgrade guide."
  exit 0
fi

if [ -z "$(git ls-files -- "${guide_spec}")" ]; then
  echo "There is no upgrade guide in website/docs to update (the next major version has not been set up), so there is nothing to check."
  exit 0
fi

cat <<'MSG'
FAIL: this pull request does not update the upgrade guide.

Changes behind a major version feature flag, deprecated properties, and
deprecated resources or data sources must be documented in that version's
upgrade guide (website/docs/<version>-upgrade-guide.html.markdown). Add or
update the section for each resource or data source listed above. See
contributing/topics/guide-breaking-changes.md.

This check cannot tell new code from code that was moved or reworded. If every
line listed above is already covered by the upgrade guide, a maintainer can
disregard this failure. To make the check pass, add a
"// lintignore: upgrade-guide" comment on the line or the line above it.
MSG
exit 1

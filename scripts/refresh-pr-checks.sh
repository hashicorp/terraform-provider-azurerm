#!/usr/bin/env bash
# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

# Re-runs CI for open PRs by applying the 'ci-refresh' label, which the
# pr-check workflows accept as a re-trigger (see pr-checks-combined.yaml).
# PR check results go stale when main moves (e.g. a newly enabled linter):
# checks run against a merge with main computed at the last push, and GitHub
# never re-runs them when main changes. Labeling fires a fresh pull_request
# event whose merge commit is recomputed against the current base branch;
# the workflow removes the label again so it acts as a stateless button.
#
# By default only PRs whose checks currently all pass are refreshed (their
# green is what silently goes stale - red PRs already advertise themselves)
# and PRs with merge conflicts are skipped (no merge commit can be built).
#
# Must be run with a user token (gh auth login): label events created by a
# workflow's GITHUB_TOKEN would not trigger the check workflows.
#
# Usage: refresh-pr-checks.sh [-a] [-n] [-l limit] [-s seconds] [-R owner/repo]
#   -a  refresh all open PRs, not just those with passing checks
#   -n  dry run: list the PRs that would be refreshed, change nothing
#   -l  max number of PRs to refresh (default 500)
#   -s  seconds to sleep between PRs, spreads CI load (default 8)
#   -R  repository (default hashicorp/terraform-provider-azurerm)

set -euo pipefail

REPO="hashicorp/terraform-provider-azurerm"
ALL=false
DRY_RUN=false
LIMIT=500
SLEEP_SECONDS=8
LABEL="ci-refresh"

while getopts "anl:s:R:h" opt; do
  case "$opt" in
    a) ALL=true ;;
    n) DRY_RUN=true ;;
    l) LIMIT="$OPTARG" ;;
    s) SLEEP_SECONDS="$OPTARG" ;;
    R) REPO="$OPTARG" ;;
    h)
      grep '^#' "$0" | tail -n +4 | sed 's/^# \{0,1\}//'
      exit 0
      ;;
    *) exit 1 ;;
  esac
done

# the label must exist before it can be applied
if [ "$DRY_RUN" != "true" ] && ! gh label list --repo "$REPO" --json name --jq '.[].name' | grep -qx "$LABEL"; then
  echo "creating the '$LABEL' label..."
  gh label create "$LABEL" --repo "$REPO" \
    --color "1D76DB" \
    --description "Re-runs the PR checks against a fresh merge with the base branch (auto-removed)"
fi

# The search API's status qualifier finds PRs whose checks currently pass;
# it is far cheaper than querying every PR's check rollup (which times out at
# this repo's PR count). PRs with no checks at all (e.g. fork PRs awaiting
# workflow approval) don't match status:success, which is what we want:
# refreshing them would still leave them waiting on approval.
QUERY="repo:$REPO is:pr is:open"
if [ "$ALL" != "true" ]; then
  QUERY="$QUERY status:success"
fi

PRS=$(gh api -X GET search/issues --paginate \
  -f q="$QUERY" -F per_page=100 --jq '.items[].number' | head -n "$LIMIT")

if [ -z "$PRS" ]; then
  echo "no PRs match"
  exit 0
fi

COUNT=$(echo "$PRS" | wc -l | tr -d ' ')
echo "refreshing checks on $COUNT PRs..."

for n in $PRS; do
  if [ "$DRY_RUN" = "true" ]; then
    echo "would refresh #$n"
    continue
  fi

  # reading mergeable makes GitHub recompute the speculative merge commit
  # against the current base branch, so the labeled event tests fresh state.
  # It also filters conflicted PRs, which have no merge commit to check.
  state="null"
  for _ in 1 2 3 4 5; do
    state=$(gh api "repos/$REPO/pulls/$n" --jq '.mergeable') || break
    [ "$state" != "null" ] && break
    sleep 2
  done
  if [ "$state" = "false" ]; then
    echo "skipping #$n (merge conflict)"
    continue
  fi

  echo "refreshing #$n"
  gh pr edit "$n" --repo "$REPO" --add-label "$LABEL"
  sleep "$SLEEP_SECONDS"
done

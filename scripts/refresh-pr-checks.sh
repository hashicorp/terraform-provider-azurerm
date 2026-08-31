#!/usr/bin/env bash
# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

# Re-runs CI for open PRs by merging the current base branch into them via the
# update-branch API (the "Update branch" button). PR check results go stale when
# main moves (e.g. a newly enabled linter): checks run against a merge with main
# computed at the last push, and GitHub never re-runs them when main changes.
# Updating the branch pushes a real merge commit, which fires a genuine
# 'synchronize' event, so every check workflow re-runs natively against current
# main - no workflow changes, labels, or skipped-run noise involved.
#
# By default only PRs whose checks currently all pass are refreshed (their
# green is what silently goes stale - red PRs already advertise themselves).
# Skipped automatically: PRs with merge conflicts (no merge commit can be
# built), PRs already up to date with the base branch (nothing is stale), and
# fork PRs without "maintainer can modify" (the API returns 403).
#
# Must be run with a user token (gh auth login): synchronize events pushed by
# a workflow's GITHUB_TOKEN would not trigger the check workflows.
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
echo "refreshing checks on up to $COUNT PRs..."

for n in $PRS; do
  # reading mergeable makes GitHub compute the merge preview against the current
  # base branch, and filters conflicted PRs, which have no merge commit to build
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

  if [ "$DRY_RUN" = "true" ]; then
    echo "would refresh #$n"
    continue
  fi

  echo "refreshing #$n"
  # 422 = already up to date with base; 403 = fork without maintainer-edit -
  # both are fine to skip, anything else is worth seeing on stderr
  if ! gh api -X PUT "repos/$REPO/pulls/$n/update-branch" > /dev/null; then
    echo "skipping #$n (up to date, or branch not updatable)"
    continue
  fi
  sleep "$SLEEP_SECONDS"
done

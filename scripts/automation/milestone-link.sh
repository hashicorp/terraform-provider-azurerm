#!/usr/bin/env bash
# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

# Assigns a merged PR and all issues it closes to the next unreleased milestone
# (the lowest open vX.Y.0 milestone). Run from the milestone-link workflow.
#
# Required environment:
#   GH_TOKEN          - token with issues:write and pull-requests:write
#   GITHUB_REPOSITORY - owner/repo (set automatically by GitHub Actions)
#   PR_NUMBER         - number of the merged pull request

set -euo pipefail

repo="${GITHUB_REPOSITORY:?GITHUB_REPOSITORY must be set}"
pr_number="${PR_NUMBER:?PR_NUMBER must be set}"

milestone=$(gh api "repos/${repo}/milestones?state=open&per_page=100" \
  --jq '[.[] | select(.title | test("^v[0-9]+\\.[0-9]+\\.0$"))]
        | sort_by(.title | ltrimstr("v") | split(".") | map(tonumber))
        | first // empty')

if [[ -z "${milestone}" ]]; then
  echo "No open vX.Y.0 milestones found, nothing to do"
  exit 0
fi

milestone_number=$(jq -r '.number' <<<"${milestone}")
milestone_title=$(jq -r '.title' <<<"${milestone}")
echo "Next unreleased milestone: ${milestone_title}"

assign_milestone() {
  local number="$1"
  local existing
  existing=$(gh api "repos/${repo}/issues/${number}" --jq '.milestone.title // empty')
  if [[ -n "${existing}" ]]; then
    echo "#${number} is already in milestone ${existing}, skipping"
    return
  fi
  gh api --silent -X PATCH "repos/${repo}/issues/${number}" -F milestone="${milestone_number}"
  echo "#${number} added to milestone ${milestone_title}"
}

assign_milestone "${pr_number}"

# every issue this PR closes, whether via closing keywords or linked manually in
# the UI - not just the first one mentioned in the PR body
linked_issues=$(gh api graphql \
  -F owner="${repo%/*}" -F repo="${repo#*/}" -F pr="${pr_number}" \
  -f query='
    query($owner: String!, $repo: String!, $pr: Int!) {
      repository(owner: $owner, name: $repo) {
        pullRequest(number: $pr) {
          closingIssuesReferences(first: 100) {
            nodes { number repository { nameWithOwner } }
          }
        }
      }
    }' \
  --jq ".data.repository.pullRequest.closingIssuesReferences.nodes[]
        | select(.repository.nameWithOwner == \"${repo}\") | .number")

for issue in ${linked_issues}; do
  assign_milestone "${issue}"
done

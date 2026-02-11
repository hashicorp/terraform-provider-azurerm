#!/usr/bin/env bash
# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

# shellcheck disable=SC2086

# as the normal stale bot does not work for draft PRs, we need to do it manually until this pr is merged https://github.com/actions/stale/pull/1314

set -euo pipefail

DRY_RUN=false

while getopts o:r:t:d flag
do
  case "${flag}" in
    o) owner=${OPTARG};;
    r) repo=${OPTARG};;
    t) token=${OPTARG};;
    d) DRY_RUN=true;;
    *) echo "Usage: $0 -o owner -r repo [-t token] [-d]"; exit 1;;
  esac
done

# Use token from env if not provided via flag
token="${token:-${GH_TOKEN:-${GITHUB_TOKEN:-}}}"
if [[ -z "$token" ]]; then
  echo "Error: No token provided. Use -t flag or set GH_TOKEN/GITHUB_TOKEN env var."
  exit 1
fi

API_BASE="https://api.github.com/repos/${owner}/${repo}"
# Handle both GNU date (Linux) and BSD date (macOS)
SIXTY_DAYS_AGO=$(date -u -d "60 days ago" +%Y-%m-%dT%H:%M:%SZ 2>/dev/null || date -u -v-60d +%Y-%m-%dT%H:%M:%SZ)

echo "=== Close Inactive Draft PRs ==="
echo "Repository: ${owner}/${repo}"
echo "Close threshold: 60 days (before ${SIXTY_DAYS_AGO})"
if [[ "$DRY_RUN" == "true" ]]; then
  echo "Mode: DRY RUN (no changes will be made)"
else
  echo "Mode: LIVE"
fi
echo ""

# Collect all open PRs
echo "Fetching open PRs..."
all_prs=()
page=1
while :; do
  prs_json=$(curl -s -L \
    -H "Accept: application/vnd.github+json" \
    -H "Authorization: Bearer $token" \
    -H "X-GitHub-Api-Version: 2022-11-28" \
    "${API_BASE}/pulls?state=open&per_page=100&page=${page}")
  
  count=$(echo "$prs_json" | jq length)
  if [[ "$count" -eq 0 ]]; then
    break
  fi
  
  while IFS= read -r pr; do
    all_prs+=("$pr")
  done < <(echo "$prs_json" | jq -c '.[]')
  
  page=$((page + 1))
done

total_prs=${#all_prs[@]}
echo "Found ${total_prs} open PR(s)"
echo ""

# Process each PR
drafts_found=0
close_count=0

for pr in "${all_prs[@]}"; do
  draft=$(echo "$pr" | jq -r '.draft')
  
  # Skip non-draft PRs silently
  if [[ "$draft" != "true" ]]; then
    continue
  fi
  
  pr_number=$(echo "$pr" | jq -r '.number')
  pr_title=$(echo "$pr" | jq -r '.title')
  updated_at=$(echo "$pr" | jq -r '.updated_at')
  labels=$(echo "$pr" | jq -r '.labels[].name' 2>/dev/null || echo "")
  
  drafts_found=$((drafts_found + 1))
  echo "Draft PR #${pr_number} \"${pr_title}\" (updated: ${updated_at})"
  
  # Check for keep-draft label
  if echo "$labels" | grep -q "^keep-draft$"; then
    echo "  ↳ Has 'keep-draft' label, skipping"
    continue
  fi
  
  # Close if 60+ days inactive
  if [[ "$updated_at" < "$SIXTY_DAYS_AGO" ]]; then
    close_count=$((close_count + 1))
    echo "  ↳ Inactive for 60+ days → CLOSING"
    
    if [[ "$DRY_RUN" == "false" ]]; then
      curl -s -L -X PATCH \
        -H "Accept: application/vnd.github+json" \
        -H "Authorization: Bearer $token" \
        -H "X-GitHub-Api-Version: 2022-11-28" \
        "${API_BASE}/pulls/${pr_number}" \
        -d '{"state":"closed"}' > /dev/null
      
      curl -s -L -X POST \
        -H "Accept: application/vnd.github+json" \
        -H "Authorization: Bearer $token" \
        -H "X-GitHub-Api-Version: 2022-11-28" \
        "${API_BASE}/issues/${pr_number}/comments" \
        -d '{"body":"I'\''m going to close this draft pull request because it has been inactive for _60 days_ ⏳. This helps our maintainers find and focus on the active contributions.\n\nIf you would like to continue working on this, please reopen the pull request and mark it as ready for review when complete. Thank you!"}' > /dev/null
      echo "  ↳ Closed and commented"
    else
      echo "  ↳ (dry run - would close and comment)"
    fi
  fi
done

echo ""
echo "=== Summary ==="
echo "Total PRs checked: ${total_prs}"
echo "Draft PRs found: ${drafts_found}"
echo "Closed: ${close_count}"
if [[ "$DRY_RUN" == "true" ]]; then
  echo "(dry run - no actual changes made)"
fi
echo "Done."

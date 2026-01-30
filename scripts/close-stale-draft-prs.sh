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
THIRTY_DAYS_AGO=$(date -u -d "30 days ago" +%Y-%m-%dT%H:%M:%SZ 2>/dev/null || date -u -v-30d +%Y-%m-%dT%H:%M:%SZ)
SIXTY_DAYS_AGO=$(date -u -d "60 days ago" +%Y-%m-%dT%H:%M:%SZ 2>/dev/null || date -u -v-60d +%Y-%m-%dT%H:%M:%SZ)

echo "=== Stale Draft PR Processor ==="
echo "Repository: ${owner}/${repo}"
echo "Stale threshold: 30 days (before ${THIRTY_DAYS_AGO})"
echo "Close threshold: 60 days (before ${SIXTY_DAYS_AGO})"
if [[ "$DRY_RUN" == "true" ]]; then
  echo "Mode: DRY RUN (no changes will be made)"
else
  echo "Mode: LIVE"
fi
echo ""

# First, collect all open PRs to get total count
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
current=0
drafts_found=0
stale_count=0
close_count=0

for pr in "${all_prs[@]}"; do
  current=$((current + 1))
  
  pr_number=$(echo "$pr" | jq -r '.number')
  pr_title=$(echo "$pr" | jq -r '.title' | head -c 60)
  draft=$(echo "$pr" | jq -r '.draft')
  updated_at=$(echo "$pr" | jq -r '.updated_at')
  labels=$(echo "$pr" | jq -r '.labels[].name' 2>/dev/null || echo "")
  
  echo "PR ${current}/${total_prs} - #${pr_number} \"${pr_title}\""
  
  # Check if draft
  if [[ "$draft" != "true" ]]; then
    echo "  ↳ Not a draft, skipping"
    continue
  fi
  
  drafts_found=$((drafts_found + 1))
  echo "  ↳ Is draft, last updated: ${updated_at}"
  
  # Check for keep-draft label
  if echo "$labels" | grep -q "^keep-draft$"; then
    echo "  ↳ Has 'keep-draft' label, skipping"
    continue
  fi
  
  # Check for existing stale-draft label
  has_stale_label=false
  if echo "$labels" | grep -q "^stale-draft$"; then
    has_stale_label=true
    echo "  ↳ Already has 'stale-draft' label"
  fi
  
  # Close if 60+ days stale
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
  
  # Comment and label if 30+ days stale without stale-draft label
  elif [[ "$updated_at" < "$THIRTY_DAYS_AGO" && "$has_stale_label" == "false" ]]; then
    stale_count=$((stale_count + 1))
    echo "  ↳ Inactive for 30+ days → MARKING STALE"
    
    if [[ "$DRY_RUN" == "false" ]]; then
      curl -s -L -X POST \
        -H "Accept: application/vnd.github+json" \
        -H "Authorization: Bearer $token" \
        -H "X-GitHub-Api-Version: 2022-11-28" \
        "${API_BASE}/issues/${pr_number}/labels" \
        -d '{"labels":["stale-draft"]}' > /dev/null
      
      curl -s -L -X POST \
        -H "Accept: application/vnd.github+json" \
        -H "Authorization: Bearer $token" \
        -H "X-GitHub-Api-Version: 2022-11-28" \
        "${API_BASE}/issues/${pr_number}/comments" \
        -d '{"body":"This draft pull request is being labeled as \"stale-draft\" because it has not been updated for _30 days_ ⏳.\n\nIf this draft is still valid, please remove the \"stale-draft\" label. To prevent automatic closure after 60 days, add the `keep-draft` label or mark it as ready for review.\n\nIf you need some help completing this draft, please leave a comment letting us know. Thank you!"}' > /dev/null
      echo "  ↳ Labeled and commented"
    else
      echo "  ↳ (dry run - would label and comment)"
    fi
  
  elif [[ "$has_stale_label" == "true" ]]; then
    echo "  ↳ Already stale, not yet 60 days old"
  else
    echo "  ↳ Not stale yet (updated within 30 days)"
  fi
done

echo ""
echo "=== Summary ==="
echo "Total PRs checked: ${total_prs}"
echo "Draft PRs found: ${drafts_found}"
echo "Newly marked stale: ${stale_count}"
echo "Closed: ${close_count}"
if [[ "$DRY_RUN" == "true" ]]; then
  echo "(dry run - no actual changes made)"
fi
echo "Done."

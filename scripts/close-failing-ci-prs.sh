#!/usr/bin/env bash
# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

# shellcheck disable=SC2086

# Warns and closes PRs that have had failing CI for an extended period.
# Checks CI status directly via the GitHub API (does not depend on labels).
#
# - After 7 days of failing CI: leaves a warning comment
# - After 14 days of failing CI: closes the PR with a polite message
# - PRs with "ci-ignore-failure" label are skipped

set -euo pipefail

DRY_RUN=false
WARN_DAYS=7
CLOSE_DAYS=14
IGNORE_LABEL="ci-ignore-failure"
WARNING_MARKER="<!-- ci-failure-warning -->"

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

echo "=== Close PRs With Failing CI ==="
echo "Repository: ${owner}/${repo}"
echo "Warn after: ${WARN_DAYS} days of failing CI"
echo "Close after: ${CLOSE_DAYS} days of failing CI"
echo "Ignore label: ${IGNORE_LABEL}"
if [[ "$DRY_RUN" == "true" ]]; then
  echo "Mode: DRY RUN (no changes will be made)"
else
  echo "Mode: LIVE"
fi
echo ""

# Fetch all open PRs
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

# Check CI status for a commit SHA.
# Outputs "status|timestamp|failed_checks" where:
#   status = "failure" or "ok"
#   timestamp = earliest failure time (or PR updated_at as fallback)
#   failed_checks = comma-separated list of failed check/status names
check_ci_status() {
  local sha=$1
  local pr_updated_at=$2

  local has_failure=false
  local earliest_failure=""
  local failed_names=""

  # Check combined commit status (legacy status API)
  local status_json
  status_json=$(curl -s -L \
    -H "Accept: application/vnd.github+json" \
    -H "Authorization: Bearer $token" \
    -H "X-GitHub-Api-Version: 2022-11-28" \
    "${API_BASE}/commits/${sha}/status")

  local combined_state
  combined_state=$(echo "$status_json" | jq -r '.state')

  # Collect failure timestamps and names from legacy statuses
  if [[ "$combined_state" == "failure" ]]; then
    has_failure=true
    local status_failure_time
    status_failure_time=$(echo "$status_json" | jq -r '[.statuses[] | select(.state == "failure" or .state == "error") | .updated_at] | sort | first // empty')
    if [[ -n "$status_failure_time" ]]; then
      earliest_failure="$status_failure_time"
    fi
    local status_names
    status_names=$(echo "$status_json" | jq -r '[.statuses[] | select(.state == "failure" or .state == "error") | .context] | join(",")')
    if [[ -n "$status_names" ]]; then
      failed_names="$status_names"
    fi
  fi

  # Check check runs (checks API) with pagination
  local check_page=1
  while :; do
    local check_runs_json
    check_runs_json=$(curl -s -L \
      -H "Accept: application/vnd.github+json" \
      -H "Authorization: Bearer $token" \
      -H "X-GitHub-Api-Version: 2022-11-28" \
      "${API_BASE}/commits/${sha}/check-runs?per_page=100&page=${check_page}")

    local page_count
    page_count=$(echo "$check_runs_json" | jq '.check_runs | length')
    if [[ "$page_count" -eq 0 ]]; then
      break
    fi

    # Find failed check runs - times and names
    local failed_time
    failed_time=$(echo "$check_runs_json" | jq -r '[.check_runs[] | select(.conclusion == "failure" or .conclusion == "cancelled") | .completed_at // empty] | map(select(. != "")) | sort | first // empty')

    local check_names
    check_names=$(echo "$check_runs_json" | jq -r '[.check_runs[] | select(.conclusion == "failure" or .conclusion == "cancelled") | .name] | join(",")')

    if [[ -n "$failed_time" ]]; then
      has_failure=true
      if [[ -z "$earliest_failure" ]] || [[ "$failed_time" < "$earliest_failure" ]]; then
        earliest_failure="$failed_time"
      fi
    fi

    if [[ -n "$check_names" ]]; then
      has_failure=true
      if [[ -n "$failed_names" ]]; then
        failed_names="${failed_names},${check_names}"
      else
        failed_names="$check_names"
      fi
    fi

    if [[ "$page_count" -lt 100 ]]; then
      break
    fi
    check_page=$((check_page + 1))
  done

  if [[ "$has_failure" == "true" ]]; then
    echo "failure|${earliest_failure:-$pr_updated_at}|${failed_names}"
  else
    echo "ok||"
  fi
}

# Map a failed check name to actionable guidance
get_check_guidance() {
  local check_name=$1
  case "$check_name" in
    depscheck|"Vendor Dependencies Check")
      echo "Run \`make depscheck\`. Do not modify files in the \`vendor/\` directory directly - instead update dependencies in \`go.mod\` and run \`go mod vendor\`."
      ;;
    website-lint|"Website Linting")
      echo "Run \`make website-lint\` and \`make document-validate\` locally. Check your documentation files under \`website/\` for formatting issues."
      ;;
    gencheck|"Generation Check")
      echo "Run \`make generate\` to regenerate any auto-generated code, then commit the changes."
      ;;
    golint|"GoLang Linting")
      echo "Run the Go linter locally with \`golangci-lint run ./internal/...\` and fix any reported issues."
      ;;
    tflint|"Terraform Schema Linting")
      echo "Run \`make tflint\` locally and fix any Terraform schema issues in your resource/data source definitions."
      ;;
    detect|"Breaking Schema Changes")
      echo "Your changes contain breaking schema changes. Please review the [breaking changes guide](contributing/topics/guide-breaking-changes.md) and ensure any breaking changes are behind the appropriate feature flag."
      ;;
    test|"Unit Tests")
      echo "Run \`make test\` locally to reproduce and fix the failing unit tests."
      ;;
    "Static Analysis")
      echo "Run \`bash ./scripts/run-static-analysis.sh\` locally and fix any reported issues."
      ;;
    shellcheck|"ShellCheck Scripts")
      echo "Run \`make shellcheck\` to check shell scripts for issues."
      ;;
    "Validate Examples")
      echo "Run \`make validate-examples\` to check that your example configurations are valid."
      ;;
    compatibility-32bit-test|"32 Bit Build")
      echo "The 32-bit build is failing. Run \`GOARCH=386 GOOS=linux go build ./...\` locally to diagnose."
      ;;
    "Preview API Version Linter")
      echo "Check that any API version references are not using preview versions unless explicitly required."
      ;;
    *)
      echo "Check the CI logs for details on this failure."
      ;;
  esac
}

# Build a guidance section from a comma-separated list of failed check names
build_guidance() {
  local failed_checks=$1
  if [[ -z "$failed_checks" ]]; then
    echo ""
    return
  fi

  local guidance="\\n\\n**Failing checks and how to fix them:**\\n"
  local seen=""

  IFS=',' read -ra checks <<< "$failed_checks"
  for check in "${checks[@]}"; do
    # Deduplicate
    if echo "$seen" | grep -qF "|${check}|"; then
      continue
    fi
    seen="${seen}|${check}|"

    local fix
    fix=$(get_check_guidance "$check")
    guidance="${guidance}\\n- **${check}**: ${fix}"
  done

  echo "$guidance"
}

# Check if we already left a warning comment (uses a hidden HTML marker)
has_warning_comment() {
  local pr_number=$1

  local comments_json
  comments_json=$(curl -s -L \
    -H "Accept: application/vnd.github+json" \
    -H "Authorization: Bearer $token" \
    -H "X-GitHub-Api-Version: 2022-11-28" \
    "${API_BASE}/issues/${pr_number}/comments?per_page=100")

  echo "$comments_json" | jq -r '[.[] | select(.body | contains("'"${WARNING_MARKER}"'"))] | length'
}

# Process each PR
warn_count=0
close_count=0
skip_count=0
failing_count=0

for pr in "${all_prs[@]}"; do
  draft=$(echo "$pr" | jq -r '.draft')

  # Skip draft PRs
  if [[ "$draft" == "true" ]]; then
    continue
  fi

  pr_number=$(echo "$pr" | jq -r '.number')
  pr_title=$(echo "$pr" | jq -r '.title')
  head_sha=$(echo "$pr" | jq -r '.head.sha')
  updated_at=$(echo "$pr" | jq -r '.updated_at')
  labels=$(echo "$pr" | jq -r '.labels[].name' 2>/dev/null || echo "")
  pr_author=$(echo "$pr" | jq -r '.user.login')

  # Skip PRs with ignore label
  if echo "$labels" | grep -q "^${IGNORE_LABEL}$"; then
    continue
  fi

  # Check CI status - returns "status|timestamp|failed_checks"
  ci_result=$(check_ci_status "$head_sha" "$updated_at")
  ci_status=$(echo "$ci_result" | cut -d'|' -f1)
  ci_failed_since=$(echo "$ci_result" | cut -d'|' -f2)
  ci_failed_checks=$(echo "$ci_result" | cut -d'|' -f3)

  if [[ "$ci_status" != "failure" ]]; then
    continue
  fi

  failing_count=$((failing_count + 1))
  echo "PR #${pr_number} \"${pr_title}\""

  if [[ -z "$ci_failed_since" ]]; then
    echo "  ↳ CI failing but could not determine since when, skipping"
    skip_count=$((skip_count + 1))
    continue
  fi

  # Calculate days since CI started failing
  # Handle both GNU date (Linux) and BSD date (macOS)
  failed_epoch=$(date -d "$ci_failed_since" +%s 2>/dev/null || date -jf "%Y-%m-%dT%H:%M:%SZ" "$ci_failed_since" +%s 2>/dev/null || date -jf "%Y-%m-%dT%T%z" "$ci_failed_since" +%s)
  now_epoch=$(date -u +%s)
  days_since=$(( (now_epoch - failed_epoch) / 86400 ))

  echo "  ↳ CI failing since: ${ci_failed_since} (${days_since} days)"
  if [[ -n "$ci_failed_checks" ]]; then
    echo "  ↳ Failed checks: ${ci_failed_checks}"
  fi

  # Build guidance text for the comment
  guidance=$(build_guidance "$ci_failed_checks")

  # Close if past close threshold
  if [[ "$days_since" -ge "$CLOSE_DAYS" ]]; then
    close_count=$((close_count + 1))
    echo "  ↳ CI failing for ${days_since} days -> CLOSING"

    if [[ "$DRY_RUN" == "false" ]]; then
      comment_body="${WARNING_MARKER}\nThank you for your contribution @${pr_author}. Unfortunately, we are unable to review or merge this pull request as the CI checks have been failing for more than 14 days.\n\nPlease feel free to reopen this PR once the CI issues have been resolved.${guidance}\n\nThank you for your understanding!"

      # Use jq to safely build JSON
      json_payload=$(jq -n --arg body "$comment_body" '{"body": $body}')

      curl -s -L -X POST \
        -H "Accept: application/vnd.github+json" \
        -H "Authorization: Bearer $token" \
        -H "X-GitHub-Api-Version: 2022-11-28" \
        "${API_BASE}/issues/${pr_number}/comments" \
        -d "$json_payload" > /dev/null

      curl -s -L -X PATCH \
        -H "Accept: application/vnd.github+json" \
        -H "Authorization: Bearer $token" \
        -H "X-GitHub-Api-Version: 2022-11-28" \
        "${API_BASE}/pulls/${pr_number}" \
        -d '{"state":"closed"}' > /dev/null
      echo "  ↳ Closed and commented"
    else
      echo "  ↳ (dry run - would close and comment)"
    fi

  # Warn if past warn threshold and not already warned
  elif [[ "$days_since" -ge "$WARN_DAYS" ]]; then
    existing_warnings=$(has_warning_comment "$pr_number")
    if [[ "$existing_warnings" -gt 0 ]]; then
      echo "  ↳ Already warned, waiting for close threshold"
      skip_count=$((skip_count + 1))
      continue
    fi

    warn_count=$((warn_count + 1))
    echo "  ↳ CI failing for ${days_since} days -> WARNING"

    if [[ "$DRY_RUN" == "false" ]]; then
      comment_body="${WARNING_MARKER}\nHi @${pr_author}, we have noticed that the CI on this pull request has been failing for 7 days.\n\nIf the CI failures are not resolved within the next 7 days, we will close this pull request.${guidance}\n\nIf you need help, please leave a comment and we will do our best to assist. Thank you!"

      json_payload=$(jq -n --arg body "$comment_body" '{"body": $body}')

      curl -s -L -X POST \
        -H "Accept: application/vnd.github+json" \
        -H "Authorization: Bearer $token" \
        -H "X-GitHub-Api-Version: 2022-11-28" \
        "${API_BASE}/issues/${pr_number}/comments" \
        -d "$json_payload" > /dev/null
      echo "  ↳ Warning comment posted"
    else
      echo "  ↳ (dry run - would post warning comment)"
    fi
  else
    echo "  ↳ Under threshold (${days_since}/${WARN_DAYS} days), skipping"
    skip_count=$((skip_count + 1))
  fi
done

echo ""
echo "=== Summary ==="
echo "Total PRs checked: ${total_prs}"
echo "PRs with failing CI: ${failing_count}"
echo "Warned: ${warn_count}"
echo "Closed: ${close_count}"
echo "Skipped: ${skip_count}"
if [[ "$DRY_RUN" == "true" ]]; then
  echo "(dry run - no actual changes made)"
fi
echo "Done."

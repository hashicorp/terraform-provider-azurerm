#!/bin/bash

POST_GITHUB_COMMENT="%POST_GITHUB_COMMENT%"
GITHUB_REPO="%env.GITHUB_REPO%"
GIT_PAT="%env.GIT_PAT%"
TEAMCITY_TOKEN="%env.TEAMCITY_TOKEN%"
TEAMCITY_SERVER_URL="%teamcity.serverUrl%"
# BUILD_START_TIME is set in the first build step: $(date +%s)
BUILD_ID="%teamcity.build.id%"
BUILD_TYPE_ID="%system.teamcity.buildType.id%"
BUILD_START_TIME="%env.BUILD_START_TIME%"
BETA_VERSION_ENV_VAR="%env.BETA_VERSION_ENV_VAR%"
TEAMCITY_BUILD_BRANCH="%teamcity.build.branch%"
LABEL_SUCCESS="%env.LABEL_SUCCESS%"
LABEL_FAILURE="%env.LABEL_FAILURE%"
LABEL_OUTDATED="%env.LABEL_OUTDATED%"
LABEL_NEW_FAILURE="%env.LABEL_NEW_FAILURE%"
APPLY_TESTING_LABELS_ENABLED="%env.APPLY_TESTING_LABELS_ENABLED%"
TRACKING_ID="%TRACKING_ID%"

if [ "$POST_GITHUB_COMMENT" != "true" ]; then
  echo "GitHub commenting disabled — skipping."
  exit 0
fi

if [[ "$TEAMCITY_BUILD_BRANCH" =~ refs/pull/([0-9]+)/merge ]]; then
  PR_NUMBER="${BASH_REMATCH[1]}"
else
  echo "Not a PR merge branch: $TEAMCITY_BUILD_BRANCH"
  exit 0
fi

echo "Tracking ID: $TRACKING_ID"

github_api_request() {
  local endpoint="$1"
  curl -s \
    -H "Authorization: Bearer $GIT_PAT" \
    -H "Accept: application/vnd.github+json" \
    "https://api.github.com/repos/$GITHUB_REPO${endpoint}"
}

# Apply a label to the PR
apply_label() {
  local label="$1"
  echo "Applying label: $label"
  curl -s -X POST \
    -H "Authorization: Bearer $GIT_PAT" \
    -H "Accept: application/vnd.github+json" \
    -H "Content-Type: application/json" \
    "https://api.github.com/repos/$GITHUB_REPO/issues/${PR_NUMBER}/labels" \
    -d "{\"labels\":[\"$label\"]}"
}

# Remove a label from the PR
remove_label() {
  local label="$1"
  echo "Removing label: $label"
  curl -s -X DELETE \
    -H "Authorization: Bearer $GIT_PAT" \
    -H "Accept: application/vnd.github+json" \
    "https://api.github.com/repos/$GITHUB_REPO/issues/${PR_NUMBER}/labels/$label" \
    2>/dev/null || true
}

set_testing_label() {
  local label="$1"
  if [ "$label" = "$LABEL_SUCCESS" ]; then
    remove_label "$LABEL_OUTDATED"
    remove_label "$LABEL_FAILURE"
    remove_label "$LABEL_NEW_FAILURE"
    apply_label "$LABEL_SUCCESS"
  elif [ "$label" = "$LABEL_FAILURE" ]; then
    remove_label "$LABEL_OUTDATED"
    remove_label "$LABEL_SUCCESS"
    remove_label "$LABEL_NEW_FAILURE"
    apply_label "$LABEL_FAILURE"
  elif [ "$label" = "$LABEL_NEW_FAILURE" ]; then
    remove_label "$LABEL_OUTDATED"
    remove_label "$LABEL_SUCCESS"
    remove_label "$LABEL_FAILURE"
    apply_label "$LABEL_NEW_FAILURE"
  fi
}

# Fetch test results for this build from the TeamCity REST API
TEAMCITY_ERROR=""
RAW_TEST_RESULTS_JSON=$(curl -sS -f \
  -H "Authorization: Bearer $TEAMCITY_TOKEN" \
  -H "Accept: application/json" \
  "$TEAMCITY_SERVER_URL/app/rest/testOccurrences?locator=build:(id:$BUILD_ID),count:100000&fields=testOccurrence(name,status,duration,newFailure,test(id),firstFailed(build(id,number,branchName)))")

if [ $? -ne 0 ] || [ -z "$RAW_TEST_RESULTS_JSON" ]; then
  TEAMCITY_ERROR="Failed to fetch test results from TeamCity for build $BUILD_ID."
  TEST_RESULTS=""
else
  TEST_RESULTS=$(echo "$RAW_TEST_RESULTS_JSON" | jq -r '(.testOccurrence // [])[]
      | select(.status == "SUCCESS" or .status == "FAILURE")
      | "\(.name)|\(if .status == "SUCCESS" then "PASS" else "FAIL" end)|\((.duration // 0) / 1000)|"')

  if [ $? -ne 0 ]; then
    TEAMCITY_ERROR="Failed to parse TeamCity test results for build $BUILD_ID."
    TEST_RESULTS=""
  fi
fi

PASS_COUNT=$(echo "$TEST_RESULTS" | awk -F'|' 'BEGIN{c=0} $1!="" && $2=="PASS"{c++} END{print c}')
FAIL_COUNT=$(echo "$TEST_RESULTS" | awk -F'|' 'BEGIN{c=0} $1!="" && $2=="FAIL"{c++} END{print c}')
TOTAL=$((PASS_COUNT + FAIL_COUNT))

# Fetch main branch test results early to identify new failures for comment marking
NEW_FAILURES=""
MAIN_TEST_RESULTS=""
MAIN_TEST_ID_MAP=""
if [ -z "$TEAMCITY_ERROR" ] && [ "$FAIL_COUNT" -gt 0 ]; then
  echo "Fetching main branch test results for comparison..."

  # Get the latest successful build on main branch
  MAIN_BUILD_INFO=$(curl -s \
    -H "Authorization: Bearer $TEAMCITY_TOKEN" \
    -H "Accept: application/json" \
    "$TEAMCITY_SERVER_URL/app/rest/builds?locator=buildType:(id:$BUILD_TYPE_ID),branch:refs/heads/main,status:SUCCESS,count:1")

  MAIN_BUILD_ID=$(echo "$MAIN_BUILD_INFO" | jq -r '.build[0].id // empty')

  if [ -n "$MAIN_BUILD_ID" ]; then
    echo "Found main branch build: $MAIN_BUILD_ID"

    # Fetch test results from main branch build via the TeamCity REST API.
    MAIN_RAW_JSON=$(curl -s \
      -H "Authorization: Bearer $TEAMCITY_TOKEN" \
      -H "Accept: application/json" \
      "$TEAMCITY_SERVER_URL/app/rest/testOccurrences?locator=build:(id:$MAIN_BUILD_ID),count:100000&fields=testOccurrence(name,status,test(id))")

    # Build name→status pairs (for new-failure detection)
    MAIN_TEST_RESULTS=$(echo "$MAIN_RAW_JSON" \
      | jq -r '(.testOccurrence // [])[]
          | select(.status == "SUCCESS" or .status == "FAILURE")
          | "\(.name)|\(if .status == "SUCCESS" then "PASS" else "FAIL" end)"' 2>/dev/null || echo "")

    # Build name→test-entity-id pairs (for history queries)
    MAIN_TEST_ID_MAP=$(echo "$MAIN_RAW_JSON" \
      | jq -r '(.testOccurrence // [])[]
          | select(.test.id != null)
          | "\(.name)|\(.test.id)"' 2>/dev/null || echo "")

    if [ -n "$MAIN_TEST_RESULTS" ]; then
      # Extract failed test names from current PR
      PR_FAILED_TESTS=$(echo "$TEST_RESULTS" | awk -F'|' '{if($2=="FAIL") print $1}' | sort)

      # Extract failed test names from main branch
      MAIN_FAILED_TESTS=$(echo "$MAIN_TEST_RESULTS" | awk -F'|' '{if($2=="FAIL") print $1}' | sort)

      # Find tests that failed in PR but not in main
      NEW_FAILURES=$(comm -23 <(echo "$PR_FAILED_TESTS") <(echo "$MAIN_FAILED_TESTS"))

      if [ -n "$NEW_FAILURES" ]; then
        echo "Identified new test failures not in main branch"
      fi
    fi
  fi
fi

CURRENT_TIME_S=$(date +%s)
BUILD_DURATION=$((CURRENT_TIME_S - BUILD_START_TIME))

BUILD_HOURS=$((BUILD_DURATION / 3600))
BUILD_MINUTES=$(((BUILD_DURATION % 3600) / 60))
BUILD_SECONDS=$((BUILD_DURATION % 60))

TEST_HISTORY=""
HAS_NEW_FAILURES="false"
if [ -z "$TEAMCITY_ERROR" ]; then
  echo "Fetching per-test main branch history..."

  ALL_TEST_NAMES=$(echo "$RAW_TEST_RESULTS_JSON" \
    | jq -r '(.testOccurrence // [])[] | select(.status == "SUCCESS" or .status == "FAILURE") | .name' \
    2>/dev/null || echo "")

  while IFS= read -r test_name; do
    [ -z "$test_name" ] && continue

    TEST_ID=$(echo "$RAW_TEST_RESULTS_JSON" \
      | TEST_NAME="$test_name" jq -r '
          (.testOccurrence // [])[]
          | select(.name == env.TEST_NAME)
          | .test.id // empty
        ' 2>/dev/null | head -1)

    PR_STATUS=$(echo "$RAW_TEST_RESULTS_JSON" \
      | TEST_NAME="$test_name" jq -r '
          (.testOccurrence // [])[]
          | select(.name == env.TEST_NAME)
          | .status
        ' 2>/dev/null | head -1)

    FIRST_FAILED_BRANCH=$(echo "$RAW_TEST_RESULTS_JSON" \
      | TEST_NAME="$test_name" jq -r '
          (.testOccurrence // [])[]
          | select(.name == env.TEST_NAME)
          | .firstFailed.build.branchName // ""
        ' 2>/dev/null | head -1)

    IS_NEW="false"
    if [ "$PR_STATUS" = "FAILURE" ] && [ "$FIRST_FAILED_BRANCH" != "refs/heads/main" ] && [ -n "$FIRST_FAILED_BRANCH" ]; then
      IS_NEW="true"
      HAS_NEW_FAILURES="true"
    fi

    if [ -z "$TEST_ID" ]; then
      TEST_HISTORY+="${test_name}|N/A|N/A|N/A|${IS_NEW}"$'\n'
      continue
    fi

    MAIN_HISTORY_JSON=$(curl -s \
      -H "Authorization: Bearer $TEAMCITY_TOKEN" \
      -H "Accept: application/json" \
      "$TEAMCITY_SERVER_URL/app/rest/testOccurrences?locator=test:(id:${TEST_ID}),branch:refs/heads/main,count:1000&fields=testOccurrence(status,build(startDate))")

    CUTOFF_S=$(( CURRENT_TIME_S - 100 * 86400 ))

    PCT="%"
    TC_DATE_FMT="${PCT}Y${PCT}m${PCT}dT${PCT}H${PCT}M${PCT}SZ"
    MAIN_STATS=$(echo "$MAIN_HISTORY_JSON" | CUTOFF_S="$CUTOFF_S" CURRENT_TIME_S="$CURRENT_TIME_S" jq -r --arg fmt "$TC_DATE_FMT" '
      [ .testOccurrence[]?
        | select(
            (.build.startDate // "") != "" and
            (.build.startDate | gsub("\\+(?:[0-9]{4})$"; "Z") | strptime($fmt) | mktime) >= (env.CUTOFF_S | tonumber)
          )
      ] as $recent |
      ($recent | map(select(.status=="SUCCESS")) | length) as $s |
      ($recent | map(select(.status=="FAILURE")) | length) as $f |
      ($s + $f) as $total |
      (if $total > 0 then (($f * 100 / $total | floor | tostring) + "%") else "N/A" end) as $rate |
      ([ .testOccurrence[]?
          | select(.status=="FAILURE" and (.build.startDate // "") != "")
          | (.build.startDate | gsub("\\+(?:[0-9]{4})$"; "Z") | strptime($fmt) | mktime)
        ]) as $fail_ts_list |
      ($fail_ts_list | min // null) as $oldest_fail_ts |
      ($fail_ts_list | max // null) as $newest_fail_ts |
      (if $oldest_fail_ts != null
        then (((env.CURRENT_TIME_S | tonumber) - $oldest_fail_ts) / 86400 | floor | tostring) + "d"
        else ""
      end) as $first_ago |
      (if $newest_fail_ts != null
        then (((env.CURRENT_TIME_S | tonumber) - $newest_fail_ts) / 86400 | floor | tostring) + "d"
        else ""
      end) as $last_ago |
      "\($rate)|\($first_ago)|\($last_ago)"
    ' 2>/dev/null || echo "N/A|")

    # Parse the three pipe-delimited fields from HISTORY_STATS
    HIST_RATE=$(echo "$MAIN_STATS" | cut -d'|' -f1)
    FIRST_AGO=$(echo "$MAIN_STATS" | cut -d'|' -f2)
    LAST_AGO=$(echo "$MAIN_STATS"  | cut -d'|' -f3)

    HIST_RATE_LINK="[$HIST_RATE]($TEAMCITY_SERVER_URL/test/${TEST_ID}?currentProjectId=TF_AzureRM)"

    FIRST_DISPLAY="${FIRST_AGO:-N/A}"
    LAST_DISPLAY="${LAST_AGO:-N/A}"

    TEST_HISTORY+="${test_name}|${HIST_RATE_LINK}|${FIRST_DISPLAY}|${LAST_DISPLAY}|${IS_NEW}"$'\n'
  done <<< "$ALL_TEST_NAMES"
fi

if [ -n "$TEAMCITY_ERROR" ]; then
  COMMENT="Build: [$BUILD_ID]($TEAMCITY_SERVER_URL/viewLog.html?buildId=$BUILD_ID)
PR: #$PR_NUMBER

**TeamCity Error:** $TEAMCITY_ERROR

Unable to collect test details for this run. Please check TeamCity build logs.
"
else
  TABLE_ROWS=$(echo "$TEST_RESULTS" | TEST_HISTORY="$TEST_HISTORY" awk -F'|' '
BEGIN {
    m = split(ENVIRON["TEST_HISTORY"], h_array, "\n")
    for (i = 1; i <= m; i++) {
        if (h_array[i] == "") continue
        line = h_array[i]
        name_part = ""
        pipe_count = 0
        for (j = length(line); j >= 1; j--) {
            if (substr(line, j, 1) == "|") {
                pipe_count++
                if (pipe_count == 4) {
                    name_part = substr(line, 1, j - 1)
                    last4 = substr(line, j + 1)
                    break
                }
            }
        }
        if (name_part != "") hist[name_part] = last4
    }
}
$1 == "" { next }
{
    test_name = $1
    status = $2
    duration = $3

    if (test_name in hist) {
        n_hist = split(hist[test_name], hf, "|")
        fail_rate  = (n_hist >= 1) ? hf[1] : "N/A"
        first_fail = (n_hist >= 2) ? hf[2] : "N/A"
        last_fail  = (n_hist >= 3) ? hf[3] : "N/A"
        is_new     = (n_hist >= 4) ? hf[4] : "true"
    } else {
        fail_rate  = "N/A"
        first_fail = "N/A"
        last_fail  = "N/A"
        is_new     = "true"
    }

    if (status == "PASS") {
      label = "✅ PASS"
    } else if (is_new == "true") {
      label = "❌ NEW"
    } else {
      label = "❌ FAIL"
    }

    print "| " label " | " test_name " | " duration "s | " fail_rate " | " first_fail " | " last_fail " |"
}
')

  COMMENT="Build: [$BUILD_ID]($TEAMCITY_SERVER_URL/viewLog.html?buildId=$BUILD_ID)
PR: #$PR_NUMBER

**Total:** $TOTAL
**Passed:** $PASS_COUNT
**Failed:** $FAIL_COUNT
**Test Duration:** ${BUILD_HOURS}h ${BUILD_MINUTES}m ${BUILD_SECONDS}s

<details>
<summary>Test Details</summary>

| Status | Test Name | Duration | %❌ | First | Last |
| :--- | :--- | ---: | --- | --- | --- |
${TABLE_ROWS}
</details>
"
fi

# Fetch PR author if there are failures
AUTHOR_MESSAGE=""
if [ -z "$TEAMCITY_ERROR" ] && [ "$FAIL_COUNT" -gt 0 ]; then
  PR_AUTHOR=$(github_api_request "/pulls/${PR_NUMBER}" \
  | jq -r '.user.login')

  if [ -z "$PR_AUTHOR" ] || [ "$PR_AUTHOR" = "null" ]; then
    echo "Warning: Could not fetch PR author"
  else
    if [ "$HAS_NEW_FAILURES" = "true" ]; then
      AUTHOR_MESSAGE="@${PR_AUTHOR} - One or more tests newly failed in this PR. Please review the failures.
      "
    else
      AUTHOR_MESSAGE="@${PR_AUTHOR} - One or more tests failed in this PR. Please review the failures.
      "
    fi
  fi
fi

# Add a unique identifier to track comments from this script
# Include tracking ID (hidden in HTML comment) to prevent minimizing current run's comments
COMMENT_IDENTIFIER="<!-- teamcity-test-results -->"

TRACKING_COMMENT=""
if [ "$TRACKING_ID" != "0" ]; then
  TRACKING_COMMENT="<!-- tracking-id:${TRACKING_ID} -->"
fi

BETA_ENV_VAR_NAME="${BETA_VERSION_ENV_VAR#env.}"
BETA_MODE_MESSAGE=""
if [ "${!BETA_ENV_VAR_NAME}" == "true" ]; then
  BETA_MODE_MESSAGE="**Testing in Beta version enabled**
  "
fi

echo "Fetching existing comments..."
COMMENTS_JSON=$(github_api_request "/issues/${PR_NUMBER}/comments")

# Filter comments that should be minimized (teamcity-test-results or /test comments)
# but exclude those with the current tracking ID
COMMENT_IDS=$(echo "$COMMENTS_JSON" | jq -r --arg tracking_id "$TRACKING_ID" '
  .[] |
  select(.body | type == "string" and (contains("<!-- teamcity-test-results -->") or startswith("/test"))) |
  select(.body | contains("tracking-id:" + $tracking_id) | not) |
  .node_id
' 2>&1 | grep -v "^jq:")

if [ -n "$COMMENT_IDS" ]; then
  echo "Found previous comments to minimize"
  while IFS= read -r COMMENT_NODE_ID; do
    if [ -n "$COMMENT_NODE_ID" ]; then
      echo "Minimizing comment: $COMMENT_NODE_ID"
      RESPONSE=$(curl -s -X POST \
        -H "Authorization: bearer $GIT_PAT" \
        -H "Content-Type: application/json" \
        https://api.github.com/graphql \
        -d "{\"query\": \"mutation { minimizeComment(input: {subjectId: \\\"$COMMENT_NODE_ID\\\", classifier: OUTDATED}) { minimizedComment { isMinimized } } }\"}")

      # Check if minimization was successful
      if echo "$RESPONSE" | jq -e '.data.minimizeComment.minimizedComment.isMinimized' > /dev/null 2>&1; then
        echo "Successfully minimized comment: $COMMENT_NODE_ID"
      else
        echo "Warning: Failed to minimize comment: $COMMENT_NODE_ID"
        echo "Response: $RESPONSE"
      fi
    fi
  done <<< "$COMMENT_IDS"
else
  echo "No previous comments found to minimize"
fi

COMMENT="${COMMENT_IDENTIFIER}
${TRACKING_COMMENT}
${AUTHOR_MESSAGE}
${BETA_MODE_MESSAGE}
${COMMENT}"

echo "Posting new comment to GitHub..."

curl -s -X POST \
  -H "Authorization: Bearer $GIT_PAT" \
  -H "Accept: application/vnd.github+json" \
  -H "Content-Type: application/json" \
  "https://api.github.com/repos/$GITHUB_REPO/issues/${PR_NUMBER}/comments" \
  -d "{\"body\": $(jq -Rs . <<< "$COMMENT")}"

if [ "$APPLY_TESTING_LABELS_ENABLED" = "true" ]; then
  echo "Applying labels..."

  if [ -n "$TEAMCITY_ERROR" ]; then
    echo "Skipping label application due to TeamCity error"
    exit 0
  fi

  # If no failures, apply teamcity-passed label
  if [ "$FAIL_COUNT" -eq 0 ]; then
    echo "No test failures detected"
    set_testing_label "$LABEL_SUCCESS"
    exit 0
  fi

  # If there are failures, determine label based on earlier analysis
  if [ -z "$NEW_FAILURES" ]; then
    echo "All failed tests also exist in main branch"
    set_testing_label "$LABEL_FAILURE"
  else
    echo "Found new test failures not present in main branch"
    set_testing_label "$LABEL_NEW_FAILURE"
  fi

  echo "Label application complete"
fi

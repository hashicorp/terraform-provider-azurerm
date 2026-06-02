#!/bin/bash

POST_GITHUB_COMMENT="%POST_GITHUB_COMMENT%"
GITHUB_REPO="%env.GITHUB_REPO%"
GIT_PAT="%env.GIT_PAT%"
TEAMCITY_TOKEN="%env.TEAMCITY_TOKEN%"
# BUILD_START_TIME is set in the first build step: $(date +%s)
BUILD_START_TIME="%env.BUILD_START_TIME%"
BETA_VERSION_ENV_VAR="%env.BETA_VERSION_ENV_VAR%"
TEAMCITY_BUILD_BRANCH="%teamcity.build.branch%"

if [ "$POST_GITHUB_COMMENT" != "true" ]; then
  echo "GitHub commenting disabled — skipping."
  exit 0
fi

if [[ "$TEAMCITY_BUILD_BRANCH" =~ refs/pull/([0-9]+)/merge ]]; then
  PR_NUMBER="${BASH_REMATCH[1]}"
else
  echo "Not a PR merge branch: %teamcity.build.branch%"
  exit 0
fi

TRACKING_ID="%TRACKING_ID%"
echo "Tracking ID: $TRACKING_ID"

BUILD_ID="%teamcity.build.id%"

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
  if [ "$label" = "testing-succeeded" ]; then
    remove_label "testing-failed"
    apply_label "testing-succeeded"
  elif [ "$label" = "testing-failed" ]; then
    remove_label "testing-succeeded"
    apply_label "testing-failed"
  fi
}

TEST_RESULTS=$(file="results.txt"

awk '
# Parse TeamCity testStdOut messages
/##teamcity\[testStdOut/ {
    line = $0

    # Extract test name between name= and the next single quote
    if (match(line, /name=.([^'"'"']+)./)) {
        test_name = substr(line, RSTART+6, RLENGTH-7)
    }

    # Extract output content between out= and the closing bracket
    if (match(line, /out=.([^'"'"']+)./)) {
        output = substr(line, RSTART+5, RLENGTH-6)

        # Replace |n with newlines
        gsub(/\|n/, "\n", output)

        # Extract status (PASS or FAIL)
        status = ""
        if (match(output, /--- PASS:/)) {
            status = "PASS"
        } else if (match(output, /--- FAIL:/)) {
            status = "FAIL"
        }

        # Extract duration
        duration = ""
        if (match(output, /\(([0-9.]+)s\)/)) {
            duration_str = substr(output, RSTART, RLENGTH)
            gsub(/[()s]/, "", duration_str)
            duration = duration_str
        }

        # Print result
        if (test_name != "" && status != "" && duration != "") {
            print test_name "|" status "|" duration "|" output
        }
    }
}
' "$file"
)

PASS_COUNT=$(echo "$TEST_RESULTS" | awk -F'|' '{if($2=="PASS") print}' | wc -l | tr -d ' ')
FAIL_COUNT=$(echo "$TEST_RESULTS" | awk -F'|' '{if($2=="FAIL") print}' | wc -l | tr -d ' ')
TOTAL=$((PASS_COUNT + FAIL_COUNT))

# Fetch main branch test results early to identify new failures for comment marking
NEW_FAILURES=""
MAIN_TEST_RESULTS=""
if [ "$FAIL_COUNT" -gt 0 ]; then
  echo "Fetching main branch test results for comparison..."

  # Get the latest successful build on main branch
  MAIN_BUILD_INFO=$(curl -s \
    -H "Authorization: Bearer $TEAMCITY_TOKEN" \
    -H "Accept: application/json" \
    "%teamcity.serverUrl%/app/rest/builds?locator=buildType:(id:%system.teamcity.buildType.id%),branch:refs/heads/main,status:SUCCESS,count:1")

  MAIN_BUILD_ID=$(echo "$MAIN_BUILD_INFO" | jq -r '.build[0].id // empty')

  if [ -n "$MAIN_BUILD_ID" ]; then
    echo "Found main branch build: $MAIN_BUILD_ID"

    # Download test results from main branch build
    MAIN_RESULTS_URL="%teamcity.serverUrl%/app/rest/builds/id:$MAIN_BUILD_ID/artifacts/content/results.txt"
    MAIN_TEST_RESULTS=$(curl -s \
      -H "Authorization: Bearer $TEAMCITY_TOKEN" \
      "$MAIN_RESULTS_URL" 2>/dev/null || echo "")

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

COMMENT="Build: [$BUILD_ID](%teamcity.serverUrl%/viewLog.html?buildId=$BUILD_ID)
PR: #$PR_NUMBER

**Total:** $TOTAL
**Passed:** $PASS_COUNT
**Failed:** $FAIL_COUNT
**Test Duration:** ${BUILD_HOURS}h ${BUILD_MINUTES}m ${BUILD_SECONDS}s

<details>
<summary>Test Details</summary>

<table>
<tr><td><b>Status</b></td><td><b>Test Name</b></td><td><b>Duration</b></td></tr>
"

TABLE_ROWS=$(echo "$TEST_RESULTS" | awk -v RS='\nTestAcc' -v new_failures="$NEW_FAILURES" '
BEGIN {
    # Build array of new failure test names
    split(new_failures, nf_array, "\n")
    for (i in nf_array) {
        new_fail[nf_array[i]] = 1
    }
}
NR==1 && /^TestAcc/ {
    record = $0
}
NR>1 {
    record = "TestAcc" $0
}
record != "" {
    # Find positions of first 3 pipes
    pipe1 = index(record, "|")
    pipe2 = index(substr(record, pipe1+1), "|") + pipe1
    pipe3 = index(substr(record, pipe2+1), "|") + pipe2

    test_name = substr(record, 1, pipe1-1)
    status = substr(record, pipe1+1, pipe2-pipe1-1)
    duration = substr(record, pipe2+1, pipe3-pipe2-1)
    output = substr(record, pipe3+1)

    if (status == "PASS") {
        print "<tr><td>✅ PASS</td><td>" test_name "</td><td>" duration "s</td></tr>"
    } else if (status == "FAIL") {
        # Check if this is a new failure
        if (test_name in new_fail) {
            print "<tr><td>❌ FAIL 🆕</td><td>" test_name "</td><td>" duration "s</td></tr>"
        } else {
            print "<tr><td>❌ FAIL</td><td>" test_name "</td><td>" duration "s</td></tr>"
        }
    }
}
')

COMMENT+="$TABLE_ROWS"

COMMENT+="</table>
</details>
"

# Fetch PR author if there are failures
AUTHOR_MESSAGE=""
if [ "$FAIL_COUNT" -gt 0 ]; then
  PR_AUTHOR=$(github_api_request "/pulls/${PR_NUMBER}" \
  | jq -r '.user.login')

  if [ -z "$PR_AUTHOR" ] || [ "$PR_AUTHOR" = "null" ]; then
    echo "Warning: Could not fetch PR author"
  else
    if [ -z "$NEW_FAILURES" ]; then
      AUTHOR_MESSAGE="@${PR_AUTHOR} - One or more tests failed in this PR. Please review the failures.
      "
    else
      AUTHOR_MESSAGE="@${PR_AUTHOR} - One or more tests newly failed in this PR. Please review the failures.
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

echo "Applying labels..."

# If no failures, apply testing-succeeded label
if [ "$FAIL_COUNT" -eq 0 ]; then
  echo "No test failures detected"
  set_testing_label "testing-succeeded"
  exit 0
fi

# If there are failures, determine label based on earlier analysis
if [ -z "$MAIN_TEST_RESULTS" ]; then
  echo "Could not fetch main branch results - applying 'testing-failed' label as precaution..."
  set_testing_label "testing-failed"
elif [ -z "$NEW_FAILURES" ]; then
  echo "All failed tests also exist in main branch"
  set_testing_label "testing-succeeded"
else
  echo "Found new test failures not present in main branch"
  set_testing_label "testing-failed"
fi

echo "Label application complete"

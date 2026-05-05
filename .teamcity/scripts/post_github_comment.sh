#!/bin/bash

POST_GITHUB_COMMENT="%POST_GITHUB_COMMENT%"
POST_GITHUB_COMMENT_DETAILED="%POST_GITHUB_COMMENT_DETAILED%"

if [ "$POST_GITHUB_COMMENT" != "true" ] && [ "$POST_GITHUB_COMMENT_DETAILED" != "true" ]; then
  echo "GitHub commenting disabled — skipping."
  exit 0
fi

TEAMCITY_BUILD_BRANCH="%teamcity.build.branch%"

if [[ "$TEAMCITY_BUILD_BRANCH" =~ refs/pull/([0-9]+)/merge ]]; then
  PR_NUMBER="${BASH_REMATCH[1]}"
else
  echo "Not a PR merge branch: %teamcity.build.branch%"
  exit 0
fi

detailed=false
#if [ "$POST_GITHUB_COMMENT_DETAILED" = "true" ]; then
#  echo "Detailed GitHub commenting enabled."
#  detailed=true
#fi

BUILD_ID="%teamcity.build.id%"

github_api_request() {
  local endpoint="$1"
  curl -s \
    -H "Authorization: Bearer %env.GIT_PAT%" \
    -H "Accept: application/vnd.github+json" \
    "https://api.github.com/repos/%env.GITHUB_REPO%${endpoint}"
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

# Parse results to get counts
PASS_COUNT=$(echo "$TEST_RESULTS" | awk -F'|' '{if($2=="PASS") print}' | wc -l | tr -d ' ')
FAIL_COUNT=$(echo "$TEST_RESULTS" | awk -F'|' '{if($2=="FAIL") print}' | wc -l | tr -d ' ')
TOTAL=$((PASS_COUNT + FAIL_COUNT))

# Calculate build duration from start time to now
# BUILD_START_TIME should is set in the first build step: $(date +%s)
BUILD_START_TIME="%env.BUILD_START_TIME%"
CURRENT_TIME_S=$(date +%s)
BUILD_DURATION=$((CURRENT_TIME_S - BUILD_START_TIME))

# Convert to hours, minutes and seconds for better readability
BUILD_HOURS=$((BUILD_DURATION / 3600))
BUILD_MINUTES=$(((BUILD_DURATION % 3600) / 60))
BUILD_SECONDS=$((BUILD_DURATION % 60))


# Fetch PR author if there are failures
PREFIX=""
if [ "$FAIL_COUNT" -gt 0 ]; then
  PR_AUTHOR=$(github_api_request "/pulls/${PR_NUMBER}" \
  | jq -r '.user.login')

  if [ -z "$PR_AUTHOR" ] || [ "$PR_AUTHOR" = "null" ]; then
    echo "Warning: Could not fetch PR author"
  else
    PREFIX="@${PR_AUTHOR} - One or more tests failed in this PR. Please review the failures.

"
  fi
fi

# Build GitHub comment
COMMENT="${PREFIX}Build: [$BUILD_ID](%teamcity.serverUrl%/viewLog.html?buildId=$BUILD_ID)
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

TABLE_ROWS=$(echo "$TEST_RESULTS" | awk -v RS='\nTestAcc' -v detailed="$detailed" '
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
        if (detailed == "true") {
            # Escape HTML special chars
            gsub(/&/, "\\&", output)
            gsub(/</, "\\<", output)
            gsub(/>/, "\\>", output)
            print "<tr><td>❌ FAIL</td><td>" test_name "<br/><details><summary>Error Details</summary><pre><code>" output "</code></pre></details></td><td>" duration "s</td></tr>"
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

# Add a unique identifier to track comments from this script
COMMENT_IDENTIFIER="<!-- teamcity-test-results -->"
COMMENT="${COMMENT_IDENTIFIER}
${COMMENT}"

echo "Minimizing previous comments from this run..."
# Fetch existing comments on the PR
echo "Fetching existing comments..."
COMMENT_IDS=$(github_api_request "/issues/${PR_NUMBER}/comments" \
  | jq -r '.[] | select(.body | type == "string" and (contains("<!-- teamcity-test-results -->") or startswith("/test"))) | .node_id' 2>&1 | grep -v "^jq:")


# Minimize each previous comment using GitHub's GraphQL API
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

echo "Posting new comment to GitHub..."

curl -s -X POST \
  -H "Authorization: Bearer %env.GIT_PAT%" \
  -H "Accept: application/vnd.github+json" \
  -H "Content-Type: application/json" \
  "https://api.github.com/repos/%env.GITHUB_REPO%/issues/${PR_NUMBER}/comments" \
  -d "{\"body\": $(jq -Rs . <<< "$COMMENT")}"

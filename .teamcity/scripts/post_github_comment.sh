#!/bin/bash

if [ "%POST_GITHUB_COMMENT%" != "true" ] && [ "%POST_GITHUB_COMMENT_DETAILED%" != "true" ]; then
  echo "GitHub commenting disabled — skipping."
  exit 0
fi

if [[ "%teamcity.build.branch%" =~ refs/pull/([0-9]+)/merge ]]; then
  PR_NUMBER="${BASH_REMATCH[1]}"
else
  echo "Not a PR merge branch: %teamcity.build.branch%"
  exit 0
fi

detailed=false
if [ "%POST_GITHUB_COMMENT_DETAILED%" = "true" ]; then
  echo "Detailed GitHub commenting enabled."
  detailed=true
fi

BUILD_ID="%teamcity.build.id%"

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


# Build GitHub comment
COMMENT="Build: [$BUILD_ID](%teamcity.serverUrl%/viewLog.html?buildId=$BUILD_ID)
PR: #$PR_NUMBER

**Total:** $TOTAL
**Passed:** $PASS_COUNT
**Failed:** $FAIL_COUNT

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

echo "Posting comment to GitHub..."

curl -s \
-H "Authorization: Bearer %env.GIT_PAT%" \
-H "Accept: application/vnd.github+json" \
https://api.github.com/repos/%env.GITHUB_REPO%/issues/${PR_NUMBER}/comments \
-d "{\"body\": $(jq -Rs . <<< "$COMMENT")}"

echo "Comment posted successfully."

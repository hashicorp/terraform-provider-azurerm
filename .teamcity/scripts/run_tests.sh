#!/bin/bash

# Runs the acceptance tests for one service, first saying what state this agent's Go build cache is in
# so that the build log shows whether the tests are compiled from scratch or from what earlier builds
# left behind.
#
# Each agent keeps its own cache between builds (see GoCache() in build_components.kt) and nothing is
# fetched: the first build on an agent compiles everything, later builds only what has changed.

SERVICE_PATH="%SERVICE_PATH%"

GO_CACHE="$(go env GOCACHE)"
# kept beside the cache rather than in it, so that it isn't emptied or published along with it
LAST_USED_FILE="$(dirname "$GO_CACHE")/last-used"

# nothing in here fails the build
describe_go_cache() {
  local size="not there yet"
  if [ -d "$GO_CACHE" ]; then
    # given up on rather than holding up the tests, should the cache have grown very large
    size="$(timeout 20 du -sh "$GO_CACHE" 2>/dev/null | awk '{print $1}')"
  fi
  echo "$(go env GOVERSION), build cache at $GO_CACHE (${size:-size unknown})."

  local last_used_at last_used_by
  if [ -f "$LAST_USED_FILE" ] && read -r last_used_at last_used_by < "$LAST_USED_FILE" && [[ "$last_used_at" =~ ^[0-9]+$ ]]; then
    local minutes=$(( ($(date +%s) - last_used_at) / 60 ))
    if [ "$minutes" -lt 120 ]; then
      echo "Last used on this agent $minutes minutes ago, by $last_used_by."
    else
      echo "Last used on this agent $((minutes / 60)) hours ago, by $last_used_by."
    fi
  else
    echo "Not used on this agent before."
  fi
  mkdir -p "$(dirname "$LAST_USED_FILE")"
  echo "$(date +%s) ${TEAMCITY_BUILDCONF_NAME:-an unnamed build} #${BUILD_NUMBER:-?}" > "$LAST_USED_FILE"

  # asks Go which of the packages the tests are built from it already has compiled, which compiles
  # nothing itself ("unsafe" is left out as it's never compiled)
  local packages
  if packages="$(go list -deps -f '{{if ne .ImportPath "unsafe"}}{{.Stale}}{{end}}' "$SERVICE_PATH/..." 2>/dev/null)"; then
    local total compiled
    total="$(grep -c . <<< "$packages")"
    compiled="$(grep -c '^false$' <<< "$packages")"
    if [ "$compiled" -eq "$total" ]; then
      echo "All $total packages the tests are built from are already compiled."
    else
      echo "$compiled of the $total packages the tests are built from are already compiled, the rest are compiled now."
    fi
  fi
}

describe_go_cache

# -count=1 has the tests run every time: the cache outlives a build, and without it Go answers from the
# cache rather than run tests which it has already seen pass
#
# the reporter is told when `go test` started so that it can say how long the tests took to compile
go test -v "$SERVICE_PATH/..." -count=1 -timeout="%TIMEOUT%h" -test.parallel="%PARALLELISM%" -run="%TEST_PREFIX%" -json | go run ./internal/tools/teamcity-test-reporter -started-at="$(date +%s)"

#!/bin/bash

# Says what state this agent's Go build cache is in, so that the build log shows whether the tests are
# compiled from scratch or from what earlier builds left behind, and drops what nothing has used lately.
#
# Each agent keeps its own cache between builds (see GoCache() in build_components.kt) and nothing is
# fetched: the first build on an agent compiles everything, later builds only what has changed.
#
# Nothing in here fails the build.

SERVICE_PATH="%SERVICE_PATH%"

GO_CACHE="$(go env GOCACHE)"
# compiled packages which no build has used for this many days are deleted
RETENTION_DAYS=2
# kept beside the cache rather than in it, so that it isn't emptied or published along with it
LAST_USED_FILE="$(dirname "$GO_CACHE")/last-used"

size="not there yet"
if [ -d "$GO_CACHE" ]; then
  # Go itself only drops entries unused for five days, which is longer than most agents live. Entries are
  # files or directories named <hash>-a and <hash>-d, and Go rebuilds anything it finds missing.
  find "$GO_CACHE" -mindepth 2 -maxdepth 2 \( -name '*-a' -o -name '*-d' \) -mmin +"$((RETENTION_DAYS * 24 * 60))" -exec rm -rf {} + 2>/dev/null
  # given up on rather than holding up the build, should the cache have grown very large
  size="$(timeout 20 du -sh "$GO_CACHE" 2>/dev/null | awk '{print $1}')"
fi
echo "$(go env GOVERSION), build cache at $GO_CACHE (${size:-size unknown})."

if [ -f "$LAST_USED_FILE" ] && read -r last_used_at last_used_by < "$LAST_USED_FILE" && [[ "$last_used_at" =~ ^[0-9]+$ ]]; then
  minutes=$(( ($(date +%s) - last_used_at) / 60 ))
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
# nothing itself. "unsafe" is never compiled, and the test binaries ("foo.test" and the "[foo.test]"
# variants of packages) are never cached, so they're left out.
if packages="$(go list -deps -test -f '{{.Stale}} {{.ImportPath}}' "$SERVICE_PATH/..." 2>/dev/null | grep -v -e ' unsafe$' -e '\.test$' -e '\.test\]$')"; then
  total="$(grep -c . <<< "$packages")"
  compiled="$(grep -c '^false ' <<< "$packages")"
  if [ "$compiled" -eq "$total" ]; then
    echo "All $total packages the tests are built from are already compiled."
  else
    echo "$compiled of the $total packages the tests are built from are already compiled, the rest are compiled by the next step."
  fi
fi

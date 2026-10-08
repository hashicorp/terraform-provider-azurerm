#!/bin/bash

# Runs the acceptance tests for one service.

SERVICE_PATH="%SERVICE_PATH%"

# -count=1 has the tests run every time: the Go cache outlives a build (see go_cache.sh), and without it Go
# answers from the cache rather than run tests which it has already seen pass
#
# the reporter is told when `go test` started so that it can say how long the tests took to compile
go test -v "$SERVICE_PATH/..." -count=1 -timeout="%TIMEOUT%h" -test.parallel="%PARALLELISM%" -run="%TEST_PREFIX%" -json | go run ./internal/tools/teamcity-test-reporter -started-at="$(date +%s)"

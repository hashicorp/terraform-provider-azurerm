#!/usr/bin/env bash

function checkForConditionalRun {
  if [ "$TRAVIS" == "ci" ];
  then
    echo "Checking if this should be conditionally run.."
    result=$(git diff --name-only origin/master | grep azurerm/)
    if [ "$result" = "" ];
    then
      echo "No changes committed to ./azurerm - nothing to lint - exiting"
      exit 0
    fi
  fi
}

function runLinters {
  echo "==> Checking source code against linters..."
  (while true; do sleep 300; echo "(I'm still alive and linting!)"; done) & PID="$!"; echo "Watcher subprocess: $PID"; \
  golangci-lint run ./... -v --concurrency 1 --config .golangci-travis.yml ; ES="$?"; kill -9 "$PID"; exit "$ES"
}

function main {
  checkForConditionalRun
  runLinters
}

main

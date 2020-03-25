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

function runTests {
  echo "==> Running Unit Tests..."
  go test -i $TEST || exit 1
  go test -v $TEST "$TESTARGS" -timeout=30s -parallel=4
}

function main {
  checkForConditionalRun
  runTests
}

main

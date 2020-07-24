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
  golangci-lint run ./...
}

function main {
  checkForConditionalRun
  runLinters
}

main

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

function runDeprecatedFunctions {
  echo "==> Checking for use of gradually deprecated functions..."
  result=$(git diff --name-only origin/master | grep azurerm/)
  if [ "$result" = "" ];
  then
    # No changes committed to ./azurerm - so nothing to check
    exit 0
  fi

  # require resources to be imported is now hard-coded on - but only checking for additions
  result=$(git diff origin/master | grep + | grep -R "features\.ShouldResourcesBeImported")
  if [ "$result" = "" ];
  then
    echo "The Feature Flag for 'ShouldResourcesBeImported' will be deprecated in the future"
    echo "and shouldn't be used in new resources - please remove new usages of the"
    echo "'ShouldResourcesBeImported' function from these changes - since this is now enabled"
    echo "by default."
    echo ""
    echo "In the future this function will be marked as Deprecated - however it's not for"
    echo "the moment to not conflict with open Pull Requests."
  fi
}

function main {
  checkForConditionalRun
  runDeprecatedFunctions
}

main

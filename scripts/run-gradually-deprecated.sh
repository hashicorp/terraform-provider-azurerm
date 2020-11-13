#!/usr/bin/env bash

function runGraduallyDeprecatedFunctions {
  echo "==> Checking for use of gradually deprecated functions..."
  
  # require resources to be imported is now hard-coded on - but only checking for additions
  result=$(git diff origin/master | grep + | grep -R "features\.ShouldResourcesBeImported")
  if [ "$result" != "" ];
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

function runDeprecatedFunctions {
  echo "==> Checking for use of deprecated functions..."
  result=$(grep -Ril "d.setid(\"\")" ./azurerm/internal/services/**/data_source_*.go)
  if [ "$result" != "" ];
  then
    echo "Data Sources should return an error when a resource cannot be found rather than"
    echo "setting an empty ID (by calling 'd.SetId("")'."
    echo ""
    echo "Please remove the references to 'd.SetId("") from the Data Sources listed below"
    echo "and raise an error instead:"
    echo ""
    echo $result
  fi
}

function main {
  runGraduallyDeprecatedFunctions
  runDeprecatedFunctions
}

main

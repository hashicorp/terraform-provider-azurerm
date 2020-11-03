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
  echo "==> Checking source code against terraform provider linters..."
	tfproviderlint \
        -AT001\
        -AT001.ignored-filename-suffixes _data_source_test.go\
        -AT005 -AT006 -AT007\
        -R001 -R002 -R003 -R004 -R006\
        -S001 -S002 -S003 -S004 -S005 -S006 -S007 -S008 -S009 -S010 -S011 -S012 -S013 -S014 -S015 -S016 -S017 -S018 -S019 -S020\
        -S021 -S022 -S023 -S024 -S025 -S026 -S027 -S028 -S029 -S030 -S031 -S032 -S033\
        ./azurerm/...
	bash ./scripts/terrafmt-acctests.sh
}

function main {
  checkForConditionalRun
  runTests
}

main

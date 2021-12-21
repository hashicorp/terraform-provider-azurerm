#!/usr/bin/env bash

function compileAndUnitTests {
  local servicePackage=$1

  echo "Running Unit Tests ('go test ./internal/services/$servicePackage/...').."
  go test "./internal/services/$servicePackage/..."
  if (($? > 1)); then
      echo "Build Failed."
      exit 1
  fi
}

function runLinters {
  local servicePackage=$1

  echo "Running Go Linting ('golangci-lint')"
  golangci-lint run -v "./internal/services/$servicePackage/..."
  if (($? > 1)); then
      echo "Go Linting Failed."
      exit 1
  fi

  echo "Running Terraform Provider Linting ('tfproviderlint')"
  tfproviderlint \
          -AT001\
          -AT001.ignored-filename-suffixes _data_source_test.go\
          -AT005 -AT006 -AT007\
          -R001 -R002 -R003 -R004 -R006 -R006.package-aliases pluginsdk\
          -S001 -S002 -S003 -S004 -S005 -S006 -S007 -S008 -S009 -S010 -S011 -S012 -S013 -S014 -S015 -S016 -S017 -S018 -S019 -S020\
          -S021 -S022 -S023 -S024 -S025 -S026 -S027 -S028 -S029 -S030 -S031 -S032 -S033\
          "./internal/services/$servicePackage/..."
  if (($? > 1)); then
      echo "Terraform Provider Linting Failed."
      exit 1
  fi
}

function main {
  local servicePackage=$1

  compileAndUnitTests "$servicePackage"
  if (($? > 1)); then
      exit 1
  fi

  runLinters "$servicePackage"
  if (($? > 1)); then
      exit 1
  fi
}

main $1



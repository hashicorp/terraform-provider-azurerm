#!/usr/bin/env sh
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


echo "==> Checking that acceptance test packages are used..."

invalid_files=$(find ./internal/services -maxdepth 2 -type f -regex '.*\(resource\|data_source\).*_test\.go$' -print0 | \
  xargs -0 awk '
    FNR > 5 {
      nextfile  # Skip to next file after processing first 5 lines
    }
    $1 == "package" && $0 !~ /_test$/ {
      print FILENAME
      nextfile  # Skip to next file once we find an invalid package
    }
  ')

if [ -n "$invalid_files" ]; then
  echo "$invalid_files"
  echo ""
  echo "------------------------------------------------"
  echo ""
  echo "The acceptance test files listed above are using the same package as the resource or data source code."
  echo "They must use a test package to prevent a circular dependency. To fix this change the first line:"
  echo ""
  echo "> package service"
  echo ""
  echo "to"
  echo ""
  echo "> package service_test"
  echo ""
  exit 1
fi

exit 0

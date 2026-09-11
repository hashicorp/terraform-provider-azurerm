#!/usr/bin/env bash
# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0


echo "==> Checking acceptance test terraform blocks are formatted..."

if ! terrafmt diff -c -q -f -p "*_test.go" ./internal; then
  echo "------------------------------------------------"
  echo ""
  echo "The preceding files contain terraform blocks that are not correctly formatted or contain errors."
  echo "You can fix this by running make tools and then terrafmt on them."
  echo ""
  echo "to easily fix all terraform blocks:"
  echo "$ make terrafmt"
  echo ""
  echo "format only acceptance test config blocks:"
  echo "$ terrafmt fmt -f -p \"*_test.go\" ./internal"
  echo ""
  echo "format a single test file:"
  echo "$ terrafmt fmt -f ./internal/services/service/tests/resource_test.go"
  echo ""
  echo "on windows:"
  echo "$ terrafmt fmt -f -p \"*_test.go\" ./internal"
  echo ""
  exit 1
fi

exit 0

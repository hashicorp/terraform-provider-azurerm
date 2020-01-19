#!/usr/bin/env bash

echo "==> Checking acceptance test terraform blocks are formatted..."

files=$(find ./azurerm/internal -type f -name "*_test.go")
error=false

for f in $files; do
  terrafmt diff -c -q -f "$f" || error=true
done

if ${error}; then
  echo "------------------------------------------------"
  echo ""
  echo "The preceding files contain terraform blocks that are not correctly formatted or contain errors."
  echo "You can fix this by running make tools and then terrafmt on them."
  echo ""
  echo "format a single file:"
  echo "$ terrafmt fmt ./azurerm/internal/services/service/tests/resource_test.go"
  echo ""
  echo "format all website files:"
  echo "$ find azurerm | egrep \"_test.go\" | sort | while read f; do terrafmt fmt -f \$f; done"
  echo ""
  echo "on windows:"
  echo "$ Get-ChildItem -Path . -Recurse -Filter \"_test.go\" | foreach {terrafmt fmt -f $.name}"
  echo ""
  exit 1
fi

exit 0

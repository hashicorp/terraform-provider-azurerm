#!/usr/bin/env bash

echo "==> Checking documentation terraform blocks are formatted..."

files=$(find ./website -type f -name "*.html.markdown")
error=false

for f in $files; do
  terrafmt diff -c -q "$f" || error=true
done

if ${error}; then
  echo "------------------------------------------------"
  echo ""
  echo "The preceding files contain terraform blocks that are not correctly formatted or contain errors."
  echo "You can fix this by running make tools and then terrafmt on them."
  echo ""
  echo "to easily fix all terraform blocks:"
  echo "$ make terrafmt"
  echo ""
  echo "format only website config blocks:"
  echo "$ find . | egrep html.markdown | sort | while read f; do terrafmt fmt \$f; done"
  echo ""
  echo "format a single website file:"
  echo "$ terrafmt fmt ./website/path/to/file.html.markdown"
  echo ""
  echo "on windows:"
  echo "$ Get-ChildItem -Path . -Recurse -Filter \"*html.markdown\" | foreach {terrafmt fmt $_.fullName}"
  echo ""
  exit 1
fi

exit 0

#!/usr/bin/env bash

echo "==> Checking website Terraform blocks are formatted..."

files=$(find ./website -type f -name "*.html.markdown")
tofmt=()
error=false

for f in $files; do
  terrafmt diff -c -q "$f" || error=true
done

echo $tofmt

if $error; then
  echo "------------------------------------------------"
  echo ""
  echo "The preceding files contain terraform blocks that are not correctly formatted:"
  echo "You can fix this by running terrafmt on them:"
  echo ""
  echo "format a single file:"
  echo "$ terrafmt fmt ./website/path/to/file.html.markdown"
  echo ""
  echo "format all website files:"
  echo "$ find . | egrep html.markdown | sort | while read f; do terrafmt fmt \$f; done"
  echo ""
  echo ""
  exit 1
fi

exit 0

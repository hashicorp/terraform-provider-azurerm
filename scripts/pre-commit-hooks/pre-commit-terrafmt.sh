#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

start_time=$(date +%s.%N)

echo "==> ⌛ pre-commit-terrafmt.sh..."

if ! command -v terrafmt &> /dev/null; then
    echo "Error: terrafmt command not found."
    echo "Please run 'make tools' first to install required dependencies."
    exit 1
fi

staged_test_files=$(git diff --cached --name-only --diff-filter=ACM | grep "_test\.go$" || true)
staged_doc_files=$(git diff --cached --name-only --diff-filter=ACM | grep "\.html\.markdown$" || true)

error=false

if [ -n "$staged_test_files" ]; then
  echo "Checking staged test files..."
  for f in $staged_test_files; do
    if [ -f "$f" ]; then
      terrafmt diff -c -q -f "$f" || error=true
    fi
  done
fi

if [ -n "$staged_doc_files" ]; then
  echo "Checking staged documentation files..."
  for f in $staged_doc_files; do
    if [ -f "$f" ]; then
      terrafmt diff -c -q "$f" || error=true
    fi
  done
fi

end_time=$(date +%s.%N)
execution_time=$(echo "$end_time - $start_time" | bc -l)

if ${error}; then
  echo "------------------------------------------------"
  echo ""
  echo "❌ The preceding staged files contain terraform blocks that are not correctly formatted or contain errors."
  echo "You can fix this by running one of the below commands:"
  echo ""
  echo "format staged test files:"
  echo "$ git diff --cached --name-only --diff-filter=ACM | grep '_test\.go$' | while read f; do terrafmt fmt -f \"\$f\"; done"
  echo ""
  echo "format staged documentation files:"
  echo "$ git diff --cached --name-only --diff-filter=ACM | grep '\.html\.markdown$' | while read f; do terrafmt fmt \"\$f\"; done"
  echo ""
  echo "format a single file:"
  echo "$ terrafmt fmt -f ./path/to/file"
  echo ""
  echo "After fixing the formatting, you'll need to stage the changes again:"
  echo "$ git add <fixed-files>"
  echo ""
  printf "Script execution time: %.3f seconds\n" "$execution_time"
  
  exit 1
fi

if [ -z "$staged_test_files" ] && [ -z "$staged_doc_files" ]; then
  echo "✅ No _test.go or .html.markdown files staged for commit."
else
  echo "✅ All staged _test.go and .html.markdown files are correctly formatted."
fi

printf "Script execution time: %.3f seconds\n" "$execution_time"
exit 0

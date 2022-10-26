#!/usr/bin/env bash

echo "==> Checking documentation Markdown content are formatted..."

files=$(find ./website -type f -name "*.html.markdown")
error=false

for f in $files; do
    markdownlint "$f" --disable MD013 || error=true
done

if ${error}; then
    echo "------------------------------------------------"
    echo ""
    echo "The preceding files contain Markdown content that are not correctly formatted or contain errors."
    echo "You can fix this by running make tools and then markdownlint on them."
    echo ""
    echo "format a single website file:"
    echo "$ markdownlint --disable MD013 --fix ./website/path/to/file.html.markdown"
    exit 1
fi

exit 0

#!/usr/bin/env bash

docs=$(ls website/docs/**/*.markdown)
error=false

for doc in $docs; do
  dirname=$(dirname "$doc")
  category=$(basename "$dirname")

  case "$category" in
    "guides")
      # Guides require a page_title
      if ! grep "^page_title: " "$doc" > /dev/null; then
        echo "Website guide file is missing a 'page_title' line at the top: $doc"
        error=true
      fi
      ;;

    "d" | "r")
      # Resources and data sources require a subcategory
      if ! grep "^subcategory: " "$doc" > /dev/null; then
        echo "Website documentation file is missing a 'subcategory' line at the top: $doc"
        error=true
      fi
      ;;

    *)
      error=true
      echo "Unknown category \"$category\". " \
        "Docs can only exist in r/, d/, or guides/ folders."
      ;;
  esac
done

if $error; then
  exit 1
fi

exit 0

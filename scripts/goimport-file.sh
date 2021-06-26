#!/bin/bash

# remove blank lines in go imports then run goimports

if [ $# != 1 ] ; then
  echo "usage: $0 <filename>"
  exit 1
fi

sed -i '' -e '
  /^import/,/)/ {
    /^$/d
  }
' $1

goimports -w $1

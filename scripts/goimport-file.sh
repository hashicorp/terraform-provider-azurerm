#!/bin/bash

# remove blank lines in go imports then run goimports

if [ $# != 1 ] ; then
  echo "usage: $0 <filename>"
  exit 1
fi

# remove empty lines inside import block via sed
sed_expression='
  /^import/,/)/ {
    /^$/d
  }
'

case "$OSTYPE" in
    "linux-gnu"*)
        sed -i -e "$sed_expression" $1
        ;;
    "darwin"*)
        sed_version="$(sed --version)"

        if [[ $sed_version == *GNU* ]]; then
            # For MacOS users who aliases the `sed` from `gsed`
            sed -i -e "$sed_expression" $1
        else
            # The Posix sed
            sed -i '' -e "$sed_expression" $1
        fi

        ;;
esac 

goimports -w $1

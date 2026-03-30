#!/bin/bash
# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0


# remove blank lines in go imports then run goimports

if [ $# != 1 ] ; then
  echo "usage: $0 <filename>"
  exit 1
fi

# remove empty lines inside import block via sed
sed_expression='
  /^import (/,/)/ {
    /^$/d
  }
'

case "$OSTYPE" in
    "linux-gnu"*)
        sed -i -e "$sed_expression" "$1"
        ;;
    "darwin"*)
        # Check if we have GNU sed (e.g. MacOS users who alias sed to gsed)
        if sed --version >/dev/null 2>&1; then
            sed -i -e "$sed_expression" "$1"
        else
            # The Posix sed (default on MacOS)
            sed -i '' -e "$sed_expression" "$1"
        fi

        ;;
esac 

goimports -w "$1"

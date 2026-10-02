#!/usr/bin/env bash
# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0


echo "==> Checking documentation terraform blocks are formatted..."

if ! terrafmt diff -c -q -p "*.html.markdown" ./website; then
  echo "------------------------------------------------"
  echo ""
  echo "The preceding files contain terraform blocks that are not correctly formatted or contain errors."
  echo "You can fix this by running make tools and then terrafmt on them."
  echo ""
  echo "to easily fix all terraform blocks:"
  echo "$ make terrafmt"
  echo ""
  echo "format only website config blocks:"
  echo "$ terrafmt fmt -p \"*.html.markdown\" ./website"
  echo ""
  echo "format a single website file:"
  echo "$ terrafmt fmt ./website/path/to/file.html.markdown"
  echo ""
  echo "on windows:"
  echo "$ terrafmt fmt -p \"*.html.markdown\" ./website"
  echo ""
  exit 1
fi

exit 0

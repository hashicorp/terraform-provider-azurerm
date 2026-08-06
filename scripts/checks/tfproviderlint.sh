#!/usr/bin/env bash
# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0




function on_failure {
  echo ""
  echo "==> tfproviderlint failed!"
  echo "    Common causes:"
  echo "    - Schema definition issues (e.g. missing Required/Optional, invalid types)"
  echo "    - Acceptance test formatting issues (run: terrafmt fmt -f <file>)"
  echo "    - tfproviderlint rule violations (see https://github.com/bflad/tfproviderlint)"
  echo ""
}

function runTests {
  echo "==> Checking source code against terraform provider linters..."
	tfproviderlintx \
        -AT001\
        -AT001.ignored-filename-suffixes _data_source_test.go\
        -AT003 -AT005 -AT006 -AT007 -AT008 -AT009 -AT010 -AT011\
        -R001 -R002 -R003 -R004 -R006 -R006.package-aliases pluginsdk\
        -R007 -R010 -R011 -R012 -R013 -R014 -R015 -R016 -R017 -R019\
        -S001 -S002 -S003 -S004 -S005 -S006 -S007 -S008 -S009 -S010 -S011 -S012 -S013 -S014 -S015 -S016 -S017 -S018 -S019 -S020\
        -S021 -S022 -S023 -S024 -S025 -S026 -S027 -S028 -S029 -S030 -S031 -S032 -S033 -S035 -S036 -S037\
        -V009 -V010\
        -XAT001 -XR006 -XR008\
        ./internal/...
	bash ./scripts/checks/terrafmt-acctests.sh
}

function main {
  runTests || { on_failure; exit 1; }
}

main


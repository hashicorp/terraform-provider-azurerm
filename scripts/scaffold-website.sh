#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


function verifyVars {
  echo "==> Verifying required variables are set..."
  if [ "${BRAND_NAME}" == "" ]; then
    echo -e "\nBRAND_NAME is unset, exiting"
    exit 1
  fi
  if [ "${RESOURCE_NAME}" == "" ]; then
    echo -e "\nRESOURCE_NAME is unset, exiting"
    exit 1
  fi
  if [ "${RESOURCE_TYPE}" == "" ]; then
    echo -e "\nRESOURCE_TYPE is unset, exiting"
    exit 1
  fi
  if [ "${RESOURCE_TYPE}" == "resource" ]; then
    if [ "${RESOURCE_ID}" == "" ]; then
      echo -e "\nRESOURCE_ID is unset, exiting"
      exit 1
    fi
  fi
  if [[ -n "${EXAMPLE}" ]]; then
      if [[ -z "${SERVICE_TYPE}" ]]; then
          echo -e "\nSERVICE_TYPE is unset, exiting"
          exit 1
      fi
      if [[ -z "${TESTCASE}" ]]; then
          echo -e "\nTESTCASE is unset, exiting"
          exit 1
      fi
  fi

  echo "==> Validated."
}

function scaffoldDocumentation {
  echo "==> Scaffolding Documentation..."
  if [[ -z "$EXAMPLE" ]]; then
      go run ./internal/tools/website-scaffold/main.go -name "${RESOURCE_NAME}" -brand-name "${BRAND_NAME}" -type "${RESOURCE_TYPE}" -resource-id "${RESOURCE_ID}" -website-path ./website/
  else
      go run ./internal/tools/website-scaffold/main.go -name "${RESOURCE_NAME}" -brand-name "${BRAND_NAME}" -type "${RESOURCE_TYPE}" -resource-id "${RESOURCE_ID}" -website-path ./website/ -example -root-dir . -service-pkg ./internal/services/"${SERVICE_TYPE}" -testcase "${TESTCASE}"
  fi
  echo "==> Done."
}

function main {
  verifyVars
  scaffoldDocumentation
}

main

#!/usr/bin/env bash

function verifyVars {
  echo "==> Verifying required variables are set..."
  if [ "${BRAND_NAME}" == "" ]; then
    echo "\nBRAND_NAME is unset, exiting"
    exit 1
  fi
  if [ "${RESOURCE_NAME}" == "" ]; then
    echo "\nRESOURCE_NAME is unset, exiting"
    exit 1
  fi
  if [ "${RESOURCE_TYPE}" == "" ]; then
    echo "\nRESOURCE_TYPE is unset, exiting"
    exit 1
  fi
  if [ "${RESOURCE_TYPE}" == "resource" ]; then
    if [ "${RESOURCE_ID}" == "" ]; then
      echo "\nRESOURCE_ID is unset, exiting"
      exit 1
    fi
  fi

  echo "==> Validated."
}

function scaffoldDocumentation {
  echo "==> Scaffolding Documentation..."
  go run azurerm/internal/tools/website-scaffold/main.go -name "${RESOURCE_NAME}" -brand-name "${BRAND_NAME}" -type "${RESOURCE_TYPE}" -resource-id "${RESOURCE_ID}" -website-path ./website/
  echo "==> Done."
}

function main {
  verifyVars
  scaffoldDocumentation
}

main

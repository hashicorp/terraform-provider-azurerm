#!/usr/bin/env bash
set -e

if [ -z "${1-}" ]; then
  echo "Usage: $0 <service_name> [branch_name]"
  echo "Example: $0 compute"
  echo "Example: $0 compute my-feature-branch"
  exit 1
fi

SERVICE="$1"
BRANCH="${2-}"

if [ -z "${BRANCH}" ]; then
  if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    BRANCH=$(git rev-parse --abbrev-ref HEAD)
  else
    BRANCH="main"
  fi
fi

if [ -z "${ARM_BLOCKBUSTER_VIDEO-}" ]; then
  echo "Error: ARM_BLOCKBUSTER_VIDEO environment variable is not set."
  echo "This must be set to the Azure Storage Account name containing the cassettes."
  exit 1
fi

if ! command -v azcopy >/dev/null 2>&1; then
  echo "Error: azcopy is not installed or not in PATH."
  exit 1
fi

SOURCE_URL="https://${ARM_BLOCKBUSTER_VIDEO}.blob.core.windows.net/cassettes/${BRANCH}/${SERVICE}/"
DEST_DIR="./internal/services/${SERVICE}/vcrtestdata/"

# Ensure the destination directory exists
mkdir -p "${DEST_DIR}"

echo "Pulling cassettes for service '${SERVICE}' from branch '${BRANCH}'..."
echo "Source: ${SOURCE_URL}"
echo "Destination: ${DEST_DIR}"

azcopy copy "${SOURCE_URL}*" "${DEST_DIR}" --recursive

echo "Successfully pulled cassettes for ${SERVICE}."

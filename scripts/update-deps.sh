#!/bin/bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


REPO_DIR="$(cd "$(dirname "$0")"/.. && pwd)"
cd "${REPO_DIR}"

PROVIDER_REPO="hashicorp/terraform-provider-azurerm"
TRUNK="main"

usage() {
  echo "Usage: $0 dependent-module-name [version] [-p] [-f]" >&2
  echo >&2
  echo " -p  Optionally open a PR. Requires GITHUB_TOKEN env var." >&2
  echo " -f  Force dependency update when \`${TRUNK}\` branch is not checked out" >&2
  echo >&2
}

which jq 2>/dev/null 1>/dev/null || jq() {
  cat -
}

while [ $# -gt 0 ]; do
  while getopts ':pfh' opt; do
    case "$opt" in
      p)
        PR=1
        ;;
      f)
        FORCE=1
        ;;
      *)
        usage
        exit 1
        ;;
    esac
  done

  # shift options to get positional arguments
  [ $? -eq 0 ] || exit 1
  [ $OPTIND -gt $# ] && break
  shift $[$OPTIND-1]
  OPTIND=1
  ARGS[${#ARGS[*]}]="${1}"
  shift
done

MODULE_PATH="${ARGS[0]}"
if [[ -z "${MODULE_PATH}" ]]; then
  echo "Error: no module name specified" >&2
  exit 1
fi

_mod=(${MODULE_PATH//\// })
for (( i=${#_mod[@]}-1 ; i>=0 ; i-- )) ; do
  if ! [[ "${_mod[${i}]}" =~ ^v[0-9]$ ]]; then
    MODULE_NAME="${_mod[${i}]}"
    break
  fi
done

if [[ -z "${MODULE_NAME}" ]]; then
  echo "Error: could not determine module name from path: ${MODULE_PATH}" >&2
  exit 2
fi

if [[ -n "${PR}" ]] && [[ -z "${GITHUB_TOKEN}" ]];then
  echo "Error: must set GITHUB_TOKEN when \`-p\` specified" >&2
  exit 2
fi

BRANCH_NAME="dependencies/${MODULE_NAME}"

CURRENT_BRANCH="$(git rev-parse --abbrev-ref HEAD)"
if [[ "${CURRENT_BRANCH}" != "${TRUNK}" ]]; then
  if [[ "${FORCE}" == "1" ]]; then
    echo "Caution: Proceeding with dependencies update from current branch ${CURRENT_BRANCH}"
  else
    echo "Dependency updates should be based on \`${TRUNK}\` branch. Specify \`-f\` to override." >&2
    exit 1
  fi
fi


if git branch -l ${BRANCH_NAME} | grep -q ${BRANCH_NAME}; then
  if [[ "${FORCE}" == "1" ]]; then
    echo "Caution: Deleting existing branch ${BRANCH_NAME} as \`-f\` was specified."
    git branch -D ${BRANCH_NAME}
  else
    echo "The branch \`${BRANCH_NAME}\` already exists. Specify \`-f\` to delete it." >&2
    exit 1
  fi
fi

if [[ -n "$(git status --short)" ]]; then
  echo "Error: working tree is dirty" >&2
  exit 2
fi

set -e

pwd
# Ensure latest changes are checked out
( set -x; git pull origin "${TRUNK}" )

echo "Checking out a new branch..."
(
  set -x
  git checkout -b ${BRANCH_NAME}
)

VERSION="${ARGS[1]}"

echo "Updating dependency..."
(
  set -x
  go get ${MODULE_PATH}@${VERSION:-latest}
  go get ./...
  go mod vendor
  go mod tidy
)

if [[ -z "${VERSION}" ]]; then
  _mod=($(grep "${MODULE_PATH}" "${REPO_DIR}"/go.mod))
  if [[ "${#_mod[@]}" == "2" ]]; then
    VERSION="${_mod[1]}"
  else
      echo "Could not determine latest version of ${MODULE_PATH}" >&2
      exit 3
  fi
fi

echo "Committing..."

COMMIT_MSG="dependencies: updating to \`${VERSION}\` of \`${MODULE_PATH}\`"

(
  set -x
  git add "${REPO_DIR}"/go.* "${REPO_DIR}"/vendor
  git commit -m "${COMMIT_MSG}"
  git push -fu origin "${BRANCH_NAME}"
)

if [[ -n "${PR}" ]]; then
  echo "Opening pull request..."
  (
    set -x
    curl -X POST -H "Authorization: Bearer ${GITHUB_TOKEN}" "https://api.github.com/repos/${PROVIDER_REPO}/pulls" \
      -d "{\"title\":\"${COMMIT_MSG}\",\"body\":\"\",\"base\":\"${TRUNK}\",\"head\":\"${BRANCH_NAME}\"}" | jq .html_url
  )
fi

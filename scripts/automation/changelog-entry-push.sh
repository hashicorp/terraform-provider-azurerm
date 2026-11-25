#!/usr/bin/env bash
# Copyright IBM Corp. 2014, 2026
# SPDX-License-Identifier: MPL-2.0

# Pushes a changelog entry onto a pull request's branch, after checking the commit is that and nothing else. Run from the
# changelog-entry-add-to-pr workflow, in the one job that holds a token allowed to push.
#
# The entry is written by an earlier job that has the pull request's files on disk and no such token. It hands the commit over as a git
# bundle. This script pushes it only if it is one commit on top of the pull request's latest commit that changes .changelog/<number> alone,
# and never writes a file from the pull request to disk.
#
# Required environment:
#   BUNDLE_PATH - the git bundle holding the commit, on a branch named changelog-entry
#   PR_NUMBER   - number of the pull request
#   FETCH_URL   - repository the pull request was opened against
#   PUSH_URL    - repository the pull request's branch lives in, which is a fork for a pull request from a fork
#   HEAD_REF    - name of the pull request's branch
#
# Optional environment:
#   PUSH_TOKEN  - token to push with; without it git's own credentials are used

set -euo pipefail

bundle="${BUNDLE_PATH:?BUNDLE_PATH must be set}"
pr_number="${PR_NUMBER:?PR_NUMBER must be set}"
fetch_url="${FETCH_URL:?FETCH_URL must be set}"
push_url="${PUSH_URL:?PUSH_URL must be set}"
head_ref="${HEAD_REF:?HEAD_REF must be set}"

fail() {
  echo "::error::$1" >&2
  exit 1
}

[[ "${pr_number}" =~ ^[0-9]+$ ]] || fail "Invalid pull request number: ${pr_number}"
git check-ref-format "refs/heads/${head_ref}" || fail "Invalid branch name: ${head_ref}"
[[ -f "${bundle}" ]] || fail "No bundle found at ${bundle}"

bundle=$(cd "$(dirname "${bundle}")" && pwd)/$(basename "${bundle}")

# a repository with no working folder, so nothing from the pull request is written to disk beside the token
work=$(mktemp -d)
trap 'rm -rf "${work}"' EXIT
git init --quiet --bare "${work}"
cd "${work}"

git fetch --quiet --no-tags --depth=1 "${fetch_url}" "refs/pull/${pr_number}/head:refs/heads/pr-head"
head_sha=$(git rev-parse refs/heads/pr-head)

git bundle verify --quiet "${bundle}" || fail "The entry was not written on top of the pull request's latest commit. Comment /changelog again."
git fetch --quiet --no-tags "${bundle}" "refs/heads/changelog-entry:refs/heads/changelog-entry"
new_sha=$(git rev-parse refs/heads/changelog-entry)

[[ "$(git rev-list --count "${head_sha}..${new_sha}")" == "1" ]] || fail "Refusing to push, expected one commit on top of the pull request."
[[ "$(git rev-parse "${new_sha}^")" == "${head_sha}" ]] || fail "The entry was not written on top of the pull request's latest commit. Comment /changelog again."

changed=0
while IFS= read -r path; do
  changed=$((changed + 1))
  [[ "${path}" =~ ^\.changelog/${pr_number}\.(hcl|md|yml)$ ]] || fail "Refusing to push, the commit changes a file that is not this pull request's changelog entry: ${path}"

  # an entry is a plain file: a link or a program under that name is refused. A removed file has no mode.
  mode=$(git ls-tree "${new_sha}" -- "${path}" | cut -d' ' -f1)
  [[ -z "${mode}" || "${mode}" == "100644" ]] || fail "Refusing to push, the changelog entry is not a plain file: ${path}"
done < <(git diff --name-only --no-renames "${head_sha}" "${new_sha}")

((changed > 0)) || fail "Refusing to push, the commit changes nothing."

# the token goes to git through the environment, so it is not on a command line or in a config file
if [[ -n "${PUSH_TOKEN:-}" ]]; then
  export GIT_CONFIG_COUNT=1
  export GIT_CONFIG_KEY_0="http.https://github.com/.extraheader"
  GIT_CONFIG_VALUE_0="AUTHORIZATION: basic $(printf 'x-access-token:%s' "${PUSH_TOKEN}" | base64 | tr -d '\n')"
  export GIT_CONFIG_VALUE_0
fi

git push --quiet "${push_url}" "${new_sha}:refs/heads/${head_ref}"
echo "Pushed changelog entry ${new_sha} to ${head_ref}"

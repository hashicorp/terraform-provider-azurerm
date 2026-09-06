#!/usr/bin/env bash
# Copyright IBM Corp. 2014, 2026
# SPDX-License-Identifier: MPL-2.0

# Lint the terraform blocks embedded in website docs (or any file terrafmt understands) with
# tflint. terrafmt extracts every block, each block is written to its own directory and
# tflint --recursive lints them in one process per file. Findings are reported as file:line.
# (Acceptance tests are linted from their rendered configs instead - see tflint-acctests.sh.)
#
#   scripts/checks/tflint-embedded.sh <dir> <file pattern> [extra tflint args...]
#   scripts/checks/tflint-embedded.sh --files [extra tflint args...]   # newline-separated paths on stdin
#
#   scripts/checks/tflint-embedded.sh ./website '*.html.markdown'
#   git diff --name-only main -- '*.html.markdown' | scripts/checks/tflint-embedded.sh --files
#
# JOBS controls parallelism (default: number of cores). Each file is handled by re-invoking
# this script in --process-file mode via xargs -P.

set -eo pipefail
trap 'exit 130' INT TERM

TFLINT_CONFIG="${TFLINT_CONFIG:-$(pwd -P)/.tflint.hcl}"
export TFLINT_CONFIG

process_one_file() {
  local filename=$1
  shift
  local blocks td exit_code=0

  # capture first so a terrafmt/jq failure propagates instead of reading as "no blocks"
  blocks=$(terrafmt blocks --fmtcompat --json "${filename}" | jq --compact-output '.blocks[]?')
  [[ -z "${blocks}" ]] && return 0

  td=$(mktemp -d)
  local n=0
  while IFS= read -r block; do
    n=$((n + 1))
    mkdir -p "${td}/${n}"
    jq -r '.text' <<< "${block}" > "${td}/${n}/main.tf"
    jq -r '.start_line' <<< "${block}" > "${td}/${n}/.start"
    # doc snippets sometimes reference a variable/local they do not define; declare those (as
    # unknown values) rather than failing the whole file
    {
      echo 'variable "tfmt_placeholder" {}'
      while read -r v; do
        grep -qE "^\s*variable\s+\"${v}\"" "${td}/${n}/main.tf" || echo "variable \"${v}\" {}"
      done < <(grep -oE '\bvar\.[A-Za-z0-9_-]+' "${td}/${n}/main.tf" | sort -u | sed 's/^var\.//' || true)
      while read -r l; do
        grep -qE "^\s*${l}\s*=" "${td}/${n}/main.tf" || echo "locals { ${l} = var.tfmt_placeholder }"
      done < <(grep -oE '\blocal\.[A-Za-z0-9_-]+' "${td}/${n}/main.tf" | sort -u | sed 's/^local\.//' || true)
    } > "${td}/${n}/tfmt_placeholders.tf"
  done <<< "${blocks}"

  local output
  set +e
  output=$(tflint --config "${TFLINT_CONFIG}" --chdir="${td}" --recursive --format=compact "$@" 2>&1)
  exit_code=$?
  set -e

  if [[ ${exit_code} -ne 0 ]]; then
    # compact lines look like "<block>/main.tf:<line>:<col>: <severity> - <message> (<rule>)";
    # map them back to the source file line (block start line + line within the block)
    while IFS= read -r line; do
      if [[ "${line}" =~ ^(.*/)?([0-9]+)/main\.tf:([0-9]+):([0-9]+):\ (.*)$ ]]; then
        local b="${BASH_REMATCH[2]}" l="${BASH_REMATCH[3]}" msg="${BASH_REMATCH[5]}"
        echo "${filename}:$(( $(cat "${td}/${b}/.start") + l )): ${msg}"
      elif [[ -n "${line}" && "${line}" != *"issue(s) found"* ]]; then
        echo "${filename}: ${line}"
      fi
    done <<< "${output}"
  fi

  rm -rf "${td}"
  return "${exit_code}"
}

if [[ "${1:-}" == "--process-file" ]]; then
  shift
  process_one_file "$@"
  exit $?
fi

if [[ "${1:-}" == "--files" ]]; then
  shift
  list_files() { grep . || true; }
  echo "==> Checking terraform blocks in the given files with tflint..."
else
  dir=${1:-}
  pattern=${2:-}
  if [[ -z "${dir}" || -z "${pattern}" ]]; then
    echo "usage: $0 <dir> <file pattern> [tflint args...]  |  $0 --files [tflint args...] < paths" >&2
    exit 64
  fi
  shift 2
  list_files() { find "${dir}" -type f -name "${pattern}" | sort; }
  echo "==> Checking ${pattern} terraform blocks under ${dir} with tflint..."
fi

for tool in terrafmt tflint jq; do
  command -v "${tool}" >/dev/null || { echo "${tool} not installed (see 'make tools')" >&2; exit 1; }
done

tflint --config "${TFLINT_CONFIG}" --init >/dev/null

JOBS="${JOBS:-$(nproc 2>/dev/null || sysctl -n hw.ncpu 2>/dev/null || echo 4)}"
# xargs exits 123 when any worker fails, which fails the make target / CI step
list_files | xargs -P "${JOBS}" -I {} "${BASH_SOURCE[0]}" --process-file {} "$@"

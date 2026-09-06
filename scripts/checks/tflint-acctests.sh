#!/usr/bin/env bash
# Copyright IBM Corp. 2014, 2026
# SPDX-License-Identifier: MPL-2.0

# Lint the terraform used by the acceptance tests with tflint.
#
# Rather than scraping HCL blocks out of _test.go files, the tests themselves render every step
# config: with ARM_ACCTEST_EXPORT_CONFIGS_DIR set the acceptance framework writes each step's
# config to <dir>/<test name>/<step>/main.tf and skips the test (see exportStepConfigs in
# internal/acceptance/testcase.go), so composed templates, Sprintf placeholders and HCL built in
# Go loops all come out exactly as the test would apply them. tflint then lints the lot.
#
#   scripts/checks/tflint-acctests.sh [go package patterns...]   # default ./internal/services/...
#
#   scripts/checks/tflint-acctests.sh ./internal/services/compute/...
#
# JOBS controls tflint parallelism (default: number of cores). Findings are reported as
# "<test name> step <n> line <l>: <message>".

set -eo pipefail
trap 'exit 130' INT TERM

TFLINT_CONFIG="${TFLINT_CONFIG:-$(pwd -P)/.tflint.hcl}"
JOBS="${JOBS:-$(nproc 2>/dev/null || sysctl -n hw.ncpu 2>/dev/null || echo 4)}"
export TFLINT_CONFIG

for tool in tflint go; do
  command -v "${tool}" >/dev/null || { echo "${tool} not installed (see 'make tools')" >&2; exit 1; }
done

pkgs=("$@")
[[ ${#pkgs[@]} -eq 0 ]] && pkgs=(./internal/services/...)

work_dir=$(mktemp -d)
trap 'rm -rf "${work_dir}"' EXIT
export_dir="${work_dir}/export"
shard_dir="${work_dir}/shards"
mkdir -p "${export_dir}" "${shard_dir}"

# Tests skip themselves when the env they need is missing, and the framework's PreCheck insists
# on credentials, so give every ARM_* variable the tests read a placeholder value (real values in
# the environment win). Nothing is run against Azure - the export hook skips before any client
# is built.
export TF_ACC=1
export ARM_ACCTEST_EXPORT_CONFIGS_DIR="${export_dir}"
: "${ARM_TEST_LOCATION:=westeurope}" "${ARM_TEST_LOCATION_ALT:=northeurope}" "${ARM_TEST_LOCATION_ALT2:=eastus2}"
export ARM_TEST_LOCATION ARM_TEST_LOCATION_ALT ARM_TEST_LOCATION_ALT2
while read -r name; do
  [[ -n "${!name:-}" ]] && continue
  case "${name}" in
    *SUBSCRIPTION_ID*|*CLIENT_ID*|*TENANT_ID*|*PRINCIPAL_ID*|*OBJECT_ID*) export "${name}=00000000-0000-0000-0000-000000000000" ;;
    *) export "${name}=placeholder" ;;
  esac
done < <(grep -rhoE '"ARM_[A-Z0-9_]+"' --include='*.go' internal | tr -d '"' | sort -u)

echo "==> Exporting acceptance test configs (${pkgs[*]})..."
# every test skips after exporting, so a non-zero exit is a compile error or a test that could
# not render its config
if ! go test -count=1 -run '^TestAcc' "${pkgs[@]}" > "${work_dir}/go-test.log" 2>&1; then
  grep -E '^(FAIL|--- FAIL|panic:|\s+\S+\.go:[0-9]+:)' "${work_dir}/go-test.log" || tail -20 "${work_dir}/go-test.log"
  echo "ERROR: exporting acceptance test configs failed (see above)" >&2
  exit 1
fi

configs=$(find "${export_dir}" -name main.tf | wc -l | tr -d ' ')
echo "==> Checking ${configs} exported step configs with tflint..."
tflint --config "${TFLINT_CONFIG}" --init >/dev/null

# one tflint process per shard of tests, JOBS shards in parallel
i=0
while IFS= read -r test_dir; do
  mkdir -p "${shard_dir}/$((i % JOBS))"
  mv "${test_dir}" "${shard_dir}/$((i % JOBS))/"
  i=$((i + 1))
done < <(find "${export_dir}" -mindepth 1 -maxdepth 1 -type d | sort)

# shellcheck disable=SC2329 # invoked via xargs below
lint_shard() {
  local output code
  cd -P "$1"
  set +e
  output=$(tflint --config "${TFLINT_CONFIG}" --recursive --format=compact 2>&1)
  code=$?
  set -e
  # compact lines look like "<test name>/<step>/main.tf:<line>:<col>: <severity> - <message> (<rule>)"
  while IFS= read -r line; do
    if [[ "${line}" =~ ^(.+)/([0-9]+)/main\.tf:([0-9]+):[0-9]+:\ (.*)$ ]]; then
      echo "${BASH_REMATCH[1]} step ${BASH_REMATCH[2]} line ${BASH_REMATCH[3]}: ${BASH_REMATCH[4]}"
    elif [[ -n "${line}" && "${line}" != *"issue(s) found"* ]]; then
      echo "${line}"
    fi
  done <<< "${output}"
  return "${code}"
}
export -f lint_shard

# xargs exits 123 when any shard fails, which fails the make target / CI step
code=0
# shellcheck disable=SC2016 # $1 is expanded by the inner bash
find "${shard_dir}" -mindepth 1 -maxdepth 1 -type d | sort | xargs -P "${JOBS}" -I {} bash -c 'lint_shard "$1"' _ {} > "${work_dir}/findings" || code=$?
sort "${work_dir}/findings"
exit "${code}"

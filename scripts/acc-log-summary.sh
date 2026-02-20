#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Usage:
  scripts/acc-log-summary.sh [--dir DIR] [--file PATH] [--glob GLOB] [--context N] [--regex REGEX]

Summarizes a Terraform acceptance log (typically written via terraform-plugin-testing).

Selection order when --file is not given:
  1) TF_ACC_LOG_PATH (single log file)
  2) TF_LOG_PATH_MASK (per-test logs) -> uses directory + inferred glob
  3) Latest file in --dir matching --glob (defaults: tf-*.log)

Output:
  - First likely root-cause marker + surrounding context
  - Tail of the file for quick sanity checking

Examples:
  scripts/acc-log-summary.sh
  scripts/acc-log-summary.sh --dir /tmp --glob '*.log'
  scripts/acc-log-summary.sh --file /tmp/tf-TestAccSomething.log
  scripts/acc-log-summary.sh --context 120
EOF
}

log_dir="."
log_glob="tf-*.log"
log_file=""
context_lines=80

default_regex='^(Error:|panic:|--- FAIL:)|context deadline exceeded|StatusCode=409|ScopeLocked|TooManyRequests|StatusCode=429|429 Too Many Requests|Retry-After|X-Ms-Failure-Cause'
regex="$default_regex"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --dir)
      log_dir="$2"
      shift 2
      ;;
    --file)
      log_file="$2"
      shift 2
      ;;
    --glob)
      log_glob="$2"
      shift 2
      ;;
    --context)
      context_lines="$2"
      shift 2
      ;;
    --regex)
      regex="$2"
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "Unknown arg: $1" >&2
      usage >&2
      exit 2
      ;;
  esac
done

if [[ -z "$log_file" && -n "${TF_ACC_LOG_PATH:-}" ]]; then
  log_file="$TF_ACC_LOG_PATH"
fi

if [[ -z "$log_file" && -n "${TF_LOG_PATH_MASK:-}" ]]; then
  # TF_LOG_PATH_MASK example: /tmp/tf-%s.log
  mask="$TF_LOG_PATH_MASK"
  inferred_dir=$(dirname "$mask")
  inferred_glob=$(basename "$mask")
  inferred_glob=${inferred_glob//%s/*}
  # Respect explicit --dir/--glob, but if user left defaults, use inferred values.
  if [[ "$log_dir" == "." ]]; then
    log_dir="$inferred_dir"
  fi
  if [[ "$log_glob" == "tf-*.log" ]]; then
    log_glob="$inferred_glob"
  fi
fi

if [[ -z "$log_file" ]]; then
  shopt -s nullglob
  candidates=("$log_dir"/$log_glob)
  shopt -u nullglob

  if [[ ${#candidates[@]} -eq 0 ]]; then
    echo "No logs found in $log_dir matching $log_glob" >&2
    echo "Hint: pass --file explicitly, or set TF_LOG_PATH_MASK / TF_ACC_LOG_PATH." >&2
    exit 1
  fi

  # ls -t sorts by mtime desc.
  log_file=$(ls -1t -- "${candidates[@]}" 2>/dev/null | head -n 1)
fi

if [[ ! -f "$log_file" ]]; then
  echo "Log file not found: $log_file" >&2
  exit 1
fi

printf 'Log: %s\n' "$log_file"
printf 'Size: %s bytes\n' "$(wc -c <"$log_file" | tr -d ' ')"

# Portable-ish mtime printing (GNU stat vs BSD stat). If it fails, skip.
if mtime=$(stat -c '%y' "$log_file" 2>/dev/null); then
  printf 'Last write: %s\n\n' "$mtime"
elif mtime=$(stat -f '%Sm' -t '%Y-%m-%d %H:%M:%S %z' "$log_file" 2>/dev/null); then
  printf 'Last write: %s\n\n' "$mtime"
else
  printf '\n'
fi

match=$(grep -n -E -m 1 "$regex" "$log_file" || true)

if [[ -z "$match" ]]; then
  echo "No obvious error markers found with pattern:" >&2
  echo "  $regex" >&2
  echo
  echo "Tail (last 120 lines):"
  tail -n 120 "$log_file"
  exit 0
fi

line_no=${match%%:*}

start=$(( line_no - context_lines ))
end=$(( line_no + context_lines ))
if (( start < 1 )); then start=1; fi

printf 'First match: line %s\n' "$line_no"
printf '%s\n\n' "$match"

printf 'Context (lines %s-%s):\n' "$start" "$end"
sed -n "${start},${end}p" "$log_file"

printf '\nTail (last 120 lines):\n'
tail -n 120 "$log_file"

#!/usr/bin/env bash
# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0
#
# serve-data-api.sh — clone (or reuse) the Pandora repository and run its Data API.
#
# The scaff tool resolves resource schemas from the Pandora Data API, which by
# default listens on http://localhost:8080. This script fetches the Pandora
# source and launches `data-api serve` against its api-definitions/ directory.
#
# Usage:
#   ./serve-data-api.sh [--port PORT] [--services svc1,svc2] [--update] [--help]
#
# Environment variables (flags take precedence):
#   PANDORA_DIR        Directory to clone into / reuse. Defaults to a sibling
#                      "pandora" checkout next to this provider repo when present,
#                      otherwise "${HOME}/.cache/scaff/pandora".
#   PANDORA_REPO_URL   Git URL to clone (default: https://github.com/hashicorp/pandora.git).
#   PANDORA_BRANCH     Branch to clone (default: the repository default branch).
#   PANDORA_API_PORT   Port for the Data API (default: 8080). Honoured natively
#                      by data-api as well.
#   PANDORA_SERVICES   Comma-separated Service names to load (default: all).
#   PANDORA_UPDATE     When set to 1, `git pull --ff-only` an existing checkout
#                      before serving (off by default so a working checkout is
#                      never disturbed).
#
set -euo pipefail

log() { printf '\033[1;34m==>\033[0m %s\n' "$*" >&2; }
warn() { printf '\033[1;33mwarning:\033[0m %s\n' "$*" >&2; }
die() { printf '\033[1;31merror:\033[0m %s\n' "$*" >&2; exit 1; }

# --- Resolve paths -----------------------------------------------------------
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# scripts -> scaff -> tools -> internal -> <provider repo root>
REPO_ROOT="$(cd "${SCRIPT_DIR}/../../../.." && pwd)"
WORKSPACE_ROOT="$(cd "${REPO_ROOT}/.." && pwd)"

# --- Configuration (env with defaults) ---------------------------------------
PANDORA_REPO_URL="${PANDORA_REPO_URL:-https://github.com/hashicorp/pandora.git}"
PANDORA_BRANCH="${PANDORA_BRANCH:-}"
PANDORA_API_PORT="${PANDORA_API_PORT:-8080}"
PANDORA_SERVICES="${PANDORA_SERVICES:-}"
PANDORA_UPDATE="${PANDORA_UPDATE:-}"

# Default PANDORA_DIR: prefer an existing sibling checkout, else a cache dir.
if [[ -z "${PANDORA_DIR:-}" ]]; then
  if [[ -d "${WORKSPACE_ROOT}/pandora/tools/data-api" ]]; then
    PANDORA_DIR="${WORKSPACE_ROOT}/pandora"
  else
    PANDORA_DIR="${HOME}/.cache/scaff/pandora"
  fi
fi

# --- Parse arguments (override env) ------------------------------------------
while [[ $# -gt 0 ]]; do
  case "$1" in
    --port) PANDORA_API_PORT="${2:?--port requires a value}"; shift 2 ;;
    --port=*) PANDORA_API_PORT="${1#*=}"; shift ;;
    --services) PANDORA_SERVICES="${2:?--services requires a value}"; shift 2 ;;
    --services=*) PANDORA_SERVICES="${1#*=}"; shift ;;
    --dir) PANDORA_DIR="${2:?--dir requires a value}"; shift 2 ;;
    --dir=*) PANDORA_DIR="${1#*=}"; shift ;;
    --update) PANDORA_UPDATE=1; shift ;;
    -h|--help)
      sed -n '3,23p' "${BASH_SOURCE[0]}" | sed 's/^# \{0,1\}//'
      exit 0
      ;;
    *) die "unknown argument: $1 (use --help)" ;;
  esac
done

# --- Preflight ---------------------------------------------------------------
command -v git >/dev/null 2>&1 || die "git is required but was not found on PATH"
command -v go >/dev/null 2>&1 || die "go is required but was not found on PATH"

DATA_API_DIR="${PANDORA_DIR}/tools/data-api"
DATA_DIRECTORY="${PANDORA_DIR}/api-definitions"

# --- Clone or reuse ----------------------------------------------------------
if [[ -d "${DATA_API_DIR}" ]]; then
  log "Using existing Pandora checkout at ${PANDORA_DIR}"
  if [[ "${PANDORA_UPDATE}" == "1" ]]; then
    log "Updating checkout (git pull --ff-only)"
    git -C "${PANDORA_DIR}" pull --ff-only || warn "could not fast-forward the existing checkout; continuing with it as-is"
  fi
else
  log "Cloning Pandora into ${PANDORA_DIR}"
  mkdir -p "$(dirname "${PANDORA_DIR}")"
  # Shallow clone without submodules: the Data API only needs api-definitions/,
  # not the (large) rest-api-specs / msgraph-metadata submodules.
  clone_args=(--depth 1)
  [[ -n "${PANDORA_BRANCH}" ]] && clone_args+=(--branch "${PANDORA_BRANCH}")
  git clone "${clone_args[@]}" "${PANDORA_REPO_URL}" "${PANDORA_DIR}"
fi

[[ -d "${DATA_API_DIR}" ]] || die "data-api directory not found at ${DATA_API_DIR}"
[[ -d "${DATA_DIRECTORY}" ]] || die "api-definitions directory not found at ${DATA_DIRECTORY}"

# --- Build -------------------------------------------------------------------
log "Building the Data API"
( cd "${DATA_API_DIR}" && go build -o data-api . )

# --- Serve -------------------------------------------------------------------
serve_args=(serve --port "${PANDORA_API_PORT}" --data-directory "${DATA_DIRECTORY}")
[[ -n "${PANDORA_SERVICES}" ]] && serve_args+=(--services "${PANDORA_SERVICES}")

log "Serving the Data API at http://localhost:${PANDORA_API_PORT} (Ctrl-C to stop)"
log "Data directory: ${DATA_DIRECTORY}"
cd "${DATA_API_DIR}"
exec ./data-api "${serve_args[@]}"

#!/bin/bash

# Fills the agent's shared provider directory with the provider versions this build's tests will ask
# for, so that tests link to one copy rather than each downloading and unpacking their own.
#
# Tests only ever read the directory - it is their filesystem mirror (see TF_CLI_CONFIG_FILE) and a
# test downloads anything missing from it as normal. This step is the directory's only writer, since
# Terraform's own plugin cache is not safe to fill from tests running in parallel.
#
# Nothing here fails the build: a provider which can't be fetched is downloaded by the tests instead.

TERRAFORM="%env.TF_ACC_TERRAFORM_PATH%"
MIRROR_DIR="%env.TF_ACC_TERRAFORM_PROVIDER_CACHE_MIRROR_PATH%"
CLI_CONFIG_FILE="%env.TF_ACC_TERRAFORM_PROVIDER_CACHE_CONFIG_FILE%"
SERVICE_PATH="%SERVICE_PATH%"
# provider versions which no build has used for this many days are deleted
RETENTION_DAYS=14

CACHE_DIR="$(dirname "$MIRROR_DIR")"

mkdir -p "$MIRROR_DIR"

cat > "$CLI_CONFIG_FILE" <<EOF
provider_installation {
  filesystem_mirror {
    path = "$MIRROR_DIR"
  }
  direct {}
}
EOF

# that config is for the tests - this step has to download from the registry
unset TF_CLI_CONFIG_FILE
export TF_IN_AUTOMATION=1

# prints "<source> <version>" for each provider version the tests will ask for
needed_providers() {
  # the providers every acceptance test uses alongside azurerm
  awk -F'"' '
    /VersionConstraint:/ { version = ""; if ($2 ~ /^=[0-9.]+$/) version = substr($2, 2) }
    /Source:/ { if (version != "" && $2 ~ /^registry\.terraform\.io\//) print $2, version; version = "" }
  ' internal/acceptance/testcase.go

  # the released azurerm which regression tests start from - the current release unless the test pins another
  {
    tr -d '[:space:]' < version/VERSION
    echo
    grep -rhoE --include='*_test.go' '\}, "[0-9]+\.[0-9]+\.[0-9]+"\)$' "$SERVICE_PATH" | grep -oE '[0-9]+\.[0-9]+\.[0-9]+'
  } | sort -u | sed 's|^|registry.terraform.io/hashicorp/azurerm |'
}

NEEDED="$(needed_providers)"

# mark what this build needs as in use, then drop the versions nothing has needed recently
while read -r source version; do
  if [ -n "$version" ] && [ -d "$MIRROR_DIR/$source/$version" ]; then
    touch "$MIRROR_DIR/$source/$version"
  fi
done <<< "$NEEDED"
find "$MIRROR_DIR" -mindepth 4 -maxdepth 4 -type d -mtime +"$RETENTION_DAYS" -exec rm -rf {} +

# left behind by interrupted builds
rm -rf "$CACHE_DIR"/.staging.*

while read -r source version; do
  if [ -z "$version" ]; then
    continue
  fi

  if [ -d "$MIRROR_DIR/$source/$version" ]; then
    echo "$source $version is already on this agent."
    continue
  fi

  echo "Fetching $source $version.."

  # fetched elsewhere and moved into place, so an interrupted build never leaves a partial provider behind
  staging="$(mktemp -d "$CACHE_DIR/.staging.XXXXXX")" || continue
  mkdir "$staging/cache" "$staging/config"

  cat > "$staging/config/main.tf" <<EOF
terraform {
  required_providers {
    $(basename "$source") = {
      source  = "$source"
      version = "=$version"
    }
  }
}
EOF

  if TF_PLUGIN_CACHE_DIR="$staging/cache" "$TERRAFORM" -chdir="$staging/config" init -backend=false -input=false -no-color 2>&1 \
    && mkdir -p "$MIRROR_DIR/$source" \
    && mv "$staging/cache/$source/$version" "$MIRROR_DIR/$source/$version"; then
    echo "Fetched $source $version."
  else
    echo "Could not fetch $source $version - the tests which need it will download it themselves."
  fi

  rm -rf "$staging"
done <<< "$NEEDED"

exit 0

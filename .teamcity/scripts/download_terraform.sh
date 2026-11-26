#!/bin/bash

# Keeps one copy of each Terraform Core version on the agent, shared by every build, rather than
# downloading it into each build's checkout.

set -e

TERRAFORM_VERSION="%env.TERRAFORM_CORE_VERSION%"
TERRAFORM_PATH="%env.TF_ACC_TERRAFORM_PATH%"
# versions which no build has used for this many days are deleted
RETENTION_DAYS=14

VERSION_DIR="$(dirname "$TERRAFORM_PATH")"
CACHE_DIR="$(dirname "$VERSION_DIR")"

mkdir -p "$CACHE_DIR"

# mark this version as in use, then drop the versions nothing has used recently
if [ -x "$TERRAFORM_PATH" ]; then
  touch "$VERSION_DIR"
fi
find "$CACHE_DIR" -mindepth 1 -maxdepth 1 -type d -atime +"$RETENTION_DAYS" -exec rm -rf {} +

if [ -x "$TERRAFORM_PATH" ]; then
  echo "Terraform Core v$TERRAFORM_VERSION is already on this agent."
  exit 0
fi

# downloaded elsewhere and moved into place, so an interrupted build never leaves a partial binary behind
DOWNLOAD_DIR="$(mktemp -d "$CACHE_DIR/.download.XXXXXX")"
trap 'rm -rf "$DOWNLOAD_DIR"' EXIT

# https://releases.hashicorp.com/terraform/0.12.28/terraform_0.12.28_linux_amd64.zip
wget -nv -O"$DOWNLOAD_DIR/tf.zip" "https://releases.hashicorp.com/terraform/$TERRAFORM_VERSION/terraform_${TERRAFORM_VERSION}_linux_amd64.zip"
unzip "$DOWNLOAD_DIR/tf.zip" terraform -d "$DOWNLOAD_DIR"
mkdir -p "$VERSION_DIR"
mv "$DOWNLOAD_DIR/terraform" "$TERRAFORM_PATH"

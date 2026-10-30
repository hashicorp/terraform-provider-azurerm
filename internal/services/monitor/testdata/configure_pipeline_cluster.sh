#!/bin/bash
# Copyright IBM Corp. 2014, 2026
# SPDX-License-Identifier: MPL-2.0

# Installs kubectl if missing, and works around a cert-manager extension issue: its
# default ClusterIssuers reference "*-current"-suffixed CA secrets that nothing creates
# automatically, which otherwise leaves the pipeline controller's own TLS certificate stuck
# "Issuing" forever and the pipelineGroups ARM resource stuck at provisioningState "Creating".
#
# Also creates the Kubernetes prerequisites the pipeline group's tlsConfigurations and
# persistent volume reference: the mTLS secrets and an Azure Files-backed PersistentVolume.
#
# KUBECONFIG and the other environment variables below are supplied by the calling
# local-exec provisioner's environment.

set -euo pipefail

mkdir -p "$HOME/go/bin"
export PATH="$HOME/go/bin:$PATH"

command -v kubectl >/dev/null || (
  case "$(uname -s)" in
    Linux) operating_system=linux ;;
    Darwin) operating_system=darwin ;;
    *) echo "Unsupported operating system for kubectl: $(uname -s)" >&2; exit 1 ;;
  esac

  case "$(uname -m)" in
    x86_64|amd64) architecture=amd64 ;;
    arm64|aarch64) architecture=arm64 ;;
    *) echo "Unsupported architecture for kubectl: $(uname -m)" >&2; exit 1 ;;
  esac

  version=$(curl -fsSL https://dl.k8s.io/release/stable.txt)
  curl -fsSL -o "$HOME/go/bin/kubectl" "https://dl.k8s.io/release/$version/bin/$operating_system/$architecture/kubectl"
  chmod +x "$HOME/go/bin/kubectl"
)

if ! kubectl version --client; then
  echo "kubectl is not runnable; remove the incompatible binary from PATH and install kubectl for this operating system and architecture" >&2
  exit 1
fi

kubectl wait --for=condition=Ready certificate/arc-amp-root-ca -n cert-manager --timeout=300s
kubectl wait --for=condition=Ready certificate/arc-amp-client-root-ca -n cert-manager --timeout=300s

for pair in 'arc-amp-root-ca arc-amp-root-ca-current' 'arc-amp-client-root-ca arc-amp-client-root-ca-current'; do
  set -- $pair
  kubectl get secret -n cert-manager "$1" -o json | python3 -c "import json,sys; d=json.load(sys.stdin); d['metadata']={'name':'$2','namespace':'cert-manager'}; print(json.dumps(d))" | kubectl apply -f -
done

kubectl create namespace "$PIPELINE_NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -

kubectl create secret generic client-ca-bundle \
  --namespace "$PIPELINE_NAMESPACE" \
  --from-literal=ca.crt="$CLIENT_CA_CERT" \
  --dry-run=client -o yaml | kubectl apply -f -

kubectl create secret tls pipeline-tls-cert \
  --namespace "$PIPELINE_NAMESPACE" \
  --cert=<(printf '%s' "$PIPELINE_TLS_CERT") \
  --key=<(printf '%s' "$PIPELINE_TLS_KEY") \
  --dry-run=client -o yaml | kubectl apply -f -

kubectl create secret generic azure-file-secret \
  --namespace "$PIPELINE_NAMESPACE" \
  --from-literal=azurestorageaccountname="$STORAGE_ACCOUNT" \
  --from-literal=azurestorageaccountkey="$STORAGE_ACCOUNT_KEY" \
  --dry-run=client -o yaml | kubectl apply -f -

kubectl apply -f "$PIPELINE_PERSISTENT_VOLUME_CONFIG"

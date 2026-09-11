#!/bin/bash

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

set -e

mkdir -p "$HOME/go/bin"
export PATH="$HOME/go/bin:$PATH"

command -v kubectl >/dev/null || (
  curl -sL -o "$HOME/go/bin/kubectl" "https://dl.k8s.io/release/$(curl -sL https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
  chmod +x "$HOME/go/bin/kubectl"
)

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


#!/bin/bash

set -e

# This script creates nessesary certificates to execute integration
# tests for Azure automation account certificates

KEY_FILE="automation_certificate_test.key"
CERT_FILE="automation_certificate_test.crt"
CERT_PFX_FILE="automation_certificate_test.pfx"
CERT_THUMBPRINT_FILE="automation_certificate_test.thumb"

openssl req \
  -x509 \
  -nodes \
  -sha256 \
  -days 3650 \
  -subj "/CN=Local" \
  -newkey rsa:2048 \
  -keyout ${KEY_FILE} \
  -out ${CERT_FILE}

openssl pkcs12 \
  -export \
  -in ${CERT_FILE} \
  -inkey ${KEY_FILE} \
  -passout pass: \
  -CSP "Microsoft Enhanced RSA and AES Cryptographic Provider" \
  -out $CERT_PFX_FILE

openssl x509 -in ${CERT_FILE} -fingerprint -noout | cut -d = -f 2 | sed 's/://g' > ${CERT_THUMBPRINT_FILE}
#!/bin/bash

# This script creates nessesary certificates to execute integration
# tests for Azure API Management

# api_management_api_test.pfx
openssl req \
    -newkey rsa:2048 \
    -x509 \
    -nodes \
    -keyout api_management_api_test.key \
    -new \
    -out api_management_api_test.crt \
    -subj /CN=api.terraform.io \
    -sha256 \
    -days 3650

openssl pkcs12 \
    -export \
    -out api_management_api_test.pfx \
    -inkey api_management_api_test.key \
    -in api_management_api_test.crt \
    -password pass:terraform

rm -f api_management_api_test.key api_management_api_test.crt

# api_management_api2_test.pfx
openssl req \
    -newkey rsa:2048 \
    -x509 \
    -nodes \
    -keyout api_management_api2_test.key \
    -new \
    -out api_management_api2_test.crt \
    -subj /CN=api2.terraform.io \
    -sha256 \
    -days 3650

openssl pkcs12 \
    -export \
    -out api_management_api2_test.pfx \
    -inkey api_management_api2_test.key \
    -in api_management_api2_test.crt \
    -password pass:terraform

rm -f api_management_api2_test.key api_management_api2_test.crt

# api_management_portal_test.pfx
openssl req \
    -newkey rsa:2048 \
    -x509 \
    -nodes \
    -keyout api_management_portal_test.key \
    -new \
    -out api_management_portal_test.crt \
    -subj /CN=portal.terraform.io \
    -sha256 \
    -days 3650

openssl pkcs12 \
    -export \
    -out api_management_portal_test.pfx \
    -inkey api_management_portal_test.key \
    -in api_management_portal_test.crt \
    -password pass:terraform

rm -f api_management_portal_test.key api_management_portal_test.crt

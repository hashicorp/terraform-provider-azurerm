#!/bin/bash

# This script creates nessesary certificates to execute integration
# tests for Azure API Management

declare -a certs=("api_management_api" "api_management_api2" "api_management_portal")
declare -a domains=("api.terraform.io" "api2.terraform.io" "portal.terraform.io")

arraylength=${#certs[@]}

for (( i=0; i<${arraylength}; i++ ));
do
  openssl req \
      -newkey rsa:2048 \
      -x509 \
      -nodes \
      -keyout ${certs[$i]}_test.key \
      -new \
      -out ${certs[$i]}_test.crt \
      -subj /CN=${domains[$i]} \
      -sha256 \
      -days 3650

  openssl pkcs12 \
      -export \
      -out ${certs[$i]}_test.pfx \
      -inkey ${certs[$i]}_test.key \
      -in ${certs[$i]}_test.crt \
      -password pass:terraform

  rm -f ${certs[$i]}_test.key ${certs[$i]}_test.crt
done

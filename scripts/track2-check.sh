#!/usr/bin/env bash

function checkForAzCoreUsages {
  result=$(grep -R "github.com/Azure/azure-sdk-for-go/sdk/azcore" go.mod go.sum)
  if [ "$result" != "" ];
  then
    echo "Detected usages of AzCore of the Azure SDK for Go"
    echo ""
    echo "At this time Terraform makes use of both Track1 of the Azure SDK"
    echo "for Go and some Embedded SDK's - however for various reasons we"
    echo "do not make use of AzCore and Track2 of the Azure SDK for Go."
    echo ""
    echo "We've detected a usage of AzCore in the go.mod/go.sum with the"
    echo "import path of github.com/Azure/azure-sdk-for-go/sdk/azcore"
    echo "which is likely coming from a dependency, either from Track2 of"
    echo "the Azure SDK for Go - or another Azure SDK library."
    echo ""
    echo "Rather than importing an SDK which has a reliance on the Track2"
    echo "libraries, please use either github.com/hashicorp/go-azure-sdk"
    echo "(please open an issue on that repository if you need a"
    echo "Service/API Version which isn't supported) - or the 'Track1'"
    echo "Azure SDK for Go."
    exit 1
  fi
}

function checkForAzIdentityUsages {
  result=$(grep -R "github.com/Azure/azure-sdk-for-go/sdk/azidentity" go.mod go.sum)
  if [ "$result" != "" ];
  then
    echo "Detected usages of AzIdentity of the Azure SDK for Go"
    echo ""
    echo "At this time Terraform makes use of both Track1 of the Azure SDK"
    echo "for Go and some Embedded SDK's - however for various reasons we"
    echo "do not make use of AzIdentity and Track2 of the Azure SDK for Go."
    echo ""
    echo "We've detected a usage of AzIdentity in the go.mod/go.sum with the"
    echo "import path of github.com/Azure/azure-sdk-for-go/sdk/azidentity"
    echo "which is likely coming from a dependency, either from Track2 of"
    echo "the Azure SDK for Go - or another Azure SDK library."
    echo ""
    echo "Rather than importing an SDK which has a reliance on the Track2"
    echo "libraries, please use either github.com/hashicorp/go-azure-sdk"
    echo "(please open an issue on that repository if you need a"
    echo "Service/API Version which isn't supported) - or the 'Track1'"
    echo "Azure SDK for Go."
    exit 1
  fi
}

function main {
  checkForAzCoreUsages
  checkForAzIdentityUsages
}

main

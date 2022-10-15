#! /bin/bash
GIT_REPO_LOCATION=~/git
DATAPLANE_LOCATION="${GIT_REPO_LOCATION}/azure-rest-api-specs/specification/keyvault/data-plane"
AZURERM_LOCATION="${GIT_REPO_LOCATION}/terraform-provider-azurerm"
SDK_LOCATION="${AZURERM_LOCATION}/internal/services/keyvault/sdk/v7.3/keyvault"

mkdir -p "${SDK_LOCATION}"
cp "${SDK_LOCATION}/../readme.go.md" "${DATAPLANE_LOCATION}"/readme.go.md
cp "${SDK_LOCATION}/../readme.md" "${DATAPLANE_LOCATION}"/readme.md 
cd "${DATAPLANE_LOCATION}/" || exit 1

# Install autorest before running this
autorest --use=@microsoft.azure/autorest.go@2.1.183 --tag=package-7.3 --go --openapi-type=data-plane --use-onever --version=V2 --go-sdk-folder="${SDK_LOCATION}"

# Remove the interfaces
rm -r "${SDK_LOCATION}/keyvaultapi"

# Format seems off ..
cd "${AZURERM_LOCATION}" || exit 1
make fmt
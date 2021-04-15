#!/usr/bin/env bash

PACKAGE_NAME="$1"
SERVICE_DIRECTORY="./azurerm/internal/services/$PACKAGE_NAME"

function updateTests {
  echo "Updating Tests.."

  # references to `resource.` should be `acceptance.`
  find "$SERVICE_DIRECTORY" -type f -iname "*_test.go" -exec sed -i '' -e 's/github\.com\/hashicorp\/terraform-plugin-sdk\/helper\/resource/github\.com\/terraform-providers\/terraform-provider-azurerm\/azurerm\/internal\/acceptance/g' {} \;
  find "$SERVICE_DIRECTORY" -type f -iname "*_test.go" -exec sed -i '' -e 's/resource\./acceptance\./g' {} \;

  # references to `terraform.` should be `pluginsdk.`
  find "$SERVICE_DIRECTORY" -type f -iname "*_test.go" -exec sed -i '' -e 's/github\.com\/hashicorp\/terraform-plugin-sdk\/terraform/github\.com\/terraform-providers\/terraform-provider-azurerm\/azurerm\/internal\/tf\/pluginsdk/g' {} \;
  find "$SERVICE_DIRECTORY" -type f -iname "*_test.go" -exec sed -i '' -e 's/terraform\./pluginsdk\./g' {} \;
}

function updateCode {
  echo "Updating Resources.."
  # references to `schema.` should be `pluginsdk.`
  find "$SERVICE_DIRECTORY" -type f -exec sed -i '' -e 's/github\.com\/hashicorp\/terraform-plugin-sdk\/helper\/schema/github\.com\/terraform-providers\/terraform-provider-azurerm\/azurerm\/internal\/tf\/pluginsdk/g' {} \;
  find "$SERVICE_DIRECTORY" -type f -exec sed -i '' -e 's/schema\./pluginsdk\./g' {} \;

  # update the import for validation, the alised/unaliased part can remain the same
  find "$SERVICE_DIRECTORY" -type f -exec sed -i '' -e 's/github\.com\/hashicorp\/terraform-plugin-sdk\/helper\/validation/github\.com\/terraform-providers\/terraform-provider-azurerm\/azurerm\/internal\/tf\/validation/g' {} \;
}

function format {
  echo "Running goimports.."
  goimports -w "$SERVICE_DIRECTORY"
}

updateTests
updateCode
format
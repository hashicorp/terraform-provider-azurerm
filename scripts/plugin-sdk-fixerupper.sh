#!/usr/bin/env bash

PACKAGE_NAME="$1"
SERVICE_DIRECTORY="./azurerm/internal/services/$PACKAGE_NAME"

function updateTests {
  echo "Updating Tests.."

  # references to `resource.` should be `acceptance.`
  find "$SERVICE_DIRECTORY" -type f -iname "*_test.go" -exec sed -i '' -e 's/github\.com\/hashicorp\/terraform-plugin-sdk\/helper\/resource/github\.com\/terraform-providers\/terraform-provider-azurerm\/azurerm\/internal\/acceptance/g' {} \;
  find "$SERVICE_DIRECTORY" -type f -iname "*_test.go" -exec sed -i '' -e 's/resource\./acceptance\./g' {} \;

  # references to `acctest.` should be `acceptance.`
  find "$SERVICE_DIRECTORY" -type f -iname "*_test.go" -exec sed -i '' -e 's/github\.com\/hashicorp\/terraform-plugin-sdk\/helper\/acctest/github\.com\/terraform-providers\/terraform-provider-azurerm\/azurerm\/internal\/acceptance/g' {} \;
  find "$SERVICE_DIRECTORY" -type f -iname "*_test.go" -exec sed -i '' -e 's/acctest\./acceptance\./g' {} \;

  # references to `terraform.` should be `pluginsdk.`
  find "$SERVICE_DIRECTORY" -type f -iname "*_test.go" -exec sed -i '' -e 's/github\.com\/hashicorp\/terraform-plugin-sdk\/terraform/github\.com\/terraform-providers\/terraform-provider-azurerm\/azurerm\/internal\/tf\/pluginsdk/g' {} \;
  find "$SERVICE_DIRECTORY" -type f -iname "*_test.go" -exec sed -i '' -e 's/terraform\./pluginsdk\./g' {} \;
}

function updateCode {
  echo "Updating Resources.."
  # references to `schema.` should be `pluginsdk.`
  find "$SERVICE_DIRECTORY" -type f -iname "*.go"  -exec sed -i '' -e 's/github\.com\/hashicorp\/terraform-plugin-sdk\/helper\/schema/github\.com\/terraform-providers\/terraform-provider-azurerm\/azurerm\/internal\/tf\/pluginsdk/g' {} \;
  find "$SERVICE_DIRECTORY" -type f -iname "*.go"  -exec sed -i '' -e 's/schema\./pluginsdk\./g' {} \;

  # references to `structure.` should be `pluginsdk.`
  find "$SERVICE_DIRECTORY" -type f -iname "*.go"  -exec sed -i '' -e 's/github\.com\/hashicorp\/terraform-plugin-sdk\/helper\/structure/github\.com\/terraform-providers\/terraform-provider-azurerm\/azurerm\/internal\/tf\/pluginsdk/g' {} \;
  find "$SERVICE_DIRECTORY" -type f -iname "*.go"  -exec sed -i '' -e 's/structure\./pluginsdk\./g' {} \;

  # update the import for validation, the alised/unaliased part can remain the same
  find "$SERVICE_DIRECTORY" -type f -iname "*.go"  -exec sed -i '' -e 's/github\.com\/hashicorp\/terraform-plugin-sdk\/helper\/validation/github\.com\/terraform-providers\/terraform-provider-azurerm\/azurerm\/internal\/tf\/validation/g' {} \;

  # `resource.` can become `pluginsdk`.
  find "$SERVICE_DIRECTORY" -type f -iname "*.go"  -exec sed -i '' -e 's/github\.com\/hashicorp\/terraform-plugin-sdk\/helper\/resource/github\.com\/terraform-providers\/terraform-provider-azurerm\/azurerm\/internal\/tf\/pluginsdk/g' {} \;
  find "$SERVICE_DIRECTORY" -type f -iname "*.go"  -exec sed -i '' -e 's/resource\./pluginsdk\./g' {} \;
  find "$SERVICE_DIRECTORY" -type f -iname "*.go"  -exec sed -i '' -e 's/\&resource\./\&pluginsdk\./g' {} \;

  # import functions should use the pluginsdk function
  find "$SERVICE_DIRECTORY" -type f -iname "*.go"  -exec sed -i '' -e 's/azSchema\.ValidateResourceIDPriorToImportThen/pluginsdk\.ImporterValidatingResourceIdThen/g' {} \;
  find "$SERVICE_DIRECTORY" -type f -iname "*.go"  -exec sed -i '' -e 's/azSchema\.ValidateResourceIDPriorToImport/pluginsdk\.ImporterValidatingResourceId/g' {} \;
}

function revertValidationFuncs {
  echo "Reverting unintentional changes to client functions.."
  git checkout -- "$SERVICE_DIRECTORY/client" &> /dev/null
  echo "Reverting unintentional changes to parse functions.."
  git checkout -- "$SERVICE_DIRECTORY/parse" &> /dev/null
}

function format {
  echo "Running fmt.."
  gofmt -w "$SERVICE_DIRECTORY"
  echo "Running goimports.."
  goimports -w "$SERVICE_DIRECTORY"
}

updateTests
updateCode
revertValidationFuncs
format

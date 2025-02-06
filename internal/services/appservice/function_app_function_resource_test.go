// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type FunctionAppFunctionResource struct{}

func TestAccFunctionAppFunction_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_function", "test")
	r := FunctionAppFunctionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("language"),
	})
}

func TestAccFunctionAppFunction_basicUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_function", "test")
	r := FunctionAppFunctionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("language"),
		{
			Config: r.basicUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("language"),
	})
}

func TestAccFunctionAppFunction_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_function", "test")
	r := FunctionAppFunctionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccFunctionAppFunction_withFiles(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_function", "test")
	r := FunctionAppFunctionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withLocalFiles(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("language", "file"),
	})
}

func (r FunctionAppFunctionResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := webapps.ParseFunctionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.AppService.WebAppsClient.GetFunction(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if response.WasNotFound(resp.HttpResponse) {
		return utils.Bool(false), nil
	}

	return utils.Bool(true), nil
}

func (r FunctionAppFunctionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_function" "test" {
  name            = "testAcc-FnAppFn-%[2]d"
  function_app_id = azurerm_linux_function_app.test.id
  language        = "Python"
  test_data = jsonencode({
    "name" = "Azure"
  })
  config_json = jsonencode({
    "bindings" = [
      {
        "authLevel" = "function"
        "direction" = "in"
        "methods" = [
          "get",
          "post",
        ]
        "name" = "req"
        "type" = "httpTrigger"
      },
      {
        "direction" = "out"
        "name"      = "$return"
        "type"      = "http"
      },
    ]
  })
}
`, r.templateLinux(data), data.RandomInteger)
}

func (r FunctionAppFunctionResource) basicUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_function" "test" {
  name            = "testAcc-FnAppFn-%[2]d"
  function_app_id = azurerm_linux_function_app.test.id
  language        = "Python"
  test_data = jsonencode({
    "name" = "AzureRM"
  })
  config_json = jsonencode({
    "bindings" = [
      {
        "authLevel" = "function"
        "direction" = "in"
        "methods" = [
          "get",
          "post",
          "put",
        ]
        "name" = "req"
        "type" = "httpTrigger"
      },
      {
        "direction" = "out"
        "name"      = "$return"
        "type"      = "http"
      },
    ]
  })
}
`, r.templateLinux(data), data.RandomInteger)
}

func (r FunctionAppFunctionResource) withLocalFiles(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_function" "test" {
  name            = "testAcc-FnAppFn-%[2]d"
  function_app_id = azurerm_windows_function_app.test.id
  language        = "CSharp"
  file {
    name    = "run.csx"
    content = file("testdata/run.csx")
  }
  test_data = jsonencode({
    "name" = "Azure"
  })
  config_json = jsonencode({
    "bindings" = [
      {
        "authLevel" = "function"
        "direction" = "in"
        "methods" = [
          "get",
          "post",
        ]
        "name" = "req"
        "type" = "httpTrigger"
      },
      {
        "direction" = "out"
        "name"      = "$return"
        "type"      = "http"
      },
    ]
  })
}
`, r.templateWindows(data), data.RandomInteger)
}

func (r FunctionAppFunctionResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_function_app_function" "import" {
  name            = azurerm_function_app_function.test.name
  function_app_id = azurerm_function_app_function.test.function_app_id
  language        = azurerm_function_app_function.test.language
  test_data       = azurerm_function_app_function.test.test_data
  config_json     = azurerm_function_app_function.test.config_json
}

`, r.basic(data))
}

func (r FunctionAppFunctionResource) templateLinux(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-LFA-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  os_type             = "Linux"
  sku_name            = "S1"
}

resource "azurerm_linux_function_app" "test" {
  name                = "acctest-WFA-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    application_stack {
      python_version = "3.9"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r FunctionAppFunctionResource) templateWindows(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-LFA-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  os_type             = "Windows"
  sku_name            = "S1"
}

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    application_stack {
      dotnet_version = "v6.0"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

package appservice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AppFunctionResource struct{}

func TestAccAppFunction_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_function", "test")
	r := AppFunctionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r AppFunctionResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.AppFunctionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.AppService.WebAppsClient.GetFunction(ctx, id.ResourceGroup, id.SiteName, id.FunctionName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if utils.ResponseWasNotFound(resp.Response) {
		return utils.Bool(false), nil
	}

	return utils.Bool(true), nil
}

func (r AppFunctionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_app_function" "test" {
  name            = "testAcc-FnAppFn-%[2]d"
  function_app_id = "/subscriptions/1a6092a6-137e-4025-9a7c-ef77f76f2c02/resourceGroups/stedev-app-functions-beta/providers/Microsoft.Web/sites/stedev-win-app-function-beta"
  language        = "CSharp"
  test_data       = jsonencode({
  "name" = "Azure"
})
  config_json     = jsonencode({
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
      "name" = "$return"
      "type" = "http"
    },
  ]
})
}

`, r.template(data), data.RandomInteger)
}

func (r AppFunctionResource) template(data acceptance.TestData) string {
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

resource "azurerm_linux_function_app" "test" {
  name                = "acctest-WFA-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    application_stack {
      dotnet_version = "6"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

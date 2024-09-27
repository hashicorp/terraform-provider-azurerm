// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package web_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type FunctionAppHostKeysDataSource struct{}

func TestAccFunctionAppHostKeysDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_function_app_host_keys", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: FunctionAppHostKeysDataSource{}.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("primary_key").Exists(),
				check.That(data.ResourceName).Key("default_function_key").Exists(),
				check.That(data.ResourceName).Key("event_grid_extension_config_key").Exists(),
			),
		},
	})
}

func (d FunctionAppHostKeysDataSource) basic(data acceptance.TestData) string {
	template := FunctionAppResource{}.basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_function_app_host_keys" "test" {
  name                = azurerm_function_app.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}

func TestAccFunctionAppHostKeysDataSource_linuxEventGridTrigger(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_function_app_host_keys", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: FunctionAppHostKeysDataSource{}.linuxEventGridTrigger(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("primary_key").Exists(),
				check.That(data.ResourceName).Key("default_function_key").Exists(),
				check.That(data.ResourceName).Key("event_grid_extension_config_key").Exists(),
			),
		},
	})
}

func (d FunctionAppHostKeysDataSource) linuxEventGridTrigger(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
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

  os_type  = "Linux"
  sku_name = "EP1"
}

resource "azurerm_linux_function_app" "test" {
  name                       = "acctest-%[1]d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  service_plan_id            = azurerm_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  zip_deploy_file = abspath("testdata/test_trigger.zip")

  app_settings = {
    WEBSITE_RUN_FROM_PACKAGE = 1
  }

  identity {
    type = "SystemAssigned"
  }

  site_config {
    application_stack {
      python_version = "3.11"
    }
  }
}

// The key is not always present when azurerm_linux_function_app.test creation completes.
resource "time_sleep" "wait_for_event_grid_key" {
  depends_on = [azurerm_linux_function_app.test]

  create_duration = "30s"
}

data "azurerm_function_app_host_keys" "test" {
  depends_on = [time_sleep.wait_for_event_grid_key]

  name                = azurerm_linux_function_app.test.name
  resource_group_name = azurerm_linux_function_app.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

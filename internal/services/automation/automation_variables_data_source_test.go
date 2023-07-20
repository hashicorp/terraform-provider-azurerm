// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type AutomationVariablesDataSource struct{}

func TestAccDataSourceAzureRMAutomationVariables_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_automation_variables", "test")
	r := AutomationVariablesDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("bool.#").HasValue("3"),
				check.That(data.ResourceName).Key("bool.0.name").HasValue(fmt.Sprintf("acctestBoolVar-%d-0", data.RandomInteger)),
				check.That(data.ResourceName).Key("bool.0.value").HasValue("true"),
				check.That(data.ResourceName).Key("datetime.#").HasValue("3"),
				check.That(data.ResourceName).Key("datetime.0.name").HasValue(fmt.Sprintf("acctestDateTimeVar-%d-0", data.RandomInteger)),
				check.That(data.ResourceName).Key("datetime.0.value").HasValue("2019-04-20T08:40:04.02Z"),
				check.That(data.ResourceName).Key("encrypted.#").HasValue("3"),
				check.That(data.ResourceName).Key("encrypted.0.name").HasValue(fmt.Sprintf("acctestEncryptedVar-%d-0", data.RandomInteger)),
				check.That(data.ResourceName).Key("int.#").HasValue("3"),
				check.That(data.ResourceName).Key("int.0.name").HasValue(fmt.Sprintf("acctestIntVar-%d-0", data.RandomInteger)),
				check.That(data.ResourceName).Key("int.0.value").HasValue("0"),
				check.That(data.ResourceName).Key("string.#").HasValue("3"),
				check.That(data.ResourceName).Key("string.0.name").HasValue(fmt.Sprintf("acctestStringVar-%d-0", data.RandomInteger)),
				check.That(data.ResourceName).Key("string.0.value").HasValue("Hello, Terraform Variables Test."),
			),
		},
	})
}

func (AutomationVariablesDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-auto-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctestAutoAcct-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_automation_variable_bool" "test" {
  count = 3

  name                    = "acctestBoolVar-%d-${count.index}"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  description             = "This variable is created by Terraform acceptance test."
  value                   = true
}

resource "azurerm_automation_variable_datetime" "test" {
  count = 3

  name                    = "acctestDateTimeVar-%d-${count.index}"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  description             = "This variable is created by Terraform acceptance test."
  value                   = "2019-04-20T08:40:04.02Z"
}

resource "azurerm_automation_variable_int" "test" {
  count = 3

  name                    = "acctestIntVar-%d-${count.index}"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  description             = "This variable is created by Terraform acceptance test."
  value                   = count.index
}

resource "azurerm_automation_variable_string" "test" {
  count = 3

  name                    = "acctestStringVar-%d-${count.index}"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  description             = "This variable is created by Terraform acceptance test."
  value                   = "Hello, Terraform Variables Test."
}

resource "azurerm_automation_variable_string" "encrypted" {
  count = 3

  name                    = "acctestEncryptedVar-%d-${count.index}"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  description             = "This variable is created by Terraform acceptance test."
  value                   = "Hello, Terraform Variables Test."
  encrypted               = true
}

data "azurerm_automation_variables" "test" {
  automation_account_id = azurerm_automation_account.test.id

  depends_on = [
    azurerm_automation_variable_bool.test,
    azurerm_automation_variable_datetime.test,
    azurerm_automation_variable_int.test,
    azurerm_automation_variable_string.test,
    azurerm_automation_variable_string.encrypted,
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

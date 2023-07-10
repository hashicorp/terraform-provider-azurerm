// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type AutomationAccountDataSource struct{}

func TestAccDataSourceAutomationAccount(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_automation_account", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: AutomationAccountDataSource{}.complete(data),
			Check:  acceptance.ComposeTestCheckFunc(),
		},
	})
}

func (AutomationAccountDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-auto-%[1]d"
  location = "%[2]s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctestautomationAccount-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"

  identity {
    type = "SystemAssigned"
  }
}

data "azurerm_automation_account" "test" {
  resource_group_name = azurerm_resource_group.test.name
  name                = azurerm_automation_account.test.name
}

output "automation_account_system_managed_identity_principal_id" {
  value = data.azurerm_automation_account.test.identity[0].principal_id
}
`, data.RandomInteger, data.Locations.Primary)
}

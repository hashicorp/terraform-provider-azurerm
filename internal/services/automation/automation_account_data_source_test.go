// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type AutomationAccountDataSource struct{}

func TestAccDataSourceAutomationAccount_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_automation_account", "test")
	resource := acceptance.BuildTestData(t, "azurerm_automation_account", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: AutomationAccountDataSource{}.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").MatchesOtherKey(check.That(resource.ResourceName).Key("location")),
				check.That(data.ResourceName).Key("sku_name").MatchesOtherKey(check.That(resource.ResourceName).Key("sku_name")),
				check.That(data.ResourceName).Key("tags.%").MatchesOtherKey(check.That(resource.ResourceName).Key("tags.%")),
				check.That(data.ResourceName).Key("public_network_access_enabled").MatchesOtherKey(check.That(resource.ResourceName).Key("public_network_access_enabled")),
				check.That(data.ResourceName).Key("local_authentication_enabled").MatchesOtherKey(check.That(resource.ResourceName).Key("local_authentication_enabled")),
				check.That(data.ResourceName).Key("dsc_server_endpoint").MatchesOtherKey(check.That(resource.ResourceName).Key("dsc_server_endpoint")),
				check.That(data.ResourceName).Key("dsc_primary_access_key").MatchesOtherKey(check.That(resource.ResourceName).Key("dsc_primary_access_key")),
				check.That(data.ResourceName).Key("dsc_secondary_access_key").MatchesOtherKey(check.That(resource.ResourceName).Key("dsc_secondary_access_key")),
			),
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

tags = {
  "Hello" = "World"
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

func TestAccDataSourceAutomationAccount_encryption(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_automation_account", "test")
	resource := acceptance.BuildTestData(t, "azurerm_automation_account", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: AutomationAccountDataSource{}.encryption(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("encryption.#").MatchesOtherKey(check.That(resource.ResourceName).Key("encryption.#")),
			),
		},
	})
}

func (AutomationAccountDataSource) encryption(data acceptance.TestData) string {
	return fmt.Sprintf(`
  %s

  data "azurerm_automation_account" "test" {
    resource_group_name = azurerm_resource_group.test.name
    name                = azurerm_automation_account.test.name
  }
    `, AutomationAccountResource{}.encryption_userIdentity(data))
}

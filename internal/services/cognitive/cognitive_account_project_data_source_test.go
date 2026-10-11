// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type CognitiveAccountProjectDataSource struct{}

func TestAccCognitiveAccountProjectDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_cognitive_account_project", "test")
	r := CognitiveAccountProjectDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").IsUUID(),
				check.That(data.ResourceName).Key("default").Exists(),
				check.That(data.ResourceName).Key("endpoints.%").Exists(),
			),
		},
	})
}

func TestAccCognitiveAccountProjectDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_cognitive_account_project", "test")
	r := CognitiveAccountProjectDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("description").HasValue("Project Description"),
				check.That(data.ResourceName).Key("display_name").HasValue("Project Display Name"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned, UserAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").IsUUID(),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("default").Exists(),
				check.That(data.ResourceName).Key("endpoints.%").Exists(),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
				check.That(data.ResourceName).Key("tags.Environment").HasValue("foo"),
				check.That(data.ResourceName).Key("tags.Purpose").HasValue("AcceptanceTest"),
			),
		},
	})
}

func (CognitiveAccountProjectDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cognitive-%[1]d"
  location = "%[2]s"
}

resource "azurerm_cognitive_account" "test" {
  name                       = "acctest-cog-%[1]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  kind                       = "AIServices"
  sku_name                   = "S0"
  project_management_enabled = true
  custom_subdomain_name      = "acctest-cog-%[1]d"
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_cognitive_account_project" "test" {
  name                 = "acctest-cog-project-%[1]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  location             = azurerm_resource_group.test.location
  identity {
    type = "SystemAssigned"
  }
}

data "azurerm_cognitive_account_project" "test" {
  name                   = azurerm_cognitive_account_project.test.name
  cognitive_account_name = azurerm_cognitive_account.test.name
  resource_group_name    = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary)
}

func (CognitiveAccountProjectDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cognitive-%[1]d"
  location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "acctest-uai-%[1]d"
}

resource "azurerm_cognitive_account" "test" {
  name                       = "acctest-cog-%[1]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  kind                       = "AIServices"
  sku_name                   = "S0"
  project_management_enabled = true
  custom_subdomain_name      = "acctest-cog-%[1]d"
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_cognitive_account_project" "test" {
  name                 = "acctest-cog-project-%[1]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  location             = azurerm_resource_group.test.location
  description          = "Project Description"
  display_name         = "Project Display Name"

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  tags = {
    Environment = "foo"
    Purpose     = "AcceptanceTest"
  }
}

data "azurerm_cognitive_account_project" "test" {
  name                   = azurerm_cognitive_account_project.test.name
  cognitive_account_name = azurerm_cognitive_account.test.name
  resource_group_name    = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary)
}

// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package loadtestservice_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type PlaywrightWorkspaceDataSource struct{}

func TestAccPlaywrightWorkspaceDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_playwright_workspace", "test")
	r := PlaywrightWorkspaceDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").HasValue(location.Normalize(data.Locations.Primary)),
				check.That(data.ResourceName).Key("dataplane_uri").Exists(),
				check.That(data.ResourceName).Key("local_auth_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("regional_affinity_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("workspace_id").Exists(),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
				check.That(data.ResourceName).Key("tags.Environment").HasValue("Sandbox"),
				check.That(data.ResourceName).Key("tags.Label").HasValue("Test"),
			),
		},
	})
}

func (PlaywrightWorkspaceDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-pww-%d"
  location = "%s"
}

resource "azurerm_playwright_workspace" "test" {
  name                = "acctest-pww-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  tags = {
    Environment = "Sandbox"
    Label       = "Test"
  }
}

data "azurerm_playwright_workspace" "test" {
  name                = azurerm_playwright_workspace.test.name
  resource_group_name = azurerm_playwright_workspace.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(8))
}

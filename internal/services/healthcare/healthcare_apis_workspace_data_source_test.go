package healthcare_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type HealthCareWorkspaceDataSource struct{}

func TestAccHealthCareWorkspaceDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_healthcare_workspace", "test")
	r := HealthCareWorkspaceDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists()),
		},
	})
}

func (HealthCareWorkspaceDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-health-%d"
  location = "%s"
}

resource "azurerm_healthcare_workspace" "test" {
  name                = "acctestwk%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

data "azurerm_healthcare_workspace" "test" {
  name                = azurerm_healthcare_workspace.test.name
  resource_group_name = azurerm_healthcare_workspace.test.resource_group_name
  location            = azurerm_healthcare_workspace.test.location
}`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(8))
}

package apimanagement_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ApiManagementWorkspaceDataSource struct{}

func TestAccDataSourceApiManagementWorkspace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_api_management_workspace", "test")
	r := ApiManagementWorkspaceDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("workspace_name").Exists(),
			),
		},
	})
}

func (r ApiManagementWorkspaceDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_api_management_workspace" "test" {
  name                = azurerm_api_management_workspace.test.name
  service_name        = azurerm_api_management_workspace.test.service_name
  resource_group_name = azurerm_api_management_workspace.test.resource_group_name
}
`, r.template(data))
}

func (r ApiManagementWorkspaceDataSource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-apim-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name           = "Developer_1"
}

resource "azurerm_api_management_workspace" "test" {
  name                = "acctestws%d"
  service_name        = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  workspace_name      = "acctest-workspace-%d"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

package databricks_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type DatabricksWorkspaceDataSource struct {
}

func TestAccDatabricksWorkspaceDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_databricks_workspace", "test")
	r := DatabricksWorkspaceDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				acceptance.TestMatchResourceAttr(data.ResourceName, "workspace_url", regexp.MustCompile("azuredatabricks.net")),
				check.That(data.ResourceName).Key("workspace_id").Exists(),
			),
		},
	})
}

func (DatabricksWorkspaceDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-databricks-%d"
  location = "%s"
}

resource "azurerm_databricks_workspace" "test" {
  name                = "acctestDBW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "standard"
}

data "azurerm_databricks_workspace" "test" {
  name                = azurerm_databricks_workspace.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

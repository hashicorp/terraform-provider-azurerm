package databricks_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDatabricksWorkspaceDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_databricks_workspace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDatabricksWorkspace_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(data.ResourceName, "workspace_url", regexp.MustCompile("azuredatabricks.net")),
					resource.TestCheckResourceAttrSet(data.ResourceName, "workspace_id"),
				),
			},
		},
	})
}

func testAccDataSourceDatabricksWorkspace_basic(data acceptance.TestData) string {
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

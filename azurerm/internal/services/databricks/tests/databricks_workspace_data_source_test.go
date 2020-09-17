package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMDatabricksWorkspace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_databricks_workspace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDatabricksWorkspace_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "sku", "premium"),
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

resource "azurerm_resource_group" "example" {
  name     = "acctestRG-databricks-%[1]d"
  location = "%[2]s"
}

resource "azurerm_databricks_workspace" "example" {
  name     = "acctestRG-databricks-workspace-%[1]d"
  location = azurerm_resource_group.example.location
  sku = "premium"
}

data "azurerm_databricks_workspace" "example" {
  name                = azurerm_databricks_workspace.example.name
  resource_group_name = azurerm_resource_group.example.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

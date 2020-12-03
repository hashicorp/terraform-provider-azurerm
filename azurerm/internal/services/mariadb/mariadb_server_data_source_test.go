package mariadb_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMMariaDbServer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_mariadb_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMariaDbServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMMariaDbServer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMariaDbServerExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "administrator_login", "acctestun"),
					resource.TestCheckResourceAttr(data.ResourceName, "version", "10.2"),
					resource.TestCheckResourceAttr(data.ResourceName, "ssl_enforcement", "Enabled"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMMariaDbServer_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-maria-%d"
  location = "%s"
}

resource "azurerm_mariadb_server" "test" {
  name                = "acctestmariadbsvr-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "B_Gen5_2"

  storage_profile {
    storage_mb            = 51200
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "acctestun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "10.2"
  ssl_enforcement_enabled      = true
}

data "azurerm_mariadb_server" "test" {
  name                = azurerm_mariadb_server.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

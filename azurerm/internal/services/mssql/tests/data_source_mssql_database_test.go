package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMMsSqlDatabase_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_mssql_database", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMMsSqlDatabase_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMMsSqlDatabase_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssqldb-%[1]d"
  location = "%s"
}

resource "azurerm_sql_server" "test" {
  name                         = "acctest%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "4dm1n157r470r"
  administrator_login_password = "4-v3ry-53cr37-p455w0rd"
}

resource "azurerm_mssql_database" "test" {
  name            = "acctest-db-%[1]d"
  mssql_server_id = azurerm_sql_server.test.id
}

data "azurerm_mssql_database" "test" {
  name            = azurerm_mssql_database.test.name
  mssql_server_id = azurerm_sql_server.test.id
}

`, data.RandomInteger, data.Locations.Primary)
}

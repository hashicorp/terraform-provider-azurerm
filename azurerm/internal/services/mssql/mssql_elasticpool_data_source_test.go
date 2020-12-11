package mssql_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMMsSqlElasticPool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_mssql_elasticpool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlElasticPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMMsSqlElasticPool_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlElasticPoolExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "server_name"),
					resource.TestCheckResourceAttr(data.ResourceName, "location", data.Locations.Primary),
					resource.TestCheckResourceAttr(data.ResourceName, "max_size_gb", "50"),
					resource.TestCheckResourceAttr(data.ResourceName, "per_db_min_capacity", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "per_db_max_capacity", "4"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "zone_redundant", "false"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMMsSqlElasticPool_licenseType(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_mssql_elasticpool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlElasticPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMMsSqlElasticPool_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlElasticPoolExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "license_type", "LicenseIncluded"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMMsSqlElasticPool_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
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

resource "azurerm_mssql_elasticpool" "test" {
  name                = "acctest-pool-dtu-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  server_name         = azurerm_sql_server.test.name
  max_size_gb         = 50
  zone_redundant      = false

  sku {
    name     = "GP_Gen5"
    tier     = "GeneralPurpose"
    capacity = 4
    family   = "Gen5"
  }

  per_database_settings {
    min_capacity = 0
    max_capacity = 4
  }
}

data "azurerm_mssql_elasticpool" "test" {
  name                = azurerm_mssql_elasticpool.test.name
  resource_group_name = azurerm_resource_group.test.name
  server_name         = azurerm_sql_server.test.name
}
`, data.RandomInteger, data.Locations.Primary)
}

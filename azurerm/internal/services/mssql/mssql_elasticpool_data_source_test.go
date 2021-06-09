package mssql_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type MsSqlElasticPoolDataSource struct{}

func TestAccDataSourceMsSqlElasticPool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_mssql_elasticpool", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: MsSqlElasticPoolDataSource{}.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("server_name").Exists(),
				check.That(data.ResourceName).Key("location").HasValue(data.Locations.Primary),
				check.That(data.ResourceName).Key("license_type").HasValue("LicenseIncluded"),
				check.That(data.ResourceName).Key("max_size_gb").HasValue("50"),
				check.That(data.ResourceName).Key("per_db_min_capacity").HasValue("0"),
				check.That(data.ResourceName).Key("per_db_max_capacity").HasValue("4"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("zone_redundant").HasValue("false"),
			),
		},
	})
}

func (MsSqlElasticPoolDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
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

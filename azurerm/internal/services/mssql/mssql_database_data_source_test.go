package mssql_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type MsSqlDatabaseDataSource struct{}

func TestAccDataSourceMsSqlDatabase_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_mssql_database", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: MsSqlDatabaseDataSource{}.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctest-db-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("server_id").Exists(),
			),
		},
	})
}

func TestAccDataSourceMsSqlDatabase_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_mssql_database", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: MsSqlDatabaseDataSource{}.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctest-db-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("server_id").Exists(),
				check.That(data.ResourceName).Key("collation").HasValue("SQL_AltDiction_CP850_CI_AI"),
				check.That(data.ResourceName).Key("license_type").HasValue("BasePrice"),
				check.That(data.ResourceName).Key("max_size_gb").HasValue("1"),
				check.That(data.ResourceName).Key("sku_name").HasValue("GP_Gen5_2"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.ENV").HasValue("Test"),
			),
		},
	})
}

func (MsSqlDatabaseDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

data "azurerm_mssql_database" "test" {
  name      = azurerm_mssql_database.test.name
  server_id = azurerm_sql_server.test.id
}

`, MsSqlDatabaseResource{}.basic(data))
}

func (MsSqlDatabaseDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

data "azurerm_mssql_database" "test" {
  name      = azurerm_mssql_database.test.name
  server_id = azurerm_sql_server.test.id
}

`, MsSqlDatabaseResource{}.complete(data))
}

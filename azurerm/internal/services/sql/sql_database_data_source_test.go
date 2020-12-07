package sql_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type SqlDatabaseDataSource struct{}

func TestAccDataSourceSqlDatabase_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_sql_database", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: SqlDatabaseDataSource{}.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("collation").Exists(),
				check.That(data.ResourceName).Key("edition").Exists(),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("read_scale").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("server_name").Exists(),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func TestAccDataSourceSqlDatabase_elasticPool(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_sql_database", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: SqlDatabaseDataSource{}.elasticPool(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("elastic_pool_name").Exists(),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("server_name").Exists(),
			),
		},
	})
}

func TestAccDataSourceSqlDatabase_readScale(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_sql_database", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: SqlDatabaseDataSource{}.readScale(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("read_scale").HasValue("true"),
				check.That(data.ResourceName).Key("server_name").Exists(),
			),
		},
	})
}

func (d SqlDatabaseDataSource) basic(data acceptance.TestData) string {
	template := SqlDatabaseResource{}.basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_sql_database" "test" {
  name                = azurerm_sql_database.test.name
  server_name         = azurerm_sql_database.test.server_name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}

func (d SqlDatabaseDataSource) elasticPool(data acceptance.TestData) string {
	template := SqlDatabaseResource{}.elasticPool(data)
	return fmt.Sprintf(`
%s

data "azurerm_sql_database" "test" {
  name                = azurerm_sql_database.test.name
  server_name         = azurerm_sql_database.test.server_name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}

func (d SqlDatabaseDataSource) readScale(data acceptance.TestData, readScale bool) string {
	template := SqlDatabaseResource{}.readScale(data, readScale)
	return fmt.Sprintf(`
%s

data "azurerm_sql_database" "test" {
  name                = azurerm_sql_database.test.name
  server_name         = azurerm_sql_database.test.server_name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}

package sql_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

type SqlDatabaseDataSource struct {}

func TestAccDataSourceSqlDatabase_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_sql_database", "test")

	data.DataSourceTest(t, []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMSqlDatabase_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlDatabaseExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "collation"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "edition"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "location"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "read_scale"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "server_name"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
	})
}

func TestAccDataSourceSqlDatabase_elasticPool(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_sql_database", "test")

	data.DataSourceTest(t, []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMSqlDatabase_elasticPool(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlDatabaseExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "elastic_pool_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "location"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "server_name"),
				),
			},
	})
}

func TestAccDataSourceSqlDatabase_readScale(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_sql_database", "test")

	data.DataSourceTest(t, []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMSqlDatabase_readScale(data, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlDatabaseExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "location"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttr(data.ResourceName, "read_scale", "true"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "server_name"),
				),
			},
	})
}

func testAccDataSourceAzureRMSqlDatabase_basic(data acceptance.TestData) string {
	template := testAccAzureRMSqlDatabase_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_sql_database" "test" {
  name                = azurerm_sql_database.test.name
  server_name         = azurerm_sql_database.test.server_name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}

func testAccDataSourceAzureRMSqlDatabase_elasticPool(data acceptance.TestData) string {
	template := testAccAzureRMSqlDatabase_elasticPool(data)
	return fmt.Sprintf(`
%s

data "azurerm_sql_database" "test" {
  name                = azurerm_sql_database.test.name
  server_name         = azurerm_sql_database.test.server_name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}

func testAccDataSourceAzureRMSqlDatabase_readScale(data acceptance.TestData, readScale bool) string {
	template := testAccAzureRMSqlDatabase_readScale(data, readScale)
	return fmt.Sprintf(`
%s

data "azurerm_sql_database" "test" {
  name                = azurerm_sql_database.test.name
  server_name         = azurerm_sql_database.test.server_name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}

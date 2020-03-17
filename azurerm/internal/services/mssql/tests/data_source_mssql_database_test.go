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

func TestAccDataSourceAzureRMMsSqlDatabase_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_mssql_database", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMMsSqlDatabase_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "collation", "SQL_AltDiction_CP850_CI_AI"),
					resource.TestCheckResourceAttr(data.ResourceName, "license_type", "BasePrice"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "GP_Gen4_2"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.ENV", "Test"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMMsSqlDatabase_elasticPool(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_mssql_database", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMMsSqlDatabase_elasticPool(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "elastic_pool_id"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMMsSqlDatabase_BC(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_mssql_database", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMMsSqlDatabase_BC(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlDatabaseExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "read_scale", "Enabled"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "BC_Gen5_2"),
					resource.TestCheckResourceAttr(data.ResourceName, "zone_redundant", "true"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMMsSqlDatabase_basic(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlDatabase_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_mssql_database" "test" {
  name          = azurerm_mssql_database.test.name
  sql_server_id = azurerm_sql_server.test.id
}

`, template)
}

func testAccDataSourceAzureRMMsSqlDatabase_complete(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlDatabase_complete(data)
	return fmt.Sprintf(`
%s

data "azurerm_mssql_database" "test" {
  name          = azurerm_mssql_database.test.name
  sql_server_id = azurerm_sql_server.test.id
}

`, template)
}

func testAccDataSourceAzureRMMsSqlDatabase_elasticPool(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlDatabase_elasticPool(data)
	return fmt.Sprintf(`
%s

data "azurerm_mssql_database" "test" {
  name          = azurerm_mssql_database.test.name
  sql_server_id = azurerm_sql_server.test.id
}

`, template)
}

func testAccDataSourceAzureRMMsSqlDatabase_BC(data acceptance.TestData) string {
	template := testAccAzureRMMsSqlDatabase_BC(data)
	return fmt.Sprintf(`
%s

data "azurerm_mssql_database" "test" {
  name          = azurerm_mssql_database.test.name
  sql_server_id = azurerm_sql_server.test.id
}

`, template)
}

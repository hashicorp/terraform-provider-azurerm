package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMSqlDatabase_basic(t *testing.T) {
	dataSourceName := "data.azurerm_sql_database.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMSqlDatabase_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlDatabaseExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "collation"),
					resource.TestCheckResourceAttrSet(dataSourceName, "edition"),
					resource.TestCheckResourceAttrSet(dataSourceName, "location"),
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "read_scale"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "server_name"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMSqlDatabase_elasticPool(t *testing.T) {
	dataSourceName := "data.azurerm_sql_database.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMSqlDatabase_elasticPool(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlDatabaseExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "elastic_pool_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "location"),
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "server_name"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMSqlDatabase_readScale(t *testing.T) {
	dataSourceName := "data.azurerm_sql_database.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMSqlDatabase_readScale(ri, location, true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlDatabaseExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "location"),
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_group_name"),
					resource.TestCheckResourceAttr(dataSourceName, "read_scale", "true"),
					resource.TestCheckResourceAttrSet(dataSourceName, "server_name"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMSqlDatabase_basic(rInt int, location string) string {
	template := testAccAzureRMSqlDatabase_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_sql_database" "test" {
  name                = "${azurerm_sql_database.test.name}"
  server_name         = "${azurerm_sql_database.test.server_name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, template)
}

func testAccDataSourceAzureRMSqlDatabase_elasticPool(rInt int, location string) string {
	template := testAccAzureRMSqlDatabase_elasticPool(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_sql_database" "test" {
  name                = "${azurerm_sql_database.test.name}"
  server_name         = "${azurerm_sql_database.test.server_name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, template)
}

func testAccDataSourceAzureRMSqlDatabase_readScale(rInt int, location string, readScale bool) string {
	template := testAccAzureRMSqlDatabase_readScale(rInt, location, readScale)
	return fmt.Sprintf(`
%s

data "azurerm_sql_database" "test" {
  name                = "${azurerm_sql_database.test.name}"
  server_name         = "${azurerm_sql_database.test.server_name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, template)
}

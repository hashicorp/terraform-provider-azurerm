package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMMsSqlManagedDatabase_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_mssql_managed_database", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMsSqlManagedInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMMsSqlManagedDatabaseBasic_template(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMsSqlManagedInstanceExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttr(data.ResourceName, "location", data.Locations.Primary),
					resource.TestCheckResourceAttr(data.ResourceName, "collation", "Lithuanian_100_BIN"),
					resource.TestCheckResourceAttr(data.ResourceName, "catalog_collation", "SQL_Latin1_General_CP1_CI_AS"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMMsSqlManagedDatabaseBasic_template(data acceptance.TestData) string {
	template := testAccDataSourceAzureRMMsSqlManagedInstance_basic(data)
	return fmt.Sprintf(`%s

	resource "azurerm_mssql_managed_database" "test" {
		name                         = "acctest-db-%[2]d"
		managed_instance_id          = azurerm_mssql_managed_instance.test.id
		collation					 = "Lithuanian_100_BIN"
		catalog_collation           =  "SQL_Latin1_General_CP1_CI_AS"
	  }
	
	  data "azurerm_mssql_managed_database" "example" {
		name                         = azurerm_mssql_managed_database.test.name
		managed_instance_name = azurerm_mssql_managed_instance.test.name
		resource_group_name =  azurerm_mssql_managed_instance.test.resource_group_name
	   
	   }`, template, data.RandomInteger)
}

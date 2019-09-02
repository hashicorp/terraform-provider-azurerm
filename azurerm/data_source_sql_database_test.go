package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMSqlDatabase_basic(t *testing.T) {
	dataSourceName := "data.azurerm_sql_database.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMSqlDatabase_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlDatabaseExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "server_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_group_name"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "0"),
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

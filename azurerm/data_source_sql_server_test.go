package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMSqlServer_basic(t *testing.T) {
	dataSourceName := "data.azurerm_sql_server.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSqlServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMSqlServer_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSqlServerExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "location"),
					resource.TestCheckResourceAttrSet(dataSourceName, "fqdn"),
					resource.TestCheckResourceAttrSet(dataSourceName, "version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "administrator_login"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMSqlServer_basic(rInt int, location string) string {
	template := testAccAzureRMSqlServer_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_sql_server" "test" {
  name                = "${azurerm_sql_server.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, template)
}

package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMSqlServer(t *testing.T) {
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
				),
			},
		},
	})
}

//func TestAccDataSourceAzureRMServiceBusNamespace_premium(t *testing.T) {
//	dataSourceName := "data.azurerm_servicebus_namespace.test"
//	ri := tf.AccRandTimeInt()
//	location := testLocation()
//
//	resource.ParallelTest(t, resource.TestCase{
//		PreCheck:     func() { testAccPreCheck(t) },
//		Providers:    testAccProviders,
//		CheckDestroy: testCheckAzureRMSqlServerDestroy,
//		Steps: []resource.TestStep{
//			{
//				Config: testAccDataSourceAzureRMServiceBusNamespace_premium(ri, location),
//				Check: resource.ComposeTestCheckFunc(
//					testCheckAzureRMServiceBusNamespaceExists(dataSourceName),
//					resource.TestCheckResourceAttrSet(dataSourceName, "location"),
//					resource.TestCheckResourceAttrSet(dataSourceName, "sku"),
//					resource.TestCheckResourceAttrSet(dataSourceName, "capacity"),
//					resource.TestCheckResourceAttrSet(dataSourceName, "default_primary_connection_string"),
//					resource.TestCheckResourceAttrSet(dataSourceName, "default_secondary_connection_string"),
//					resource.TestCheckResourceAttrSet(dataSourceName, "default_primary_key"),
//					resource.TestCheckResourceAttrSet(dataSourceName, "default_secondary_key"),
//				),
//			},
//		},
//	})
//}

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

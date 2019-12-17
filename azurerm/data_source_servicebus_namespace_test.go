package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMServiceBusNamespace_basic(t *testing.T) {
	dataSourceName := "data.azurerm_servicebus_namespace.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMServiceBusNamespace_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusNamespaceExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "location"),
					resource.TestCheckResourceAttrSet(dataSourceName, "sku"),
					resource.TestCheckResourceAttrSet(dataSourceName, "capacity"),
					resource.TestCheckResourceAttrSet(dataSourceName, "default_primary_connection_string"),
					resource.TestCheckResourceAttrSet(dataSourceName, "default_secondary_connection_string"),
					resource.TestCheckResourceAttrSet(dataSourceName, "default_primary_key"),
					resource.TestCheckResourceAttrSet(dataSourceName, "default_secondary_key"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMServiceBusNamespace_premium(t *testing.T) {
	dataSourceName := "data.azurerm_servicebus_namespace.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMServiceBusNamespace_premium(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusNamespaceExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "location"),
					resource.TestCheckResourceAttrSet(dataSourceName, "sku"),
					resource.TestCheckResourceAttrSet(dataSourceName, "capacity"),
					resource.TestCheckResourceAttrSet(dataSourceName, "default_primary_connection_string"),
					resource.TestCheckResourceAttrSet(dataSourceName, "default_secondary_connection_string"),
					resource.TestCheckResourceAttrSet(dataSourceName, "default_primary_key"),
					resource.TestCheckResourceAttrSet(dataSourceName, "default_secondary_key"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMServiceBusNamespace_basic(rInt int, location string) string {
	template := testAccAzureRMServiceBusNamespace_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_servicebus_namespace" "test" {
  name                = "${azurerm_servicebus_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, template)
}

func testAccDataSourceAzureRMServiceBusNamespace_premium(rInt int, location string) string {
	template := testAccAzureRMServiceBusNamespace_premium(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_servicebus_namespace" "test" {
  name                = "${azurerm_servicebus_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, template)
}

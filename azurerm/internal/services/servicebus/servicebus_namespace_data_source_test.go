package servicebus_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMServiceBusNamespace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_servicebus_namespace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMServiceBusNamespace_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusNamespaceExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "location"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "sku"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "capacity"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "default_primary_connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "default_secondary_connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "default_primary_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "default_secondary_key"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMServiceBusNamespace_premium(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_servicebus_namespace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMServiceBusNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMServiceBusNamespace_premium(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServiceBusNamespaceExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "location"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "sku"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "capacity"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "default_primary_connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "default_secondary_connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "default_primary_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "default_secondary_key"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMServiceBusNamespace_basic(data acceptance.TestData) string {
	template := testAccAzureRMServiceBusNamespace_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_servicebus_namespace" "test" {
  name                = azurerm_servicebus_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}

func testAccDataSourceAzureRMServiceBusNamespace_premium(data acceptance.TestData) string {
	template := testAccAzureRMServiceBusNamespace_premium(data)
	return fmt.Sprintf(`
%s

data "azurerm_servicebus_namespace" "test" {
  name                = azurerm_servicebus_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}

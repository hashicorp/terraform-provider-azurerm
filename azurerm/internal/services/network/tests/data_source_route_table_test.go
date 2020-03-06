package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMRouteTable_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_route_table", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMRouteTable_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "route.#", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMRouteTable_singleRoute(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_route_table", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMRouteTable_singleRoute(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "route.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "route.0.name", "route1"),
					resource.TestCheckResourceAttr(data.ResourceName, "route.0.address_prefix", "10.1.0.0/16"),
					resource.TestCheckResourceAttr(data.ResourceName, "route.0.next_hop_type", "VnetLocal"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMRouteTable_multipleRoutes(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_route_table", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMRouteTable_multipleRoutes(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "route.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "route.0.name", "route1"),
					resource.TestCheckResourceAttr(data.ResourceName, "route.0.address_prefix", "10.1.0.0/16"),
					resource.TestCheckResourceAttr(data.ResourceName, "route.0.next_hop_type", "VnetLocal"),
					resource.TestCheckResourceAttr(data.ResourceName, "route.1.name", "route2"),
					resource.TestCheckResourceAttr(data.ResourceName, "route.1.address_prefix", "10.2.0.0/16"),
					resource.TestCheckResourceAttr(data.ResourceName, "route.1.next_hop_type", "VnetLocal"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMRouteTable_basic(data acceptance.TestData) string {
	r := testAccAzureRMRouteTable_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_route_table" "test" {
  name                = azurerm_route_table.test.name
  resource_group_name = azurerm_route_table.test.resource_group_name
}
`, r)
}

func testAccDataSourceAzureRMRouteTable_singleRoute(data acceptance.TestData) string {
	r := testAccAzureRMRouteTable_singleRoute(data)
	return fmt.Sprintf(`
%s

data "azurerm_route_table" "test" {
  name                = azurerm_route_table.test.name
  resource_group_name = azurerm_route_table.test.resource_group_name
}
`, r)
}

func testAccDataSourceAzureRMRouteTable_multipleRoutes(data acceptance.TestData) string {
	r := testAccAzureRMRouteTable_multipleRoutes(data)
	return fmt.Sprintf(`
%s

data "azurerm_route_table" "test" {
  name                = azurerm_route_table.test.name
  resource_group_name = azurerm_route_table.test.resource_group_name
}
`, r)
}

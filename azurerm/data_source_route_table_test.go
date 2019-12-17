package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMRouteTable_basic(t *testing.T) {
	dataSourceName := "data.azurerm_route_table.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	config := testAccDataSourceAzureRMRouteTable_basic(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "route.#", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMRouteTable_singleRoute(t *testing.T) {
	dataSourceName := "data.azurerm_route_table.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	config := testAccDataSourceAzureRMRouteTable_singleRoute(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "route.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "route.0.name", "route1"),
					resource.TestCheckResourceAttr(dataSourceName, "route.0.address_prefix", "10.1.0.0/16"),
					resource.TestCheckResourceAttr(dataSourceName, "route.0.next_hop_type", "VnetLocal"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMRouteTable_multipleRoutes(t *testing.T) {
	dataSourceName := "data.azurerm_route_table.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	config := testAccDataSourceAzureRMRouteTable_multipleRoutes(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRouteTableExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "route.#", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "route.0.name", "route1"),
					resource.TestCheckResourceAttr(dataSourceName, "route.0.address_prefix", "10.1.0.0/16"),
					resource.TestCheckResourceAttr(dataSourceName, "route.0.next_hop_type", "VnetLocal"),
					resource.TestCheckResourceAttr(dataSourceName, "route.1.name", "route2"),
					resource.TestCheckResourceAttr(dataSourceName, "route.1.address_prefix", "10.2.0.0/16"),
					resource.TestCheckResourceAttr(dataSourceName, "route.1.next_hop_type", "VnetLocal"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMRouteTable_basic(rInt int, location string) string {
	r := testAccAzureRMRouteTable_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_route_table" "test" {
  name                = "${azurerm_route_table.test.name}"
  resource_group_name = "${azurerm_route_table.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMRouteTable_singleRoute(rInt int, location string) string {
	r := testAccAzureRMRouteTable_singleRoute(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_route_table" "test" {
  name                = "${azurerm_route_table.test.name}"
  resource_group_name = "${azurerm_route_table.test.resource_group_name}"
}
`, r)
}

func testAccDataSourceAzureRMRouteTable_multipleRoutes(rInt int, location string) string {
	r := testAccAzureRMRouteTable_multipleRoutes(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_route_table" "test" {
  name                = "${azurerm_route_table.test.name}"
  resource_group_name = "${azurerm_route_table.test.resource_group_name}"
}
`, r)
}

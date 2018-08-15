package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMRouteTable_basic(t *testing.T) {
	dataSourceName := "data.azurerm_route_table.test"
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccDataSourceAzureRMRouteTable_basic(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
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
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccDataSourceAzureRMRouteTable_singleRoute(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
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
	ri := acctest.RandInt()
	location := testLocation()
	config := testAccDataSourceAzureRMRouteTable_multipleRoutes(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
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
	resource := testAccAzureRMRouteTable_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_route_table" "test" {
  name                = "${azurerm_route_table.test.name}"
  resource_group_name = "${azurerm_route_table.test.resource_group_name}"
}
`, resource)
}

func testAccDataSourceAzureRMRouteTable_singleRoute(rInt int, location string) string {
	resource := testAccAzureRMRouteTable_singleRoute(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_route_table" "test" {
  name                = "${azurerm_route_table.test.name}"
  resource_group_name = "${azurerm_route_table.test.resource_group_name}"
}
`, resource)
}

func testAccDataSourceAzureRMRouteTable_multipleRoutes(rInt int, location string) string {
	resource := testAccAzureRMRouteTable_multipleRoutes(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_route_table" "test" {
  name                = "${azurerm_route_table.test.name}"
  resource_group_name = "${azurerm_route_table.test.resource_group_name}"
}
`, resource)
}

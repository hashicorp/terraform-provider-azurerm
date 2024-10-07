// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type RouteTableDataSource struct{}

func TestAccDataSourceRouteTable_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_route_table", "test")
	r := RouteTableDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("bgp_route_propagation_enabled").Exists(),
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("route.#").Exists(),
				check.That(data.ResourceName).Key("route.#").HasValue("0"),
				check.That(data.ResourceName).Key("subnets.#").Exists(),
				check.That(data.ResourceName).Key("tags.%").Exists(),
			),
		},
	})
}

func TestAccDataSourceRouteTable_singleRoute(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_route_table", "test")
	r := RouteTableDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.singleRoute(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("route.#").HasValue("1"),
				check.That(data.ResourceName).Key("route.0.name").HasValue("route1"),
				check.That(data.ResourceName).Key("route.0.address_prefix").HasValue("10.1.0.0/16"),
				check.That(data.ResourceName).Key("route.0.next_hop_type").HasValue("VnetLocal"),
			),
		},
	})
}

func TestAccDataSourceRouteTable_multipleRoutes(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_route_table", "test")
	r := RouteTableDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.multipleRoutes(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("route.#").HasValue("2"),
			),
		},
	})
}

func (RouteTableDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_route_table" "test" {
  name                = azurerm_route_table.test.name
  resource_group_name = azurerm_route_table.test.resource_group_name
}
`, RouteTableResource{}.basic(data))
}

func (RouteTableDataSource) singleRoute(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_route_table" "test" {
  name                = azurerm_route_table.test.name
  resource_group_name = azurerm_route_table.test.resource_group_name
}
`, RouteTableResource{}.singleRoute(data))
}

func (RouteTableDataSource) multipleRoutes(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_route_table" "test" {
  name                = azurerm_route_table.test.name
  resource_group_name = azurerm_route_table.test.resource_group_name
}
`, RouteTableResource{}.multipleRoutes(data))
}

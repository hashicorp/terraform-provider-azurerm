// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/routetables"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type RouteTableResource struct{}

func TestAccRouteTable_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_table", "test")
	r := RouteTableResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("disable_bgp_route_propagation").HasValue("false"),
				check.That(data.ResourceName).Key("route.#").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRouteTable_basicNilNextHopIPAddress(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_table", "test")
	r := RouteTableResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.nilNextHopIPAddess(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("disable_bgp_route_propagation").HasValue("false"),
				check.That(data.ResourceName).Key("route.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRouteTable_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_table", "test")
	r := RouteTableResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("disable_bgp_route_propagation").HasValue("false"),
				check.That(data.ResourceName).Key("route.#").HasValue("0"),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_route_table"),
		},
	})
}

func TestAccRouteTable_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_table", "test")
	r := RouteTableResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("disable_bgp_route_propagation").HasValue("true"),
				check.That(data.ResourceName).Key("route.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRouteTable_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_table", "test")
	r := RouteTableResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("disable_bgp_route_propagation").HasValue("false"),
				check.That(data.ResourceName).Key("route.#").HasValue("0"),
			),
		},
		{
			Config: r.basicAppliance(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("disable_bgp_route_propagation").HasValue("false"),
				check.That(data.ResourceName).Key("route.#").HasValue("1"),
			),
		},
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("disable_bgp_route_propagation").HasValue("true"),
				check.That(data.ResourceName).Key("route.#").HasValue("1"),
			),
		},
	})
}

func TestAccRouteTable_singleRoute(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_table", "test")
	r := RouteTableResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.singleRoute(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccRouteTable_removeRoute(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_table", "test")
	r := RouteTableResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// This configuration includes a single explicit route block
			Config: r.singleRoute(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("route.#").HasValue("1"),
			),
		},
		{
			// This configuration has no route blocks at all.
			Config: r.noRouteBlocks(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				// The route from the first step is preserved because no
				// blocks at all means "ignore existing blocks".
				check.That(data.ResourceName).Key("route.#").HasValue("1"),
			),
		},
		{
			// This configuration sets route to [] explicitly using the
			// attribute syntax.
			Config: r.singleRouteRemoved(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				// The route from the first step is now removed, leaving us
				// with no routes at all.
				check.That(data.ResourceName).Key("route.#").HasValue("0"),
			),
		},
	})
}

func TestAccRouteTable_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_table", "test")
	r := RouteTableResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func TestAccRouteTable_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_table", "test")
	r := RouteTableResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("Production"),
				check.That(data.ResourceName).Key("tags.cost_center").HasValue("MSFT"),
			),
		},
		{
			Config: r.withTagsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("staging"),
			),
		},
	})
}

func TestAccRouteTable_multipleRoutes(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_table", "test")
	r := RouteTableResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.singleRoute(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("route.#").HasValue("1"),
			),
		},
		{
			Config: r.multipleRoutes(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("route.#").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func (t RouteTableResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := routetables.ParseRouteTableID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.RouteTables.Get(ctx, *id, routetables.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("reading Route Table (%s): %+v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (RouteTableResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := routetables.ParseRouteTableID(state.ID)
	if err != nil {
		return nil, err
	}

	if err = client.Network.RouteTables.DeleteThenPoll(ctx, *id); err != nil {
		return nil, fmt.Errorf("deleting Route Table %q: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (RouteTableResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_route_table" "test" {
  name                = "acctestrt%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RouteTableResource) nilNextHopIPAddess(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_route_table" "test" {
  name                = "acctestrt%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  route {
    name                   = "route1"
    address_prefix         = "101.1.0.0/16"
    next_hop_type          = "Internet"
    next_hop_in_ip_address = null
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r RouteTableResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_route_table" "import" {
  name                = azurerm_route_table.test.name
  location            = azurerm_route_table.test.location
  resource_group_name = azurerm_route_table.test.resource_group_name
}
`, r.basic(data))
}

func (RouteTableResource) basicAppliance(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_route_table" "test" {
  name                = "acctestrt%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  route {
    name                   = "route1"
    address_prefix         = "10.1.0.0/16"
    next_hop_type          = "VirtualAppliance"
    next_hop_in_ip_address = "192.168.0.1"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RouteTableResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_route_table" "test" {
  name                = "acctestrt%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  route {
    name           = "acctestRoute"
    address_prefix = "10.1.0.0/16"
    next_hop_type  = "VnetLocal"
  }

  disable_bgp_route_propagation = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RouteTableResource) singleRoute(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_route_table" "test" {
  name                = "acctestrt%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  route {
    name           = "route1"
    address_prefix = "10.1.0.0/16"
    next_hop_type  = "VnetLocal"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RouteTableResource) noRouteBlocks(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_route_table" "test" {
  name                = "acctestrt%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RouteTableResource) singleRouteRemoved(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_route_table" "test" {
  name                = "acctestrt%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  route = []
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RouteTableResource) multipleRoutes(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_route_table" "test" {
  name                = "acctestrt%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  route {
    name           = "route1"
    address_prefix = "10.1.0.0/16"
    next_hop_type  = "VnetLocal"
  }

  route {
    name           = "route2"
    address_prefix = "10.2.0.0/16"
    next_hop_type  = "VnetLocal"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RouteTableResource) withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_route_table" "test" {
  name                = "acctestrt%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  route {
    name           = "route1"
    address_prefix = "10.1.0.0/16"
    next_hop_type  = "VnetLocal"
  }

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (RouteTableResource) withTagsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_route_table" "test" {
  name                = "acctestrt%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  route {
    name           = "route1"
    address_prefix = "10.1.0.0/16"
    next_hop_type  = "VnetLocal"
  }

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

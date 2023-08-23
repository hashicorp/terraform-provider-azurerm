// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type RouteServerResource struct{}

func TestAccRouteServer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_server", "test")
	r := RouteServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func TestAccRouteServer_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_server", "test")
	r := RouteServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_route_server"),
		},
	})
}

func TestAccRouteServer_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_server", "test")
	r := RouteServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func TestAccRouteServer_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_route_server", "test")
	r := RouteServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}
func (r RouteServerResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	routeServerId, err := parse.VirtualHubID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.VirtualHubClient.Get(ctx, routeServerId.ResourceGroup, routeServerId.Name)
	if err != nil {
		return nil, fmt.Errorf("reading Route Server %s: %+v", routeServerId, err)
	}

	ipConfig, err := clients.Network.VirtualHubIPClient.List(ctx, routeServerId.ResourceGroup, routeServerId.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Ip Config for Route Server %s: %+v", routeServerId, err)
	}
	if ipConfig.Values() == nil {
		return nil, fmt.Errorf("no IP Config is set for the Route Server %s: %+v", routeServerId, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r RouteServerResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_route_server" "test" {
  name                 = "acctestrs-%d"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  sku                  = "Standard"
  public_ip_address_id = azurerm_public_ip.test.id
  subnet_id            = azurerm_subnet.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r RouteServerResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_route_server" "import" {
  name                 = azurerm_route_server.test.name
  resource_group_name  = azurerm_route_server.test.resource_group_name
  location             = azurerm_route_server.test.location
  sku                  = azurerm_route_server.test.sku
  public_ip_address_id = azurerm_route_server.test.public_ip_address_id
  subnet_id            = azurerm_route_server.test.subnet_id
}
`, r.basic(data))
}

func (r RouteServerResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_route_server" "test" {
  name                             = "acctestrs-%d"
  resource_group_name              = azurerm_resource_group.test.name
  location                         = azurerm_resource_group.test.location
  sku                              = "Standard"
  public_ip_address_id             = azurerm_public_ip.test.id
  subnet_id                        = azurerm_subnet.test.id
  branch_to_branch_traffic_enabled = true
}
`, r.template(data), data.RandomInteger)
}

func (r RouteServerResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "RouteServerSubnet"
  virtual_network_name = azurerm_virtual_network.test.name
  resource_group_name  = azurerm_resource_group.test.name
  address_prefixes     = ["10.0.0.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  allocation_method   = "Static"
  sku                 = "Standard"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

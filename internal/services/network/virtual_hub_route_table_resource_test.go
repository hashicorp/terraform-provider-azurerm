// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type VirtualHubRouteTableResource struct{}

func TestAccVirtualHubRouteTable_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_route_table", "test")
	r := VirtualHubRouteTableResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualHubRouteTable_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_route_table", "test")
	r := VirtualHubRouteTableResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccVirtualHubRouteTable_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_route_table", "test")
	r := VirtualHubRouteTableResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualHubRouteTable_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_route_table", "test")
	r := VirtualHubRouteTableResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t VirtualHubRouteTableResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := virtualwans.ParseHubRouteTableID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.VirtualWANs.HubRouteTablesGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (VirtualHubRouteTableResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-VHUB-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VNET-%d"
  address_space       = ["10.5.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_network_security_group" "test" {
  name                = "acctest-NSG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-SUBNET-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.5.1.0/24"]
}

resource "azurerm_subnet_network_security_group_association" "test" {
  subnet_id                 = azurerm_subnet.test.id
  network_security_group_id = azurerm_network_security_group.test.id
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctest-VWAN-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_hub" "test" {
  name                = "acctest-VHUB-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_wan_id      = azurerm_virtual_wan.test.id
  address_prefix      = "10.0.2.0/24"
}

resource "azurerm_virtual_hub_connection" "test" {
  name                      = "acctest-VHUBCONN-%d"
  virtual_hub_id            = azurerm_virtual_hub.test.id
  remote_virtual_network_id = azurerm_virtual_network.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r VirtualHubRouteTableResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_route_table" "test" {
  name           = "acctest-RouteTable-%d"
  virtual_hub_id = azurerm_virtual_hub.test.id
  labels         = ["Label1"]
}
`, r.template(data), data.RandomInteger)
}

func (r VirtualHubRouteTableResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_route_table" "import" {
  name           = azurerm_virtual_hub_route_table.test.name
  virtual_hub_id = azurerm_virtual_hub_route_table.test.virtual_hub_id
  labels         = azurerm_virtual_hub_route_table.test.labels
}
`, r.basic(data))
}

func (r VirtualHubRouteTableResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_route_table" "test" {
  name           = "acctest-RouteTable-%d"
  virtual_hub_id = azurerm_virtual_hub.test.id
  labels         = ["labeL1", "AnotherLabel"]

  route {
    name              = "VHub-Route-Test"
    destinations_type = "CIDR"
    destinations      = ["10.0.0.0/16"]
    next_hop_type     = "ResourceId"
    next_hop          = azurerm_virtual_hub_connection.test.id
  }
}
`, r.template(data), data.RandomInteger)
}

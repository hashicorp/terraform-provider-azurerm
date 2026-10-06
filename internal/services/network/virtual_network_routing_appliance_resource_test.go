// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworkappliances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type VirtualNetworkRoutingApplianceResource struct{}

func TestAccVirtualNetworkRoutingAppliance_sequential(t *testing.T) {
	// These tests run sequentially because Azure limits routing appliances to 2 per region per subscription.
	// https://learn.microsoft.com/azure/virtual-network/virtual-network-routing-appliance-overview#limitations
	testCases := map[string]map[string]func(t *testing.T){
		"RoutingAppliance": {
			"basic":          testAccVirtualNetworkRoutingAppliance_basic,
			"complete":       testAccVirtualNetworkRoutingAppliance_complete,
			"update":         testAccVirtualNetworkRoutingAppliance_update,
			"requiresImport": testAccVirtualNetworkRoutingAppliance_requiresImport,
			"dualStack":      testAccVirtualNetworkRoutingAppliance_dualStack,
			"identity":       testAccVirtualNetworkRoutingAppliance_resourceIdentity,
			"list":           testAccVirtualNetworkRoutingAppliance_list_basic,
		},
	}

	for group, tests := range testCases {
		t.Run(group, func(t *testing.T) {
			for name, test := range tests {
				t.Run(name, test)
			}
		})
	}
}

func testAccVirtualNetworkRoutingAppliance_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_routing_appliance", "test")
	r := VirtualNetworkRoutingApplianceResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_guid").IsNotEmpty(),
				check.That(data.ResourceName).Key("ip_configuration.0.name").IsNotEmpty(),
				check.That(data.ResourceName).Key("ip_configuration.0.primary").Exists(),
				check.That(data.ResourceName).Key("ip_configuration.0.private_ip_address").IsNotEmpty(),
				check.That(data.ResourceName).Key("ip_configuration.0.private_ip_address_allocation").IsNotEmpty(),
			),
		},
		data.ImportStep(),
	})
}

func testAccVirtualNetworkRoutingAppliance_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_routing_appliance", "test")
	r := VirtualNetworkRoutingApplianceResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep(),
	})
}

func testAccVirtualNetworkRoutingAppliance_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_routing_appliance", "test")
	r := VirtualNetworkRoutingApplianceResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep(),
	})
}

func testAccVirtualNetworkRoutingAppliance_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_routing_appliance", "test")
	r := VirtualNetworkRoutingApplianceResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func testAccVirtualNetworkRoutingAppliance_dualStack(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_routing_appliance", "test")
	r := VirtualNetworkRoutingApplianceResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.dualStack(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep(),
	})
}

func (VirtualNetworkRoutingApplianceResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := virtualnetworkappliances.ParseVirtualNetworkApplianceID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Network.VirtualNetworkAppliances.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r VirtualNetworkRoutingApplianceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network_routing_appliance" "test" {
  name                = "acctestVNRA-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  bandwidth_in_gbps   = 10
  subnet_id           = azurerm_subnet.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r VirtualNetworkRoutingApplianceResource) dualStack(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network_routing_appliance" "test" {
  name                       = "acctestVNRA-%d"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  bandwidth_in_gbps          = 10
  subnet_id                  = azurerm_subnet.test.id
  private_ip_address_version = "DualStack"
}
`, r.template(data), data.RandomInteger)
}

func (r VirtualNetworkRoutingApplianceResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network_routing_appliance" "test" {
  name                       = "acctestVNRA-%d"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  bandwidth_in_gbps          = 10
  subnet_id                  = azurerm_subnet.test.id
  private_ip_address_version = "IPv4"

  tags = {
    environment = "test"
    purpose     = "routing"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r VirtualNetworkRoutingApplianceResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network_routing_appliance" "import" {
  name                = azurerm_virtual_network_routing_appliance.test.name
  resource_group_name = azurerm_virtual_network_routing_appliance.test.resource_group_name
  location            = azurerm_virtual_network_routing_appliance.test.location
  bandwidth_in_gbps   = azurerm_virtual_network_routing_appliance.test.bandwidth_in_gbps
  subnet_id           = azurerm_virtual_network_routing_appliance.test.subnet_id
}
`, r.basic(data))
}

func (VirtualNetworkRoutingApplianceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestVNet-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.0.0.0/16", "fd00:db8::/48"]
}

resource "azurerm_subnet" "test" {
  name                 = "VirtualNetworkApplianceSubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.0.0/24", "fd00:db8:0:1::/64"]
}
`, data.RandomInteger, data.Locations.Primary)
}

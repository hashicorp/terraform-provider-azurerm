// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualnetworkgateways"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type VirtualNetworkGatewayNatRuleResource struct{}

func TestAccVirtualNetworkGatewayNatRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway_nat_rule", "test")
	r := VirtualNetworkGatewayNatRuleResource{}

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

func TestAccVirtualNetworkGatewayNatRule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway_nat_rule", "test")
	r := VirtualNetworkGatewayNatRuleResource{}

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

func TestAccVirtualNetworkGatewayNatRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway_nat_rule", "test")
	r := VirtualNetworkGatewayNatRuleResource{}

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

func TestAccVirtualNetworkGatewayNatRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway_nat_rule", "test")
	r := VirtualNetworkGatewayNatRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualNetworkGatewayNatRule_updatePortRange(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway_nat_rule", "test")
	r := VirtualNetworkGatewayNatRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.updatePortRange(data, "10.1.0.0/26", "100", "10.2.0.0/26", "200"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updatePortRange(data, "10.3.0.0/26", "300", "10.4.0.0/26", "400"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r VirtualNetworkGatewayNatRuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := virtualnetworkgateways.ParseVirtualNetworkGatewayNatRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.VirtualNetworkGateways.VirtualNetworkGatewayNatRulesGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r VirtualNetworkGatewayNatRuleResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network_gateway_nat_rule" "test" {
  name                       = "acctest-vnetgwnatrule-%d"
  resource_group_name        = azurerm_resource_group.test.name
  virtual_network_gateway_id = data.azurerm_virtual_network_gateway.test.id

  external_mapping {
    address_space = "10.1.0.0/26"
  }

  internal_mapping {
    address_space = "10.3.0.0/26"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r VirtualNetworkGatewayNatRuleResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network_gateway_nat_rule" "test" {
  name                       = "acctest-vnetgwnatrule-%d"
  resource_group_name        = azurerm_resource_group.test.name
  virtual_network_gateway_id = data.azurerm_virtual_network_gateway.test.id
  mode                       = "EgressSnat"
  type                       = "Dynamic"
  ip_configuration_id        = data.azurerm_virtual_network_gateway.test.ip_configuration.0.id

  external_mapping {
    address_space = "10.1.0.0/26"
  }

  external_mapping {
    address_space = "10.2.0.0/26"
  }

  internal_mapping {
    address_space = "10.3.0.0/26"
  }

  internal_mapping {
    address_space = "10.4.0.0/26"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r VirtualNetworkGatewayNatRuleResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network_gateway_nat_rule" "test" {
  name                       = "acctest-vnetgwnatrule-%d"
  resource_group_name        = azurerm_resource_group.test.name
  virtual_network_gateway_id = data.azurerm_virtual_network_gateway.test.id
  mode                       = "EgressSnat"
  type                       = "Dynamic"
  ip_configuration_id        = data.azurerm_virtual_network_gateway.test.ip_configuration.0.id

  external_mapping {
    address_space = "10.2.0.0/26"
  }

  internal_mapping {
    address_space = "10.4.0.0/26"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r VirtualNetworkGatewayNatRuleResource) updatePortRange(data acceptance.TestData, externalAddressSpace, externalPortRange, internalAddressSpace, internalPortRange string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network_gateway_nat_rule" "test" {
  name                       = "acctest-vnetgwnatrule-%d"
  resource_group_name        = azurerm_resource_group.test.name
  virtual_network_gateway_id = data.azurerm_virtual_network_gateway.test.id
  mode                       = "EgressSnat"
  type                       = "Static"

  external_mapping {
    address_space = "%s"
    port_range    = "%s"
  }

  internal_mapping {
    address_space = "%s"
    port_range    = "%s"
  }
}
`, r.template(data), data.RandomInteger, externalAddressSpace, externalPortRange, internalAddressSpace, internalPortRange)
}

func (r VirtualNetworkGatewayNatRuleResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network_gateway_nat_rule" "import" {
  name                       = azurerm_virtual_network_gateway_nat_rule.test.name
  resource_group_name        = azurerm_virtual_network_gateway_nat_rule.test.resource_group_name
  virtual_network_gateway_id = azurerm_virtual_network_gateway_nat_rule.test.virtual_network_gateway_id

  external_mapping {
    address_space = "10.1.0.0/26"
  }

  internal_mapping {
    address_space = "10.3.0.0/26"
  }
}
`, r.basic(data))
}

func (VirtualNetworkGatewayNatRuleResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-vnetgwnatrule-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
  sku                 = "Basic"
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctestvng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "Basic"

  ip_configuration {
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }
}

data "azurerm_virtual_network_gateway" "test" {
  name                = azurerm_virtual_network_gateway.test.name
  resource_group_name = azurerm_virtual_network_gateway.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

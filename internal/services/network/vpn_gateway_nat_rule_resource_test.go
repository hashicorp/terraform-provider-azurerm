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
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type VPNGatewayNatRuleResource struct{}

func TestAccVpnGatewayNatRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_gateway_nat_rule", "test")
	r := VPNGatewayNatRuleResource{}

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

func TestAccVpnGatewayNatRule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_gateway_nat_rule", "test")
	r := VPNGatewayNatRuleResource{}

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

func TestAccVpnGatewayNatRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_gateway_nat_rule", "test")
	r := VPNGatewayNatRuleResource{}

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

func TestAccVpnGatewayNatRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_gateway_nat_rule", "test")
	r := VPNGatewayNatRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
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

func TestAccVpnGatewayNatRule_externalMappingAndInternalMapping(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_gateway_nat_rule", "test")
	r := VPNGatewayNatRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.externalMappingAndInternalMapping(data, "10.2.0.0/26", "200", "10.4.0.0/26", "400"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.externalMappingAndInternalMapping(data, "10.3.0.0/26", "300", "10.5.0.0/26", "500"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r VPNGatewayNatRuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.VpnGatewayNatRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.NatRuleClient.Get(ctx, id.ResourceGroup, id.VpnGatewayName, id.NatRuleName)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r VPNGatewayNatRuleResource) basic(data acceptance.TestData) string {
	if !features.FourPointOhBeta() {
		return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway_nat_rule" "test" {
  name                            = "acctest-vpnnatrule-%d"
  resource_group_name             = azurerm_resource_group.test.name
  vpn_gateway_id                  = azurerm_vpn_gateway.test.id
  external_address_space_mappings = ["192.168.21.0/26"]
  internal_address_space_mappings = ["10.4.0.0/26"]
}
`, r.template(data), data.RandomInteger)
	}

	return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway_nat_rule" "test" {
  name                = "acctest-vpnnatrule-%d"
  resource_group_name = azurerm_resource_group.test.name
  vpn_gateway_id      = azurerm_vpn_gateway.test.id

  external_mapping {
    address_space = "192.168.21.0/26"
  }

  internal_mapping {
    address_space = "10.4.0.0/26"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r VPNGatewayNatRuleResource) complete(data acceptance.TestData) string {
	if !features.FourPointOhBeta() {
		return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway_nat_rule" "test" {
  name                            = "acctest-vpnnatrule-%d"
  resource_group_name             = azurerm_resource_group.test.name
  vpn_gateway_id                  = azurerm_vpn_gateway.test.id
  external_address_space_mappings = ["192.168.21.0/26"]
  internal_address_space_mappings = ["10.4.0.0/26"]
  mode                            = "EgressSnat"
  type                            = "Dynamic"
  ip_configuration_id             = "Instance0"
}
`, r.template(data), data.RandomInteger)
	}

	return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway_nat_rule" "test" {
  name                = "acctest-vpnnatrule-%d"
  resource_group_name = azurerm_resource_group.test.name
  vpn_gateway_id      = azurerm_vpn_gateway.test.id
  mode                = "EgressSnat"
  type                = "Dynamic"
  ip_configuration_id = "Instance0"

  external_mapping {
    address_space = "192.168.21.0/26"
  }

  internal_mapping {
    address_space = "10.4.0.0/26"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r VPNGatewayNatRuleResource) update(data acceptance.TestData) string {
	if !features.FourPointOhBeta() {
		return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway_nat_rule" "test" {
  name                            = "acctest-vpnnatrule-%d"
  resource_group_name             = azurerm_resource_group.test.name
  vpn_gateway_id                  = azurerm_vpn_gateway.test.id
  external_address_space_mappings = ["192.168.22.0/26"]
  internal_address_space_mappings = ["10.5.0.0/26"]
  mode                            = "EgressSnat"
  type                            = "Dynamic"
  ip_configuration_id             = "Instance1"
}
`, r.template(data), data.RandomInteger)
	}

	return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway_nat_rule" "test" {
  name                = "acctest-vpnnatrule-%d"
  resource_group_name = azurerm_resource_group.test.name
  vpn_gateway_id      = azurerm_vpn_gateway.test.id
  mode                = "EgressSnat"
  type                = "Dynamic"
  ip_configuration_id = "Instance1"

  external_mapping {
    address_space = "192.168.22.0/26"
  }

  internal_mapping {
    address_space = "10.5.0.0/26"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r VPNGatewayNatRuleResource) requiresImport(data acceptance.TestData) string {
	if !features.FourPointOhBeta() {
		return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway_nat_rule" "import" {
  name                            = azurerm_vpn_gateway_nat_rule.test.name
  resource_group_name             = azurerm_vpn_gateway_nat_rule.test.resource_group_name
  vpn_gateway_id                  = azurerm_vpn_gateway_nat_rule.test.vpn_gateway_id
  external_address_space_mappings = azurerm_vpn_gateway_nat_rule.test.external_address_space_mappings
  internal_address_space_mappings = azurerm_vpn_gateway_nat_rule.test.internal_address_space_mappings
  mode                            = azurerm_vpn_gateway_nat_rule.test.mode
  type                            = azurerm_vpn_gateway_nat_rule.test.type
}
`, r.basic(data))
	}

	return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway_nat_rule" "import" {
  name                = azurerm_vpn_gateway_nat_rule.test.name
  resource_group_name = azurerm_vpn_gateway_nat_rule.test.resource_group_name
  vpn_gateway_id      = azurerm_vpn_gateway_nat_rule.test.vpn_gateway_id
  mode                = azurerm_vpn_gateway_nat_rule.test.mode
  type                = azurerm_vpn_gateway_nat_rule.test.type

  external_mapping {
    address_space = "192.168.21.0/26"
  }

  internal_mapping {
    address_space = "10.4.0.0/26"
  }
}
`, r.basic(data))
}

func (r VPNGatewayNatRuleResource) externalMappingAndInternalMapping(data acceptance.TestData, externalAddressSpace, externalPortRange, internalAddressSpace, internalPortRange string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway_nat_rule" "test" {
  name                = "acctest-vpnnatrule-%d"
  resource_group_name = azurerm_resource_group.test.name
  vpn_gateway_id      = azurerm_vpn_gateway.test.id

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

func (VPNGatewayNatRuleResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-vpnnatrule-%d"
  location = "%s"
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctest-vwan-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_hub" "test" {
  name                = "acctest-vhub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_prefix      = "10.0.2.0/24"
  virtual_wan_id      = azurerm_virtual_wan.test.id
}

resource "azurerm_vpn_gateway" "test" {
  name                = "acctest-vpngateway-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  virtual_hub_id      = azurerm_virtual_hub.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

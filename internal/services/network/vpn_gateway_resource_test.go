// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type VPNGatewayResource struct{}

func TestAccVPNGateway_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_gateway", "test")
	r := VPNGatewayResource{}

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

func TestAccVPNGateway_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_gateway", "test")
	r := VPNGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_vpn_gateway"),
		},
	})
}

func TestAccVPNGateway_bgpSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_gateway", "test")
	r := VPNGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.bgpSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVPNGateway_scaleUnit(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_gateway", "test")
	r := VPNGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.scaleUnit(data, 2),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.scaleUnit(data, 3),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVPNGateway_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_gateway", "test")
	r := VPNGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.tags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.tagsUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.tags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVPNGateway_routingPreferenceMicrosoftNetwork(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_gateway", "test")
	r := VPNGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.routingPreference(data, "Microsoft Network"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVPNGateway_routingPreferenceInternet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_gateway", "test")
	r := VPNGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.routingPreference(data, "Internet"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVPNGateway_bgpRouteTranslationForNatEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_gateway", "test")
	r := VPNGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.bgpRouteTranslationForNatEnabled(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.bgpRouteTranslationForNatEnabled(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t VPNGatewayResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := virtualwans.ParseVpnGatewayID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.VirtualWANs.VpnGatewaysGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r VPNGatewayResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway" "test" {
  name                = "acctestVPNG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  virtual_hub_id      = azurerm_virtual_hub.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r VPNGatewayResource) routingPreference(data acceptance.TestData, routingPreference string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway" "test" {
  name                = "acctestVPNG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  virtual_hub_id      = azurerm_virtual_hub.test.id
  routing_preference  = "%s"
}
`, r.template(data), data.RandomInteger, routingPreference)
}

func (r VPNGatewayResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway" "import" {
  name                = azurerm_vpn_gateway.test.name
  location            = azurerm_vpn_gateway.test.location
  resource_group_name = azurerm_vpn_gateway.test.resource_group_name
  virtual_hub_id      = azurerm_vpn_gateway.test.virtual_hub_id
}
`, r.basic(data))
}

func (r VPNGatewayResource) bgpSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway" "test" {
  name                = "acctestVPNG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  virtual_hub_id      = azurerm_virtual_hub.test.id

  bgp_settings {
    asn         = 65515
    peer_weight = 0

    instance_0_bgp_peering_address {
      custom_ips = ["169.254.21.5"]
    }

    instance_1_bgp_peering_address {
      custom_ips = ["169.254.21.10"]
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r VPNGatewayResource) scaleUnit(data acceptance.TestData, scaleUnit int) string {
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway" "test" {
  name                = "acctestVPNG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  virtual_hub_id      = azurerm_virtual_hub.test.id
  scale_unit          = %d
}
`, r.template(data), data.RandomInteger, scaleUnit)
}

func (r VPNGatewayResource) tags(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway" "test" {
  name                = "acctestVPNG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  virtual_hub_id      = azurerm_virtual_hub.test.id

  tags = {
    Hello = "World"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r VPNGatewayResource) tagsUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway" "test" {
  name                = "acctestVPNG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  virtual_hub_id      = azurerm_virtual_hub.test.id

  tags = {
    Hello = "World"
    Rick  = "C-137"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r VPNGatewayResource) bgpRouteTranslationForNatEnabled(data acceptance.TestData, bgpRouteTranslationForNatEnabled bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_gateway" "test" {
  name                                  = "acctestVPNG-%d"
  location                              = azurerm_resource_group.test.location
  resource_group_name                   = azurerm_resource_group.test.name
  virtual_hub_id                        = azurerm_virtual_hub.test.id
  bgp_route_translation_for_nat_enabled = %t
}
`, r.template(data), data.RandomInteger, bgpRouteTranslationForNatEnabled)
}

func (VPNGatewayResource) template(data acceptance.TestData) string {
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
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctestvwan-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_hub" "test" {
  name                = "acctestvh-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_prefix      = "10.0.1.0/24"
  virtual_wan_id      = azurerm_virtual_wan.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

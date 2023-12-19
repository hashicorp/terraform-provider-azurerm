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
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type VirtualNetworkGatewayResource struct{}

func TestAccVirtualNetworkGateway_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway", "test")
	r := VirtualNetworkGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku").HasValue("Basic"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualNetworkGateway_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway", "test")
	r := VirtualNetworkGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_virtual_network_gateway"),
		},
	})
}

func TestAccVirtualNetworkGateway_lowerCaseSubnetName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway", "test")
	r := VirtualNetworkGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.lowerCaseSubnetName(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku").HasValue("Basic"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualNetworkGateway_vpnGw1(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway", "test")
	r := VirtualNetworkGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.vpnGw1(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccVirtualNetworkGateway_activeActive(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway", "test")
	r := VirtualNetworkGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.activeActive(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccVirtualNetworkGateway_activeActiveZoneRedundantWithP2S(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway", "test")
	r := VirtualNetworkGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.activeActiveZoneRedundantWithP2S(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccVirtualNetworkGateway_standard(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway", "test")
	r := VirtualNetworkGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sku(data, "Standard"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku").HasValue("Standard"),
			),
		},
	})
}

func TestAccVirtualNetworkGateway_vpnGw2(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway", "test")
	r := VirtualNetworkGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sku(data, "VpnGw2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku").HasValue("VpnGw2"),
			),
		},
	})
}

func TestAccVirtualNetworkGateway_vpnGw3(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway", "test")
	r := VirtualNetworkGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sku(data, "VpnGw3"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku").HasValue("VpnGw3"),
			),
		},
	})
}

func TestAccVirtualNetworkGateway_generation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway", "test")
	r := VirtualNetworkGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.generation(data, "Generation2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("generation").HasValue("Generation2"),
			),
		},
	})
}

func TestAccVirtualNetworkGateway_vpnClientConfig(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway", "test")
	r := VirtualNetworkGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.vpnClientConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("vpn_client_configuration.0.radius_server_address").HasValue("1.2.3.4"),
				check.That(data.ResourceName).Key("vpn_client_configuration.0.vpn_client_protocols.#").HasValue("1"),
			),
		},
	})
}

func TestAccVirtualNetworkGateway_vpnClientConfigAzureAdAuth(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway", "test")
	r := VirtualNetworkGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.vpnClientConfigAzureAdAuth(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualNetworkGateway_vpnClientConfigOpenVPN(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway", "test")
	r := VirtualNetworkGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.vpnClientConfigOpenVPN(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("vpn_client_configuration.0.vpn_client_protocols.#").HasValue("1"),
			),
		},
	})
}

func TestAccVirtualNetworkGateway_vpnClientConfigMultipleAuthTypes(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway", "test")
	r := VirtualNetworkGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.vpnClientConfigMultipleAuthTypes(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("vpn_client_configuration.0.vpn_client_protocols.#").HasValue("1"),
				check.That(data.ResourceName).Key("vpn_client_configuration.0.vpn_auth_types.0").HasValue("AAD"),
				check.That(data.ResourceName).Key("vpn_client_configuration.0.vpn_auth_types.1").HasValue("Radius"),
				check.That(data.ResourceName).Key("vpn_client_configuration.0.aad_tenant").IsSet(),
				check.That(data.ResourceName).Key("vpn_client_configuration.0.aad_audience").HasValue("41b23e61-6c1e-4545-b367-cd054e0ed4b4"),
				check.That(data.ResourceName).Key("vpn_client_configuration.0.aad_issuer").IsSet(),
				check.That(data.ResourceName).Key("vpn_client_configuration.0.radius_server.0.address").HasValue("1.2.3.4"),
				check.That(data.ResourceName).Key("vpn_client_configuration.0.radius_server.0.secret").HasValue("1234"),
			),
		},
	})
}

func TestAccVirtualNetworkGateway_enableBgp(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway", "test")
	r := VirtualNetworkGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.enableBgp(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("bgp_settings.#").HasValue("1"),
				check.That(data.ResourceName).Key("bgp_settings.0.peering_addresses.#").HasValue("1"),
				check.That(data.ResourceName).Key("bgp_settings.0.peering_addresses.0.default_addresses.#").HasValue("1"),
				check.That(data.ResourceName).Key("bgp_settings.0.peering_addresses.0.tunnel_ip_addresses.#").HasValue("1"),
				check.That(data.ResourceName).Key("bgp_settings.0.peering_addresses.0.ip_configuration_name").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualNetworkGateway_enableBgpWithAPIPA(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway", "test")
	r := VirtualNetworkGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.enableBgpWithAPIPA(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("bgp_settings.#").HasValue("1"),
				check.That(data.ResourceName).Key("bgp_settings.0.peering_addresses.#").HasValue("1"),
				check.That(data.ResourceName).Key("bgp_settings.0.peering_addresses.0.default_addresses.#").HasValue("1"),
				check.That(data.ResourceName).Key("bgp_settings.0.peering_addresses.0.tunnel_ip_addresses.#").HasValue("1"),
				check.That(data.ResourceName).Key("bgp_settings.0.peering_addresses.0.apipa_addresses.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualNetworkGateway_activeActiveEnableBgpWithAPIPA(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway", "test")
	r := VirtualNetworkGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.activeActiveEnableBgpWithAPIPA(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("bgp_settings.#").HasValue("1"),
				check.That(data.ResourceName).Key("bgp_settings.0.peering_addresses.#").HasValue("2"),
				check.That(data.ResourceName).Key("bgp_settings.0.peering_addresses.0.default_addresses.#").HasValue("1"),
				check.That(data.ResourceName).Key("bgp_settings.0.peering_addresses.1.default_addresses.#").HasValue("1"),
				check.That(data.ResourceName).Key("bgp_settings.0.peering_addresses.0.tunnel_ip_addresses.#").HasValue("1"),
				check.That(data.ResourceName).Key("bgp_settings.0.peering_addresses.1.tunnel_ip_addresses.#").HasValue("1"),
				check.That(data.ResourceName).Key("bgp_settings.0.peering_addresses.0.apipa_addresses.#").HasValue("1"),
				check.That(data.ResourceName).Key("bgp_settings.0.peering_addresses.1.apipa_addresses.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualNetworkGateway_expressRoute(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway", "test")
	r := VirtualNetworkGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.expressRoute(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("type").HasValue("ExpressRoute"),
				check.That(data.ResourceName).Key("bgp_settings.#").HasValue("0"),
			),
		},
	})
}

func TestAccVirtualNetworkGateway_privateIpAddressEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway", "test")
	r := VirtualNetworkGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.privateIpAddressEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.privateIpAddressEnabledUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualNetworkGateway_customRoute(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway", "test")
	r := VirtualNetworkGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.customRoute(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.customRouteUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualNetworkGateway_edgeZone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway", "test")
	r := VirtualNetworkGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.edgeZone(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualNetworkGateway_updateTagsWithBgpSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway", "test")
	r := VirtualNetworkGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.activeActiveEnableBgpWithAPIPA(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.activeActiveEnableBgpWithAPIPAAndTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t VirtualNetworkGatewayResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	gatewayName := state.Attributes["name"]
	resourceGroup := state.Attributes["resource_group_name"]

	resp, err := clients.Network.VnetGatewayClient.Get(ctx, resourceGroup, gatewayName)
	if err != nil {
		return nil, fmt.Errorf("reading Virtual Network Gateway (%s): %+v", state.ID, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (VirtualNetworkGatewayResource) basic(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r VirtualNetworkGatewayResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network_gateway" "import" {
  name                = azurerm_virtual_network_gateway.test.name
  location            = azurerm_virtual_network_gateway.test.location
  resource_group_name = azurerm_virtual_network_gateway.test.resource_group_name
  type                = azurerm_virtual_network_gateway.test.type
  vpn_type            = azurerm_virtual_network_gateway.test.vpn_type
  sku                 = azurerm_virtual_network_gateway.test.sku

  ip_configuration {
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }
}
`, r.basic(data))
}

func (VirtualNetworkGatewayResource) lowerCaseSubnetName(data acceptance.TestData) string {
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

resource "azurerm_subnet" "test" {
  name                 = "gatewaySubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (VirtualNetworkGatewayResource) vpnGw1(data acceptance.TestData) string {
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
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctestvng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "VpnGw1"

  remote_vnet_traffic_enabled           = true
  virtual_wan_traffic_enabled           = true
  bgp_route_translation_for_nat_enabled = true
  ip_sec_replay_protection_enabled      = false

  ip_configuration {
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (VirtualNetworkGatewayResource) activeActive(data acceptance.TestData) string {
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

resource "azurerm_subnet" "test" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "first" {
  name                = "acctestpip1-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_public_ip" "second" {
  name = "acctestpip2-%d"

  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test" {
  depends_on = [
    azurerm_public_ip.first,
    azurerm_public_ip.second,
  ]
  name                = "acctestvng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "VpnGw1"

  active_active = true
  enable_bgp    = true

  ip_configuration {
    name                 = "gw-ip1"
    public_ip_address_id = azurerm_public_ip.first.id

    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }

  ip_configuration {
    name                          = "gw-ip2"
    public_ip_address_id          = azurerm_public_ip.second.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }

  bgp_settings {
    asn = "65010"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (VirtualNetworkGatewayResource) activeActiveZoneRedundantWithP2S(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%[1]d"
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

resource "azurerm_public_ip" "first" {
  name = "acctestpip1-%[1]d"

  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  zones               = ["1", "2", "3"]
}

resource "azurerm_public_ip" "second" {
  name = "acctestpip2-%[1]d"

  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  zones               = ["1", "2", "3"]
}

resource "azurerm_public_ip" "thirth" {
  name = "acctestpip3-%[1]d"

  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  zones               = ["1", "2", "3"]
}

resource "azurerm_virtual_network_gateway" "test" {
  depends_on = [
    azurerm_public_ip.first,
    azurerm_public_ip.second,
    azurerm_public_ip.thirth,
  ]
  name                = "acctestvng-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "VpnGw1AZ"

  active_active = true
  enable_bgp    = true

  ip_configuration {
    name                 = "gw-ip1"
    public_ip_address_id = azurerm_public_ip.first.id

    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }

  ip_configuration {
    name                 = "gw-ip2"
    public_ip_address_id = azurerm_public_ip.second.id

    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }

  ip_configuration {
    name                 = "gw-ip3"
    public_ip_address_id = azurerm_public_ip.thirth.id

    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }

  vpn_client_configuration {
    address_space        = ["10.2.0.0/24"]
    vpn_client_protocols = ["OpenVPN"]

    aad_tenant   = "https://login.microsoftonline.com/%[3]s/"
    aad_audience = "41b23e61-6c1e-4545-b367-cd054e0ed4b4"
    aad_issuer   = "https://sts.windows.net/%[3]s/"
  }

  bgp_settings {
    asn = "65010"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.Client().TenantID)
}

func (VirtualNetworkGatewayResource) vpnClientConfig(data acceptance.TestData) string {
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
}

resource "azurerm_virtual_network_gateway" "test" {
  depends_on          = [azurerm_public_ip.test]
  name                = "acctestvng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "VpnGw1"

  ip_configuration {
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }

  policy_group {
    name       = "TestPolicyGroup"
    is_default = true
    priority   = 1

    policy_member {
      name  = "TestPolicyMember"
      type  = "RadiusAzureGroupId"
      value = "6ad1bd08"
    }
  }

  vpn_client_configuration {
    address_space        = ["10.2.0.0/24"]
    vpn_client_protocols = ["IkeV2"]

    radius_server_address = "1.2.3.4"
    radius_server_secret  = "1234"

    virtual_network_gateway_client_connection {
      name               = "TestConnection"
      policy_group_names = ["TestPolicyGroup"]
      address_prefixes   = ["10.2.0.0/24"]
    }

    ipsec_policy {
      sa_lifetime_in_seconds    = 300
      sa_data_size_in_kilobytes = 1024
      ipsec_encryption          = "AES256"
      ipsec_integrity           = "SHA256"
      ike_encryption            = "AES128"
      ike_integrity             = "SHA256"
      dh_group                  = "DHGroup14"
      pfs_group                 = "PFS14"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (VirtualNetworkGatewayResource) vpnClientConfigAzureAdAuth(data acceptance.TestData) string {
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
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctestvng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "VpnGw1"

  ip_configuration {
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }

  vpn_client_configuration {
    address_space        = ["10.2.0.0/24"]
    vpn_client_protocols = ["OpenVPN"]

    aad_tenant   = "https://login.microsoftonline.com/%s/"
    aad_audience = "41b23e61-6c1e-4545-b367-cd054e0ed4b4"
    aad_issuer   = "https://sts.windows.net/%s/"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.Client().TenantID, data.Client().TenantID)
}

func (VirtualNetworkGatewayResource) vpnClientConfigOpenVPN(data acceptance.TestData) string {
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
}

resource "azurerm_virtual_network_gateway" "test" {
  depends_on          = [azurerm_public_ip.test]
  name                = "acctestvng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "VpnGw1"

  ip_configuration {
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }

  vpn_client_configuration {
    address_space        = ["10.2.0.0/24"]
    vpn_client_protocols = ["OpenVPN"]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (VirtualNetworkGatewayResource) sku(data acceptance.TestData, sku string) string {
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
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctestvng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "%s"

  ip_configuration {
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, sku)
}

func (VirtualNetworkGatewayResource) enableBgp(data acceptance.TestData) string {
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

resource "azurerm_subnet" "test" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip1-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctestvng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type       = "Vpn"
  vpn_type   = "RouteBased"
  sku        = "VpnGw1"
  enable_bgp = true

  ip_configuration {
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (VirtualNetworkGatewayResource) enableBgpWithAPIPA(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ngw-%d"
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
  name                = "acctestpip1-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctestvng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type       = "Vpn"
  vpn_type   = "RouteBased"
  sku        = "VpnGw1"
  enable_bgp = true

  ip_configuration {
    name                          = "gw-ip"
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }

  bgp_settings {
    asn = "65010"
    peering_addresses {
      apipa_addresses = ["169.254.21.1"]
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (VirtualNetworkGatewayResource) expressRoute(data acceptance.TestData) string {
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

resource "azurerm_subnet" "test" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip1-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctestvng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type     = "ExpressRoute"
  vpn_type = "PolicyBased"
  sku      = "Standard"

  ip_configuration {
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (VirtualNetworkGatewayResource) generation(data acceptance.TestData, generation string) string {
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
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctestvng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type       = "Vpn"
  vpn_type   = "RouteBased"
  sku        = "VpnGw2"
  generation = "%s"

  ip_configuration {
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, generation)
}

func (VirtualNetworkGatewayResource) privateIpAddressEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
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

resource "azurerm_subnet" "test" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  zones               = ["1", "2", "3"]
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type                       = "Vpn"
  vpn_type                   = "RouteBased"
  sku                        = "VpnGw1AZ"
  private_ip_address_enabled = true

  custom_route {
    address_prefixes = [
      "101.168.0.6/32"
    ]
  }

  ip_configuration {
    name                          = "vnetGatewayConfig"
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (VirtualNetworkGatewayResource) privateIpAddressEnabledUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
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

resource "azurerm_subnet" "test" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  zones               = ["1", "2", "3"]
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type                       = "Vpn"
  vpn_type                   = "RouteBased"
  sku                        = "VpnGw1AZ"
  private_ip_address_enabled = false

  custom_route {
    address_prefixes = [
      "101.168.0.6/32"
    ]
  }

  ip_configuration {
    name                          = "vnetGatewayConfig"
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (VirtualNetworkGatewayResource) customRoute(data acceptance.TestData) string {
	return fmt.Sprintf(`
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

resource "azurerm_subnet" "test" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  zones               = ["1", "2", "3"]
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type                       = "Vpn"
  vpn_type                   = "RouteBased"
  sku                        = "VpnGw1AZ"
  private_ip_address_enabled = true

  custom_route {
    address_prefixes = [
      "101.168.0.6/32"
    ]
  }

  ip_configuration {
    name                          = "vnetGatewayConfig"
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (VirtualNetworkGatewayResource) customRouteUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
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

resource "azurerm_subnet" "test" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  zones               = ["1", "2", "3"]
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type                       = "Vpn"
  vpn_type                   = "RouteBased"
  sku                        = "VpnGw1AZ"
  private_ip_address_enabled = false

  custom_route {
    address_prefixes = [
      "101.168.0.6/32",
      "101.168.0.7/32"
    ]
  }

  ip_configuration {
    name                          = "vnetGatewayConfig"
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (VirtualNetworkGatewayResource) activeActiveEnableBgpWithAPIPA(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ngw-%d"
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

resource "azurerm_public_ip" "first" {
  name                = "acctestpip1-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_public_ip" "second" {
  name = "acctestpip2-%d"

  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test" {
  depends_on = [
    azurerm_public_ip.first,
    azurerm_public_ip.second,
  ]
  name                = "acctestvng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "VpnGw1"

  active_active = true
  enable_bgp    = true

  ip_configuration {
    name                 = "gw-ip1"
    public_ip_address_id = azurerm_public_ip.first.id

    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }

  ip_configuration {
    name                          = "gw-ip2"
    public_ip_address_id          = azurerm_public_ip.second.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }

  bgp_settings {
    asn = "65010"
    peering_addresses {
      ip_configuration_name = "gw-ip1"
      apipa_addresses       = ["169.254.21.1"]
    }
    peering_addresses {
      ip_configuration_name = "gw-ip2"
      apipa_addresses       = ["169.254.21.2"]
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (VirtualNetworkGatewayResource) vpnClientConfigMultipleAuthTypes(data acceptance.TestData) string {
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
}

resource "azurerm_virtual_network_gateway" "test" {
  depends_on          = [azurerm_public_ip.test]
  name                = "acctestvng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "VpnGw1"

  ip_configuration {
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }

  vpn_client_configuration {
    address_space        = ["10.2.0.0/24"]
    vpn_client_protocols = ["OpenVPN"]
    vpn_auth_types       = ["AAD", "Radius"]

    aad_tenant   = "https://login.microsoftonline.com/%s/"
    aad_audience = "41b23e61-6c1e-4545-b367-cd054e0ed4b4"
    aad_issuer   = "https://sts.windows.net/%s/"

    radius_server {
      address = "1.2.3.4"
      secret  = "1234"
      score   = 2
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.Client().TenantID, data.Client().TenantID)
}

func (VirtualNetworkGatewayResource) edgeZone(data acceptance.TestData) string {
	// @tombuildsstuff: WestUS has an edge zone available - so hard-code to that for now
	data.Locations.Primary = "westus"

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
}

data "azurerm_extended_locations" "test" {
  location = azurerm_resource_group.test.location
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctestvng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  edge_zone           = data.azurerm_extended_locations.test.extended_locations[0]

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "Basic"

  ip_configuration {
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (VirtualNetworkGatewayResource) activeActiveEnableBgpWithAPIPAAndTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ngw-%d"
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

resource "azurerm_public_ip" "first" {
  name                = "acctestpip1-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_public_ip" "second" {
  name = "acctestpip2-%d"

  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test" {
  depends_on = [
    azurerm_public_ip.first,
    azurerm_public_ip.second,
  ]
  name                = "acctestvng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "VpnGw1"

  active_active = true
  enable_bgp    = true

  ip_configuration {
    name                 = "gw-ip1"
    public_ip_address_id = azurerm_public_ip.first.id

    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }

  ip_configuration {
    name                          = "gw-ip2"
    public_ip_address_id          = azurerm_public_ip.second.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }

  bgp_settings {
    asn = "65010"
    peering_addresses {
      ip_configuration_name = "gw-ip1"
      apipa_addresses       = ["169.254.21.1"]
    }
    peering_addresses {
      ip_configuration_name = "gw-ip2"
      apipa_addresses       = ["169.254.21.2"]
    }
  }

  tags = {
    env = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

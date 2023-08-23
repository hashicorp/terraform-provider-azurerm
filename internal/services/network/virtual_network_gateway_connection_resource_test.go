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

type VirtualNetworkGatewayConnectionResource struct{}

func TestAccVirtualNetworkGatewayConnection_sitetosite(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway_connection", "test")
	r := VirtualNetworkGatewayConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sitetosite(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualNetworkGatewayConnection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway_connection", "test")
	r := VirtualNetworkGatewayConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sitetosite(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_virtual_network_gateway_connection"),
		},
	})
}

func TestAccVirtualNetworkGatewayConnection_sitetositeWithoutSharedKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway_connection", "test")
	r := VirtualNetworkGatewayConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sitetositeWithoutSharedKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("shared_key"),
	})
}

func TestAccVirtualNetworkGatewayConnection_vnettonet(t *testing.T) {
	data1 := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway_connection", "test_1")
	data2 := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway_connection", "test_2")
	r := VirtualNetworkGatewayConnectionResource{}

	sharedKey := "4-v3ry-53cr37-1p53c-5h4r3d-k3y"

	data1.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.vnettovnet(data1, data2.RandomInteger, sharedKey),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data1.ResourceName).ExistsInAzure(r),
				acceptance.TestCheckResourceAttr(data1.ResourceName, "shared_key", sharedKey),
				acceptance.TestCheckResourceAttr(data2.ResourceName, "shared_key", sharedKey),
			),
		},
	})
}

func TestAccVirtualNetworkGatewayConnection_ipsecpolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway_connection", "test")
	r := VirtualNetworkGatewayConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.ipsecpolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccVirtualNetworkGatewayConnection_trafficSelectorPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway_connection", "test")
	r := VirtualNetworkGatewayConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.trafficselectorpolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("traffic_selector_policy.0.local_address_cidrs.0").HasValue("10.66.18.0/24"),
				check.That(data.ResourceName).Key("traffic_selector_policy.0.local_address_cidrs.1").HasValue("10.66.17.0/24"),
				check.That(data.ResourceName).Key("traffic_selector_policy.0.remote_address_cidrs.0").HasValue("10.1.1.0/24"),
			),
		},
	})
}

func TestAccVirtualNetworkGatewayConnection_trafficSelectorPolicyMultiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway_connection", "test")
	r := VirtualNetworkGatewayConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.trafficselectorpolicymultiple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("traffic_selector_policy.0.local_address_cidrs.0").HasValue("10.66.18.0/24"),
				check.That(data.ResourceName).Key("traffic_selector_policy.0.local_address_cidrs.1").HasValue("10.66.17.0/24"),
				check.That(data.ResourceName).Key("traffic_selector_policy.0.remote_address_cidrs.0").HasValue("10.1.1.0/24"),
				check.That(data.ResourceName).Key("traffic_selector_policy.1.local_address_cidrs.0").HasValue("10.66.20.0/24"),
				check.That(data.ResourceName).Key("traffic_selector_policy.1.local_address_cidrs.1").HasValue("10.66.19.0/24"),
				check.That(data.ResourceName).Key("traffic_selector_policy.1.remote_address_cidrs.0").HasValue("10.1.2.0/24"),
			),
		},
	})
}

func TestAccVirtualNetworkGatewayConnection_connectionprotocol(t *testing.T) {
	expectedConnectionProtocol := "IKEv1"
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway_connection", "test")
	r := VirtualNetworkGatewayConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.connectionprotocol(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("connection_protocol").HasValue(expectedConnectionProtocol),
			),
		},
	})
}

func TestAccVirtualNetworkGatewayConnection_ConnectionMode(t *testing.T) {
	expectedConnectionMode := "InitiatorOnly"
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway_connection", "test")
	r := VirtualNetworkGatewayConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.connectionMode(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("connection_mode").HasValue(expectedConnectionMode),
			),
		},
	})
}

func TestAccVirtualNetworkGatewayConnection_updatingSharedKey(t *testing.T) {
	data1 := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway_connection", "test_1")
	data2 := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway_connection", "test_2")
	r := VirtualNetworkGatewayConnectionResource{}

	firstSharedKey := "4-v3ry-53cr37-1p53c-5h4r3d-k3y"
	secondSharedKey := "4-r33ly-53cr37-1p53c-5h4r3d-k3y"

	data1.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.vnettovnet(data1, data2.RandomInteger, firstSharedKey),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data1.ResourceName).ExistsInAzure(r),
				check.That(data2.ResourceName).ExistsInAzure(r),
				acceptance.TestCheckResourceAttr(data1.ResourceName, "shared_key", firstSharedKey),
				acceptance.TestCheckResourceAttr(data2.ResourceName, "shared_key", firstSharedKey),
			),
		},
		{
			Config: r.vnettovnet(data1, data2.RandomInteger, secondSharedKey),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data1.ResourceName).ExistsInAzure(r),
				check.That(data2.ResourceName).ExistsInAzure(r),
				acceptance.TestCheckResourceAttr(data1.ResourceName, "shared_key", secondSharedKey),
				acceptance.TestCheckResourceAttr(data2.ResourceName, "shared_key", secondSharedKey),
			),
		},
	})
}

func TestAccVirtualNetworkGatewayConnection_useLocalAzureIpAddressEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway_connection", "test")
	r := VirtualNetworkGatewayConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.useLocalAzureIpAddressEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.useLocalAzureIpAddressEnabledUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualNetworkGatewayConnection_useCustomBgpAddresses(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway_connection", "test")
	r := VirtualNetworkGatewayConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.useCustomBgpAddresses(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("custom_bgp_addresses.0.primary").HasValue("169.254.21.2"),
				check.That(data.ResourceName).Key("custom_bgp_addresses.0.secondary").HasValue("169.254.21.6"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualNetworkGatewayConnection_primaryCustomBgpAddress(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway_connection", "test")
	r := VirtualNetworkGatewayConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.primaryCustomBgpAddress(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccVirtualNetworkGatewayConnection_natRuleIds(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway_connection", "test")
	r := VirtualNetworkGatewayConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.natRuleIds(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t VirtualNetworkGatewayConnectionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	gatewayName := state.Attributes["name"]
	resourceGroup := state.Attributes["resource_group_name"]

	resp, err := clients.Network.VnetGatewayConnectionsClient.Get(ctx, resourceGroup, gatewayName)
	if err != nil {
		return nil, fmt.Errorf("reading Virtual Network Gateway Connection (%s): %+v", state.ID, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (VirtualNetworkGatewayConnectionResource) sitetosite(data acceptance.TestData) string {
	return fmt.Sprintf(`
variable "random" {
  default = "%d"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-${var.random}"
  location = "%s"
}

resource "azurerm_resource_group" "test2" {
  name     = "acctestRG2-${var.random}"
  location = azurerm_resource_group.test.location
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-${var.random}"
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
  name                = "acctest-${var.random}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctest-${var.random}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "Basic"

  ip_configuration {
    name                          = "vnetGatewayConfig"
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }
}

resource "azurerm_local_network_gateway" "test" {
  name                = "acctest-${var.random}"
  location            = azurerm_resource_group.test2.location
  resource_group_name = azurerm_resource_group.test2.name

  gateway_address = "168.62.225.23"
  address_space   = ["10.1.1.0/24"]
}

resource "azurerm_virtual_network_gateway_connection" "test" {
  name                = "acctest-${var.random}"
  location            = azurerm_resource_group.test2.location
  resource_group_name = azurerm_resource_group.test2.name

  type                       = "IPsec"
  virtual_network_gateway_id = azurerm_virtual_network_gateway.test.id
  local_network_gateway_id   = azurerm_local_network_gateway.test.id

  shared_key = "4-v3ry-53cr37-1p53c-5h4r3d-k3y"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (VirtualNetworkGatewayConnectionResource) sitetositeWithoutSharedKey(data acceptance.TestData) string {
	return fmt.Sprintf(`
variable "random" {
  default = "%d"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-${var.random}"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-${var.random}"
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
  name                = "acctest-${var.random}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctest-${var.random}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "Basic"

  ip_configuration {
    name                          = "vnetGatewayConfig"
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }
}

resource "azurerm_local_network_gateway" "test" {
  name                = "acctest-${var.random}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  gateway_address = "168.62.225.23"
  address_space   = ["10.1.1.0/24"]
}

resource "azurerm_virtual_network_gateway_connection" "test" {
  lifecycle {
    ignore_changes = ["shared_key"]
  }
  name                = "acctest-${var.random}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type                       = "IPsec"
  virtual_network_gateway_id = azurerm_virtual_network_gateway.test.id
  local_network_gateway_id   = azurerm_local_network_gateway.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r VirtualNetworkGatewayConnectionResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network_gateway_connection" "import" {
  name                       = azurerm_virtual_network_gateway_connection.test.name
  location                   = azurerm_virtual_network_gateway_connection.test.location
  resource_group_name        = azurerm_virtual_network_gateway_connection.test.resource_group_name
  type                       = azurerm_virtual_network_gateway_connection.test.type
  virtual_network_gateway_id = azurerm_virtual_network_gateway_connection.test.virtual_network_gateway_id
  local_network_gateway_id   = azurerm_virtual_network_gateway_connection.test.local_network_gateway_id
  shared_key                 = azurerm_virtual_network_gateway_connection.test.shared_key
}
`, r.sitetosite(data))
}

func (VirtualNetworkGatewayConnectionResource) vnettovnet(data acceptance.TestData, rInt2 int, sharedKey string) string {
	return fmt.Sprintf(`
variable "random1" {
  default = "%d"
}

variable "random2" {
  default = "%d"
}

variable "shared_key" {
  default = "%s"
}

resource "azurerm_resource_group" "test_1" {
  name     = "acctestRG-${var.random1}"
  location = "%s"
}

resource "azurerm_virtual_network" "test_1" {
  name                = "acctestvn-${var.random1}"
  location            = azurerm_resource_group.test_1.location
  resource_group_name = azurerm_resource_group.test_1.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test_1" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.test_1.name
  virtual_network_name = azurerm_virtual_network.test_1.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "test_1" {
  name                = "acctest-${var.random1}"
  location            = azurerm_resource_group.test_1.location
  resource_group_name = azurerm_resource_group.test_1.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test_1" {
  name                = "acctest-${var.random1}"
  location            = azurerm_resource_group.test_1.location
  resource_group_name = azurerm_resource_group.test_1.name

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "Basic"

  ip_configuration {
    name                          = "vnetGatewayConfig"
    public_ip_address_id          = azurerm_public_ip.test_1.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test_1.id
  }
}

resource "azurerm_virtual_network_gateway_connection" "test_1" {
  name                = "acctest-${var.random1}"
  location            = azurerm_resource_group.test_1.location
  resource_group_name = azurerm_resource_group.test_1.name

  type                            = "Vnet2Vnet"
  virtual_network_gateway_id      = azurerm_virtual_network_gateway.test_1.id
  peer_virtual_network_gateway_id = azurerm_virtual_network_gateway.test_2.id

  shared_key = var.shared_key
}

resource "azurerm_resource_group" "test_2" {
  name     = "acctestRG-${var.random2}"
  location = "%s"
}

resource "azurerm_virtual_network" "test_2" {
  name                = "acctest-${var.random2}"
  location            = azurerm_resource_group.test_2.location
  resource_group_name = azurerm_resource_group.test_2.name
  address_space       = ["10.1.0.0/16"]
}

resource "azurerm_subnet" "test_2" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.test_2.name
  virtual_network_name = azurerm_virtual_network.test_2.name
  address_prefixes     = ["10.1.1.0/24"]
}

resource "azurerm_public_ip" "test_2" {
  name                = "acctest-${var.random2}"
  location            = azurerm_resource_group.test_2.location
  resource_group_name = azurerm_resource_group.test_2.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test_2" {
  name                = "acctest-${var.random2}"
  location            = azurerm_resource_group.test_2.location
  resource_group_name = azurerm_resource_group.test_2.name

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "Basic"

  ip_configuration {
    name                          = "vnetGatewayConfig"
    public_ip_address_id          = azurerm_public_ip.test_2.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test_2.id
  }
}

resource "azurerm_virtual_network_gateway_connection" "test_2" {
  name                = "acctest-${var.random2}"
  location            = azurerm_resource_group.test_2.location
  resource_group_name = azurerm_resource_group.test_2.name

  type                            = "Vnet2Vnet"
  virtual_network_gateway_id      = azurerm_virtual_network_gateway.test_2.id
  peer_virtual_network_gateway_id = azurerm_virtual_network_gateway.test_1.id

  shared_key = var.shared_key
}
`, data.RandomInteger, rInt2, sharedKey, data.Locations.Primary, data.Locations.Secondary)
}

func (VirtualNetworkGatewayConnectionResource) ipsecpolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
variable "random" {
  default = "%d"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-${var.random}"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-${var.random}"
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
  name                = "acctest-${var.random}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctest-${var.random}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "VpnGw1"

  ip_configuration {
    name                          = "vnetGatewayConfig"
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }
}

resource "azurerm_local_network_gateway" "test" {
  name                = "acctest-${var.random}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  gateway_address = "168.62.225.23"
  address_space   = ["10.1.1.0/24"]
}

resource "azurerm_virtual_network_gateway_connection" "test" {
  name                = "acctest-${var.random}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type                       = "IPsec"
  virtual_network_gateway_id = azurerm_virtual_network_gateway.test.id
  local_network_gateway_id   = azurerm_local_network_gateway.test.id

  use_policy_based_traffic_selectors = true
  routing_weight                     = 20

  ipsec_policy {
    dh_group         = "DHGroup14"
    ike_encryption   = "GCMAES256"
    ike_integrity    = "SHA256"
    ipsec_encryption = "AES256"
    ipsec_integrity  = "SHA256"
    pfs_group        = "PFS14"
    sa_datasize      = 102400000
    sa_lifetime      = 27000
  }

  shared_key = "4-v3ry-53cr37-1p53c-5h4r3d-k3y"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (VirtualNetworkGatewayConnectionResource) connectionMode(data acceptance.TestData) string {
	return fmt.Sprintf(`
variable "random" {
  default = "%d"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-${var.random}"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-${var.random}"
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
  name                = "acctest-${var.random}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctest-${var.random}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "VpnGw1"

  ip_configuration {
    name                          = "vnetGatewayConfig"
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }
}

resource "azurerm_local_network_gateway" "test" {
  name                = "acctest-${var.random}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  gateway_address = "168.62.225.23"
  address_space   = ["10.1.1.0/24"]
}

resource "azurerm_virtual_network_gateway_connection" "test" {
  name                = "acctest-${var.random}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type                       = "IPsec"
  virtual_network_gateway_id = azurerm_virtual_network_gateway.test.id
  local_network_gateway_id   = azurerm_local_network_gateway.test.id

  connection_mode = "InitiatorOnly"

  ipsec_policy {
    dh_group         = "DHGroup14"
    ike_encryption   = "AES256"
    ike_integrity    = "SHA256"
    ipsec_encryption = "AES256"
    ipsec_integrity  = "SHA256"
    pfs_group        = "PFS2048"
    sa_datasize      = 102400000
    sa_lifetime      = 27000
  }

  shared_key = "4-v3ry-53cr37-1p53c-5h4r3d-k3y"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (VirtualNetworkGatewayConnectionResource) connectionprotocol(data acceptance.TestData) string {
	return fmt.Sprintf(`
variable "random" {
  default = "%d"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-${var.random}"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-${var.random}"
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
  name                = "acctest-${var.random}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctest-${var.random}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "VpnGw1"

  ip_configuration {
    name                          = "vnetGatewayConfig"
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }
}

resource "azurerm_local_network_gateway" "test" {
  name                = "acctest-${var.random}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  gateway_address = "168.62.225.23"
  address_space   = ["10.1.1.0/24"]
}

resource "azurerm_virtual_network_gateway_connection" "test" {
  name                = "acctest-${var.random}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type                       = "IPsec"
  virtual_network_gateway_id = azurerm_virtual_network_gateway.test.id
  local_network_gateway_id   = azurerm_local_network_gateway.test.id

  connection_protocol = "IKEv1"

  ipsec_policy {
    dh_group         = "DHGroup14"
    ike_encryption   = "AES256"
    ike_integrity    = "SHA256"
    ipsec_encryption = "AES256"
    ipsec_integrity  = "SHA256"
    pfs_group        = "PFS2048"
    sa_datasize      = 102400000
    sa_lifetime      = 27000
  }

  shared_key = "4-v3ry-53cr37-1p53c-5h4r3d-k3y"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (VirtualNetworkGatewayConnectionResource) trafficselectorpolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
variable "random" {
  default = "%d"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-${var.random}"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-${var.random}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.66.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.66.1.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctest-${var.random}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctest-${var.random}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "VpnGw1"

  ip_configuration {
    name                          = "vnetGatewayConfig"
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }
}

resource "azurerm_local_network_gateway" "test" {
  name                = "acctest"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  gateway_address = "168.62.225.23"
  address_space   = ["10.1.1.0/24"]
}

resource "azurerm_virtual_network_gateway_connection" "test" {
  name                = "acctest-${var.random}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type                       = "IPsec"
  virtual_network_gateway_id = azurerm_virtual_network_gateway.test.id
  local_network_gateway_id   = azurerm_local_network_gateway.test.id

  use_policy_based_traffic_selectors = true
  routing_weight                     = 20

  ipsec_policy {
    dh_group         = "DHGroup14"
    ike_encryption   = "AES256"
    ike_integrity    = "SHA256"
    ipsec_encryption = "AES256"
    ipsec_integrity  = "SHA256"
    pfs_group        = "PFS2048"
    sa_datasize      = 102400000
    sa_lifetime      = 27000
  }

  shared_key = "4-v3ry-53cr37-1p53c-5h4r3d-k3y"

  traffic_selector_policy {
    local_address_cidrs  = ["10.66.18.0/24", "10.66.17.0/24"]
    remote_address_cidrs = ["10.1.1.0/24"]
  }

}
`, data.RandomInteger, data.Locations.Primary)
}

func (VirtualNetworkGatewayConnectionResource) trafficselectorpolicymultiple(data acceptance.TestData) string {
	return fmt.Sprintf(`
variable "random" {
  default = "%d"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-${var.random}"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-${var.random}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.66.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.66.1.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctest-${var.random}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctest-${var.random}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "VpnGw1"

  ip_configuration {
    name                          = "vnetGatewayConfig"
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }
}

resource "azurerm_local_network_gateway" "test" {
  name                = "acctest"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  gateway_address = "168.62.225.23"
  address_space   = ["10.1.1.0/24"]
}

resource "azurerm_virtual_network_gateway_connection" "test" {
  name                = "acctest-${var.random}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type                       = "IPsec"
  virtual_network_gateway_id = azurerm_virtual_network_gateway.test.id
  local_network_gateway_id   = azurerm_local_network_gateway.test.id

  use_policy_based_traffic_selectors = true
  routing_weight                     = 20

  ipsec_policy {
    dh_group         = "DHGroup14"
    ike_encryption   = "AES256"
    ike_integrity    = "SHA256"
    ipsec_encryption = "AES256"
    ipsec_integrity  = "SHA256"
    pfs_group        = "PFS2048"
    sa_datasize      = 102400000
    sa_lifetime      = 27000
  }

  shared_key = "4-v3ry-53cr37-1p53c-5h4r3d-k3y"

  traffic_selector_policy {
    local_address_cidrs  = ["10.66.18.0/24", "10.66.17.0/24"]
    remote_address_cidrs = ["10.1.1.0/24"]
  }

  traffic_selector_policy {
    local_address_cidrs  = ["10.66.20.0/24", "10.66.19.0/24"]
    remote_address_cidrs = ["10.1.2.0/24"]
  }

}
`, data.RandomInteger, data.Locations.Primary)
}

func (VirtualNetworkGatewayConnectionResource) useLocalAzureIpAddressEnabled(data acceptance.TestData) string {
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
  name                = "acctestip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  zones               = ["3"]
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctestgw-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type                       = "Vpn"
  vpn_type                   = "RouteBased"
  sku                        = "VpnGw1AZ"
  private_ip_address_enabled = true
  ip_configuration {
    name                          = "vnetGatewayConfig"
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }
}

resource "azurerm_local_network_gateway" "test" {
  name                = "acctestlgw-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  gateway_address = "168.62.225.23"
  address_space   = ["10.1.1.0/24"]
}

resource "azurerm_virtual_network_gateway_connection" "test" {
  name                           = "acctestgwc-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name
  local_azure_ip_address_enabled = true

  type                       = "IPsec"
  virtual_network_gateway_id = azurerm_virtual_network_gateway.test.id
  local_network_gateway_id   = azurerm_local_network_gateway.test.id

  shared_key = "4-v3ry-53cr37-1p53c-5h4r3d-k3y"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (VirtualNetworkGatewayConnectionResource) useLocalAzureIpAddressEnabledUpdate(data acceptance.TestData) string {
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
  name                = "acctestip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  zones               = ["3"]
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctestgw-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "VpnGw1AZ"
  ip_configuration {
    name                          = "vnetGatewayConfig"
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }
}

resource "azurerm_local_network_gateway" "test" {
  name                = "acctestlgw-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  gateway_address = "168.62.225.23"
  address_space   = ["10.1.1.0/24"]
}

resource "azurerm_virtual_network_gateway_connection" "test" {
  name                           = "acctestgwc-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name
  local_azure_ip_address_enabled = false

  type                       = "IPsec"
  virtual_network_gateway_id = azurerm_virtual_network_gateway.test.id
  local_network_gateway_id   = azurerm_local_network_gateway.test.id
  dpd_timeout_seconds        = 30

  shared_key = "4-v3ry-53cr37-1p53c-5h4r3d-k3y"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (VirtualNetworkGatewayConnectionResource) useCustomBgpAddresses(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.1.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.1.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_public_ip" "test2" {
  name                = "acctestip2-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctestgw-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type                       = "Vpn"
  vpn_type                   = "RouteBased"
  enable_bgp                 = true
  active_active              = true
  private_ip_address_enabled = false
  sku                        = "VpnGw2"
  generation                 = "Generation2"

  ip_configuration {
    name                          = "default"
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }

  ip_configuration {
    name                          = "activeactive"
    public_ip_address_id          = azurerm_public_ip.test2.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }

  bgp_settings {
    asn = "65000"

    peering_addresses {
      ip_configuration_name = "default"
      apipa_addresses = [
        "169.254.21.2",
        "169.254.22.2"
      ]
    }

    peering_addresses {
      ip_configuration_name = "activeActive"
      apipa_addresses = [
        "169.254.21.6",
        "169.254.22.6"
      ]
    }
  }
}

resource "azurerm_local_network_gateway" "test" {
  name                = "acctestlgw-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  gateway_address = "168.62.225.23"

  bgp_settings {
    asn                 = "64512"
    bgp_peering_address = "169.254.21.1"
  }
}

resource "azurerm_virtual_network_gateway_connection" "test" {
  name                           = "acctestgwc-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name
  local_azure_ip_address_enabled = false

  type                       = "IPsec"
  virtual_network_gateway_id = azurerm_virtual_network_gateway.test.id
  local_network_gateway_id   = azurerm_local_network_gateway.test.id
  dpd_timeout_seconds        = 30

  shared_key = "4-v3ry-53cr37-1p53c-5h4r3d-k3y"

  enable_bgp = true

  custom_bgp_addresses {
    primary   = "169.254.21.2"
    secondary = "169.254.21.6"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (VirtualNetworkGatewayConnectionResource) natRuleIds(data acceptance.TestData) string {
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
  name                = "acctestvnetgw-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type     = "Vpn"
  vpn_type = "RouteBased"
  sku      = "Basic"

  ip_configuration {
    name                          = "vnetGatewayConfig"
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }
}

data "azurerm_virtual_network_gateway" "test" {
  name                = azurerm_virtual_network_gateway.test.name
  resource_group_name = azurerm_virtual_network_gateway.test.resource_group_name
}

resource "azurerm_local_network_gateway" "test" {
  name                = "acctestlocalnetworkgw-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  gateway_address = "168.62.225.23"
  address_space   = ["10.1.1.0/24"]
}

resource "azurerm_virtual_network_gateway_nat_rule" "test" {
  name                       = "acctestvnetgwegressnatrule-%d"
  resource_group_name        = azurerm_resource_group.test.name
  virtual_network_gateway_id = data.azurerm_virtual_network_gateway.test.id
  mode                       = "EgressSnat"
  type                       = "Dynamic"
  ip_configuration_id        = data.azurerm_virtual_network_gateway.test.ip_configuration.0.id

  external_mapping {
    address_space = "10.1.0.0/26"
  }

  internal_mapping {
    address_space = "10.2.0.0/26"
  }
}

resource "azurerm_virtual_network_gateway_nat_rule" "test2" {
  name                       = "acctestvnetgwingressnatrule-%d"
  resource_group_name        = azurerm_resource_group.test.name
  virtual_network_gateway_id = data.azurerm_virtual_network_gateway.test.id
  mode                       = "IngressSnat"
  type                       = "Dynamic"
  ip_configuration_id        = data.azurerm_virtual_network_gateway.test.ip_configuration.0.id

  external_mapping {
    address_space = "10.7.0.0/26"
  }

  internal_mapping {
    address_space = "10.8.0.0/26"
  }
}

resource "azurerm_virtual_network_gateway_connection" "test" {
  name                = "acctestvnetgwconn-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type                       = "IPsec"
  virtual_network_gateway_id = azurerm_virtual_network_gateway.test.id
  local_network_gateway_id   = azurerm_local_network_gateway.test.id

  egress_nat_rule_ids  = [azurerm_virtual_network_gateway_nat_rule.test.id]
  ingress_nat_rule_ids = [azurerm_virtual_network_gateway_nat_rule.test2.id]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (VirtualNetworkGatewayConnectionResource) primaryCustomBgpAddress(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-vnetgwconn-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.1.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.1.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctestip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctestgw-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type                       = "Vpn"
  vpn_type                   = "RouteBased"
  enable_bgp                 = true
  active_active              = false
  private_ip_address_enabled = false
  sku                        = "VpnGw2"
  generation                 = "Generation2"

  ip_configuration {
    name                          = "default"
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }

  bgp_settings {
    asn = "65000"

    peering_addresses {
      ip_configuration_name = "default"
      apipa_addresses = [
        "169.254.21.2",
        "169.254.22.2"
      ]
    }
  }
}

resource "azurerm_local_network_gateway" "test" {
  name                = "acctestlgw-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  gateway_address = "168.62.225.23"

  bgp_settings {
    asn                 = "64512"
    bgp_peering_address = "169.254.21.1"
  }
}

resource "azurerm_virtual_network_gateway_connection" "test" {
  name                           = "acctestgwc-%d"
  location                       = azurerm_resource_group.test.location
  resource_group_name            = azurerm_resource_group.test.name
  local_azure_ip_address_enabled = false

  type                       = "IPsec"
  virtual_network_gateway_id = azurerm_virtual_network_gateway.test.id
  local_network_gateway_id   = azurerm_local_network_gateway.test.id
  dpd_timeout_seconds        = 30

  shared_key = "4-v3ry-53cr37-1p53c-5h4r3d-k3y"

  enable_bgp = true

  custom_bgp_addresses {
    primary = "169.254.21.2"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualnetworkgatewayconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type VirtualNetworkGatewayConnectionDataSource struct{}

func TestAccDataSourceVirtualNetworkGatewayConnection_siteToSite(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_virtual_network_gateway_connection", "test")
	r := VirtualNetworkGatewayConnectionDataSource{}
	sharedKey := "4-v3ry-53cr37-1p53c-5h4r3d-k3y"

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.siteToSite(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("shared_key").HasValue(sharedKey),
				check.That(data.ResourceName).Key("type").HasValue(string(virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionTypeIPsec)),
			),
		},
	})
}

func TestAccDataSourceVirtualNetworkGatewayConnection_vnetToVnet(t *testing.T) {
	data1 := acceptance.BuildTestData(t, "data.azurerm_virtual_network_gateway_connection", "test_1")
	data2 := acceptance.BuildTestData(t, "data.azurerm_virtual_network_gateway_connection", "test_2")
	r := VirtualNetworkGatewayConnectionDataSource{}

	sharedKey := "4-v3ry-53cr37-1p53c-5h4r3d-k3y"

	data1.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.vnetToVnet(data1, data2.RandomInteger, sharedKey),
			Check: acceptance.ComposeTestCheckFunc(
				acceptance.TestCheckResourceAttr(data1.ResourceName, "shared_key", sharedKey),
				acceptance.TestCheckResourceAttr(data2.ResourceName, "shared_key", sharedKey),
				acceptance.TestCheckResourceAttr(data1.ResourceName, "type", string(virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionTypeVnetTwoVnet)),
				acceptance.TestCheckResourceAttr(data2.ResourceName, "type", string(virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionTypeVnetTwoVnet)),
			),
		},
	})
}

func TestAccDataSourceVirtualNetworkGatewayConnection_ipsecPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_virtual_network_gateway_connection", "test")
	r := VirtualNetworkGatewayConnectionDataSource{}
	sharedKey := "4-v3ry-53cr37-1p53c-5h4r3d-k3y"

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.ipsecPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("shared_key").HasValue(sharedKey),
				check.That(data.ResourceName).Key("type").HasValue(string(virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionTypeIPsec)),
				check.That(data.ResourceName).Key("routing_weight").HasValue("20"),
				check.That(data.ResourceName).Key("ipsec_policy.0.dh_group").HasValue(string(virtualnetworkgatewayconnections.DhGroupDHGroupOneFour)),
				check.That(data.ResourceName).Key("ipsec_policy.0.ike_encryption").HasValue(string(virtualnetworkgatewayconnections.IkeEncryptionAESTwoFiveSix)),
				check.That(data.ResourceName).Key("ipsec_policy.0.ike_integrity").HasValue(string(virtualnetworkgatewayconnections.IkeIntegritySHATwoFiveSix)),
				check.That(data.ResourceName).Key("ipsec_policy.0.ipsec_encryption").HasValue(string(virtualnetworkgatewayconnections.IPsecEncryptionAESTwoFiveSix)),
				check.That(data.ResourceName).Key("ipsec_policy.0.ipsec_integrity").HasValue(string(virtualnetworkgatewayconnections.IPsecIntegritySHATwoFiveSix)),
				check.That(data.ResourceName).Key("ipsec_policy.0.pfs_group").HasValue(string(virtualnetworkgatewayconnections.PfsGroupPFSTwoZeroFourEight)),
				check.That(data.ResourceName).Key("ipsec_policy.0.sa_datasize").HasValue("102400000"),
				check.That(data.ResourceName).Key("ipsec_policy.0.sa_lifetime").HasValue("27000"),
			),
		},
	})
}

func (VirtualNetworkGatewayConnectionDataSource) siteToSite(data acceptance.TestData) string {
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
  type                = "Vpn"
  vpn_type            = "RouteBased"
  sku                 = "Basic"

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
  gateway_address     = "168.62.225.23"
  address_space       = ["10.1.1.0/24"]
}

resource "azurerm_virtual_network_gateway_connection" "test" {
  name                       = "acctest-${var.random}"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  type                       = "IPsec"
  virtual_network_gateway_id = azurerm_virtual_network_gateway.test.id
  local_network_gateway_id   = azurerm_local_network_gateway.test.id
  shared_key                 = "4-v3ry-53cr37-1p53c-5h4r3d-k3y"
}

data "azurerm_virtual_network_gateway_connection" "test" {
  name                = azurerm_virtual_network_gateway_connection.test.name
  resource_group_name = azurerm_virtual_network_gateway_connection.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary)
}

func (VirtualNetworkGatewayConnectionDataSource) vnetToVnet(data acceptance.TestData, rInt2 int, sharedKey string) string {
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

data "azurerm_virtual_network_gateway_connection" "test_1" {
  name                = azurerm_virtual_network_gateway_connection.test_1.name
  resource_group_name = azurerm_virtual_network_gateway_connection.test_1.resource_group_name
}

data "azurerm_virtual_network_gateway_connection" "test_2" {
  name                = azurerm_virtual_network_gateway_connection.test_2.name
  resource_group_name = azurerm_virtual_network_gateway_connection.test_2.resource_group_name
}
`, data.RandomInteger, rInt2, sharedKey, data.Locations.Primary, data.Locations.Secondary)
}

func (VirtualNetworkGatewayConnectionDataSource) ipsecPolicy(data acceptance.TestData) string {
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
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctest-${var.random}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  type                = "Vpn"
  vpn_type            = "RouteBased"
  sku                 = "VpnGw1"

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
  gateway_address     = "168.62.225.23"
  address_space       = ["10.1.1.0/24"]
}

resource "azurerm_virtual_network_gateway_connection" "test" {
  name                               = "acctest-${var.random}"
  location                           = azurerm_resource_group.test.location
  resource_group_name                = azurerm_resource_group.test.name
  type                               = "IPsec"
  virtual_network_gateway_id         = azurerm_virtual_network_gateway.test.id
  local_network_gateway_id           = azurerm_local_network_gateway.test.id
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
}

data "azurerm_virtual_network_gateway_connection" "test" {
  name                = azurerm_virtual_network_gateway_connection.test.name
  resource_group_name = azurerm_virtual_network_gateway_connection.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary)
}

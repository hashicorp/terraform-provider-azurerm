package tests

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccAzureRMDataSourceVirtualNetworkGatewayConnection_sitetosite(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_virtual_network_gateway_connection", "test")
	sharedKey := "4-v3ry-53cr37-1p53c-5h4r3d-k3y"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkGatewayConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataSourceVirtualNetworkGatewayConnection_sitetosite(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayConnectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "shared_key", sharedKey),
					resource.TestCheckResourceAttr(data.ResourceName, "type", string(network.IPsec)),
				),
			},
		},
	})
}

func TestAccAzureRMDataSourceVirtualNetworkGatewayConnection_vnettovnet(t *testing.T) {
	data1 := acceptance.BuildTestData(t, "data.azurerm_virtual_network_gateway_connection", "test_1")
	data2 := acceptance.BuildTestData(t, "data.azurerm_virtual_network_gateway_connection", "test_2")

	sharedKey := "4-v3ry-53cr37-1p53c-5h4r3d-k3y"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkGatewayConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataSourceVirtualNetworkGatewayConnection_vnettovnet(data1, data2.RandomInteger, sharedKey),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayConnectionExists(data1.ResourceName),
					testCheckAzureRMVirtualNetworkGatewayConnectionExists(data2.ResourceName),
					resource.TestCheckResourceAttr(data1.ResourceName, "shared_key", sharedKey),
					resource.TestCheckResourceAttr(data2.ResourceName, "shared_key", sharedKey),
					resource.TestCheckResourceAttr(data1.ResourceName, "type", string(network.Vnet2Vnet)),
					resource.TestCheckResourceAttr(data2.ResourceName, "type", string(network.Vnet2Vnet)),
				),
			},
		},
	})
}

func TestAccAzureRMDataSourceVirtualNetworkGatewayConnection_ipsecpolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_virtual_network_gateway_connection", "test")
	sharedKey := "4-v3ry-53cr37-1p53c-5h4r3d-k3y"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkGatewayConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataSourceVirtualNetworkGatewayConnection_ipsecpolicy(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayConnectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "shared_key", sharedKey),
					resource.TestCheckResourceAttr(data.ResourceName, "type", string(network.IPsec)),
					resource.TestCheckResourceAttr(data.ResourceName, "routing_weight", "20"),
					resource.TestCheckResourceAttr(data.ResourceName, "ipsec_policy.0.dh_group", string(network.DHGroup14)),
					resource.TestCheckResourceAttr(data.ResourceName, "ipsec_policy.0.ike_encryption", string(network.AES256)),
					resource.TestCheckResourceAttr(data.ResourceName, "ipsec_policy.0.ike_integrity", string(network.IkeIntegritySHA256)),
					resource.TestCheckResourceAttr(data.ResourceName, "ipsec_policy.0.ipsec_encryption", string(network.IpsecEncryptionAES256)),
					resource.TestCheckResourceAttr(data.ResourceName, "ipsec_policy.0.ipsec_integrity", string(network.IpsecIntegritySHA256)),
					resource.TestCheckResourceAttr(data.ResourceName, "ipsec_policy.0.pfs_group", string(network.PfsGroupPFS2048)),
					resource.TestCheckResourceAttr(data.ResourceName, "ipsec_policy.0.sa_datasize", "102400000"),
					resource.TestCheckResourceAttr(data.ResourceName, "ipsec_policy.0.sa_lifetime", "27000"),
				),
			},
		},
	})
}

func testAccAzureRMDataSourceVirtualNetworkGatewayConnection_sitetosite(data acceptance.TestData) string {
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
  address_prefix       = "10.0.1.0/24"
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

func testAccAzureRMDataSourceVirtualNetworkGatewayConnection_vnettovnet(data acceptance.TestData, rInt2 int, sharedKey string) string {
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
  address_prefix       = "10.0.1.0/24"
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
  address_prefix       = "10.1.1.0/24"
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

func testAccAzureRMDataSourceVirtualNetworkGatewayConnection_ipsecpolicy(data acceptance.TestData) string {
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
  address_prefix       = "10.0.1.0/24"
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

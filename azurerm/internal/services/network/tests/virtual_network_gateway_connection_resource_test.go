package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMVirtualNetworkGatewayConnection_sitetosite(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway_connection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkGatewayConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetworkGatewayConnection_sitetosite(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVirtualNetworkGatewayConnection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway_connection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkGatewayConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetworkGatewayConnection_sitetosite(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayConnectionExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMVirtualNetworkGatewayConnection_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_virtual_network_gateway_connection"),
			},
		},
	})
}

func TestAccAzureRMVirtualNetworkGatewayConnection_sitetositeWithoutSharedKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway_connection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkGatewayConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetworkGatewayConnection_sitetositeWithoutSharedKey(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVirtualNetworkGatewayConnection_vnettonet(t *testing.T) {
	data1 := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway_connection", "test_1")
	data2 := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway_connection", "test_2")

	sharedKey := "4-v3ry-53cr37-1p53c-5h4r3d-k3y"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkGatewayConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetworkGatewayConnection_vnettovnet(data1, data2.RandomInteger, sharedKey),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayConnectionExists(data1.ResourceName),
					testCheckAzureRMVirtualNetworkGatewayConnectionExists(data2.ResourceName),
					resource.TestCheckResourceAttr(data1.ResourceName, "shared_key", sharedKey),
					resource.TestCheckResourceAttr(data2.ResourceName, "shared_key", sharedKey),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualNetworkGatewayConnection_ipsecpolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway_connection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkGatewayConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetworkGatewayConnection_ipsecpolicy(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayConnectionExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualNetworkGatewayConnection_trafficSelectorPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway_connection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkGatewayConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetworkGatewayConnection_trafficselectorpolicy(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayConnectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "traffic_selector_policy.0.local_address_cidrs.0", "10.66.18.0/24"),
					resource.TestCheckResourceAttr(data.ResourceName, "traffic_selector_policy.0.local_address_cidrs.1", "10.66.17.0/24"),
					resource.TestCheckResourceAttr(data.ResourceName, "traffic_selector_policy.0.remote_address_cidrs.0", "10.1.1.0/24"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualNetworkGatewayConnection_connectionprotocol(t *testing.T) {
	expectedConnectionProtocol := "IKEv1"
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway_connection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkGatewayConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetworkGatewayConnection_connectionprotocol(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayConnectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_protocol", expectedConnectionProtocol),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualNetworkGatewayConnection_updatingSharedKey(t *testing.T) {
	data1 := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway_connection", "test_1")
	data2 := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway_connection", "test_2")

	firstSharedKey := "4-v3ry-53cr37-1p53c-5h4r3d-k3y"
	secondSharedKey := "4-r33ly-53cr37-1p53c-5h4r3d-k3y"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkGatewayConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetworkGatewayConnection_vnettovnet(data1, data2.RandomInteger, firstSharedKey),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayConnectionExists(data1.ResourceName),
					testCheckAzureRMVirtualNetworkGatewayConnectionExists(data2.ResourceName),
					resource.TestCheckResourceAttr(data1.ResourceName, "shared_key", firstSharedKey),
					resource.TestCheckResourceAttr(data2.ResourceName, "shared_key", firstSharedKey),
				),
			},
			{
				Config: testAccAzureRMVirtualNetworkGatewayConnection_vnettovnet(data1, data2.RandomInteger, secondSharedKey),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayConnectionExists(data1.ResourceName),
					testCheckAzureRMVirtualNetworkGatewayConnectionExists(data2.ResourceName),
					resource.TestCheckResourceAttr(data1.ResourceName, "shared_key", secondSharedKey),
					resource.TestCheckResourceAttr(data2.ResourceName, "shared_key", secondSharedKey),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualNetworkGatewayConnection_useLocalAzureIpAddressEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_gateway_connection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkGatewayConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetworkGatewayConnection_useLocalAzureIpAddressEnabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMVirtualNetworkGatewayConnection_useLocalAzureIpAddressEnabledUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayConnectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMVirtualNetworkGatewayConnectionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.VnetGatewayConnectionsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		connectionName := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, connectionName)
		if err != nil {
			return fmt.Errorf("Bad: Get on vnetGatewayConnectionsClient: %+v", err)
		}

		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Bad: Virtual Network Gateway Connection %q (resource group: %q) does not exist", connectionName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMVirtualNetworkGatewayConnectionDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.VnetGatewayConnectionsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_virtual_network_gateway_connection" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			return nil
		}

		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Virtual Network Gateway Connection still exists: %#v", resp.VirtualNetworkGatewayConnectionPropertiesFormat)
		}
	}

	return nil
}

func testAccAzureRMVirtualNetworkGatewayConnection_sitetosite(data acceptance.TestData) string {
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
  name                = "acctest-${var.random}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type                       = "IPsec"
  virtual_network_gateway_id = azurerm_virtual_network_gateway.test.id
  local_network_gateway_id   = azurerm_local_network_gateway.test.id

  shared_key = "4-v3ry-53cr37-1p53c-5h4r3d-k3y"
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualNetworkGatewayConnection_sitetositeWithoutSharedKey(data acceptance.TestData) string {
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
  name                = "acctest-${var.random}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type                       = "IPsec"
  virtual_network_gateway_id = azurerm_virtual_network_gateway.test.id
  local_network_gateway_id   = azurerm_local_network_gateway.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMVirtualNetworkGatewayConnection_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMVirtualNetworkGatewayConnection_sitetosite(data)
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
`, template)
}

func testAccAzureRMVirtualNetworkGatewayConnection_vnettovnet(data acceptance.TestData, rInt2 int, sharedKey string) string {
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
`, data.RandomInteger, rInt2, sharedKey, data.Locations.Primary, data.Locations.Secondary)
}

func testAccAzureRMVirtualNetworkGatewayConnection_ipsecpolicy(data acceptance.TestData) string {
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

func testAccAzureRMVirtualNetworkGatewayConnection_connectionprotocol(data acceptance.TestData) string {
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

func testAccAzureRMVirtualNetworkGatewayConnection_trafficselectorpolicy(data acceptance.TestData) string {
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
  address_prefix       = "10.66.1.0/24"
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

func testAccAzureRMVirtualNetworkGatewayConnection_useLocalAzureIpAddressEnabled(data acceptance.TestData) string {
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
  address_prefix       = "10.0.1.0/24"
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
  sku                        = "VpnGw1"
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
  name                       = "acctestgwc-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  local_azure_ip_address_enabled = true

  type                       = "IPsec"
  virtual_network_gateway_id = azurerm_virtual_network_gateway.test.id
  local_network_gateway_id   = azurerm_local_network_gateway.test.id

  shared_key = "4-v3ry-53cr37-1p53c-5h4r3d-k3y"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMVirtualNetworkGatewayConnection_useLocalAzureIpAddressEnabledUpdate(data acceptance.TestData) string {
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
  address_prefix       = "10.0.1.0/24"
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
  name                = "acctestlgw-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  gateway_address = "168.62.225.23"
  address_space   = ["10.1.1.0/24"]
}

resource "azurerm_virtual_network_gateway_connection" "test" {
  name                       = "acctestgwc-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  local_azure_ip_address_enabled = false

  type                       = "IPsec"
  virtual_network_gateway_id = azurerm_virtual_network_gateway.test.id
  local_network_gateway_id   = azurerm_local_network_gateway.test.id
  dpd_timeout_seconds        = 30

  shared_key = "4-v3ry-53cr37-1p53c-5h4r3d-k3y"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

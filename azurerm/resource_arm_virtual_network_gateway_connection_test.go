package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMVirtualNetworkGatewayConnection_sitetosite(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualNetworkGatewayConnection_sitetosite(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkGatewayConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayConnectionExists("azurerm_virtual_network_gateway_connection.test"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualNetworkGatewayConnection_vnettonet(t *testing.T) {
	firstResourceName := "azurerm_virtual_network_gateway_connection.test_1"
	secondResourceName := "azurerm_virtual_network_gateway_connection.test_2"

	ri := acctest.RandInt()
	ri2 := acctest.RandInt()
	sharedKey := "4-v3ry-53cr37-1p53c-5h4r3d-k3y"
	config := testAccAzureRMVirtualNetworkGatewayConnection_vnettovnet(ri, ri2, sharedKey, testLocation(), testAltLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkGatewayConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayConnectionExists(firstResourceName),
					testCheckAzureRMVirtualNetworkGatewayConnectionExists(secondResourceName),
					resource.TestCheckResourceAttr(firstResourceName, "shared_key", sharedKey),
					resource.TestCheckResourceAttr(secondResourceName, "shared_key", sharedKey),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualNetworkGatewayConnection_ipsecpolicy(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualNetworkGatewayConnection_ipsecpolicy(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkGatewayConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayConnectionExists("azurerm_virtual_network_gateway_connection.test"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualNetworkGatewayConnection_updatingSharedKey(t *testing.T) {
	firstResourceName := "azurerm_virtual_network_gateway_connection.test_1"
	secondResourceName := "azurerm_virtual_network_gateway_connection.test_2"

	ri := acctest.RandInt()
	ri2 := acctest.RandInt()
	loc1 := testLocation()
	loc2 := testAltLocation()

	firstSharedKey := "4-v3ry-53cr37-1p53c-5h4r3d-k3y"
	secondSharedKey := "4-r33ly-53cr37-1p53c-5h4r3d-k3y"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkGatewayConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetworkGatewayConnection_vnettovnet(ri, ri2, firstSharedKey, loc1, loc2),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayConnectionExists(firstResourceName),
					testCheckAzureRMVirtualNetworkGatewayConnectionExists(secondResourceName),
					resource.TestCheckResourceAttr(firstResourceName, "shared_key", firstSharedKey),
					resource.TestCheckResourceAttr(secondResourceName, "shared_key", firstSharedKey),
				),
			},
			{
				Config: testAccAzureRMVirtualNetworkGatewayConnection_vnettovnet(ri, ri2, secondSharedKey, loc1, loc2),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayConnectionExists(firstResourceName),
					testCheckAzureRMVirtualNetworkGatewayConnectionExists(secondResourceName),
					resource.TestCheckResourceAttr(firstResourceName, "shared_key", secondSharedKey),
					resource.TestCheckResourceAttr(secondResourceName, "shared_key", secondSharedKey),
				),
			},
		},
	})
}

func testCheckAzureRMVirtualNetworkGatewayConnectionExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		name, resourceGroup, err := getArmResourceNameAndGroup(s, name)
		if err != nil {
			return err
		}

		client := testAccProvider.Meta().(*ArmClient).vnetGatewayConnectionsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on vnetGatewayConnectionsClient: %+v", err)
		}

		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Bad: Virtual Network Gateway Connection %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMVirtualNetworkGatewayConnectionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).vnetGatewayConnectionsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

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

func testAccAzureRMVirtualNetworkGatewayConnection_sitetosite(rInt int, location string) string {
	return fmt.Sprintf(`
variable "random" {
  default = "%d"
}

resource "azurerm_resource_group" "test" {
  name = "acctestRG-${var.random}"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name = "acctestvn-${var.random}"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name = "GatewaySubnet"
  resource_group_name = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix = "10.0.1.0/24"
}

resource "azurerm_public_ip" "test" {
  name = "acctest-${var.random}"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test" {
  name = "acctest-${var.random}"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  type = "Vpn"
  vpn_type = "RouteBased"
  sku = "Basic"

  ip_configuration {
    name = "vnetGatewayConfig"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
    private_ip_address_allocation = "Dynamic"
    subnet_id = "${azurerm_subnet.test.id}"
  }
}

resource "azurerm_local_network_gateway" "test" {
  name = "acctest-${var.random}"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  gateway_address = "168.62.225.23"
  address_space = ["10.1.1.0/24"]
}

resource "azurerm_virtual_network_gateway_connection" "test" {
  name = "acctest-${var.random}"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  type = "IPsec"
  virtual_network_gateway_id = "${azurerm_virtual_network_gateway.test.id}"
  local_network_gateway_id = "${azurerm_local_network_gateway.test.id}"

  shared_key = "4-v3ry-53cr37-1p53c-5h4r3d-k3y"
}
`, rInt, location)
}

func testAccAzureRMVirtualNetworkGatewayConnection_vnettovnet(rInt, rInt2 int, sharedKey, location, altLocation string) string {
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
  name = "acctestRG-${var.random1}"
  location = "%s"
}

resource "azurerm_virtual_network" "test_1" {
  name = "acctestvn-${var.random1}"
  location = "${azurerm_resource_group.test_1.location}"
  resource_group_name = "${azurerm_resource_group.test_1.name}"
  address_space = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test_1" {
  name = "GatewaySubnet"
  resource_group_name = "${azurerm_resource_group.test_1.name}"
  virtual_network_name = "${azurerm_virtual_network.test_1.name}"
  address_prefix = "10.0.1.0/24"
}

resource "azurerm_public_ip" "test_1" {
  name = "acctest-${var.random1}"
  location = "${azurerm_resource_group.test_1.location}"
  resource_group_name = "${azurerm_resource_group.test_1.name}"
  public_ip_address_allocation = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test_1" {
  name = "acctest-${var.random1}"
  location = "${azurerm_resource_group.test_1.location}"
  resource_group_name = "${azurerm_resource_group.test_1.name}"

  type = "Vpn"
  vpn_type = "RouteBased"
  sku = "Basic"

  ip_configuration {
    name = "vnetGatewayConfig"
    public_ip_address_id = "${azurerm_public_ip.test_1.id}"
    private_ip_address_allocation = "Dynamic"
    subnet_id = "${azurerm_subnet.test_1.id}"
  }
}

resource "azurerm_virtual_network_gateway_connection" "test_1" {
  name = "acctest-${var.random1}"
  location = "${azurerm_resource_group.test_1.location}"
  resource_group_name = "${azurerm_resource_group.test_1.name}"

  type = "Vnet2Vnet"
  virtual_network_gateway_id = "${azurerm_virtual_network_gateway.test_1.id}"
  peer_virtual_network_gateway_id = "${azurerm_virtual_network_gateway.test_2.id}"

  shared_key = "${var.shared_key}"
}

resource "azurerm_resource_group" "test_2" {
  name = "acctestRG-${var.random2}"
  location = "%s"
}

resource "azurerm_virtual_network" "test_2" {
  name = "acctest-${var.random2}"
  location = "${azurerm_resource_group.test_2.location}"
  resource_group_name = "${azurerm_resource_group.test_2.name}"
  address_space = ["10.1.0.0/16"]
}

resource "azurerm_subnet" "test_2" {
  name = "GatewaySubnet"
  resource_group_name = "${azurerm_resource_group.test_2.name}"
  virtual_network_name = "${azurerm_virtual_network.test_2.name}"
  address_prefix = "10.1.1.0/24"
}

resource "azurerm_public_ip" "test_2" {
  name = "acctest-${var.random2}"
  location = "${azurerm_resource_group.test_2.location}"
  resource_group_name = "${azurerm_resource_group.test_2.name}"
  public_ip_address_allocation = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test_2" {
  name = "acctest-${var.random2}"
  location = "${azurerm_resource_group.test_2.location}"
  resource_group_name = "${azurerm_resource_group.test_2.name}"

  type = "Vpn"
  vpn_type = "RouteBased"
  sku = "Basic"

  ip_configuration {
    name = "vnetGatewayConfig"
    public_ip_address_id = "${azurerm_public_ip.test_2.id}"
    private_ip_address_allocation = "Dynamic"
    subnet_id = "${azurerm_subnet.test_2.id}"
  }
}

resource "azurerm_virtual_network_gateway_connection" "test_2" {
  name = "acctest-${var.random2}"
  location = "${azurerm_resource_group.test_2.location}"
  resource_group_name = "${azurerm_resource_group.test_2.name}"

  type = "Vnet2Vnet"
  virtual_network_gateway_id = "${azurerm_virtual_network_gateway.test_2.id}"
  peer_virtual_network_gateway_id = "${azurerm_virtual_network_gateway.test_1.id}"

  shared_key = "${var.shared_key}"
}
`, rInt, rInt2, sharedKey, location, altLocation)
}

func testAccAzureRMVirtualNetworkGatewayConnection_ipsecpolicy(rInt int, location string) string {
	return fmt.Sprintf(`
variable "random" {
  default = "%d"
}

resource "azurerm_resource_group" "test" {
  name = "acctestRG-${var.random}"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name = "acctestvn-${var.random}"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name = "GatewaySubnet"
  resource_group_name = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix = "10.0.1.0/24"
}

resource "azurerm_public_ip" "test" {
  name = "acctest-${var.random}"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test" {
  name = "acctest-${var.random}"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  type = "Vpn"
  vpn_type = "RouteBased"
  sku = "VpnGw1"

  ip_configuration {
    name = "vnetGatewayConfig"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
    private_ip_address_allocation = "Dynamic"
    subnet_id = "${azurerm_subnet.test.id}"
  }
}

resource "azurerm_local_network_gateway" "test" {
  name = "acctest-${var.random}"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  gateway_address = "168.62.225.23"
  address_space = ["10.1.1.0/24"]
}

resource "azurerm_virtual_network_gateway_connection" "test" {
  name = "acctest-${var.random}"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  type = "IPsec"
  virtual_network_gateway_id = "${azurerm_virtual_network_gateway.test.id}"
  local_network_gateway_id = "${azurerm_local_network_gateway.test.id}"

  use_policy_based_traffic_selectors = true
  routing_weight = 20

  ipsec_policy {
    dh_group = "DHGroup14"
    ike_encryption = "AES256"
    ike_integrity = "SHA256"
    ipsec_encryption = "AES256"
    ipsec_integrity = "SHA256"
    pfs_group = "PFS2048"
    sa_datasize = 102400000
    sa_lifetime = 27000
  }

  shared_key = "4-v3ry-53cr37-1p53c-5h4r3d-k3y"
}
`, rInt, location)
}

package azurerm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"testing"
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

func TestAccAzureRMVirtualNetworkGatewayConnection_vnettovnet(t *testing.T) {
	ri := acctest.RandInt()
	ri2 := acctest.RandInt()
	config := testAccAzureRMVirtualNetworkGatewayConnection_vnettovnet(ri, ri2, testLocation(), testAltLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkGatewayConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayConnectionExists("azurerm_virtual_network_gateway_connection.test_1"),
					testCheckAzureRMVirtualNetworkGatewayConnectionExists("azurerm_virtual_network_gateway_connection.test_2"),
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
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%[1]d"
    location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name = "test-%[1]d"
  location = "%[2]s"
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
    name = "test-%[1]d"
    location = "%[2]s"
    resource_group_name = "${azurerm_resource_group.test.name}"
    public_ip_address_allocation = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test" {
  name = "test-%[1]d"
  location = "%[2]s"
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
    name = "test-%[1]d"
    location = "%[2]s"
    resource_group_name = "${azurerm_resource_group.test.name}"

    gateway_address = "168.62.225.23"
    address_space = ["10.1.1.0/24"]
}

resource "azurerm_virtual_network_gateway_connection" "test" {
    name = "test-%[1]d"
    location = "%[2]s"
    resource_group_name = "${azurerm_resource_group.test.name}"

    type = "IPsec"
    virtual_network_gateway_id = "${azurerm_virtual_network_gateway.test.id}"
    local_network_gateway_id = "${azurerm_local_network_gateway.test.id}"

    shared_key = "4-v3ry-53cr37-1p53c-5h4r3d-k3y"
}
`, rInt, location)
}

func testAccAzureRMVirtualNetworkGatewayConnection_vnettovnet(rInt int, rInt2 int, location string, altLocation string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test_1" {
    name = "acctestRG-%[1]d"
    location = "%[3]s"
}

resource "azurerm_virtual_network" "test_1" {
  name = "acctest-%[1]d"
  location = "%[3]s"
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
  name = "acctest-%[1]d"
  location = "%[3]s"
  resource_group_name = "${azurerm_resource_group.test_1.name}"
  public_ip_address_allocation = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test_1" {
  name = "acctest-%[1]d"
  location = "%[3]s"
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
  name = "acctest-%[1]d"
  location = "%[3]s"
  resource_group_name = "${azurerm_resource_group.test_1.name}"

  type = "Vnet2Vnet"
  virtual_network_gateway_id = "${azurerm_virtual_network_gateway.test_1.id}"
  peer_virtual_network_gateway_id = "${azurerm_virtual_network_gateway.test_2.id}"

  shared_key = "4-v3ry-53cr37-1p53c-5h4r3d-k3y"
}

resource "azurerm_resource_group" "test_2" {
    name = "acctestRG-%[2]d"
    location = "%[4]s"
}

resource "azurerm_virtual_network" "test_2" {
  name = "acctest-%[2]d"
  location = "%[4]s"
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
  name = "acctest-%[2]d"
  location = "%[4]s"
  resource_group_name = "${azurerm_resource_group.test_2.name}"
  public_ip_address_allocation = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test_2" {
  name = "acctest-%[2]d"
  location = "%[4]s"
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
  name = "acctest-%[2]d"
  location = "%[4]s"
  resource_group_name = "${azurerm_resource_group.test_2.name}"

  type = "Vnet2Vnet"
  virtual_network_gateway_id = "${azurerm_virtual_network_gateway.test_2.id}"
  peer_virtual_network_gateway_id = "${azurerm_virtual_network_gateway.test_1.id}"

  shared_key = "4-v3ry-53cr37-1p53c-5h4r3d-k3y"
}
`, rInt, rInt2, location, altLocation)
}

package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMVirtualNetworkGateway_basic(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := "azurerm_virtual_network_gateway.test"
	config := testAccAzureRMVirtualNetworkGateway_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku", "Basic"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualNetworkGateway_lowerCaseSubnetName(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := "azurerm_virtual_network_gateway.test"
	config := testAccAzureRMVirtualNetworkGateway_lowerCaseSubnetName(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "sku", "Basic"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualNetworkGateway_vpnGw1(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualNetworkGateway_vpnGw1(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayExists("azurerm_virtual_network_gateway.test"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualNetworkGateway_activeActive(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualNetworkGateway_activeActive(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayExists("azurerm_virtual_network_gateway.test"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualNetworkGateway_vpnClientConfig(t *testing.T) {
	ri := acctest.RandInt()
	resourceName := "azurerm_virtual_network_gateway.test"
	config := testAccAzureRMVirtualNetworkGateway_vpnClientConfig(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "vpn_client_configuration.0.radius_server_address", "1.2.3.4"),
					resource.TestCheckResourceAttr(resourceName, "vpn_client_configuration.0.vpn_client_protocols.#", "2"),
				),
			},
		},
	})
}

func testCheckAzureRMVirtualNetworkGatewayExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		name, resourceGroup, err := getArmResourceNameAndGroup(s, name)
		if err != nil {
			return err
		}

		client := testAccProvider.Meta().(*ArmClient).vnetGatewayClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on vnetGatewayClient: %+v", err)
		}

		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Bad: Virtual Network Gateway %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMVirtualNetworkGatewayDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).vnetGatewayClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_virtual_network_gateway" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			return nil
		}

		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Virtual Network Gateway still exists:\n%#v", resp.VirtualNetworkGatewayPropertiesFormat)
		}
	}

	return nil
}

func testAccAzureRMVirtualNetworkGateway_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name = "acctestvn-%d"
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
  name = "acctestpip-%d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test" {
  name = "acctestvng-%d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  type = "Vpn"
  vpn_type = "RouteBased"
  sku = "Basic"

  ip_configuration {
    public_ip_address_id = "${azurerm_public_ip.test.id}"
    private_ip_address_allocation = "Dynamic"
    subnet_id = "${azurerm_subnet.test.id}"
  }
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMVirtualNetworkGateway_lowerCaseSubnetName(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name = "acctestvn-%d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name = "gatewaySubnet"
  resource_group_name = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix = "10.0.1.0/24"
}

resource "azurerm_public_ip" "test" {
  name = "acctestpip-%d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test" {
  name = "acctestvng-%d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  type = "Vpn"
  vpn_type = "RouteBased"
  sku = "Basic"

  ip_configuration {
    public_ip_address_id = "${azurerm_public_ip.test.id}"
    private_ip_address_allocation = "Dynamic"
    subnet_id = "${azurerm_subnet.test.id}"
  }
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMVirtualNetworkGateway_vpnGw1(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name = "acctestvn-%d"
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
  name = "acctestpip-%d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test" {
  name = "acctestvng-%d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  type = "Vpn"
  vpn_type = "RouteBased"
  sku = "VpnGw1"

  ip_configuration {
    public_ip_address_id = "${azurerm_public_ip.test.id}"
    private_ip_address_allocation = "Dynamic"
    subnet_id = "${azurerm_subnet.test.id}"
  }
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMVirtualNetworkGateway_activeActive(rInt int, location string) string {

	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name = "acctestvn-%d"
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


resource "azurerm_public_ip" "first" {
  name = "acctestpip1-%d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "Dynamic"
}

resource "azurerm_public_ip" "second" {
  name = "acctestpip2-%d"

  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test" {
  depends_on = ["azurerm_public_ip.first", "azurerm_public_ip.second"]
  name = "acctestvng-%d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  type = "Vpn"
  vpn_type = "RouteBased"
  sku = "VpnGw1"

  active_active = true
  enable_bgp = true

  ip_configuration {
    name = "gw-ip1"
    public_ip_address_id = "${azurerm_public_ip.first.id}"

    private_ip_address_allocation = "Dynamic"
    subnet_id = "${azurerm_subnet.test.id}"
  }


  ip_configuration {
    name = "gw-ip2"
    public_ip_address_id = "${azurerm_public_ip.second.id}"
    private_ip_address_allocation = "Dynamic"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  bgp_settings {
    asn = "65010"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt)

}

func testAccAzureRMVirtualNetworkGateway_vpnClientConfig(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name = "acctestvn-%d"
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
  name = "acctestpip-%d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test" {
  depends_on = ["azurerm_public_ip.test"]
  name = "acctestvng-%d"
  location = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  type = "Vpn"
  vpn_type = "RouteBased"
  sku = "VpnGw1"

  ip_configuration {
    public_ip_address_id = "${azurerm_public_ip.test.id}"
    private_ip_address_allocation = "Dynamic"
    subnet_id = "${azurerm_subnet.test.id}"
  }

  vpn_client_configuration {
	address_space = ["10.2.0.0/24"]
	vpn_client_protocols = ["SSTP", "IkeV2"]

	radius_server_address = "1.2.3.4"
    radius_server_secret = "1234"
  }
}
`, rInt, location, rInt, rInt, rInt)
}

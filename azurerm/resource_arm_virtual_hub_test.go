package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMVirtualHub_basic(t *testing.T) {
	resourceName := "azurerm_virtual_hub.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHub_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMVirtualHub_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_virtual_hub.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHub_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMVirtualHub_requiresImport(ri, testLocation()),
				ExpectError: testRequiresImportError("azurerm_virtual_hub"),
			},
		},
	})
}

func TestAccAzureRMVirtualHub_complete(t *testing.T) {
	resourceName := "azurerm_virtual_hub.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHub_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "address_prefix", "10.0.2.0/24"),
					resource.TestCheckResourceAttr(resourceName, "virtual_network_connection.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "virtual_network_connection.0.name", "testConnection"),
					resource.TestCheckResourceAttr(resourceName, "virtual_network_connection.0.allow_hub_to_remote_vnet_transit", "false"),
					resource.TestCheckResourceAttr(resourceName, "virtual_network_connection.0.allow_remote_vnet_to_use_hub_vnet_gateways", "false"),
					resource.TestCheckResourceAttr(resourceName, "virtual_network_connection.0.enable_internet_security", "false"),
					resource.TestCheckResourceAttr(resourceName, "route.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "route.0.address_prefixes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "route.0.address_prefixes.0", "10.0.3.0/24"),
					resource.TestCheckResourceAttr(resourceName, "route.0.next_hop_ip_address", "10.0.5.6"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.ENV", "prod"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMVirtualHub_update(t *testing.T) {
	resourceName := "azurerm_virtual_hub.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHub_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "address_prefix", "10.0.1.0/24"),
					resource.TestCheckResourceAttr(resourceName, "virtual_network_connection.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "route.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMVirtualHub_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "address_prefix", "10.0.2.0/24"),
					resource.TestCheckResourceAttr(resourceName, "virtual_network_connection.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "route.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "route.0.address_prefixes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "route.0.address_prefixes.0", "10.0.3.0/24"),
					resource.TestCheckResourceAttr(resourceName, "route.0.next_hop_ip_address", "10.0.5.6"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMVirtualHubExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Virtual Hub not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := testAccProvider.Meta().(*ArmClient).Network.VirtualHubClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		if resp, err := client.Get(ctx, resourceGroup, name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Virtual Hub %q (Resource Group %q) does not exist", name, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on network.VirtualHubClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMVirtualHubDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).Network.VirtualHubClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_virtual_hub" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on network.VirtualHubClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMVirtualHub_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-%d"
  location = "%s"
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctest-VirtualWan-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_virtual_hub" "test" {
  name                = "acctest-VirtualHub-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  address_prefix      = "10.0.1.0/24"
  virtual_wan_id      = "${azurerm_virtual_wan.test.id}"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMVirtualHub_requiresImport(rInt int, location string) string {
	return fmt.Sprintf(`
%s
resource "azurerm_virtual_hub" "import" {
  name                = "${azurerm_virtual_hub.test.name}"
  location            = "${azurerm_virtual_hub.test.location}"
  resource_group_name = "${azurerm_virtual_hub.test.name}"
}
}
`, testAccAzureRMVirtualHub_basic(rInt, location))
}

func testAccAzureRMVirtualHub_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VirtualNetwork-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_network_security_group" "test" {
  name                = "acctest-NetworkSecurityGroup-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                      = "acctest-Subnet-%d"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  virtual_network_name      = "${azurerm_virtual_network.test.name}"
  address_prefix            = "10.5.1.0/24"
  network_security_group_id = "${azurerm_network_security_group.test.id}"
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctest-VirtualWan-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_virtual_hub" "test" {
  name                = "acctest-VirtualHub-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  address_prefix      = "10.0.2.0/24"

  virtual_wan_id  = "${azurerm_virtual_wan.test.id}"

  virtual_network_connection {
    name                                       = "testConnection"
    remote_virtual_network_id                  = "${azurerm_virtual_network.test.id}"
    allow_hub_to_remote_vnet_transit           = "false"
    allow_remote_vnet_to_use_hub_vnet_gateways = "false"
    enable_internet_security                   = "false"
  }

  route {
    address_prefixes    = ["10.0.3.0/24"]
    next_hop_ip_address = "10.0.5.6"
  }

  tags = {
    ENV = "prod"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt)
}

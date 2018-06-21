package azurerm

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
)

func init() {
	resource.AddTestSweepers("azurerm_virtual_network", &resource.Sweeper{
		Name: "azurerm_virtual_network",
		F:    testSweepVirtualNetworks,
		Dependencies: []string{
			"azurerm_application_gateway",
			"azurerm_subnet",
			"azurerm_network_interface",
			"azurerm_virtual_machine",
		},
	})
}

func testSweepVirtualNetworks(region string) error {
	armClient, err := buildConfigForSweepers()
	if err != nil {
		return err
	}

	client := (*armClient).vnetClient
	ctx := (*armClient).StopContext

	log.Printf("Retrieving the Virtual Networks..")
	results, err := client.ListAll(ctx)
	if err != nil {
		return fmt.Errorf("Error Listing on Virtual Networks: %+v", err)
	}

	for _, network := range results.Values() {
		id, err := parseAzureResourceID(*network.ID)
		if err != nil {
			return fmt.Errorf("Error parsing Azure Resource ID %q", id)
		}

		resourceGroupName := id.ResourceGroup
		name := *network.Name
		location := *network.Location

		if !shouldSweepAcceptanceTestResource(name, location, region) {
			continue
		}

		log.Printf("Deleting Virtual Network %q", name)
		future, err := client.Delete(ctx, resourceGroupName, name)
		if err != nil {
			if response.WasNotFound(future.Response()) {
				continue
			}

			return err
		}

		err = future.WaitForCompletion(ctx, client.Client)
		if err != nil {
			if response.WasNotFound(future.Response()) {
				continue
			}

			return err
		}
	}

	return nil
}

func TestAccAzureRMVirtualNetwork_basic(t *testing.T) {
	resourceName := "azurerm_virtual_network.test"
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualNetwork_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualNetwork_disappears(t *testing.T) {
	resourceName := "azurerm_virtual_network.test"
	ri := acctest.RandInt()
	config := testAccAzureRMVirtualNetwork_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkExists(resourceName),
					testCheckAzureRMVirtualNetworkDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMVirtualNetwork_withTags(t *testing.T) {
	resourceName := "azurerm_virtual_network.test"
	location := testLocation()
	ri := acctest.RandInt()
	preConfig := testAccAzureRMVirtualNetwork_withTags(ri, location)
	postConfig := testAccAzureRMVirtualNetwork_withTagsUpdated(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkExists(resourceName),
					resource.TestCheckResourceAttr(
						resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(
						resourceName, "tags.environment", "Production"),
					resource.TestCheckResourceAttr(resourceName, "tags.cost_center", "MSFT"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkExists(resourceName),
					resource.TestCheckResourceAttr(
						resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(
						resourceName, "tags.environment", "staging"),
				),
			},
		},
	})
}

func TestAccAzureRMVirtualNetwork_bug373(t *testing.T) {
	resourceName := "azurerm_virtual_network.test"
	rs := acctest.RandString(6)
	config := testAccAzureRMVirtualNetwork_bug373(rs, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkExists(resourceName),
				),
			},
		},
	})
}

func testCheckAzureRMVirtualNetworkExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		virtualNetworkName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for virtual network: %s", virtualNetworkName)
		}

		// Ensure resource group/virtual network combination exists in API
		client := testAccProvider.Meta().(*ArmClient).vnetClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, virtualNetworkName, "")
		if err != nil {
			return fmt.Errorf("Bad: Get on vnetClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Virtual Network %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMVirtualNetworkDisappears(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		virtualNetworkName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for virtual network: %s", virtualNetworkName)
		}

		// Ensure resource group/virtual network combination exists in API
		client := testAccProvider.Meta().(*ArmClient).vnetClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		future, err := client.Delete(ctx, resourceGroup, virtualNetworkName)
		if err != nil {
			return fmt.Errorf("Error deleting Virtual Network %q (RG %q): %+v", virtualNetworkName, resourceGroup, err)
		}

		err = future.WaitForCompletion(ctx, client.Client)
		if err != nil {
			return fmt.Errorf("Error waiting for deletion of Virtual Network %q (RG %q): %+v", virtualNetworkName, resourceGroup, err)
		}

		return nil
	}
}

func testCheckAzureRMVirtualNetworkDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).vnetClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_virtual_network" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name, "")

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Virtual Network still exists:\n%#v", resp.VirtualNetworkPropertiesFormat)
		}
	}

	return nil
}

func testAccAzureRMVirtualNetwork_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctestvirtnet%d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    subnet {
        name = "subnet1"
        address_prefix = "10.0.1.0/24"
    }
}
`, rInt, location, rInt)
}

func testAccAzureRMVirtualNetwork_withTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctestvirtnet%d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    subnet {
        name = "subnet1"
        address_prefix = "10.0.1.0/24"
    }

    tags {
	environment = "Production"
	cost_center = "MSFT"
    }
}
`, rInt, location, rInt)
}

func testAccAzureRMVirtualNetwork_withTagsUpdated(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "%s"
}

resource "azurerm_virtual_network" "test" {
    name = "acctestvirtnet%d"
    address_space = ["10.0.0.0/16"]
    location = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"

    subnet {
        name = "subnet1"
        address_prefix = "10.0.1.0/24"
    }

    tags {
	environment = "staging"
    }
}
`, rInt, location, rInt)
}

func testAccAzureRMVirtualNetwork_bug373(rString string, location string) string {
	return fmt.Sprintf(`
variable "environment" {
  default = "TestVirtualNetworkBug373"
}

variable "network_cidr" {
  default = "10.0.0.0/16"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG%s"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "${azurerm_resource_group.test.name}-vnet"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["${var.network_cidr}"]
  location            = "${azurerm_resource_group.test.location}"

  tags {
    environment = "${var.environment}"
  }
}

resource "azurerm_subnet" "public" {
  name                      = "${azurerm_resource_group.test.name}-subnet-public"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  virtual_network_name      = "${azurerm_virtual_network.test.name}"
  address_prefix            = "10.0.1.0/24"
  network_security_group_id = "${azurerm_network_security_group.test.id}"
}

resource "azurerm_subnet" "private" {
  name                      = "${azurerm_resource_group.test.name}-subnet-private"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  virtual_network_name      = "${azurerm_virtual_network.test.name}"
  address_prefix            = "10.0.2.0/24"
  network_security_group_id = "${azurerm_network_security_group.test.id}"
}

resource "azurerm_network_security_group" "test" {
  name                = "default-network-sg"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  security_rule {
    name                       = "default-allow-all"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "*"
    source_port_range          = "*"
    destination_port_range     = "*"
    source_address_prefix      = "${var.network_cidr}"
    destination_address_prefix = "*"
  }

  tags {
    environment = "${var.environment}"
  }
}
`, rString, location)
}

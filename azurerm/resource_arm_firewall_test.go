package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAzureRMFirewall_basic(t *testing.T) {
	resourceName := "azurerm_firewall.test"
	ri := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewall_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "ip_configuration.0.name", "configuration"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_configuration.0.private_ip_address"),
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

func TestAccAzureRMFirewall_withTags(t *testing.T) {
	resourceName := "azurerm_firewall.test"
	ri := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewall_withTags(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "Production"),
					resource.TestCheckResourceAttr(resourceName, "tags.cost_center", "MSFT"),
				),
			},
			{
				Config: testAccAzureRMFirewall_withUpdatedTags(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "staging"),
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

func TestAccAzureRMFirewall_disappears(t *testing.T) {
	resourceName := "azurerm_firewall.test"
	ri := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFirewall_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFirewallExists(resourceName),
					testCheckAzureRMFirewallDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMFirewallExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Azure Firewall: %q", name)
		}

		client := testAccProvider.Meta().(*ArmClient).azureFirewallsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Azure Firewall %q (Resource Group: %q) does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on azureFirewallsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMFirewallDisappears(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Azure Firewall: %q", name)
		}

		client := testAccProvider.Meta().(*ArmClient).azureFirewallsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		future, err := client.Delete(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Delete on azureFirewallsClient: %+v", err)
		}
		err = future.WaitForCompletionRef(ctx, client.Client)
		if err != nil {
			return fmt.Errorf("Bad: waiting for Deletion on azureFirewallsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMFirewallDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).azureFirewallsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_firewall" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Firewall still exists:\n%#v", resp.AzureFirewallPropertiesFormat)
	}

	return nil
}

func testAccAzureRMFirewall_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "AzureFirewallSubnet"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_public_ip" "test" {
  name                         = "acctestpip%d"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "Static"
  sku                          = "Standard"
}

resource "azurerm_firewall" "test" {
  name                = "acctestfirewall%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "configuration"
    subnet_id                     = "${azurerm_subnet.test.id}"
    internal_public_ip_address_id = "${azurerm_public_ip.test.id}"
  }
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMFirewall_withTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "AzureFirewallSubnet"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_public_ip" "test" {
  name                         = "acctestpip%d"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "Static"
  sku                          = "Standard"
}

resource "azurerm_firewall" "test" {
  name                = "acctestfirewall%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "configuration"
    subnet_id                     = "${azurerm_subnet.test.id}"
    internal_public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  tags {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMFirewall_withUpdatedTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "AzureFirewallSubnet"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_public_ip" "test" {
  name                         = "acctestpip%d"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "Static"
  sku                          = "Standard"
}

resource "azurerm_firewall" "test" {
  name                = "acctestfirewall%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  ip_configuration {
    name                          = "configuration"
    subnet_id                     = "${azurerm_subnet.test.id}"
    internal_public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  tags {
    environment = "staging"
  }
}
`, rInt, location, rInt, rInt, rInt)
}

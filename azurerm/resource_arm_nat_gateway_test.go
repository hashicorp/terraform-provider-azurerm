package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMNatGateway_basic(t *testing.T) {
	resourceName := "azurerm_nat_gateway.test"
	ri := tf.AccRandTimeInt()
	// Using alt location because the resource currently in private preview and is only available in eastus2.
	location := testAltLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNatGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNatGateway_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNatGatewayExists(resourceName),
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

func TestAccAzureRMNatGateway_complete(t *testing.T) {
	resourceName := "azurerm_nat_gateway.test"
	ri := tf.AccRandTimeInt()
	// Using alt location because the resource currently in private preview and is only available in eastus2.
	location := testAltLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNatGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNatGateway_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNatGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "public_ip_address_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "public_ip_prefix_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "sku_name", "Standard"),
					resource.TestCheckResourceAttr(resourceName, "idle_timeout_in_minutes", "10"),
					resource.TestCheckResourceAttr(resourceName, "zones.#", "1"),
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

func TestAccAzureRMNatGateway_update(t *testing.T) {
	resourceName := "azurerm_nat_gateway.test"
	ri := tf.AccRandTimeInt()
	// Using alt location because the resource currently in private preview and is only available in eastus2.
	location := testAltLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNatGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNatGateway_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNatGatewayExists(resourceName),
				),
			},
			{
				Config: testAccAzureRMNatGateway_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNatGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "public_ip_address_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "public_ip_prefix_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "sku_name", "Standard"),
					resource.TestCheckResourceAttr(resourceName, "idle_timeout_in_minutes", "10"),
					resource.TestCheckResourceAttr(resourceName, "zones.#", "1"),
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

func testCheckAzureRMNatGatewayExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Nat Gateway not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := testAccProvider.Meta().(*ArmClient).Network.NatGatewayClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		if resp, err := client.Get(ctx, resourceGroup, name, ""); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Nat Gateway %q (Resource Group %q) does not exist", name, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on network.NatGatewayClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMNatGatewayDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).Network.NatGatewayClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_nat_gateway" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, name, ""); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on network.NatGatewayClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMNatGateway_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-%d"
  location = "%s"
}

resource "azurerm_nat_gateway" "test" {
  name                = "acctestnatGateway-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt)
}

func testAccAzureRMNatGateway_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicIP-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
  sku                 = "Standard"
  zones               = ["1"]
}

resource "azurerm_public_ip_prefix" "test" {
  name                = "acctestpublicIPPrefix-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  prefix_length       = 30
  zones               = ["1"]
}

resource "azurerm_nat_gateway" "test" {
  name                    = "acctestnatGateway-%d"
  location                = "${azurerm_resource_group.test.location}"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  public_ip_address_ids   = ["${azurerm_public_ip.test.id}"]
  public_ip_prefix_ids    = ["${azurerm_public_ip_prefix.test.id}"]
  sku_name                = "Standard"
  idle_timeout_in_minutes = 10
  zones                   = ["1"]
}
`, rInt, location, rInt, rInt, rInt)
}

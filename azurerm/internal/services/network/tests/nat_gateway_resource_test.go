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

func TestAccAzureRMNatGateway_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_nat_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNatGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNatGateway_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNatGatewayExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNatGateway_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_nat_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNatGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNatGateway_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNatGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "public_ip_address_ids.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "public_ip_prefix_ids.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "Standard"),
					resource.TestCheckResourceAttr(data.ResourceName, "idle_timeout_in_minutes", "10"),
					resource.TestCheckResourceAttr(data.ResourceName, "zones.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNatGateway_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_nat_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNatGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNatGateway_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNatGatewayExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMNatGateway_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNatGatewayExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "public_ip_address_ids.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "public_ip_prefix_ids.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "Standard"),
					resource.TestCheckResourceAttr(data.ResourceName, "idle_timeout_in_minutes", "10"),
					resource.TestCheckResourceAttr(data.ResourceName, "zones.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMNatGatewayExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.NatGatewayClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Nat Gateway not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

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
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.NatGatewayClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

// Using alt location because the resource currently in private preview and is only available in eastus2.
func testAccAzureRMNatGateway_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-%d"
  location = "%s"
}

resource "azurerm_nat_gateway" "test" {
  name                = "acctestnatGateway-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Secondary, data.RandomInteger)
}

// Using alt location because the resource currently in private preview and is only available in eastus2.
func testAccAzureRMNatGateway_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicIP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  zones               = ["1"]
}

resource "azurerm_public_ip_prefix" "test" {
  name                = "acctestpublicIPPrefix-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  prefix_length       = 30
  zones               = ["1"]
}

resource "azurerm_nat_gateway" "test" {
  name                    = "acctestnatGateway-%d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  public_ip_address_ids   = [azurerm_public_ip.test.id]
  public_ip_prefix_ids    = [azurerm_public_ip_prefix.test.id]
  sku_name                = "Standard"
  idle_timeout_in_minutes = 10
  zones                   = ["1"]
}
`, data.RandomInteger, data.Locations.Secondary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

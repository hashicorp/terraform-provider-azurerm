package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMVirtualNetworkPeering_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_peering", "test1")
	secondResourceName := "azurerm_virtual_network_peering.test2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetworkPeering_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkPeeringExists(data.ResourceName),
					testCheckAzureRMVirtualNetworkPeeringExists(secondResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "allow_virtual_network_access", "true"),
					resource.TestCheckResourceAttr(secondResourceName, "allow_virtual_network_access", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVirtualNetworkPeering_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_peering", "test1")
	secondResourceName := "azurerm_virtual_network_peering.test2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetworkPeering_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkPeeringExists(data.ResourceName),
					testCheckAzureRMVirtualNetworkPeeringExists(secondResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMVirtualNetworkPeering_requiresImport),
		},
	})
}

func TestAccAzureRMVirtualNetworkPeering_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_peering", "test1")
	secondResourceName := "azurerm_virtual_network_peering.test2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetworkPeering_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkPeeringExists(data.ResourceName),
					testCheckAzureRMVirtualNetworkPeeringExists(secondResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "allow_virtual_network_access", "true"),
					resource.TestCheckResourceAttr(secondResourceName, "allow_virtual_network_access", "true"),
					testCheckAzureRMVirtualNetworkPeeringDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMVirtualNetworkPeering_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_network_peering", "test1")
	secondResourceName := "azurerm_virtual_network_peering.test2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualNetworkPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualNetworkPeering_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkPeeringExists(data.ResourceName),
					testCheckAzureRMVirtualNetworkPeeringExists(secondResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "allow_virtual_network_access", "true"),
					resource.TestCheckResourceAttr(secondResourceName, "allow_virtual_network_access", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "allow_forwarded_traffic", "false"),
					resource.TestCheckResourceAttr(secondResourceName, "allow_forwarded_traffic", "false"),
				),
			},

			{
				Config: testAccAzureRMVirtualNetworkPeering_basicUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualNetworkPeeringExists(data.ResourceName),
					testCheckAzureRMVirtualNetworkPeeringExists(secondResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "allow_virtual_network_access", "true"),
					resource.TestCheckResourceAttr(secondResourceName, "allow_virtual_network_access", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "allow_forwarded_traffic", "true"),
					resource.TestCheckResourceAttr(secondResourceName, "allow_forwarded_traffic", "true"),
				),
			},
		},
	})
}

func testCheckAzureRMVirtualNetworkPeeringExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.VnetPeeringsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		vnetName := rs.Primary.Attributes["virtual_network_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for virtual network peering: %s", name)
		}

		// Ensure resource group/virtual network peering combination exists in API
		resp, err := client.Get(ctx, resourceGroup, vnetName, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on vnetPeeringsClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Virtual Network Peering %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMVirtualNetworkPeeringDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.VnetPeeringsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		vnetName := rs.Primary.Attributes["virtual_network_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for virtual network peering: %s", name)
		}

		// Ensure resource group/virtual network peering combination exists in API
		future, err := client.Delete(ctx, resourceGroup, vnetName, name)
		if err != nil {
			return fmt.Errorf("Error deleting Peering %q (NW %q / RG %q): %+v", name, vnetName, resourceGroup, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for deletion of Peering %q (NW %q / RG %q): %+v", name, vnetName, resourceGroup, err)
		}

		return nil
	}
}

func testCheckAzureRMVirtualNetworkPeeringDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.VnetPeeringsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_virtual_network_peering" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		vnetName := rs.Primary.Attributes["virtual_network_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, vnetName, name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Virtual Network Peering sitll exists:\n%#v", resp.VirtualNetworkPeeringPropertiesFormat)
		}
	}

	return nil
}

func testAccAzureRMVirtualNetworkPeering_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test1" {
  name                = "acctestvirtnet-1-%d"
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.1.0/24"]
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_network" "test2" {
  name                = "acctestvirtnet-2-%d"
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.2.0/24"]
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_network_peering" "test1" {
  name                         = "acctestpeer-1-%d"
  resource_group_name          = azurerm_resource_group.test.name
  virtual_network_name         = azurerm_virtual_network.test1.name
  remote_virtual_network_id    = azurerm_virtual_network.test2.id
  allow_virtual_network_access = true
}

resource "azurerm_virtual_network_peering" "test2" {
  name                         = "acctestpeer-2-%d"
  resource_group_name          = azurerm_resource_group.test.name
  virtual_network_name         = azurerm_virtual_network.test2.name
  remote_virtual_network_id    = azurerm_virtual_network.test1.id
  allow_virtual_network_access = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMVirtualNetworkPeering_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMVirtualNetworkPeering_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network_peering" "import" {
  name                         = azurerm_virtual_network_peering.test1.name
  resource_group_name          = azurerm_virtual_network_peering.test1.resource_group_name
  virtual_network_name         = azurerm_virtual_network_peering.test1.virtual_network_name
  remote_virtual_network_id    = azurerm_virtual_network_peering.test1.remote_virtual_network_id
  allow_virtual_network_access = azurerm_virtual_network_peering.test1.allow_virtual_network_access
}
`, template)
}

func testAccAzureRMVirtualNetworkPeering_basicUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test1" {
  name                = "acctestvirtnet-1-%d"
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.1.0/24"]
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_network" "test2" {
  name                = "acctestvirtnet-2-%d"
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.2.0/24"]
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_network_peering" "test1" {
  name                         = "acctestpeer-1-%d"
  resource_group_name          = azurerm_resource_group.test.name
  virtual_network_name         = azurerm_virtual_network.test1.name
  remote_virtual_network_id    = azurerm_virtual_network.test2.id
  allow_forwarded_traffic      = true
  allow_virtual_network_access = true
}

resource "azurerm_virtual_network_peering" "test2" {
  name                         = "acctestpeer-2-%d"
  resource_group_name          = azurerm_resource_group.test.name
  virtual_network_name         = azurerm_virtual_network.test2.name
  remote_virtual_network_id    = azurerm_virtual_network.test1.id
  allow_forwarded_traffic      = true
  allow_virtual_network_access = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMPrivateDnsZoneVirtualNetworkLink_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_zone_virtual_network_link", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsZoneVirtualNetworkLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateDnsZoneVirtualNetworkLink_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsZoneVirtualNetworkLinkExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPrivateDnsZoneVirtualNetworkLink_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_zone_virtual_network_link", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsZoneVirtualNetworkLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateDnsZoneVirtualNetworkLink_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsZoneVirtualNetworkLinkExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMPrivateDnsZoneVirtualNetworkLink_requiresImport),
		},
	})
}

func TestAccAzureRMPrivateDnsZoneVirtualNetworkLink_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_zone_virtual_network_link", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsZoneVirtualNetworkLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateDnsZoneVirtualNetworkLink_withTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsZoneVirtualNetworkLinkExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
				),
			},
			{
				Config: testAccAzureRMPrivateDnsZoneVirtualNetworkLink_withTagsUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsZoneVirtualNetworkLinkExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMPrivateDnsZoneVirtualNetworkLinkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).PrivateDns.VirtualNetworkLinksClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		dnsZoneName := rs.Primary.Attributes["private_dns_zone_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]

		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Private DNS zone virtual network link: %s", name)
		}

		resp, err := client.Get(ctx, resourceGroup, dnsZoneName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: virtual network link %q (Private DNS zone %q / resource group: %s) does not exist", name, dnsZoneName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get Private DNS zone virtual network link: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMPrivateDnsZoneVirtualNetworkLinkDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).PrivateDns.VirtualNetworkLinksClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_private_dns_zone_virtual_network_link" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		dnsZoneName := rs.Primary.Attributes["private_dns_zone_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, dnsZoneName, name)
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}

			return err
		}

		return fmt.Errorf("Private DNS zone virtual network link still exists:\n%#v", resp)
	}

	return nil
}

func testAccAzureRMPrivateDnsZoneVirtualNetworkLink_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "vnet%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]

  subnet {
    name           = "subnet1"
    address_prefix = "10.0.1.0/24"
  }
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_dns_zone_virtual_network_link" "test" {
  name                  = "acctestVnetZone%d.com"
  private_dns_zone_name = azurerm_private_dns_zone.test.name
  virtual_network_id    = azurerm_virtual_network.test.id
  resource_group_name   = azurerm_resource_group.test.name
}

`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMPrivateDnsZoneVirtualNetworkLink_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMPrivateDnsZoneVirtualNetworkLink_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_private_dns_zone_virtual_network_link" "import" {
  name                  = azurerm_private_dns_zone_virtual_network_link.test.name
  private_dns_zone_name = azurerm_private_dns_zone_virtual_network_link.test.private_dns_zone_name
  virtual_network_id    = azurerm_private_dns_zone_virtual_network_link.test.virtual_network_id
  resource_group_name   = azurerm_private_dns_zone_virtual_network_link.test.resource_group_name
}
`, template)
}

func testAccAzureRMPrivateDnsZoneVirtualNetworkLink_withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "vnet%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]

  subnet {
    name           = "subnet1"
    address_prefix = "10.0.1.0/24"
  }
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_private_dns_zone_virtual_network_link" "test" {
  name                  = "acctestVnetZone%d.com"
  private_dns_zone_name = azurerm_private_dns_zone.test.name
  virtual_network_id    = azurerm_virtual_network.test.id
  resource_group_name   = azurerm_resource_group.test.name

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMPrivateDnsZoneVirtualNetworkLink_withTagsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "vnet%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]

  subnet {
    name           = "subnet1"
    address_prefix = "10.0.1.0/24"
  }
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_private_dns_zone_virtual_network_link" "test" {
  name                  = "acctestVnetZone%d.com"
  private_dns_zone_name = azurerm_private_dns_zone.test.name
  virtual_network_id    = azurerm_virtual_network.test.id
  resource_group_name   = azurerm_resource_group.test.name

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

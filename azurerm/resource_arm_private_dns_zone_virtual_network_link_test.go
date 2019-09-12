package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMPrivateDnsZoneVirtualNetworkLink_basic(t *testing.T) {
	resourceName := "azurerm_private_dns_zone_virtual_network_link.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMPrivateDnsZoneVirtualNetworkLink_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsZoneVirtualNetworkLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsZoneVirtualNetworkLinkExists(resourceName),
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

func TestAccAzureRMPrivateDnsZoneVirtualNetworkLink_requiresImport(t *testing.T) {
	if !requireResourcesToBeImported {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_private_dns_zone_virtual_network_link.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsZoneVirtualNetworkLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateDnsZoneVirtualNetworkLink_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsZoneVirtualNetworkLinkExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMPrivateDnsZoneVirtualNetworkLink_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_private_dns_zone_virtual_network_link"),
			},
		},
	})
}

func TestAccAzureRMPrivateDnsZoneVirtualNetworkLink_withTags(t *testing.T) {
	resourceName := "azurerm_private_dns_zone_virtual_network_link.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	preConfig := testAccAzureRMPrivateDnsZoneVirtualNetworkLink_withTags(ri, location)
	postConfig := testAccAzureRMPrivateDnsZoneVirtualNetworkLink_withTagsUpdate(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsZoneVirtualNetworkLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsZoneVirtualNetworkLinkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsZoneVirtualNetworkLinkExists(resourceName),
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

func testCheckAzureRMPrivateDnsZoneVirtualNetworkLinkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
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

		client := testAccProvider.Meta().(*ArmClient).privateDns.VirtualNetworkLinksClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

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
	conn := testAccProvider.Meta().(*ArmClient).privateDns.VirtualNetworkLinksClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

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

func testAccAzureRMPrivateDnsZoneVirtualNetworkLink_basic(rInt int, location string) string {
	return fmt.Sprintf(`
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
  name                	= "acctest%d"
  private_dns_zone_name = "${azurerm_private_dns_zone.test.name}"
  virtual_network_id 	= "${azurerm_virtual_network.test.id}"
  resource_group_name 	= "${azurerm_resource_group.test.name}"
}

`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMPrivateDnsZoneVirtualNetworkLink_requiresImport(rInt int, location string) string {
	template := testAccAzureRMPrivateDnsZoneVirtualNetworkLink_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_private_dns_zone_virtual_network_link" "import" {
  name                	= "${azurerm_private_dns_zone_virtual_network_link.test.name}
  private_dns_zone_name = "${azurerm_private_dns_zone_virtual_network_link.test.private_dns_zone_name}"
  virtual_network_id 	= "${azurerm_private_dns_zone_virtual_network_link.test.virtual_network_id}"
  resource_group_name 	= "${azurerm_private_dns_zone_virtual_network_link.test.resource_group_name}"
}
`, template)
}

func testAccAzureRMPrivateDnsZoneVirtualNetworkLink_withTags(rInt int, location string) string {
	return fmt.Sprintf(`
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
  name                	= "acctest%d"
  private_dns_zone_name = "${azurerm_private_dns_zone.test.name}"
  virtual_network_id 	= "${azurerm_virtual_network.test.id}"
  resource_group_name 	= "${azurerm_resource_group.test.name}"
	
  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMPrivateDnsZoneVirtualNetworkLink_withTagsUpdate(rInt int, location string) string {
	return fmt.Sprintf(`
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
  name                	= "acctestzone%d.com"
  private_dns_zone_name = "${azurerm_private_dns_zone.test.name}"
  virtual_network_id 	= "${azurerm_virtual_network.test.id}"
  resource_group_name 	= "${azurerm_resource_group.test.name}"
	
  tags = {
    environment = "staging"
  }
}
`, rInt, location, rInt, rInt, rInt)
}

package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMDnsZone_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_zone", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDnsZone_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsZoneExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDnsZone_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_dns_zone", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDnsZone_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsZoneExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMDnsZone_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_dns_zone"),
			},
		},
	})
}

func TestAccAzureRMDnsZone_withVNets(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_zone", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDnsZone_withVNets(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsZoneExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDnsZone_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_zone", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDnsZone_withTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsZoneExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
				),
			},
			{
				Config: testAccAzureRMDnsZone_withTagsUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsZoneExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMDnsZoneExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		zoneName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for DNS zone: %s", zoneName)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Dns.ZonesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		resp, err := client.Get(ctx, resourceGroup, zoneName)
		if err != nil {
			return fmt.Errorf("Bad: Get DNS zone: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: DNS zone %s (resource group: %s) does not exist", zoneName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDnsZoneDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Dns.ZonesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_dns_zone" {
			continue
		}

		zoneName := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, zoneName)
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}

			return err
		}

		return fmt.Errorf("DNS Zone still exists:\n%#v", resp)
	}

	return nil
}

func testAccAzureRMDnsZone_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMDnsZone_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMDnsZone_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_dns_zone" "import" {
  name                = "${azurerm_dns_zone.test.name}"
  resource_group_name = "${azurerm_dns_zone.test.resource_group_name}"
}
`, template)
}

func testAccAzureRMDnsZone_withVNets(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG_%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet%d"
  location            = "%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.0.0.0/16"]
  dns_servers         = ["168.63.129.16"]
}

resource "azurerm_dns_zone" "test" {
  name                             = "acctestzone%d.com"
  resource_group_name              = "${azurerm_resource_group.test.name}"
  zone_type                        = "Private"
  registration_virtual_network_ids = ["${azurerm_virtual_network.test.id}"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMDnsZone_withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMDnsZone_withTagsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

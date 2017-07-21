package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMDnsZone_basic(t *testing.T) {
	ri := acctest.RandInt()
	config := fmt.Sprintf(testAccAzureRMDnsZone_basic, ri, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDnsZoneDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsZoneExists("azurerm_dns_zone.test"),
				),
			},
		},
	})
}

func TestAccAzureRMDnsZone_withTags(t *testing.T) {
	ri := acctest.RandInt()
	preConfig := fmt.Sprintf(testAccAzureRMDnsZone_withTags, ri, ri)
	postConfig := fmt.Sprintf(testAccAzureRMDnsZone_withTagsUupdate, ri, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDnsZoneDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsZoneExists("azurerm_dns_zone.test"),
					resource.TestCheckResourceAttr(
						"azurerm_dns_zone.test", "tags.%", "2"),
				),
			},

			resource.TestStep{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsZoneExists("azurerm_dns_zone.test"),
					resource.TestCheckResourceAttr(
						"azurerm_dns_zone.test", "tags.%", "1"),
				),
			},
		},
	})
}

func testCheckAzureRMDnsZoneExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		zoneName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for DNS zone: %s", zoneName)
		}

		conn := testAccProvider.Meta().(*ArmClient).zonesClient
		resp, err := conn.Get(resourceGroup, zoneName)
		if err != nil {
			return fmt.Errorf("Bad: Get zone: %v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: DNS zone %s (resource group: %s) does not exist", zoneName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDnsZoneDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).zonesClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_dns_zone_record" {
			continue
		}

		zoneName := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(resourceGroup, zoneName)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("DNS zone still exists:\n%#v", resp.ZoneProperties)
		}

	}

	return nil
}

var testAccAzureRMDnsZone_basic = `
resource "azurerm_resource_group" "test" {
    name = "acctestRG_%d"
    location = "West US"
}
resource "azurerm_dns_zone" "test" {
    name = "acctestzone%d.com"
    resource_group_name = "${azurerm_resource_group.test.name}"
}
`

var testAccAzureRMDnsZone_withTags = `
resource "azurerm_resource_group" "test" {
    name = "acctestRG_%d"
    location = "West US"
}
resource "azurerm_dns_zone" "test" {
    name = "acctestzone%d.com"
    resource_group_name = "${azurerm_resource_group.test.name}"
	tags {
		environment = "Production"
		cost_center = "MSFT"
    }
}
`

var testAccAzureRMDnsZone_withTagsUupdate = `
resource "azurerm_resource_group" "test" {
    name = "acctestRG_%d"
    location = "West US"
}
resource "azurerm_dns_zone" "test" {
    name = "acctestzone%d.com"
    resource_group_name = "${azurerm_resource_group.test.name}"
	tags {
		environment = "staging"
    }
}
`

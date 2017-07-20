package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/arm/dns"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMDnsARecord_basic(t *testing.T) {
	ri := acctest.RandInt()
	config := fmt.Sprintf(testAccAzureRMDnsARecord_basic, ri, ri, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDnsARecordDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsARecordExists("azurerm_dns_a_record.test"),
				),
			},
		},
	})
}

func TestAccAzureRMDnsARecord_updateRecords(t *testing.T) {
	ri := acctest.RandInt()
	preConfig := fmt.Sprintf(testAccAzureRMDnsARecord_basic, ri, ri, ri)
	postConfig := fmt.Sprintf(testAccAzureRMDnsARecord_updateRecords, ri, ri, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDnsARecordDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsARecordExists("azurerm_dns_a_record.test"),
					resource.TestCheckResourceAttr(
						"azurerm_dns_a_record.test", "records.#", "2"),
				),
			},

			resource.TestStep{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsARecordExists("azurerm_dns_a_record.test"),
					resource.TestCheckResourceAttr(
						"azurerm_dns_a_record.test", "records.#", "3"),
				),
			},
		},
	})
}

func TestAccAzureRMDnsARecord_withTags(t *testing.T) {
	ri := acctest.RandInt()
	preConfig := fmt.Sprintf(testAccAzureRMDnsARecord_withTags, ri, ri, ri)
	postConfig := fmt.Sprintf(testAccAzureRMDnsARecord_withTagsUpdate, ri, ri, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDnsARecordDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsARecordExists("azurerm_dns_a_record.test"),
					resource.TestCheckResourceAttr(
						"azurerm_dns_a_record.test", "tags.%", "2"),
				),
			},

			resource.TestStep{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsARecordExists("azurerm_dns_a_record.test"),
					resource.TestCheckResourceAttr(
						"azurerm_dns_a_record.test", "tags.%", "1"),
				),
			},
		},
	})
}

func testCheckAzureRMDnsARecordExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		aName := rs.Primary.Attributes["name"]
		zoneName := rs.Primary.Attributes["zone_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for DNS A record: %s", aName)
		}

		conn := testAccProvider.Meta().(*ArmClient).dnsClient
		resp, err := conn.Get(resourceGroup, zoneName, aName, dns.A)
		if err != nil {
			return fmt.Errorf("Bad: Get A RecordSet: %v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: DNS A record %s (resource group: %s) does not exist", aName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDnsARecordDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).dnsClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_dns_a_record" {
			continue
		}

		aName := rs.Primary.Attributes["name"]
		zoneName := rs.Primary.Attributes["zone_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(resourceGroup, zoneName, aName, dns.A)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("DNS A record still exists:\n%#v", resp.RecordSetProperties)
		}

	}

	return nil
}

var testAccAzureRMDnsARecord_basic = `
resource "azurerm_resource_group" "test" {
    name = "acctestRG_%d"
    location = "West US"
}
resource "azurerm_dns_zone" "test" {
    name = "acctestzone%d.com"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_a_record" "test" {
    name = "myarecord%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    zone_name = "${azurerm_dns_zone.test.name}"
    ttl = "300"
    records = ["1.2.3.4", "1.2.4.5"]
}
`

var testAccAzureRMDnsARecord_updateRecords = `
resource "azurerm_resource_group" "test" {
    name = "acctestRG_%d"
    location = "West US"
}
resource "azurerm_dns_zone" "test" {
    name = "acctestzone%d.com"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_a_record" "test" {
    name = "myarecord%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    zone_name = "${azurerm_dns_zone.test.name}"
    ttl = "300"
    records = ["1.2.3.4", "1.2.4.5", "1.2.3.7"]
}
`

var testAccAzureRMDnsARecord_withTags = `
resource "azurerm_resource_group" "test" {
    name = "acctestRG_%d"
    location = "West US"
}
resource "azurerm_dns_zone" "test" {
    name = "acctestzone%d.com"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_a_record" "test" {
    name = "myarecord%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    zone_name = "${azurerm_dns_zone.test.name}"
    ttl = "300"
    records = ["1.2.3.4", "1.2.4.5"]

    tags {
	environment = "Production"
	cost_center = "MSFT"
    }
}
`

var testAccAzureRMDnsARecord_withTagsUpdate = `
resource "azurerm_resource_group" "test" {
    name = "acctestRG_%d"
    location = "West US"
}
resource "azurerm_dns_zone" "test" {
    name = "acctestzone%d.com"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_a_record" "test" {
    name = "myarecord%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    zone_name = "${azurerm_dns_zone.test.name}"
    ttl = "300"
    records = ["1.2.3.4", "1.2.4.5"]

    tags {
	environment = "staging"
    }
}
`

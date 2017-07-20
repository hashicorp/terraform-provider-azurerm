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

func TestAccAzureRMDnsNsRecord_basic(t *testing.T) {
	ri := acctest.RandInt()
	config := fmt.Sprintf(testAccAzureRMDnsNsRecord_basic, ri, ri, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDnsNsRecordDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsNsRecordExists("azurerm_dns_ns_record.test"),
				),
			},
		},
	})
}

func TestAccAzureRMDnsNsRecord_updateRecords(t *testing.T) {
	ri := acctest.RandInt()
	preConfig := fmt.Sprintf(testAccAzureRMDnsNsRecord_basic, ri, ri, ri)
	postConfig := fmt.Sprintf(testAccAzureRMDnsNsRecord_updateRecords, ri, ri, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDnsNsRecordDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsNsRecordExists("azurerm_dns_ns_record.test"),
					resource.TestCheckResourceAttr(
						"azurerm_dns_ns_record.test", "record.#", "2"),
				),
			},

			resource.TestStep{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsNsRecordExists("azurerm_dns_ns_record.test"),
					resource.TestCheckResourceAttr(
						"azurerm_dns_ns_record.test", "record.#", "3"),
				),
			},
		},
	})
}

func TestAccAzureRMDnsNsRecord_withTags(t *testing.T) {
	ri := acctest.RandInt()
	preConfig := fmt.Sprintf(testAccAzureRMDnsNsRecord_withTags, ri, ri, ri)
	postConfig := fmt.Sprintf(testAccAzureRMDnsNsRecord_withTagsUpdate, ri, ri, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDnsNsRecordDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsNsRecordExists("azurerm_dns_ns_record.test"),
					resource.TestCheckResourceAttr(
						"azurerm_dns_ns_record.test", "tags.%", "2"),
				),
			},

			resource.TestStep{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsNsRecordExists("azurerm_dns_ns_record.test"),
					resource.TestCheckResourceAttr(
						"azurerm_dns_ns_record.test", "tags.%", "1"),
				),
			},
		},
	})
}

func testCheckAzureRMDnsNsRecordExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		nsName := rs.Primary.Attributes["name"]
		zoneName := rs.Primary.Attributes["zone_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for DNS NS record: %s", nsName)
		}

		conn := testAccProvider.Meta().(*ArmClient).dnsClient
		resp, err := conn.Get(resourceGroup, zoneName, nsName, dns.NS)
		if err != nil {
			return fmt.Errorf("Bad: Get NS recordSet: %v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: DNS NS record %s (resource group: %s) does not exist", nsName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDnsNsRecordDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).dnsClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_dns_ns_record" {
			continue
		}

		nsName := rs.Primary.Attributes["name"]
		zoneName := rs.Primary.Attributes["zone_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(resourceGroup, zoneName, nsName, dns.NS)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("DNS MSNSecord still exists:\n%#v", resp.RecordSetProperties)
		}

	}

	return nil
}

var testAccAzureRMDnsNsRecord_basic = `
resource "azurerm_resource_group" "test" {
    name = "acctestRG_%d"
    location = "West US"
}
resource "azurerm_dns_zone" "test" {
    name = "acctestzone%d.com"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_ns_record" "test" {
    name = "mynsrecord%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    zone_name = "${azurerm_dns_zone.test.name}"
    ttl = "300"

    record {
    	nsdname = "ns1.contoso.com"
    }

    record {
    	nsdname = "ns2.contoso.com"
    }
}
`

var testAccAzureRMDnsNsRecord_updateRecords = `
resource "azurerm_resource_group" "test" {
    name = "acctestRG_%d"
    location = "West US"
}
resource "azurerm_dns_zone" "test" {
    name = "acctestzone%d.com"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_ns_record" "test" {
    name = "mynsrecord%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    zone_name = "${azurerm_dns_zone.test.name}"
    ttl = "300"

    record {
    	nsdname = "ns1.contoso.com"
    }

    record {
    	nsdname = "ns2.contoso.com"
    }

    record {
    	nsdname = "ns3.contoso.com"
    }
}
`

var testAccAzureRMDnsNsRecord_withTags = `
resource "azurerm_resource_group" "test" {
    name = "acctestRG_%d"
    location = "West US"
}
resource "azurerm_dns_zone" "test" {
    name = "acctestzone%d.com"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_ns_record" "test" {
    name = "mynsrecord%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    zone_name = "${azurerm_dns_zone.test.name}"
    ttl = "300"

    record {
    	nsdname = "ns1.contoso.com"
    }

    record {
    	nsdname = "ns2.contoso.com"
    }

    tags {
	environment = "Production"
	cost_center = "MSFT"
    }
}
`

var testAccAzureRMDnsNsRecord_withTagsUpdate = `
resource "azurerm_resource_group" "test" {
    name = "acctestRG_%d"
    location = "West US"
}
resource "azurerm_dns_zone" "test" {
    name = "acctestzone%d.com"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_ns_record" "test" {
    name = "mynsrecord%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    zone_name = "${azurerm_dns_zone.test.name}"
    ttl = "300"
    record {
    	nsdname = "ns1.contoso.com"
    }

    record {
    	nsdname = "ns2.contoso.com"
    }

    tags {
	environment = "staging"
    }
}
`

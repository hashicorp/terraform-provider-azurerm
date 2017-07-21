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

func TestAccAzureRMDnsSrvRecord_basic(t *testing.T) {
	ri := acctest.RandInt()
	config := fmt.Sprintf(testAccAzureRMDnsSrvRecord_basic, ri, ri, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDnsSrvRecordDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsSrvRecordExists("azurerm_dns_srv_record.test"),
				),
			},
		},
	})
}

func TestAccAzureRMDnsSrvRecord_updateRecords(t *testing.T) {
	ri := acctest.RandInt()
	preConfig := fmt.Sprintf(testAccAzureRMDnsSrvRecord_basic, ri, ri, ri)
	postConfig := fmt.Sprintf(testAccAzureRMDnsSrvRecord_updateRecords, ri, ri, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDnsSrvRecordDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsSrvRecordExists("azurerm_dns_srv_record.test"),
					resource.TestCheckResourceAttr(
						"azurerm_dns_srv_record.test", "record.#", "2"),
				),
			},

			resource.TestStep{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsSrvRecordExists("azurerm_dns_srv_record.test"),
					resource.TestCheckResourceAttr(
						"azurerm_dns_srv_record.test", "record.#", "3"),
				),
			},
		},
	})
}

func TestAccAzureRMDnsSrvRecord_withTags(t *testing.T) {
	ri := acctest.RandInt()
	preConfig := fmt.Sprintf(testAccAzureRMDnsSrvRecord_withTags, ri, ri, ri)
	postConfig := fmt.Sprintf(testAccAzureRMDnsSrvRecord_withTagsUpdate, ri, ri, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDnsSrvRecordDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsSrvRecordExists("azurerm_dns_srv_record.test"),
					resource.TestCheckResourceAttr(
						"azurerm_dns_srv_record.test", "tags.%", "2"),
				),
			},

			resource.TestStep{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsSrvRecordExists("azurerm_dns_srv_record.test"),
					resource.TestCheckResourceAttr(
						"azurerm_dns_srv_record.test", "tags.%", "1"),
				),
			},
		},
	})
}

func testCheckAzureRMDnsSrvRecordExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		srvName := rs.Primary.Attributes["name"]
		zoneName := rs.Primary.Attributes["zone_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for DNS SRV record: %s", srvName)
		}

		conn := testAccProvider.Meta().(*ArmClient).dnsClient
		resp, err := conn.Get(resourceGroup, zoneName, srvName, dns.SRV)
		if err != nil {
			return fmt.Errorf("Bad: Get SRV RecordSet: %v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: DNS SRV record %s (resource group: %s) does not exist", srvName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDnsSrvRecordDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).dnsClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_dns_srv_record" {
			continue
		}

		srvName := rs.Primary.Attributes["name"]
		zoneName := rs.Primary.Attributes["zone_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(resourceGroup, zoneName, srvName, dns.SRV)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("DNS SRV record still exists:\n%#v", resp.RecordSetProperties)
		}

	}

	return nil
}

var testAccAzureRMDnsSrvRecord_basic = `
resource "azurerm_resource_group" "test" {
    name = "acctestRG_%d"
    location = "West US"
}
resource "azurerm_dns_zone" "test" {
    name = "acctestzone%d.com"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_srv_record" "test" {
    name = "myarecord%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    zone_name = "${azurerm_dns_zone.test.name}"
    ttl = "300"

    record {
	priority = 1
	weight = 5
	port = 8080
	target = "target1.contoso.com"
    }

    record {
	priority = 2
	weight = 25
	port = 8080
	target = "target2.contoso.com"
    }
}
`

var testAccAzureRMDnsSrvRecord_updateRecords = `
resource "azurerm_resource_group" "test" {
    name = "acctestRG_%d"
    location = "West US"
}
resource "azurerm_dns_zone" "test" {
    name = "acctestzone%d.com"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_srv_record" "test" {
    name = "myarecord%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    zone_name = "${azurerm_dns_zone.test.name}"
    ttl = "300"

    record {
	priority = 1
	weight = 5
	port = 8080
	target = "target1.contoso.com"
    }

    record {
	priority = 2
	weight = 25
	port = 8080
	target = "target2.contoso.com"
    }

    record {
	priority = 3
	weight = 100
	port = 8080
	target = "target3.contoso.com"
    }
}
`

var testAccAzureRMDnsSrvRecord_withTags = `
resource "azurerm_resource_group" "test" {
    name = "acctestRG_%d"
    location = "West US"
}
resource "azurerm_dns_zone" "test" {
    name = "acctestzone%d.com"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_srv_record" "test" {
    name = "myarecord%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    zone_name = "${azurerm_dns_zone.test.name}"
    ttl = "300"

    record {
	priority = 1
	weight = 5
	port = 8080
	target = "target1.contoso.com"
    }

    record {
	priority = 2
	weight = 25
	port = 8080
	target = "target2.contoso.com"
    }

    tags {
	environment = "Production"
	cost_center = "MSFT"
    }
}
`

var testAccAzureRMDnsSrvRecord_withTagsUpdate = `
resource "azurerm_resource_group" "test" {
    name = "acctestRG_%d"
    location = "West US"
}
resource "azurerm_dns_zone" "test" {
    name = "acctestzone%d.com"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_srv_record" "test" {
    name = "myarecord%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    zone_name = "${azurerm_dns_zone.test.name}"
    ttl = "300"

    record {
	priority = 1
	weight = 5
	port = 8080
	target = "target1.contoso.com"
    }

    record {
	priority = 2
	weight = 25
	port = 8080
	target = "target2.contoso.com"
    }

    tags {
	environment = "staging"
    }
}
`

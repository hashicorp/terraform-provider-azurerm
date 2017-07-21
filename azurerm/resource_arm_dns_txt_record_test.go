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

func TestAccAzureRMDnsTxtRecord_basic(t *testing.T) {
	ri := acctest.RandInt()
	config := fmt.Sprintf(testAccAzureRMDnsTxtRecord_basic, ri, ri, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDnsTxtRecordDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsTxtRecordExists("azurerm_dns_txt_record.test"),
				),
			},
		},
	})
}

func TestAccAzureRMDnsTxtRecord_updateRecords(t *testing.T) {
	ri := acctest.RandInt()
	preConfig := fmt.Sprintf(testAccAzureRMDnsTxtRecord_basic, ri, ri, ri)
	postConfig := fmt.Sprintf(testAccAzureRMDnsTxtRecord_updateRecords, ri, ri, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDnsTxtRecordDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsTxtRecordExists("azurerm_dns_txt_record.test"),
					resource.TestCheckResourceAttr(
						"azurerm_dns_txt_record.test", "record.#", "2"),
				),
			},

			resource.TestStep{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsTxtRecordExists("azurerm_dns_txt_record.test"),
					resource.TestCheckResourceAttr(
						"azurerm_dns_txt_record.test", "record.#", "3"),
				),
			},
		},
	})
}

func TestAccAzureRMDnsTxtRecord_withTags(t *testing.T) {
	ri := acctest.RandInt()
	preConfig := fmt.Sprintf(testAccAzureRMDnsTxtRecord_withTags, ri, ri, ri)
	postConfig := fmt.Sprintf(testAccAzureRMDnsTxtRecord_withTagsUpdate, ri, ri, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDnsTxtRecordDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsTxtRecordExists("azurerm_dns_txt_record.test"),
					resource.TestCheckResourceAttr(
						"azurerm_dns_txt_record.test", "tags.%", "2"),
				),
			},

			resource.TestStep{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsTxtRecordExists("azurerm_dns_txt_record.test"),
					resource.TestCheckResourceAttr(
						"azurerm_dns_txt_record.test", "tags.%", "1"),
				),
			},
		},
	})
}

func testCheckAzureRMDnsTxtRecordExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		txtName := rs.Primary.Attributes["name"]
		zoneName := rs.Primary.Attributes["zone_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for DNS TXT record: %s", txtName)
		}

		conn := testAccProvider.Meta().(*ArmClient).dnsClient
		resp, err := conn.Get(resourceGroup, zoneName, txtName, dns.TXT)
		if err != nil {
			return fmt.Errorf("Bad: Get TXT RecordSet: %v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: DNS TXT record %s (resource group: %s) does not exist", txtName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDnsTxtRecordDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).dnsClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_dns_txt_record" {
			continue
		}

		txtName := rs.Primary.Attributes["name"]
		zoneName := rs.Primary.Attributes["zone_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(resourceGroup, zoneName, txtName, dns.TXT)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("DNS TXT record still exists:\n%#v", resp.RecordSetProperties)
		}

	}

	return nil
}

var testAccAzureRMDnsTxtRecord_basic = `
resource "azurerm_resource_group" "test" {
    name = "acctestRG_%d"
    location = "West US"
}
resource "azurerm_dns_zone" "test" {
    name = "acctestzone%d.com"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_txt_record" "test" {
    name = "myarecord%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    zone_name = "${azurerm_dns_zone.test.name}"
    ttl = "300"

    record {
    	value = "Quick brown fox"
    }

    record {
    	value = "Another test txt string"
    }
}
`

var testAccAzureRMDnsTxtRecord_updateRecords = `
resource "azurerm_resource_group" "test" {
    name = "acctestRG_%d"
    location = "West US"
}
resource "azurerm_dns_zone" "test" {
    name = "acctestzone%d.com"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_txt_record" "test" {
    name = "myarecord%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    zone_name = "${azurerm_dns_zone.test.name}"
    ttl = "300"

    record {
    	value = "Quick brown fox"
    }

    record {
    	value = "Another test txt string"
    }

    record {
    	value = "A wild 3rd record appears"
    }
}
`

var testAccAzureRMDnsTxtRecord_withTags = `
resource "azurerm_resource_group" "test" {
    name = "acctestRG_%d"
    location = "West US"
}
resource "azurerm_dns_zone" "test" {
    name = "acctestzone%d.com"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_txt_record" "test" {
    name = "myarecord%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    zone_name = "${azurerm_dns_zone.test.name}"
    ttl = "300"

    record {
    	value = "Quick brown fox"
    }

    record {
    	value = "Another test txt string"
    }

    tags {
	environment = "Production"
	cost_center = "MSFT"
    }
}
`

var testAccAzureRMDnsTxtRecord_withTagsUpdate = `
resource "azurerm_resource_group" "test" {
    name = "acctestRG_%d"
    location = "West US"
}
resource "azurerm_dns_zone" "test" {
    name = "acctestzone%d.com"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_txt_record" "test" {
    name = "myarecord%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    zone_name = "${azurerm_dns_zone.test.name}"
    ttl = "300"
    record {
    	value = "Quick brown fox"
    }

    record {
    	value = "Another test txt string"
    }

    tags {
	environment = "staging"
    }
}
`

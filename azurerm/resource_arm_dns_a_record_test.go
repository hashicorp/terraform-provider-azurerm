package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/dns/mgmt/2016-04-01/dns"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMDnsARecord_basic(t *testing.T) {
	resourceName := "azurerm_dns_a_record.test"
	ri := acctest.RandInt()
	config := testAccAzureRMDnsARecord_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDnsARecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsARecordExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMDnsARecord_updateRecords(t *testing.T) {
	resourceName := "azurerm_dns_a_record.test"
	ri := acctest.RandInt()
	location := testLocation()
	preConfig := testAccAzureRMDnsARecord_basic(ri, location)
	postConfig := testAccAzureRMDnsARecord_updateRecords(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDnsARecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsARecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "records.#", "2"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsARecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "records.#", "3"),
				),
			},
		},
	})
}

func TestAccAzureRMDnsARecord_withTags(t *testing.T) {
	resourceName := "azurerm_dns_a_record.test"
	ri := acctest.RandInt()
	location := testLocation()
	preConfig := testAccAzureRMDnsARecord_withTags(ri, location)
	postConfig := testAccAzureRMDnsARecord_withTagsUpdate(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDnsARecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsARecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsARecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
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
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := conn.Get(ctx, resourceGroup, zoneName, aName, dns.A)
		if err != nil {
			return fmt.Errorf("Bad: Get A RecordSet: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: DNS A record %s (resource group: %s) does not exist", aName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDnsARecordDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).dnsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_dns_a_record" {
			continue
		}

		aName := rs.Primary.Attributes["name"]
		zoneName := rs.Primary.Attributes["zone_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, zoneName, aName, dns.A)

		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}

			return err
		}

		return fmt.Errorf("DNS A record still exists:\n%#v", resp.RecordSetProperties)
	}

	return nil
}

func testAccAzureRMDnsARecord_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_a_record" "test" {
  name                = "myarecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300
  records             = ["1.2.3.4", "1.2.4.5"]
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMDnsARecord_updateRecords(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_a_record" "test" {
  name                = "myarecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300
  records             = ["1.2.3.4", "1.2.4.5", "1.2.3.7"]
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMDnsARecord_withTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_a_record" "test" {
  name                = "myarecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300
  records             = ["1.2.3.4", "1.2.4.5"]

  tags {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMDnsARecord_withTagsUpdate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_a_record" "test" {
  name                = "myarecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300
  records             = ["1.2.3.4", "1.2.4.5"]

  tags {
    environment = "staging"
  }
}
`, rInt, location, rInt, rInt)
}

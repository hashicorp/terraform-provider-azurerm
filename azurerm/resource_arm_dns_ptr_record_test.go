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

func TestAccAzureRMDnsPtrRecord_basic(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMDnsPtrRecord_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDnsPtrRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsPtrRecordExists("azurerm_dns_ptr_record.test"),
				),
			},
		},
	})
}

func TestAccAzureRMDnsPtrRecord_updateRecords(t *testing.T) {
	ri := acctest.RandInt()
	location := testLocation()
	preConfig := testAccAzureRMDnsPtrRecord_basic(ri, location)
	postConfig := testAccAzureRMDnsPtrRecord_updateRecords(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDnsPtrRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsPtrRecordExists("azurerm_dns_ptr_record.test"),
					resource.TestCheckResourceAttr("azurerm_dns_ptr_record.test", "records.#", "2"),
				),
			},

			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsPtrRecordExists("azurerm_dns_ptr_record.test"),
					resource.TestCheckResourceAttr("azurerm_dns_ptr_record.test", "records.#", "3"),
				),
			},
		},
	})
}

func TestAccAzureRMDnsPtrRecord_withTags(t *testing.T) {
	ri := acctest.RandInt()
	location := testLocation()
	preConfig := testAccAzureRMDnsPtrRecord_withTags(ri, location)
	postConfig := testAccAzureRMDnsPtrRecord_withTagsUpdate(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDnsPtrRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsPtrRecordExists("azurerm_dns_ptr_record.test"),
					resource.TestCheckResourceAttr("azurerm_dns_ptr_record.test", "tags.%", "2"),
				),
			},

			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsPtrRecordExists("azurerm_dns_ptr_record.test"),
					resource.TestCheckResourceAttr(
						"azurerm_dns_ptr_record.test", "tags.%", "1"),
				),
			},
		},
	})
}

func testCheckAzureRMDnsPtrRecordExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		ptrName := rs.Primary.Attributes["name"]
		zoneName := rs.Primary.Attributes["zone_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for DNS PTR record: %s", ptrName)
		}

		conn := testAccProvider.Meta().(*ArmClient).dnsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := conn.Get(ctx, resourceGroup, zoneName, ptrName, dns.PTR)
		if err != nil {
			return fmt.Errorf("Bad: Get PTR RecordSet: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: DNS PTR record %s (resource group: %s) does not exist", ptrName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDnsPtrRecordDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).dnsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_dns_ptr_record" {
			continue
		}

		ptrName := rs.Primary.Attributes["name"]
		zoneName := rs.Primary.Attributes["zone_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, zoneName, ptrName, dns.PTR)

		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}

			return err
		}

		return fmt.Errorf("DNS PTR record still exists:\n%#v", resp)
	}

	return nil
}

func testAccAzureRMDnsPtrRecord_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_ptr_record" "test" {
  name                = "testptrrecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300
  records             = ["hashicorp.com", "microsoft.com"]
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMDnsPtrRecord_updateRecords(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_ptr_record" "test" {
  name                = "testptrrecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300
  records             = ["hashicorp.com", "microsoft.com", "reddit.com"]
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMDnsPtrRecord_withTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_ptr_record" "test" {
  name                = "testptrrecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300
  records             = ["hashicorp.com", "microsoft.com"]

  tags {
    environment = "Dev"
    cost_center = "Ops"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMDnsPtrRecord_withTagsUpdate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_ptr_record" "test" {
  name                = "testptrrecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300
  records             = ["hashicorp.com", "microsoft.com"]

  tags {
    environment = "Stage"
  }
}
`, rInt, location, rInt, rInt)
}

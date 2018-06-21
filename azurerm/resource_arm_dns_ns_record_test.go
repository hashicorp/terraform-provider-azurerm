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

//TODO: remove this once we remove the `record` attribute
func TestAccAzureRMDnsNsRecord_deprecatedBasic(t *testing.T) {
	resourceName := "azurerm_dns_ns_record.test"
	ri := acctest.RandInt()
	config := testAccAzureRMDnsNsRecord_deprecatedBasic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDnsNsRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsNsRecordExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMDnsNsRecord_basic(t *testing.T) {
	resourceName := "azurerm_dns_ns_record.test"
	ri := acctest.RandInt()
	config := testAccAzureRMDnsNsRecord_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDnsNsRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsNsRecordExists(resourceName),
				),
			},
		},
	})
}

//TODO: remove this once we remove the `record` attribute
func TestAccAzureRMDnsNsRecord_deprecatedUpdateRecords(t *testing.T) {
	resourceName := "azurerm_dns_ns_record.test"
	ri := acctest.RandInt()
	location := testLocation()
	preConfig := testAccAzureRMDnsNsRecord_deprecatedBasic(ri, location)
	postConfig := testAccAzureRMDnsNsRecord_deprecatedUpdateRecords(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDnsNsRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsNsRecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "record.#", "2"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsNsRecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "record.#", "3"),
				),
			},
		},
	})
}

func TestAccAzureRMDnsNsRecord_updateRecords(t *testing.T) {
	resourceName := "azurerm_dns_ns_record.test"
	ri := acctest.RandInt()
	location := testLocation()
	preConfig := testAccAzureRMDnsNsRecord_basic(ri, location)
	postConfig := testAccAzureRMDnsNsRecord_updateRecords(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDnsNsRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsNsRecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "records.#", "2"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsNsRecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "records.#", "3"),
				),
			},
		},
	})
}

//TODO: remove this once we remove the `record` attribute
func TestAccAzureRMDnsNsRecord_deprecatedChangeRecordToRecords(t *testing.T) {
	resourceName := "azurerm_dns_ns_record.test"
	ri := acctest.RandInt()
	location := testLocation()
	preConfig := testAccAzureRMDnsNsRecord_deprecatedBasic(ri, location)
	postConfig := testAccAzureRMDnsNsRecord_deprecatedBasicNewRecords(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDnsNsRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsNsRecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "records.#", "2"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsNsRecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "records.#", "2"),
				),
			},
		},
	})
}

//TODO: remove this once we remove the `record` attribute
func TestAccAzureRMDnsNsRecord_deprecatedWithTags(t *testing.T) {
	resourceName := "azurerm_dns_ns_record.test"
	ri := acctest.RandInt()
	location := testLocation()
	preConfig := testAccAzureRMDnsNsRecord_deprecatedWithTags(ri, location)
	postConfig := testAccAzureRMDnsNsRecord_deprecatedWithTagsUpdate(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDnsNsRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsNsRecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsNsRecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
				),
			},
		},
	})
}

func TestAccAzureRMDnsNsRecord_withTags(t *testing.T) {
	resourceName := "azurerm_dns_ns_record.test"
	ri := acctest.RandInt()
	location := testLocation()
	preConfig := testAccAzureRMDnsNsRecord_withTags(ri, location)
	postConfig := testAccAzureRMDnsNsRecord_withTagsUpdate(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDnsNsRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsNsRecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsNsRecordExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
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
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := conn.Get(ctx, resourceGroup, zoneName, nsName, dns.NS)
		if err != nil {
			return fmt.Errorf("Bad: Get DNS NS Record: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: DNS NS record %s (resource group: %s) does not exist", nsName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDnsNsRecordDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).dnsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_dns_ns_record" {
			continue
		}

		nsName := rs.Primary.Attributes["name"]
		zoneName := rs.Primary.Attributes["zone_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, zoneName, nsName, dns.NS)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("DNS NS Record still exists:\n%#v", resp.RecordSetProperties)
		}

	}

	return nil
}

func testAccAzureRMDnsNsRecord_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_ns_record" "test" {
  name                = "mynsrecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300

  records = ["ns1.contoso.com", "ns2.contoso.com"]
}
`, rInt, location, rInt, rInt)
}

//TODO: remove this once we remove the `record` attribute
func testAccAzureRMDnsNsRecord_deprecatedBasic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_ns_record" "test" {
  name                = "mynsrecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300

  record {
	  nsdname = "ns1.contoso.com"
  }

  record {
	  nsdname = "ns2.contoso.com"
  }
}
`, rInt, location, rInt, rInt)
}

//TODO: remove this once we remove the `record` attribute
func testAccAzureRMDnsNsRecord_deprecatedBasicNewRecords(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_ns_record" "test" {
  name                = "mynsrecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300

  records = ["ns2.contoso.com", "ns1.contoso.com"]
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMDnsNsRecord_updateRecords(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_ns_record" "test" {
  name                = "mynsrecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300

  records = ["ns1.contoso.com", "ns2.contoso.com", "ns3.contoso.com"]
}
`, rInt, location, rInt, rInt)
}

//TODO: remove this once we remove the `record` attribute
func testAccAzureRMDnsNsRecord_deprecatedUpdateRecords(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_ns_record" "test" {
  name                = "mynsrecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300

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
`, rInt, location, rInt, rInt)
}

func testAccAzureRMDnsNsRecord_withTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_ns_record" "test" {
  name                = "mynsrecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300

  records = ["ns1.contoso.com", "ns2.contoso.com"]

  tags {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, rInt, location, rInt, rInt)
}

//TODO: remove this once we remove the `record` attribute
func testAccAzureRMDnsNsRecord_deprecatedWithTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_ns_record" "test" {
  name                = "mynsrecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300

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
`, rInt, location, rInt, rInt)
}

func testAccAzureRMDnsNsRecord_withTagsUpdate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_ns_record" "test" {
  name                = "mynsrecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300

  records = ["ns1.contoso.com", "ns2.contoso.com"]

  tags {
    environment = "staging"
  }
}
`, rInt, location, rInt, rInt)
}

//TODO: remove this once we remove the `record` attribute
func testAccAzureRMDnsNsRecord_deprecatedWithTagsUpdate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_ns_record" "test" {
  name                = "mynsrecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300

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
`, rInt, location, rInt, rInt)
}

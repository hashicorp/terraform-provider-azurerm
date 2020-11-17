package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/dns/mgmt/2018-05-01/dns"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/dns/parse"
)

func TestAccAzureRMDnsNsRecord_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_ns_record", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsNsRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDnsNsRecord_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsNsRecordExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "fqdn"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDnsNsRecord_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_ns_record", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsNsRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDnsNsRecord_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsNsRecordExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMDnsNsRecord_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_dns_ns_record"),
			},
		},
	})
}

func TestAccAzureRMDnsNsRecord_updateRecords(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_ns_record", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsNsRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDnsNsRecord_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsNsRecordExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "records.#", "2"),
				),
			},
			{
				Config: testAccAzureRMDnsNsRecord_updateRecords(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsNsRecordExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "records.#", "3"),
				),
			},
		},
	})
}

func TestAccAzureRMDnsNsRecord_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_ns_record", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsNsRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDnsNsRecord_withTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsNsRecordExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
				),
			},
			{
				Config: testAccAzureRMDnsNsRecord_withTagsUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsNsRecordExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMDnsNsRecordExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Dns.RecordSetsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.DnsNsRecordID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := conn.Get(ctx, id.ResourceGroup, id.ZoneName, id.Name, dns.NS)
		if err != nil {
			return fmt.Errorf("Bad: Get DNS NS Record: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: DNS NS record %s (resource group: %s) does not exist", id.Name, id.ResourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDnsNsRecordDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Dns.RecordSetsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_dns_ns_record" {
			continue
		}

		id, err := parse.DnsNsRecordID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := conn.Get(ctx, id.ResourceGroup, id.ZoneName, id.Name, dns.NS)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("DNS NS Record still exists:\n%#v", resp.RecordSetProperties)
		}
	}

	return nil
}

func testAccAzureRMDnsNsRecord_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_dns_ns_record" "test" {
  name                = "mynsrecord%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_dns_zone.test.name
  ttl                 = 300

  records = ["ns1.contoso.com", "ns2.contoso.com"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDnsNsRecord_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMDnsNsRecord_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_dns_ns_record" "import" {
  name                = azurerm_dns_ns_record.test.name
  resource_group_name = azurerm_dns_ns_record.test.resource_group_name
  zone_name           = azurerm_dns_ns_record.test.zone_name
  ttl                 = 300

  records = ["ns1.contoso.com", "ns2.contoso.com"]
}
`, template)
}

func testAccAzureRMDnsNsRecord_updateRecords(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_dns_ns_record" "test" {
  name                = "mynsrecord%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_dns_zone.test.name
  ttl                 = 300

  records = ["ns1.contoso.com", "ns2.contoso.com", "ns3.contoso.com"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDnsNsRecord_withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_dns_ns_record" "test" {
  name                = "mynsrecord%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_dns_zone.test.name
  ttl                 = 300

  records = ["ns1.contoso.com", "ns2.contoso.com"]

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDnsNsRecord_withTagsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_dns_ns_record" "test" {
  name                = "mynsrecord%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_dns_zone.test.name
  ttl                 = 300

  records = ["ns1.contoso.com", "ns2.contoso.com"]

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

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

func TestAccAzureRMDnsAAAARecord_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_aaaa_record", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsAaaaRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDnsAAAARecord_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsAaaaRecordExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "fqdn"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDnsAAAARecord_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_aaaa_record", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsAaaaRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDnsAAAARecord_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsAaaaRecordExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMDnsAAAARecord_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_dns_aaaa_record"),
			},
		},
	})
}

func TestAccAzureRMDnsAAAARecord_updateRecords(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_aaaa_record", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsAaaaRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDnsAAAARecord_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsAaaaRecordExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "records.#", "2"),
				),
			},
			{
				Config: testAccAzureRMDnsAAAARecord_updateRecords(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsAaaaRecordExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "records.#", "3"),
				),
			},
		},
	})
}

func TestAccAzureRMDnsAAAARecord_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_aaaa_record", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsAaaaRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDnsAAAARecord_withTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsAaaaRecordExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
				),
			},
			{
				Config: testAccAzureRMDnsAAAARecord_withTagsUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsAaaaRecordExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDnsAAAARecord_withAlias(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_aaaa_record", "test")
	targetResourceName := "azurerm_public_ip.test"
	targetResourceName2 := "azurerm_public_ip.test2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsAaaaRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDnsAAAARecord_withAlias(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsAaaaRecordExists(data.ResourceName),
					resource.TestCheckResourceAttrPair(data.ResourceName, "target_resource_id", targetResourceName, "id"),
				),
			},
			{
				Config: testAccAzureRMDnsAAAARecord_withAliasUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsAaaaRecordExists(data.ResourceName),
					resource.TestCheckResourceAttrPair(data.ResourceName, "target_resource_id", targetResourceName2, "id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDnsAAAARecord_RecordsToAlias(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_aaaa_record", "test")
	targetResourceName := "azurerm_public_ip.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsAaaaRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDnsAaaaRecord_AliasToRecordsUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsAaaaRecordExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "records.#", "2"),
				),
			},
			{
				Config: testAccAzureRMDnsAaaaRecord_AliasToRecords(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsAaaaRecordExists(data.ResourceName),
					resource.TestCheckResourceAttrPair(data.ResourceName, "target_resource_id", targetResourceName, "id"),
					resource.TestCheckNoResourceAttr(data.ResourceName, "records"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDnsAaaaRecord_AliasToRecords(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_aaaa_record", "test")
	targetResourceName := "azurerm_public_ip.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsAaaaRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDnsAaaaRecord_AliasToRecords(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsAaaaRecordExists(data.ResourceName),
					resource.TestCheckResourceAttrPair(data.ResourceName, "target_resource_id", targetResourceName, "id"),
				),
			},
			{
				Config: testAccAzureRMDnsAaaaRecord_AliasToRecordsUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsAaaaRecordExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "records.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "target_resource_id", ""),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDnsAAAARecord_uncompressed(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_aaaa_record", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsAaaaRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDnsAAAARecord_uncompressed(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsAaaaRecordExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "fqdn"),
				),
			},
			{
				Config: testAccAzureRMDnsAAAARecord_uncompressed(data), // just use the same for updating
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsAaaaRecordExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "records.#", "2"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMDnsAaaaRecordExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Dns.RecordSetsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.AaaaRecordID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := conn.Get(ctx, id.ResourceGroup, id.DnszoneName, id.AAAAName, dns.AAAA)
		if err != nil {
			return fmt.Errorf("Bad: Get AAAA RecordSet: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: DNS AAAA record %s (resource group: %s) does not exist", id.AAAAName, id.ResourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDnsAaaaRecordDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Dns.RecordSetsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_dns_aaaa_record" {
			continue
		}

		id, err := parse.AaaaRecordID(rs.Primary.ID)
		if err != nil {
			return err
		}
		resp, err := conn.Get(ctx, id.ResourceGroup, id.DnszoneName, id.AAAAName, dns.AAAA)
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}

			return err
		}

		return fmt.Errorf("DNS AAAA record still exists:\n%#v", resp.RecordSetProperties)
	}

	return nil
}

func testAccAzureRMDnsAAAARecord_basic(data acceptance.TestData) string {
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

resource "azurerm_dns_aaaa_record" "test" {
  name                = "myarecord%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_dns_zone.test.name
  ttl                 = 300
  records             = ["2607:f8b0:4009:1803::1005", "2607:f8b0:4009:1803::1006"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDnsAAAARecord_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMDnsAAAARecord_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_dns_aaaa_record" "import" {
  name                = azurerm_dns_aaaa_record.test.name
  resource_group_name = azurerm_dns_aaaa_record.test.resource_group_name
  zone_name           = azurerm_dns_aaaa_record.test.zone_name
  ttl                 = 300
  records             = ["2607:f8b0:4009:1803::1005", "2607:f8b0:4009:1803::1006"]
}
`, template)
}

func testAccAzureRMDnsAAAARecord_updateRecords(data acceptance.TestData) string {
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

resource "azurerm_dns_aaaa_record" "test" {
  name                = "myarecord%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_dns_zone.test.name
  ttl                 = 300
  records             = ["2607:f8b0:4009:1803::1005", "2607:f8b0:4009:1803::1006", "::1"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDnsAAAARecord_withTags(data acceptance.TestData) string {
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

resource "azurerm_dns_aaaa_record" "test" {
  name                = "myarecord%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_dns_zone.test.name
  ttl                 = 300
  records             = ["2607:f8b0:4009:1803::1005", "2607:f8b0:4009:1803::1006"]

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDnsAAAARecord_withTagsUpdate(data acceptance.TestData) string {
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

resource "azurerm_dns_aaaa_record" "test" {
  name                = "myarecord%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_dns_zone.test.name
  ttl                 = 300
  records             = ["2607:f8b0:4009:1803::1005", "2607:f8b0:4009:1803::1006"]

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDnsAAAARecord_withAlias(data acceptance.TestData) string {
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

resource "azurerm_public_ip" "test" {
  name                = "mypublicip%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
  ip_version          = "IPv6"
}

resource "azurerm_dns_aaaa_record" "test" {
  name                = "myaaaarecord%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_dns_zone.test.name
  ttl                 = 300
  target_resource_id  = azurerm_public_ip.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDnsAAAARecord_withAliasUpdate(data acceptance.TestData) string {
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

resource "azurerm_public_ip" "test2" {
  name                = "mypublicip%d2"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
  ip_version          = "IPv6"
}

resource "azurerm_dns_aaaa_record" "test" {
  name                = "myaaaarecord%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_dns_zone.test.name
  ttl                 = 300
  target_resource_id  = azurerm_public_ip.test2.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDnsAaaaRecord_AliasToRecords(data acceptance.TestData) string {
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

resource "azurerm_public_ip" "test" {
  name                = "mypublicip%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
  ip_version          = "IPv6"
}

resource "azurerm_dns_aaaa_record" "test" {
  name                = "myarecord%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_dns_zone.test.name
  ttl                 = 300
  target_resource_id  = azurerm_public_ip.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDnsAaaaRecord_AliasToRecordsUpdate(data acceptance.TestData) string {
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

resource "azurerm_dns_aaaa_record" "test" {
  name                = "myarecord%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_dns_zone.test.name
  ttl                 = 300
  records             = ["3a62:353:8885:293c:a218:45cc:9ee9:4e27", "3a62:353:8885:293c:a218:45cc:9ee9:4e28"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDnsAAAARecord_uncompressed(data acceptance.TestData) string {
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

resource "azurerm_dns_aaaa_record" "test" {
  name                = "myarecord%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_dns_zone.test.name
  ttl                 = 300
  records             = ["2607:f8b0:4005:0800:0000:0000:0000:1003", "2201:1234:1234::1"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

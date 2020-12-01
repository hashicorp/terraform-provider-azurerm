package dns_test

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

func TestAccAzureRMDnsARecord_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_a_record", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsARecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDnsARecord_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsARecordExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "fqdn"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDnsARecord_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_a_record", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsARecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDnsARecord_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsARecordExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMDnsARecord_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_dns_a_record"),
			},
		},
	})
}

func TestAccAzureRMDnsARecord_updateRecords(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_a_record", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsARecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDnsARecord_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsARecordExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "records.#", "2"),
				),
			},
			{
				Config: testAccAzureRMDnsARecord_updateRecords(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsARecordExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "records.#", "3"),
				),
			},
		},
	})
}

func TestAccAzureRMDnsARecord_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_a_record", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsARecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDnsARecord_withTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsARecordExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
				),
			},
			{
				Config: testAccAzureRMDnsARecord_withTagsUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsARecordExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDnsARecord_withAlias(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_a_record", "test")
	targetResourceName := "azurerm_public_ip.test"
	targetResourceName2 := "azurerm_public_ip.test2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsARecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDnsARecord_withAlias(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsARecordExists(data.ResourceName),
					resource.TestCheckResourceAttrPair(data.ResourceName, "target_resource_id", targetResourceName, "id"),
				),
			},
			{
				Config: testAccAzureRMDnsARecord_withAliasUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsARecordExists(data.ResourceName),
					resource.TestCheckResourceAttrPair(data.ResourceName, "target_resource_id", targetResourceName2, "id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDnsARecord_RecordsToAlias(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_a_record", "test")
	targetResourceName := "azurerm_public_ip.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsARecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDnsARecord_AliasToRecordsUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsARecordExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "records.#", "2"),
				),
			},
			{
				Config: testAccAzureRMDnsARecord_AliasToRecords(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsARecordExists(data.ResourceName),
					resource.TestCheckResourceAttrPair(data.ResourceName, "target_resource_id", targetResourceName, "id"),
					resource.TestCheckResourceAttr(data.ResourceName, "records.#", "0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDnsARecord_AliasToRecords(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_a_record", "test")
	targetResourceName := "azurerm_public_ip.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsARecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDnsARecord_AliasToRecords(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsARecordExists(data.ResourceName),
					resource.TestCheckResourceAttrPair(data.ResourceName, "target_resource_id", targetResourceName, "id"),
				),
			},
			{
				Config: testAccAzureRMDnsARecord_AliasToRecordsUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsARecordExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "records.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "target_resource_id", ""),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMDnsARecordExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Dns.RecordSetsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.ARecordID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := conn.Get(ctx, id.ResourceGroup, id.DnszoneName, id.AName, dns.A)
		if err != nil {
			return fmt.Errorf("Bad: Get A RecordSet: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: DNS A record %s (resource group: %s) does not exist", id.AName, id.ResourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDnsARecordDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Dns.RecordSetsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_dns_a_record" {
			continue
		}

		id, err := parse.ARecordID(rs.Primary.ID)
		if err != nil {
			return err
		}
		resp, err := conn.Get(ctx, id.ResourceGroup, id.DnszoneName, id.AName, dns.A)
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

func testAccAzureRMDnsARecord_basic(data acceptance.TestData) string {
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

resource "azurerm_dns_a_record" "test" {
  name                = "myarecord%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_dns_zone.test.name
  ttl                 = 300
  records             = ["1.2.3.4", "1.2.4.5"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDnsARecord_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMDnsARecord_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_dns_a_record" "import" {
  name                = azurerm_dns_a_record.test.name
  resource_group_name = azurerm_dns_a_record.test.resource_group_name
  zone_name           = azurerm_dns_a_record.test.zone_name
  ttl                 = 300
  records             = ["1.2.3.4", "1.2.4.5"]
}
`, template)
}

func testAccAzureRMDnsARecord_updateRecords(data acceptance.TestData) string {
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

resource "azurerm_dns_a_record" "test" {
  name                = "myarecord%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_dns_zone.test.name
  ttl                 = 300
  records             = ["1.2.3.4", "1.2.4.5", "1.2.3.7"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDnsARecord_withTags(data acceptance.TestData) string {
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

resource "azurerm_dns_a_record" "test" {
  name                = "myarecord%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_dns_zone.test.name
  ttl                 = 300
  records             = ["1.2.3.4", "1.2.4.5"]

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDnsARecord_withTagsUpdate(data acceptance.TestData) string {
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

resource "azurerm_dns_a_record" "test" {
  name                = "myarecord%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_dns_zone.test.name
  ttl                 = 300
  records             = ["1.2.3.4", "1.2.4.5"]

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDnsARecord_withAlias(data acceptance.TestData) string {
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
  ip_version          = "IPv4"
}

resource "azurerm_dns_a_record" "test" {
  name                = "myarecord%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_dns_zone.test.name
  ttl                 = 300
  target_resource_id  = azurerm_public_ip.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDnsARecord_withAliasUpdate(data acceptance.TestData) string {
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
  ip_version          = "IPv4"
}

resource "azurerm_dns_a_record" "test" {
  name                = "myarecord%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_dns_zone.test.name
  ttl                 = 300
  target_resource_id  = azurerm_public_ip.test2.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDnsARecord_AliasToRecords(data acceptance.TestData) string {
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
  ip_version          = "IPv4"
}

resource "azurerm_dns_a_record" "test" {
  name                = "myarecord%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_dns_zone.test.name
  ttl                 = 300
  target_resource_id  = azurerm_public_ip.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDnsARecord_AliasToRecordsUpdate(data acceptance.TestData) string {
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

resource "azurerm_dns_a_record" "test" {
  name                = "myarecord%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_dns_zone.test.name
  ttl                 = 300
  records             = ["1.2.3.4", "1.2.4.5"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

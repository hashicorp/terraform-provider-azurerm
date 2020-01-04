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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMDnsCaaRecord_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_caa_record", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsCaaRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDnsCaaRecord_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsCaaRecordExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "fqdn"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDnsCaaRecord_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_dns_caa_record", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsCaaRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDnsCaaRecord_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsCaaRecordExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMDnsCaaRecord_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_dns_caa_record"),
			},
		},
	})
}

func TestAccAzureRMDnsCaaRecord_updateRecords(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_caa_record", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsCaaRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDnsCaaRecord_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsCaaRecordExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "record.#", "4"),
				),
			},
			{
				Config: testAccAzureRMDnsCaaRecord_updateRecords(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsCaaRecordExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "record.#", "5"),
				),
			},
		},
	})
}

func TestAccAzureRMDnsCaaRecord_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_caa_record", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsCaaRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDnsCaaRecord_withTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsCaaRecordExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
				),
			},
			{
				Config: testAccAzureRMDnsCaaRecord_withTagsUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDnsCaaRecordExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMDnsCaaRecordExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		caaName := rs.Primary.Attributes["name"]
		zoneName := rs.Primary.Attributes["zone_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for DNS CAA record: %s", caaName)
		}

		conn := acceptance.AzureProvider.Meta().(*clients.Client).Dns.RecordSetsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		resp, err := conn.Get(ctx, resourceGroup, zoneName, caaName, dns.CAA)
		if err != nil {
			return fmt.Errorf("Bad: Get CAA RecordSet: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: DNS CAA record %s (resource group: %s) does not exist", caaName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDnsCaaRecordDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Dns.RecordSetsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_dns_caa_record" {
			continue
		}

		caaName := rs.Primary.Attributes["name"]
		zoneName := rs.Primary.Attributes["zone_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, zoneName, caaName, dns.CAA)

		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}

			return err
		}

		return fmt.Errorf("DNS CAA record still exists:\n%#v", resp.RecordSetProperties)
	}

	return nil
}

func testAccAzureRMDnsCaaRecord_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_caa_record" "test" {
  name                = "myarecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300

  record {
    flags = 0
    tag   = "issue"
    value = "example.com"
  }

  record {
    flags = 0
    tag   = "issue"
    value = "example.net"
  }

  record {
    flags = 1
    tag   = "issuewild"
    value = ";"
  }

  record {
    flags = 0
    tag   = "iodef"
    value = "mailto:terraform@nonexist.tld"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDnsCaaRecord_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMDnsCaaRecord_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_dns_caa_record" "import" {
  name                = "${azurerm_dns_caa_record.test.name}"
  resource_group_name = "${azurerm_dns_caa_record.test.resource_group_name}"
  zone_name           = "${azurerm_dns_caa_record.test.zone_name}"
  ttl                 = 300

  record {
    flags = 0
    tag   = "issue"
    value = "example.com"
  }

  record {
    flags = 0
    tag   = "issue"
    value = "example.net"
  }

  record {
    flags = 1
    tag   = "issuewild"
    value = ";"
  }

  record {
    flags = 0
    tag   = "iodef"
    value = "mailto:terraform@nonexist.tld"
  }
}
`, template)
}

func testAccAzureRMDnsCaaRecord_updateRecords(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_caa_record" "test" {
  name                = "myarecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300

  record {
    flags = 0
    tag   = "issue"
    value = "example.com"
  }

  record {
    flags = 0
    tag   = "issue"
    value = "example.net"
  }

  record {
    flags = 1
    tag   = "issuewild"
    value = ";"
  }

  record {
    flags = 0
    tag   = "iodef"
    value = "mailto:terraform@nonexist.tld"
  }

  record {
    flags = 0
    tag   = "issue"
    value = "letsencrypt.org"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDnsCaaRecord_withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_caa_record" "test" {
  name                = "myarecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300

  record {
    flags = 0
    tag   = "issue"
    value = "example.net"
  }

  record {
    flags = 1
    tag   = "issuewild"
    value = ";"
  }

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDnsCaaRecord_withTagsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_caa_record" "test" {
  name                = "myarecord%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  zone_name           = "${azurerm_dns_zone.test.name}"
  ttl                 = 300

  record {
    flags = 0
    tag   = "issue"
    value = "example.net"
  }

  record {
    flags = 1
    tag   = "issuewild"
    value = ";"
  }

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

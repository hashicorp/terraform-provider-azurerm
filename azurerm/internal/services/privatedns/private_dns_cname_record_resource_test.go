package privatedns_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/privatedns/mgmt/2018-09-01/privatedns"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMPrivateDnsCNameRecord_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_cname_record", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsCNameRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateDnsCNameRecord_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsCNameRecordExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "fqdn"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPrivateDnsCNameRecord_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_cname_record", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsCNameRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateDnsCNameRecord_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsCNameRecordExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMPrivateDnsCNameRecord_requiresImport),
		},
	})
}

func TestAccAzureRMPrivateDnsCNameRecord_subdomain(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_cname_record", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsCNameRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateDnsCNameRecord_subdomain(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsCNameRecordExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "record", "test.contoso.com"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPrivateDnsCNameRecord_updateRecords(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_cname_record", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsCNameRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateDnsCNameRecord_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsCNameRecordExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMPrivateDnsCNameRecord_updateRecords(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsCNameRecordExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMPrivateDnsCNameRecord_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_cname_record", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsCNameRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateDnsCNameRecord_withTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsCNameRecordExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
				),
			},
			{
				Config: testAccAzureRMPrivateDnsCNameRecord_withTagsUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsCNameRecordExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMPrivateDnsCNameRecordExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).PrivateDns.RecordSetsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		aName := rs.Primary.Attributes["name"]
		zoneName := rs.Primary.Attributes["zone_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Private DNS CNAME record: %s", aName)
		}

		resp, err := conn.Get(ctx, resourceGroup, zoneName, privatedns.CNAME, aName)
		if err != nil {
			return fmt.Errorf("Bad: Get CNAME RecordSet: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Private DNS CNAME record %s (resource group: %s) does not exist", aName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMPrivateDnsCNameRecordDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).PrivateDns.RecordSetsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_private_dns_cname_record" {
			continue
		}

		aName := rs.Primary.Attributes["name"]
		zoneName := rs.Primary.Attributes["zone_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, zoneName, privatedns.CNAME, aName)
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}

			return err
		}

		return fmt.Errorf("Private DNS CNAME record still exists:\n%#v", resp.RecordSetProperties)
	}

	return nil
}

func testAccAzureRMPrivateDnsCNameRecord_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_dns_cname_record" "test" {
  name                = "acctestcname%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_private_dns_zone.test.name
  ttl                 = 300
  record              = "contoso.com"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMPrivateDnsCNameRecord_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMPrivateDnsCNameRecord_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_private_dns_cname_record" "import" {
  name                = azurerm_private_dns_cname_record.test.name
  resource_group_name = azurerm_private_dns_cname_record.test.resource_group_name
  zone_name           = azurerm_private_dns_cname_record.test.zone_name
  ttl                 = 300
  record              = "contoso.com"
}
`, template)
}

func testAccAzureRMPrivateDnsCNameRecord_subdomain(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_dns_cname_record" "test" {
  name                = "acctestcname%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_private_dns_zone.test.name
  ttl                 = 300
  record              = "test.contoso.com"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMPrivateDnsCNameRecord_updateRecords(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_dns_cname_record" "test" {
  name                = "acctestcname%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_private_dns_zone.test.name
  ttl                 = 300
  record              = "contoso.com"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMPrivateDnsCNameRecord_withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_dns_cname_record" "test" {
  name                = "acctestcname%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_private_dns_zone.test.name
  ttl                 = 300
  record              = "contoso.com"

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMPrivateDnsCNameRecord_withTagsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_dns_cname_record" "test" {
  name                = "acctestcname%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_private_dns_zone.test.name
  ttl                 = 300
  record              = "contoso.com"

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

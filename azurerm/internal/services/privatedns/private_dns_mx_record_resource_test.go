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

func TestAccAzureRMPrivateDnsMxRecord_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_mx_record", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsMxRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateDnsMxRecord_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsMxRecordExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "fqdn"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPrivateDnsMxRecord_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_mx_record", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsMxRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateDnsMxRecord_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsMxRecordExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMPrivateDnsMxRecord_requiresImport),
		},
	})
}

func TestAccAzureRMPrivateDnsMxRecord_updateRecords(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_mx_record", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsMxRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateDnsMxRecord_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsMxRecordExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "record.#", "2"),
				),
			},
			{
				Config: testAccAzureRMPrivateDnsMxRecord_updateRecords(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsMxRecordExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "record.#", "3"),
				),
			},
		},
	})
}

func TestAccAzureRMPrivateDnsMxRecord_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_mx_record", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsMxRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateDnsMxRecord_withTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsMxRecordExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
				),
			},
			{
				Config: testAccAzureRMPrivateDnsMxRecord_withTagsUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsMxRecordExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPrivateDnsMxRecord_emptyName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_mx_record", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsMxRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateDnsMxRecord_emptyName(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsMxRecordExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMPrivateDnsMxRecordExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).PrivateDns.RecordSetsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		mxName := rs.Primary.Attributes["name"]
		zoneName := rs.Primary.Attributes["zone_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Private DNS MX record: %s", mxName)
		}

		resp, err := conn.Get(ctx, resourceGroup, zoneName, privatedns.MX, mxName)
		if err != nil {
			return fmt.Errorf("Bad: Get MX RecordSet: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Private DNS MX record %s (resource group: %s) does not exist", mxName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMPrivateDnsMxRecordDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).PrivateDns.RecordSetsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_private_dns_mx_record" {
			continue
		}

		mxName := rs.Primary.Attributes["name"]
		zoneName := rs.Primary.Attributes["zone_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, zoneName, privatedns.MX, mxName)
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}

			return err
		}

		return fmt.Errorf("Private DNS MX record still exists:\n%#v", resp.RecordSetProperties)
	}

	return nil
}

func testAccAzureRMPrivateDnsMxRecord_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-prvdns-%d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "testzone%d.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_dns_mx_record" "test" {
  name                = "testaccmx%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_private_dns_zone.test.name
  ttl                 = 300
  record {
    preference = 10
    exchange   = "mx1.contoso.com"
  }

  record {
    preference = 10
    exchange   = "mx2.contoso.com"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMPrivateDnsMxRecord_emptyName(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-prvdns-%d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "testzone%d.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_dns_mx_record" "test" {
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_private_dns_zone.test.name
  ttl                 = 300
  record {
    preference = 10
    exchange   = "mx1.contoso.com"
  }

  record {
    preference = 10
    exchange   = "mx2.contoso.com"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMPrivateDnsMxRecord_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMPrivateDnsMxRecord_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_private_dns_mx_record" "import" {
  name                = azurerm_private_dns_mx_record.test.name
  resource_group_name = azurerm_private_dns_mx_record.test.resource_group_name
  zone_name           = azurerm_private_dns_mx_record.test.zone_name
  ttl                 = 300
  record {
    preference = 10
    exchange   = "mx1.contoso.com"
  }
  record {
    preference = 10
    exchange   = "mx2.contoso.com"
  }
}
`, template)
}

func testAccAzureRMPrivateDnsMxRecord_updateRecords(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-prvdns-%d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "testzone%d.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_dns_mx_record" "test" {
  name                = "testaccmx%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_private_dns_zone.test.name
  ttl                 = 300
  record {
    preference = 10
    exchange   = "mx1.contoso.com"
  }
  record {
    preference = 10
    exchange   = "mx2.contoso.com"
  }
  record {
    preference = 20
    exchange   = "backupmx.contoso.com"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMPrivateDnsMxRecord_withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-prvdns-%d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "testzone%d.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_dns_mx_record" "test" {
  name                = "testaccmx%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_private_dns_zone.test.name
  ttl                 = 300
  record {
    preference = 10
    exchange   = "mx1.contoso.com"
  }
  record {
    preference = 10
    exchange   = "mx2.contoso.com"
  }

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMPrivateDnsMxRecord_withTagsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-prvdns-%d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "testzone%d.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_dns_mx_record" "test" {
  name                = "testaccmx%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_private_dns_zone.test.name
  ttl                 = 300
  record {
    preference = 10
    exchange   = "mx1.contoso.com"
  }
  record {
    preference = 10
    exchange   = "mx2.contoso.com"
  }

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

package tests

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

func TestAccAzureRMPrivateDnsTxtRecord_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_txt_record", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsTxtRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateDnsTxtRecord_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsTxtRecordExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "fqdn"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPrivateDnsTxtRecord_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_txt_record", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsTxtRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateDnsTxtRecord_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsTxtRecordExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMPrivateDnsTxtRecord_requiresImport),
		},
	})
}

func TestAccAzureRMPrivateDnsTxtRecord_updateRecords(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_txt_record", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsTxtRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateDnsTxtRecord_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsTxtRecordExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "record.#", "2"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMPrivateDnsTxtRecord_updateRecords(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsTxtRecordExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "record.#", "3"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMPrivateDnsTxtRecord_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsTxtRecordExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "record.#", "2"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPrivateDnsTxtRecord_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_txt_record", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPrivateDnsTxtRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPrivateDnsTxtRecord_withTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsTxtRecordExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMPrivateDnsTxtRecord_withTagsUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPrivateDnsTxtRecordExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMPrivateDnsTxtRecordExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).PrivateDns.RecordSetsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		txtName := rs.Primary.Attributes["name"]
		zoneName := rs.Primary.Attributes["zone_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Private DNS TXT record: %s", txtName)
		}

		resp, err := conn.Get(ctx, resourceGroup, zoneName, privatedns.TXT, txtName)
		if err != nil {
			return fmt.Errorf("Bad: Get TXT RecordSet: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Private DNS TXT record %s (resource group: %s) does not exist", txtName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMPrivateDnsTxtRecordDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).PrivateDns.RecordSetsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_private_dns_txt_record" {
			continue
		}

		txtName := rs.Primary.Attributes["name"]
		zoneName := rs.Primary.Attributes["zone_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, zoneName, privatedns.TXT, txtName)
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}

			return err
		}

		return fmt.Errorf("Private DNS TXT record still exists:\n%#v", resp.RecordSetProperties)
	}

	return nil
}

func testAccAzureRMPrivateDnsTxtRecord_basic(data acceptance.TestData) string {
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

resource "azurerm_private_dns_txt_record" "test" {
  name                = "testacctxt%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_private_dns_zone.test.name
  ttl                 = 300

  record {
    value = "Quick brown fox"
  }

  record {
    value = "A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMPrivateDnsTxtRecord_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMPrivateDnsTxtRecord_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_private_dns_txt_record" "import" {
  name                = azurerm_private_dns_txt_record.test.name
  resource_group_name = azurerm_private_dns_txt_record.test.resource_group_name
  zone_name           = azurerm_private_dns_txt_record.test.zone_name
  ttl                 = 300

  record {
    value = "Quick brown fox"
  }

  record {
    value = "A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......"
  }
}
`, template)
}

func testAccAzureRMPrivateDnsTxtRecord_updateRecords(data acceptance.TestData) string {
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

resource "azurerm_private_dns_txt_record" "test" {
  name                = "test%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_private_dns_zone.test.name
  ttl                 = 300

  record {
    value = "Quick brown fox"
  }

  record {
    value = "A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......"
  }

  record {
    value = "I'm a record too'"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMPrivateDnsTxtRecord_withTags(data acceptance.TestData) string {
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

resource "azurerm_private_dns_txt_record" "test" {
  name                = "test%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_private_dns_zone.test.name
  ttl                 = 300

  record {
    value = "Quick brown fox"
  }

  record {
    value = "A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......"
  }

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMPrivateDnsTxtRecord_withTagsUpdate(data acceptance.TestData) string {
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

resource "azurerm_private_dns_txt_record" "test" {
  name                = "test%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_private_dns_zone.test.name
  ttl                 = 300

  record {
    value = "Quick brown fox"
  }

  record {
    value = "A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......A long text......"
  }

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

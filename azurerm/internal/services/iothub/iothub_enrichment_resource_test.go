package iothub_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMIotHubEnrichment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_enrichment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotHubEnrichmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHubEnrichment_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubEnrichmentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMIotHubEnrichment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_enrichment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotHubEnrichmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHubEnrichment_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubEnrichmentExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMIotHubEnrichment_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_iothub_enrichment"),
			},
		},
	})
}

func TestAccAzureRMIotHubEnrichment_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_enrichment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotHubEnrichmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHubEnrichment_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubEnrichmentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMIotHubEnrichment_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubEnrichmentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMIotHubEnrichmentDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).IoTHub.ResourceClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_iothub_enrichment" {
			continue
		}

		enrichmentKey := rs.Primary.Attributes["key"]
		iothubName := rs.Primary.Attributes["iothub_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		iothub, err := client.Get(ctx, resourceGroup, iothubName)
		if err != nil {
			if utils.ResponseWasNotFound(iothub.Response) {
				return nil
			}

			return fmt.Errorf("Bad: Get on iothubResourceClient: %+v", err)
		}
		if iothub.Properties == nil || iothub.Properties.Routing == nil {
			return nil
		}
		enrichments := iothub.Properties.Routing.Enrichments

		if enrichments == nil {
			return nil
		}

		for _, enrichment := range *enrichments {
			if strings.EqualFold(*enrichment.Key, enrichmentKey) {
				return fmt.Errorf("Bad: enrichment %s still exists on IoTHb %s", enrichmentKey, iothubName)
			}
		}
	}
	return nil
}

func testCheckAzureRMIotHubEnrichmentExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).IoTHub.ResourceClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		parsedIothubId, err := azure.ParseAzureResourceID(rs.Primary.ID)
		if err != nil {
			return err
		}
		iothubName := parsedIothubId.Path["IotHubs"]
		enrichmentKey := parsedIothubId.Path["Enrichments"]
		resourceGroup := parsedIothubId.ResourceGroup

		iothub, err := client.Get(ctx, resourceGroup, iothubName)
		if err != nil {
			if utils.ResponseWasNotFound(iothub.Response) {
				return fmt.Errorf("IotHub %q (Resource Group %q) was not found", iothubName, resourceGroup)
			}

			return fmt.Errorf("Error loading IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
		}

		if iothub.Properties == nil || iothub.Properties.Routing == nil {
			return fmt.Errorf("Bad: No Enrichment %s defined for IotHub %s", enrichmentKey, iothubName)
		}
		enrichments := iothub.Properties.Routing.Enrichments

		if enrichments == nil {
			return fmt.Errorf("Bad: No enrichment %s defined for IotHub %s", enrichmentKey, iothubName)
		}

		for _, enrichment := range *enrichments {
			if strings.EqualFold(*enrichment.Key, enrichmentKey) {
				return nil
			}
		}

		return fmt.Errorf("Bad: No enrichment %s defined for IotHub %s", enrichmentKey, iothubName)
	}
}

func testAccAzureRMIotHubEnrichment_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMIotHubEnrichment_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_enrichment" "import" {
  resource_group_name = azurerm_resource_group.test.name
  iothub_name         = azurerm_iothub.test.name
  key                 = "acctest"

  value          = "$twin.tags.DeviceType"
  endpoint_names = [azurerm_iothub_endpoint_storage_container.test.name]
}
`, template)
}

func testAccAzureRMIotHubEnrichment_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iothub-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "test%[1]d"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }

  tags = {
    purpose = "testing"
  }
}

resource "azurerm_iothub_endpoint_storage_container" "test" {
  resource_group_name = azurerm_resource_group.test.name
  iothub_name         = azurerm_iothub.test.name
  name                = "acctest"

  connection_string          = azurerm_storage_account.test.primary_blob_connection_string
  batch_frequency_in_seconds = 60
  max_chunk_size_in_bytes    = 10485760
  container_name             = azurerm_storage_container.test.name
  encoding                   = "Avro"
  file_name_format           = "{iothub}/{partition}_{YYYY}_{MM}_{DD}_{HH}_{mm}"
}

resource "azurerm_iothub_enrichment" "test" {
  resource_group_name = azurerm_resource_group.test.name
  iothub_name         = azurerm_iothub.test.name
  key                 = "acctest"

  value          = "$twin.tags.DeviceType"
  endpoint_names = [azurerm_iothub_endpoint_storage_container.test.name]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMIotHubEnrichment_update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iothub-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "test%[1]d"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }

  tags = {
    purpose = "testing"
  }
}

resource "azurerm_iothub_endpoint_storage_container" "test" {
  resource_group_name = azurerm_resource_group.test.name
  iothub_name         = azurerm_iothub.test.name
  name                = "acctest"

  connection_string          = azurerm_storage_account.test.primary_blob_connection_string
  batch_frequency_in_seconds = 60
  max_chunk_size_in_bytes    = 10485760
  container_name             = azurerm_storage_container.test.name
  encoding                   = "Avro"
  file_name_format           = "{iothub}/{partition}_{YYYY}_{MM}_{DD}_{HH}_{mm}"
}

resource "azurerm_iothub_enrichment" "test" {
  resource_group_name = azurerm_resource_group.test.name
  iothub_name         = azurerm_iothub.test.name
  key                 = "acctest"

  value          = "$twin.tags.Tenant"
  endpoint_names = [azurerm_iothub_endpoint_storage_container.test.name]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

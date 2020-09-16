package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parsers"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMStorageSync_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_sync", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageSyncDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageSync_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageSyncExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageSync_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_sync", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageSyncDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageSync_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageSyncExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMStorageSync_requiresImport),
		},
	})
}

func TestAccAzureRMStorageSync_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_sync", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageSyncDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageSync_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageSyncExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageSync_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_sync", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageSyncDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageSync_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageSyncExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageSync_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageSyncExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMStorageSync_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageSyncExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMStorageSyncExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("storage Sync Service not found: %s", resourceName)
		}

		id, err := parsers.ParseStorageSyncID(rs.Primary.ID)
		if err != nil {
			return err
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Storage.SyncServiceClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Storage Sync Service (Storage Sync Service Name %q / Resource Group %q) does not exist", id.Name, id.ResourceGroup)
			}
			return fmt.Errorf("bad: Get on StorageSyncsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMStorageSyncDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Storage.SyncServiceClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_storage_sync" {
			continue
		}

		id, err := parsers.ParseStorageSyncID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on StorageSyncsClient: %+v", err)
			}
		}

		return nil
	}
	return nil
}

func testAccAzureRMStorageSync_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ss-%d"
  location = "%s"
}

resource "azurerm_storage_sync" "test" {
  name                = "acctest-SS-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  tags = {
    ENV = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMStorageSync_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMStorageSync_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_sync" "import" {
  name                = azurerm_storage_sync.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, template)
}

func testAccAzureRMStorageSync_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_sync" "test" {
  name                    = "acctest-SS-%d"
  resource_group_name     = azurerm_resource_group.test.name
  location                = azurerm_resource_group.test.location
  incoming_traffic_policy = "AllowVirtualNetworksOnly"
  tags = {
    ENV = "Staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

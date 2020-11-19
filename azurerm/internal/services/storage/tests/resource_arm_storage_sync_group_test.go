package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMStorageSyncGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_sync_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageSyncGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageSyncGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageSyncGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageSyncGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_sync_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageSyncGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMStorageSyncGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageSyncGroupExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMStorageSyncGroup_requiresImport),
		},
	})
}

func testCheckAzureRMStorageSyncGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("storage Sync Group not found: %s", resourceName)
		}
		id, err := parse.StorageSyncGroupID(rs.Primary.ID)
		if err != nil {
			return err
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Storage.SyncGroupsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		if resp, err := client.Get(ctx, id.ResourceGroup, id.StorageSyncName, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Storage Sync Group (Storage Sync Group Name %q / Resource Group %q) does not exist", id.Name, id.ResourceGroup)
			}
			return fmt.Errorf("bad: Get on StorageSyncGroupsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMStorageSyncGroupDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Storage.SyncGroupsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_storage_sync_group" {
			continue
		}

		id, err := parse.StorageSyncGroupID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.StorageSyncName, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on StorageSyncGroupsClient: %+v", err)
			}
		}

		return nil
	}
	return nil
}

func testAccAzureRMStorageSyncGroup_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-SS-%[1]d"
  location = "%s"
}

resource "azurerm_storage_sync" "test" {
  name                = "acctest-StorageSync-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_storage_sync_group" "test" {
  name            = "acctest-StorageSyncGroup-%[1]d"
  storage_sync_id = azurerm_storage_sync.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMStorageSyncGroup_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMStorageSyncGroup_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_sync_group" "import" {
  name            = azurerm_storage_sync_group.test.name
  storage_sync_id = azurerm_storage_sync.test.id
}
`, template)
}

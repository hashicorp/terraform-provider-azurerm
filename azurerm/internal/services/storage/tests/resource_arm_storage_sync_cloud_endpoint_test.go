package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parsers"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMStorageSyncCloudEndpoint_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_sync_cloud_endpoint", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageSyncCloudEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCloudEndpoint_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageSyncCloudEndpointExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMStorageSyncCloudEndpoint_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_sync_cloud_endpoint", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMStorageSyncCloudEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCloudEndpoint_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageSyncCloudEndpointExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMCloudEndpoint_requiresImport),
		},
	})
}

func testCheckAzureRMStorageSyncCloudEndpointExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("storage Sync Cloud Endpoint not found: %s", resourceName)
		}

		id, err := parsers.SyncCloudEndpointID(rs.Primary.Attributes["id"])
		if err != nil {
			return err
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Storage.CloudEndpointsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		if resp, err := client.Get(ctx, id.ResourceGroup, id.StorageSyncName, id.StorageSyncGroup, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Storage Sync Cloud Endpoint %q (Storage Sync Group %q / Storage Sync Service Name %q / Resource Group %q) does not exist", id.Name, id.ResourceGroup, id.StorageSyncName, id.ResourceGroup)
			}
			return fmt.Errorf("bad: Get on CloudEndpointsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMStorageSyncCloudEndpointDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Storage.CloudEndpointsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_storage_sync_cloud_endpoint" {
			continue
		}

		id, err := parsers.SyncCloudEndpointID(rs.Primary.Attributes["id"])
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.StorageSyncName, id.StorageSyncGroup, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on CloudEndpointsClient: %+v", err)
			}
		}

		return nil
	}
	return nil
}

func testAccAzureRMCloudEndpoint_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-StorageSync-%[1]d"
  location = "%[2]s"
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

resource "azurerm_storage_account" "test" {
  name                     = "accstr%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_share" "test" {
  name                 = "acctest-share-%[1]d"
  storage_account_name = azurerm_storage_account.test.name

  acl {
    id = "GhostedRecall"
    access_policy {
      permissions = "r"
    }
  }
}

data "azurerm_subscription" "primary" {
}

data "azuread_service_principal" "test" {
  display_name = "Microsoft.StorageSync"
}

resource "azurerm_role_assignment" "test" {
  scope                = data.azurerm_subscription.primary.id
  role_definition_name = "Reader and Data Access"
  principal_id         = data.azuread_service_principal.test.object_id
}

resource "azurerm_storage_sync_cloud_endpoint" "test" {
  name                      = "acctest-CEP-%[1]d"
  storage_sync_group_id     = azurerm_storage_sync_group.test.id
  storage_account_id        = azurerm_storage_account.test.id
  storage_account_tenant_id = "%[4]s"
  file_share_name           = azurerm_storage_share.test.name

  depends_on = [azurerm_role_assignment.test]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, os.Getenv("ARM_TENANT_ID"))
}

func testAccAzureRMCloudEndpoint_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMCloudEndpoint_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_sync_cloud_endpoint" "import" {
  name                      = azurerm_storage_sync_cloud_endpoint.test.name
  storage_sync_group_id     = azurerm_storage_sync_cloud_endpoint.test.storage_sync_group_id
  storage_account_id        = azurerm_storage_sync_cloud_endpoint.test.storage_account_id
  storage_account_tenant_id = azurerm_storage_sync_cloud_endpoint.test.storage_account_tenant_id
  file_share_name           = azurerm_storage_sync_cloud_endpoint.test.file_share_name
}
`, template)
}

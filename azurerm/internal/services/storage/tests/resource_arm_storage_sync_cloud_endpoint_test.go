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

func TestAccAzureRMCloudEndpoint_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_sync_cloud_endpoint", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCloudEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCloudEndpoint_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCloudEndpointExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCloudEndpoint_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_sync_cloud_endpoint", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCloudEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCloudEndpoint_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCloudEndpointExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMCloudEndpoint_requiresImport),
		},
	})
}

func testCheckAzureRMCloudEndpointExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("storage Sync Service not found: %s", resourceName)
		}

		id, err := parsers.CloudEndpointID(rs.Primary.Attributes["id"])
		if err != nil {
			return err
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Storage.CloudEndpointsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		if resp, err := client.Get(ctx, id.ResourceGroup, id.StorageSyncName, id.StorageSyncGroup, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Storage Sync Service %q (Storage Sync Group %q / Storage Sync Service Name %q / Resource Group %q) does not exist", id.Name, id.ResourceGroup, id.StorageSyncName, id.ResourceGroup)
			}
			return fmt.Errorf("bad: Get on CloudEndpointsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMCloudEndpointDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Storage.CloudEndpointsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_storage_sync_cloud_endpoint" {
			continue
		}

		id, err := parsers.CloudEndpointID(rs.Primary.Attributes["id"])
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
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_sync" "test" {
  name                = "acctest-storagesync-%[1]d"
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
  lifecycle {
    ignore_changes = [
      acl, metadata
    ]
  }
}

resource "azurerm_storage_sync_cloud_endpoint" "test" {
  name                  = "acctest-cep-%[1]d"
  storage_sync_group_id = azurerm_storage_sync_group.test.id
  storage_account_id    = azurerm_storage_account.test.id
  file_share_name       = azurerm_storage_share.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccAzureRMCloudEndpoint_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMCloudEndpoint_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_sync_cloud_endpoint" "import" {
  name                  = azurerm_storage_sync_cloud_endpoint.test.name
  storage_sync_group_id = azurerm_storage_sync_group.test.id
  storage_account_id    = azurerm_storage_account.test.id
  file_share_name       = azurerm_storage_share.test.name
}
`, template)
}

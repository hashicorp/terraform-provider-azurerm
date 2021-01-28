package storage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type StorageSyncCloudEndpointResource struct{}

func TestAccAzureRMStorageSyncCloudEndpoint_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_sync_cloud_endpoint", "test")
	r := StorageSyncCloudEndpointResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMStorageSyncCloudEndpoint_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_sync_cloud_endpoint", "test")
	r := StorageSyncCloudEndpointResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMStorageSyncCloudEndpoint_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_sync_cloud_endpoint", "test")
	r := StorageSyncCloudEndpointResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r StorageSyncCloudEndpointResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.StorageSyncCloudEndpointID(state.Attributes["id"])
	if err != nil {
		return nil, err
	}

	resp, err := client.Storage.CloudEndpointsClient.Get(ctx, id.ResourceGroup, id.StorageSyncServiceName, id.SyncGroupName, id.CloudEndpointName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("bad: Get on CloudEndpointsClient: %+v", err)
	}

	return utils.Bool(resp.CloudEndpointProperties != nil), nil
}

func (r StorageSyncCloudEndpointResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_sync_cloud_endpoint" "test" {
  name                  = "acctest-CEP-%d"
  storage_sync_group_id = azurerm_storage_sync_group.test.id
  storage_account_id    = azurerm_storage_account.test.id
  file_share_name       = azurerm_storage_share.test.name
}
`, template, data.RandomInteger)
}

func (r StorageSyncCloudEndpointResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_sync_cloud_endpoint" "test" {
  name                      = "acctest-CEP-%d"
  storage_sync_group_id     = azurerm_storage_sync_group.test.id
  storage_account_id        = azurerm_storage_account.test.id
  storage_account_tenant_id = "%s"
  file_share_name           = azurerm_storage_share.test.name
}
`, template, data.RandomInteger, data.Client().TenantID)
}

func (r StorageSyncCloudEndpointResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_sync_cloud_endpoint" "import" {
  name                  = azurerm_storage_sync_cloud_endpoint.test.name
  storage_sync_group_id = azurerm_storage_sync_cloud_endpoint.test.storage_sync_group_id
  storage_account_id    = azurerm_storage_sync_cloud_endpoint.test.storage_account_id
  file_share_name       = azurerm_storage_sync_cloud_endpoint.test.file_share_name
}
`, template)
}

func (r StorageSyncCloudEndpointResource) template(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

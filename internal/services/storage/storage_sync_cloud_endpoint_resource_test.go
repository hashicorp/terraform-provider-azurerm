// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/cloudendpointresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type StorageSyncCloudEndpointResource struct{}

func TestAccAzureRMStorageSyncCloudEndpoint_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_sync_cloud_endpoint", "test")
	r := StorageSyncCloudEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMStorageSyncCloudEndpoint_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_sync_cloud_endpoint", "test")
	r := StorageSyncCloudEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMStorageSyncCloudEndpoint_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_sync_cloud_endpoint", "test")
	r := StorageSyncCloudEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r StorageSyncCloudEndpointResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := cloudendpointresource.ParseCloudEndpointID(state.Attributes["id"])
	if err != nil {
		return nil, err
	}

	resp, err := client.Storage.SyncCloudEndpointsClient.CloudEndpointsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
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
  quota                = 1

  acl {
    id = "GhostedRecall"
    access_policy {
      permissions = "r"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

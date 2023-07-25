// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datashare_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/dataset"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DataShareDataSetBlobStorageResource struct{}

func TestAccDataShareDataSetBlobStorage_basicFile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share_dataset_blob_storage", "test")
	r := DataShareDataSetBlobStorageResource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basicFile(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("display_name").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataShareDataSetBlobStorage_basicFolder(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share_dataset_blob_storage", "test")
	r := DataShareDataSetBlobStorageResource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basicFolder(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("display_name").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataShareDataSetBlobStorage_basicContainer(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share_dataset_blob_storage", "test")
	r := DataShareDataSetBlobStorageResource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basicContainer(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("display_name").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataShareDataSetBlobStorage_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share_dataset_blob_storage", "test")
	r := DataShareDataSetBlobStorageResource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basicFile(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (t DataShareDataSetBlobStorageResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := dataset.ParseDataSetID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DataShare.DataSetClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if model := resp.Model; model != nil {
		ds := *model
		if _, ok := ds.(dataset.BlobDataSet); ok {
			return utils.Bool(true), nil
		}
		if _, ok := ds.(dataset.BlobFolderDataSet); ok {
			return utils.Bool(true), nil
		}
		if _, ok := ds.(dataset.BlobContainerDataSet); ok {
			return utils.Bool(true), nil
		}
	}

	return nil, fmt.Errorf("%s is not a blob storage dataset", *id)
}

func (DataShareDataSetBlobStorageResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-datashare-%[1]d"
  location = "%[2]s"
}

resource "azurerm_data_share_account" "test" {
  name                = "acctest-DSA-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_data_share" "test" {
  name       = "acctest_DS_%[1]d"
  account_id = azurerm_data_share_account.test.id
  kind       = "CopyBased"
}

resource "azurerm_storage_account" "test" {
  name                            = "acctest%[3]s"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  account_tier                    = "Standard"
  account_replication_type        = "RAGRS"
  allow_nested_items_to_be_public = true
}

resource "azurerm_storage_container" "test" {
  name                  = "acctest-sc-%[1]d"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "container"
}

data "azuread_service_principal" "test" {
  display_name = azurerm_data_share_account.test.name
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Storage Blob Data Reader"
  principal_id         = data.azuread_service_principal.test.object_id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r DataShareDataSetBlobStorageResource) basicFile(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_data_share_dataset_blob_storage" "test" {
  name           = "acctest-DSDSBS-file-%[2]d"
  data_share_id  = azurerm_data_share.test.id
  container_name = azurerm_storage_container.test.name
  storage_account {
    name                = azurerm_storage_account.test.name
    resource_group_name = azurerm_storage_account.test.resource_group_name
    subscription_id     = "%[3]s"
  }
  file_path = "myfile.txt"
  depends_on = [
    azurerm_role_assignment.test,
  ]
}
`, r.template(data), data.RandomInteger, os.Getenv("ARM_SUBSCRIPTION_ID"))
}

func (r DataShareDataSetBlobStorageResource) basicFolder(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_data_share_dataset_blob_storage" "test" {
  name           = "acctest-DSDSBS-folder-%[2]d"
  data_share_id  = azurerm_data_share.test.id
  container_name = azurerm_storage_container.test.name
  storage_account {
    name                = azurerm_storage_account.test.name
    resource_group_name = azurerm_storage_account.test.resource_group_name
    subscription_id     = "%[3]s"
  }
  folder_path = "test"
  depends_on = [
    azurerm_role_assignment.test,
  ]
}
`, r.template(data), data.RandomInteger, os.Getenv("ARM_SUBSCRIPTION_ID"))
}

func (r DataShareDataSetBlobStorageResource) basicContainer(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_data_share_dataset_blob_storage" "test" {
  name           = "acctest-DSDSBS-folder-%[2]d"
  data_share_id  = azurerm_data_share.test.id
  container_name = azurerm_storage_container.test.name
  storage_account {
    name                = azurerm_storage_account.test.name
    resource_group_name = azurerm_storage_account.test.resource_group_name
    subscription_id     = "%[3]s"
  }
  depends_on = [
    azurerm_role_assignment.test,
  ]
}
`, r.template(data), data.RandomInteger, os.Getenv("ARM_SUBSCRIPTION_ID"))
}

func (r DataShareDataSetBlobStorageResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_share_dataset_blob_storage" "import" {
  name           = azurerm_data_share_dataset_blob_storage.test.name
  data_share_id  = azurerm_data_share.test.id
  container_name = azurerm_data_share_dataset_blob_storage.test.container_name
  storage_account {
    name                = azurerm_data_share_dataset_blob_storage.test.storage_account.0.name
    resource_group_name = azurerm_data_share_dataset_blob_storage.test.storage_account.0.resource_group_name
    subscription_id     = azurerm_data_share_dataset_blob_storage.test.storage_account.0.subscription_id
  }
}
`, r.basicFile(data))
}

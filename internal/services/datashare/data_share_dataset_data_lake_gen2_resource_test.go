// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datashare_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/dataset"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DataShareDataSetDataLakeGen2Resource struct{}

func TestAccDataShareDataSetDataLakeGen2_basicFile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share_dataset_data_lake_gen2", "test")
	r := DataShareDataSetDataLakeGen2Resource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
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

func TestAccDataShareDataSetDataLakeGen2_basicFolder(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share_dataset_data_lake_gen2", "test")
	r := DataShareDataSetDataLakeGen2Resource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
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

func TestAccDataShareDataSetDataLakeGen2File_basicSystem(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share_dataset_data_lake_gen2", "test")
	r := DataShareDataSetDataLakeGen2Resource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicSystem(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("display_name").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataShareDataLakeGen2DataSet_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share_dataset_data_lake_gen2", "test")
	r := DataShareDataSetDataLakeGen2Resource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicFile(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (t DataShareDataSetDataLakeGen2Resource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := dataset.ParseDataSetID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DataShare.DataSetClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if model := resp.Model; model != nil {
		if _, ok := model.(dataset.ADLSGen2FileDataSet); ok {
			return utils.Bool(true), nil
		}
		if _, ok := model.(dataset.ADLSGen2FolderDataSet); ok {
			return utils.Bool(true), nil
		}
		if _, ok := model.(dataset.ADLSGen2FileSystemDataSet); ok {
			return utils.Bool(true), nil
		}
	}

	return nil, fmt.Errorf("%s is not a datalake store gen2 dataset", *id)
}

func (DataShareDataSetDataLakeGen2Resource) template(data acceptance.TestData) string {
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
  name                = "acctest-dsa-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_data_share" "test" {
  name       = "acctest_ds_%[1]d"
  account_id = azurerm_data_share_account.test.id
  kind       = "CopyBased"
}

resource "azurerm_storage_account" "test" {
  name                     = "accteststr%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_kind             = "BlobStorage"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_data_lake_gen2_filesystem" "test" {
  name               = "acctest-%[1]d"
  storage_account_id = azurerm_storage_account.test.id
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

func (r DataShareDataSetDataLakeGen2Resource) basicFile(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_share_dataset_data_lake_gen2" "test" {
  name               = "acctest-dlds-%d"
  share_id           = azurerm_data_share.test.id
  storage_account_id = azurerm_storage_account.test.id
  file_system_name   = azurerm_storage_data_lake_gen2_filesystem.test.name
  file_path          = "myfile.txt"
  depends_on = [
    azurerm_role_assignment.test,
  ]
}
`, r.template(data), data.RandomInteger)
}

func (r DataShareDataSetDataLakeGen2Resource) basicFolder(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_share_dataset_data_lake_gen2" "test" {
  name               = "acctest-dlds-%d"
  share_id           = azurerm_data_share.test.id
  storage_account_id = azurerm_storage_account.test.id
  file_system_name   = azurerm_storage_data_lake_gen2_filesystem.test.name
  folder_path        = "test"
  depends_on = [
    azurerm_role_assignment.test,
  ]
}
`, r.template(data), data.RandomInteger)
}

func (r DataShareDataSetDataLakeGen2Resource) basicSystem(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_share_dataset_data_lake_gen2" "test" {
  name               = "acctest-dlds-%d"
  share_id           = azurerm_data_share.test.id
  storage_account_id = azurerm_storage_account.test.id
  file_system_name   = azurerm_storage_data_lake_gen2_filesystem.test.name
  depends_on = [
    azurerm_role_assignment.test,
  ]
}
`, r.template(data), data.RandomInteger)
}

func (r DataShareDataSetDataLakeGen2Resource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_share_dataset_data_lake_gen2" "import" {
  name               = azurerm_data_share_dataset_data_lake_gen2.test.name
  share_id           = azurerm_data_share.test.id
  storage_account_id = azurerm_data_share_dataset_data_lake_gen2.test.storage_account_id
  file_system_name   = azurerm_data_share_dataset_data_lake_gen2.test.file_system_name
  file_path          = azurerm_data_share_dataset_data_lake_gen2.test.file_path
}
`, r.basicFile(data))
}

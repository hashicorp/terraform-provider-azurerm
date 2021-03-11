package datashare_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/datashare/mgmt/2019-11-01/datashare"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datashare/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type DataShareDataSetBlobStorageResource struct {
}

func TestAccDataShareDataSetBlobStorage_basicFile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share_dataset_blob_storage", "test")
	r := DataShareDataSetBlobStorageResource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basicFile(data),
			Check: resource.ComposeTestCheckFunc(
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

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basicFolder(data),
			Check: resource.ComposeTestCheckFunc(
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

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basicContainer(data),
			Check: resource.ComposeTestCheckFunc(
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

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basicFile(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (t DataShareDataSetBlobStorageResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.DataSetID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DataShare.DataSetClient.Get(ctx, id.ResourceGroup, id.AccountName, id.ShareName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Data Share Data Set %q (resource group: %q): %+v", id.Name, id.ResourceGroup, err)
	}

	switch resp := resp.Value.(type) {
	case datashare.BlobDataSet:
		return utils.Bool(resp.BlobProperties != nil), nil

	case datashare.BlobFolderDataSet:
		return utils.Bool(resp.BlobFolderProperties != nil), nil

	case datashare.BlobContainerDataSet:
		return utils.Bool(resp.BlobContainerProperties != nil), nil
	}

	return nil, fmt.Errorf("Data Share Data %q (Resource Group %q / accountName %q / shareName %q) is not a datalake store gen2 dataset", id.Name, id.ResourceGroup, id.AccountName, id.ShareName)
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
  name                     = "acctest%[3]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "RAGRS"
  allow_blob_public_access = true
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
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(12))
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

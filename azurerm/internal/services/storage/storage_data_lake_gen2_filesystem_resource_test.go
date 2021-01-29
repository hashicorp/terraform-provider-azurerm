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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/datalakestore/filesystems"
)

type StorageDataLakeGen2FileSystemResource struct{}

func TestAccStorageDataLakeGen2FileSystem_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_data_lake_gen2_filesystem", "test")
	r := StorageDataLakeGen2FileSystemResource{}

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

func TestAccStorageDataLakeGen2FileSystem_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_data_lake_gen2_filesystem", "test")
	r := StorageDataLakeGen2FileSystemResource{}

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

func TestAccStorageDataLakeGen2FileSystem_withDefaultACL(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_data_lake_gen2_filesystem", "test")
	r := StorageDataLakeGen2FileSystemResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withDefaultACL(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccStorageDataLakeGen2FileSystem_UpdateDefaultACL(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_data_lake_gen2_filesystem", "test")
	r := StorageDataLakeGen2FileSystemResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withDefaultACL(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withExecuteACLForSPN(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageDataLakeGen2FileSystem_properties(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_data_lake_gen2_filesystem", "test")
	r := StorageDataLakeGen2FileSystemResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.properties(data, "aGVsbG8="),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.properties(data, "ZXll"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageDataLakeGen2FileSystem_handlesStorageAccountDeletion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_data_lake_gen2_filesystem", "test")
	r := StorageDataLakeGen2FileSystemResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func (r StorageDataLakeGen2FileSystemResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := filesystems.ParseResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Storage.FileSystemsClient.GetProperties(ctx, id.AccountName, id.DirectoryName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving File System %q (Account %q): %+v", id.DirectoryName, id.AccountName, err)
	}
	return utils.Bool(true), nil
}

func (r StorageDataLakeGen2FileSystemResource) Destroy(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := filesystems.ParseResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	if _, err := client.Storage.FileSystemsClient.Delete(ctx, id.AccountName, id.DirectoryName); err != nil {
		return nil, fmt.Errorf("deleting File System %q (Account %q): %+v", id.DirectoryName, id.AccountName, err)
	}
	return utils.Bool(true), nil
}

func (r StorageDataLakeGen2FileSystemResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_data_lake_gen2_filesystem" "test" {
  name               = "acctest-%d"
  storage_account_id = azurerm_storage_account.test.id
}
`, template, data.RandomInteger)
}

func (r StorageDataLakeGen2FileSystemResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_data_lake_gen2_filesystem" "import" {
  name               = azurerm_storage_data_lake_gen2_filesystem.test.name
  storage_account_id = azurerm_storage_data_lake_gen2_filesystem.test.storage_account_id
}
`, template)
}

func (r StorageDataLakeGen2FileSystemResource) properties(data acceptance.TestData, value string) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_data_lake_gen2_filesystem" "test" {
  name               = "acctest-%d"
  storage_account_id = azurerm_storage_account.test.id

  properties = {
    key = "%s"
  }
}
`, template, data.RandomInteger, value)
}

func (r StorageDataLakeGen2FileSystemResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_kind             = "BlobStorage"
  account_tier             = "Standard"
  account_replication_type = "LRS"
  is_hns_enabled           = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r StorageDataLakeGen2FileSystemResource) withDefaultACL(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "current" {
}

resource "azurerm_role_assignment" "storageAccountRoleAssignment" {
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = data.azurerm_client_config.current.object_id
}

resource "azurerm_storage_data_lake_gen2_filesystem" "test" {
  name               = "acctest-%[2]d"
  storage_account_id = azurerm_storage_account.test.id
  ace {
    type        = "user"
    permissions = "rwx"
  }
  ace {
    type        = "group"
    permissions = "r-x"
  }
  ace {
    type        = "other"
    permissions = "---"
  }
  depends_on = [
    azurerm_role_assignment.storageAccountRoleAssignment
  ]
}
`, template, data.RandomInteger)
}

func (r StorageDataLakeGen2FileSystemResource) withExecuteACLForSPN(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "current" {
}

resource "azurerm_role_assignment" "storageAccountRoleAssignment" {
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = data.azurerm_client_config.current.object_id
}

resource "azuread_application" "test" {
  name = "acctestspa%[2]d"
}

resource "azuread_service_principal" "test" {
  application_id = azuread_application.test.application_id
}

resource "azurerm_storage_data_lake_gen2_filesystem" "test" {
  name               = "acctest-%[2]d"
  storage_account_id = azurerm_storage_account.test.id
  ace {
    type        = "user"
    permissions = "rwx"
  }
  ace {
    type        = "user"
    id          = azuread_service_principal.test.object_id
    permissions = "--x"
  }
  ace {
    type        = "group"
    permissions = "r-x"
  }
  ace {
    type        = "mask"
    permissions = "r-x"
  }
  ace {
    type        = "other"
    permissions = "---"
  }
  depends_on = [
    azurerm_role_assignment.storageAccountRoleAssignment,
    azuread_service_principal.test
  ]
}
`, template, data.RandomInteger)
}

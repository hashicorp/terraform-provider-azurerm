// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/jackofallops/giovanni/storage/2023-11-03/datalakestore/filesystems"
)

type StorageDataLakeGen2FileSystemResource struct{}

func TestAccStorageDataLakeGen2FileSystem_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_data_lake_gen2_filesystem", "test")
	r := StorageDataLakeGen2FileSystemResource{}

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

func TestAccStorageDataLakeGen2FileSystem_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_data_lake_gen2_filesystem", "test")
	r := StorageDataLakeGen2FileSystemResource{}

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

func TestAccStorageDataLakeGen2FileSystem_withDefaultACL(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_data_lake_gen2_filesystem", "test")
	r := StorageDataLakeGen2FileSystemResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withDefaultACL(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccStorageDataLakeGen2FileSystem_UpdateDefaultACL(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_data_lake_gen2_filesystem", "test")
	r := StorageDataLakeGen2FileSystemResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withDefaultACL(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withExecuteACLForSPN(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageDataLakeGen2FileSystem_encryptionScope(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_data_lake_gen2_filesystem", "test")
	r := StorageDataLakeGen2FileSystemResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.encryptionScope(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageDataLakeGen2FileSystem_properties(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_data_lake_gen2_filesystem", "test")
	r := StorageDataLakeGen2FileSystemResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.properties(data, "aGVsbG8="),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.properties(data, "ZXll"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageDataLakeGen2FileSystem_handlesStorageAccountDeletion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_data_lake_gen2_filesystem", "test")
	r := StorageDataLakeGen2FileSystemResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func TestAccStorageDataLakeGen2FileSystem_withOwnerGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_data_lake_gen2_filesystem", "test")
	r := StorageDataLakeGen2FileSystemResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.withOwnerGroup,
			TestResource: r,
		}),
	})
}

func TestAccStorageDataLakeGen2FileSystem_withSuperUsers(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_data_lake_gen2_filesystem", "test")
	r := StorageDataLakeGen2FileSystemResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withSuperUsers(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r StorageDataLakeGen2FileSystemResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := filesystems.ParseFileSystemID(state.ID, client.Storage.StorageDomainSuffix)
	if err != nil {
		return nil, err
	}

	account, err := client.Storage.FindAccount(ctx, client.Account.SubscriptionId, id.AccountId.AccountName)
	if err != nil {
		return nil, fmt.Errorf("retrieving Account %q for Queue %q: %+v", id.AccountId, id.FileSystemName, err)
	}
	if account == nil {
		return nil, fmt.Errorf("unable to determine Resource Group for Storage Queue %q (Account %q)", id.FileSystemName, id.AccountId.AccountName)
	}

	filesystemsClient, err := client.Storage.DataLakeFilesystemsDataPlaneClient(ctx, *account, client.Storage.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return nil, fmt.Errorf("building Data Lake Gen2 Filesystems Client: %+v", err)
	}

	resp, err := filesystemsClient.GetProperties(ctx, id.FileSystemName)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving File System %q (Account %q): %+v", id.FileSystemName, id.AccountId.AccountName, err)
	}

	return utils.Bool(true), nil
}

func (r StorageDataLakeGen2FileSystemResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := filesystems.ParseFileSystemID(state.ID, client.Storage.StorageDomainSuffix)
	if err != nil {
		return nil, err
	}

	account, err := client.Storage.FindAccount(ctx, client.Account.SubscriptionId, id.AccountId.AccountName)
	if err != nil {
		return nil, fmt.Errorf("retrieving Account %q for Queue %q: %+v", id.AccountId, id.FileSystemName, err)
	}
	if account == nil {
		return nil, fmt.Errorf("unable to determine Resource Group for Storage Queue %q (Account %q)", id.FileSystemName, id.AccountId.AccountName)
	}

	filesystemsClient, err := client.Storage.DataLakeFilesystemsDataPlaneClient(ctx, *account, client.Storage.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return nil, fmt.Errorf("building Data Lake Gen2 Filesystems Client: %+v", err)
	}

	if _, err = filesystemsClient.Delete(ctx, id.FileSystemName); err != nil {
		return nil, fmt.Errorf("deleting File System %q (Account %q): %+v", id.FileSystemName, id.AccountId.AccountName, err)
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

func (r StorageDataLakeGen2FileSystemResource) encryptionScope(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_encryption_scope" "test" {
  name               = "acctestEScontainer%[2]d"
  storage_account_id = azurerm_storage_account.test.id
  source             = "Microsoft.Storage"
}

resource "azurerm_storage_data_lake_gen2_filesystem" "test" {
  name               = "acctest-%[2]d"
  storage_account_id = azurerm_storage_account.test.id

  default_encryption_scope = azurerm_storage_encryption_scope.test.name
}
`, template, data.RandomInteger)
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

provider "azuread" {}

data "azurerm_client_config" "current" {
}

resource "azurerm_role_assignment" "storageAccountRoleAssignment" {
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = data.azurerm_client_config.current.object_id
}

resource "azuread_application" "test" {
  display_name = "acctestspa%[2]d"
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

func (r StorageDataLakeGen2FileSystemResource) withOwnerGroup(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azuread" {}

data "azurerm_client_config" "current" {}

resource "azurerm_role_assignment" "storage_blob_owner" {
  role_definition_name = "Storage Blob Data Owner"
  scope                = azurerm_resource_group.test.id
  principal_id         = data.azurerm_client_config.current.object_id
}

resource "azuread_application" "test" {
  display_name = "acctestspa%[2]d"
}

resource "azuread_service_principal" "test" {
  application_id = azuread_application.test.application_id
}

resource "azurerm_storage_data_lake_gen2_filesystem" "test" {
  name               = "acctest-%[2]d"
  storage_account_id = azurerm_storage_account.test.id
  owner              = azuread_service_principal.test.object_id
  group              = azuread_service_principal.test.object_id
}
`, template, data.RandomInteger)
}

func (r StorageDataLakeGen2FileSystemResource) withSuperUsers(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azuread" {}

resource "azuread_application" "test" {
  display_name = "acctestspa%[2]d"
}

resource "azuread_service_principal" "test" {
  application_id = azuread_application.test.application_id
}

resource "azurerm_storage_data_lake_gen2_filesystem" "test" {
  name               = "acctest-%[2]d"
  storage_account_id = azurerm_storage_account.test.id
  owner              = "$superuser"
  group              = "$superuser"
}
`, template, data.RandomInteger)
}

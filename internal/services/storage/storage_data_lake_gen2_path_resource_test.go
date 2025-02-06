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
	"github.com/jackofallops/giovanni/storage/2023-11-03/datalakestore/paths"
)

type StorageDataLakeGen2PathResource struct{}

func TestAccStorageDataLakeGen2Path_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_data_lake_gen2_path", "test")
	r := StorageDataLakeGen2PathResource{}

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

func TestAccStorageDataLakeGen2Path_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_data_lake_gen2_path", "test")
	r := StorageDataLakeGen2PathResource{}

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

func TestAccStorageDataLakeGen2Path_withSimpleACLAndUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_data_lake_gen2_path", "test")
	r := StorageDataLakeGen2PathResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withSimpleACL(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withSimpleACLUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageDataLakeGen2Path_withSimpleACL(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_data_lake_gen2_path", "test")
	r := StorageDataLakeGen2PathResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withSimpleACL(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageDataLakeGen2Path_withACLWithSpecificUserAndDefaults(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_data_lake_gen2_path", "test")
	r := StorageDataLakeGen2PathResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withACLWithSpecificUserAndDefaults(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageDataLakeGen2Path_withOwner(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_data_lake_gen2_path", "test")
	r := StorageDataLakeGen2PathResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withOwner(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageDataLakeGen2Path_withSuperUsers(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_data_lake_gen2_path", "test")
	r := StorageDataLakeGen2PathResource{}

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

func (r StorageDataLakeGen2PathResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := paths.ParsePathID(state.ID, client.Storage.StorageDomainSuffix)
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

	pathsClient, err := client.Storage.DataLakePathsDataPlaneClient(ctx, *account, client.Storage.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return nil, fmt.Errorf("building Data Lake Gen2 Paths Client: %+v", err)
	}

	resp, err := pathsClient.GetProperties(ctx, id.FileSystemName, id.Path, paths.GetPropertiesInput{Action: paths.GetPropertiesActionGetStatus})
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Path %q (File System %q / Account %q): %+v", id.Path, id.FileSystemName, id.AccountId.AccountName, err)
	}

	return utils.Bool(true), nil
}

func (r StorageDataLakeGen2PathResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_data_lake_gen2_path" "test" {
  storage_account_id = azurerm_storage_account.test.id
  filesystem_name    = azurerm_storage_data_lake_gen2_filesystem.test.name
  path               = "testpath"
  resource           = "directory"
}
`, template)
}

func (r StorageDataLakeGen2PathResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_data_lake_gen2_path" "import" {
  path               = azurerm_storage_data_lake_gen2_path.test.path
  filesystem_name    = azurerm_storage_data_lake_gen2_path.test.filesystem_name
  storage_account_id = azurerm_storage_data_lake_gen2_path.test.storage_account_id
  resource           = azurerm_storage_data_lake_gen2_path.test.resource
}
`, template)
}

func (r StorageDataLakeGen2PathResource) withSimpleACL(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_role_assignment" "storage_blob_owner" {
  role_definition_name = "Storage Blob Data Owner"
  scope                = azurerm_resource_group.test.id
  principal_id         = data.azurerm_client_config.current.object_id
}

resource "azurerm_storage_data_lake_gen2_path" "test" {
  storage_account_id = azurerm_storage_account.test.id
  filesystem_name    = azurerm_storage_data_lake_gen2_filesystem.test.name
  path               = "testpath"
  resource           = "directory"
  ace {
    type        = "user"
    permissions = "r-x"
  }
  ace {
    type        = "group"
    permissions = "-wx"
  }
  ace {
    type        = "other"
    permissions = "--x"
  }
}
`, template)
}

func (r StorageDataLakeGen2PathResource) withSimpleACLUpdated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_role_assignment" "storage_blob_owner" {
  role_definition_name = "Storage Blob Data Owner"
  scope                = azurerm_resource_group.test.id
  principal_id         = data.azurerm_client_config.current.object_id
}

resource "azurerm_storage_data_lake_gen2_path" "test" {
  storage_account_id = azurerm_storage_account.test.id
  filesystem_name    = azurerm_storage_data_lake_gen2_filesystem.test.name
  path               = "testpath"
  resource           = "directory"
  ace {
    type        = "user"
    permissions = "rwx"
  }
  ace {
    type        = "group"
    permissions = "-wx"
  }
  ace {
    type        = "other"
    permissions = "--x"
  }
}
`, template)
}

func (r StorageDataLakeGen2PathResource) withACLWithSpecificUserAndDefaults(data acceptance.TestData) string {
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

resource "azurerm_storage_data_lake_gen2_path" "test" {
  storage_account_id = azurerm_storage_account.test.id
  filesystem_name    = azurerm_storage_data_lake_gen2_filesystem.test.name
  path               = "testpath"
  resource           = "directory"
  ace {
    type        = "user"
    permissions = "r-x"
  }
  ace {
    type        = "user"
    id          = azuread_service_principal.test.object_id
    permissions = "r-x"
  }
  ace {
    type        = "group"
    permissions = "-wx"
  }
  ace {
    type        = "mask"
    permissions = "--x"
  }
  ace {
    type        = "other"
    permissions = "--x"
  }
  ace {
    scope       = "default"
    type        = "user"
    permissions = "r-x"
  }
  ace {
    scope       = "default"
    type        = "user"
    id          = azuread_service_principal.test.object_id
    permissions = "r-x"
  }
  ace {
    scope       = "default"
    type        = "group"
    permissions = "-wx"
  }
  ace {
    scope       = "default"
    type        = "mask"
    permissions = "--x"
  }
  ace {
    scope       = "default"
    type        = "other"
    permissions = "--x"
  }
}
`, template, data.RandomInteger)
}

func (r StorageDataLakeGen2PathResource) withOwner(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azuread" {}

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

resource "azurerm_storage_data_lake_gen2_path" "test" {
  storage_account_id = azurerm_storage_account.test.id
  filesystem_name    = azurerm_storage_data_lake_gen2_filesystem.test.name
  path               = "testpath"
  resource           = "directory"
  owner              = azuread_service_principal.test.object_id
}
`, template, data.RandomInteger)
}

func (r StorageDataLakeGen2PathResource) withSuperUsers(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azuread" {}

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

resource "azurerm_storage_data_lake_gen2_path" "test" {
  storage_account_id = azurerm_storage_account.test.id
  filesystem_name    = azurerm_storage_data_lake_gen2_filesystem.test.name
  path               = "testpath"
  resource           = "directory"
  owner              = "$superuser"
  group              = "$superuser"
}
`, template, data.RandomInteger)
}

func (r StorageDataLakeGen2PathResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_kind             = "BlobStorage"
  account_tier             = "Standard"
  account_replication_type = "LRS"
  is_hns_enabled           = true
}

data "azurerm_client_config" "current" {
}

resource "azurerm_role_assignment" "storageAccountRoleAssignment" {
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = data.azurerm_client_config.current.object_id
}


resource "azurerm_storage_data_lake_gen2_filesystem" "test" {
  name               = "fstest"
  storage_account_id = azurerm_storage_account.test.id
  depends_on = [
    azurerm_role_assignment.storageAccountRoleAssignment
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

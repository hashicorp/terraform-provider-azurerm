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
	"github.com/jackofallops/giovanni/storage/2023-11-03/file/directories"
)

type StorageShareDirectoryResource struct{}

func TestAccStorageShareDirectory_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_directory", "test")
	r := StorageShareDirectoryResource{}

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

func TestAccStorageShareDirectory_basicAzureADAuth(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_directory", "test")
	r := StorageShareDirectoryResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicAzureADAuth(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageShareDirectory_uppercase(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_directory", "test")
	r := StorageShareDirectoryResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.uppercase(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageShareDirectory_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_directory", "test")
	r := StorageShareDirectoryResource{}

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

func TestAccStorageShareDirectory_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_directory", "test")
	r := StorageShareDirectoryResource{}

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

func TestAccStorageShareDirectory_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_directory", "test")
	r := StorageShareDirectoryResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageShareDirectory_nested(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_directory", "parent")
	r := StorageShareDirectoryResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.nested(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That("azurerm_storage_share_directory.child_one").ExistsInAzure(r),
				check.That("azurerm_storage_share_directory.child_two").ExistsInAzure(r),
				check.That("azurerm_storage_share_directory.multiple_child_one").ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageShareDirectory_nestedWithBackslashes(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_directory", "dos")
	r := StorageShareDirectoryResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.nestedWithBackslashes(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That("azurerm_storage_share_directory.c").ExistsInAzure(r),
				check.That("azurerm_storage_share_directory.dos").ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r StorageShareDirectoryResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := directories.ParseDirectoryID(state.ID, client.Storage.StorageDomainSuffix)
	if err != nil {
		return nil, err
	}
	account, err := client.Storage.FindAccount(ctx, client.Account.SubscriptionId, id.AccountId.AccountName)
	if err != nil {
		return nil, fmt.Errorf("retrieving Account %q for Directory %q (Share %q): %s", id.AccountId.AccountName, id.DirectoryPath, id.ShareName, err)
	}
	if account == nil {
		return nil, fmt.Errorf("unable to determine Resource Group for Storage Share Directory %q (Share %q / Account %q)", id.DirectoryPath, id.ShareName, id.AccountId.AccountName)
	}
	dirClient, err := client.Storage.FileShareDirectoriesDataPlaneClient(ctx, *account, client.Storage.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return nil, fmt.Errorf("building File Share client for %s: %+v", account.StorageAccountId, err)
	}
	resp, err := dirClient.Get(ctx, id.ShareName, id.DirectoryPath)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Storage Share %q (File Share %q in %s): %+v", id.DirectoryPath, id.ShareName, account.StorageAccountId, err)
	}
	return utils.Bool(true), nil
}

func (r StorageShareDirectoryResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_directory" "test" {
  name             = "dir"
  storage_share_id = azurerm_storage_share.test.id
}
`, template)
}

func (r StorageShareDirectoryResource) basicAzureADAuth(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  storage_use_azuread = true
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_share" "test" {
  name                 = "fileshare"
  storage_account_name = azurerm_storage_account.test.name
  quota                = 50
}

resource "azurerm_storage_share_directory" "test" {
  name             = "dir"
  storage_share_id = azurerm_storage_share.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r StorageShareDirectoryResource) uppercase(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_directory" "test" {
  name             = "UpperCaseCharacterS"
  storage_share_id = azurerm_storage_share.test.id
}
`, template)
}

func (r StorageShareDirectoryResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_directory" "import" {
  name             = azurerm_storage_share_directory.test.name
  storage_share_id = azurerm_storage_share.test.id
}
`, template)
}

func (r StorageShareDirectoryResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_directory" "test" {
  name             = "dir"
  storage_share_id = azurerm_storage_share.test.id

  metadata = {
    hello = "world"
  }
}
`, template)
}

func (r StorageShareDirectoryResource) updated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_directory" "test" {
  name             = "dir"
  storage_share_id = azurerm_storage_share.test.id

  metadata = {
    hello    = "world"
    sunshine = "at dawn"
  }
}
`, template)
}

func (r StorageShareDirectoryResource) nested(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_directory" "parent" {
  name             = "123--parent-dir"
  storage_share_id = azurerm_storage_share.test.id
}

resource "azurerm_storage_share_directory" "child_one" {
  name             = "${azurerm_storage_share_directory.parent.name}/child1"
  storage_share_id = azurerm_storage_share.test.id
}

resource "azurerm_storage_share_directory" "child_two" {
  name             = "${azurerm_storage_share_directory.child_one.name}/childtwo--123"
  storage_share_id = azurerm_storage_share.test.id
}

resource "azurerm_storage_share_directory" "multiple_child_one" {
  name             = "${azurerm_storage_share_directory.parent.name}/c"
  storage_share_id = azurerm_storage_share.test.id
}
`, template)
}

func (r StorageShareDirectoryResource) nestedWithBackslashes(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_directory" "c" {
  name             = "c"
  storage_share_id = azurerm_storage_share.test.id
}

resource "azurerm_storage_share_directory" "dos" {
  name             = "c\\dos"
  storage_share_id = azurerm_storage_share.test.id
  depends_on       = [azurerm_storage_share_directory.c]
}
`, template)
}

func (r StorageShareDirectoryResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_share" "test" {
  name                 = "fileshare"
  storage_account_name = azurerm_storage_account.test.name
  quota                = 50
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

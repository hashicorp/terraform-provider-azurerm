// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/jackofallops/giovanni/storage/2023-11-03/file/directories"
)

type StorageShareDirectoryResourceDeprecated struct{}

func TestAccStorageShareDirectory_basic_deprecated(t *testing.T) {
	if features.FivePointOh() {
		t.Skip("skipping as test is not valid in 5.0")
	}

	data := acceptance.BuildTestData(t, "azurerm_storage_share_directory", "test")
	r := StorageShareDirectoryResourceDeprecated{}

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

func TestAccStorageShareDirectory_basicAzureADAuth_deprecated(t *testing.T) {
	if features.FivePointOh() {
		t.Skip("skipping as test is not valid in 5.0")
	}

	data := acceptance.BuildTestData(t, "azurerm_storage_share_directory", "test")
	r := StorageShareDirectoryResourceDeprecated{}

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

func TestAccStorageShareDirectory_uppercase_deprecated(t *testing.T) {
	if features.FivePointOh() {
		t.Skip("skipping as test is not valid in 5.0")
	}

	data := acceptance.BuildTestData(t, "azurerm_storage_share_directory", "test")
	r := StorageShareDirectoryResourceDeprecated{}

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

func TestAccStorageShareDirectory_requiresImport_deprecated(t *testing.T) {
	if features.FivePointOh() {
		t.Skip("skipping as test is not valid in 5.0")
	}

	data := acceptance.BuildTestData(t, "azurerm_storage_share_directory", "test")
	r := StorageShareDirectoryResourceDeprecated{}

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

func TestAccStorageShareDirectory_complete_deprecated(t *testing.T) {
	if features.FivePointOh() {
		t.Skip("skipping as test is not valid in 5.0")
	}

	data := acceptance.BuildTestData(t, "azurerm_storage_share_directory", "test")
	r := StorageShareDirectoryResourceDeprecated{}

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

func TestAccStorageShareDirectory_update_deprecated(t *testing.T) {
	if features.FivePointOh() {
		t.Skip("skipping as test is not valid in 5.0")
	}

	data := acceptance.BuildTestData(t, "azurerm_storage_share_directory", "test")
	r := StorageShareDirectoryResourceDeprecated{}

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

func TestAccStorageShareDirectory_nested_deprecated(t *testing.T) {
	if features.FivePointOh() {
		t.Skip("skipping as test is not valid in 5.0")
	}

	data := acceptance.BuildTestData(t, "azurerm_storage_share_directory", "parent")
	r := StorageShareDirectoryResourceDeprecated{}

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

func TestAccStorageShareDirectory_nestedWithBackslashes_deprecated(t *testing.T) {
	if features.FivePointOh() {
		t.Skip("skipping as test is not valid in 5.0")
	}

	data := acceptance.BuildTestData(t, "azurerm_storage_share_directory", "dos")
	r := StorageShareDirectoryResourceDeprecated{}

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

func (r StorageShareDirectoryResourceDeprecated) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving Storage Share %q (File Share %q in %s): %+v", id.DirectoryPath, id.ShareName, account.StorageAccountId, err)
	}
	return pointer.To(true), nil
}

func (r StorageShareDirectoryResourceDeprecated) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_directory" "test" {
  name             = "dir"
  storage_share_id = azurerm_storage_share.test.id
}
`, template)
}

func (r StorageShareDirectoryResourceDeprecated) basicAzureADAuth(data acceptance.TestData) string {
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

func (r StorageShareDirectoryResourceDeprecated) uppercase(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_directory" "test" {
  name             = "UpperCaseCharacterS"
  storage_share_id = azurerm_storage_share.test.id
}
`, template)
}

func (r StorageShareDirectoryResourceDeprecated) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_directory" "import" {
  name             = azurerm_storage_share_directory.test.name
  storage_share_id = azurerm_storage_share.test.id
}
`, template)
}

func (r StorageShareDirectoryResourceDeprecated) complete(data acceptance.TestData) string {
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

func (r StorageShareDirectoryResourceDeprecated) updated(data acceptance.TestData) string {
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

func (r StorageShareDirectoryResourceDeprecated) nested(data acceptance.TestData) string {
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

func (r StorageShareDirectoryResourceDeprecated) nestedWithBackslashes(data acceptance.TestData) string {
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

func (r StorageShareDirectoryResourceDeprecated) template(data acceptance.TestData) string {
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

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/jackofallops/giovanni/storage/2023-11-03/file/files"
)

type StorageShareFileResource struct{}

func TestAccAzureRMStorageShareFile_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_file", "test")
	r := StorageShareFileResource{}

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

func TestAccAzureRMStorageShareFile_basicAzureADAuth(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_file", "test")
	r := StorageShareFileResource{}

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

func TestAccAzureRMStorageShareFile_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_file", "test")
	r := StorageShareFileResource{}

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

func TestAccAzureRMStorageShareFile_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_file", "test")
	r := StorageShareFileResource{}

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

func TestAccAzureRMStorageShareFile_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_file", "test")
	r := StorageShareFileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("content_length").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMStorageShareFile_withFile(t *testing.T) {
	sourceBlob, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("Failed to create local source blob file")
	}

	if err := populateTempFile(sourceBlob); err != nil {
		t.Fatalf("Error populating temp file: %s", err)
	}
	data := acceptance.BuildTestData(t, "azurerm_storage_share_file", "test")
	r := StorageShareFileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withFile(data, sourceBlob.Name()),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("content_length").Exists(),
			),
		},
		data.ImportStep("source"),
	})
}

func TestAccAzureRMStorageShareFile_withEmptyFile(t *testing.T) {
	sourceBlob, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("Failed to create local source blob file")
	}

	data := acceptance.BuildTestData(t, "azurerm_storage_share_file", "test")
	r := StorageShareFileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.withFile(data, sourceBlob.Name()),
			ExpectError: regexp.MustCompile(`Error: file .* is empty`),
		},
	})
}

func TestAccAzureRMStorageShareFile_withPath(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_file", "test")
	r := StorageShareFileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withPath(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMStorageShareFile_withPathUsingBackslashes(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_file", "test")
	r := StorageShareFileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withPathUsingBackslashes(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMStorageShareFile_withPathInNameUsingBackslashes(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_file", "test")
	r := StorageShareFileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withPathInNameUsingBackslashes(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (StorageShareFileResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := files.ParseFileID(state.ID, clients.Storage.StorageDomainSuffix)
	if err != nil {
		return nil, err
	}

	account, err := clients.Storage.FindAccount(ctx, clients.Account.SubscriptionId, id.AccountId.AccountName)
	if err != nil {
		return nil, fmt.Errorf("retrieving Account %q for File %q (Share %q): %s", id.AccountId.AccountName, id.FileName, id.ShareName, err)
	}
	if account == nil {
		return utils.Bool(false), nil
	}

	client, err := clients.Storage.FileShareFilesDataPlaneClient(ctx, *account, clients.Storage.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return nil, fmt.Errorf("building File Share Files Client: %s", err)
	}

	resp, err := client.GetProperties(ctx, id.ShareName, id.DirectoryPath, id.FileName)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return nil, fmt.Errorf("checking for presence of existing File %q (File Share %q in %s): %+v", id.FileName, id.ShareName, account.StorageAccountId, err)
		}
	}

	return utils.Bool(true), nil
}

func (StorageShareFileResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
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

func (r StorageShareFileResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_file" "test" {
  name             = "file"
  storage_share_id = azurerm_storage_share.test.id

  metadata = {
    hello = "world"
  }
}
`, r.template(data))
}

func (r StorageShareFileResource) basicAzureADAuth(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  storage_use_azuread = true
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%[1]d"
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

resource "azurerm_storage_share_file" "test" {
  name             = "file"
  storage_share_id = azurerm_storage_share.test.id

  metadata = {
    hello = "world"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r StorageShareFileResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_file" "import" {
  name             = azurerm_storage_share_file.test.name
  storage_share_id = azurerm_storage_share_file.test.storage_share_id

  metadata = {
    hello = "world"
  }
}
`, r.basic(data))
}

func (r StorageShareFileResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_file" "test" {
  name             = "file"
  storage_share_id = azurerm_storage_share.test.id


  content_type        = "test_content_type"
  content_encoding    = "test_encoding"
  content_disposition = "test_content_disposition"

  metadata = {
    hello = "world"
  }
}
`, r.template(data))
}

func (r StorageShareFileResource) withFile(data acceptance.TestData, fileName string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_file" "test" {
  name             = "test"
  storage_share_id = azurerm_storage_share.test.id

  source      = "%[2]s"
  content_md5 = filemd5(%[2]q)

  metadata = {
    hello = "world"
  }
}
`, r.template(data), fileName)
}

func (r StorageShareFileResource) withPath(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_directory" "parent" {
  name             = "parent"
  storage_share_id = azurerm_storage_share.test.id
}

resource "azurerm_storage_share_file" "test" {
  name             = "test"
  path             = azurerm_storage_share_directory.parent.name
  storage_share_id = azurerm_storage_share.test.id
}
`, r.template(data))
}

func (r StorageShareFileResource) withPathUsingBackslashes(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_file" "test" {
  name             = "command.com"
  path             = "c\\dos"
  storage_share_id = azurerm_storage_share.test.id
  depends_on       = [azurerm_storage_share_directory.dos]
}
`, StorageShareDirectoryResource{}.nestedWithBackslashes(data))
}

func (r StorageShareFileResource) withPathInNameUsingBackslashes(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_file" "test" {
  name             = "c\\dos\\command.com"
  storage_share_id = azurerm_storage_share.test.id
  depends_on       = [azurerm_storage_share_directory.dos]
}
`, StorageShareDirectoryResource{}.nestedWithBackslashes(data))
}

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
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/file/directories"
)

type StorageShareDirectoryResource struct{}

func TestAccStorageShareDirectory_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_directory", "test")
	r := StorageShareDirectoryResource{}

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

func TestAccStorageShareDirectory_uppercase(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_directory", "test")
	r := StorageShareDirectoryResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.uppercase(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageShareDirectory_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_directory", "test")
	r := StorageShareDirectoryResource{}

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

func TestAccStorageShareDirectory_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_directory", "test")
	r := StorageShareDirectoryResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageShareDirectory_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_directory", "test")
	r := StorageShareDirectoryResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageShareDirectory_nested(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_directory", "parent")
	r := StorageShareDirectoryResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.nested(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That("azurerm_storage_share_directory.child_one").ExistsInAzure(r),
				check.That("azurerm_storage_share_directory.child_two").ExistsInAzure(r),
				check.That("azurerm_storage_share_directory.multiple_child_one").ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r StorageShareDirectoryResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := directories.ParseResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	account, err := client.Storage.FindAccount(ctx, id.AccountName)
	if err != nil {
		return nil, fmt.Errorf("retrieving Account %q for Directory %q (Share %q): %s", id.AccountName, id.DirectoryName, id.ShareName, err)
	}
	if account == nil {
		return nil, fmt.Errorf("unable to determine Resource Group for Storage Share Directory %q (Share %q / Account %q)", id.DirectoryName, id.ShareName, id.AccountName)
	}
	dirClient, err := client.Storage.FileShareDirectoriesClient(ctx, *account)
	if err != nil {
		return nil, fmt.Errorf("building File Share client for Storage Account %q (Resource Group %q): %+v", id.AccountName, account.ResourceGroup, err)
	}
	resp, err := dirClient.Get(ctx, id.AccountName, id.ShareName, id.DirectoryName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Storage Share %q (File Share %q / Account %q / Resource Group %q): %s", id.DirectoryName, id.ShareName, id.AccountName, account.ResourceGroup, err)
	}
	return utils.Bool(true), nil
}

func (r StorageShareDirectoryResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_directory" "test" {
  name                 = "dir"
  share_name           = azurerm_storage_share.test.name
  storage_account_name = azurerm_storage_account.test.name
}
`, template)
}

func (r StorageShareDirectoryResource) uppercase(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_directory" "test" {
  name                 = "UpperCaseCharacterS"
  share_name           = azurerm_storage_share.test.name
  storage_account_name = azurerm_storage_account.test.name
}
`, template)
}

func (r StorageShareDirectoryResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_directory" "import" {
  name                 = azurerm_storage_share_directory.test.name
  share_name           = azurerm_storage_share_directory.test.share_name
  storage_account_name = azurerm_storage_share_directory.test.storage_account_name
}
`, template)
}

func (r StorageShareDirectoryResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_directory" "test" {
  name                 = "dir"
  share_name           = azurerm_storage_share.test.name
  storage_account_name = azurerm_storage_account.test.name

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
  name                 = "dir"
  share_name           = azurerm_storage_share.test.name
  storage_account_name = azurerm_storage_account.test.name

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
  name                 = "123--parent-dir"
  share_name           = azurerm_storage_share.test.name
  storage_account_name = azurerm_storage_account.test.name
}

resource "azurerm_storage_share_directory" "child_one" {
  name                 = "${azurerm_storage_share_directory.parent.name}/child1"
  share_name           = azurerm_storage_share.test.name
  storage_account_name = azurerm_storage_account.test.name
}

resource "azurerm_storage_share_directory" "child_two" {
  name                 = "${azurerm_storage_share_directory.child_one.name}/childtwo--123"
  share_name           = azurerm_storage_share.test.name
  storage_account_name = azurerm_storage_account.test.name
}

resource "azurerm_storage_share_directory" "multiple_child_one" {
  name                 = "${azurerm_storage_share_directory.parent.name}/c"
  share_name           = azurerm_storage_share.test.name
  storage_account_name = azurerm_storage_account.test.name
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

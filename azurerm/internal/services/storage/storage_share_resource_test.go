package storage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type StorageShareResource struct{}

func TestAccStorageShare_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share", "test")
	r := StorageShareResource{}

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

func TestAccStorageShare_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share", "test")
	r := StorageShareResource{}

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

func TestAccStorageShare_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share", "test")
	r := StorageShareResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func TestAccStorageShare_deleteAndRecreate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share", "test")
	r := StorageShareResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.template(data),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageShare_metaData(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share", "test")
	r := StorageShareResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.metaData(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.metaDataUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageShare_acl(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share", "test")
	r := StorageShareResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.acl(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.aclUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageShare_aclGhostedRecall(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share", "test")
	r := StorageShareResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.aclGhostedRecall(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageShare_updateQuota(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share", "test")
	r := StorageShareResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.updateQuota(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("quota").HasValue("5"),
			),
		},
	})
}

func TestAccStorageShare_largeQuota(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share", "test")
	r := StorageShareResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.largeQuota(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.largeQuotaUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r StorageShareResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.StorageShareDataPlaneID(state.ID)
	if err != nil {
		return nil, err
	}

	account, err := client.Storage.FindAccount(ctx, id.AccountName)
	if err != nil {
		return nil, fmt.Errorf("retrieving Account %q for Share %q: %+v", id.AccountName, id.Name, err)
	}
	if account == nil {
		return nil, fmt.Errorf("unable to determine Account %q for Storage Share %q", id.AccountName, id.Name)
	}

	sharesClient, err := client.Storage.FileSharesClient(ctx, *account)
	if err != nil {
		return nil, fmt.Errorf("building File Share Client for Storage Account %q (Resource Group %q): %+v", id.AccountName, account.ResourceGroup, err)
	}

	props, err := sharesClient.Get(ctx, account.ResourceGroup, id.AccountName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving File Share %q (Account %q / Resource Group %q): %+v", id.Name, id.AccountName, account.ResourceGroup, err)
	}
	return utils.Bool(props != nil), nil
}

func (r StorageShareResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.StorageShareDataPlaneID(state.ID)
	if err != nil {
		return nil, err
	}

	account, err := client.Storage.FindAccount(ctx, id.AccountName)
	if err != nil {
		return nil, fmt.Errorf("retrieving Account %q for Share %q: %+v", id.AccountName, id.Name, err)
	}
	if account == nil {
		return nil, fmt.Errorf("unable to determine Account %q for Storage Share %q", id.AccountName, id.Name)
	}

	sharesClient, err := client.Storage.FileSharesClient(ctx, *account)
	if err != nil {
		return nil, fmt.Errorf("building File Share Client for Storage Account %q (Resource Group %q): %+v", id.AccountName, account.ResourceGroup, err)
	}
	if err := sharesClient.Delete(ctx, account.ResourceGroup, id.AccountName, id.Name); err != nil {
		return nil, fmt.Errorf("deleting File Share %q (Account %q / Resource Group %q): %+v", id.Name, id.AccountName, account.ResourceGroup, err)
	}
	return utils.Bool(true), nil
}

func (r StorageShareResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share" "test" {
  name                 = "testshare%s"
  storage_account_name = azurerm_storage_account.test.name
}
`, template, data.RandomString)
}

func (r StorageShareResource) metaData(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share" "test" {
  name                 = "testshare%s"
  storage_account_name = azurerm_storage_account.test.name

  metadata = {
    hello = "world"
  }
}
`, template, data.RandomString)
}

func (r StorageShareResource) metaDataUpdated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share" "test" {
  name                 = "testshare%s"
  storage_account_name = azurerm_storage_account.test.name

  metadata = {
    hello = "world"
    happy = "birthday"
  }
}
`, template, data.RandomString)
}

func (r StorageShareResource) acl(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share" "test" {
  name                 = "testshare%s"
  storage_account_name = azurerm_storage_account.test.name

  acl {
    id = "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI"

    access_policy {
      permissions = "rwd"
      start       = "2019-07-02T09:38:21.0000000Z"
      expiry      = "2019-07-02T10:38:21.0000000Z"
    }
  }
}
`, template, data.RandomString)
}

func (r StorageShareResource) aclGhostedRecall(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share" "test" {
  name                 = "testshare%s"
  storage_account_name = azurerm_storage_account.test.name

  acl {
    id = "GhostedRecall"
    access_policy {
      permissions = "r"
    }
  }
}
`, template, data.RandomString)
}

func (r StorageShareResource) aclUpdated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share" "test" {
  name                 = "testshare%s"
  storage_account_name = azurerm_storage_account.test.name

  acl {
    id = "AAAANDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI"

    access_policy {
      permissions = "rwd"
      start       = "2019-07-02T09:38:21.0000000Z"
      expiry      = "2019-07-02T10:38:21.0000000Z"
    }
  }
  acl {
    id = "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI"

    access_policy {
      permissions = "rwd"
      start       = "2019-07-02T09:38:21.0000000Z"
      expiry      = "2019-07-02T10:38:21.0000000Z"
    }
  }
}
`, template, data.RandomString)
}

func (r StorageShareResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share" "import" {
  name                 = azurerm_storage_share.test.name
  storage_account_name = azurerm_storage_share.test.storage_account_name
}
`, template)
}

func (r StorageShareResource) updateQuota(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share" "test" {
  name                 = "testshare%s"
  storage_account_name = azurerm_storage_account.test.name
  quota                = 5
}
`, template, data.RandomString)
}

func (r StorageShareResource) largeQuota(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storageshare-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestshare%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Premium"
  account_replication_type = "LRS"
  account_kind             = "FileStorage"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_share" "test" {
  name                 = "testshare%s"
  storage_account_name = azurerm_storage_account.test.name
  quota                = 6000
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func (r StorageShareResource) largeQuotaUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storageshare-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestshare%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Premium"
  account_replication_type = "LRS"
  account_kind             = "FileStorage"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_share" "test" {
  name                 = "testshare%s"
  storage_account_name = azurerm_storage_account.test.name
  quota                = 10000
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func (r StorageShareResource) template(data acceptance.TestData) string {
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
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

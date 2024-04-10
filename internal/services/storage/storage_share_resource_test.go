// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2023-11-03/file/shares"
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
				check.That(data.ResourceName).Key("enabled_protocol").HasValue("SMB"),
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

func TestAccStorageShare_accessTierStandard(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share", "test")
	r := StorageShareResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.accessTierStandard(data, "Cool"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.accessTierStandard(data, "Hot"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.accessTierStandard(data, "TransactionOptimized"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageShare_accessTierPremium(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share", "test")
	r := StorageShareResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.accessTierPremium(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageShare_nfsProtocol(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share", "test")
	r := StorageShareResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.protocol(data, "NFS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

// TestAccStorageShare_protocolUpdate is to ensure destroy-then-create of the storage share can tolerant the "ShareBeingDeleted" issue.
func TestAccStorageShare_protocolUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share", "test")
	r := StorageShareResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.protocol(data, "NFS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.protocol(data, "SMB"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r StorageShareResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := shares.ParseShareID(state.ID, client.Storage.StorageDomainSuffix)
	if err != nil {
		return nil, err
	}

	account, err := client.Storage.FindAccount(ctx, client.Account.SubscriptionId, id.AccountId.AccountName)
	if err != nil {
		return nil, fmt.Errorf("retrieving Account %q for Share %q: %+v", id.AccountId.AccountName, id.ShareName, err)
	}
	if account == nil {
		return nil, fmt.Errorf("unable to determine Account %q for Storage Share %q", id.AccountId.AccountName, id.ShareName)
	}

	sharesClient, err := client.Storage.FileSharesDataPlaneClient(ctx, *account, client.Storage.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return nil, fmt.Errorf("building File Share Client for %s: %+v", account.StorageAccountId, err)
	}

	props, err := sharesClient.Get(ctx, id.ShareName)
	if err != nil {
		return nil, fmt.Errorf("retrieving File Share %q in %s: %+v", id.ShareName, account.StorageAccountId, err)
	}

	return utils.Bool(props != nil), nil
}

func (r StorageShareResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := shares.ParseShareID(state.ID, client.Storage.StorageDomainSuffix)
	if err != nil {
		return nil, err
	}

	account, err := client.Storage.FindAccount(ctx, client.Account.SubscriptionId, id.AccountId.AccountName)
	if err != nil {
		return nil, fmt.Errorf("retrieving Account %q for Share %q: %+v", id.AccountId.AccountName, id.ShareName, err)
	}
	if account == nil {
		return nil, fmt.Errorf("unable to determine Account %q for Storage Share %q", id.AccountId.AccountName, id.ShareName)
	}

	sharesClient, err := client.Storage.FileSharesDataPlaneClient(ctx, *account, client.Storage.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return nil, fmt.Errorf("building File Share Client for %s: %+v", account.StorageAccountId, err)
	}
	if err := sharesClient.Delete(ctx, id.ShareName); err != nil {
		return nil, fmt.Errorf("deleting File Share %q in %s: %+v", id.ShareName, account.StorageAccountId, err)
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
  quota                = 5
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
  quota                = 5

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
  quota                = 5

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
  quota                = 5

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
  quota                = 5

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
  quota                = 5

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
  quota                = 5
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

func (r StorageShareResource) accessTierStandard(data acceptance.TestData, tier string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
  storage_use_azuread = true
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
  account_kind             = "StorageV2"
}

resource "azurerm_storage_share" "test" {
  name                 = "testshare%s"
  storage_account_name = azurerm_storage_account.test.name
  quota                = 100
  enabled_protocol     = "SMB"
  access_tier          = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, tier)
}

func (r StorageShareResource) accessTierPremium(data acceptance.TestData) string {
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
  account_tier             = "Premium"
  account_replication_type = "LRS"
  account_kind             = "FileStorage"
}

resource "azurerm_storage_share" "test" {
  name                 = "testshare%s"
  storage_account_name = azurerm_storage_account.test.name
  quota                = 100
  enabled_protocol     = "SMB"
  access_tier          = "Premium"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func (r StorageShareResource) protocol(data acceptance.TestData, protocol string) string {
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
  account_kind             = "FileStorage"
  account_tier             = "Premium"
  account_replication_type = "LRS"
}

resource "azurerm_storage_share" "test" {
  name                 = "testshare%s"
  storage_account_name = azurerm_storage_account.test.name
  enabled_protocol     = "%s"
  quota                = 100
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, protocol)
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

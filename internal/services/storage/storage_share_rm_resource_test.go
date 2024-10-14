// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/fileshares"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type StorageShareResourceRm struct{}

func TestAccStorageShareRm_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_rm", "test")
	r := StorageShareResourceRm{}

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

func TestAccStorageShareRm_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_rm", "test")
	r := StorageShareResourceRm{}

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

func TestAccStorageShareRm_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_rm", "test")
	r := StorageShareResourceRm{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func TestAccStorageShareRm_deleteAndRecreate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_rm", "test")
	r := StorageShareResourceRm{}

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

func TestAccStorageShareRm_metaData(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_rm", "test")
	r := StorageShareResourceRm{}

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

func TestAccStorageShareRm_acl(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_rm", "test")
	r := StorageShareResourceRm{}

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

func TestAccStorageShareRm_aclGhostedRecall(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_rm", "test")
	r := StorageShareResourceRm{}

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

func TestAccStorageShareRm_updateQuota(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_rm", "test")
	r := StorageShareResourceRm{}

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

func TestAccStorageShareRm_largeQuota(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_rm", "test")
	r := StorageShareResourceRm{}

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

func TestAccStorageShareRm_accessTierStandard(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_rm", "test")
	r := StorageShareResourceRm{}

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
	})
}

func TestAccStorageShareRm_accessTierPremium(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_rm", "test")
	r := StorageShareResourceRm{}

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

func TestAccStorageShareRm_nfsProtocol(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_rm", "test")
	r := StorageShareResourceRm{}

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
func TestAccStorageShareRm_protocolUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_share_rm", "test")
	r := StorageShareResourceRm{}

	data.ResourceTestIgnoreRecreate(t, r, []acceptance.TestStep{
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

func (r StorageShareResourceRm) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := fileshares.ParseShareID(state.ID)
	if err != nil {
		return nil, err
	}

	existing, err := client.Storage.ResourceManager.FileShares.Get(ctx, *id, fileshares.DefaultGetOperationOptions())
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return utils.Bool(false), nil
		}
	}

	return utils.Bool(true), nil
}

func (r StorageShareResourceRm) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := fileshares.ParseShareID(state.ID)
	if err != nil {
		return nil, err
	}

	if _, err = client.Storage.ResourceManager.FileShares.Delete(ctx, *id, fileshares.DefaultDeleteOperationOptions()); err != nil {
		return nil, fmt.Errorf("deleting File Share %q in %s: %+v", id.ShareName, id.StorageAccountName, err)
	}

	return utils.Bool(true), nil
}

func (r StorageShareResourceRm) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_rm" "test" {
  name                 = "testshare%s"
  storage_account_name = azurerm_storage_account.test.name
  resource_group_name  = azurerm_resource_group.test.name
  quota                = 5
}
`, template, data.RandomString)
}

func (r StorageShareResourceRm) metaData(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_rm" "test" {
  name                 = "testshare%s"
  storage_account_name = azurerm_storage_account.test.name
  resource_group_name  = azurerm_resource_group.test.name
  quota                = 5

  metadata = {
    hello = "world"
  }
}
`, template, data.RandomString)
}

func (r StorageShareResourceRm) metaDataUpdated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_rm" "test" {
  name                 = "testshare%s"
  storage_account_name = azurerm_storage_account.test.name
  resource_group_name  = azurerm_resource_group.test.name
  quota                = 5

  metadata = {
    hello = "world"
    happy = "birthday"
  }
}
`, template, data.RandomString)
}

func (r StorageShareResourceRm) acl(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_rm" "test" {
  name                 = "testshare%s"
  storage_account_name = azurerm_storage_account.test.name
  resource_group_name  = azurerm_resource_group.test.name
  quota                = 5

  acl {
    id = "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI"

    access_policy {
      permission  = "rwd"
      start_time  = "2019-07-02T09:38:21.0000000Z"
      expiry_time = "2019-07-02T10:38:21.0000000Z"
    }
  }
}
`, template, data.RandomString)
}

func (r StorageShareResourceRm) aclGhostedRecall(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_rm" "test" {
  name                 = "testshare%s"
  storage_account_name = azurerm_storage_account.test.name
  resource_group_name  = azurerm_resource_group.test.name
  quota                = 5

  acl {
    id = "GhostedRecall"
    access_policy {
      permission = "r"
    }
  }
}
`, template, data.RandomString)
}

func (r StorageShareResourceRm) aclUpdated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_rm" "test" {
  name                 = "testshare%s"
  storage_account_name = azurerm_storage_account.test.name
  resource_group_name  = azurerm_resource_group.test.name
  quota                = 5

  acl {
    id = "AAAANDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI"

    access_policy {
      permission  = "rwd"
      start_time  = "2019-07-02T09:38:21.0000000Z"
      expiry_time = "2019-07-02T10:38:21.0000000Z"
    }
  }
  acl {
    id = "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI"

    access_policy {
      permission  = "rwd"
      start_time  = "2019-07-02T09:38:21.0000000Z"
      expiry_time = "2019-07-02T10:38:21.0000000Z"
    }
  }
}
`, template, data.RandomString)
}

func (r StorageShareResourceRm) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_rm" "import" {
  name                 = azurerm_storage_share_rm.test.name
  storage_account_name = azurerm_storage_share_rm.test.storage_account_name
  resource_group_name  = azurerm_resource_group.test.name
  quota                = 5
}
`, template)
}

func (r StorageShareResourceRm) updateQuota(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share_rm" "test" {
  name                 = "testshare%s"
  storage_account_name = azurerm_storage_account.test.name
  resource_group_name  = azurerm_resource_group.test.name
  quota                = 5
}
`, template, data.RandomString)
}

func (r StorageShareResourceRm) largeQuota(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storageshare-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                            = "acctestshare%s"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  account_tier                    = "Premium"
  account_replication_type        = "LRS"
  account_kind                    = "FileStorage"
  allow_nested_items_to_be_public = false

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_share_rm" "test" {
  name                 = "testshare%s"
  storage_account_name = azurerm_storage_account.test.name
  resource_group_name  = azurerm_resource_group.test.name
  quota                = 6000
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func (r StorageShareResourceRm) largeQuotaUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storageshare-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                            = "acctestshare%s"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  account_tier                    = "Premium"
  account_replication_type        = "LRS"
  account_kind                    = "FileStorage"
  allow_nested_items_to_be_public = false

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_share_rm" "test" {
  name                 = "testshare%s"
  storage_account_name = azurerm_storage_account.test.name
  resource_group_name  = azurerm_resource_group.test.name
  quota                = 10000
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func (r StorageShareResourceRm) accessTierStandard(data acceptance.TestData, tier string) string {
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
  name                            = "acctestacc%s"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  account_tier                    = "Standard"
  account_replication_type        = "LRS"
  account_kind                    = "StorageV2"
  allow_nested_items_to_be_public = false
}

resource "azurerm_storage_share_rm" "test" {
  name                 = "testshare%s"
  storage_account_name = azurerm_storage_account.test.name
  resource_group_name  = azurerm_resource_group.test.name
  quota                = 100
  enabled_protocol     = "SMB"
  access_tier          = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, tier)
}

func (r StorageShareResourceRm) accessTierPremium(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                            = "acctestacc%s"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  account_tier                    = "Premium"
  account_replication_type        = "LRS"
  account_kind                    = "FileStorage"
  allow_nested_items_to_be_public = false
}

resource "azurerm_storage_share_rm" "test" {
  name                 = "testshare%s"
  storage_account_name = azurerm_storage_account.test.name
  resource_group_name  = azurerm_resource_group.test.name
  quota                = 100
  enabled_protocol     = "SMB"
  access_tier          = "Premium"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func (r StorageShareResourceRm) protocol(data acceptance.TestData, protocol string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                            = "acctestacc%s"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  account_kind                    = "FileStorage"
  account_tier                    = "Premium"
  account_replication_type        = "LRS"
  allow_nested_items_to_be_public = false
}

resource "azurerm_storage_share_rm" "test" {
  name                 = "testshare%s"
  storage_account_name = azurerm_storage_account.test.name
  resource_group_name  = azurerm_resource_group.test.name
  enabled_protocol     = "%s"
  quota                = 100
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, protocol)
}

func (r StorageShareResourceRm) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                            = "acctestacc%s"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  account_tier                    = "Standard"
  account_replication_type        = "LRS"
  allow_nested_items_to_be_public = false

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

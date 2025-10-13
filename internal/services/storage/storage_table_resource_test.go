// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-05-01/tableservice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/jackofallops/giovanni/storage/2023-11-03/table/tables"
)

type StorageTableResource struct{}

func TestAccStorageTable_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_table", "test")
	r := StorageTableResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("url").HasValue(fmt.Sprintf("https://acctestacc%s.table.core.windows.net/Tables('acctestst%d')", data.RandomString, data.RandomInteger)),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageTable_basicDeprecated(t *testing.T) {
	if features.FivePointOh() {
		t.Skip("skipping as not valid in 5.0")
	}

	data := acceptance.BuildTestData(t, "azurerm_storage_table", "test")
	r := StorageTableResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicDeprecated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("url").HasValue(fmt.Sprintf("https://acctestacc%s.table.core.windows.net/Tables('acctestst%d')", data.RandomString, data.RandomInteger)),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageTable_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_table", "test")
	r := StorageTableResource{}

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

func TestAccStorageTable_requiresImportDeprecated(t *testing.T) {
	if features.FivePointOh() {
		t.Skip("skipping as not valid in 5.0")
	}

	data := acceptance.BuildTestData(t, "azurerm_storage_table", "test")
	r := StorageTableResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicDeprecated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImportDeprecated),
	})
}

func TestAccStorageTable_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_table", "test")
	r := StorageTableResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func TestAccStorageTable_disappearsDeprecated(t *testing.T) {
	if features.FivePointOh() {
		t.Skip("skipping as not valid in 5.0")
	}

	data := acceptance.BuildTestData(t, "azurerm_storage_table", "test")
	r := StorageTableResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basicDeprecated,
			TestResource: r,
		}),
	})
}

func TestAccStorageTable_acl(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_table", "test")
	r := StorageTableResource{}

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

func TestAccStorageTable_aclDeprecated(t *testing.T) {
	if features.FivePointOh() {
		t.Skip("skipping as not valid in 5.0")
	}

	data := acceptance.BuildTestData(t, "azurerm_storage_table", "test")
	r := StorageTableResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.aclDeprecated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.aclUpdatedDeprecated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageTable_migrateToStorageID(t *testing.T) {
	if features.FivePointOh() {
		t.Skip("skipping as test is not valid in 5.0")
	}
	data := acceptance.BuildTestData(t, "azurerm_storage_table", "test")
	r := StorageTableResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicDeprecated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("storage_account_name").IsSet(),
				check.That(data.ResourceName).Key("storage_account_id").DoesNotExist(),
				check.That(data.ResourceName).Key("id").MatchesRegex(regexp.MustCompile("https:*")),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("storage_account_name").IsEmpty(),
				check.That(data.ResourceName).Key("storage_account_id").IsSet(),
				check.That(data.ResourceName).Key("id").MatchesRegex(regexp.MustCompile("/subscriptions/*")),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageTable_migrateFromStorageIDShouldFail(t *testing.T) {
	if features.FivePointOh() {
		t.Skip("skipping as test is not valid in 5.0")
	}
	data := acceptance.BuildTestData(t, "azurerm_storage_table", "test")
	r := StorageTableResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("storage_account_id").IsSet(),
				check.That(data.ResourceName).Key("storage_account_name").IsEmpty(),
			),
		},
		data.ImportStep(),
		{
			Config:      r.basicDeprecated(data),
			ExpectError: regexp.MustCompile("expected action to not be Replace"),
		},
	})
}

func (r StorageTableResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	if !features.FivePointOh() && !strings.HasPrefix(state.ID, "/subscriptions") {
		id, err := tables.ParseTableID(state.ID, client.Storage.StorageDomainSuffix)
		if err != nil {
			return nil, err
		}
		account, err := client.Storage.FindAccount(ctx, client.Account.SubscriptionId, id.AccountId.AccountName)
		if err != nil {
			return nil, fmt.Errorf("retrieving Account %q for Table %q: %+v", id.AccountId.AccountName, id.TableName, err)
		}
		if account == nil {
			return nil, fmt.Errorf("unable to determine Resource Group for Storage Table %q (Account %q)", id.TableName, id.AccountId.AccountName)
		}
		tablesClient, err := client.Storage.TablesDataPlaneClient(ctx, *account, client.Storage.DataPlaneOperationSupportingAnyAuthMethod())
		if err != nil {
			return nil, fmt.Errorf("building Table Client: %+v", err)
		}

		return tablesClient.Exists(ctx, id.TableName)
	}

	id, err := tableservice.ParseTableID(state.ID)
	if err != nil {
		return nil, err
	}
	existing, err := client.Storage.ResourceManager.TableService.TableGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(existing.Model != nil), nil
}

func (r StorageTableResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	if !features.FivePointOh() && !strings.HasPrefix(state.ID, "/subscriptions") {
		id, err := tables.ParseTableID(state.ID, client.Storage.StorageDomainSuffix)
		if err != nil {
			return nil, err
		}
		account, err := client.Storage.FindAccount(ctx, client.Account.SubscriptionId, id.AccountId.AccountName)
		if err != nil {
			return nil, fmt.Errorf("retrieving Account %q for Table %q: %+v", id.AccountId.AccountName, id.TableName, err)
		}
		if account == nil {
			return nil, fmt.Errorf("unable to determine Resource Group for Storage Table %q (Account %q)", id.TableName, id.AccountId.AccountName)
		}
		tablesClient, err := client.Storage.TablesDataPlaneClient(ctx, *account, client.Storage.DataPlaneOperationSupportingAnyAuthMethod())
		if err != nil {
			return nil, fmt.Errorf("building Table Client: %+v", err)
		}

		exists, err := tablesClient.Exists(ctx, id.TableName)
		if err != nil {
			return nil, fmt.Errorf("retrieving Table %q (Account %q): %+v", id.TableName, id.AccountId.AccountName, err)
		}
		if exists == nil || !*exists {
			return nil, fmt.Errorf("table %q doesn't exist in Account %q so it can't be deleted", id.TableName, id.AccountId.AccountName)
		}
		if err := tablesClient.Delete(ctx, id.TableName); err != nil {
			return nil, fmt.Errorf("deleting Table %q (Account %q): %+v", id.TableName, id.AccountId.AccountName, err)
		}
		return utils.Bool(true), nil
	}

	id, err := tableservice.ParseTableID(state.ID)
	if err != nil {
		return nil, err
	}
	if _, err := client.Storage.ResourceManager.TableService.TableDelete(ctx, *id); err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(true), nil
}

func (r StorageTableResource) basic(data acceptance.TestData) string {
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

resource "azurerm_storage_table" "test" {
  name               = "acctestst%d"
  storage_account_id = azurerm_storage_account.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func (r StorageTableResource) basicDeprecated(data acceptance.TestData) string {
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

resource "azurerm_storage_table" "test" {
  name                 = "acctestst%d"
  storage_account_name = azurerm_storage_account.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func (r StorageTableResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_table" "import" {
  name               = azurerm_storage_table.test.name
  storage_account_id = azurerm_storage_table.test.storage_account_id
}
`, template)
}

func (r StorageTableResource) requiresImportDeprecated(data acceptance.TestData) string {
	template := r.basicDeprecated(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_table" "import" {
  name                 = azurerm_storage_table.test.name
  storage_account_name = azurerm_storage_table.test.storage_account_name
}
`, template)
}

func (r StorageTableResource) acl(data acceptance.TestData) string {
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

resource "azurerm_storage_table" "test" {
  name               = "acctestst%d"
  storage_account_id = azurerm_storage_account.test.id
  acl {
    id = "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI"

    access_policy {
      permissions = "raud"
      start       = "2020-11-26T08:49:37.0000000Z"
      expiry      = "2020-11-27T08:49:37.0000000Z"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func (r StorageTableResource) aclDeprecated(data acceptance.TestData) string {
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

resource "azurerm_storage_table" "test" {
  name                 = "acctestst%d"
  storage_account_name = azurerm_storage_account.test.name
  acl {
    id = "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI"

    access_policy {
      permissions = "raud"
      start       = "2020-11-26T08:49:37.0000000Z"
      expiry      = "2020-11-27T08:49:37.0000000Z"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func (r StorageTableResource) aclUpdated(data acceptance.TestData) string {
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

resource "azurerm_storage_table" "test" {
  name               = "acctestst%d"
  storage_account_id = azurerm_storage_account.test.id

  acl {
    id = "AAAANDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI"

    access_policy {
      permissions = "raud"
      start       = "2020-11-26T08:49:37.0000000Z"
      expiry      = "2020-11-27T08:49:37.0000000Z"
    }
  }
  acl {
    id = "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI"

    access_policy {
      permissions = "raud"
      start       = "2019-07-02T09:38:21.0000000Z"
      expiry      = "2019-07-02T10:38:21.0000000Z"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func (r StorageTableResource) aclUpdatedDeprecated(data acceptance.TestData) string {
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

resource "azurerm_storage_table" "test" {
  name                 = "acctestst%d"
  storage_account_name = azurerm_storage_account.test.name

  acl {
    id = "AAAANDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI"

    access_policy {
      permissions = "raud"
      start       = "2020-11-26T08:49:37.0000000Z"
      expiry      = "2020-11-27T08:49:37.0000000Z"
    }
  }
  acl {
    id = "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI"

    access_policy {
      permissions = "raud"
      start       = "2019-07-02T09:38:21.0000000Z"
      expiry      = "2019-07-02T10:38:21.0000000Z"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

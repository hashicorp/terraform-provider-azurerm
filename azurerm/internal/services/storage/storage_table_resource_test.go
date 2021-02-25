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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type StorageTableResource struct{}

func TestAccStorageTable_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_table", "test")
	r := StorageTableResource{}

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

func TestAccStorageTable_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_table", "test")
	r := StorageTableResource{}

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

func TestAccStorageTable_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_table", "test")
	r := StorageTableResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func TestAccStorageTable_acl(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_table", "test")
	r := StorageTableResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.acl(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.aclUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r StorageTableResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.StorageTableDataPlaneID(state.ID)
	if err != nil {
		return nil, err
	}
	account, err := client.Storage.FindAccount(ctx, id.AccountName)
	if err != nil {
		return nil, fmt.Errorf("retrieving Account %q for Table %q: %+v", id.AccountName, id.Name, err)
	}
	if account == nil {
		return nil, fmt.Errorf("unable to determine Resource Group for Storage Storage Table %q (Account %q)", id.Name, id.AccountName)
	}
	tablesClient, err := client.Storage.TablesClient(ctx, *account)
	if err != nil {
		return nil, fmt.Errorf("building Table Client: %+v", err)
	}

	return tablesClient.Exists(ctx, account.ResourceGroup, id.AccountName, id.Name)
}

func (r StorageTableResource) Destroy(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.StorageTableDataPlaneID(state.ID)
	if err != nil {
		return nil, err
	}
	account, err := client.Storage.FindAccount(ctx, id.AccountName)
	if err != nil {
		return nil, fmt.Errorf("retrieving Account %q for Table %q: %+v", id.AccountName, id.Name, err)
	}
	if account == nil {
		return nil, fmt.Errorf("unable to determine Resource Group for Storage Storage Table %q (Account %q)", id.Name, id.AccountName)
	}
	tablesClient, err := client.Storage.TablesClient(ctx, *account)
	if err != nil {
		return nil, fmt.Errorf("building Table Client: %+v", err)
	}

	exists, err := tablesClient.Exists(ctx, account.ResourceGroup, id.AccountName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Table %q (Account %q): %+v", id.Name, id.AccountName, err)
	}
	if exists == nil || !*exists {
		return nil, fmt.Errorf("table %q doesn't exist in Account %q so it can't be deleted", id.Name, id.AccountName)
	}
	if err := tablesClient.Delete(ctx, account.ResourceGroup, id.AccountName, id.Name); err != nil {
		return nil, fmt.Errorf("deleting Table %q (Account %q): %+v", id.Name, id.AccountName, err)
	}
	return utils.Bool(true), nil
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

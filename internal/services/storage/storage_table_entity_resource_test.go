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
	"github.com/jackofallops/giovanni/storage/2023-11-03/table/entities"
)

type StorageTableEntityResource struct{}

func TestAccTableEntity_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_table_entity", "test")
	r := StorageTableEntityResource{}

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

func TestAccTableEntity_basicAzureADAuth(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_table_entity", "test")
	r := StorageTableEntityResource{}

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

func TestAccTableEntity_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_table_entity", "test")
	r := StorageTableEntityResource{}

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

func TestAccTableEntity_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_table_entity", "test")
	r := StorageTableEntityResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
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

func TestAccTableEntity_updateTyped(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_table_entity", "test")
	r := StorageTableEntityResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateType(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updatedTypeInt64(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.updatedTypeDouble(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.updateTypeString(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.updateTypeBoolean(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (r StorageTableEntityResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := entities.ParseEntityID(state.ID, client.Storage.StorageDomainSuffix)
	if err != nil {
		return nil, err
	}
	account, err := client.Storage.FindAccount(ctx, client.Account.SubscriptionId, id.AccountId.AccountName)
	if err != nil {
		return nil, fmt.Errorf("retrieving Account %q for Table %q: %+v", id.AccountId.AccountName, id.TableName, err)
	}
	if account == nil {
		return nil, fmt.Errorf("storage Account %q was not found", id.AccountId.AccountName)
	}

	entitiesClient, err := client.Storage.TableEntityDataPlaneClient(ctx, *account, client.Storage.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return nil, fmt.Errorf("building Table Entity Client: %+v", err)
	}

	input := entities.GetEntityInput{
		PartitionKey:  id.PartitionKey,
		RowKey:        id.RowKey,
		MetaDataLevel: entities.NoMetaData,
	}
	resp, err := entitiesClient.Get(ctx, id.TableName, input)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Entity (Partition Key %q / Row Key %q) (Table %q in %s): %+v", id.PartitionKey, id.RowKey, id.TableName, account.StorageAccountId, err)
	}
	return utils.Bool(true), nil
}

func (r StorageTableEntityResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_table_entity" "test" {
  storage_table_id = azurerm_storage_table.test.id

  partition_key = "test_partition%[2]d"
  row_key       = "test_row%[2]d"
  entity = {
    Foo = "Bar"
  }
}
`, template, data.RandomInteger)
}

func (r StorageTableEntityResource) basicAzureADAuth(data acceptance.TestData) string {
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

resource "azurerm_storage_table" "test" {
  name                 = "acctestst%[1]d"
  storage_account_name = azurerm_storage_account.test.name
}

resource "azurerm_storage_table_entity" "test" {
  storage_table_id = azurerm_storage_table.test.id

  partition_key = "test_partition%[1]d"
  row_key       = "test_row%[1]d"
  entity = {
    Foo = "Bar"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r StorageTableEntityResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_table_entity" "import" {
  storage_table_id = azurerm_storage_table.test.id

  partition_key = "test_partition%[2]d"
  row_key       = "test_row%[2]d"
  entity = {
    Foo = "Bar"
  }
}
`, template, data.RandomInteger)
}

func (r StorageTableEntityResource) updated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_table_entity" "test" {
  storage_table_id = azurerm_storage_table.test.id

  partition_key = "test_partition%[2]d"
  row_key       = "test_row%[2]d"
  entity = {
    Foo  = "Bar"
    Test = "Updated"
  }
}
`, template, data.RandomInteger)
}

func (r StorageTableEntityResource) updateType(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_table_entity" "test" {
  storage_table_id = azurerm_storage_table.test.id

  partition_key = "test_partition%[2]d"
  row_key       = "test_row%[2]d"
  entity = {
    Foo              = 123
    "Foo@odata.type" = "Edm.Int32"
    Test             = "Updated"
  }
}
`, template, data.RandomInteger)
}

func (r StorageTableEntityResource) updatedTypeInt64(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_table_entity" "test" {
  storage_table_id = azurerm_storage_table.test.id

  partition_key = "test_partition%[2]d"
  row_key       = "test_row%[2]d"
  entity = {
    Foo              = 123
    "Foo@odata.type" = "Edm.Int64"
    Test             = "Updated"
  }
}
`, template, data.RandomInteger)
}

func (r StorageTableEntityResource) updatedTypeDouble(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_table_entity" "test" {
  storage_table_id = azurerm_storage_table.test.id

  partition_key = "test_partition%[2]d"
  row_key       = "test_row%[2]d"
  entity = {
    Foo              = 123.123
    "Foo@odata.type" = "Edm.Double"
    Test             = "Updated"
  }
}
`, template, data.RandomInteger)
}

func (r StorageTableEntityResource) updateTypeString(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_table_entity" "test" {
  storage_table_id = azurerm_storage_table.test.id

  partition_key = "test_partition%[2]d"
  row_key       = "test_row%[2]d"
  entity = {
    Foo              = "123.123"
    "Foo@odata.type" = "Edm.String"
    Test             = "Updated"
  }
  lifecycle {
    ignore_changes = [entity]
  }
}
`, template, data.RandomInteger)
}

func (r StorageTableEntityResource) updateTypeBoolean(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_table_entity" "test" {
  storage_table_id = azurerm_storage_table.test.id

  partition_key = "test_partition%[2]d"
  row_key       = "test_row%[2]d"
  entity = {
    Foo              = "true"
    "Foo@odata.type" = "Edm.Boolean"
    Test             = "Updated"
  }
}
`, template, data.RandomInteger)
}

func (r StorageTableEntityResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
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

resource "azurerm_storage_table" "test" {
  name                 = "acctestst%[1]d"
  storage_account_name = azurerm_storage_account.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

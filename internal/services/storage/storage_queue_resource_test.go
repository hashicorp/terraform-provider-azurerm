// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/storagequeues"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StorageQueueResource struct{}

func TestAccStorageQueue_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_queue", "test")
	r := StorageQueueResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("url").HasValue(fmt.Sprintf("https://acctestacc%s.queue.core.windows.net/acctestmysamplequeue-%d", data.RandomString, data.RandomInteger)),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageQueue_basicAzureADAuth(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_queue", "test")
	r := StorageQueueResource{}

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

func TestAccStorageQueue_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_queue", "test")
	r := StorageQueueResource{}

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

func TestAccStorageQueue_metaData(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_queue", "test")
	r := StorageQueueResource{}

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

func (r StorageQueueResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := storagequeues.ParseQueueID(state.ID)
	if err != nil {
		return nil, err
	}
	existing, err := client.Storage.ResourceManager.StorageQueues.QueueGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(existing.Model != nil), nil
}

func (r StorageQueueResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_queue" "test" {
  name               = "acctestmysamplequeue-%d"
  storage_account_id = azurerm_storage_account.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r StorageQueueResource) basicAzureADAuth(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  storage_use_azuread = true
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

resource "azurerm_storage_queue" "test" {
  name               = "acctestmysamplequeue-%d"
  storage_account_id = azurerm_storage_account.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func (r StorageQueueResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_queue" "import" {
  name               = azurerm_storage_queue.test.name
  storage_account_id = azurerm_storage_queue.test.storage_account_id
}
`, r.basic(data))
}

func (r StorageQueueResource) metaData(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_queue" "test" {
  name               = "acctestmysamplequeue-%d"
  storage_account_id = azurerm_storage_account.test.id

  metadata = {
    hello = "world"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r StorageQueueResource) metaDataUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_queue" "test" {
  name               = "acctestmysamplequeue-%d"
  storage_account_id = azurerm_storage_account.test.id

  metadata = {
    hello = "world"
    rick  = "M0rty"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r StorageQueueResource) template(data acceptance.TestData) string {
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

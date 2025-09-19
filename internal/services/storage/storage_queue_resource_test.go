// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-05-01/queueservice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/jackofallops/giovanni/storage/2023-11-03/queue/queues"
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

func TestAccStorageQueue_basicDeprecated(t *testing.T) {
	if features.FivePointOh() {
		t.Skip("skipping as not valid in 5.0")
	}

	data := acceptance.BuildTestData(t, "azurerm_storage_queue", "test")
	r := StorageQueueResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicDeprecated(data),
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

func TestAccStorageQueue_basicAzureADAuthDeprecated(t *testing.T) {
	if features.FivePointOh() {
		t.Skip("skipping as not valid in 5.0")
	}

	data := acceptance.BuildTestData(t, "azurerm_storage_queue", "test")
	r := StorageQueueResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicAzureADAuthDeprecated(data),
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

func TestAccStorageQueue_requiresImportDeprecated(t *testing.T) {
	if features.FivePointOh() {
		t.Skip("skipping as not valid in 5.0")
	}

	data := acceptance.BuildTestData(t, "azurerm_storage_queue", "test")
	r := StorageQueueResource{}

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

func TestAccStorageQueue_metaDataDeprecated(t *testing.T) {
	if features.FivePointOh() {
		t.Skip("skipping as not valid in 5.0")
	}

	data := acceptance.BuildTestData(t, "azurerm_storage_queue", "test")
	r := StorageQueueResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.metaDataDeprecated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.metaDataUpdatedDeprecated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r StorageQueueResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	if !features.FivePointOh() && !strings.HasPrefix(state.ID, "/subscriptions") {
		id, err := queues.ParseQueueID(state.ID, client.Storage.StorageDomainSuffix)
		if err != nil {
			return nil, err
		}
		account, err := client.Storage.FindAccount(ctx, client.Account.SubscriptionId, id.AccountId.AccountName)
		if err != nil {
			return nil, fmt.Errorf("retrieving Account %q for Queue %q: %+v", id.AccountId.AccountName, id.QueueName, err)
		}
		if account == nil {
			return nil, fmt.Errorf("unable to determine Resource Group for Storage Queue %q (Account %q)", id.QueueName, id.AccountId.AccountName)
		}
		queuesClient, err := client.Storage.QueuesDataPlaneClient(ctx, *account, client.Storage.DataPlaneOperationSupportingAnyAuthMethod())
		if err != nil {
			return nil, fmt.Errorf("building Queues Client: %+v", err)
		}
		queue, err := queuesClient.Get(ctx, id.QueueName)
		if err != nil {
			return nil, fmt.Errorf("retrieving Queue %q (Account %q): %+v", id.QueueName, id.AccountId.AccountName, err)
		}
		return pointer.To(queue != nil), nil
	}

	id, err := queueservice.ParseQueueID(state.ID)
	if err != nil {
		return nil, err
	}
	existing, err := client.Storage.ResourceManager.QueueService.QueueGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(existing.Model != nil), nil
}

func (r StorageQueueResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_queue" "test" {
  name               = "acctestmysamplequeue-%d"
  storage_account_id = azurerm_storage_account.test.id
}
`, template, data.RandomInteger)
}

func (r StorageQueueResource) basicDeprecated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_queue" "test" {
  name                 = "acctestmysamplequeue-%d"
  storage_account_name = azurerm_storage_account.test.name
}
`, template, data.RandomInteger)
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

func (r StorageQueueResource) basicAzureADAuthDeprecated(data acceptance.TestData) string {
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
  name                 = "acctestmysamplequeue-%d"
  storage_account_name = azurerm_storage_account.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func (r StorageQueueResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_queue" "import" {
  name               = azurerm_storage_queue.test.name
  storage_account_id = azurerm_storage_queue.test.storage_account_id
}
`, template)
}

func (r StorageQueueResource) requiresImportDeprecated(data acceptance.TestData) string {
	template := r.basicDeprecated(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_queue" "import" {
  name                 = azurerm_storage_queue.test.name
  storage_account_name = azurerm_storage_queue.test.storage_account_name
}
`, template)
}

func (r StorageQueueResource) metaData(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_queue" "test" {
  name               = "acctestmysamplequeue-%d"
  storage_account_id = azurerm_storage_account.test.id

  metadata = {
    hello = "world"
  }
}
`, template, data.RandomInteger)
}

func (r StorageQueueResource) metaDataDeprecated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_queue" "test" {
  name                 = "acctestmysamplequeue-%d"
  storage_account_name = azurerm_storage_account.test.name

  metadata = {
    hello = "world"
  }
}
`, template, data.RandomInteger)
}

func (r StorageQueueResource) metaDataUpdated(data acceptance.TestData) string {
	template := r.template(data)
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
`, template, data.RandomInteger)
}

func (r StorageQueueResource) metaDataUpdatedDeprecated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_queue" "test" {
  name                 = "acctestmysamplequeue-%d"
  storage_account_name = azurerm_storage_account.test.name

  metadata = {
    hello = "world"
    rick  = "M0rty"
  }
}
`, template, data.RandomInteger)
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

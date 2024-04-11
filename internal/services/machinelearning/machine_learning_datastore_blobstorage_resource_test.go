// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package machinelearning_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-10-01/datastore"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MachineLearningDataStoreBlobStorage struct{}

func TestAccMachineLearningDataStoreBlobStorage_accountKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_datastore_blobstorage", "test")
	r := MachineLearningDataStoreBlobStorage{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.blobStorageAccountKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("account_key"),
	})
}

func TestAccMachineLearningDataStoreBlobStorage_sasToken(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_datastore_blobstorage", "test")
	r := MachineLearningDataStoreBlobStorage{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.blobStorageSas(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("shared_access_signature"),
	})
}

func TestAccMachineLearningDataStoreBlobStorage_Update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_datastore_blobstorage", "test")
	r := MachineLearningDataStoreBlobStorage{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.blobStorageAccountKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("account_key"),
		{
			Config: r.blobStorageSas(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("account_key", "shared_access_signature"),
	})
}

func TestAccMachineLearningDataStoreBlobStorage_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_datastore_blobstorage", "test")
	r := MachineLearningDataStoreBlobStorage{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.blobStorageAccountKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r MachineLearningDataStoreBlobStorage) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	dataStoreClient := client.MachineLearning.Datastore
	id, err := datastore.ParseDataStoreID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := dataStoreClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Machine Learning Data Store %q: %+v", state.ID, err)
	}

	return utils.Bool(resp.Model.Properties != nil), nil
}

func (r MachineLearningDataStoreBlobStorage) blobStorageAccountKey(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "acctestcontainer%[2]d"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_machine_learning_datastore_blobstorage" "test" {
  name                 = "accdatastore%[2]d"
  workspace_id         = azurerm_machine_learning_workspace.test.id
  storage_container_id = azurerm_storage_container.test.resource_manager_id
  account_key          = azurerm_storage_account.test.primary_access_key
}
`, template, data.RandomInteger)
}

func (r MachineLearningDataStoreBlobStorage) blobStorageSas(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "acctestcontainer%[2]d"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

data "azurerm_storage_account_sas" "test" {
  connection_string = azurerm_storage_account.test.primary_connection_string
  https_only        = true
  signed_version    = "2019-10-10"

  resource_types {
    service   = true
    container = true
    object    = true
  }

  services {
    blob  = true
    queue = false
    table = false
    file  = true
  }

  start  = "2022-01-01T06:17:07Z"
  expiry = "2024-12-23T06:17:07Z"

  permissions {
    read    = true
    write   = true
    delete  = false
    list    = false
    add     = true
    create  = true
    update  = false
    process = false
    tag     = false
    filter  = false
  }
}

resource "azurerm_machine_learning_datastore_blobstorage" "test" {
  name                    = "accdatastore%[2]d"
  workspace_id            = azurerm_machine_learning_workspace.test.id
  storage_container_id    = azurerm_storage_container.test.resource_manager_id
  shared_access_signature = data.azurerm_storage_account_sas.test.sas
}
`, template, data.RandomInteger)
}

func (r MachineLearningDataStoreBlobStorage) requiresImport(data acceptance.TestData) string {
	template := r.blobStorageAccountKey(data)
	return fmt.Sprintf(`


%s

resource "azurerm_machine_learning_datastore_blobstorage" "import" {
  name                 = azurerm_machine_learning_datastore_blobstorage.test.name
  workspace_id         = azurerm_machine_learning_datastore_blobstorage.test.workspace_id
  storage_container_id = azurerm_machine_learning_datastore_blobstorage.test.storage_container_id
  account_key          = azurerm_machine_learning_datastore_blobstorage.test.account_key
}
`, template)
}

func (r MachineLearningDataStoreBlobStorage) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-ml-%[1]d"
  location = "%[2]s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestai-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestvault%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name = "standard"

  purge_protection_enabled = true
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Get",
    "Delete",
    "Purge",
  ]
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[4]d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_machine_learning_workspace" "test" {
  name                    = "acctest-MLW-%[1]d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  application_insights_id = azurerm_application_insights.test.id
  key_vault_id            = azurerm_key_vault.test.id
  storage_account_id      = azurerm_storage_account.test.id

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomIntOfLength(15))
}

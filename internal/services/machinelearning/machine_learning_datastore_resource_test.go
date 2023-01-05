package machinelearning_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2022-05-01/datastore"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MachineLearningDataStore struct{}

func TestAccMachineLearningDataStore_blobStorageAccountKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_datastore", "test")
	r := MachineLearningDataStore{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.blobStorageAccountKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("credentials.0.account_key"),
	})
}

func TestAccMachineLearningDataStore_blobStorageSasToken(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_datastore", "test")
	r := MachineLearningDataStore{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.blobStorageSas(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("credentials.0.shared_access_signature"),
	})
}

func TestAccMachineLearningDataStore_fileShareAccountKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_datastore", "test")
	r := MachineLearningDataStore{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.fileShareAccountKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("credentials.0.account_key"),
	})
}

func TestAccMachineLearningDataStore_fileShareSas(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_datastore", "test")
	r := MachineLearningDataStore{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.fileShareSas(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("credentials.0.shared_access_signature"),
	})
}

//func TestAccMachineLearningDataStore_dataLakeGen1Basic(t *testing.T) {
//	data := acceptance.BuildTestData(t, "azurerm_machine_learning_datastore", "test")
//	r := MachineLearningDataStore{}
//
//	data.ResourceTest(t, r, []acceptance.TestStep{
//		{
//			Config: r.dataLakeGen1Basic(data),
//			Check: acceptance.ComposeTestCheckFunc(
//				check.That(data.ResourceName).ExistsInAzure(r),
//			),
//		},
//		data.ImportStep(),
//	})
//}
//
//func TestAccMachineLearningDataStore_dataLakeGen1Spn(t *testing.T) {
//	data := acceptance.BuildTestData(t, "azurerm_machine_learning_datastore", "test")
//	r := MachineLearningDataStore{}
//
//	data.ResourceTest(t, r, []acceptance.TestStep{
//		{
//			Config: r.dataLakeGen1Spn(data),
//			Check: acceptance.ComposeTestCheckFunc(
//				check.That(data.ResourceName).ExistsInAzure(r),
//			),
//		},
//		data.ImportStep("credentials.0.client_secret"),
//	})
//}

func TestAccMachineLearningDataStore_dataLakeGen2Basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_machine_learning_datastore", "test")
	r := MachineLearningDataStore{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dataLakeGen2Basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMachineLearningDataStore_dataLakeGen2WithSpn(t *testing.T) {
	if os.Getenv("ARM_TENANT_ID") == "" || os.Getenv("ARM_CLIENT_ID") == "" ||
		os.Getenv("ARM_CLIENT_SECRET") == "" {
		t.Skip("Skipping as `ARM_TENANT_ID` or `ARM_CLIENT_ID` or `ARM_CLIENT_SECRET` not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_machine_learning_datastore", "test")
	r := MachineLearningDataStore{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dataLakeGen2WithSpn(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("credentials.0.client_secret"),
	})
}

func (r MachineLearningDataStore) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	dataStoreClient := client.MachineLearning.DatastoreClient
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

func (r MachineLearningDataStore) blobStorageAccountKey(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "acctestcontainer%[2]d"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_machine_learning_datastore" "test" {
  name                 = "accdatastore%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  workspace_name       = azurerm_machine_learning_workspace.test.name
  type                 = "AzureBlob"
  storage_account_name = azurerm_storage_account.test.name
  container_name       = azurerm_storage_container.test.name

  credentials {
    account_key = azurerm_storage_account.test.primary_access_key
  }
}
`, template, data.RandomInteger)
}

func (r MachineLearningDataStore) blobStorageSas(data acceptance.TestData) string {
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

resource "azurerm_machine_learning_datastore" "test" {
  name                 = "accdatastore%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  workspace_name       = azurerm_machine_learning_workspace.test.name
  type                 = "AzureBlob"
  storage_account_name = azurerm_storage_account.test.name
  container_name       = azurerm_storage_container.test.name

  credentials {
    shared_access_signature = data.azurerm_storage_account_sas.test.sas
  }
}


`, template, data.RandomInteger)
}

func (r MachineLearningDataStore) fileShareAccountKey(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share" "test" {
  name                 = "accfs%[2]d"
  storage_account_name = azurerm_storage_account.test.name
  quota                = 1
}

resource "azurerm_machine_learning_datastore" "test" {
  name                 = "accdatastore%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  workspace_name       = azurerm_machine_learning_workspace.test.name
  type                 = "AzureFile"
  storage_account_name = azurerm_storage_account.test.name
  file_share_name      = azurerm_storage_share.test.name

  credentials {
    account_key = azurerm_storage_account.test.primary_access_key
  }
}


`, template, data.RandomInteger)
}

func (r MachineLearningDataStore) fileShareSas(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_share" "test" {
  name                 = "accfs%[2]d"
  storage_account_name = azurerm_storage_account.test.name
  quota                = 1
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

resource "azurerm_machine_learning_datastore" "test" {
  name                 = "accdatastore%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  workspace_name       = azurerm_machine_learning_workspace.test.name
  type                 = "AzureFile"
  storage_account_name = azurerm_storage_account.test.name
  file_share_name      = azurerm_storage_share.test.name


  credentials {
    shared_access_signature = data.azurerm_storage_account_sas.test.sas
  }
}


`, template, data.RandomInteger)
}

func (r MachineLearningDataStore) dataLakeGen1Basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s


resource "azurerm_machine_learning_datastore" "test" {
  name                = "accdatastore%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  workspace_name      = azurerm_machine_learning_workspace.test.name
  type                = "AzureDataLakeGen1"

  store_name = "xz3gen1"

}
`, template, data.RandomInteger)
}

func (r MachineLearningDataStore) dataLakeGen1Spn(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s


resource "azurerm_machine_learning_datastore" "test" {
  name                = "accdatastore%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  workspace_name      = azurerm_machine_learning_workspace.test.name
  type                = "AzureDataLakeGen1"

  store_name = "xz3gen1"
  credentials {
    tenant_id     = "72f988bf-86f1-41af-91ab-2d7cd011db47"
    client_id     = "b78b19bf-be33-40ad-98df-a56fe3d5860a"
    client_secret = "R-U7Q~nCvXXWav3xRzwyd-sldhyfgL~3bsDcn"
  }
}
`, template, data.RandomInteger)
}

func (r MachineLearningDataStore) dataLakeGen2Basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "acctestcontainer%[2]d"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_machine_learning_datastore" "test" {
  name                 = "accdatastore%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  workspace_name       = azurerm_machine_learning_workspace.test.name
  type                 = "AzureDataLakeGen2"
  storage_account_name = azurerm_storage_account.test.name
  container_name       = azurerm_storage_container.test.name
}
`, template, data.RandomInteger)
}

func (r MachineLearningDataStore) dataLakeGen2WithSpn(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "acctestcontainer%[2]d"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_machine_learning_datastore" "test" {
  name                 = "accdatastore%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  workspace_name       = azurerm_machine_learning_workspace.test.name
  type                 = "AzureDataLakeGen2"
  storage_account_name = azurerm_storage_account.test.name
  container_name       = azurerm_storage_container.test.name

  credentials {
    tenant_id     = "%[3]s"
    client_id     = "%[4]s"
    client_secret = "%[5]s"
  }
}




`, template, data.RandomInteger, os.Getenv("ARM_TENANT_ID"), os.Getenv("ARM_CLIENT_ID"), os.Getenv("ARM_CLIENT_SECRET"))
}

func (r MachineLearningDataStore) template(data acceptance.TestData) string {
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

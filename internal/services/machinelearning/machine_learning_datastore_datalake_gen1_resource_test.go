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

type MachineLearningDataStoreDataLakeGen1 struct{}

func TestAccMachineLearningDataStoreDataLakeGen1_basic(t *testing.T) {
	if os.Getenv("ARM_DATA_LAKE_GEN1") == "" {
		t.Skip("Skipping as `ARM_DATA_LAKE_GEN1` not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_machine_learning_datastore_datalake_gen1", "test")
	r := MachineLearningDataStoreDataLakeGen1{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dataLakeGen1Basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMachineLearningDataStoreDataLakeGen1_spn(t *testing.T) {
	if os.Getenv("ARM_TENANT_ID") == "" || os.Getenv("ARM_CLIENT_ID") == "" ||
		os.Getenv("ARM_CLIENT_SECRET") == "" || os.Getenv("ARM_DATA_LAKE_GEN1") == "" {
		t.Skip("Skipping as `ARM_TENANT_ID` or `ARM_CLIENT_ID` or `ARM_CLIENT_SECRET` or `ARM_DATA_LAKE_GEN1` not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_machine_learning_datastore_datalake_gen1", "test")
	r := MachineLearningDataStoreDataLakeGen1{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dataLakeGen1WithSpn(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("client_secret"),
	})
}

func TestAccMachineLearningDataStoreDataLakeGen1_requiresImport(t *testing.T) {
	if os.Getenv("ARM_DATA_LAKE_GEN1") == "" {
		t.Skip("Skipping as `ARM_DATA_LAKE_GEN1` not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_machine_learning_datastore_datalake_gen1", "test")
	r := MachineLearningDataStoreDataLakeGen1{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dataLakeGen1Basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r MachineLearningDataStoreDataLakeGen1) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r MachineLearningDataStoreDataLakeGen1) dataLakeGen1Basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "acctestcontainer%[2]d"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_machine_learning_datastore_datalake_gen1" "test" {
  name         = "accdatastore%[2]d"
  workspace_id = azurerm_machine_learning_workspace.test.id
  store_name   = "%[3]s"
}
`, template, data.RandomInteger, os.Getenv("ARM_DATA_LAKE_GEN1"))
}

func (r MachineLearningDataStoreDataLakeGen1) dataLakeGen1WithSpn(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "acctestcontainer%[2]d"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_machine_learning_datastore_datalake_gen1" "test" {
  name         = "accdatastore%[2]d"
  workspace_id = azurerm_machine_learning_workspace.test.id
  store_name   = "%[3]s"

  tenant_id     = "%[4]s"
  client_id     = "%[5]s"
  client_secret = "%[6]s"
}


`, template, data.RandomInteger, os.Getenv("ARM_DATA_LAKE_GEN1"), os.Getenv("ARM_TENANT_ID"), os.Getenv("ARM_CLIENT_ID"), os.Getenv("ARM_CLIENT_SECRET"))
}

func (r MachineLearningDataStoreDataLakeGen1) requiresImport(data acceptance.TestData) string {
	template := r.dataLakeGen1Basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_machine_learning_datastore_datalake_gen1" "import" {
  name         = azurerm_machine_learning_datastore_datalake_gen1.test.name
  workspace_id = azurerm_machine_learning_datastore_datalake_gen1.test.workspace_id
  store_name   = azurerm_machine_learning_datastore_datalake_gen1.test.store_name
}
`, template)
}

func (r MachineLearningDataStoreDataLakeGen1) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
    resource_group {
      prevent_deletion_if_contains_resources = false
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

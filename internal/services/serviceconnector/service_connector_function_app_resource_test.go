// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package serviceconnector_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicelinker/2024-04-01/servicelinker"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type FunctionAppConnectorResource struct{}

func (r FunctionAppConnectorResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := servicelinker.ParseScopedLinkerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.ServiceConnector.ServiceLinkerClient.LinkerGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return utils.Bool(true), nil
}

func TestAccServiceConnectorFunctionAppCosmosdb_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_connection", "test")
	r := FunctionAppConnectorResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cosmosdbBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceConnectorFunctionAppCosmosdb_secretAuth(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_connection", "test")
	r := FunctionAppConnectorResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cosmosdbSecretAuth(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("authentication"),
	})
}

func TestAccServiceConnectorFunctionAppCosmosdb_servicePrincipalSecretAuth(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_connection", "test")
	r := FunctionAppConnectorResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cosmosdbServicePrincipalSecretAuth(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("authentication"),
	})
}

func TestAccServiceConnectorFunctionAppCosmosdb_userAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_connection", "test")
	r := FunctionAppConnectorResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cosmosdbWithUserAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("authentication"),
	})
}

func TestAccServiceConnectorFunctionAppStorageBlob_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_connection", "test")
	r := FunctionAppConnectorResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageBlob(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceConnectorFunctionAppStorageBlob_secretStore(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_connection", "test")
	r := FunctionAppConnectorResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.secretStore(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceConnectorFunctionAppCosmosdb_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_connection", "test")
	r := FunctionAppConnectorResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cosmosdbBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.cosmosdbUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccServiceConnectorFunctionApp_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_connection", "test")
	r := FunctionAppConnectorResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r FunctionAppConnectorResource) storageBlob(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[3]d"
  location = "%[1]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%[2]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[3]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%[3]d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  lifecycle {
    ignore_changes = [
      identity,
    ]
  }
}

resource "azurerm_function_app_connection" "test" {
  name               = "acctestserviceconnector%[3]d"
  function_app_id    = azurerm_function_app.test.id
  target_resource_id = azurerm_storage_account.test.id
  authentication {
    type = "systemAssignedIdentity"
  }
}
`, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func (r FunctionAppConnectorResource) secretStore(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      recover_soft_deleted_key_vaults    = false
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[2]d"
  location = "%[1]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_key_vault" "test" {
  name                     = "accAKV-%[3]s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  purge_protection_enabled = true
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%[2]d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  lifecycle {
    ignore_changes = [
      identity,
    ]
  }
}

resource "azurerm_function_app_connection" "test" {
  name               = "acctestserviceconnector%[2]d"
  function_app_id    = azurerm_function_app.test.id
  target_resource_id = azurerm_storage_account.test.id

  secret_store {
    key_vault_id = azurerm_key_vault.test.id
  }
  authentication {
    type = "systemAssignedIdentity"
  }
}
`, data.Locations.Primary, data.RandomInteger, data.RandomString)
}

func (r FunctionAppConnectorResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_function_app_connection" "test" {
  name               = "acctestserviceconnector%[2]d"
  function_app_id    = azurerm_function_app.test.id
  target_resource_id = azurerm_cosmosdb_account.test.id
  client_type        = "java"
  authentication {
    type = "systemAssignedIdentity"
  }
}
`, template, data.RandomInteger)
}

func (r FunctionAppConnectorResource) cosmosdbBasic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_function_app_connection" "test" {
  name               = "acctestserviceconnector%[2]d"
  function_app_id    = azurerm_function_app.test.id
  target_resource_id = azurerm_cosmosdb_account.test.id
  authentication {
    type = "systemAssignedIdentity"
  }
}
`, template, data.RandomInteger)
}

func (r FunctionAppConnectorResource) cosmosdbSecretAuth(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_function_app_connection" "test" {
  name               = "acctestserviceconnector%[2]d"
  function_app_id    = azurerm_function_app.test.id
  target_resource_id = azurerm_cosmosdb_account.test.id
  authentication {
    type   = "secret"
    name   = "foo"
    secret = "bar"
  }
}
`, template, data.RandomInteger)
}

func (r FunctionAppConnectorResource) cosmosdbServicePrincipalSecretAuth(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_function_app_connection" "test" {
  name               = "acctestserviceconnector%[3]d"
  function_app_id    = azurerm_function_app.test.id
  target_resource_id = azurerm_cosmosdb_account.test.id
  authentication {
    type         = "servicePrincipalSecret"
    client_id    = "someclientid"
    principal_id = azurerm_user_assigned_identity.test.principal_id
    secret       = "bar"
  }
}
`, template, data.RandomString, data.RandomInteger)
}

func (r FunctionAppConnectorResource) cosmosdbWithUserAssignedIdentity(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

data "azurerm_subscription" "test" {}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_function_app_connection" "test" {
  name               = "acctestserviceconnector%[3]d"
  function_app_id    = azurerm_function_app.test.id
  target_resource_id = azurerm_cosmosdb_account.test.id
  authentication {
    type            = "userAssignedIdentity"
    subscription_id = data.azurerm_subscription.test.subscription_id
    client_id       = azurerm_user_assigned_identity.test.client_id
  }
}
`, template, data.RandomString, data.RandomInteger)
}

func (r FunctionAppConnectorResource) cosmosdbUpdate(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_sql_database" "update" {
  name                = "cosmos-sql-db-update"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_cosmosdb_account.test.name
  throughput          = 400
}

resource "azurerm_cosmosdb_sql_container" "update" {
  name                = "test-containerupdate%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_sql_database.test.name
  partition_key_paths = ["/definitionupdate"]
}

resource "azurerm_service_plan" "update" {
  location            = azurerm_resource_group.test.location
  name                = "testserviceplanupdate%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "P1v2"
  os_type             = "Linux"
}

resource "azurerm_function_app_connection" "test" {
  name               = "acctestserviceconnector%[2]d"
  function_app_id    = azurerm_function_app.test.id
  target_resource_id = azurerm_cosmosdb_account.test.id
  authentication {
    type = "systemAssignedIdentity"
  }
}
`, template, data.RandomInteger, data.RandomString)
}

func (r FunctionAppConnectorResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctestcosmosdb%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  consistency_policy {
    consistency_level       = "BoundedStaleness"
    max_interval_in_seconds = 10
    max_staleness_prefix    = 200
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}

resource "azurerm_cosmosdb_sql_database" "test" {
  name                = "cosmos-sql-db"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_cosmosdb_account.test.name
  throughput          = 400
}

resource "azurerm_cosmosdb_sql_container" "test" {
  name                = "test-container%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_sql_database.test.name
  partition_key_paths = ["/definition"]
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_function_app" "test" {
  name                       = "acctest-%[1]d-func"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_plan_id        = azurerm_app_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  lifecycle {
    ignore_changes = [
      identity,
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

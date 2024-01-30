// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cosmos_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CosmosDbSQLRoleDefinitionResource struct{}

func TestAccCosmosDbSQLRoleDefinition_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_sql_role_definition", "test")
	r := CosmosDbSQLRoleDefinitionResource{}

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

func TestAccCosmosDbSQLRoleDefinition_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_sql_role_definition", "test")
	r := CosmosDbSQLRoleDefinitionResource{}

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

func TestAccCosmosDbSQLRoleDefinition_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_sql_role_definition", "test")
	r := CosmosDbSQLRoleDefinitionResource{}
	roleDefinitionId := uuid.New().String()

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, roleDefinitionId),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDbSQLRoleDefinition_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_sql_role_definition", "test")
	r := CosmosDbSQLRoleDefinitionResource{}
	roleDefinitionId := uuid.New().String()

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, roleDefinitionId),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data, roleDefinitionId),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data, roleDefinitionId),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDbSQLRoleDefinition_multiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_sql_role_definition", "test")
	r := CosmosDbSQLRoleDefinitionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multiple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r CosmosDbSQLRoleDefinitionResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.SqlRoleDefinitionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Cosmos.SqlResourceClient.GetSQLRoleDefinition(ctx, id.Name, id.ResourceGroup, id.DatabaseAccountName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %q: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r CosmosDbSQLRoleDefinitionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%d"
  location = "%s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-cosmos-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  consistency_policy {
    consistency_level = "Strong"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r CosmosDbSQLRoleDefinitionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_sql_role_definition" "test" {
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_cosmosdb_account.test.name
  name                = "acctestsqlrole%s"
  assignable_scopes   = ["/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/${azurerm_resource_group.test.name}/providers/Microsoft.DocumentDB/databaseAccounts/${azurerm_cosmosdb_account.test.name}/dbs/sales"]

  permissions {
    data_actions = ["Microsoft.DocumentDB/databaseAccounts/sqlDatabases/containers/items/read"]
  }
}
`, r.template(data), data.RandomString)
}

func (r CosmosDbSQLRoleDefinitionResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_sql_role_definition" "import" {
  role_definition_id  = azurerm_cosmosdb_sql_role_definition.test.role_definition_id
  resource_group_name = azurerm_cosmosdb_sql_role_definition.test.resource_group_name
  account_name        = azurerm_cosmosdb_sql_role_definition.test.account_name
  name                = azurerm_cosmosdb_sql_role_definition.test.name
  assignable_scopes   = azurerm_cosmosdb_sql_role_definition.test.assignable_scopes

  permissions {
    data_actions = ["Microsoft.DocumentDB/databaseAccounts/sqlDatabases/containers/items/read"]
  }
}
`, r.basic(data))
}

func (r CosmosDbSQLRoleDefinitionResource) complete(data acceptance.TestData, roleDefinitionId string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_sql_role_definition" "test" {
  role_definition_id  = "%s"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_cosmosdb_account.test.name
  name                = "acctestsqlrole%s"
  type                = "BuiltInRole"
  assignable_scopes   = ["/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/${azurerm_resource_group.test.name}/providers/Microsoft.DocumentDB/databaseAccounts/${azurerm_cosmosdb_account.test.name}/dbs/sales"]

  permissions {
    data_actions = ["Microsoft.DocumentDB/databaseAccounts/sqlDatabases/containers/items/read"]
  }
}
`, r.template(data), roleDefinitionId, data.RandomString)
}

func (r CosmosDbSQLRoleDefinitionResource) update(data acceptance.TestData, roleDefinitionId string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_sql_role_definition" "test" {
  role_definition_id  = "%s"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_cosmosdb_account.test.name
  name                = "acctestsqlrole2%s"
  type                = "BuiltInRole"
  assignable_scopes   = ["/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/${azurerm_resource_group.test.name}/providers/Microsoft.DocumentDB/databaseAccounts/${azurerm_cosmosdb_account.test.name}/dbs/purchases"]

  permissions {
    data_actions = ["Microsoft.DocumentDB/databaseAccounts/sqlDatabases/containers/items/*"]
  }
}
`, r.template(data), roleDefinitionId, data.RandomString)
}

func (r CosmosDbSQLRoleDefinitionResource) multiple(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_sql_role_definition" "test" {
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_cosmosdb_account.test.name
  name                = "acctestsqlrole%s"
  assignable_scopes   = ["/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/${azurerm_resource_group.test.name}/providers/Microsoft.DocumentDB/databaseAccounts/${azurerm_cosmosdb_account.test.name}/dbs/sales"]

  permissions {
    data_actions = ["Microsoft.DocumentDB/databaseAccounts/sqlDatabases/containers/items/read"]
  }
}

resource "azurerm_cosmosdb_sql_role_definition" "test2" {
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_cosmosdb_account.test.name
  name                = "acctestsqlrole2%s"
  assignable_scopes   = ["/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/${azurerm_resource_group.test.name}/providers/Microsoft.DocumentDB/databaseAccounts/${azurerm_cosmosdb_account.test.name}/dbs/sales"]

  permissions {
    data_actions = ["Microsoft.DocumentDB/databaseAccounts/sqlDatabases/containers/items/read"]
  }
}
`, r.template(data), data.RandomString, data.RandomString)
}

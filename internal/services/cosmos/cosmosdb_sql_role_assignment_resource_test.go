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

type CosmosDbSQLRoleAssignmentResource struct{}

func TestAccCosmosDbSQLRoleAssignment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_sql_role_assignment", "test")
	r := CosmosDbSQLRoleAssignmentResource{}

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

func TestAccCosmosDbSQLRoleAssignment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_sql_role_assignment", "test")
	r := CosmosDbSQLRoleAssignmentResource{}

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

func TestAccCosmosDbSQLRoleAssignment_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_sql_role_assignment", "test")
	r := CosmosDbSQLRoleAssignmentResource{}
	roleAssignmentId := uuid.New().String()

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, roleAssignmentId),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data, roleAssignmentId),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data, roleAssignmentId),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDbSQLRoleAssignment_multiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_sql_role_assignment", "test")
	r := CosmosDbSQLRoleAssignmentResource{}

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

func (r CosmosDbSQLRoleAssignmentResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.SqlRoleAssignmentID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Cosmos.SqlResourceClient.GetSQLRoleAssignment(ctx, id.Name, id.ResourceGroup, id.DatabaseAccountName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %q: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r CosmosDbSQLRoleAssignmentResource) template(data acceptance.TestData) string {
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

resource "azurerm_cosmosdb_sql_role_definition" "test" {
  name                = "acctestsqlrole%s"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_cosmosdb_account.test.name
  type                = "CustomRole"
  assignable_scopes   = ["/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/${azurerm_resource_group.test.name}/providers/Microsoft.DocumentDB/databaseAccounts/${azurerm_cosmosdb_account.test.name}"]

  permissions {
    data_actions = ["Microsoft.DocumentDB/databaseAccounts/sqlDatabases/containers/items/read"]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString)
}

func (r CosmosDbSQLRoleAssignmentResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_sql_role_assignment" "test" {
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_cosmosdb_account.test.name
  role_definition_id  = azurerm_cosmosdb_sql_role_definition.test.id
  principal_id        = data.azurerm_client_config.current.object_id
  scope               = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/${azurerm_resource_group.test.name}/providers/Microsoft.DocumentDB/databaseAccounts/${azurerm_cosmosdb_account.test.name}"
}
`, r.template(data))
}

func (r CosmosDbSQLRoleAssignmentResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_sql_role_assignment" "import" {
  name                = azurerm_cosmosdb_sql_role_assignment.test.name
  resource_group_name = azurerm_cosmosdb_sql_role_assignment.test.resource_group_name
  account_name        = azurerm_cosmosdb_sql_role_assignment.test.account_name
  role_definition_id  = azurerm_cosmosdb_sql_role_assignment.test.role_definition_id
  principal_id        = azurerm_cosmosdb_sql_role_assignment.test.principal_id
  scope               = azurerm_cosmosdb_sql_role_assignment.test.scope
}
`, r.basic(data))
}

func (r CosmosDbSQLRoleAssignmentResource) complete(data acceptance.TestData, roleAssignmentId string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_sql_role_assignment" "test" {
  name                = "%s"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_cosmosdb_account.test.name
  role_definition_id  = azurerm_cosmosdb_sql_role_definition.test.id
  principal_id        = data.azurerm_client_config.current.object_id
  scope               = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/${azurerm_resource_group.test.name}/providers/Microsoft.DocumentDB/databaseAccounts/${azurerm_cosmosdb_account.test.name}"
}
`, r.template(data), roleAssignmentId)
}

func (r CosmosDbSQLRoleAssignmentResource) update(data acceptance.TestData, roleAssignmentId string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_sql_role_definition" "test2" {
  name                = "acctestsqlrole2%s"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_cosmosdb_account.test.name
  type                = "CustomRole"
  assignable_scopes   = ["/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/${azurerm_resource_group.test.name}/providers/Microsoft.DocumentDB/databaseAccounts/${azurerm_cosmosdb_account.test.name}"]

  permissions {
    data_actions = ["Microsoft.DocumentDB/databaseAccounts/sqlDatabases/containers/items/*"]
  }
}

resource "azurerm_cosmosdb_sql_role_assignment" "test" {
  name                = "%s"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_cosmosdb_account.test.name
  role_definition_id  = azurerm_cosmosdb_sql_role_definition.test2.id
  principal_id        = data.azurerm_client_config.current.object_id
  scope               = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/${azurerm_resource_group.test.name}/providers/Microsoft.DocumentDB/databaseAccounts/${azurerm_cosmosdb_account.test.name}"
}
`, r.template(data), data.RandomString, roleAssignmentId)
}

func (r CosmosDbSQLRoleAssignmentResource) multiple(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_sql_role_assignment" "test" {
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_cosmosdb_account.test.name
  role_definition_id  = azurerm_cosmosdb_sql_role_definition.test.id
  principal_id        = data.azurerm_client_config.current.object_id
  scope               = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/${azurerm_resource_group.test.name}/providers/Microsoft.DocumentDB/databaseAccounts/${azurerm_cosmosdb_account.test.name}"
}

resource "azurerm_cosmosdb_sql_role_definition" "test2" {
  name                = "acctestsqlrole2%s"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_cosmosdb_account.test.name
  type                = "CustomRole"
  assignable_scopes   = ["/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/${azurerm_resource_group.test.name}/providers/Microsoft.DocumentDB/databaseAccounts/${azurerm_cosmosdb_account.test.name}"]

  permissions {
    data_actions = ["Microsoft.DocumentDB/databaseAccounts/sqlDatabases/containers/items/read"]
  }
}

resource "azurerm_cosmosdb_sql_role_assignment" "test2" {
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_cosmosdb_account.test.name
  role_definition_id  = azurerm_cosmosdb_sql_role_definition.test2.id
  principal_id        = data.azurerm_client_config.current.object_id
  scope               = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/${azurerm_resource_group.test.name}/providers/Microsoft.DocumentDB/databaseAccounts/${azurerm_cosmosdb_account.test.name}"
}
`, r.template(data), data.RandomString)
}

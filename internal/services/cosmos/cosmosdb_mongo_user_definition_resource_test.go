// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cosmos_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2022-11-15/mongorbacs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CosmosMongoUserDefinitionResource struct{}

func TestAccCosmosDbMongoUserDefinition_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_mongo_user_definition", "test")
	r := CosmosMongoUserDefinitionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password"),
	})
}

func TestAccCosmosDbMongoUserDefinition_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_mongo_user_definition", "test")
	r := CosmosMongoUserDefinitionResource{}

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

func TestAccCosmosDbMongoUserDefinition_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_mongo_user_definition", "test")
	r := CosmosMongoUserDefinitionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password"),
	})
}

func TestAccCosmosDbMongoUserDefinition_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_mongo_user_definition", "test")
	r := CosmosMongoUserDefinitionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("password"),
	})
}

func (r CosmosMongoUserDefinitionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := mongorbacs.ParseMongodbUserDefinitionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Cosmos.MongoRBACClient.MongoDBResourcesGetMongoUserDefinition(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %q: %+v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r CosmosMongoUserDefinitionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_mongo_user_definition" "test" {
  cosmos_mongo_database_id = azurerm_cosmosdb_mongo_database.test.id
  username                 = "myUserName-%d"
  password                 = "myPassword-%d"
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r CosmosMongoUserDefinitionResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_mongo_user_definition" "import" {
  cosmos_mongo_database_id = azurerm_cosmosdb_mongo_user_definition.test.cosmos_mongo_database_id
  username                 = azurerm_cosmosdb_mongo_user_definition.test.username
  password                 = azurerm_cosmosdb_mongo_user_definition.test.password
}
`, r.basic(data))
}

func (r CosmosMongoUserDefinitionResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_mongo_role_definition" "test" {
  cosmos_mongo_database_id = azurerm_cosmosdb_mongo_database.test.id
  role_name                = "acctestmongoroledef%d"
}

resource "azurerm_cosmosdb_mongo_user_definition" "test" {
  cosmos_mongo_database_id = azurerm_cosmosdb_mongo_database.test.id
  username                 = "myUserName-%d"
  password                 = "myPassword-%d"
  inherited_role_names     = [azurerm_cosmosdb_mongo_role_definition.test.role_name]
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r CosmosMongoUserDefinitionResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_mongo_role_definition" "test" {
  cosmos_mongo_database_id = azurerm_cosmosdb_mongo_database.test.id
  role_name                = "acctestmongoroledef%d"
}

resource "azurerm_cosmosdb_mongo_role_definition" "test2" {
  cosmos_mongo_database_id = azurerm_cosmosdb_mongo_database.test.id
  role_name                = "acctestmongoroledef2%d"
}

resource "azurerm_cosmosdb_mongo_user_definition" "test" {
  cosmos_mongo_database_id = azurerm_cosmosdb_mongo_database.test.id
  username                 = "myUserName-%d"
  password                 = "myPassword2-%d"
  inherited_role_names     = [azurerm_cosmosdb_mongo_role_definition.test2.role_name, azurerm_cosmosdb_mongo_role_definition.test.role_name]
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r CosmosMongoUserDefinitionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mongouserdef-%d"
  location = "%s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "MongoDB"

  capabilities {
    name = "EnableMongo"
  }

  capabilities {
    name = "EnableMongoRoleBasedAccessControl"
  }

  consistency_policy {
    consistency_level = "Strong"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}

resource "azurerm_cosmosdb_mongo_database" "test" {
  name                = "acctest-mongodb-%d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

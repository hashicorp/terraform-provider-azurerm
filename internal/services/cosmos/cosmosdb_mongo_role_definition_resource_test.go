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

type CosmosMongoRoleDefinitionResource struct{}

func TestAccCosmosDbMongoRoleDefinition_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_mongo_role_definition", "test")
	r := CosmosMongoRoleDefinitionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDbMongoRoleDefinition_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_mongo_role_definition", "test")
	r := CosmosMongoRoleDefinitionResource{}

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

func TestAccCosmosDbMongoRoleDefinition_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_mongo_role_definition", "test")
	r := CosmosMongoRoleDefinitionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDbMongoRoleDefinition_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_mongo_role_definition", "test")
	r := CosmosMongoRoleDefinitionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r CosmosMongoRoleDefinitionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := mongorbacs.ParseMongodbRoleDefinitionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Cosmos.MongoRBACClient.MongoDBResourcesGetMongoRoleDefinition(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %q: %+v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r CosmosMongoRoleDefinitionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_mongo_role_definition" "test" {
  cosmos_mongo_database_id = azurerm_cosmosdb_mongo_database.test.id
  role_name                = "acctestmongoroledef%d"
}
`, r.template(data), data.RandomInteger)
}

func (r CosmosMongoRoleDefinitionResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_mongo_role_definition" "import" {
  cosmos_mongo_database_id = azurerm_cosmosdb_mongo_role_definition.test.cosmos_mongo_database_id
  role_name                = azurerm_cosmosdb_mongo_role_definition.test.role_name
}
`, r.basic(data))
}

func (r CosmosMongoRoleDefinitionResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_mongo_collection" "test" {
  name                = "acctest-mongocoll-%d"
  resource_group_name = azurerm_cosmosdb_mongo_database.test.resource_group_name
  account_name        = azurerm_cosmosdb_mongo_database.test.account_name
  database_name       = azurerm_cosmosdb_mongo_database.test.name

  index {
    keys   = ["_id"]
    unique = true
  }
}

resource "azurerm_cosmosdb_mongo_role_definition" "base" {
  cosmos_mongo_database_id = azurerm_cosmosdb_mongo_database.test.id
  role_name                = "acctestbaseroledef%d"

  depends_on = [azurerm_cosmosdb_mongo_collection.test]
}

resource "azurerm_cosmosdb_mongo_database" "test2" {
  name                = "acctest-mongodb2-%d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
}

resource "azurerm_cosmosdb_mongo_collection" "test2" {
  name                = "acctest-mongocoll2-%d"
  resource_group_name = azurerm_cosmosdb_mongo_database.test2.resource_group_name
  account_name        = azurerm_cosmosdb_mongo_database.test2.account_name
  database_name       = azurerm_cosmosdb_mongo_database.test2.name

  index {
    keys   = ["_id"]
    unique = true
  }
}

resource "azurerm_cosmosdb_mongo_role_definition" "test" {
  cosmos_mongo_database_id = azurerm_cosmosdb_mongo_database.test.id
  role_name                = "acctestmongoroledef%d"
  inherited_role_names     = [azurerm_cosmosdb_mongo_role_definition.base.role_name]

  privilege {
    actions = ["insert", "find"]

    resource {
      collection_name = azurerm_cosmosdb_mongo_collection.test.name
      db_name         = azurerm_cosmosdb_mongo_database.test2.name
    }
  }

  depends_on = [azurerm_cosmosdb_mongo_collection.test2]
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r CosmosMongoRoleDefinitionResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_mongo_collection" "test" {
  name                = "acctest-mongocoll-%d"
  resource_group_name = azurerm_cosmosdb_mongo_database.test.resource_group_name
  account_name        = azurerm_cosmosdb_mongo_database.test.account_name
  database_name       = azurerm_cosmosdb_mongo_database.test.name

  index {
    keys   = ["_id"]
    unique = true
  }
}

resource "azurerm_cosmosdb_mongo_role_definition" "base" {
  cosmos_mongo_database_id = azurerm_cosmosdb_mongo_database.test.id
  role_name                = "acctestbaseroledef%d"

  depends_on = [azurerm_cosmosdb_mongo_collection.test]
}

resource "azurerm_cosmosdb_mongo_role_definition" "base2" {
  cosmos_mongo_database_id = azurerm_cosmosdb_mongo_database.test.id
  role_name                = "acctestbaseroledef2%d"

  depends_on = [azurerm_cosmosdb_mongo_collection.test]
}

resource "azurerm_cosmosdb_mongo_database" "test2" {
  name                = "acctest-mongodb2-%d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
}

resource "azurerm_cosmosdb_mongo_collection" "test2" {
  name                = "acctest-mongocoll2-%d"
  resource_group_name = azurerm_cosmosdb_mongo_database.test2.resource_group_name
  account_name        = azurerm_cosmosdb_mongo_database.test2.account_name
  database_name       = azurerm_cosmosdb_mongo_database.test2.name

  index {
    keys   = ["_id"]
    unique = true
  }
}

resource "azurerm_cosmosdb_mongo_database" "test3" {
  name                = "acctest-mongodb3-%d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
}

resource "azurerm_cosmosdb_mongo_collection" "test3" {
  name                = "acctest-mongocoll3-%d"
  resource_group_name = azurerm_cosmosdb_mongo_database.test3.resource_group_name
  account_name        = azurerm_cosmosdb_mongo_database.test3.account_name
  database_name       = azurerm_cosmosdb_mongo_database.test3.name

  index {
    keys   = ["_id"]
    unique = true
  }
}

resource "azurerm_cosmosdb_mongo_role_definition" "test" {
  cosmos_mongo_database_id = azurerm_cosmosdb_mongo_database.test.id
  role_name                = "acctestmongoroledef%d"
  inherited_role_names     = [azurerm_cosmosdb_mongo_role_definition.base2.role_name, azurerm_cosmosdb_mongo_role_definition.base.role_name]

  privilege {
    actions = ["insert", "find"]

    resource {
      collection_name = azurerm_cosmosdb_mongo_collection.test2.name
      db_name         = azurerm_cosmosdb_mongo_database.test3.name
    }
  }

  privilege {
    actions = ["find"]

    resource {
      collection_name = azurerm_cosmosdb_mongo_collection.test.name
      db_name         = azurerm_cosmosdb_mongo_database.test2.name
    }
  }

  depends_on = [azurerm_cosmosdb_mongo_collection.test2, azurerm_cosmosdb_mongo_collection.test3]
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r CosmosMongoRoleDefinitionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mongoroledef-%d"
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

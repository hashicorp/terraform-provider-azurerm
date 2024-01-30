// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cosmos_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-10-15/documentdb" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CosmosDbSQLFunctionResource struct{}

func TestAccCosmosDbSQLFunction_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_sql_function", "test")
	r := CosmosDbSQLFunctionResource{}
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

func TestAccCosmosDbSQLFunction_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_sql_function", "test")
	r := CosmosDbSQLFunctionResource{}
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

func TestAccCosmosDbSQLFunction_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_sql_function", "test")
	r := CosmosDbSQLFunctionResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r CosmosDbSQLFunctionResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.SqlFunctionID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Cosmos.SqlResourceClient.GetSQLUserDefinedFunction(ctx, id.ResourceGroup, id.DatabaseAccountName, id.SqlDatabaseName, id.ContainerName, id.UserDefinedFunctionName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving CosmosDb SQLFunction %q: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r CosmosDbSQLFunctionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%[2]d"
  location = "%[1]s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-cosmos-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "%[3]s"

  consistency_policy {
    consistency_level = "%[4]s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}

resource "azurerm_cosmosdb_sql_database" "test" {
  name                = "acctest-sql-database-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
}

resource "azurerm_cosmosdb_sql_container" "test" {
  name                = "acctest-sql-container-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_sql_database.test.name
  partition_key_path  = "/definition/id"
}
	`, data.Locations.Primary, data.RandomInteger, string(documentdb.DatabaseAccountKindGlobalDocumentDB), string(documentdb.DefaultConsistencyLevelSession))
}

func (r CosmosDbSQLFunctionResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_sql_function" "test" {
  name         = "acctest-dssdf-%d"
  container_id = azurerm_cosmosdb_sql_container.test.id
  body         = <<BODY
  	function test() {
		var context = getContext();
		var response = context.getResponse();
		response.setBody('Hello, World');
	}
BODY
}
`, template, data.RandomInteger)
}

func (r CosmosDbSQLFunctionResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_sql_function" "import" {
  name         = azurerm_cosmosdb_sql_function.test.name
  container_id = azurerm_cosmosdb_sql_function.test.container_id
  body         = <<BODY
  	function test() {
		var context = getContext();
		var response = context.getResponse();
		response.setBody('Hello, World');
	}
BODY
}
`, config)
}

func (r CosmosDbSQLFunctionResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_sql_function" "test" {
  name         = "acctest-dssdf-%d"
  container_id = azurerm_cosmosdb_sql_container.test.id
  body         = <<BODY
	function test() {
		var context = getContext();
		var response = context.getResponse();
		response.setBody('Welcome To Terraform');
	}
BODY
}
`, template, data.RandomInteger)
}

package kusto_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2022-12-29/dataconnections"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type KustoCosmosDBDataConnectionResource struct{}

func (k KustoCosmosDBDataConnectionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := dataconnections.ParseDataConnectionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Kusto.DataConnectionsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}

	if resp.Model != nil {
		value, ok := (*resp.Model).(dataconnections.CosmosDbDataConnection)
		if !ok {
			return nil, fmt.Errorf("%s is not an CosmosDB Data Connection", id.String())
		}
		exists := value.Properties != nil
		return &exists, nil
	} else {
		return nil, fmt.Errorf("response model is empty")
	}
}

func TestAccKustoCosmosDBDataConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cosmosdb_data_connection", "test")
	r := KustoCosmosDBDataConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cosmosdb_account_id").Exists(),
				check.That(data.ResourceName).Key("cosmosdb_database").Exists(),
				check.That(data.ResourceName).Key("cosmosdb_container").Exists(),
				check.That(data.ResourceName).Key("cluster_name").Exists(),
				check.That(data.ResourceName).Key("database_name").Exists(),
				check.That(data.ResourceName).Key("table_name").Exists(),
				check.That(data.ResourceName).Key("managed_identity_id").Exists(),
				check.That(data.ResourceName).Key("table_name").HasValue("TestTable"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKustoCosmosDBDataConnection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cosmosdb_data_connection", "test")
	r := KustoCosmosDBDataConnectionResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cosmosdb_account_id").Exists(),
				check.That(data.ResourceName).Key("cosmosdb_database").Exists(),
				check.That(data.ResourceName).Key("cosmosdb_container").Exists(),
				check.That(data.ResourceName).Key("cluster_name").Exists(),
				check.That(data.ResourceName).Key("database_name").Exists(),
				check.That(data.ResourceName).Key("table_name").Exists(),
				check.That(data.ResourceName).Key("managed_identity_id").Exists(),
				check.That(data.ResourceName).Key("table_name").HasValue("TestTable"),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cosmosdb_account_id").Exists(),
				check.That(data.ResourceName).Key("cosmosdb_database").Exists(),
				check.That(data.ResourceName).Key("cosmosdb_container").Exists(),
				check.That(data.ResourceName).Key("cluster_name").Exists(),
				check.That(data.ResourceName).Key("database_name").Exists(),
				check.That(data.ResourceName).Key("table_name").Exists(),
				check.That(data.ResourceName).Key("managed_identity_id").Exists(),
				check.That(data.ResourceName).Key("table_name").HasValue("TestTable"),
				check.That(data.ResourceName).Key("mapping_rule_name").HasValue("TestMappingRule"),
				check.That(data.ResourceName).Key("retrieval_start_date").HasValue("2023-02-29T12:00:00.6554616Z"),
			),
		},
		data.ImportStep(),
	})
}

func (k KustoCosmosDBDataConnectionResource) basic(data acceptance.TestData) string {
	template := k.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_kusto_cosmosdb_data_connection" "test" {
	  name                  = "acctestkcd%s"
	  resource_group_name   = azurerm_resource_group.test.name
      location              = azurerm_resource_group.test.location
      cosmosdb_account_id   = azurerm_cosmosdb_account.test.id
      cosmosdb_database     = azurerm_cosmosdb_sql_database.test.name
      cosmosdb_container    = azurerm_cosmosdb_sql_container.test.name
	  cluster_name          = azurerm_kusto_cluster.test.name	
      database_name         = azurerm_kusto_database.test.name
      managed_identity_id   = azurerm_user_assigned_identity.test.id
      table_name            = "TestTable"
}`, template, data.RandomString)
}

func (k KustoCosmosDBDataConnectionResource) complete(data acceptance.TestData) string {
	template := k.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_kusto_cosmosdb_data_connection" "test" {
	  name                  = "acctestkcd%s"
	  resource_group_name   = azurerm_resource_group.test.name
      location              = azurerm_resource_group.test.location
      cosmosdb_account_id   = azurerm_cosmosdb_account.test.id
      cosmosdb_database     = azurerm_cosmosdb_sql_database.test.name
      cosmosdb_container    = azurerm_cosmosdb_sql_container.test.name
	  cluster_name          = azurerm_kusto_cluster.test.name	
      database_name         = azurerm_kusto_database.test.name
      managed_identity_id   = azurerm_user_assigned_identity.test.id
      table_name            = "TestTable"
      mapping_rule_name     = "TestMappingRule"
      retrieval_start_date  = "2023-02-29T12:00:00.6554616Z"
}`, template, data.RandomString)
}

func (k KustoCosmosDBDataConnectionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "acctestuami-%d"
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  consistency_policy {
    consistency_level = "Session"
    max_interval_in_seconds = 5
    max_staleness_prefix = 100
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}

resource "azurerm_cosmosdb_sql_database" "test" {
  name                = "acctestcosmosdbsqldb-%d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
}

resource "azurerm_cosmosdb_sql_container" "test" {
  name                  = "acctestcosmosdbsqlcon-%d"
  resource_group_name   = azurerm_cosmosdb_account.test.resource_group_name
  account_name          = azurerm_cosmosdb_account.test.name
  database_name         = azurerm_cosmosdb_sql_database.test.name
  partition_key_path    = "/part"
  throughput            = 400
}

data "azurerm_cosmosdb_sql_role_definition" "test" {
  role_definition_id = "00000000-0000-0000-0000-000000000001"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_cosmosdb_account.test.name
}


resource "azurerm_cosmosdb_sql_role_assignment" "test" {
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_cosmosdb_account.test.name
  role_definition_id  = azurerm_cosmosdb_sql_role_definition.test.id
  principal_id        = azurerm_user_assigned_identity.test.principal_id
  scope               = azurerm_cosmosdb_account.test.id
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }

  identity {
    type = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}

resource "azurerm_kusto_database" "test" {
  name                = "acctestkd-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_name        = azurerm_kusto_cluster.test.name
}

resource "azurerm_kusto_script" "test" {
  name = "create-table-script"
  database_id = azurerm_kusto_database.test.id
  script_content = ".create table TestTable(Id:string, Name:string, _ts:long, _timestamp:datetime)"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomInteger)
}

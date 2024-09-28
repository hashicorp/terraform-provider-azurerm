// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kusto_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-08-15/dataconnections"
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
		value, ok := resp.Model.(dataconnections.CosmosDbDataConnection)
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
				check.That(data.ResourceName).Key("cosmosdb_container_id").Exists(),
				check.That(data.ResourceName).Key("kusto_database_id").Exists(),
				check.That(data.ResourceName).Key("table_name").Exists(),
				check.That(data.ResourceName).Key("managed_identity_id").Exists(),
				check.That(data.ResourceName).Key("table_name").HasValue("TestTable"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKustoCosmosDBDataConnection_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cosmosdb_data_connection", "test")
	r := KustoCosmosDBDataConnectionResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cosmosdb_container_id").Exists(),
				check.That(data.ResourceName).Key("kusto_database_id").Exists(),
				check.That(data.ResourceName).Key("table_name").Exists(),
				check.That(data.ResourceName).Key("managed_identity_id").Exists(),
				check.That(data.ResourceName).Key("table_name").HasValue("TestTable"),
				check.That(data.ResourceName).Key("mapping_rule_name").HasValue("TestMapping"),
				check.That(data.ResourceName).Key("retrieval_start_date").HasValue("2023-06-26T12:00:00.6554616Z"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKustoCosmosDBDataConnection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cosmosdb_data_connection", "test")
	r := KustoCosmosDBDataConnectionResource{}
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

func TestAccKustoCosmosDBDataConnection_diffResourceGroups(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_cosmosdb_data_connection", "test")
	r := KustoCosmosDBDataConnectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.diffResourceGroups(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("cosmosdb_container_id").Exists(),
				check.That(data.ResourceName).Key("kusto_database_id").Exists(),
				check.That(data.ResourceName).Key("table_name").Exists(),
				check.That(data.ResourceName).Key("managed_identity_id").Exists(),
				check.That(data.ResourceName).Key("table_name").HasValue("TestTable"),
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
  location              = azurerm_resource_group.test.location
  cosmosdb_container_id = azurerm_cosmosdb_sql_container.test.id
  kusto_database_id     = azurerm_kusto_database.test.id
  managed_identity_id   = azurerm_kusto_cluster.test.id
  table_name            = "TestTable"
  lifecycle {
    ignore_changes = [retrieval_start_date]
  }
}`, template, data.RandomString)
}

func (k KustoCosmosDBDataConnectionResource) complete(data acceptance.TestData) string {
	template := k.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_kusto_cosmosdb_data_connection" "test" {
  name                  = "acctestkcd%s"
  location              = azurerm_resource_group.test.location
  cosmosdb_container_id = azurerm_cosmosdb_sql_container.test.id
  kusto_database_id     = azurerm_kusto_database.test.id
  managed_identity_id   = azurerm_kusto_cluster.test.id
  table_name            = "TestTable"
  mapping_rule_name     = "TestMapping"
  retrieval_start_date  = "2023-06-26T12:00:00.6554616Z"
}`, template, data.RandomString)
}

func (k KustoCosmosDBDataConnectionResource) requiresImport(data acceptance.TestData) string {
	template := k.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_kusto_cosmosdb_data_connection" "import" {
  name                  = azurerm_kusto_cosmosdb_data_connection.test.name
  location              = azurerm_resource_group.test.location
  cosmosdb_container_id = azurerm_cosmosdb_sql_container.test.id
  kusto_database_id     = azurerm_kusto_database.test.id
  managed_identity_id   = azurerm_kusto_cluster.test.id
  table_name            = "TestTable"
  lifecycle {
    ignore_changes = [retrieval_start_date]
  }
}`, template)
}

func (k KustoCosmosDBDataConnectionResource) diffResourceGroups(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_resource_group" "test2" {
  name     = "acctestSecRG-%d"
  location = "%s"
}

data "azurerm_role_definition" "builtin" {
  role_definition_id = "fbdf93bf-df7d-467e-a4d2-9458aa1360c8"
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_resource_group.test2.id
  role_definition_name = data.azurerm_role_definition.builtin.name
  principal_id         = azurerm_kusto_cluster.test.identity[0].principal_id
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%d"
  location            = azurerm_resource_group.test2.location
  resource_group_name = azurerm_resource_group.test2.name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  consistency_policy {
    consistency_level       = "Session"
    max_interval_in_seconds = 5
    max_staleness_prefix    = 100
  }

  geo_location {
    location          = azurerm_resource_group.test2.location
    failover_priority = 0
  }
}

resource "azurerm_cosmosdb_sql_database" "test" {
  name                = "acctestcosmosdbsqldb-%d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
}

resource "azurerm_cosmosdb_sql_container" "test" {
  name                = "acctestcosmosdbsqlcon-%d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_sql_database.test.name
  partition_key_paths = ["/part"]
  throughput          = 400
}

data "azurerm_cosmosdb_sql_role_definition" "test" {
  role_definition_id  = "00000000-0000-0000-0000-000000000001"
  resource_group_name = azurerm_resource_group.test2.name
  account_name        = azurerm_cosmosdb_account.test.name
}


resource "azurerm_cosmosdb_sql_role_assignment" "test" {
  resource_group_name = azurerm_resource_group.test2.name
  account_name        = azurerm_cosmosdb_account.test.name
  role_definition_id  = data.azurerm_cosmosdb_sql_role_definition.test.id
  principal_id        = azurerm_kusto_cluster.test.identity[0].principal_id
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
    type = "SystemAssigned"
  }
}

resource "azurerm_kusto_database" "test" {
  name                = "acctestkd-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_name        = azurerm_kusto_cluster.test.name
}

resource "azurerm_kusto_script" "test" {
  name           = "create-table-script"
  database_id    = azurerm_kusto_database.test.id
  script_content = <<SCRIPT
.create table TestTable(Id:string, Name:string, _ts:long, _timestamp:datetime)
.create table TestTable ingestion json mapping "TestMapping"
'['
'    {"column":"Id","path":"$.id"},'
'    {"column":"Name","path":"$.name"},'
'    {"column":"_ts","path":"$._ts"},'
'    {"column":"_timestamp","path":"$._ts", "transform":"DateTimeFromUnixSeconds"}'
']'
.alter table TestTable policy ingestionbatching "{'MaximumBatchingTimeSpan': '0:0:10', 'MaximumNumberOfItems': 10000}"
SCRIPT
}

resource "azurerm_kusto_cosmosdb_data_connection" "test" {
  name                  = "acctestkcd%s"
  location              = azurerm_resource_group.test.location
  cosmosdb_container_id = azurerm_cosmosdb_sql_container.test.id
  kusto_database_id     = azurerm_kusto_database.test.id
  managed_identity_id   = azurerm_kusto_cluster.test.id
  table_name            = "TestTable"
  lifecycle {
    ignore_changes = [retrieval_start_date]
  }
}`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomInteger, data.RandomString)
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

data "azurerm_role_definition" "builtin" {
  role_definition_id = "fbdf93bf-df7d-467e-a4d2-9458aa1360c8"
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_resource_group.test.id
  role_definition_name = data.azurerm_role_definition.builtin.name
  principal_id         = azurerm_kusto_cluster.test.identity[0].principal_id
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  consistency_policy {
    consistency_level       = "Session"
    max_interval_in_seconds = 5
    max_staleness_prefix    = 100
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
  name                = "acctestcosmosdbsqlcon-%d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_sql_database.test.name
  partition_key_paths = ["/part"]
  throughput          = 400
}

data "azurerm_cosmosdb_sql_role_definition" "test" {
  role_definition_id  = "00000000-0000-0000-0000-000000000001"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_cosmosdb_account.test.name
}


resource "azurerm_cosmosdb_sql_role_assignment" "test" {
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_cosmosdb_account.test.name
  role_definition_id  = data.azurerm_cosmosdb_sql_role_definition.test.id
  principal_id        = azurerm_kusto_cluster.test.identity[0].principal_id
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
    type = "SystemAssigned"
  }
}

resource "azurerm_kusto_database" "test" {
  name                = "acctestkd-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_name        = azurerm_kusto_cluster.test.name
}

resource "azurerm_kusto_script" "test" {
  name           = "create-table-script"
  database_id    = azurerm_kusto_database.test.id
  script_content = <<SCRIPT
.create table TestTable(Id:string, Name:string, _ts:long, _timestamp:datetime)
.create table TestTable ingestion json mapping "TestMapping"
'['
'    {"column":"Id","path":"$.id"},'
'    {"column":"Name","path":"$.name"},'
'    {"column":"_ts","path":"$._ts"},'
'    {"column":"_timestamp","path":"$._ts", "transform":"DateTimeFromUnixSeconds"}'
']'
.alter table TestTable policy ingestionbatching "{'MaximumBatchingTimeSpan': '0:0:10', 'MaximumNumberOfItems': 10000}"
SCRIPT
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomInteger)
}

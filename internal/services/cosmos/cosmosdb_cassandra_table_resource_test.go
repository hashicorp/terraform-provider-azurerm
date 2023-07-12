// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cosmos_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CosmosDBCassandraTableResource struct{}

func (t CosmosDBCassandraTableResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.CassandraTableID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Cosmos.CassandraClient.GetCassandraTable(ctx, id.ResourceGroup, id.DatabaseAccountName, id.CassandraKeyspaceName, id.TableName)
	if err != nil {
		return nil, fmt.Errorf("reading Cosmos Cassandra Table (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func TestAccCosmosDbCassandraTable_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_cassandra_table", "test")
	r := CosmosDBCassandraTableResource{}

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

func TestAccCosmosDbCassandraTable_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_cassandra_table", "test")
	r := CosmosDBCassandraTableResource{}

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

func TestAccCosmosDbCassandraTable_autoScaleSetting(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_cassandra_table", "test")
	r := CosmosDBCassandraTableResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoScaleSetting(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDbCassandraTable_serverless(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_cassandra_table", "test")
	r := CosmosDBCassandraTableResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.serverless(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDbCassandraTable_analyticalStorageTTL(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_cassandra_table", "test")
	r := CosmosDBCassandraTableResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.analyticalStorageTTL(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (CosmosDBCassandraTableResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_cassandra_table" "test" {
  name                  = "acctest-CCASST-%[2]d"
  cassandra_keyspace_id = azurerm_cosmosdb_cassandra_keyspace.test.id

  schema {
    column {
      name = "test1"
      type = "ascii"
    }

    column {
      name = "test2"
      type = "int"
    }

    partition_key {
      name = "test1"
    }
  }
}
`, CosmosDbCassandraKeyspaceResource{}.basic(data), data.RandomInteger)
}

func (CosmosDBCassandraTableResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_cassandra_table" "test" {
  name                  = "acctest-CCASST-%[2]d"
  cassandra_keyspace_id = azurerm_cosmosdb_cassandra_keyspace.test.id
  default_ttl           = 100
  throughput            = 400

  schema {
    column {
      name = "test1"
      type = "ascii"
    }

    column {
      name = "test2"
      type = "int"
    }

    partition_key {
      name = "test1"
    }
  }
}
`, CosmosDbCassandraKeyspaceResource{}.basic(data), data.RandomInteger)
}

func (CosmosDBCassandraTableResource) autoScaleSetting(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_cassandra_table" "test" {
  name                  = "acctest-CCASST-%[2]d"
  cassandra_keyspace_id = azurerm_cosmosdb_cassandra_keyspace.test.id
  autoscale_settings {
    max_throughput = 4000
  }

  schema {
    column {
      name = "test1"
      type = "ascii"
    }

    column {
      name = "test2"
      type = "int"
    }

    partition_key {
      name = "test1"
    }
  }
}
`, CosmosDbCassandraKeyspaceResource{}.basic(data), data.RandomInteger)
}

func (CosmosDBCassandraTableResource) analyticalStorageTTLTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%[1]d"
  location = "%[2]s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                       = "acctest-ca-%[1]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  offer_type                 = "Standard"
  kind                       = "GlobalDocumentDB"
  analytical_storage_enabled = true

  consistency_policy {
    consistency_level = "Strong"
  }

  capabilities {
    name = "EnableCassandra"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}

resource "azurerm_cosmosdb_cassandra_keyspace" "test" {
  name                = "acctest-%[1]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r CosmosDBCassandraTableResource) analyticalStorageTTL(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_cassandra_table" "test" {
  name                   = "acctest-CCASST-%[2]d"
  cassandra_keyspace_id  = azurerm_cosmosdb_cassandra_keyspace.test.id
  analytical_storage_ttl = 1

  schema {
    column {
      name = "test1"
      type = "ascii"
    }

    column {
      name = "test2"
      type = "int"
    }

    partition_key {
      name = "test1"
    }
  }
}
`, r.analyticalStorageTTLTemplate(data), data.RandomInteger)
}

func (CosmosDBCassandraTableResource) serverless(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_cassandra_table" "test" {
  name                  = "acctest-CCASST-%[2]d"
  cassandra_keyspace_id = azurerm_cosmosdb_cassandra_keyspace.test.id

  schema {
    column {
      name = "test1"
      type = "ascii"
    }

    column {
      name = "test2"
      type = "int"
    }

    partition_key {
      name = "test1"
    }
  }
}
`, CosmosDbCassandraKeyspaceResource{}.serverless(data), data.RandomInteger)
}

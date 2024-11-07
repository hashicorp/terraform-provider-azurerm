// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cosmos_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CosmosSqlDatabaseResource struct{}

func TestAccCosmosDbSqlDatabase_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_sql_database", "test")
	r := CosmosSqlDatabaseResource{}

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

func TestAccCosmosDbSqlDatabase_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_sql_database", "test")
	r := CosmosSqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.throughput(data, 700),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("throughput").HasValue("700"),
			),
		},
		data.ImportStep(),
		{
			Config: r.throughput(data, 1700),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("throughput").HasValue("1700"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDbSqlDatabase_autoscale(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_sql_database", "test")
	r := CosmosSqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoscale(data, 4000),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("autoscale_settings.0.max_throughput").HasValue("4000"),
			),
		},
		data.ImportStep(),
		{
			Config: r.autoscale(data, 5000),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("autoscale_settings.0.max_throughput").HasValue("5000"),
			),
		},
		data.ImportStep(),
		{
			Config: r.autoscale(data, 4000),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("autoscale_settings.0.max_throughput").HasValue("4000"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDbSqlDatabase_serverless(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_sql_database", "test")
	r := CosmosSqlDatabaseResource{}

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

func (t CosmosSqlDatabaseResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.SqlDatabaseID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Cosmos.SqlClient.GetSQLDatabase(ctx, id.ResourceGroup, id.DatabaseAccountName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading Cosmos SQL Database (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (CosmosSqlDatabaseResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_sql_database" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
}
`, CosmosDBAccountResource{}.basic(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, cosmosdb.DefaultConsistencyLevelStrong), data.RandomInteger)
}

func (CosmosSqlDatabaseResource) throughput(data acceptance.TestData, throughput int) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_sql_database" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  throughput          = %[3]d
}
`, CosmosDBAccountResource{}.basic(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, cosmosdb.DefaultConsistencyLevelStrong), data.RandomInteger, throughput)
}

func (CosmosSqlDatabaseResource) autoscale(data acceptance.TestData, maxThroughput int) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_sql_database" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  autoscale_settings {
    max_throughput = %[3]d
  }
}
`, CosmosDBAccountResource{}.basic(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, cosmosdb.DefaultConsistencyLevelStrong), data.RandomInteger, maxThroughput)
}

func (CosmosSqlDatabaseResource) serverless(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s
resource "azurerm_cosmosdb_sql_database" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
}
`, CosmosDBAccountResource{}.capabilities(data, cosmosdb.DatabaseAccountKindGlobalDocumentDB, []string{"EnableServerless"}), data.RandomInteger)
}

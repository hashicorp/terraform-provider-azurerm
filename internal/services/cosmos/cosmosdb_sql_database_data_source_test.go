// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cosmos_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type CosmosDBSqlDatabaseDataSource struct{}

func TestAccDataSourceCosmosDBSqlDatabase_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_cosmosdb_sql_database", "test")
	r := CosmosDBSqlDatabaseDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
			),
		},
	})
}

func TestAccDataSourceCosmosSqlDatabase_throughput(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_cosmosdb_sql_database", "test")
	r := CosmosDBSqlDatabaseDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.throughput(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
			),
		},
	})
}

func TestAccDataSourceCosmosSqlDatabase_serverless(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_cosmosdb_sql_database", "test")
	r := CosmosDBSqlDatabaseDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.serverless(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
			),
		},
	})
}

func (CosmosDBSqlDatabaseDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_cosmosdb_sql_database" "test" {
  name                = azurerm_cosmosdb_sql_database.test.name
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_cosmosdb_sql_database.test.account_name
}
`, CosmosSqlDatabaseResource{}.basic(data))
}

func (CosmosDBSqlDatabaseDataSource) throughput(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_cosmosdb_sql_database" "test" {
  name                = azurerm_cosmosdb_sql_database.test.name
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_cosmosdb_sql_database.test.account_name
}
`, CosmosSqlDatabaseResource{}.throughput(data, 4000))
}

func (CosmosDBSqlDatabaseDataSource) serverless(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_cosmosdb_sql_database" "test" {
  name                = azurerm_cosmosdb_sql_database.test.name
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_cosmosdb_sql_database.test.account_name
}
`, CosmosSqlDatabaseResource{}.serverless(data))
}

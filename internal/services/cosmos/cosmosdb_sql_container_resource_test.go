// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cosmos_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-05-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CosmosSqlContainerResource struct{}

func TestAccCosmosDbSqlContainer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_sql_container", "test")
	r := CosmosSqlContainerResource{}

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

func TestAccCosmosDbSqlContainer_basic_serverless(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_sql_container", "test")
	r := CosmosSqlContainerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic_serverless(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDbSqlContainer_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_sql_container", "test")
	r := CosmosSqlContainerResource{}

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

func TestAccCosmosDbSqlContainer_analyticalStorageTTL(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_sql_container", "test")
	r := CosmosSqlContainerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.analyticalStorageTTL(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.analyticalStorageTTL_removed(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.analyticalStorageTTL(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDbSqlContainer_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_sql_container", "test")
	r := CosmosSqlContainerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_ttl").HasValue("500"),
				check.That(data.ResourceName).Key("throughput").HasValue("600"),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_ttl").HasValue("1000"),
				check.That(data.ResourceName).Key("throughput").HasValue("400"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDbSqlContainer_autoscale(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_sql_container", "test")
	r := CosmosSqlContainerResource{}

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

func TestAccCosmosDbSqlContainer_indexing_policy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_sql_container", "test")
	r := CosmosSqlContainerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.indexing_policy(data, "/includedPath01/*", "/excludedPath01/?"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.indexing_policy(data, "/includedPath02/*", "/excludedPath02/?"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.indexing_policy_update_spatialIndex(data, "/includedPath02/*", "/excludedPath02/?"),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.indexing_policy_update_includedPath(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDbSqlContainer_partition_key_version(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_sql_container", "test")
	r := CosmosSqlContainerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.partition_key_version(data, 2),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("partition_key_version").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDbSqlContainer_customConflictResolutionPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_sql_container", "test")
	r := CosmosSqlContainerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.conflictResolutionPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDbSqlContainer_hierarchicalPartitionKeys(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_sql_container", "test")
	r := CosmosSqlContainerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.hierarchicalPartitionKeys(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t CosmosSqlContainerResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := cosmosdb.ParseContainerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Cosmos.CosmosDBClient.SqlResourcesGetSqlContainer(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading Cosmos SQL Container (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (CosmosSqlContainerResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_sql_container" "test" {
  name                = "acctest-CSQLC-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_sql_database.test.name
  partition_key_paths = ["/definition/id"]
}
`, CosmosSqlDatabaseResource{}.basic(data), data.RandomInteger)
}

func (CosmosSqlContainerResource) basic_serverless(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_sql_container" "test" {
  name                = "acctest-CSQLC-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_sql_database.test.name
  partition_key_paths = ["/definition/id"]
}
`, CosmosSqlDatabaseResource{}.serverless(data), data.RandomInteger)
}

func (CosmosSqlContainerResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_sql_container" "test" {
  name                = "acctest-CSQLC-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_sql_database.test.name
  partition_key_paths = ["/definition/id"]
  unique_key {
    paths = ["/definition/id1", "/definition/id2"]
  }
  default_ttl = 500
  throughput  = 600
  indexing_policy {
    indexing_mode = "consistent"

    included_path {
      path = "/*"
    }

    included_path {
      path = "/testing/id1/*"
    }

    excluded_path {
      path = "/testing/id2/*"
    }
    composite_index {
      index {
        path  = "/path1"
        order = "descending"
      }
      index {
        path  = "/path2"
        order = "ascending"
      }
    }

    composite_index {
      index {
        path  = "/path3"
        order = "ascending"
      }
      index {
        path  = "/path4"
        order = "descending"
      }
    }
  }
}
`, CosmosSqlDatabaseResource{}.basic(data), data.RandomInteger)
}

func (CosmosSqlContainerResource) analyticalStorageTTL(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_sql_database" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
}

resource "azurerm_cosmosdb_sql_container" "test" {
  name                   = "acctest-CSQLC-%[2]d"
  resource_group_name    = azurerm_cosmosdb_account.test.resource_group_name
  account_name           = azurerm_cosmosdb_account.test.name
  database_name          = azurerm_cosmosdb_sql_database.test.name
  partition_key_paths    = ["/definition/id"]
  analytical_storage_ttl = 600
}
`, CosmosDBAccountResource{}.analyticalStorage(data, "GlobalDocumentDB", "Eventual", true), data.RandomInteger, data.RandomInteger)
}

func (CosmosSqlContainerResource) analyticalStorageTTL_removed(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_sql_database" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
}

resource "azurerm_cosmosdb_sql_container" "test" {
  name                = "acctest-CSQLC-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_sql_database.test.name
  partition_key_paths = ["/definition/id"]
}
`, CosmosDBAccountResource{}.analyticalStorage(data, "GlobalDocumentDB", "Eventual", true), data.RandomInteger, data.RandomInteger)
}

func (CosmosSqlContainerResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_sql_container" "test" {
  name                = "acctest-CSQLC-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_sql_database.test.name
  partition_key_paths = ["/definition/id"]
  unique_key {
    paths = ["/definition/id1", "/definition/id2"]
  }
  default_ttl = 1000
  throughput  = 400
  indexing_policy {
    indexing_mode = "consistent"

    included_path {
      path = "/*"
    }

    included_path {
      path = "/testing/id2/*"
    }

    excluded_path {
      path = "/testing/id1/*"
    }

    composite_index {
      index {
        path  = "/path1"
        order = "ascending"
      }
      index {
        path  = "/path2"
        order = "descending"
      }
    }

    composite_index {
      index {
        path  = "/path3"
        order = "ascending"
      }
      index {
        path  = "/path4"
        order = "descending"
      }
    }
  }
}
`, CosmosSqlDatabaseResource{}.basic(data), data.RandomInteger)
}

func (CosmosSqlContainerResource) autoscale(data acceptance.TestData, maxThroughput int) string {
	return fmt.Sprintf(`
%[1]s
resource "azurerm_cosmosdb_sql_container" "test" {
  name                = "acctest-CSQLC-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_sql_database.test.name
  partition_key_paths = ["/definition/id"]
  autoscale_settings {
    max_throughput = %[3]d
  }
}
`, CosmosSqlDatabaseResource{}.basic(data), data.RandomInteger, maxThroughput)
}

func (CosmosSqlContainerResource) indexing_policy(data acceptance.TestData, includedPath, excludedPath string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_sql_container" "test" {
  name                = "acctest-CSQLC-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_sql_database.test.name
  partition_key_paths = ["/definition/id"]

  indexing_policy {
    indexing_mode = "consistent"

    included_path {
      path = "/*"
    }

    included_path {
      path = "%s"
    }

    excluded_path {
      path = "%s"
    }

    composite_index {
      index {
        path  = "/path1"
        order = "ascending"
      }
      index {
        path  = "/path2"
        order = "descending"
      }
    }

    composite_index {
      index {
        path  = "/path3"
        order = "ascending"
      }
      index {
        path  = "/path4"
        order = "descending"
      }
    }
  }
}
`, CosmosSqlDatabaseResource{}.basic(data), data.RandomInteger, includedPath, excludedPath)
}

func (CosmosSqlContainerResource) indexing_policy_update_spatialIndex(data acceptance.TestData, includedPath, excludedPath string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_sql_container" "test" {
  name                = "acctest-CSQLC-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_sql_database.test.name
  partition_key_paths = ["/definition/id"]

  indexing_policy {
    indexing_mode = "consistent"

    included_path {
      path = "/*"
    }

    included_path {
      path = "%s"
    }

    excluded_path {
      path = "%s"
    }

    composite_index {
      index {
        path  = "/path1"
        order = "ascending"
      }
      index {
        path  = "/path2"
        order = "descending"
      }
    }

    composite_index {
      index {
        path  = "/path3"
        order = "ascending"
      }
      index {
        path  = "/path4"
        order = "descending"
      }
    }

    spatial_index {
      path = "/path/*"
    }

    spatial_index {
      path = "/test/to/all/?"
    }
  }
}
`, CosmosSqlDatabaseResource{}.basic(data), data.RandomInteger, includedPath, excludedPath)
}

func (CosmosSqlContainerResource) indexing_policy_update_includedPath(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_sql_container" "test" {
  name                = "acctest-CSQLC-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_sql_database.test.name
  partition_key_paths = ["/definition/id"]

  indexing_policy {
    indexing_mode = "none"
  }
}
`, CosmosSqlDatabaseResource{}.basic(data), data.RandomInteger)
}

func (CosmosSqlContainerResource) partition_key_version(data acceptance.TestData, version int) string {
	return fmt.Sprintf(`
%[1]s
resource "azurerm_cosmosdb_sql_container" "test" {
  name                  = "acctest-CSQLC-%[2]d"
  resource_group_name   = azurerm_cosmosdb_account.test.resource_group_name
  account_name          = azurerm_cosmosdb_account.test.name
  database_name         = azurerm_cosmosdb_sql_database.test.name
  partition_key_paths   = ["/definition/id"]
  partition_key_version = %[3]d
}
`, CosmosSqlDatabaseResource{}.basic(data), data.RandomInteger, version)
}

func (CosmosSqlContainerResource) conflictResolutionPolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_sql_container" "test" {
  name                = "acctest-CSQLC-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_sql_database.test.name
  partition_key_paths = ["/definition/id"]

  conflict_resolution_policy {
    mode                          = "Custom"
    conflict_resolution_procedure = "dbs/{0}/colls/{1}/sprocs/{2}"
  }
}
`, CosmosSqlDatabaseResource{}.basic(data), data.RandomInteger)
}

func (CosmosSqlContainerResource) hierarchicalPartitionKeys(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_sql_container" "test" {
  name                  = "acctest-CSQLC-%[2]d"
  resource_group_name   = azurerm_cosmosdb_account.test.resource_group_name
  account_name          = azurerm_cosmosdb_account.test.name
  database_name         = azurerm_cosmosdb_sql_database.test.name
  partition_key_kind    = "MultiHash"
  partition_key_paths   = ["/definition", "/id", "/sessionId"]
  partition_key_version = 2
}
`, CosmosSqlDatabaseResource{}.basic(data), data.RandomInteger)
}

package cosmos_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type CosmosGremlinGraphResource struct {
}

func TestAccCosmosDbGremlinGraph_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_gremlin_graph", "test")
	r := CosmosGremlinGraphResource{}

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

func TestAccCosmosDbGremlinGraph_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_gremlin_graph", "test")
	r := CosmosGremlinGraphResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_cosmosdb_gremlin_graph"),
		},
	})
}

func TestAccCosmosDbGremlinGraph_customConflictResolutionPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_gremlin_graph", "test")
	r := CosmosGremlinGraphResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.customConflictResolutionPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDbGremlinGraph_indexPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_gremlin_graph", "test")
	r := CosmosGremlinGraphResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.indexPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.indexPolicyCompositeIndex(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.indexPolicySpatialIndex(data),
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

func TestAccCosmosDbGremlinGraph_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_gremlin_graph", "test")
	r := CosmosGremlinGraphResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.update(data, 700, 900),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("throughput").HasValue("700"),
				check.That(data.ResourceName).Key("default_ttl").HasValue("900"),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data, 1700, 1900),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("throughput").HasValue("1700"),
				check.That(data.ResourceName).Key("default_ttl").HasValue("1900"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDbGremlinGraph_autoscale(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_gremlin_graph", "test")
	r := CosmosGremlinGraphResource{}

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

func TestAccCosmosDbGremlinGraph_partition_key_version(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_gremlin_graph", "test")
	r := CosmosGremlinGraphResource{}

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

func (t CosmosGremlinGraphResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.GremlinGraphID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Cosmos.GremlinClient.GetGremlinGraph(ctx, id.ResourceGroup, id.DatabaseAccountName, id.GremlinDatabaseName, id.GraphName)
	if err != nil {
		return nil, fmt.Errorf("reading Cosmos Gremlin Graph (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (CosmosGremlinGraphResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_gremlin_graph" "test" {
  name                = "acctest-CGRPC-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_gremlin_database.test.name
  partition_key_path  = "/test"
  throughput          = 400
}
`, CosmosGremlinDatabaseResource{}.basic(data), data.RandomInteger)
}

func (r CosmosGremlinGraphResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_gremlin_graph" "import" {
  name                = azurerm_cosmosdb_gremlin_graph.test.name
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_gremlin_database.test.name
  partition_key_path  = azurerm_cosmosdb_gremlin_graph.test.partition_key_path
}
`, r.basic(data))
}

func (CosmosGremlinGraphResource) customConflictResolutionPolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_gremlin_graph" "test" {
  name                = "acctest-CGRPC-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_gremlin_database.test.name
  partition_key_path  = "/test"
  throughput          = 400

  index_policy {
    automatic      = true
    indexing_mode  = "Consistent"
    included_paths = ["/*"]
    excluded_paths = ["/\"_etag\"/?"]
  }

  conflict_resolution_policy {
    mode                          = "Custom"
    conflict_resolution_procedure = "dbs/{0}/colls/{1}/sprocs/{2}"
  }
}
`, CosmosGremlinDatabaseResource{}.basic(data), data.RandomInteger)
}

func (CosmosGremlinGraphResource) indexPolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_gremlin_graph" "test" {
  name                = "acctest-CGRPC-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_gremlin_database.test.name
  partition_key_path  = "/test"
  throughput          = 400

  index_policy {
    automatic     = false
    indexing_mode = "None"
  }

  conflict_resolution_policy {
    mode                     = "LastWriterWins"
    conflict_resolution_path = "/_ts"
  }
}
`, CosmosGremlinDatabaseResource{}.basic(data), data.RandomInteger)
}

func (CosmosGremlinGraphResource) indexPolicyCompositeIndex(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_gremlin_graph" "test" {
  name                = "acctest-CGRPC-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_gremlin_database.test.name
  partition_key_path  = "/test"
  throughput          = 400

  index_policy {
    automatic     = true
    indexing_mode = "Consistent"

    composite_index {
      index {
        path  = "/path1"
        order = "Ascending"
      }
      index {
        path  = "/path2"
        order = "Descending"
      }
    }

    composite_index {
      index {
        path  = "/path3"
        order = "Ascending"
      }
      index {
        path  = "/path4"
        order = "Descending"
      }
    }

    spatial_index {
      path = "/path/*"
    }

    spatial_index {
      path = "/test/to/all/?"
    }
  }

  conflict_resolution_policy {
    mode                     = "LastWriterWins"
    conflict_resolution_path = "/_ts"
  }
}
`, CosmosGremlinDatabaseResource{}.basic(data), data.RandomInteger)
}

func (CosmosGremlinGraphResource) indexPolicySpatialIndex(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_gremlin_graph" "test" {
  name                = "acctest-CGRPC-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_gremlin_database.test.name
  partition_key_path  = "/test"
  throughput          = 400

  index_policy {
    automatic     = true
    indexing_mode = "Consistent"

    composite_index {
      index {
        path  = "/path1"
        order = "Ascending"
      }
      index {
        path  = "/path2"
        order = "Descending"
      }
    }

    composite_index {
      index {
        path  = "/path3"
        order = "Ascending"
      }
      index {
        path  = "/path4"
        order = "Descending"
      }
    }
  }

  conflict_resolution_policy {
    mode                     = "LastWriterWins"
    conflict_resolution_path = "/_ts"
  }
}
`, CosmosGremlinDatabaseResource{}.basic(data), data.RandomInteger)
}

func (CosmosGremlinGraphResource) update(data acceptance.TestData, throughput int, defaultTTL int) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_gremlin_graph" "test" {
  name                = "acctest-CGRPC-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_gremlin_database.test.name
  partition_key_path  = "/test"
  throughput          = %[3]d
  default_ttl         = %[4]d

  index_policy {
    automatic      = true
    indexing_mode  = "Consistent"
    included_paths = ["/*"]
    excluded_paths = ["/\"_etag\"/?"]
  }

  unique_key {
    paths = ["/definition/id1", "/definition/id2"]
  }
}
`, CosmosGremlinDatabaseResource{}.basic(data), data.RandomInteger, throughput, defaultTTL)
}

func (CosmosGremlinGraphResource) autoscale(data acceptance.TestData, maxThroughput int) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_gremlin_graph" "test" {
  name                = "acctest-CGRPC-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_gremlin_database.test.name
  partition_key_path  = "/test"

  autoscale_settings {
    max_throughput = %[3]d
  }

  index_policy {
    automatic      = true
    indexing_mode  = "Consistent"
    included_paths = ["/*"]
    excluded_paths = ["/\"_etag\"/?"]
  }
}
`, CosmosGremlinDatabaseResource{}.basic(data), data.RandomInteger, maxThroughput)
}

func (CosmosGremlinGraphResource) partition_key_version(data acceptance.TestData, version int) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_gremlin_graph" "test" {
  name                  = "acctest-CGRPC-%[2]d"
  resource_group_name   = azurerm_cosmosdb_account.test.resource_group_name
  account_name          = azurerm_cosmosdb_account.test.name
  database_name         = azurerm_cosmosdb_gremlin_database.test.name
  partition_key_path    = "/test"
  partition_key_version = %[3]d
}
`, CosmosGremlinDatabaseResource{}.basic(data), data.RandomInteger, version)
}

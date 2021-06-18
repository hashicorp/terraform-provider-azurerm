package cosmos_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-01-15/documentdb"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type CosmosMongoCollectionResource struct {
}

func TestAccCosmosDbMongoCollection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_mongo_collection", "test")
	r := CosmosMongoCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("throughput").HasValue("400"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDbMongoCollection_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_mongo_collection", "test")
	r := CosmosMongoCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("shard_key").HasValue("seven"),
				check.That(data.ResourceName).Key("default_ttl_seconds").HasValue("707"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDbMongoCollection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_mongo_collection", "test")
	r := CosmosMongoCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.complete(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("shard_key").HasValue("seven"),
				check.That(data.ResourceName).Key("default_ttl_seconds").HasValue("707"),
			),
		},
		data.ImportStep(),
		{
			Config: r.updated(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_ttl_seconds").HasValue("70707"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDbMongoCollection_throughput(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_mongo_collection", "test")
	r := CosmosMongoCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.throughput(data, 700),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.throughput(data, 1400),
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

func TestAccCosmosDbMongoCollection_withIndex(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_mongo_collection", "test")
	r := CosmosMongoCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withIndex(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("default_ttl_seconds").HasValue("707"),
				check.That(data.ResourceName).Key("index.#").HasValue("4"),
				check.That(data.ResourceName).Key("system_indexes.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDbMongoCollection_analyticalStorageTTL(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_mongo_collection", "test")
	r := CosmosMongoCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.analyticalStorageTTL(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("analytical_storage_ttl").HasValue("600"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDbMongoCollection_autoscale(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_mongo_collection", "test")
	r := CosmosMongoCollectionResource{}

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

func TestAccCosmosDbMongoCollection_ver36(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_mongo_collection", "test")
	r := CosmosMongoCollectionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.ver36(data),
			Check: acceptance.ComposeAggregateTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCosmosDbMongoCollection_serverless(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_mongo_collection", "test")
	r := CosmosMongoCollectionResource{}

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

func (t CosmosMongoCollectionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.MongodbCollectionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Cosmos.MongoDbClient.GetMongoDBCollection(ctx, id.ResourceGroup, id.DatabaseAccountName, id.MongodbDatabaseName, id.CollectionName)
	if err != nil {
		return nil, fmt.Errorf("reading Cosmos Mongo Collection (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (CosmosMongoCollectionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_mongo_collection" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_cosmosdb_mongo_database.test.resource_group_name
  account_name        = azurerm_cosmosdb_mongo_database.test.account_name
  database_name       = azurerm_cosmosdb_mongo_database.test.name

  index {
    keys   = ["_id"]
    unique = true
  }
}
`, CosmosMongoDatabaseResource{}.basic(data), data.RandomInteger)
}

func (CosmosMongoCollectionResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_mongo_collection" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_cosmosdb_mongo_database.test.resource_group_name
  account_name        = azurerm_cosmosdb_mongo_database.test.account_name
  database_name       = azurerm_cosmosdb_mongo_database.test.name

  index {
    keys   = ["_id"]
    unique = true
  }

  shard_key           = "seven"
  default_ttl_seconds = 707
}
`, CosmosMongoDatabaseResource{}.basic(data), data.RandomInteger)
}

func (CosmosMongoCollectionResource) updated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_mongo_collection" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_cosmosdb_mongo_database.test.resource_group_name
  account_name        = azurerm_cosmosdb_mongo_database.test.account_name
  database_name       = azurerm_cosmosdb_mongo_database.test.name

  index {
    keys   = ["_id"]
    unique = true
  }

  shard_key           = "seven"
  default_ttl_seconds = 70707
}
`, CosmosMongoDatabaseResource{}.basic(data), data.RandomInteger)
}

func (CosmosMongoCollectionResource) throughput(data acceptance.TestData, throughput int) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_mongo_collection" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_cosmosdb_mongo_database.test.resource_group_name
  account_name        = azurerm_cosmosdb_mongo_database.test.account_name
  database_name       = azurerm_cosmosdb_mongo_database.test.name

  index {
    keys   = ["_id"]
    unique = true
  }

  throughput = %[3]d
}
`, CosmosMongoDatabaseResource{}.basic(data), data.RandomInteger, throughput)
}

func (CosmosMongoCollectionResource) autoscale(data acceptance.TestData, maxThroughput int) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_mongo_collection" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_cosmosdb_mongo_database.test.resource_group_name
  account_name        = azurerm_cosmosdb_mongo_database.test.account_name
  database_name       = azurerm_cosmosdb_mongo_database.test.name
  shard_key           = "seven"

  index {
    keys   = ["_id"]
    unique = true
  }

  autoscale_settings {
    max_throughput = %[3]d
  }
}
`, CosmosMongoDatabaseResource{}.basic(data), data.RandomInteger, maxThroughput)
}

func (CosmosMongoCollectionResource) withIndex(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_mongo_collection" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_cosmosdb_mongo_database.test.resource_group_name
  account_name        = azurerm_cosmosdb_mongo_database.test.account_name
  database_name       = azurerm_cosmosdb_mongo_database.test.name
  default_ttl_seconds = 707
  throughput          = 400

  index {
    keys   = ["seven", "six"]
    unique = true
  }

  index {
    keys   = ["day"]
    unique = false
  }

  index {
    keys = ["month"]
  }

  index {
    keys   = ["_id"]
    unique = true
  }
}
`, CosmosMongoDatabaseResource{}.basic(data), data.RandomInteger)
}

func (CosmosMongoCollectionResource) ver36(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_mongo_database" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
}

resource "azurerm_cosmosdb_mongo_collection" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_cosmosdb_mongo_database.test.resource_group_name
  account_name        = azurerm_cosmosdb_mongo_database.test.account_name
  database_name       = azurerm_cosmosdb_mongo_database.test.name

  index {
    keys   = ["_id"]
    unique = true
  }
}
`, CosmosDBAccountResource{}.capabilities(data, documentdb.MongoDB, []string{"EnableMongo"}), data.RandomInteger)
}

func (CosmosMongoCollectionResource) serverless(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_mongo_database" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
}

resource "azurerm_cosmosdb_mongo_collection" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_cosmosdb_mongo_database.test.resource_group_name
  account_name        = azurerm_cosmosdb_mongo_database.test.account_name
  database_name       = azurerm_cosmosdb_mongo_database.test.name

  index {
    keys   = ["_id"]
    unique = true
  }
}
`, CosmosDBAccountResource{}.capabilities(data, documentdb.MongoDB, []string{"EnableMongo", "EnableServerless"}), data.RandomInteger)
}

func (CosmosMongoCollectionResource) analyticalStorageTTL(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_mongo_database" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
}

resource "azurerm_cosmosdb_mongo_collection" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_cosmosdb_mongo_database.test.resource_group_name
  account_name        = azurerm_cosmosdb_mongo_database.test.account_name
  database_name       = azurerm_cosmosdb_mongo_database.test.name

  index {
    keys   = ["_id"]
    unique = true
  }

  analytical_storage_ttl = 600
}
`, CosmosDBAccountResource{}.mongoAnalyticalStorage(data, documentdb.Eventual), data.RandomInteger, data.RandomInteger)
}

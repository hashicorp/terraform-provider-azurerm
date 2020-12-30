package cosmos_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/preview/cosmos-db/mgmt/2020-04-01-preview/documentdb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMCosmosDbMongoCollection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_mongo_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDbMongoCollectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDbMongoCollection_basic(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbMongoCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "throughput", "400"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDbMongoCollection_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_mongo_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDbMongoCollectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDbMongoCollection_complete(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbMongoCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "shard_key", "seven"),
					resource.TestCheckResourceAttr(data.ResourceName, "default_ttl_seconds", "707"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDbMongoCollection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_mongo_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDbMongoCollectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDbMongoCollection_basic(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbMongoCollectionExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMCosmosDbMongoCollection_complete(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbMongoCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "shard_key", "seven"),
					resource.TestCheckResourceAttr(data.ResourceName, "default_ttl_seconds", "707"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMCosmosDbMongoCollection_updated(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbMongoCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "default_ttl_seconds", "70707"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDbMongoCollection_throughput(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_mongo_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDbMongoCollectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDbMongoCollection_throughput(data, 700),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbMongoCollectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMCosmosDbMongoCollection_throughput(data, 1400),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbMongoCollectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMCosmosDbMongoCollection_basic(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbMongoCollectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDbMongoCollection_withIndex(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_mongo_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDbMongoCollectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDbMongoCollection_withIndex(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbMongoCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "default_ttl_seconds", "707"),
					resource.TestCheckResourceAttr(data.ResourceName, "index.#", "3"),
					resource.TestCheckResourceAttr(data.ResourceName, "system_indexes.#", "2"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDbMongoCollection_autoscale(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_mongo_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDbMongoCollectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDbMongoCollection_autoscale(data, 4000),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbMongoCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "autoscale_settings.0.max_throughput", "4000"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMCosmosDbMongoCollection_autoscale(data, 5000),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbMongoCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "autoscale_settings.0.max_throughput", "5000"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMCosmosDbMongoCollection_autoscale(data, 4000),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbMongoCollectionExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "autoscale_settings.0.max_throughput", "4000"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDbMongoCollection_ver36(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_mongo_collection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDbMongoCollectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDbMongoCollection_ver36(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbMongoCollectionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMCosmosDbMongoCollectionDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Cosmos.MongoDbClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_cosmosdb_mongo_collection" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		account := rs.Primary.Attributes["account_name"]
		database := rs.Primary.Attributes["database_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.GetMongoDBCollection(ctx, resourceGroup, account, database, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Error checking destroy for Cosmos Mongo Collection %s (account %s, database %s) still exists:\n%v", name, account, database, err)
			}
		}

		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Cosmos Mongo Collection %s (account %s) still exists:\n%#v", name, account, resp)
		}
	}

	return nil
}

func testCheckAzureRMCosmosDbMongoCollectionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Cosmos.MongoDbClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		account := rs.Primary.Attributes["account_name"]
		database := rs.Primary.Attributes["database_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.GetMongoDBCollection(ctx, resourceGroup, account, database, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on cosmosAccountsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Cosmos database '%s' (account: '%s', database: %s) does not exist", name, account, database)
		}

		return nil
	}
}

func testAccAzureRMCosmosDbMongoCollection_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_mongo_collection" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_cosmosdb_mongo_database.test.resource_group_name
  account_name        = azurerm_cosmosdb_mongo_database.test.account_name
  database_name       = azurerm_cosmosdb_mongo_database.test.name
}
`, testAccAzureRMCosmosDbMongoDatabase_basic(data), data.RandomInteger)
}

func testAccAzureRMCosmosDbMongoCollection_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_mongo_collection" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_cosmosdb_mongo_database.test.resource_group_name
  account_name        = azurerm_cosmosdb_mongo_database.test.account_name
  database_name       = azurerm_cosmosdb_mongo_database.test.name

  shard_key           = "seven"
  default_ttl_seconds = 707
}
`, testAccAzureRMCosmosDbMongoDatabase_basic(data), data.RandomInteger)
}

func testAccAzureRMCosmosDbMongoCollection_updated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_mongo_collection" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_cosmosdb_mongo_database.test.resource_group_name
  account_name        = azurerm_cosmosdb_mongo_database.test.account_name
  database_name       = azurerm_cosmosdb_mongo_database.test.name

  shard_key           = "seven"
  default_ttl_seconds = 70707
}
`, testAccAzureRMCosmosDbMongoDatabase_basic(data), data.RandomInteger)
}

func testAccAzureRMCosmosDbMongoCollection_throughput(data acceptance.TestData, throughput int) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_mongo_collection" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_cosmosdb_mongo_database.test.resource_group_name
  account_name        = azurerm_cosmosdb_mongo_database.test.account_name
  database_name       = azurerm_cosmosdb_mongo_database.test.name

  throughput = %[3]d
}
`, testAccAzureRMCosmosDbMongoDatabase_basic(data), data.RandomInteger, throughput)
}

func testAccAzureRMCosmosDbMongoCollection_autoscale(data acceptance.TestData, maxThroughput int) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_mongo_collection" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_cosmosdb_mongo_database.test.resource_group_name
  account_name        = azurerm_cosmosdb_mongo_database.test.account_name
  database_name       = azurerm_cosmosdb_mongo_database.test.name
  shard_key           = "seven"

  autoscale_settings {
    max_throughput = %[3]d
  }
}
`, testAccAzureRMCosmosDbMongoDatabase_basic(data), data.RandomInteger, maxThroughput)
}

func testAccAzureRMCosmosDbMongoCollection_withIndex(data acceptance.TestData) string {
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
}
`, testAccAzureRMCosmosDbMongoDatabase_basic(data), data.RandomInteger)
}

func testAccAzureRMCosmosDbMongoCollection_ver36(data acceptance.TestData) string {
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
`, testAccAzureRMCosmosDBAccount_capabilities(data, documentdb.MongoDB, []string{"EnableMongo"}), data.RandomInteger)
}

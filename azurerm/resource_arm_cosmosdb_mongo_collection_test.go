package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMCosmosDbMongoCollection_basic(t *testing.T) {
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_cosmosdb_mongo_collection.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDbMongoCollectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDbMongoCollection_basic(ri, acceptance.Location()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbMongoCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "throughput", "400"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMCosmosDbMongoCollection_complete(t *testing.T) {
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_cosmosdb_mongo_collection.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDbMongoCollectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDbMongoCollection_complete(ri, acceptance.Location()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbMongoCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "shard_key", "seven"),
					resource.TestCheckResourceAttr(resourceName, "default_ttl_seconds", "707"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMCosmosDbMongoCollection_update(t *testing.T) {
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_cosmosdb_mongo_collection.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDbMongoCollectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDbMongoCollection_basic(ri, acceptance.Location()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbMongoCollectionExists(resourceName),
				),
			},
			{
				Config: testAccAzureRMCosmosDbMongoCollection_complete(ri, acceptance.Location()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbMongoCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "shard_key", "seven"),
					resource.TestCheckResourceAttr(resourceName, "default_ttl_seconds", "707"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMCosmosDbMongoCollection_updated(ri, acceptance.Location()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbMongoCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "default_ttl_seconds", "70707"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMCosmosDbMongoCollection_throughput(t *testing.T) {
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_cosmosdb_mongo_collection.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDbMongoCollectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDbMongoCollection_throughput(ri, acceptance.Location(), 700),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbMongoCollectionExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMCosmosDbMongoCollection_throughput(ri, acceptance.Location(), 1400),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbMongoCollectionExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMCosmosDbMongoCollection_basic(ri, acceptance.Location()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbMongoCollectionExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMCosmosDbMongoCollectionDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Cosmos.DatabaseClient
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
		client := acceptance.AzureProvider.Meta().(*clients.Client).Cosmos.DatabaseClient
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

func testAccAzureRMCosmosDbMongoCollection_basic(rInt int, location string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_mongo_collection" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = "${azurerm_cosmosdb_mongo_database.test.resource_group_name}"
  account_name        = "${azurerm_cosmosdb_mongo_database.test.account_name}"
  database_name       = "${azurerm_cosmosdb_mongo_database.test.name}"
}
`, testAccAzureRMCosmosDbMongoDatabase_basic(rInt, location), rInt)
}

func testAccAzureRMCosmosDbMongoCollection_complete(rInt int, location string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_mongo_collection" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = "${azurerm_cosmosdb_mongo_database.test.resource_group_name}"
  account_name        = "${azurerm_cosmosdb_mongo_database.test.account_name}"
  database_name       = "${azurerm_cosmosdb_mongo_database.test.name}"

  shard_key           = "seven"
  default_ttl_seconds = 707
}
`, testAccAzureRMCosmosDbMongoDatabase_basic(rInt, location), rInt)
}

func testAccAzureRMCosmosDbMongoCollection_updated(rInt int, location string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_mongo_collection" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = "${azurerm_cosmosdb_mongo_database.test.resource_group_name}"
  account_name        = "${azurerm_cosmosdb_mongo_database.test.account_name}"
  database_name       = "${azurerm_cosmosdb_mongo_database.test.name}"

  shard_key           = "seven"
  default_ttl_seconds = 70707
}
`, testAccAzureRMCosmosDbMongoDatabase_basic(rInt, location), rInt)
}

func testAccAzureRMCosmosDbMongoCollection_throughput(rInt int, location string, throughput int) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_mongo_collection" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = "${azurerm_cosmosdb_mongo_database.test.resource_group_name}"
  account_name        = "${azurerm_cosmosdb_mongo_database.test.account_name}"
  database_name       = "${azurerm_cosmosdb_mongo_database.test.name}"

  throughput = %[3]d
}
`, testAccAzureRMCosmosDbMongoDatabase_basic(rInt, location), rInt, throughput)
}

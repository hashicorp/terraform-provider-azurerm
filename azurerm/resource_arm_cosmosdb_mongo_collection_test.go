package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMCosmosDbMongoCollection_basic(t *testing.T) {
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_cosmosdb_mongo_collection.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDbMongoCollectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDbMongoCollection_basic(ri, testLocation()),
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

func TestAccAzureRMCosmosDbMongoCollection_complete(t *testing.T) {
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_cosmosdb_mongo_collection.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDbMongoCollectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDbMongoCollection_complete(ri, testLocation()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbMongoCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "shard_key", "day"),
					resource.TestCheckResourceAttr(resourceName, "default_ttl_seconds", "707"),
					resource.TestCheckResourceAttr(resourceName, "indexes.#", "2"),
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
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDbMongoCollectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDbMongoCollection_basic(ri, testLocation()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbMongoCollectionExists(resourceName),
				),
			},
			{
				Config: testAccAzureRMCosmosDbMongoCollection_complete(ri, testLocation()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbMongoCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "shard_key", "day"),
					resource.TestCheckResourceAttr(resourceName, "default_ttl_seconds", "707"),
					resource.TestCheckResourceAttr(resourceName, "indexes.#", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMCosmosDbMongoCollection_updated(ri, testLocation()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbMongoCollectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "default_ttl_seconds", "70707"),
					resource.TestCheckResourceAttr(resourceName, "indexes.#", "3"),
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
	client := testAccProvider.Meta().(*ArmClient).cosmos.DatabaseClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

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
		client := testAccProvider.Meta().(*ArmClient).cosmos.DatabaseClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

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

  default_ttl_seconds = 707
  shard_key           = "day"

  indexes {
    key    = "seven"
    unique = false
  }

  indexes {
    key    = "day"
    unique = true
  }
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

  default_ttl_seconds = 70707

  indexes {
    key    = "seven"
    unique = true
  }

  indexes {
    key    = "day"
    unique = false
  }

  indexes {
    key    = "fool"
    unique = false
  }
}
`, testAccAzureRMCosmosDbMongoDatabase_basic(rInt, location), rInt)
}

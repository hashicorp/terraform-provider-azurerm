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

func TestAccAzureRMCosmosMongoDatabase_basic(t *testing.T) {
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_cosmos_mongo_database.test"
	rn := fmt.Sprintf("acctest-%[1]d", ri)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosMongoDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosMongoDatabase_basic(ri, testLocation()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosMongoDatabaseExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rn),
					resource.TestCheckResourceAttr(resourceName, "account_name", rn),
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

func TestAccAzureRMCosmosMongoDatabase_debug(t *testing.T) {
	resourceName := "azurerm_cosmos_mongo_database.test"
	ri := tf.AccRandTimeInt()
	rn := fmt.Sprintf("acctest-%[1]d", ri)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosMongoDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosMongoDatabase_debug(ri),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosMongoDatabaseExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rn),
					resource.TestCheckResourceAttr(resourceName, "account_name", rn),
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

func testCheckAzureRMCosmosMongoDatabaseDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).cosmosAccountsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for rn, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_cosmos_mongo_database" {
			continue
		}

		if err := tf.AccCheckResourceAttributes(rs.Primary.Attributes, "name", "resource_group_name", "account_name"); err != nil {
			return fmt.Errorf("resource %s is missing an attribute: %v", rn, err)
		}
		name := rs.Primary.Attributes["name"]
		account := rs.Primary.Attributes["account_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.GetMongoDatabase(ctx, resourceGroup, account, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Error checking destroy for Cosmos Mongo Database %s (account %s) still exists:\n%v", name, account, err)
			}
		}

		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Cosmos Mongo Database %s (account %s) still exists:\n%#v", name, account, resp)
		}
	}

	return nil
}

func testCheckAzureRMCosmosMongoDatabaseExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ArmClient).cosmosAccountsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if err := tf.AccCheckResourceAttributes(rs.Primary.Attributes, "name", "resource_group_name", "account_name"); err != nil {
			return fmt.Errorf("resource %s is missing an attribute: %v", resourceName, err)
		}
		name := rs.Primary.Attributes["name"]
		account := rs.Primary.Attributes["account_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.GetMongoDatabase(ctx, resourceGroup, account, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on cosmosAccountsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Cosmos database '%s' (account: '%s') does not exist", name, account)
		}

		return nil
	}
}

func testAccAzureRMCosmosMongoDatabase_basic(rInt int, location string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmos_mongo_database" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = "${azurerm_cosmosdb_account.test.resource_group_name}"
  account_name        = "${azurerm_cosmosdb_account.test.name}"
}
`, testAccAzureRMCosmosDBAccount_mongoDB(rInt, location), rInt)
}

func testAccAzureRMCosmosMongoDatabase_debug(rInt int) string {
	return fmt.Sprintf(`

resource "azurerm_cosmos_mongo_database" "test" {
  name         = "acctest-%[1]d"
  account_name = "kt-cosmos-mongo"
  resource_group_name = "kt-cosmos-201904"
}
`, rInt)
}

func testAccAzureRMCosmosMongoDatabase_debug2(rInt int) string {
	return fmt.Sprintf(`

resource "azurerm_cosmos_mongo_database" "test" {
  name         = "acctest-%[1]d2"
  account_name = "kt-cosmos-mongo"
  resource_group_name = "kt-cosmos-201904"
}
`, rInt)
}

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

func TestAccAzureRMCosmosCassandraTable_basic(t *testing.T) {
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_cosmos_cassandra_table.test"
	rn := fmt.Sprintf("acctest-%[1]d", ri)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosCassandraTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosCassandraTable_basic(ri, testLocation()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosCassandraTableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rn),
					resource.TestCheckResourceAttr(resourceName, "account_name", rn),
					resource.TestCheckResourceAttr(resourceName, "keyspace_name", rn),
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

func TestAccAzureRMCosmosCassandraTable_complete(t *testing.T) {
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_cosmos_cassandra_table.test"
	rn := fmt.Sprintf("acctest-%[1]d", ri)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosCassandraTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosCassandraTable_complete(ri, testLocation()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosCassandraTableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rn),
					resource.TestCheckResourceAttr(resourceName, "account_name", rn),
					resource.TestCheckResourceAttr(resourceName, "keyspace_name", rn),
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

func TestAccAzureRMCosmosCassandraTable_update(t *testing.T) {
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_cosmos_cassandra_table.test"
	rn := fmt.Sprintf("acctest-%[1]d", ri)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosCassandraTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosCassandraTable_basic(ri, testLocation()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosCassandraTableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rn),
					resource.TestCheckResourceAttr(resourceName, "account_name", rn),
					resource.TestCheckResourceAttr(resourceName, "keyspace_name", rn),
				),
			},
			{
				Config: testAccAzureRMCosmosCassandraTable_complete(ri, testLocation()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosCassandraTableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rn),
					resource.TestCheckResourceAttr(resourceName, "account_name", rn),
					resource.TestCheckResourceAttr(resourceName, "keyspace_name", rn),
					//todo check set values when the SDK actually reads them
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMCosmosCassandraTable_updated(ri, testLocation()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosCassandraTableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rn),
					resource.TestCheckResourceAttr(resourceName, "account_name", rn),
					resource.TestCheckResourceAttr(resourceName, "keyspace_name", rn),
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

func TestAccAzureRMCosmosCassandraTable_debug(t *testing.T) {
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_cosmos_cassandra_table.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosCassandraTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosCassandraTable_debug(ri, testLocation()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosCassandraTableExists(resourceName),
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

func testCheckAzureRMCosmosCassandraTableDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).cosmosAccountsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for rn, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_cosmos_cassandra_table" {
			continue
		}

		if err := tf.AccCheckResourceAttributes(rs.Primary.Attributes, "name", "resource_group_name", "account_name", "keyspace_name"); err != nil {
			return fmt.Errorf("resource %s is missing an attribute: %v", rn, err)
		}
		name := rs.Primary.Attributes["name"]
		account := rs.Primary.Attributes["account_name"]
		database := rs.Primary.Attributes["keyspace_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.GetCassandraTable(ctx, resourceGroup, account, database, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Error checking destroy for Cosmos Cassandra Table %s (account %s, database %s) still exists:\n%v", name, account, database, err)
			}
		}

		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Cosmos Cassandra Table %s (account %s) still exists:\n%#v", name, account, resp)
		}
	}

	return nil
}

func testCheckAzureRMCosmosCassandraTableExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ArmClient).cosmosAccountsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if err := tf.AccCheckResourceAttributes(rs.Primary.Attributes, "name", "resource_group_name", "account_name", "keyspace_name"); err != nil {
			return fmt.Errorf("resource %s is missing an attribute: %v", resourceName, err)
		}
		name := rs.Primary.Attributes["name"]
		account := rs.Primary.Attributes["account_name"]
		database := rs.Primary.Attributes["keyspace_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.GetCassandraTable(ctx, resourceGroup, account, database, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on cosmosAccountsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Cosmos database '%s' (account: '%s', database: %s) does not exist", name, account, database)
		}

		return nil
	}
}

func testAccAzureRMCosmosCassandraTable_basic(rInt int, location string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmos_cassandra_table" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = "${azurerm_cosmos_sql_database.test.resource_group_name}"
  account_name        = "${azurerm_cosmos_sql_database.test.account_name}"
  keyspace_name       = "${azurerm_cosmos_sql_database.test.name}"
}
`, testAccAzureRMCosmosSQLDatabase_basic(rInt, location), rInt)
}

func testAccAzureRMCosmosCassandraTable_complete(rInt int, location string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmos_cassandra_table" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = "${azurerm_cosmos_sql_database.test.resource_group_name}"
  account_name        = "${azurerm_cosmos_sql_database.test.account_name}"
  keyspace_name       = "${azurerm_cosmos_sql_database.test.name}"


}
`, testAccAzureRMCosmosSQLDatabase_basic(rInt, location), rInt)
}

func testAccAzureRMCosmosCassandraTable_updated(rInt int, location string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmos_cassandra_table" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = "${azurerm_cosmos_sql_database.test.resource_group_name}"
  account_name        = "${azurerm_cosmos_sql_database.test.account_name}"
  keyspace_name       = "${azurerm_cosmos_sql_database.test.name}"

}
`, testAccAzureRMCosmosSQLDatabase_basic(rInt, location), rInt)
}

func testAccAzureRMCosmosCassandraTable_debug(rInt int, location string) string {
	return fmt.Sprintf(`


resource "azurerm_cosmos_cassandra_table" "test" {
  name                = "littbobbytables2"
  resource_group_name = "kt-cosmos-201904"
  account_name        = "kt-cosmos-cassandra"
  keyspace_name       = "sevendaykeyspace"

}
`)
}

func testAccAzureRMCosmosCassandraTable_debug2(rInt int, location string) string {
	return fmt.Sprintf(`


resource "azurerm_cosmos_cassandra_table" "test" {
  name                = "littbobbytables2"
  resource_group_name = "kt-cosmos-201904"
  account_name        = "kt-cosmos-cassandra"
  keyspace_name       = "sevendaykeyspace"

}
`)
}

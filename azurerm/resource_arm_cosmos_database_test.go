package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

//consistency
func TestAccAzureRMCosmosDatabase_basic(t *testing.T) {
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_cosmosdb_database.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		//CheckDestroy: testCheckAzureRMCosmosDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosdatabase_basic(ri, testLocation()),
				Check: resource.ComposeAggregateTestCheckFunc(
					//testCheckAzureRMCosmosDatabaseExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "account_name"),
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

func TestAccAzureRMCosmosDatabase_debug(t *testing.T) {
	resourceName := "azurerm_cosmosdb_database.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		//CheckDestroy: testCheckAzureRMCosmosDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosdatabase_debug(),
				Check: resource.ComposeAggregateTestCheckFunc(
					//	testCheckAzureRMCosmosDatabaseExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "account_name"),
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

/*
func testCheckAzureRMCosmosDatabaseDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).cosmosDatabasesClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_cosmosdb_database" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		//account_name := rs.Primary.Attributes["account_name"]

		resp, err := client.Get(ctx, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("CosmosDB Database still exists:\n%#v", resp)
		}
	}

	return nil
}

func testCheckAzureRMCosmosDatabaseExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ArmClient).cosmosDatabasesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		account_name := rs.Primary.Attributes["account_name"]

		resp, err := client.Get(ctx, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on cosmosAccountsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: CosmosDB database '%s' (account: '%s') does not exist", name, account_name)
		}

		return nil
	}
}
*/
func testAccAzureRMCosmosdatabase_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  offer_type          = "Standard"

  consistency_policy {
    consistency_level = "BoundedStaleness"

  }

  geo_location {
    location          = "${azurerm_resource_group.test.location}"
    failover_priority = 0
  }
}

resource "azurerm_cosmosdb_database" "test" {
  name         = "acctest-%[1]d"
  account_name = "${azurerm_cosmosdb_account.test.name}"
  account_key  = "${azurerm_cosmosdb_account.test.primary_key}"
}
`, rInt, location)
}

func testAccAzureRMCosmosdatabase_debug() string {
	return fmt.Sprintf(`

resource "azurerm_cosmosdb_database" "test" {
  name         = "SevenDayDBs"
  account_name = "kt-cosmos-201903"
  account_key  = "yvLnqDanONZn10a2ZgUge8cA3P9hkr3elsONP4yW6qADfj1RFkteABYBYEz627UAIUPGDRIeZNjaKqE4mBieqA=="
}
`)
}

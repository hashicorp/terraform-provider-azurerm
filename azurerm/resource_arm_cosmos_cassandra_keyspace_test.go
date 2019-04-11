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

func TestAccAzureRMCosmosCassandraKeyspace_basic(t *testing.T) {
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_cosmos_cassandra_keyspace.test"
	rn := fmt.Sprintf("acctest-%[1]d", ri)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosCassandraKeyspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosCassandraKeyspace_basic(ri, testLocation()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosCassandraKeyspaceExists(resourceName),
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

func testCheckAzureRMCosmosCassandraKeyspaceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).cosmosAccountsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for rn, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_cosmos_cassandra_keyspace" {
			continue
		}

		if err := tf.AccCheckResourceAttributes(rs.Primary.Attributes, "name", "resource_group_name", "account_name"); err != nil {
			return fmt.Errorf("resource %s is missing an attribute: %v", rn, err)
		}
		name := rs.Primary.Attributes["name"]
		account := rs.Primary.Attributes["account_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.GetCassandraKeyspace(ctx, resourceGroup, account, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Error checking destroy for Cosmos Cassandra Keyspace %s (account %s) still exists:\n%v", name, account, err)
			}
		}

		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Cosmos Cassandra Keyspace %s (account %s) still exists:\n%#v", name, account, resp)
		}
	}

	return nil
}

func testCheckAzureRMCosmosCassandraKeyspaceExists(resourceName string) resource.TestCheckFunc {
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

		resp, err := client.GetCassandraKeyspace(ctx, resourceGroup, account, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on cosmosAccountsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Cosmos database '%s' (account: '%s') does not exist", name, account)
		}

		return nil
	}
}

func testAccAzureRMCosmosCassandraKeyspace_basic(rInt int, location string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmos_cassandra_keyspace" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = "${azurerm_cosmosdb_account.test.resource_group_name}"
  account_name        = "${azurerm_cosmosdb_account.test.name}"
}
`, testAccAzureRMCosmosDBAccount_capabilityCassandra(rInt, location), rInt)
}

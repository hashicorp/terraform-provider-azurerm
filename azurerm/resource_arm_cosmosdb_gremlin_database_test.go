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

func TestAccAzureRMCosmosGremlinDatabase_basic(t *testing.T) {
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_cosmosdb_gremlin_database.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosGremlinDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosGremlinDatabase_basic(ri, acceptance.Location()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosGremlinDatabaseExists(resourceName),
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

func TestAccAzureRMCosmosGremlinDatabase_complete(t *testing.T) {
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_cosmosdb_gremlin_database.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosGremlinDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosGremlinDatabase_complete(ri, acceptance.Location()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosGremlinDatabaseExists(resourceName),
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

func testCheckAzureRMCosmosGremlinDatabaseDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Cosmos.DatabaseClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_cosmosdb_gremlin_database" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		account := rs.Primary.Attributes["account_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.GetGremlinDatabase(ctx, resourceGroup, account, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Error chacking destroy for Cosmos Gremlin Database %s (Account %s) still exists:\n%v", name, account, err)
			}
		}

		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Cosmos Gremlin Database %s (Account %s): still exist:\n%#v", name, account, resp)
		}
	}

	return nil
}

func testCheckAzureRMCosmosGremlinDatabaseExists(resourceName string) resource.TestCheckFunc {
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
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.GetGremlinDatabase(ctx, resourceGroup, account, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on cosmosAccountsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Cosmos database '%s' (Account: '%s') does not exist", name, account)
		}

		return nil
	}
}

func testAccAzureRMCosmosGremlinDatabase_basic(rInt int, location string) string {
	return fmt.Sprintf(`
	%[1]s
	
	resource "azurerm_cosmosdb_gremlin_database" "test" {
		name                = "acctest-%[2]d"
		resource_group_name = "${azurerm_cosmosdb_account.test.resource_group_name}"
		account_name        = "${azurerm_cosmosdb_account.test.name}"
	  }
	`, testAccAzureRMCosmosDBAccount_capabilityGremlin(rInt, location), rInt)
}

func testAccAzureRMCosmosGremlinDatabase_complete(rInt int, location string) string {
	return fmt.Sprintf(`
	%[1]s

	resource "azurerm_cosmosdb_gremlin_database" "test" {
		name                = "acctest-%[2]d"
		resource_group_name = "${azurerm_cosmosdb_account.test.resource_group_name}"
		account_name        = "${azurerm_cosmosdb_account.test.name}"
		throughput          = 700
	  }
	`, testAccAzureRMCosmosDBAccount_capabilityGremlin(rInt, location), rInt)
}

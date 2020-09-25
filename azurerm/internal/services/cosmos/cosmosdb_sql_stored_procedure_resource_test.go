package cosmos_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2020-04-01/documentdb"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMCosmosDbSqlStoredProcedure_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_sql_stored_procedure", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDbSqlStoredProcedureDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDbSqlStoredProcedure_basic(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbSqlStoredProcedureExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDbSqlStoredProcedure_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_sql_stored_procedure", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDbSqlStoredProcedureDestroy,
		Steps: []resource.TestStep{
			{

				Config: testAccAzureRMCosmosDbSqlStoredProcedure_basic(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbSqlStoredProcedureExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{

				Config: testAccAzureRMCosmosDbSqlStoredProcedure_update(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbSqlStoredProcedureExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMCosmosDbSqlStoredProcedureDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Cosmos.SqlClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_cosmosdb_sql_stored_procedure" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["account_name"]
		resourceGroupName := rs.Primary.Attributes["resource_group_name"]
		databaseName := rs.Primary.Attributes["database_name"]
		containerName := rs.Primary.Attributes["container_name"]

		resp, err := client.GetSQLStoredProcedure(ctx, resourceGroupName, accountName, databaseName, containerName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Error checking destroy for Cosmos SQL Stored Procedure %s (account %s) still exists:\n%v", name, accountName, err)
			}
		}

		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Cosmos SQL Stored Procedure %s (account %s) still exists:\n%#v", name, accountName, resp)
		}
	}

	return nil
}

func testCheckAzureRMCosmosDbSqlStoredProcedureExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Cosmos.SqlClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["account_name"]
		resourceGroupName := rs.Primary.Attributes["resource_group_name"]
		databaseName := rs.Primary.Attributes["database_name"]
		containerName := rs.Primary.Attributes["container_name"]

		resp, err := client.GetSQLStoredProcedure(ctx, resourceGroupName, accountName, databaseName, containerName, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on cosmosAccountsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Cosmos Stored Procedure '%s' (account: '%s') does not exist", name, accountName)
		}

		return nil
	}
}

func testAccAzureRMCosmosDbSqlStoredProcedure_base(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%[2]d"
  location = "%[1]s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-cosmos-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "%[3]s"

  consistency_policy {
    consistency_level = "%[4]s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}

resource "azurerm_cosmosdb_sql_database" "test" {
  name                = "acctest-sql-database-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
}

resource "azurerm_cosmosdb_sql_container" "test" {
  name                = "acctest-sql-container-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_sql_database.test.name
}
`, data.Locations.Primary, data.RandomInteger, string(documentdb.GlobalDocumentDB), string(documentdb.Session))
}

func testAccAzureRMCosmosDbSqlStoredProcedure_basic(data acceptance.TestData) string {
	template := testAccAzureRMCosmosDbSqlStoredProcedure_base(data)

	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_sql_stored_procedure" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_sql_database.test.name
  container_name      = azurerm_cosmosdb_sql_container.test.name

  body = <<BODY
  	function () {
		var context = getContext();
		var response = context.getResponse();
		response.setBody('Hello, World');
	}
BODY
}
`, template, data.RandomInteger)
}

func testAccAzureRMCosmosDbSqlStoredProcedure_update(data acceptance.TestData) string {
	template := testAccAzureRMCosmosDbSqlStoredProcedure_base(data)

	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_sql_stored_procedure" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_sql_database.test.name
  container_name      = azurerm_cosmosdb_sql_container.test.name

  body = <<BODY
	function () {
		var context = getContext();
		var response = context.getResponse();
		response.setBody('Welcome To Sprocs in Terraform');
	}
BODY
}
`, template, data.RandomInteger)
}

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

func TestAccAzureRMCosmosDbSqlDatabase_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_sql_database", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDbSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDbSqlDatabase_basic(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbSqlDatabaseExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDbSqlDatabase_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_sql_database", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDbSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{

				Config: testAccAzureRMCosmosDbSqlDatabase_throughput(data, 700),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbSqlDatabaseExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "throughput", "700"),
				),
			},
			data.ImportStep(),
			{

				Config: testAccAzureRMCosmosDbSqlDatabase_throughput(data, 1700),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbSqlDatabaseExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "throughput", "1700"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDbSqlDatabase_autoscale(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_sql_database", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDbSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDbSqlDatabase_autoscale(data, 4000),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbSqlDatabaseExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "autoscale_settings.0.max_throughput", "4000"),
				),
			},
			data.ImportStep(),
			{

				Config: testAccAzureRMCosmosDbSqlDatabase_autoscale(data, 5000),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbSqlDatabaseExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "autoscale_settings.0.max_throughput", "5000"),
				),
			},
			data.ImportStep(),
			{

				Config: testAccAzureRMCosmosDbSqlDatabase_autoscale(data, 4000),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbSqlDatabaseExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "autoscale_settings.0.max_throughput", "4000"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDbSqlDatabase_serverless(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_sql_database", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDbSqlDatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDbSqlDatabase_serverless(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDbSqlDatabaseExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMCosmosDbSqlDatabaseDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Cosmos.SqlClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_cosmosdb_sql_database" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		account := rs.Primary.Attributes["account_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.GetSQLDatabase(ctx, resourceGroup, account, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Error checking destroy for Cosmos SQL Database %s (account %s) still exists:\n%v", name, account, err)
			}
		}

		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Cosmos SQL Database %s (account %s) still exists:\n%#v", name, account, resp)
		}
	}

	return nil
}

func testCheckAzureRMCosmosDbSqlDatabaseExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Cosmos.SqlClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		account := rs.Primary.Attributes["account_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.GetSQLDatabase(ctx, resourceGroup, account, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on cosmosAccountsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Cosmos database '%s' (account: '%s') does not exist", name, account)
		}

		return nil
	}
}

func testAccAzureRMCosmosDbSqlDatabase_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_sql_database" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
}
`, testAccAzureRMCosmosDBAccount_basic(data, documentdb.GlobalDocumentDB, documentdb.Strong), data.RandomInteger)
}

func testAccAzureRMCosmosDbSqlDatabase_throughput(data acceptance.TestData, throughput int) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_sql_database" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  throughput          = %[3]d
}
`, testAccAzureRMCosmosDBAccount_basic(data, documentdb.GlobalDocumentDB, documentdb.Strong), data.RandomInteger, throughput)
}

func testAccAzureRMCosmosDbSqlDatabase_autoscale(data acceptance.TestData, maxThroughput int) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_sql_database" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  autoscale_settings {
    max_throughput = %[3]d
  }
}
`, testAccAzureRMCosmosDBAccount_basic(data, documentdb.GlobalDocumentDB, documentdb.Strong), data.RandomInteger, maxThroughput)
}

func testAccAzureRMCosmosDbSqlDatabase_serverless(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s
resource "azurerm_cosmosdb_sql_database" "test" {
  name                = "acctest-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
}
`, testAccAzureRMCosmosDBAccount_capabilities(data, documentdb.GlobalDocumentDB, []string{"EnableServerless"}), data.RandomInteger)
}

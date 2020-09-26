package cosmos_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMCosmosDbGremlinGraph_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_gremlin_graph", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDbGremlinGraphDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDbGremlinGraph_basic(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRmCosmosDbGremlinGraphExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDbGremlinGraph_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_gremlin_graph", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDbGremlinGraphDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDbGremlinGraph_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRmCosmosDbGremlinGraphExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMCosmosDbGremlinGraph_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_cosmosdb_gremlin_graph"),
			},
		},
	})
}

func TestAccAzureRMCosmosDbGremlinGraph_customConflictResolutionPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_gremlin_graph", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDbGremlinGraphDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDbGremlinGraph_customConflictResolutionPolicy(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRmCosmosDbGremlinGraphExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDbGremlinGraph_indexPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_gremlin_graph", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDbGremlinGraphDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDbGremlinGraph_indexPolicy(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRmCosmosDbGremlinGraphExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDbGremlinGraph_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_gremlin_graph", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDbGremlinGraphDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDbGremlinGraph_update(data, 700),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRmCosmosDbGremlinGraphExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "throughput", "700"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMCosmosDbGremlinGraph_update(data, 1700),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRmCosmosDbGremlinGraphExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "throughput", "1700"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDbGremlinGraph_autoscale(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_gremlin_graph", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDbGremlinGraphDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDbGremlinGraph_autoscale(data, 4000),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRmCosmosDbGremlinGraphExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "autoscale_settings.0.max_throughput", "4000"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMCosmosDbGremlinGraph_autoscale(data, 5000),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRmCosmosDbGremlinGraphExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "autoscale_settings.0.max_throughput", "5000"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMCosmosDbGremlinGraph_autoscale(data, 4000),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRmCosmosDbGremlinGraphExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "autoscale_settings.0.max_throughput", "4000"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMCosmosDbGremlinGraphDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Cosmos.GremlinClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_cosmosdb_gremlin_graph" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		account := rs.Primary.Attributes["account_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		database := rs.Primary.Attributes["database_name"]

		resp, err := client.GetGremlinGraph(ctx, resourceGroup, account, database, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Error checking destroy for Cosmos Gremlin Graph %s (Account %s) still exists:\n%v", name, account, err)
			}
		}

		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Cosmos Gremlin Graph %s (Account %s) still exists:\n%#v", name, account, resp)
		}
	}

	return nil
}

func testCheckAzureRmCosmosDbGremlinGraphExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Cosmos.GremlinClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not fount: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		account := rs.Primary.Attributes["account_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		database := rs.Primary.Attributes["database_name"]

		resp, err := client.GetGremlinGraph(ctx, resourceGroup, account, database, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on cosmosAccountsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Cosmos Graph '%s' (Account: '%s') does not exist", name, account)
		}
		return nil
	}
}

func testAccAzureRMCosmosDbGremlinGraph_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_gremlin_graph" "test" {
  name                = "acctest-CGRPC-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_gremlin_database.test.name
  throughput          = 400

  index_policy {
    automatic      = true
    indexing_mode  = "Consistent"
    included_paths = ["/*"]
    excluded_paths = ["/\"_etag\"/?"]
  }

  conflict_resolution_policy {
    mode                     = "LastWriterWins"
    conflict_resolution_path = "/_ts"
  }
}
`, testAccAzureRMCosmosGremlinDatabase_basic(data), data.RandomInteger)
}

func testAccAzureRMCosmosDbGremlinGraph_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_gremlin_graph" "import" {
  name                = azurerm_cosmosdb_gremlin_graph.test.name
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_gremlin_database.test.name

  index_policy {
    automatic      = true
    indexing_mode  = "Consistent"
    included_paths = ["/*"]
    excluded_paths = ["/\"_etag\"/?"]
  }

  conflict_resolution_policy {
    mode                     = "LastWriterWins"
    conflict_resolution_path = "/_ts"
  }
}
`, testAccAzureRMCosmosDbGremlinGraph_basic(data))
}

func testAccAzureRMCosmosDbGremlinGraph_customConflictResolutionPolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_gremlin_graph" "test" {
  name                = "acctest-CGRPC-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_gremlin_database.test.name
  throughput          = 400

  index_policy {
    automatic      = true
    indexing_mode  = "Consistent"
    included_paths = ["/*"]
    excluded_paths = ["/\"_etag\"/?"]
  }

  conflict_resolution_policy {
    mode                          = "Custom"
    conflict_resolution_procedure = "dbs/{0}/colls/{1}/sprocs/{2}"
  }
}
`, testAccAzureRMCosmosGremlinDatabase_basic(data), data.RandomInteger)
}

func testAccAzureRMCosmosDbGremlinGraph_indexPolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_gremlin_graph" "test" {
  name                = "acctest-CGRPC-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_gremlin_database.test.name
  throughput          = 400

  index_policy {
    automatic     = false
    indexing_mode = "None"
  }

  conflict_resolution_policy {
    mode                     = "LastWriterWins"
    conflict_resolution_path = "/_ts"
  }
}
`, testAccAzureRMCosmosGremlinDatabase_basic(data), data.RandomInteger)
}

func testAccAzureRMCosmosDbGremlinGraph_update(data acceptance.TestData, throughput int) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_gremlin_graph" "test" {
  name                = "acctest-CGRPC-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_gremlin_database.test.name
  partition_key_path  = "/test"
  throughput          = %[3]d

  index_policy {
    automatic      = true
    indexing_mode  = "Consistent"
    included_paths = ["/*"]
    excluded_paths = ["/\"_etag\"/?"]
  }

  conflict_resolution_policy {
    mode                     = "LastWriterWins"
    conflict_resolution_path = "/_ts"
  }

  unique_key {
    paths = ["/definition/id1", "/definition/id2"]
  }
}
`, testAccAzureRMCosmosGremlinDatabase_basic(data), data.RandomInteger, throughput)
}

func testAccAzureRMCosmosDbGremlinGraph_autoscale(data acceptance.TestData, maxThroughput int) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_gremlin_graph" "test" {
  name                = "acctest-CGRPC-%[2]d"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_gremlin_database.test.name
  partition_key_path  = "/test"

  autoscale_settings {
    max_throughput = %[3]d
  }

  index_policy {
    automatic      = true
    indexing_mode  = "Consistent"
    included_paths = ["/*"]
    excluded_paths = ["/\"_etag\"/?"]
  }

  conflict_resolution_policy {
    mode                     = "LastWriterWins"
    conflict_resolution_path = "/_ts"
  }
}
`, testAccAzureRMCosmosGremlinDatabase_basic(data), data.RandomInteger, maxThroughput)
}

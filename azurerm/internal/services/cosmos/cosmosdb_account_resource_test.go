package cosmos_test

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/preview/cosmos-db/mgmt/2020-04-01-preview/documentdb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMCosmosDBAccount_basic_global_boundedStaleness(t *testing.T) {
	testAccAzureRMCosmosDBAccount_basicWith(t, documentdb.GlobalDocumentDB, documentdb.BoundedStaleness)
}

func TestAccAzureRMCosmosDBAccount_basic_global_consistentPrefix(t *testing.T) {
	testAccAzureRMCosmosDBAccount_basicWith(t, documentdb.GlobalDocumentDB, documentdb.ConsistentPrefix)
}

func TestAccAzureRMCosmosDBAccount_basic_global_eventual(t *testing.T) {
	testAccAzureRMCosmosDBAccount_basicWith(t, documentdb.GlobalDocumentDB, documentdb.Eventual)
}

func TestAccAzureRMCosmosDBAccount_basic_global_session(t *testing.T) {
	testAccAzureRMCosmosDBAccount_basicWith(t, documentdb.GlobalDocumentDB, documentdb.Session)
}

func TestAccAzureRMCosmosDBAccount_basic_global_strong(t *testing.T) {
	testAccAzureRMCosmosDBAccount_basicWith(t, documentdb.GlobalDocumentDB, documentdb.Strong)
}

func TestAccAzureRMCosmosDBAccount_basic_mongo_boundedStaleness(t *testing.T) {
	testAccAzureRMCosmosDBAccount_basicWith(t, documentdb.MongoDB, documentdb.BoundedStaleness)
}

func TestAccAzureRMCosmosDBAccount_basic_mongo_consistentPrefix(t *testing.T) {
	testAccAzureRMCosmosDBAccount_basicWith(t, documentdb.MongoDB, documentdb.ConsistentPrefix)
}

func TestAccAzureRMCosmosDBAccount_basic_mongo_eventual(t *testing.T) {
	testAccAzureRMCosmosDBAccount_basicWith(t, documentdb.MongoDB, documentdb.Eventual)
}

func TestAccAzureRMCosmosDBAccount_basic_mongo_session(t *testing.T) {
	testAccAzureRMCosmosDBAccount_basicWith(t, documentdb.MongoDB, documentdb.Session)
}

func TestAccAzureRMCosmosDBAccount_basic_mongo_strong(t *testing.T) {
	testAccAzureRMCosmosDBAccount_basicWith(t, documentdb.MongoDB, documentdb.Strong)
}

func TestAccAzureRMCosmosDBAccount_basic_parse_boundedStaleness(t *testing.T) {
	testAccAzureRMCosmosDBAccount_basicWith(t, documentdb.MongoDB, documentdb.BoundedStaleness)
}

func TestAccAzureRMCosmosDBAccount_basic_parse_consistentPrefix(t *testing.T) {
	testAccAzureRMCosmosDBAccount_basicWith(t, documentdb.MongoDB, documentdb.ConsistentPrefix)
}

func TestAccAzureRMCosmosDBAccount_basic_parse_eventual(t *testing.T) {
	testAccAzureRMCosmosDBAccount_basicWith(t, documentdb.MongoDB, documentdb.Eventual)
}

func TestAccAzureRMCosmosDBAccount_basic_parse_session(t *testing.T) {
	testAccAzureRMCosmosDBAccount_basicWith(t, documentdb.MongoDB, documentdb.Session)
}

func TestAccAzureRMCosmosDBAccount_basic_parse_strong(t *testing.T) {
	testAccAzureRMCosmosDBAccount_basicWith(t, documentdb.MongoDB, documentdb.Strong)
}

func TestAccAzureRMCosmosDBAccount_public_network_access_enabled(t *testing.T) {
	testAccAzureRMCosmosDBAccount_public_network_access_enabled(t, documentdb.MongoDB, documentdb.Strong)
}

func testAccAzureRMCosmosDBAccount_public_network_access_enabled(t *testing.T, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: checkAccAzureRMCosmosDBAccount_network_access_enabled(data, kind, consistency),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, consistency, 1),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDBAccount_key_vault_uri(t *testing.T) {
	testAccAzureRMCosmosDBAccount_key_vault_uri(t, documentdb.MongoDB, documentdb.Strong)
}

func testAccAzureRMCosmosDBAccount_key_vault_uri(t *testing.T, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: checkAccAzureRMCosmosDBAccount_key_vault_uri(data, kind, consistency),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, consistency, 1),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMCosmosDBAccount_basicWith(t *testing.T, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_basic(data, kind, consistency),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, consistency, 1),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDBAccount_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_basic(data, "GlobalDocumentDB", documentdb.Eventual),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, documentdb.Eventual, 1),
				),
			},
			{
				Config:      testAccAzureRMCosmosDBAccount_requiresImport(data, documentdb.Eventual),
				ExpectError: acceptance.RequiresImportError("azurerm_cosmosdb_account"),
			},
		},
	})
}

func TestAccAzureRMCosmosDBAccount_updateConsistency_global(t *testing.T) {
	testAccAzureRMCosmosDBAccount_updateConsistency(t, documentdb.GlobalDocumentDB)
}

func TestAccAzureRMCosmosDBAccount_updateConsistency_mongo(t *testing.T) {
	testAccAzureRMCosmosDBAccount_updateConsistency(t, documentdb.MongoDB)
}

func testAccAzureRMCosmosDBAccount_updateConsistency(t *testing.T, kind documentdb.DatabaseAccountKind) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_basic(data, kind, documentdb.Strong),
				Check:  checkAccAzureRMCosmosDBAccount_basic(data, documentdb.Strong, 1),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMCosmosDBAccount_consistency(data, kind, documentdb.Strong, 8, 880),
				Check:  checkAccAzureRMCosmosDBAccount_basic(data, documentdb.Strong, 1),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMCosmosDBAccount_basic(data, kind, documentdb.BoundedStaleness),
				Check:  checkAccAzureRMCosmosDBAccount_basic(data, documentdb.BoundedStaleness, 1),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMCosmosDBAccount_consistency(data, kind, documentdb.BoundedStaleness, 7, 770),
				Check:  checkAccAzureRMCosmosDBAccount_basic(data, documentdb.BoundedStaleness, 1),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMCosmosDBAccount_consistency(data, kind, documentdb.BoundedStaleness, 77, 700),
				Check:  checkAccAzureRMCosmosDBAccount_basic(data, documentdb.BoundedStaleness, 1),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMCosmosDBAccount_basic(data, kind, documentdb.ConsistentPrefix),
				Check:  checkAccAzureRMCosmosDBAccount_basic(data, documentdb.ConsistentPrefix, 1),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDBAccount_complete_mongo(t *testing.T) {
	testAccAzureRMCosmosDBAccount_completeWith(t, documentdb.MongoDB)
}

func TestAccAzureRMCosmosDBAccount_complete_global(t *testing.T) {
	testAccAzureRMCosmosDBAccount_completeWith(t, documentdb.GlobalDocumentDB)
}

func TestAccAzureRMCosmosDBAccount_complete_parse(t *testing.T) {
	testAccAzureRMCosmosDBAccount_completeWith(t, documentdb.Parse)
}

func testAccAzureRMCosmosDBAccount_completeWith(t *testing.T, kind documentdb.DatabaseAccountKind) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_complete(data, kind, documentdb.Eventual),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, documentdb.Eventual, 3),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDBAccount_completeZoneRedundant_mongo(t *testing.T) {
	testAccAzureRMCosmosDBAccount_zoneRedundantWith(t, documentdb.MongoDB)
}

func TestAccAzureRMCosmosDBAccount_completeZoneRedundant_global(t *testing.T) {
	testAccAzureRMCosmosDBAccount_zoneRedundantWith(t, documentdb.GlobalDocumentDB)
}

func TestAccAzureRMCosmosDBAccount_completeZoneRedundant_parse(t *testing.T) {
	testAccAzureRMCosmosDBAccount_zoneRedundantWith(t, documentdb.Parse)
}

func testAccAzureRMCosmosDBAccount_zoneRedundantWith(t *testing.T, kind documentdb.DatabaseAccountKind) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_zoneRedundant(data, kind),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDBAccountExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMCosmosDBAccount_zoneRedundant_updateWith(t *testing.T, kind documentdb.DatabaseAccountKind) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_basic(data, "GlobalDocumentDB", documentdb.Eventual),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, documentdb.Eventual, 1),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMCosmosDBAccount_zoneRedundantUpdate(data, "GlobalDocumentDB", documentdb.Eventual),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, documentdb.Eventual, 2),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDBAccount_zoneRedundant_update_mongo(t *testing.T) {
	testAccAzureRMCosmosDBAccount_zoneRedundant_updateWith(t, documentdb.MongoDB)
}

func TestAccAzureRMCosmosDBAccount_update_mongo(t *testing.T) {
	testAccAzureRMCosmosDBAccount_updateWith(t, documentdb.MongoDB)
}

func TestAccAzureRMCosmosDBAccount_update_global(t *testing.T) {
	testAccAzureRMCosmosDBAccount_updateWith(t, documentdb.GlobalDocumentDB)
}

func TestAccAzureRMCosmosDBAccount_update_parse(t *testing.T) {
	testAccAzureRMCosmosDBAccount_updateWith(t, documentdb.Parse)
}
func testAccAzureRMCosmosDBAccount_updateWith(t *testing.T, kind documentdb.DatabaseAccountKind) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_basic(data, kind, documentdb.Eventual),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, documentdb.Eventual, 1),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMCosmosDBAccount_complete(data, kind, documentdb.Eventual),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, documentdb.Eventual, 3),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMCosmosDBAccount_completeUpdated(data, kind, documentdb.Eventual),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, documentdb.Eventual, 3),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMCosmosDBAccount_basicWithResources(data, kind, documentdb.Eventual),
				Check:  resource.ComposeAggregateTestCheckFunc(
				// checkAccAzureRMCosmosDBAccount_basic(data, documentdb.Eventual, 1),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDBAccount_capabilities_EnableAggregationPipeline(t *testing.T) {
	testAccAzureRMCosmosDBAccount_capabilitiesWith(t, documentdb.GlobalDocumentDB, []string{"EnableAggregationPipeline"})
}

func TestAccAzureRMCosmosDBAccount_capabilities_EnableCassandra(t *testing.T) {
	testAccAzureRMCosmosDBAccount_capabilitiesWith(t, documentdb.GlobalDocumentDB, []string{"EnableCassandra"})
}

func TestAccAzureRMCosmosDBAccount_capabilities_EnableGremlin(t *testing.T) {
	testAccAzureRMCosmosDBAccount_capabilitiesWith(t, documentdb.GlobalDocumentDB, []string{"EnableGremlin"})
}

func TestAccAzureRMCosmosDBAccount_capabilities_EnableTable(t *testing.T) {
	testAccAzureRMCosmosDBAccount_capabilitiesWith(t, documentdb.GlobalDocumentDB, []string{"EnableTable"})
}

func TestAccAzureRMCosmosDBAccount_capabilities_EnableServerless(t *testing.T) {
	testAccAzureRMCosmosDBAccount_capabilitiesWith(t, documentdb.GlobalDocumentDB, []string{"EnableServerless"})
}

func TestAccAzureRMCosmosDBAccount_capabilities_EnableMongo(t *testing.T) {
	testAccAzureRMCosmosDBAccount_capabilitiesWith(t, documentdb.MongoDB, []string{"EnableMongo"})
}

func TestAccAzureRMCosmosDBAccount_capabilities_MongoDBv34(t *testing.T) {
	testAccAzureRMCosmosDBAccount_capabilitiesWith(t, documentdb.MongoDB, []string{"MongoDBv3.4"})
}

func TestAccAzureRMCosmosDBAccount_capabilities_mongoEnableDocLevelTTL(t *testing.T) {
	testAccAzureRMCosmosDBAccount_capabilitiesWith(t, documentdb.MongoDB, []string{"mongoEnableDocLevelTTL"})
}

func TestAccAzureRMCosmosDBAccount_capabilities_DisableRateLimitingResponses(t *testing.T) {
	testAccAzureRMCosmosDBAccount_capabilitiesWith(t, documentdb.MongoDB, []string{"DisableRateLimitingResponses"})
}

func TestAccAzureRMCosmosDBAccount_capabilities_AllowSelfServeUpgradeToMongo36(t *testing.T) {
	testAccAzureRMCosmosDBAccount_capabilitiesWith(t, documentdb.MongoDB, []string{"AllowSelfServeUpgradeToMongo36"})
}

func testAccAzureRMCosmosDBAccount_capabilitiesWith(t *testing.T, kind documentdb.DatabaseAccountKind, capabilities []string) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_capabilities(data, kind, capabilities),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, documentdb.Strong, 1),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDBAccount_capabilitiesAdd(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_capabilities(data, documentdb.GlobalDocumentDB, []string{"EnableCassandra"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, documentdb.Strong, 1),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMCosmosDBAccount_capabilities(data, documentdb.GlobalDocumentDB, []string{"EnableCassandra", "EnableAggregationPipeline"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, documentdb.Strong, 1),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDBAccount_capabilitiesUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_capabilities(data, documentdb.GlobalDocumentDB, []string{"EnableCassandra"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, documentdb.Strong, 1),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMCosmosDBAccount_capabilities(data, documentdb.GlobalDocumentDB, []string{"EnableTable", "EnableAggregationPipeline"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, documentdb.Strong, 1),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDBAccount_geoLocationsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_basic(data, "GlobalDocumentDB", documentdb.Eventual),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, documentdb.Eventual, 1),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMCosmosDBAccount_geoLocationUpdate(data, "GlobalDocumentDB", documentdb.Eventual),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, documentdb.Eventual, 2),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMCosmosDBAccount_basic(data, "GlobalDocumentDB", documentdb.Eventual),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, documentdb.Eventual, 1),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDBAccount_freeTier(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_freeTier(data, "GlobalDocumentDB", documentdb.Eventual),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, documentdb.Eventual, 1),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_free_tier", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDBAccount_vNetFilters(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_vNetFilters(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckAzureRMCosmosDBAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "is_virtual_network_filter_enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "virtual_network_rule.#", "2"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMCosmosDBAccountDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Cosmos.DatabaseClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_cosmosdb_account" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("CosmosDB Account still exists:\n%#v", resp)
		}
	}

	return nil
}

func testCheckAzureRMCosmosDBAccountExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Cosmos.DatabaseClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on cosmosAccountsClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: CosmosDB Account '%s' (resource group: '%s') does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testAccAzureRMCosmosDBAccount_basic(data acceptance.TestData, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%d"
  location = "%s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "%s"

  consistency_policy {
    consistency_level = "%s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(kind), string(consistency))
}

func testAccAzureRMCosmosDBAccount_requiresImport(data acceptance.TestData, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_account" "import" {
  name                = azurerm_cosmosdb_account.test.name
  location            = azurerm_cosmosdb_account.test.location
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  offer_type          = azurerm_cosmosdb_account.test.offer_type

  consistency_policy {
    consistency_level = azurerm_cosmosdb_account.test.consistency_policy[0].consistency_level
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, testAccAzureRMCosmosDBAccount_basic(data, "GlobalDocumentDB", consistency))
}

func testAccAzureRMCosmosDBAccount_consistency(data acceptance.TestData, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel, interval, staleness int) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%d"
  location = "%s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "%s"

  consistency_policy {
    consistency_level       = "%s"
    max_interval_in_seconds = %d
    max_staleness_prefix    = %d
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(kind), string(consistency), interval, staleness)
}

func testAccAzureRMCosmosDBAccount_completePreReqs(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VNET-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  dns_servers         = ["10.0.0.4", "10.0.0.5"]
}

resource "azurerm_subnet" "subnet1" {
  name                 = "acctest-SN1-%[1]d-1"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.1.0/24"
  service_endpoints    = ["Microsoft.AzureCosmosDB"]
}

resource "azurerm_subnet" "subnet2" {
  name                 = "acctest-SN2-%[1]d-2"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
  service_endpoints    = ["Microsoft.AzureCosmosDB"]
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMCosmosDBAccount_complete(data acceptance.TestData, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "%[3]s"

  consistency_policy {
    consistency_level       = "%[4]s"
    max_interval_in_seconds = 300
    max_staleness_prefix    = 170000
  }

  is_virtual_network_filter_enabled = true

  virtual_network_rule {
    id = azurerm_subnet.subnet1.id
  }

  virtual_network_rule {
    id = azurerm_subnet.subnet2.id
  }

  enable_multiple_write_locations = true

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }

  geo_location {
    location          = "%[5]s"
    failover_priority = 1
  }

  geo_location {
    location          = "%[6]s"
    failover_priority = 2
  }
}
`, testAccAzureRMCosmosDBAccount_completePreReqs(data), data.RandomInteger, string(kind), string(consistency), data.Locations.Secondary, data.Locations.Ternary)
}

func testAccAzureRMCosmosDBAccount_zoneRedundant(data acceptance.TestData, kind documentdb.DatabaseAccountKind) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%d"
  location = "%s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "%s"

  enable_multiple_write_locations = true

  consistency_policy {
    consistency_level       = "BoundedStaleness"
    max_interval_in_seconds = 300
    max_staleness_prefix    = 100000
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }

  geo_location {
    location          = "%s"
    failover_priority = 1
    zone_redundant    = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(kind), data.Locations.Secondary)
}

func testAccAzureRMCosmosDBAccount_completeUpdated(data acceptance.TestData, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "%[3]s"

  consistency_policy {
    consistency_level       = "%[4]s"
    max_interval_in_seconds = 360
    max_staleness_prefix    = 170000
  }

  is_virtual_network_filter_enabled = true

  virtual_network_rule {
    id = azurerm_subnet.subnet2.id
  }

  enable_multiple_write_locations = true

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }

  geo_location {
    location          = "%[5]s"
    failover_priority = 1
  }

  geo_location {
    location          = "%[6]s"
    failover_priority = 2
  }
}
`, testAccAzureRMCosmosDBAccount_completePreReqs(data), data.RandomInteger, string(kind), string(consistency), data.Locations.Secondary, data.Locations.Ternary)
}

func testAccAzureRMCosmosDBAccount_basicWithResources(data acceptance.TestData, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "%s"

  consistency_policy {
    consistency_level = "%s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, testAccAzureRMCosmosDBAccount_completePreReqs(data), data.RandomInteger, string(kind), string(consistency))
}

func testAccAzureRMCosmosDBAccount_capabilities(data acceptance.TestData, kind documentdb.DatabaseAccountKind, capabilities []string) string {
	capeTf := ""
	for _, c := range capabilities {
		capeTf += fmt.Sprintf("capabilities {name = \"%s\"}\n", c)
	}

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%d"
  location = "%s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "%s"

  consistency_policy {
    consistency_level = "Strong"
  }

  %s

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(kind), capeTf)
}

func testAccAzureRMCosmosDBAccount_geoLocationUpdate(data acceptance.TestData, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%d"
  location = "%s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "%s"

  consistency_policy {
    consistency_level = "%s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }

  geo_location {
    location          = "%s"
    failover_priority = 1
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(kind), string(consistency), data.Locations.Secondary)
}

func testAccAzureRMCosmosDBAccount_zoneRedundantUpdate(data acceptance.TestData, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
variable "geo_location" {
  type = list(object({
    location          = string
    failover_priority = string
    zone_redundant    = bool
  }))
  default = [
    {
      location          = "%s"
      failover_priority = 0
      zone_redundant    = false
    },
    {
      location          = "%s"
      failover_priority = 1
      zone_redundant    = true
    }
  ]
}

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%d"
  location = "%s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "%s"
  
  enable_multiple_write_locations = true
  enable_automatic_failover       = true

  consistency_policy {
    consistency_level = "%s"
  }

  dynamic "geo_location" {
    for_each = var.geo_location
    content {
      location          = geo_location.value.location
      failover_priority = geo_location.value.failover_priority
      zone_redundant    = geo_location.value.zone_redundant
    }
  }
}
`, data.Locations.Primary, data.Locations.Secondary, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(kind), string(consistency))
}

func testAccAzureRMCosmosDBAccount_vNetFiltersPreReqs(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VNET-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  dns_servers         = ["10.0.0.4", "10.0.0.5"]
}

resource "azurerm_subnet" "subnet1" {
  name                                           = "acctest-SN1-%[1]d-1"
  resource_group_name                            = azurerm_resource_group.test.name
  virtual_network_name                           = azurerm_virtual_network.test.name
  address_prefixes                               = ["10.0.1.0/24"]
  enforce_private_link_endpoint_network_policies = false
  enforce_private_link_service_network_policies  = false
}

resource "azurerm_subnet" "subnet2" {
  name                                           = "acctest-SN2-%[1]d-2"
  resource_group_name                            = azurerm_resource_group.test.name
  virtual_network_name                           = azurerm_virtual_network.test.name
  address_prefixes                               = ["10.0.2.0/24"]
  service_endpoints                              = ["Microsoft.AzureCosmosDB"]
  enforce_private_link_endpoint_network_policies = false
  enforce_private_link_service_network_policies  = false
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMCosmosDBAccount_vNetFilters(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  enable_multiple_write_locations = false
  enable_automatic_failover       = false

  consistency_policy {
    consistency_level       = "Eventual"
    max_interval_in_seconds = 5
    max_staleness_prefix    = 100
  }

  is_virtual_network_filter_enabled = true
  ip_range_filter                   = ""

  virtual_network_rule {
    id                                   = azurerm_subnet.subnet1.id
    ignore_missing_vnet_service_endpoint = true
  }

  virtual_network_rule {
    id                                   = azurerm_subnet.subnet2.id
    ignore_missing_vnet_service_endpoint = false
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, testAccAzureRMCosmosDBAccount_vNetFiltersPreReqs(data), data.RandomInteger)
}

func testAccAzureRMCosmosDBAccount_freeTier(data acceptance.TestData, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%d"
  location = "%s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "%s"

  enable_free_tier = true

  consistency_policy {
    consistency_level = "%s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(kind), string(consistency))
}

func checkAccAzureRMCosmosDBAccount_basic(data acceptance.TestData, consistency documentdb.DefaultConsistencyLevel, locationCount int) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		testCheckAzureRMCosmosDBAccountExists(data.ResourceName),
		resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
		resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
		resource.TestCheckResourceAttr(data.ResourceName, "location", azure.NormalizeLocation(data.Locations.Primary)),
		resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
		resource.TestCheckResourceAttr(data.ResourceName, "offer_type", string(documentdb.Standard)),
		resource.TestCheckResourceAttr(data.ResourceName, "consistency_policy.0.consistency_level", string(consistency)),
		resource.TestCheckResourceAttr(data.ResourceName, "geo_location.#", strconv.Itoa(locationCount)),
		resource.TestCheckResourceAttrSet(data.ResourceName, "endpoint"),
		resource.TestCheckResourceAttr(data.ResourceName, "read_endpoints.#", strconv.Itoa(locationCount)),
		resource.TestCheckResourceAttrSet(data.ResourceName, "primary_key"),
		resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_key"),
		resource.TestCheckResourceAttrSet(data.ResourceName, "primary_readonly_key"),
		resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_readonly_key"),
	)
}

func checkAccAzureRMCosmosDBAccount_network_access_enabled(data acceptance.TestData, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%d"
  location = "%s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                          = "acctest-ca-%d"
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  offer_type                    = "Standard"
  kind                          = "%s"
  public_network_access_enabled = true

  consistency_policy {
    consistency_level = "%s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, string(kind), string(consistency))
}

func checkAccAzureRMCosmosDBAccount_key_vault_uri(data acceptance.TestData, kind documentdb.DatabaseAccountKind, consistency documentdb.DefaultConsistencyLevel) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cosmos-%d"
  location = "%s"
}

data "azuread_service_principal" "cosmosdb" {
  display_name = "Azure Cosmos DB"
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "premium"

  purge_protection_enabled = true
  soft_delete_enabled      = true

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "list",
      "create",
      "delete",
      "get",
      "update",
    ]

    secret_permissions = [
      "get",
      "delete",
      "set",
    ]
  }

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azuread_service_principal.cosmosdb.id

    key_permissions = [
      "list",
      "create",
      "delete",
      "get",
      "update",
      "unwrapKey",
      "wrapKey",
    ]

    secret_permissions = [
      "get",
      "delete",
      "set",
    ]
  }
}

resource "azurerm_key_vault_key" "test" {
  name         = "key-%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-ca-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "%s"
  key_vault_key_id    = azurerm_key_vault_key.test.id

  consistency_policy {
    consistency_level = "%s"
  }

  geo_location {
    location          = azurerm_resource_group.test.location
    failover_priority = 0
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomInteger, string(kind), string(consistency))
}

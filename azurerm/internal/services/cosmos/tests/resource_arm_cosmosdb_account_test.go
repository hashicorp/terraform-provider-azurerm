package tests

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2015-04-08/documentdb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

// TODO: refactor the test configs

//consistency
func TestAccAzureRMCosmosDBAccount_eventualConsistency(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_basic(data, string(documentdb.Eventual), "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, string(documentdb.Eventual), 1),
				),
			},
			data.ImportStep(),
		},
	})
}
func TestAccAzureRMCosmosDBAccount_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_basic(data, string(documentdb.Eventual), "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, string(documentdb.Eventual), 1),
				),
			},
			{
				Config:      testAccAzureRMCosmosDBAccount_requiresImport(data, string(documentdb.Eventual), "", ""),
				ExpectError: acceptance.RequiresImportError("azurerm_cosmosdb_account"),
			},
		},
	})
}

func TestAccAzureRMCosmosDBAccount_session(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_basic(data, string(documentdb.Session), "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, string(documentdb.Session), 1),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDBAccount_strong(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_basic(data, string(documentdb.Strong), "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, string(documentdb.Strong), 1),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDBAccount_consistentPrefix(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_basic(data, string(documentdb.ConsistentPrefix), "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, string(documentdb.ConsistentPrefix), 1),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDBAccount_boundedStaleness(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), 1),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDBAccount_boundedStaleness_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_boundedStaleness_complete(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), 1),
					resource.TestCheckResourceAttr(data.ResourceName, "consistency_policy.0.max_interval_in_seconds", "10"),
					resource.TestCheckResourceAttr(data.ResourceName, "consistency_policy.0.max_staleness_prefix", "200"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDBAccount_consistency_change(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_basic(data, string(documentdb.Session), "", ""),
				Check:  checkAccAzureRMCosmosDBAccount_basic(data, string(documentdb.Session), 1),
			},
			{
				Config: testAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), "", ""),
				Check:  checkAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), 1),
			},
			data.ImportStep(),
		},
	})
}

//DB kinds
func TestAccAzureRMCosmosDBAccount_mongoDB(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_mongoDB(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), 1),
					resource.TestCheckResourceAttr(data.ResourceName, "kind", "MongoDB"),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_strings.#", "4"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDBAccount_capabilityGremlin(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_capabilityGremlin(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), 1),
					resource.TestCheckResourceAttr(data.ResourceName, "kind", "GlobalDocumentDB"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDBAccount_capabilityTable(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_capabilityTable(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), 1),
					resource.TestCheckResourceAttr(data.ResourceName, "kind", "GlobalDocumentDB"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDBAccount_capabilityCassandra(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_capabilityCassandra(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), 1),
					resource.TestCheckResourceAttr(data.ResourceName, "kind", "GlobalDocumentDB"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDBAccount_capabilityAggregationPipeline(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_capabilityAggregationPipeline(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), 1),
					resource.TestCheckResourceAttr(data.ResourceName, "kind", "GlobalDocumentDB"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDBAccount_capabilityMongo35(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_capabilityMongo34(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), 1),
					resource.TestCheckResourceAttr(data.ResourceName, "kind", "MongoDB"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDBAccount_capabilityDocLevelTTL(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_capabilityDocLevelTTL(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), 1),
					resource.TestCheckResourceAttr(data.ResourceName, "kind", "MongoDB"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDBAccount_capabilityUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_capabilityDocLevelTTL(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), 1),
					resource.TestCheckResourceAttr(data.ResourceName, "kind", "MongoDB"),
				),
			},
			{
				Config: testAccAzureRMCosmosDBAccount_capabilityMongo34(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), 1),
					resource.TestCheckResourceAttr(data.ResourceName, "kind", "MongoDB"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAbcAzureRMCosmosDBAccount_updatePropertiesAndLocation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_basic(data, string(documentdb.Session), "", ""),
				Check:  checkAccAzureRMCosmosDBAccount_basic(data, string(documentdb.Session), 1),
			},
			{
				Config: testAccAzureRMCosmosDBAccount_geoReplicated_customId(data),
				Check:  checkAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), 2),
			},
			data.ImportStep(),
		},
	})
}

//replication
func TestAccAzureRMCosmosDBAccount_geoReplicated_customId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_geoReplicated_customId(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), 2),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDBAccount_geoReplicated_non_boundedStaleness_cp(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_geoReplicated_customConsistencyLevel(data, documentdb.Session),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, string(documentdb.Session), 2),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDBAccount_geoReplicated_add_remove(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")
	configBasic := testAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), "", "")
	configReplicated := testAccAzureRMCosmosDBAccount_geoReplicated_customId(data)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic,
				Check:  checkAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), 1),
			},
			{
				Config: configReplicated,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), 2),
				),
			},
			{
				Config: configBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), 1),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDBAccount_geoReplicated_rename(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_geoReplicated(data),
				Check:  checkAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), 2),
			},
			{
				Config: testAccAzureRMCosmosDBAccount_geoReplicated_customId(data),
				Check:  checkAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), 2),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDBAccount_virtualNetworkFilter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_virtualNetworkFilter(data),
				Check:  checkAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), 1),
			},
			data.ImportStep(),
		},
	})
}

//basic --> complete (
func TestAccAzureRMCosmosDBAccount_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_basic(data, string(documentdb.Session), "", ""),
				Check:  checkAccAzureRMCosmosDBAccount_basic(data, string(documentdb.Session), 1),
			},
			{
				Config: testAccAzureRMCosmosDBAccount_complete(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), 2),
					resource.TestCheckResourceAttr(data.ResourceName, "ip_range_filter", "104.42.195.92,40.76.54.131,52.176.6.30,52.169.50.45/32,52.187.184.26,10.20.0.0/16"),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_automatic_failover", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDBAccount_emptyIpFilter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_emptyIpFilter(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), 2),
					resource.TestCheckResourceAttr(data.ResourceName, "ip_range_filter", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_automatic_failover", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDBAccount_multiMaster(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_multiMaster(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "geo_location.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "write_endpoints.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_multiple_write_locations", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMCosmosDBAccount_multiMaster_geoReplicated_zoned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_multiMaster_geoReplicated_zoned(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "geo_location.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "write_endpoints.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "enable_multiple_write_locations", "true"),
					testCheckZoneRedundantCount(data.ResourceName, 1, 1),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckZoneRedundantCount(resourceName string, enabledCount int, disabledCount int) func(*terraform.State) error {
	return func(state *terraform.State) error {
		r := regexp.MustCompile(`(geo_location)(.)\\d+(.)(zone_redundant)`)
		root := state.RootModule().Resources[resourceName].Primary

		foundEnabled := 0
		foundDisabled := 0

		for k, v := range root.Attributes {
			if r.MatchString(k) {
				switch v {
				case "true":
					foundEnabled++
				case "false":
					foundDisabled++
				default:
					return fmt.Errorf("unexpected boolean value found: %s", v)
				}
			}
		}

		if foundEnabled != enabledCount || foundDisabled != disabledCount {
			return fmt.Errorf("unexpected number of enabled(%d) and disabled(%d) `zone_redudant` flags", foundEnabled, foundDisabled)
		}

		return nil
	}
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

func testAccAzureRMCosmosDBAccount_basic(data acceptance.TestData, consistency string, consistencyOptions string, additional string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  offer_type          = "Standard"

  consistency_policy {
    consistency_level = "%s"

    %s
  }

  geo_location {
    location          = "${azurerm_resource_group.test.location}"
    failover_priority = 0
  }

  %s
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, consistency, consistencyOptions, additional)
}

func testAccAzureRMCosmosDBAccount_requiresImport(data acceptance.TestData, consistency string, consistencyOptions string, additional string) string {
	template := testAccAzureRMCosmosDBAccount_basic(data, consistency, consistencyOptions, additional)
	return fmt.Sprintf(`
%s

resource "azurerm_cosmosdb_account" "import" {
  name                = "${azurerm_cosmosdb_account.test.name}"
  location            = "${azurerm_cosmosdb_account.test.location}"
  resource_group_name = "${azurerm_cosmosdb_account.test.resource_group_name}"
  offer_type          = "Standard"

  consistency_policy {
    consistency_level = "%s"

    %s
  }

  geo_location {
    location          = "${azurerm_resource_group.test.location}"
    failover_priority = 0
  }

  %s
}
`, template, consistency, consistencyOptions, additional)
}

func testAccAzureRMCosmosDBAccount_boundedStaleness_complete(data acceptance.TestData) string {
	return testAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), `
        max_interval_in_seconds = 10
        max_staleness_prefix    = 200
`, "")
}

func testAccAzureRMCosmosDBAccount_mongoDB(data acceptance.TestData) string {
	return testAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), "", `
        kind = "MongoDB"
    `)
}

func testAccAzureRMCosmosDBAccount_capabilityGremlin(data acceptance.TestData) string {
	return testAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), "", `
        kind = "GlobalDocumentDB"

        capabilities {
          name = "EnableGremlin"
        }
    `)
}

func testAccAzureRMCosmosDBAccount_capabilityTable(data acceptance.TestData) string {
	return testAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), "", `
        kind = "GlobalDocumentDB"

        capabilities {
          name = "EnableTable"
        }
    `)
}

func testAccAzureRMCosmosDBAccount_capabilityCassandra(data acceptance.TestData) string {
	return testAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), "", `
        kind = "GlobalDocumentDB"

        capabilities {
          name = "EnableCassandra"
        }
    `)
}

func testAccAzureRMCosmosDBAccount_capabilityAggregationPipeline(data acceptance.TestData) string {
	return testAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), "", `
        kind = "GlobalDocumentDB"

        capabilities {
          name = "EnableAggregationPipeline"
        }
    `)
}

func testAccAzureRMCosmosDBAccount_capabilityMongo34(data acceptance.TestData) string {
	return testAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), "", `
        kind = "MongoDB"

        capabilities {
          name = "MongoDBv3.4"
        }
    `)
}

func testAccAzureRMCosmosDBAccount_capabilityDocLevelTTL(data acceptance.TestData) string {
	return testAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), "", `
        kind = "MongoDB"

        capabilities {
          name = "mongoEnableDocLevelTTL"
        }
    `)
}

func testAccAzureRMCosmosDBAccount_geoReplicated(data acceptance.TestData) string {
	co := `
	max_interval_in_seconds = 373
	max_staleness_prefix    = 100001
`

	return testAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), co, fmt.Sprintf(`
        geo_location {
            location          = "%s"
            failover_priority = 1
        }

  `, data.Locations.Secondary))
}

func testAccAzureRMCosmosDBAccount_multiMaster(data acceptance.TestData) string {
	co := `
	max_interval_in_seconds = 373
	max_staleness_prefix    = 100001
`

	return testAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), co, fmt.Sprintf(`
        enable_multiple_write_locations = true

        geo_location {
            location          = "%s"
            failover_priority = 1
        }

  `, data.Locations.Secondary))
}

func testAccAzureRMCosmosDBAccount_multiMaster_geoReplicated_zoned(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  consistency_policy {
    consistency_level = "Eventual"
  }

  enable_multiple_write_locations = true

  geo_location {
    location          = "%s"
    failover_priority = 0
    zone_redundant    = true
  }

  geo_location {
    location          = "%s"
    failover_priority = 1
    zone_redundant    = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func testAccAzureRMCosmosDBAccount_geoReplicated_customId(data acceptance.TestData) string {
	co := `
	max_interval_in_seconds = 373
	max_staleness_prefix    = 100001
`

	return testAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), co, fmt.Sprintf(`
        geo_location {
            prefix            = "acctest-%d-custom-id"
            location          = "%s"
            failover_priority = 1
        }

  `, data.RandomInteger, data.Locations.Secondary))
}

func testAccAzureRMCosmosDBAccount_geoReplicated_customConsistencyLevel(data acceptance.TestData, cLevel documentdb.DefaultConsistencyLevel) string {
	return testAccAzureRMCosmosDBAccount_basic(data, string(cLevel), "", fmt.Sprintf(`
        geo_location {
            prefix            = "acctest-%d-custom-id"
            location          = "%s"
            failover_priority = 1
        }

  `, data.RandomInteger, data.Locations.Secondary))
}

func testAccAzureRMCosmosDBAccount_complete(data acceptance.TestData) string {
	co := `
	max_interval_in_seconds = 373
	max_staleness_prefix    = 100001
`
	return testAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), co, fmt.Sprintf(`
		ip_range_filter				= "104.42.195.92,40.76.54.131,52.176.6.30,52.169.50.45/32,52.187.184.26,10.20.0.0/16"
		enable_automatic_failover	= true

        geo_location {
            prefix            = "acctest-%d-custom-id"
            location          = "%s"
            failover_priority = 1
        }
  `, data.RandomInteger, data.Locations.Secondary))
}

func testAccAzureRMCosmosDBAccount_emptyIpFilter(data acceptance.TestData) string {
	co := `
	max_interval_in_seconds = 373
	max_staleness_prefix    = 100001
`

	return testAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), co, fmt.Sprintf(`
		ip_range_filter				= ""
		enable_automatic_failover	= true

        geo_location {
            prefix            = "acctest-%d-custom-id"
            location          = "%s"
            failover_priority = 1
        }
  `, data.RandomInteger, data.Locations.Secondary))
}

func testAccAzureRMCosmosDBAccount_virtualNetworkFilter(data acceptance.TestData) string {
	vnetConfig := fmt.Sprintf(`
resource "azurerm_virtual_network" "test" {
  name                = "acctest-%[1]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  dns_servers         = ["10.0.0.4", "10.0.0.5"]
}

resource "azurerm_subnet" "subnet1" {
  name                 = "acctest-%[1]d-1"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.1.0/24"
  service_endpoints    = ["Microsoft.AzureCosmosDB"]
}

resource "azurerm_subnet" "subnet2" {
  name                 = "acctest-%[1]d-2"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
  service_endpoints    = ["Microsoft.AzureCosmosDB"]
}
`, data.RandomInteger)

	basic := testAccAzureRMCosmosDBAccount_basic(data, string(documentdb.BoundedStaleness), "", `
        is_virtual_network_filter_enabled = true

        virtual_network_rule {
          id = "${azurerm_subnet.subnet1.id}"
        }

        virtual_network_rule {
          id = "${azurerm_subnet.subnet2.id}"
        }
	`)

	return vnetConfig + basic
}

func checkAccAzureRMCosmosDBAccount_basic(data acceptance.TestData, consistency string, locationCount int) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		testCheckAzureRMCosmosDBAccountExists(data.ResourceName),
		resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
		resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
		resource.TestCheckResourceAttr(data.ResourceName, "location", azure.NormalizeLocation(data.Locations.Primary)),
		resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
		resource.TestCheckResourceAttr(data.ResourceName, "offer_type", string(documentdb.Standard)),
		resource.TestCheckResourceAttr(data.ResourceName, "consistency_policy.0.consistency_level", consistency),
		resource.TestCheckResourceAttr(data.ResourceName, "geo_location.#", strconv.Itoa(locationCount)),
		resource.TestCheckResourceAttrSet(data.ResourceName, "endpoint"),
		resource.TestCheckResourceAttr(data.ResourceName, "read_endpoints.#", strconv.Itoa(locationCount)),
		resource.TestCheckResourceAttr(data.ResourceName, "write_endpoints.#", "1"),
		resource.TestCheckResourceAttr(data.ResourceName, "enable_multiple_write_locations", "false"),
		resource.TestCheckResourceAttrSet(data.ResourceName, "primary_master_key"),
		resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_master_key"),
		resource.TestCheckResourceAttrSet(data.ResourceName, "primary_readonly_master_key"),
		resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_readonly_master_key"),
	)
}

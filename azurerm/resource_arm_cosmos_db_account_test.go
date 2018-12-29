package azurerm

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2015-04-08/documentdb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

// TODO: refactor the test configs

func init() {
	resource.AddTestSweepers("azurerm_cosmosdb_account", &resource.Sweeper{
		Name: "azurerm_cosmosdb_account",
		F:    testSweepCosmosDBAccount,
	})
}

func testSweepCosmosDBAccount(region string) error {
	armClient, err := buildConfigForSweepers()
	if err != nil {
		return err
	}

	client := (*armClient).cosmosDBClient
	ctx := (*armClient).StopContext

	log.Printf("Retrieving the CosmosDB Accounts..")
	results, err := client.List(ctx)
	if err != nil {
		return fmt.Errorf("Error Listing on CosmosDB Accounts: %+v", err)
	}

	for _, account := range *results.Value {
		if !shouldSweepAcceptanceTestResource(*account.Name, *account.Location, region) {
			continue
		}

		resourceId, err := parseAzureResourceID(*account.ID)
		if err != nil {
			return err
		}

		resourceGroup := resourceId.ResourceGroup
		name := resourceId.Path["databaseAccounts"]

		log.Printf("Deleting CosmosDB Account '%s' in Resource Group '%s'", name, resourceGroup)
		future, err := client.Delete(ctx, resourceGroup, name)
		if err != nil {
			return err
		}
		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return err
		}
	}

	return nil
}

//consistency
func TestAccAzureRMCosmosDBAccount_eventualConsistency(t *testing.T) {
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_cosmosdb_account.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_basic(ri, testLocation(), string(documentdb.Eventual), "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(resourceName, testLocation(), string(documentdb.Eventual), 1),
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
func TestAccAzureRMCosmosDBAccount_requiresImport(t *testing.T) {
	if !requireResourcesToBeImported {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_cosmosdb_account.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_basic(ri, location, string(documentdb.Eventual), "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(resourceName, testLocation(), string(documentdb.Eventual), 1),
				),
			},
			{
				Config:      testAccAzureRMCosmosDBAccount_requiresImport(ri, location, string(documentdb.Eventual), "", ""),
				ExpectError: testRequiresImportError("azurerm_cosmosdb_account"),
			},
		},
	})
}

func TestAccAzureRMCosmosDBAccount_session(t *testing.T) {
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_cosmosdb_account.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_basic(ri, testLocation(), string(documentdb.Session), "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(resourceName, testLocation(), string(documentdb.Session), 1),
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

func TestAccAzureRMCosmosDBAccount_strong(t *testing.T) {
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_cosmosdb_account.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_basic(ri, testLocation(), string(documentdb.Strong), "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(resourceName, testLocation(), string(documentdb.Strong), 1),
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

func TestAccAzureRMCosmosDBAccount_consistentPrefix(t *testing.T) {
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_cosmosdb_account.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_basic(ri, testLocation(), string(documentdb.ConsistentPrefix), "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(resourceName, testLocation(), string(documentdb.ConsistentPrefix), 1),
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

func TestAccAzureRMCosmosDBAccount_boundedStaleness(t *testing.T) {
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_cosmosdb_account.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_basic(ri, testLocation(), string(documentdb.BoundedStaleness), "", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(resourceName, testLocation(), string(documentdb.BoundedStaleness), 1),
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

func TestAccAzureRMCosmosDBAccount_boundedStaleness_complete(t *testing.T) {
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_cosmosdb_account.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_boundedStaleness_complete(ri, testLocation()),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(resourceName, testLocation(), string(documentdb.BoundedStaleness), 1),
					resource.TestCheckResourceAttr(resourceName, "consistency_policy.0.max_interval_in_seconds", "10"),
					resource.TestCheckResourceAttr(resourceName, "consistency_policy.0.max_staleness_prefix", "200"),
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

func TestAccAzureRMCosmosDBAccount_consistency_change(t *testing.T) {
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_cosmosdb_account.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_basic(ri, testLocation(), string(documentdb.Session), "", ""),
				Check:  checkAccAzureRMCosmosDBAccount_basic(resourceName, testLocation(), string(documentdb.Session), 1),
			},
			{
				Config: testAccAzureRMCosmosDBAccount_basic(ri, testLocation(), string(documentdb.BoundedStaleness), "", ""),
				Check:  checkAccAzureRMCosmosDBAccount_basic(resourceName, testLocation(), string(documentdb.BoundedStaleness), 1),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

//DB kinds
func TestAccAzureRMCosmosDBAccount_mongoDB(t *testing.T) {
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_cosmosdb_account.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_mongoDB(ri, testLocation()),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(resourceName, testLocation(), string(documentdb.BoundedStaleness), 1),
					resource.TestCheckResourceAttr(resourceName, "kind", "MongoDB"),
					resource.TestCheckResourceAttr(resourceName, "connection_strings.#", "1"),
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

func TestAccAzureRMCosmosDBAccount_gremlin(t *testing.T) {
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_cosmosdb_account.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_gremlin(ri, testLocation()),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(resourceName, testLocation(), string(documentdb.BoundedStaleness), 1),
					resource.TestCheckResourceAttr(resourceName, "kind", "GlobalDocumentDB"),
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

func TestAccAzureRMCosmosDBAccount_table(t *testing.T) {
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_cosmosdb_account.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_table(ri, testLocation()),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(resourceName, testLocation(), string(documentdb.BoundedStaleness), 1),
					resource.TestCheckResourceAttr(resourceName, "kind", "GlobalDocumentDB"),
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

func TestAccAzureRMCosmosDBAccount_updatePropertiesAndLocation(t *testing.T) {
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_cosmosdb_account.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_basic(ri, testLocation(), string(documentdb.Session), "", ""),
				Check:  checkAccAzureRMCosmosDBAccount_basic(resourceName, testLocation(), string(documentdb.Session), 1),
			},
			{
				Config: testAccAzureRMCosmosDBAccount_geoReplicated_customId(ri, testLocation(), testAltLocation()),
				Check:  checkAccAzureRMCosmosDBAccount_basic(resourceName, testLocation(), string(documentdb.BoundedStaleness), 2),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

//replication
func TestAccAzureRMCosmosDBAccount_geoReplicated_customId(t *testing.T) {
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_cosmosdb_account.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_geoReplicated_customId(ri, testLocation(), testAltLocation()),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(resourceName, testLocation(), string(documentdb.BoundedStaleness), 2),
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

func TestAccAzureRMCosmosDBAccount_geoReplicated_add_remove(t *testing.T) {
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_cosmosdb_account.test"
	configBasic := testAccAzureRMCosmosDBAccount_basic(ri, testLocation(), string(documentdb.BoundedStaleness), "", "")
	configReplicated := testAccAzureRMCosmosDBAccount_geoReplicated_customId(ri, testLocation(), testAltLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic,
				Check:  checkAccAzureRMCosmosDBAccount_basic(resourceName, testLocation(), string(documentdb.BoundedStaleness), 1),
			},
			{
				Config: configReplicated,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(resourceName, testLocation(), string(documentdb.BoundedStaleness), 2),
				),
			},
			{
				Config: configBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(resourceName, testLocation(), string(documentdb.BoundedStaleness), 1),
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

func TestAccAzureRMCosmosDBAccount_geoReplicated_rename(t *testing.T) {
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_cosmosdb_account.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_geoReplicated(ri, testLocation(), testAltLocation()),
				Check:  checkAccAzureRMCosmosDBAccount_basic(resourceName, testLocation(), string(documentdb.BoundedStaleness), 2),
			},
			{
				Config: testAccAzureRMCosmosDBAccount_geoReplicated_customId(ri, testLocation(), testAltLocation()),
				Check:  checkAccAzureRMCosmosDBAccount_basic(resourceName, testLocation(), string(documentdb.BoundedStaleness), 2),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMCosmosDBAccount_virtualNetworkFilter(t *testing.T) {
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_cosmosdb_account.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_virtualNetworkFilter(ri, testLocation()),
				Check:  checkAccAzureRMCosmosDBAccount_basic(resourceName, testLocation(), string(documentdb.BoundedStaleness), 1),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

//basic --> complete (
func TestAccAzureRMCosmosDBAccount_complete(t *testing.T) {
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_cosmosdb_account.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_basic(ri, testLocation(), string(documentdb.Session), "", ""),
				Check:  checkAccAzureRMCosmosDBAccount_basic(resourceName, testLocation(), string(documentdb.Session), 1),
			},
			{
				Config: testAccAzureRMCosmosDBAccount_complete(ri, testLocation(), testAltLocation()),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkAccAzureRMCosmosDBAccount_basic(resourceName, testLocation(), string(documentdb.BoundedStaleness), 2),
					resource.TestCheckResourceAttr(resourceName, "ip_range_filter", "104.42.195.92,40.76.54.131,52.176.6.30,52.169.50.45/32,52.187.184.26,10.20.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "enable_automatic_failover", "true"),
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

func TestAccAzureRMCosmosDBAccount_multiMaster(t *testing.T) {
	ri := tf.AccRandTimeInt()
	resourceName := "azurerm_cosmosdb_account.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_multiMaster(ri, testLocation(), testAltLocation()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "geo_location.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "write_endpoints.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "enable_multiple_write_locations", "true"),
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

func testCheckAzureRMCosmosDBAccountDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).cosmosDBClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

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

func testCheckAzureRMCosmosDBAccountExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for CosmosDB Account: '%s'", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).cosmosDBClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on cosmosDBClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: CosmosDB Account '%s' (resource group: '%s') does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testAccAzureRMCosmosDBAccount_basic(rInt int, location string, consistency string, consistencyOptions string, additional string) string {
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
`, rInt, location, rInt, consistency, consistencyOptions, additional)
}

func testAccAzureRMCosmosDBAccount_requiresImport(rInt int, location string, consistency string, consistencyOptions string, additional string) string {
	template := testAccAzureRMCosmosDBAccount_basic(rInt, location, consistency, consistencyOptions, additional)
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

func testAccAzureRMCosmosDBAccount_boundedStaleness_complete(rInt int, location string) string {
	return testAccAzureRMCosmosDBAccount_basic(rInt, location, string(documentdb.BoundedStaleness), `
        max_interval_in_seconds = 10
        max_staleness_prefix    = 200
`, "")
}

func testAccAzureRMCosmosDBAccount_mongoDB(rInt int, location string) string {
	return testAccAzureRMCosmosDBAccount_basic(rInt, location, string(documentdb.BoundedStaleness), "", `
        kind = "MongoDB"
    `)
}

func testAccAzureRMCosmosDBAccount_gremlin(rInt int, location string) string {
	return testAccAzureRMCosmosDBAccount_basic(rInt, location, string(documentdb.BoundedStaleness), "", `
        kind = "GlobalDocumentDB"

        capabilities = {
          name = "EnableGremlin"
        }
    `)
}

func testAccAzureRMCosmosDBAccount_table(rInt int, location string) string {
	return testAccAzureRMCosmosDBAccount_basic(rInt, location, string(documentdb.BoundedStaleness), "", `
        kind = "GlobalDocumentDB"

        capabilities = {
          name = "EnableTable"
        }
    `)
}

func testAccAzureRMCosmosDBAccount_geoReplicated(rInt int, location string, altLocation string) string {
	return testAccAzureRMCosmosDBAccount_basic(rInt, location, string(documentdb.BoundedStaleness), "", fmt.Sprintf(`
        geo_location {
            location          = "%s"
            failover_priority = 1
        }

    `, altLocation))
}

func testAccAzureRMCosmosDBAccount_multiMaster(rInt int, location string, altLocation string) string {
	return testAccAzureRMCosmosDBAccount_basic(rInt, location, string(documentdb.BoundedStaleness), "", fmt.Sprintf(`
        enable_multiple_write_locations = true

        geo_location {
            location          = "%s"
            failover_priority = 1
        }

    `, altLocation))
}

func testAccAzureRMCosmosDBAccount_geoReplicated_customId(rInt int, location string, altLocation string) string {
	return testAccAzureRMCosmosDBAccount_basic(rInt, location, string(documentdb.BoundedStaleness), "", fmt.Sprintf(`
        geo_location {
            prefix            = "acctest-%d-custom-id"
            location          = "%s"
            failover_priority = 1
        }

    `, rInt, altLocation))
}

func testAccAzureRMCosmosDBAccount_complete(rInt int, location string, altLocation string) string {
	return testAccAzureRMCosmosDBAccount_basic(rInt, location, string(documentdb.BoundedStaleness), "", fmt.Sprintf(`
		ip_range_filter				= "104.42.195.92,40.76.54.131,52.176.6.30,52.169.50.45/32,52.187.184.26,10.20.0.0/16"
		enable_automatic_failover	= true

        geo_location {
            prefix            = "acctest-%d-custom-id"
            location          = "%s"
            failover_priority = 1
        }
    `, rInt, altLocation))
}

func testAccAzureRMCosmosDBAccount_virtualNetworkFilter(rInt int, location string) string {
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
`, rInt)

	basic := testAccAzureRMCosmosDBAccount_basic(rInt, location, string(documentdb.BoundedStaleness), "", `
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

func checkAccAzureRMCosmosDBAccount_basic(resourceName string, location string, consistency string, locationCount int) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		testCheckAzureRMCosmosDBAccountExists(resourceName),
		resource.TestCheckResourceAttrSet(resourceName, "name"),
		resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
		resource.TestCheckResourceAttr(resourceName, "location", azureRMNormalizeLocation(location)),
		resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
		resource.TestCheckResourceAttr(resourceName, "offer_type", string(documentdb.Standard)),
		resource.TestCheckResourceAttr(resourceName, "consistency_policy.0.consistency_level", consistency),
		resource.TestCheckResourceAttr(resourceName, "geo_location.#", strconv.Itoa(locationCount)),
		resource.TestCheckResourceAttrSet(resourceName, "endpoint"),
		resource.TestCheckResourceAttr(resourceName, "read_endpoints.#", strconv.Itoa(locationCount)),
		resource.TestCheckResourceAttr(resourceName, "write_endpoints.#", "1"),
		resource.TestCheckResourceAttr(resourceName, "enable_multiple_write_locations", "false"),
		resource.TestCheckResourceAttrSet(resourceName, "primary_master_key"),
		resource.TestCheckResourceAttrSet(resourceName, "secondary_master_key"),
		resource.TestCheckResourceAttrSet(resourceName, "primary_readonly_master_key"),
		resource.TestCheckResourceAttrSet(resourceName, "secondary_readonly_master_key"),
	)
}

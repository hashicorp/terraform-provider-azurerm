package azurerm

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

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
		err = future.WaitForCompletion(ctx, client.Client)
		if err != nil {
			return err
		}
	}

	return nil
}

func TestAccAzureRMCosmosDBAccount_boundedStaleness(t *testing.T) {
	resourceName := "azurerm_cosmosdb_account.test"
	ri := acctest.RandInt()
	config := testAccAzureRMCosmosDBAccount_boundedStaleness(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCosmosDBAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "kind", "GlobalDocumentDB"),
				),
			},
		},
	})
}

func TestAccAzureRMCosmosDBAccount_boundedStalenessComplete(t *testing.T) {

	ri := acctest.RandInt()
	config := testAccAzureRMCosmosDBAccount_boundedStalenessComplete(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCosmosDBAccountExists("azurerm_cosmosdb_account.test"),
				),
			},
		},
	})
}

func TestAccAzureRMCosmosDBAccount_eventualConsistency(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMCosmosDBAccount_eventualConsistency(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCosmosDBAccountExists("azurerm_cosmosdb_account.test"),
				),
			},
		},
	})
}

func TestAccAzureRMCosmosDBAccount_mongoDB(t *testing.T) {
	resourceName := "azurerm_cosmosdb_account.test"
	ri := acctest.RandInt()
	config := testAccAzureRMCosmosDBAccount_mongoDB(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCosmosDBAccountExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "kind", "MongoDB"),
				),
			},
		},
	})
}

func TestAccAzureRMCosmosDBAccount_session(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMCosmosDBAccount_session(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCosmosDBAccountExists("azurerm_cosmosdb_account.test"),
				),
			},
		},
	})
}

func TestAccAzureRMCosmosDBAccount_strong(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMCosmosDBAccount_strong(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCosmosDBAccountExists("azurerm_cosmosdb_account.test"),
				),
			},
		},
	})
}

func TestAccAzureRMCosmosDBAccount_geoReplicated(t *testing.T) {

	ri := acctest.RandInt()
	config := testAccAzureRMCosmosDBAccount_geoReplicated(ri, testLocation(), testAltLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCosmosDBAccountExists("azurerm_cosmosdb_account.test"),
				),
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

func testAccAzureRMCosmosDBAccount_boundedStaleness(rInt int, location string) string {
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
    consistency_level = "BoundedStaleness"
  }

  failover_policy {
    location = "${azurerm_resource_group.test.location}"
    priority = 0
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMCosmosDBAccount_boundedStalenessComplete(rInt int, location string) string {
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
    consistency_level       = "BoundedStaleness"
    max_interval_in_seconds = 10
    max_staleness_prefix    = 200
  }

  failover_policy {
    location = "${azurerm_resource_group.test.location}"
    priority = 0
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMCosmosDBAccount_eventualConsistency(rInt int, location string) string {
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
    consistency_level = "Eventual"
  }

  failover_policy {
    location = "${azurerm_resource_group.test.location}"
    priority = 0
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMCosmosDBAccount_session(rInt int, location string) string {
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
    consistency_level = "Session"
  }

  failover_policy {
    location = "${azurerm_resource_group.test.location}"
    priority = 0
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMCosmosDBAccount_mongoDB(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cosmosdb_account" "test" {
  name                = "acctest-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  kind                = "MongoDB"
  offer_type          = "Standard"

  consistency_policy {
    consistency_level = "BoundedStaleness"
  }

  failover_policy {
    location = "${azurerm_resource_group.test.location}"
    priority = 0
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMCosmosDBAccount_strong(rInt int, location string) string {
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
    consistency_level = "Strong"
  }

  failover_policy {
    location = "${azurerm_resource_group.test.location}"
    priority = 0
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMCosmosDBAccount_geoReplicated(rInt int, location string, altLocation string) string {
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
    consistency_level       = "BoundedStaleness"
    max_interval_in_seconds = 10
    max_staleness_prefix    = 200
  }

  failover_policy {
    location = "${azurerm_resource_group.test.location}"
    priority = 0
  }

  failover_policy {
    location = "%s"
    priority = 1
  }
}
`, rInt, location, rInt, altLocation)
}

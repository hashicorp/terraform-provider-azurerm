package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccAzureRMCosmosDBAccount_failover_boundedStaleness(t *testing.T) {
	resourceName := "azurerm_cosmosdb_account.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMCosmosDBAccount_failover_boundedStaleness(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
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

func TestAccAzureRMCosmosDBAccount_failover_boundedStalenessComplete(t *testing.T) {

	ri := tf.AccRandTimeInt()
	config := testAccAzureRMCosmosDBAccount_failover_boundedStalenessComplete(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
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

func TestAccAzureRMCosmosDBAccount_failover_eventualConsistency(t *testing.T) {
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMCosmosDBAccount_failover_eventualConsistency(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
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

func TestAccAzureRMCosmosDBAccount_failover_mongoDB(t *testing.T) {
	resourceName := "azurerm_cosmosdb_account.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMCosmosDBAccount_failover_mongoDB(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
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

func TestAccAzureRMCosmosDBAccount_failover_session(t *testing.T) {
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMCosmosDBAccount_failover_session(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
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

func TestAccAzureRMCosmosDBAccount_failover_strong(t *testing.T) {
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMCosmosDBAccount_failover_strong(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
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

func TestAccAzureRMCosmosDBAccount_failover_geoReplicated(t *testing.T) {

	ri := tf.AccRandTimeInt()
	config := testAccAzureRMCosmosDBAccount_failover_geoReplicated(ri, testLocation(), testAltLocation())

	resource.ParallelTest(t, resource.TestCase{
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

func testAccAzureRMCosmosDBAccount_failover_boundedStaleness(rInt int, location string) string {
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

func testAccAzureRMCosmosDBAccount_failover_boundedStalenessComplete(rInt int, location string) string {
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

func testAccAzureRMCosmosDBAccount_failover_eventualConsistency(rInt int, location string) string {
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

func testAccAzureRMCosmosDBAccount_failover_session(rInt int, location string) string {
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

func testAccAzureRMCosmosDBAccount_failover_mongoDB(rInt int, location string) string {
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

func testAccAzureRMCosmosDBAccount_failover_strong(rInt int, location string) string {
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

func testAccAzureRMCosmosDBAccount_failover_geoReplicated(rInt int, location string, altLocation string) string {
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
    max_interval_in_seconds = 333
    max_staleness_prefix    = 101101
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

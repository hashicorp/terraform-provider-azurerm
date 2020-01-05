package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccAzureRMCosmosDBAccount_failover_boundedStaleness(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_failover_boundedStaleness(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCosmosDBAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "kind", "GlobalDocumentDB"),
				),
			},
		},
	})
}

func TestAccAzureRMCosmosDBAccount_failover_boundedStalenessComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_failover_boundedStalenessComplete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCosmosDBAccountExists("azurerm_cosmosdb_account.test"),
				),
			},
		},
	})
}

func TestAccAzureRMCosmosDBAccount_failover_eventualConsistency(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_failover_eventualConsistency(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCosmosDBAccountExists("azurerm_cosmosdb_account.test"),
				),
			},
		},
	})
}

func TestAccAzureRMCosmosDBAccount_failover_mongoDB(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_failover_mongoDB(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCosmosDBAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "kind", "MongoDB"),
				),
			},
		},
	})
}

func TestAccAzureRMCosmosDBAccount_failover_session(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_failover_session(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCosmosDBAccountExists("azurerm_cosmosdb_account.test"),
				),
			},
		},
	})
}

func TestAccAzureRMCosmosDBAccount_failover_strong(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_failover_strong(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCosmosDBAccountExists("azurerm_cosmosdb_account.test"),
				),
			},
		},
	})
}

func TestAccAzureRMCosmosDBAccount_failover_geoReplicated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cosmosdb_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCosmosDBAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMCosmosDBAccount_failover_geoReplicated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCosmosDBAccountExists("azurerm_cosmosdb_account.test"),
				),
			},
		},
	})
}

func testAccAzureRMCosmosDBAccount_failover_boundedStaleness(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMCosmosDBAccount_failover_boundedStalenessComplete(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMCosmosDBAccount_failover_eventualConsistency(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMCosmosDBAccount_failover_session(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMCosmosDBAccount_failover_mongoDB(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMCosmosDBAccount_failover_strong(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMCosmosDBAccount_failover_geoReplicated(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Secondary)
}

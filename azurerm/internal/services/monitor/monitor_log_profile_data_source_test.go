package monitor_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type MonitorLogProfileDataSource struct {
}

// These tests are actually run as part of the resoure ones due to
// Azure only being happy about provisioning one per subscription at once
// (which our test suite can't easily workaround)

func testAccDataSourceMonitorLogProfile_storageaccount(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_monitor_log_profile", "test")
	r := MonitorLogProfileDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.storageaccountConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("categories.#").Exists(),
				check.That(data.ResourceName).Key("locations.#").Exists(),
				check.That(data.ResourceName).Key("storage_account_id").Exists(),
				check.That(data.ResourceName).Key("servicebus_rule_id").HasValue(""),
				check.That(data.ResourceName).Key("retention_policy.#").HasValue("1"),
				check.That(data.ResourceName).Key("retention_policy.0.enabled").Exists(),
				check.That(data.ResourceName).Key("retention_policy.0.days").Exists(),
			),
		},
	})
}

func testAccDataSourceMonitorLogProfile_eventhub(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_monitor_log_profile", "test")
	r := MonitorLogProfileDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.eventhubConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("categories.#").Exists(),
				check.That(data.ResourceName).Key("locations.#").Exists(),
				check.That(data.ResourceName).Key("storage_account_id").HasValue(""),
				check.That(data.ResourceName).Key("servicebus_rule_id").Exists(),
				check.That(data.ResourceName).Key("retention_policy.#").HasValue("1"),
				check.That(data.ResourceName).Key("retention_policy.0.enabled").Exists(),
				check.That(data.ResourceName).Key("retention_policy.0.days").Exists(),
			),
		},
	})
}

func (MonitorLogProfileDataSource) storageaccountConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_monitor_log_profile" "test" {
  name = "acctestlp-%d"

  categories = [
    "Action",
  ]

  locations = [
    "%s",
  ]

  storage_account_id = azurerm_storage_account.test.id

  retention_policy {
    enabled = true
    days    = 7
  }
}

data "azurerm_monitor_log_profile" "test" {
  name = azurerm_monitor_log_profile.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.Locations.Primary)
}

func (MonitorLogProfileDataSource) eventhubConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctestehns-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  capacity            = 2
}

resource "azurerm_monitor_log_profile" "test" {
  name = "acctestlp-%d"

  categories = [
    "Action",
  ]

  locations = [
    "%s",
  ]

  # RootManageSharedAccessKey is created by default with listen, send, manage permissions
  servicebus_rule_id = "${azurerm_eventhub_namespace.test.id}/authorizationrules/RootManageSharedAccessKey"

  retention_policy {
    enabled = true
    days    = 7
  }
}

data "azurerm_monitor_log_profile" "test" {
  name = azurerm_monitor_log_profile.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.Locations.Primary)
}

package monitor_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

// These tests are actually run as part of the resoure ones due to
// Azure only being happy about provisioning one per subscription at once
// (which our test suite can't easily workaround)

func testAccDataSourceAzureRMMonitorLogProfile_storageaccount(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_monitor_log_profile", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMMonitorLogProfile_storageaccountConfig(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "categories.#"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "locations.#"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "servicebus_rule_id", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "retention_policy.0.enabled"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "retention_policy.0.days"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMMonitorLogProfile_eventhub(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_monitor_log_profile", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMMonitorLogProfile_eventhubConfig(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "categories.#"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "locations.#"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_account_id", ""),
					resource.TestCheckResourceAttrSet(data.ResourceName, "servicebus_rule_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_policy.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "retention_policy.0.enabled"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "retention_policy.0.days"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMMonitorLogProfile_storageaccountConfig(data acceptance.TestData) string {
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

func testAccDataSourceAzureRMMonitorLogProfile_eventhubConfig(data acceptance.TestData) string {
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

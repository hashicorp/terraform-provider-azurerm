package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

// These tests are actually run as part of the resoure ones due to
// Azure only being happy about provisioning one per subscription at once
// (which our test suite can't easily workaround)

func testAccDataSourceAzureRMMonitorLogProfile_storageaccount(t *testing.T) {
	dataSourceName := "data.azurerm_monitor_log_profile.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMMonitorLogProfile_storageaccountConfig(ri, rs, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "categories.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "locations.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(dataSourceName, "servicebus_rule_id", ""),
					resource.TestCheckResourceAttr(dataSourceName, "retention_policy.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "retention_policy.0.enabled"),
					resource.TestCheckResourceAttrSet(dataSourceName, "retention_policy.0.days"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMMonitorLogProfile_eventhub(t *testing.T) {
	dataSourceName := "data.azurerm_monitor_log_profile.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMMonitorLogProfile_eventhubConfig(ri, rs, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "categories.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "locations.#"),
					resource.TestCheckResourceAttr(dataSourceName, "storage_account_id", ""),
					resource.TestCheckResourceAttrSet(dataSourceName, "servicebus_rule_id"),
					resource.TestCheckResourceAttr(dataSourceName, "retention_policy.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "retention_policy.0.enabled"),
					resource.TestCheckResourceAttrSet(dataSourceName, "retention_policy.0.days"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMMonitorLogProfile_storageaccountConfig(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
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

  storage_account_id = "${azurerm_storage_account.test.id}"

  retention_policy {
    enabled = true
    days    = 7
  }
}

data "azurerm_monitor_log_profile" "test" {
  name = "${azurerm_monitor_log_profile.test.name}"
}
`, rInt, location, rString, rInt, location)
}

func testAccDataSourceAzureRMMonitorLogProfile_eventhubConfig(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctestehns-%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
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
  name = "${azurerm_monitor_log_profile.test.name}"
}
`, rInt, location, rString, rInt, location)
}

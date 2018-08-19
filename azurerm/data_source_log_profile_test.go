package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMLogProfile_storageaccount(t *testing.T) {
	dataSourceName := "data.azurerm_log_profile.test"
	ri := acctest.RandInt()
	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMLogProfile_storageaccount(ri, rs, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "categories.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "locations.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(dataSourceName, "service_bus_rule_id", ""),
					resource.TestCheckResourceAttr(dataSourceName, "retention_policy.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "retention_policy.0.enabled"),
					resource.TestCheckResourceAttrSet(dataSourceName, "retention_policy.0.days"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMLogProfile_eventhub(t *testing.T) {
	dataSourceName := "data.azurerm_log_profile.test"
	ri := acctest.RandInt()
	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMLogProfile_eventhub(ri, rs, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "categories.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "locations.#"),
					resource.TestCheckResourceAttr(dataSourceName, "storage_account_id", ""),
					resource.TestCheckResourceAttrSet(dataSourceName, "service_bus_rule_id"),
					resource.TestCheckResourceAttr(dataSourceName, "retention_policy.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "retention_policy.0.enabled"),
					resource.TestCheckResourceAttrSet(dataSourceName, "retention_policy.0.days"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMLogProfile_storageaccount(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
		resource "azurerm_resource_group" "test" {
			name     = "acctest%d-rg"
			location = "%s"
		}
		
		resource "azurerm_storage_account" "test" {
			name                     = "%s"
			resource_group_name      = "${azurerm_resource_group.test.name}"
			location                 = "${azurerm_resource_group.test.location}"
			account_tier             = "Standard"
			account_replication_type = "GRS"
		}
			
		resource "azurerm_log_profile" "test" {
			name = "storageaccounttest-logprofile"
		
			categories = [
				"Action",
			]
			
			locations = [
				"%s"
			]
			
			storage_account_id = "${azurerm_storage_account.test.id}"
		
			retention_policy {
				enabled = true
				days    = 7
			}
		}

		data "azurerm_log_profile" "test" {
			name = "${azurerm_log_profile.test.name}"
		}
	`, rInt, location, rString, location)
}

func testAccDataSourceAzureRMLogProfile_eventhub(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
		resource "azurerm_resource_group" "test" {
			name     = "acctest%d-rg"
			location = "%s"
		}
		
		resource "azurerm_eventhub_namespace" "test" {
			name                = "a%s"
			location            = "${azurerm_resource_group.test.location}"
			resource_group_name = "${azurerm_resource_group.test.name}"
			sku                 = "Standard"
			capacity            = 2
		}
			
		resource "azurerm_log_profile" "test" {
			name = "eventhubtest-logprofile"
		
			categories = [
				"Action",
			]
			
			locations = [
				"%s"
			]
			
			# RootManageSharedAccessKey is created by default with listen, send, manage permissions
			service_bus_rule_id = "${azurerm_eventhub_namespace.test.id}/authorizationrules/RootManageSharedAccessKey"
		
			retention_policy {
				enabled = true
				days    = 7
			}
		}

		data "azurerm_log_profile" "test" {
			name = "${azurerm_log_profile.test.name}"
		}
	`, rInt, location, rString, location)
}

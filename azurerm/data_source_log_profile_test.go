package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMLogProfile(t *testing.T) {
	// NOTE: this is a combined test rather than separate split out tests due to
	// Azure only being happy about provisioning one per subscription at once
	// (which our test suite can't easily workaround)
	testCases := map[string]map[string]func(t *testing.T){
		"basic": {
			"eventhub":       testAccDataSourceAzureRMLogProfile_eventhub,
			"storageaccount": testAccDataSourceAzureRMLogProfile_storageaccount,
		},
	}

	for group, m := range testCases {
		m := m
		t.Run(group, func(t *testing.T) {
			for name, tc := range m {
				tc := tc
				t.Run(name, func(t *testing.T) {
					tc(t)
				})
			}
		})
	}
}

func testAccDataSourceAzureRMLogProfile_storageaccount(t *testing.T) {
	dataSourceName := "data.azurerm_log_profile.test"
	ri := acctest.RandInt()
	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMLogProfile_storageaccountConfig(ri, rs, testLocation()),
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

func testAccDataSourceAzureRMLogProfile_eventhub(t *testing.T) {
	dataSourceName := "data.azurerm_log_profile.test"
	ri := acctest.RandInt()
	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLogProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMLogProfile_eventhubConfig(ri, rs, testLocation()),
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

func testAccDataSourceAzureRMLogProfile_storageaccountConfig(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
		resource "azurerm_resource_group" "test" {
			name     = "acctestrg-%d"
			location = "%s"
		}
		
		resource "azurerm_storage_account" "test" {
			name                     = "acctestsa%s"
			resource_group_name      = "${azurerm_resource_group.test.name}"
			location                 = "${azurerm_resource_group.test.location}"
			account_tier             = "Standard"
			account_replication_type = "GRS"
		}
			
		resource "azurerm_log_profile" "test" {
			name = "acctestlp-%d"
		
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
	`, rInt, location, rString, rInt, location)
}

func testAccDataSourceAzureRMLogProfile_eventhubConfig(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
		resource "azurerm_resource_group" "test" {
			name     = "acctestrg-%d"
			location = "%s"
		}
		
		resource "azurerm_eventhub_namespace" "test" {
			name                = "acctestehns-%s"
			location            = "${azurerm_resource_group.test.location}"
			resource_group_name = "${azurerm_resource_group.test.name}"
			sku                 = "Standard"
			capacity            = 2
		}
			
		resource "azurerm_log_profile" "test" {
			name = "acctestlp-%d"
		
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
	`, rInt, location, rString, rInt, location)
}

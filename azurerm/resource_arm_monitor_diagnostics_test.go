package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMMonitorDiagnostics_basic(t *testing.T) {
	resourceName := "azurerm_monitor_diagnostics.test"
	ri := acctest.RandIntRange(10000, 99999)
	objectName := fmt.Sprintf("acctest%d", ri)
	location := testLocation()

	config := testAccAzureRMMonitorDiagnostics_basic(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorDiagnosticsExists(resourceName, objectName),
					resource.TestCheckResourceAttr(resourceName, "disabled_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "disabled_settings.0", "AuditEvent"),
				),
			},
		},
	})
}
func TestAccAzureRMMonitorDiagnostics_complete(t *testing.T) {
	resourceName := "azurerm_monitor_diagnostics.test"
	ri := acctest.RandIntRange(10000, 99999)
	objectName := fmt.Sprintf("acctest%d", ri)
	location := testLocation()

	config := testAccAzureRMMonitorDiagnostics_complete(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorDiagnosticsExists(resourceName, objectName),
					resource.TestCheckResourceAttr(resourceName, "disabled_settings.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "retention_days", "90"),
				),
			},
		},
	})
}

func TestAccAzureRMMonitorDiagnostics_update(t *testing.T) {
	resourceName := "azurerm_monitor_diagnostics.test"
	ri := acctest.RandIntRange(10000, 99999)
	objectName := fmt.Sprintf("acctest%d", ri)
	location := testLocation()

	configBasic := testAccAzureRMMonitorDiagnostics_basic(ri, location)
	configUpdate := testAccAzureRMMonitorDiagnostics_update(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: configBasic,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorDiagnosticsExists(resourceName, objectName),
					resource.TestCheckResourceAttr(resourceName, "disabled_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "disabled_settings.0", "AuditEvent"),
				),
			},
			{
				Config:             configUpdate,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMMonitorDiagnosticsExists(name, objectName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		stateName := rs.Primary.Attributes["name"]
		if stateName != objectName {
			return fmt.Errorf("State inconsistent, %s does not match state name %s", stateName, objectName)
		}

		return nil
	}
}

func testAccAzureRMMonitorDiagnostics_basic(randomInt int, location string) string {
	return fmt.Sprintf(`
	data "azurerm_client_config" "current" {}

	resource "azurerm_resource_group" "test_rg" {
		name 		= "acctest%d"
		location 	= "%s"
	}

	resource "azurerm_storage_account" "test_storage_logs" {
	name 						 = "acctestlogs%d"
		resource_group_name 	 = "${azurerm_resource_group.test_rg.name}"
		location 				 = "${azurerm_resource_group.test_rg.location}"
		account_replication_type = "LRS"
		account_tier 			 = "Standard"
	}

	resource "azurerm_key_vault" "test_vault" {
		name                = "vault%d"
		location            = "${azurerm_resource_group.test_rg.location}"
		resource_group_name = "${azurerm_resource_group.test_rg.name}"
		tenant_id           = "${data.azurerm_client_config.current.tenant_id}"
	
		sku {
			name = "standard"
		}
	}
		
	resource "azurerm_monitor_diagnostics" "test" {
		name 			   = "acctest%d"
		target_resource_id = "${azurerm_key_vault.test_vault.id}"
		storage_account_id = "${azurerm_storage_account.test_storage_logs.id}"
		disabled_settings  = ["AuditEvent"]
	}`, randomInt, location, randomInt, randomInt, randomInt)
}

func testAccAzureRMMonitorDiagnostics_complete(randomInt int, location string) string {
	return fmt.Sprintf(`
	data "azurerm_client_config" "current" {}

	resource "azurerm_resource_group" "test_rg" {
		name = "acctest%d"
		location = "%s"
	}

	resource "azurerm_log_analytics_workspace" "analytics" {
		name                = "acctest%d"
		location            = "${azurerm_resource_group.test_rg.location}"
		resource_group_name = "${azurerm_resource_group.test_rg.name}"
		sku                 = "Standard"
		retention_in_days   = 30
	}

	resource "azurerm_storage_account" "test_storage_logs" {
		name 					 = "acctestlogs%d"
		resource_group_name 	 = "${azurerm_resource_group.test_rg.name}"
		location 				 = "${azurerm_resource_group.test_rg.location}"
		account_replication_type = "LRS"
		account_tier 			 = "Standard"
	}

	resource "azurerm_key_vault" "test_vault" {
		name                = "vault%d"
		location            = "${azurerm_resource_group.test_rg.location}"
		resource_group_name = "${azurerm_resource_group.test_rg.name}"
		tenant_id           = "${data.azurerm_client_config.current.tenant_id}"
	
		sku {
			name = "standard"
		}
	}

	resource "azurerm_monitor_diagnostics" "test" {
		name 			   = "acctest%d"
		target_resource_id = "${azurerm_key_vault.test_vault.id}"
		storage_account_id = "${azurerm_storage_account.test_storage_logs.id}"
		workspace_id 	   = "${azurerm_log_analytics_workspace.analytics.id}"
		retention_days     = 90
	}`, randomInt, location, randomInt, randomInt, randomInt, randomInt)
}

func testAccAzureRMMonitorDiagnostics_update(randomInt int, location string) string {
	return fmt.Sprintf(`
	data "azurerm_client_config" "current" {}

	resource "azurerm_resource_group" "test_rg" {
		name 		= "acctest%d"
		location 	= "%s"
	}

	resource "azurerm_storage_account" "test_storage_logs" {
		name 					 = "acctestlogs%d"
		resource_group_name 	 = "${azurerm_resource_group.test_rg.name}"
		location 				 = "${azurerm_resource_group.test_rg.location}"
		account_replication_type = "LRS"
		account_tier 			 = "Standard"
	}

	resource "azurerm_key_vault" "test_vault" {
		name                = "vault%d"
		location            = "${azurerm_resource_group.test_rg.location}"
		resource_group_name = "${azurerm_resource_group.test_rg.name}"
		tenant_id           = "${data.azurerm_client_config.current.tenant_id}"
	
		sku {
			name = "standard"
		}
	}
		
	resource "azurerm_monitor_diagnostics" "test" {
		name 			   = "acctest%d"
		target_resource_id = "${azurerm_key_vault.test_vault.id}"
		storage_account_id = "${azurerm_storage_account.test_storage_logs.id}"
		retention_days 	   = 30
	}`, randomInt, location, randomInt, randomInt, randomInt)
}

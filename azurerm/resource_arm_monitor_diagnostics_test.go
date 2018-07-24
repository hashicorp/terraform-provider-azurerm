package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMMonitorDiagnostics_basic(t *testing.T) {
	resourceName := "azurerm_monitor_diagnostics.test"
	ri := acctest.RandIntRange(10000, 99999)
	objectName := fmt.Sprintf("acctest%d", ri)
	location := testLocation()

	config := testAccAzureRMMonitorDiagnostics_basic(ri, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorDiagnosticsDestroy,
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
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorDiagnosticsDestroy,
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
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorDiagnosticsDestroy,
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
				Config: configUpdate,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorDiagnosticsExists(resourceName, objectName),
					resource.TestCheckResourceAttr(resourceName, "disabled_settings.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "retention_days", "30"),
				),
			},
		},
	})
}

func testCheckAzureRMMonitorDiagnosticsDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).monitorDiagnosticSettingsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_monitor_diagnostics" {
			continue
		}

		diagSettingName := rs.Primary.Attributes["name"]
		targetResourceId := rs.Primary.Attributes["target_resource_id"]
		resp, err := conn.Get(ctx, targetResourceId, diagSettingName)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return fmt.Errorf("Error while retrieving DiagnosticSettings %v", err)
		}
		return fmt.Errorf("DiagnosticSettings %s still exists", diagSettingName)
	}
	return nil
}

func testCheckAzureRMMonitorDiagnosticsExists(name, objectName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Did not find resource %s in state", name)
		}

		stateName := rs.Primary.Attributes["name"]
		if stateName != objectName {
			return fmt.Errorf("State inconsistent, %s does not match state name %s", stateName, objectName)
		}

		conn := testAccProvider.Meta().(*ArmClient).monitorDiagnosticSettingsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		targetResourceId := rs.Primary.Attributes["target_resource_id"]

		resp, err := conn.Get(ctx, targetResourceId, objectName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("DiagnosticSettings %s did not exist for resource %s", objectName, targetResourceId)
			}
			return fmt.Errorf("Error while retrieving DiagnosticSettings %v", err)
		}

		return nil
	}
}

func testAccAzureRMMonitorDiagnostics_basic(rInt int, location string) string {
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
	}`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMMonitorDiagnostics_complete(rInt int, location string) string {
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
	}`, rInt, location, rInt, rInt, rInt, rInt)
}

func testAccAzureRMMonitorDiagnostics_update(rInt int, location string) string {
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
	}`, rInt, location, rInt, rInt, rInt)
}

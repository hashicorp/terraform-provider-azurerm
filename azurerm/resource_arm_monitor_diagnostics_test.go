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
				),
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
		name = "acctest%d"
		location = "%s"
	}

	resource "azurerm_key_vault" "test_vault" {
		name                = "vault%d"
		location            = "${azurerm_resource_group.test_rg.location}"
		resource_group_name = "${azurerm_resource_group.test_rg.name}"
		tenant_id           = "${data.azurerm_client_config.current.tenant_id}"
	
		sku {
		  name = "premium"
		}
	}

	resource "azurerm_storage_account" "test_storage_logs" {
		name = "acctestlogs%d"
		resource_group_name = "${azurerm_resource_group.test_rg.name}"
		location = "${azurerm_resource_group.test_rg.location}"
		account_replication_type = "LRS"
		account_tier = "Standard"
	}
		
	resource "azurerm_monitor_diagnostics" "test" {
		name = "acctest%d"
		resource_id = "${azurerm_key_vault.test_vault.id}"
		storage_account_id = "${azurerm_storage_account.test_storage_logs.id}"
		metric_settings = {
			category = "AllMetrics"
			retention_days = 2
		}
		log_settings = {
			category = "AuditEvent"
			retention_days = 2
		}
	}`, randomInt, location, randomInt, randomInt, randomInt)
}

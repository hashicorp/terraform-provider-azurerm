package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccAzureRMRecoveryProtectionContainerMapping_basic(t *testing.T) {
	resourceGroupName := "azurerm_resource_group.test1"
	vaultName := "azurerm_recovery_services_vault.test"
	fabricName := "azurerm_recovery_services_fabric.test1"
	protectionContainerName := "azurerm_recovery_services_protection_container.test1"
	resourceName := "azurerm_recovery_services_protection_container_mapping.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMResourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRecoveryProtectionContainerMapping_basic(ri, testLocation(), testAltLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRecoveryProtectionContainerMappingExists(resourceGroupName, vaultName, fabricName, protectionContainerName, resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAzureRMRecoveryProtectionContainerMapping_basic(rInt int, location string, altLocation string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test1" {
  name     = "acctestRG1-%d"
  location = "%s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-vault-%d"
  location            = "${azurerm_resource_group.test1.location}"
  resource_group_name = "${azurerm_resource_group.test1.name}"
  sku                 = "Standard"
}

resource "azurerm_recovery_services_fabric" "test1" {
  resource_group_name          = "${azurerm_resource_group.test1.name}"
  recovery_vault_name          = "${azurerm_recovery_services_vault.test.name}"
  name                         = "acctest-fabric1-%d"
  location                     = "${azurerm_resource_group.test1.location}"
}

resource "azurerm_recovery_services_fabric" "test2" {
  resource_group_name          = "${azurerm_resource_group.test1.name}"
  recovery_vault_name          = "${azurerm_recovery_services_vault.test.name}"
  name                         = "acctest-fabric2-%d"
  location                     = "%s"
  depends_on                   = ["azurerm_recovery_services_fabric.test1"]
}

resource "azurerm_recovery_services_protection_container" "test1" {
  resource_group_name           = "${azurerm_resource_group.test1.name}"
  recovery_vault_name           = "${azurerm_recovery_services_vault.test.name}"
  recovery_fabric_name          = "${azurerm_recovery_services_fabric.test1.name}"
  name                          = "acctest-protection-cont1-%d"
}

resource "azurerm_recovery_services_protection_container" "test2" {
  resource_group_name           = "${azurerm_resource_group.test1.name}"
  recovery_vault_name           = "${azurerm_recovery_services_vault.test.name}"
  recovery_fabric_name          = "${azurerm_recovery_services_fabric.test2.name}"
  name                          = "acctest-protection-cont2-%d"
}

resource "azurerm_recovery_services_replication_policy" "test" {
  resource_group_name           = "${azurerm_resource_group.test1.name}"
  recovery_vault_name           = "${azurerm_recovery_services_vault.test.name}"
  name                          = "acctest-policy-%d"
  recovery_point_retention_in_minutes = "${24 * 60}"
  application_consistent_snapshot_frequency_in_minutes = "${4 * 60}"
}

resource "azurerm_recovery_services_protection_container_mapping" "test" {
  resource_group_name            = "${azurerm_resource_group.test1.name}"
  recovery_vault_name            = "${azurerm_recovery_services_vault.test.name}"
  recovery_fabric_name           = "${azurerm_recovery_services_fabric.test1.name}"
  recovery_source_protection_container_name = "${azurerm_recovery_services_protection_container.test1.name}"
  recovery_target_protection_container_id = "${azurerm_recovery_services_protection_container.test2.id}"
  recovery_replication_policy_id = "${azurerm_recovery_services_replication_policy.test.id}"
  name                           = "mapping-%d"
}
`, rInt, location, rInt, rInt, rInt, altLocation, rInt, rInt, rInt, rInt)
}

func testCheckAzureRMRecoveryProtectionContainerMappingExists(resourceGroupStateName, vaultStateName string, resourceStateName string, protectionContainerStateName string, protectionContainerStateMappingName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		resourceGroupState, ok := s.RootModule().Resources[resourceGroupStateName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceGroupStateName)
		}
		vaultState, ok := s.RootModule().Resources[vaultStateName]
		if !ok {
			return fmt.Errorf("Not found: %s", vaultStateName)
		}
		fabricState, ok := s.RootModule().Resources[resourceStateName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceStateName)
		}
		protectionContainerState, ok := s.RootModule().Resources[protectionContainerStateName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceStateName)
		}
		protectionContainerMappingState, ok := s.RootModule().Resources[protectionContainerStateMappingName]
		if !ok {
			return fmt.Errorf("Not found: %s", protectionContainerStateMappingName)
		}

		resourceGroupName := resourceGroupState.Primary.Attributes["name"]
		vaultName := vaultState.Primary.Attributes["name"]
		fabricName := fabricState.Primary.Attributes["name"]
		protectionContainerName := protectionContainerState.Primary.Attributes["name"]
		mappingName := protectionContainerMappingState.Primary.Attributes["name"]

		// Ensure mapping exists in API
		client := testAccProvider.Meta().(*ArmClient).recoveryServices.ContainerMappingClient(resourceGroupName, vaultName)
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, fabricName, protectionContainerName, mappingName)
		if err != nil {
			return fmt.Errorf("Bad: Get on fabricClient: %+v", err)
		}

		if resp.Response.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: fabric: %q does not exist", fabricName)
		}

		return nil
	}
}

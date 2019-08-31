package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccAzureRMRecoveryNetworkMapping_basic(t *testing.T) {
	resourceGroupName := "azurerm_resource_group.test"
	vaultName := "azurerm_recovery_services_vault.test"
	fabricName := "azurerm_recovery_services_fabric.test1"
	networkName := "azurerm_virtual_network.test1"
	resourceName := "azurerm_recovery_network_mapping.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMResourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRecoveryNetworkMapping_basic(ri, testLocation(), testAltLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRecoveryNetworkMappingExists(resourceGroupName, vaultName, fabricName, networkName, resourceName),
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

func testAccAzureRMRecoveryNetworkMapping_basic(rInt int, location string, altLocation string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name                        = "acctestRG1-%d"
  location                    = "%s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                        = "acctest-vault-%d"
  location                    = "${azurerm_resource_group.test.location}"
  resource_group_name         = "${azurerm_resource_group.test.name}"
  sku                         = "Standard"
}

resource "azurerm_recovery_services_fabric" "test1" {
  resource_group_name         = "${azurerm_resource_group.test.name}"
  recovery_vault_name         = "${azurerm_recovery_services_vault.test.name}"
  name                        = "acctest-fabric1-%d"
  location                    = "${azurerm_resource_group.test.location}"
}

resource "azurerm_recovery_services_fabric" "test2" {
  resource_group_name         = "${azurerm_resource_group.test.name}"
  recovery_vault_name         = "${azurerm_recovery_services_vault.test.name}"
  name                        = "acctest-fabric2-%d"
  location                    = "%s"
  depends_on                  = ["azurerm_recovery_services_fabric.test1"]
}

resource "azurerm_virtual_network" "test1" {
  name                        = "network1-%d"
  resource_group_name         = "${azurerm_resource_group.test.name}"
  address_space               = [ "192.168.1.0/24" ]
  location                    = "${azurerm_recovery_services_fabric.test1.location}"
}

resource "azurerm_virtual_network" "test2" {
  name                        = "network2-%d"
  resource_group_name         = "${azurerm_resource_group.test.name}"
  address_space               = [ "192.168.2.0/24" ]
  location                    = "${azurerm_recovery_services_fabric.test2.location}"
}

resource "azurerm_recovery_network_mapping" "test" {
  resource_group_name         = "${azurerm_resource_group.test.name}"
  recovery_vault_name         = "${azurerm_recovery_services_vault.test.name}"
  name                        = "mapping-%d"
  source_recovery_fabric_name = "${azurerm_recovery_services_fabric.test1.name}"
  target_recovery_fabric_name = "${azurerm_recovery_services_fabric.test2.name}"
  source_network_id           = "${azurerm_virtual_network.test1.id}"
  target_network_id           = "${azurerm_virtual_network.test2.id}"
}
`, rInt, location, rInt, rInt, rInt, altLocation, rInt, rInt, rInt)
}

func testCheckAzureRMRecoveryNetworkMappingExists(resourceGroupStateName, vaultStateName string, fabricStateName string, networkStateName string, networkStateMappingName string) resource.TestCheckFunc {
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
		fabricState, ok := s.RootModule().Resources[fabricStateName]
		if !ok {
			return fmt.Errorf("Not found: %s", fabricStateName)
		}
		networkState, ok := s.RootModule().Resources[networkStateName]
		if !ok {
			return fmt.Errorf("Not found: %s", fabricStateName)
		}
		networkMappingState, ok := s.RootModule().Resources[networkStateMappingName]
		if !ok {
			return fmt.Errorf("Not found: %s", networkStateMappingName)
		}

		resourceGroupName := resourceGroupState.Primary.Attributes["name"]
		vaultName := vaultState.Primary.Attributes["name"]
		fabricName := fabricState.Primary.Attributes["name"]
		networkName := networkState.Primary.Attributes["name"]
		mappingName := networkMappingState.Primary.Attributes["name"]

		// Ensure mapping exists in API
		client := testAccProvider.Meta().(*ArmClient).recoveryServices.NetworkMappingClient(resourceGroupName, vaultName)
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, fabricName, networkName, mappingName)
		if err != nil {
			return fmt.Errorf("Bad: Get on networkMappingClient: %+v", err)
		}

		if resp.Response.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: networkMapping: %q does not exist", mappingName)
		}

		return nil
	}
}

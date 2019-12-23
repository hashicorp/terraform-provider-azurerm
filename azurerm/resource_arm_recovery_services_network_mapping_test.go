package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMRecoveryNetworkMapping_basic(t *testing.T) {
	resourceName := "azurerm_recovery_network_mapping.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRecoveryNetworkMappingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRecoveryNetworkMapping_basic(ri, acceptance.Location(), acceptance.AltLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRecoveryNetworkMappingExists(resourceName),
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
  name     = "acctestRG-recovery1-%d"
  location = "%s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-vault-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
}

resource "azurerm_recovery_services_fabric" "test1" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.test.name}"
  name                = "acctest-fabric1-%d"
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_recovery_services_fabric" "test2" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.test.name}"
  name                = "acctest-fabric2-%d"
  location            = "%s"
  depends_on          = ["azurerm_recovery_services_fabric.test1"]
}

resource "azurerm_virtual_network" "test1" {
  name                = "network1-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["192.168.1.0/24"]
  location            = "${azurerm_recovery_services_fabric.test1.location}"
}

resource "azurerm_virtual_network" "test2" {
  name                = "network2-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["192.168.2.0/24"]
  location            = "${azurerm_recovery_services_fabric.test2.location}"
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

func testCheckAzureRMRecoveryNetworkMappingExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		state, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroupName := state.Primary.Attributes["resource_group_name"]
		vaultName := state.Primary.Attributes["recovery_vault_name"]
		fabricName := state.Primary.Attributes["source_recovery_fabric_name"]
		networkId := state.Primary.Attributes["source_network_id"]
		mappingName := state.Primary.Attributes["name"]

		id, err := azure.ParseAzureResourceID(networkId)
		if err != nil {
			return err
		}
		networkName := id.Path["virtualNetworks"]

		client := acceptance.AzureProvider.Meta().(*clients.Client).RecoveryServices.NetworkMappingClient(resourceGroupName, vaultName)
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// TODO Fix Bad: networkMapping error
		resp, err := client.Get(ctx, fabricName, networkName, mappingName)
		if err != nil {
			if resp.Response.StatusCode == http.StatusNotFound {
				return fmt.Errorf("Bad: networkMapping: %q (network %q) does not exist", mappingName, networkName)
			}

			return fmt.Errorf("Bad: Get on networkMappingClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMRecoveryNetworkMappingDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_recovery_network_mapping" {
			continue
		}

		resourceGroupName := rs.Primary.Attributes["resource_group_name"]
		vaultName := rs.Primary.Attributes["recovery_vault_name"]
		fabricName := rs.Primary.Attributes["source_recovery_fabric_name"]
		networkId := rs.Primary.Attributes["source_network_id"]
		mappingName := rs.Primary.Attributes["name"]

		id, err := azure.ParseAzureResourceID(networkId)
		if err != nil {
			return err
		}
		networkName := id.Path["virtualNetworks"]

		client := acceptance.AzureProvider.Meta().(*clients.Client).RecoveryServices.NetworkMappingClient(resourceGroupName, vaultName)
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := client.Get(ctx, fabricName, networkName, mappingName)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Network Mapping still exists:\n%#v", resp.Properties)
		}
	}

	return nil
}

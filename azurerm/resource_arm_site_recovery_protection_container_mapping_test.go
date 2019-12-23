package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMSiteRecoveryProtectionContainerMapping_basic(t *testing.T) {
	resourceName := "azurerm_site_recovery_protection_container_mapping.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSiteRecoveryProtectionContainerMappingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSiteRecoveryProtectionContainerMapping_basic(ri, acceptance.Location(), acceptance.AltLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSiteRecoveryProtectionContainerMappingExists(resourceName),
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

func testAccAzureRMSiteRecoveryProtectionContainerMapping_basic(rInt int, location string, altLocation string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test1" {
  name     = "acctestRG-recovery1-%d"
  location = "%s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-vault-%d"
  location            = "${azurerm_resource_group.test1.location}"
  resource_group_name = "${azurerm_resource_group.test1.name}"
  sku                 = "Standard"
}

resource "azurerm_site_recovery_fabric" "test1" {
  resource_group_name = "${azurerm_resource_group.test1.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.test.name}"
  name                = "acctest-fabric1-%d"
  location            = "${azurerm_resource_group.test1.location}"
}

resource "azurerm_site_recovery_fabric" "test2" {
  resource_group_name = "${azurerm_resource_group.test1.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.test.name}"
  name                = "acctest-fabric2-%d"
  location            = "%s"
  depends_on          = ["azurerm_site_recovery_fabric.test1"]
}

resource "azurerm_site_recovery_protection_container" "test1" {
  resource_group_name  = "${azurerm_resource_group.test1.name}"
  recovery_vault_name  = "${azurerm_recovery_services_vault.test.name}"
  recovery_fabric_name = "${azurerm_site_recovery_fabric.test1.name}"
  name                 = "acctest-protection-cont1-%d"
}

resource "azurerm_site_recovery_protection_container" "test2" {
  resource_group_name  = "${azurerm_resource_group.test1.name}"
  recovery_vault_name  = "${azurerm_recovery_services_vault.test.name}"
  recovery_fabric_name = "${azurerm_site_recovery_fabric.test2.name}"
  name                 = "acctest-protection-cont2-%d"
}

resource "azurerm_site_recovery_replication_policy" "test" {
  resource_group_name                                  = "${azurerm_resource_group.test1.name}"
  recovery_vault_name                                  = "${azurerm_recovery_services_vault.test.name}"
  name                                                 = "acctest-policy-%d"
  recovery_point_retention_in_minutes                  = "${24 * 60}"
  application_consistent_snapshot_frequency_in_minutes = "${4 * 60}"
}

resource "azurerm_site_recovery_protection_container_mapping" "test" {
  resource_group_name                       = "${azurerm_resource_group.test1.name}"
  recovery_vault_name                       = "${azurerm_recovery_services_vault.test.name}"
  recovery_fabric_name                      = "${azurerm_site_recovery_fabric.test1.name}"
  recovery_source_protection_container_name = "${azurerm_site_recovery_protection_container.test1.name}"
  recovery_target_protection_container_id   = "${azurerm_site_recovery_protection_container.test2.id}"
  recovery_replication_policy_id            = "${azurerm_site_recovery_replication_policy.test.id}"
  name                                      = "mapping-%d"
}
`, rInt, location, rInt, rInt, rInt, altLocation, rInt, rInt, rInt, rInt)
}

func testCheckAzureRMSiteRecoveryProtectionContainerMappingExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		state, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroupName := state.Primary.Attributes["resource_group_name"]
		vaultName := state.Primary.Attributes["recovery_vault_name"]
		fabricName := state.Primary.Attributes["recovery_fabric_name"]
		protectionContainerName := state.Primary.Attributes["recovery_source_protection_container_name"]
		mappingName := state.Primary.Attributes["name"]

		client := acceptance.AzureProvider.Meta().(*clients.Client).RecoveryServices.ContainerMappingClient(resourceGroupName, vaultName)
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

func testCheckAzureRMSiteRecoveryProtectionContainerMappingDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_site_recovery_protection_container_mapping" {
			continue
		}

		resourceGroupName := rs.Primary.Attributes["resource_group_name"]
		vaultName := rs.Primary.Attributes["recovery_vault_name"]
		fabricName := rs.Primary.Attributes["recovery_fabric_name"]
		protectionContainerName := rs.Primary.Attributes["recovery_source_protection_container_name"]
		mappingName := rs.Primary.Attributes["name"]

		client := acceptance.AzureProvider.Meta().(*clients.Client).RecoveryServices.ContainerMappingClient(resourceGroupName, vaultName)
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := client.Get(ctx, fabricName, protectionContainerName, mappingName)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Container Mapping still exists:\n%#v", resp.Properties)
		}
	}

	return nil
}

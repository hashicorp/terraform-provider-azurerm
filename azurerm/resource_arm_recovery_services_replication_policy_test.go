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

func TestAccAzureRMRecoveryReplicationPolicy_basic(t *testing.T) {
	resourceName := "azurerm_recovery_services_replication_policy.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRecoveryReplicationPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRecoveryReplicationPolicy_basic(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRecoveryReplicationPolicyExists(resourceName),
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

func testAccAzureRMRecoveryReplicationPolicy_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_recovery_services_vault" "test" {
  name                = "acctest-vault-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"
}

resource "azurerm_recovery_services_replication_policy" "test" {
  resource_group_name                                  = "${azurerm_resource_group.test.name}"
  recovery_vault_name                                  = "${azurerm_recovery_services_vault.test.name}"
  name                                                 = "acctest-policy-%d"
  recovery_point_retention_in_minutes                  = "${24 * 60}"
  application_consistent_snapshot_frequency_in_minutes = "${4 * 60}"
}
`, rInt, location, rInt, rInt)
}

func testCheckAzureRMRecoveryReplicationPolicyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		state, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceGroupName := state.Primary.Attributes["resource_group_name"]
		vaultName := state.Primary.Attributes["recovery_vault_name"]
		policyName := state.Primary.Attributes["name"]

		client := acceptance.AzureProvider.Meta().(*clients.Client).RecoveryServices.ReplicationPoliciesClient(resourceGroupName, vaultName)
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := client.Get(ctx, policyName)
		if err != nil {
			return fmt.Errorf("Bad: Get on RecoveryServices.ReplicationPoliciesClient: %+v", err)
		}

		if resp.Response.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Replication Policies: %q does not exist", policyName)
		}

		return nil
	}
}

func testCheckAzureRMRecoveryReplicationPolicyDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_recovery_services_replication_policy" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		vaultName := rs.Primary.Attributes["recovery_vault_name"]
		policyName := rs.Primary.Attributes["name"]

		client := acceptance.AzureProvider.Meta().(*clients.Client).RecoveryServices.ReplicationPoliciesClient(resourceGroup, vaultName)
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := client.Get(ctx, policyName)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Replication Policy still exists:\n%#v", resp.Properties)
		}
	}

	return nil
}

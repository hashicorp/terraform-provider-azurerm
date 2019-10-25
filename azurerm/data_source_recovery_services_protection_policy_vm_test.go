package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMRecoveryServicesProtectionPolicyVm_basic(t *testing.T) {
	dataSourceName := "data.azurerm_recovery_services_protection_policy_vm.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRecoveryServicesProtectionPolicyVm_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRecoveryServicesProtectionPolicyVmExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "recovery_vault_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_group_name"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func testAccDataSourceRecoveryServicesProtectionPolicyVm_basic(rInt int, location string) string {
	return fmt.Sprintf(` 
%s

data "azurerm_recovery_services_protection_policy_vm" "test" {
  name                = "${azurerm_recovery_services_protection_policy_vm.test.name}"
  recovery_vault_name = "${azurerm_recovery_services_vault.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, testAccAzureRMRecoveryServicesProtectionPolicyVm_basicDaily(rInt, location))
}

package azurerm

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2016-06-01/recoveryservices"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMRecoveryServicesVault_basic(t *testing.T) {
	dataSourceName := "data.azurerm_recovery_services_vault.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRecoveryServicesVault_basic(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRecoveryServicesVaultExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "location"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_group_name"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "sku", string(recoveryservices.Standard)),
				),
			},
		},
	})
}

func testAccDataSourceRecoveryServicesVault_basic(rInt int, location string) string {
	return fmt.Sprintf(` 
%s

data "azurerm_recovery_services_vault" "test" {
  name                = "${azurerm_recovery_services_vault.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, testAccAzureRMRecoveryServicesVault_basic(rInt, location))
}

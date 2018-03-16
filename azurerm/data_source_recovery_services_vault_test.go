package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMRecoveryServicesVault_basic(t *testing.T) {
	dataSourceName := "data.azurerm_recovery_services_vault.test"
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRecoveryServicesVault_basic(ri, testLocation()),
				Check:  checkAccAzureRMRecoveryServicesVault_basic(dataSourceName),
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

package azurerm

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2016-06-01/recoveryservices"
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

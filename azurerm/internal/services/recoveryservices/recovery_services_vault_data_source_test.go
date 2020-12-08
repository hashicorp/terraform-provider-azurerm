package recoveryservices_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMRecoveryServicesVault_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_recovery_services_vault", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRecoveryServicesVault_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRecoveryServicesVaultExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "location"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku", "Standard"),
				),
			},
		},
	})
}

func testAccDataSourceRecoveryServicesVault_basic(data acceptance.TestData) string {
	template := testAccAzureRMRecoveryServicesVault_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_recovery_services_vault" "test" {
  name                = azurerm_recovery_services_vault.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}

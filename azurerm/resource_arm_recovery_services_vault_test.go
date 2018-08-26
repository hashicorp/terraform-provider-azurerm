package azurerm

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2016-06-01/recoveryservices"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMRecoveryServicesVault_basic(t *testing.T) {
	ri := acctest.RandInt()

	resourceName := "azurerm_recovery_services_vault.test"
	config := testAccAzureRMRecoveryServicesVault_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRecoveryServicesVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRecoveryServicesVaultExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "location"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "sku", string(recoveryservices.Standard)),
				),
			},
		},
	})
}

func testCheckAzureRMRecoveryServicesVaultDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_recovery_services_vault" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := testAccProvider.Meta().(*ArmClient).recoveryServicesVaultsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Recovery Services Vault still exists:\n%#v", resp)
	}

	return nil
}

func testCheckAzureRMRecoveryServicesVaultExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %q", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Recovery Services Vault: %q", name)
		}

		client := testAccProvider.Meta().(*ArmClient).recoveryServicesVaultsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Recovery Services Vault %q (resource group: %q) was not found: %+v", name, resourceGroup, err)
			}

			return fmt.Errorf("Bad: Get on recoveryServicesVaultsClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMRecoveryServicesVault_basic(rInt int, location string) string {
	return fmt.Sprintf(` 
resource "azurerm_resource_group" "test" { 
  name     = "acctestRG-%d" 
  location = "%s" 
} 

resource "azurerm_recovery_services_vault" "test" {
    name                = "acctest-%d"
    location            = "${azurerm_resource_group.test.location}"
    resource_group_name = "${azurerm_resource_group.test.name}"
    sku                 = "Standard"
}
 
`, rInt, location, rInt)
}

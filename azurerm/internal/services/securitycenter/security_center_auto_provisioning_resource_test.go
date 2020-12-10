package securitycenter_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccSecurityCenterAutoProvision_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_auto_provisioning", "test")

	// lintignore:AT001
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterAutoProvisioning_setting("On"),
				Check: resource.ComposeTestCheckFunc(
					testCheckSecurityCenterAutoProvisioningExists(data.ResourceName, "On"),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_provision", "On"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccSecurityCenterAutoProvisioning_setting("Off"),
				Check: resource.ComposeTestCheckFunc(
					testCheckSecurityCenterAutoProvisioningExists(data.ResourceName, "Off"),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_provision", "Off"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckSecurityCenterAutoProvisioningExists(resourceName string, expectedSetting string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).SecurityCenter.AutoProvisioningClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		idSplit := strings.Split(rs.Primary.Attributes["id"], "/")
		autoProvisionResourceName := idSplit[len(idSplit)-1]

		resp, err := client.Get(ctx, autoProvisionResourceName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Security Center auto provision %q was not found: %+v", autoProvisionResourceName, err)
			}

			return fmt.Errorf("Bad: GetAutoProvisioning: %+v", err)
		}

		// Check expected value
		if string(resp.AutoProvision) != expectedSetting {
			return fmt.Errorf("Security Center auto provision not expected, wanted %s, but got %s", expectedSetting, string(resp.AutoProvision))
		}

		return nil
	}
}

func testAccSecurityCenterAutoProvisioning_setting(setting string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_security_center_auto_provisioning" "test" {
  auto_provision = "%s"
}
`, setting)
}

package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSecurityCenterSetting_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_setting", "test")

	// lintignore:AT001
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSecurityCenterSetting_cfg("MCAS", true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSecurityCenterSettingExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "setting_name", "MCAS"),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMSecurityCenterSetting_cfg("MCAS", false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSecurityCenterSettingExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "setting_name", "MCAS"),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "false"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMSecurityCenterSetting_cfg("WDATP", true),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSecurityCenterSettingExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "setting_name", "WDATP"),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMSecurityCenterSetting_cfg("WDATP", false),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSecurityCenterSettingExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "setting_name", "WDATP"),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "false"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMSecurityCenterSettingExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).SecurityCenter.SettingClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		settingName := rs.Primary.Attributes["setting_name"]

		resp, err := client.Get(ctx, settingName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Security Center Setting %q was not found: %+v", settingName, err)
			}

			return fmt.Errorf("Bad: Get: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMSecurityCenterSetting_cfg(settingName string, enabled bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_security_center_setting" "test" {
  setting_name = "%s"
  enabled      = "%t"
}
`, settingName, enabled)
}

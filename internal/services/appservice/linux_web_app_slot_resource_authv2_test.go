package appservice_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccLinuxWebAppSlot_withAuthV2AzureActiveDirectory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withAuthV2AzureActiveDirectory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r LinuxWebAppSlotResource) withAuthV2AzureActiveDirectory(data acceptance.TestData) string {
	secretSettingName := "MICROSOFT_PROVIDER_AUTHENTICATION_SECRET"
	secretSettingValue := "902D17F6-FD6B-4E44-BABB-58E788DCD907"

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_linux_web_app.test.id

  site_config {}

  app_settings = {
    "%[3]s" = "%[4]s"
  }

  auth_v2_settings {
    auth_enabled           = true
    unauthenticated_action = "Return401"
    active_directory {
      client_id                  = data.azurerm_client_config.current.client_id
      client_secret_setting_name = "%[3]s"
      tenant_auth_endpoint       = "https://sts.windows.net/%[5]s/v2.0"
    }
    login {}
  }
}
`, r.baseTemplate(data), data.RandomInteger, secretSettingName, secretSettingValue, data.Client().TenantID)
}

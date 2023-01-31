package appservice_test

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"os"
	"testing"
)

func TestAccLinuxWebApp_authV2Basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authV2(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("app,linux"),
			),
		},
		data.ImportStep(),
	})
}
func TestAccLinuxWebApp_authV2WithApple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authV2Apple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("app,linux"),
			),
		},
		data.ImportStep(),
	})
}

func (r LinuxWebAppResource) authV2(data acceptance.TestData) string {
	secretSettingName := "MICROSOFT_PROVIDER_AUTHENTICATION_SECRET"
	secretSettingValue := os.Getenv("ARM_CLIENT_SECRET")
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_client_config" "current" {}

resource "azurerm_linux_web_app" "test" {
  name                = "acctestLWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {}

  app_settings = {
    "%[3]s" = "%[4]s" 
  }

  auth_v2_settings {
    auth_enabled = true
    unauthenticated_action = "Return401"
    active_directory {
      client_id = data.azurerm_client_config.current.client_id
      client_secret_setting_name = "%[3]s"      
      tenant_auth_endpoint = "https://sts.windows.net/%[5]s/v2.0"
    }
    login {}
  }
}
`, r.baseTemplate(data), data.RandomInteger, secretSettingName, secretSettingValue, data.Client().TenantID)
}

func (r LinuxWebAppResource) authV2Apple(data acceptance.TestData) string {
	secretSettingName := "APPLE_PROVIDER_AUTHENTICATION_SECRET"
	secretSettingValue := "902D17F6-FD6B-4E44-BABB-58E788DCD907"
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_client_config" "current" {}

resource "azurerm_linux_web_app" "test" {
  name                = "acctestLWA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {}

  app_settings = {
    "%[3]s" = "%[4]s" 
  }

  sticky_settings {
    app_setting_names = ["%[3]s"]
  }

  auth_v2_settings {
    auth_enabled = true
    unauthenticated_action = "Return401"
    
    apple {
      client_id = "testAppleID"
      client_secret_setting_name = "%[3]s"
    }

    login {}
  }
}
`, r.baseTemplate(data), data.RandomInteger, secretSettingName, secretSettingValue)
}

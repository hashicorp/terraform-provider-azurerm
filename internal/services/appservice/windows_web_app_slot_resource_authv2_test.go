// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccWindowsWebAppSlot_withAuthV2AzureActiveDirectory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authV2AzureActiveDirectory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsWebAppSlot_withAuthV2AzureActiveDirectoryNoSecretName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authV2AzureActiveDirectoryNoSecretName(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsWebAppSlot_withAuthV2Apple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authV2Apple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsWebAppSlot_withAuthV2CustomOIDC(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authV2CustomOIDC(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsWebAppSlot_withAuthV2Facebook(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authV2Facebook(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsWebAppSlot_withAuthV2Github(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authV2Github(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsWebAppSlot_withAuthV2Google(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authV2Google(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsWebAppSlot_withAuthV2Microsoft(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authV2Microsoft(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsWebAppSlot_withAuthV2Twitter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authV2Twitter(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsWebAppSlot_withAuthV2MultipleAuths(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app_slot", "test")
	r := WindowsWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authV2Multi(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func (r WindowsWebAppSlotResource) authV2AzureActiveDirectory(data acceptance.TestData) string {
	secretSettingName := "MICROSOFT_PROVIDER_AUTHENTICATION_SECRET"
	secretSettingValue := "902D17F6-FD6B-4E44-BABB-58E788DCD907"

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_client_config" "current" {}

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_windows_web_app.test.id

  site_config {}

  app_settings = {
    "%[3]s" = "%[4]s"
  }

  auth_settings_v2 {
    auth_enabled           = true
    unauthenticated_action = "Return401"
    active_directory_v2 {
      client_id                  = data.azurerm_client_config.current.client_id
      client_secret_setting_name = "%[3]s"
      tenant_auth_endpoint       = "https://sts.windows.net/%[5]s/v2.0"
      allowed_applications       = ["WhoopsMissedThisOne"]
    }
    login {}
  }
}
`, r.baseTemplate(data), data.RandomInteger, secretSettingName, secretSettingValue, data.Client().TenantID)
}

func (r WindowsWebAppSlotResource) authV2AzureActiveDirectoryNoSecretName(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_client_config" "current" {}

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_windows_web_app.test.id

  site_config {}

  auth_settings_v2 {
    auth_enabled           = true
    unauthenticated_action = "Return401"
    active_directory_v2 {
      client_id            = data.azurerm_client_config.current.client_id
      tenant_auth_endpoint = "https://sts.windows.net/%[3]s/v2.0"
    }
    login {}
  }
}
`, r.baseTemplate(data), data.RandomInteger, data.Client().TenantID)
}

func (r WindowsWebAppSlotResource) authV2Apple(data acceptance.TestData) string {
	secretSettingName := "APPLE_PROVIDER_AUTHENTICATION_SECRET"
	secretSettingValue := "902D17F6-FD6B-4E44-BABB-58E788DCD907"

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_windows_web_app.test.id

  site_config {}

  app_settings = {
    "%[3]s" = "%[4]s"
  }

  auth_settings_v2 {
    auth_enabled           = true
    unauthenticated_action = "Return401"

    apple_v2 {
      client_id                  = "testAppleID"
      client_secret_setting_name = "%[3]s"
    }

    login {}
  }
}
`, r.baseTemplate(data), data.RandomInteger, secretSettingName, secretSettingValue, data.Client().TenantID)
}

func (r WindowsWebAppSlotResource) authV2CustomOIDC(data acceptance.TestData) string {
	secretSettingName := "TESTCUSTOM_PROVIDER_AUTHENTICATION_SECRET"
	secretSettingValue := "902D17F6-FD6B-4E44-BABB-58E788DCD907"

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_windows_web_app.test.id

  site_config {}

  app_settings = {
    "%[3]s" = "%[4]s"
  }

  auth_settings_v2 {
    auth_enabled           = true
    unauthenticated_action = "Return401"

    custom_oidc_v2 {
      name                          = "testcustom"
      client_id                     = "testCustomID"
      openid_configuration_endpoint = "https://oidc.testcustom.contoso.com/auth"
    }

    login {}
  }
}
`, r.baseTemplate(data), data.RandomInteger, secretSettingName, secretSettingValue, data.Client().TenantID)
}

func (r WindowsWebAppSlotResource) authV2Facebook(data acceptance.TestData) string {
	secretSettingName := "FACEBOOK_PROVIDER_AUTHENTICATION_SECRET"
	secretSettingValue := "902D17F6-FD6B-4E44-BABB-58E788DCD907"

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_windows_web_app.test.id

  site_config {}

  app_settings = {
    "%[3]s" = "%[4]s"
  }

  auth_settings_v2 {
    auth_enabled           = true
    unauthenticated_action = "RedirectToLoginPage"

    facebook_v2 {
      app_id                  = "testFacebookID"
      app_secret_setting_name = "%[3]s"
    }

    login {}
  }
}
`, r.baseTemplate(data), data.RandomInteger, secretSettingName, secretSettingValue, data.Client().TenantID)
}

func (r WindowsWebAppSlotResource) authV2Github(data acceptance.TestData) string {
	secretSettingName := "GITHUB_PROVIDER_AUTHENTICATION_SECRET"
	secretSettingValue := "902D17F6-FD6B-4E44-BABB-58E788DCD907"

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_windows_web_app.test.id

  site_config {}

  app_settings = {
    "%[3]s" = "%[4]s"
  }

  auth_settings_v2 {
    auth_enabled           = true
    unauthenticated_action = "RedirectToLoginPage"

    github_v2 {
      client_id                  = "testGithubID"
      client_secret_setting_name = "%[3]s"
    }

    login {}
  }
}
`, r.baseTemplate(data), data.RandomInteger, secretSettingName, secretSettingValue, data.Client().TenantID)
}

func (r WindowsWebAppSlotResource) authV2Google(data acceptance.TestData) string {
	secretSettingName := "GOOGLE_PROVIDER_AUTHENTICATION_SECRET"
	secretSettingValue := "902D17F6-FD6B-4E44-BABB-58E788DCD907"

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_windows_web_app.test.id

  site_config {}

  app_settings = {
    "%[3]s" = "%[4]s"
  }

  auth_settings_v2 {
    auth_enabled           = true
    unauthenticated_action = "RedirectToLoginPage"

    google_v2 {
      client_id                  = "testGoogleID"
      client_secret_setting_name = "%[3]s"
    }

    login {}
  }
}
`, r.baseTemplate(data), data.RandomInteger, secretSettingName, secretSettingValue, data.Client().TenantID)
}

func (r WindowsWebAppSlotResource) authV2Microsoft(data acceptance.TestData) string {
	secretSettingName := "MICROSOFT_PROVIDER_AUTHENTICATION_SECRET"
	secretSettingValue := "902D17F6-FD6B-4E44-BABB-58E788DCD907"

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_windows_web_app.test.id

  site_config {}

  app_settings = {
    "%[3]s" = "%[4]s"
  }

  auth_settings_v2 {
    auth_enabled           = true
    unauthenticated_action = "RedirectToLoginPage"

    microsoft_v2 {
      client_id                  = "testMSFTID"
      client_secret_setting_name = "%[3]s"
    }

    login {}
  }
}
`, r.baseTemplate(data), data.RandomInteger, secretSettingName, secretSettingValue, data.Client().TenantID)
}

func (r WindowsWebAppSlotResource) authV2Twitter(data acceptance.TestData) string {
	secretSettingName := "TWITTER_PROVIDER_AUTHENTICATION_SECRET"
	secretSettingValue := "902D17F6-FD6B-4E44-BABB-58E788DCD907"

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_windows_web_app.test.id

  site_config {}

  app_settings = {
    "%[3]s" = "%[4]s"
  }

  auth_settings_v2 {
    auth_enabled           = true
    unauthenticated_action = "RedirectToLoginPage"

    twitter_v2 {
      consumer_key                 = "testTwitterKey"
      consumer_secret_setting_name = "%[3]s"
    }

    login {}
  }
}
`, r.baseTemplate(data), data.RandomInteger, secretSettingName, secretSettingValue, data.Client().TenantID)
}

func (r WindowsWebAppSlotResource) authV2Multi(data acceptance.TestData) string {
	secretSettingValue := "902D17F6-FD6B-4E44-BABB-58E788DCD907"

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app_slot" "test" {
  name           = "acctestWAS-%[2]d"
  app_service_id = azurerm_windows_web_app.test.id

  site_config {}

  app_settings = {
    "APPLE_PROVIDER_AUTHENTICATION_SECRET"     = "%[3]s"
    "FACEBOOK_PROVIDER_AUTHENTICATION_SECRET"  = "%[3]s"
    "GITHUB_PROVIDER_AUTHENTICATION_SECRET"    = "%[3]s"
    "GOOGLE_PROVIDER_AUTHENTICATION_SECRET"    = "%[3]s"
    "MICROSOFT_PROVIDER_AUTHENTICATION_SECRET" = "%[3]s"
    "TWITTER_PROVIDER_AUTHENTICATION_SECRET"   = "%[3]s"
  }

  auth_settings_v2 {
    auth_enabled           = true
    unauthenticated_action = "RedirectToLoginPage"

    apple_v2 {
      client_id                  = "testAppleID"
      client_secret_setting_name = "APPLE_PROVIDER_AUTHENTICATION_SECRET"
    }

    facebook_v2 {
      app_id                  = "testFacebookID"
      app_secret_setting_name = "FACEBOOK_PROVIDER_AUTHENTICATION_SECRET"
    }

    github_v2 {
      client_id                  = "testGithubID"
      client_secret_setting_name = "GITHUB_PROVIDER_AUTHENTICATION_SECRET"
    }

    google_v2 {
      client_id                  = "testGoogleID"
      client_secret_setting_name = "GOOGLE_PROVIDER_AUTHENTICATION_SECRET"
    }

    microsoft_v2 {
      client_id                  = "testMSFTID"
      client_secret_setting_name = "MICROSOFT_PROVIDER_AUTHENTICATION_SECRET"
    }

    twitter_v2 {
      consumer_key                 = "testTwitterKey"
      consumer_secret_setting_name = "TWITTER_PROVIDER_AUTHENTICATION_SECRET"
    }

    login {}
  }
}
`, r.baseTemplate(data), data.RandomInteger, secretSettingValue)
}

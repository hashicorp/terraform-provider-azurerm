// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccLinuxFunctionApp_authV2AzureActiveDirectory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_function_app", "test")
	r := LinuxFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authV2AzureActiveDirectory(data, SkuBasicPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxFunctionApp_authV2AzureActiveDirectoryNoSecretName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_function_app", "test")
	r := LinuxFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authV2AzureActiveDirectoryNoSecretName(data, SkuBasicPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxFunctionApp_authV2AzureActiveDirectoryRemove(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_function_app", "test")
	r := LinuxFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authV2AzureActiveDirectory(data, SkuBasicPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.basic(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxFunctionApp_authV2Apple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_function_app", "test")
	r := LinuxFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authV2Apple(data, SkuBasicPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxFunctionApp_authV2AppleCustomSettingName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_function_app", "test")
	r := LinuxFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authV2AppleCustomSetting(data, SkuBasicPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxFunctionApp_authV2CustomOIDC(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_function_app", "test")
	r := LinuxFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authV2CustomOIDC(data, SkuBasicPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxFunctionApp_authV2Facebook(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_function_app", "test")
	r := LinuxFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authV2Facebook(data, SkuBasicPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxFunctionApp_authV2Github(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_function_app", "test")
	r := LinuxFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authV2Github(data, SkuBasicPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxFunctionApp_authV2Google(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_function_app", "test")
	r := LinuxFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authV2Google(data, SkuBasicPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxFunctionApp_authV2Microsoft(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_function_app", "test")
	r := LinuxFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authV2Microsoft(data, SkuBasicPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxFunctionApp_authV2Twitter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_function_app", "test")
	r := LinuxFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authV2Twitter(data, SkuBasicPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxFunctionApp_authV2MultipleAuths(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_function_app", "test")
	r := LinuxFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authV2Multi(data, SkuBasicPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxFunctionApp_authV2Update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_function_app", "test")
	r := LinuxFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authV2AzureActiveDirectory(data, SkuBasicPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.authV2Google(data, SkuBasicPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxFunctionApp_authV2UpgradeFromV1(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_function_app", "test")
	r := LinuxFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.standardComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.completeAuthV2(data, SkuBasicPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func (r LinuxFunctionAppResource) authV2AzureActiveDirectory(data acceptance.TestData, planSku string) string {
	secretSettingName := "MICROSOFT_PROVIDER_AUTHENTICATION_SECRET"
	secretSettingValue := "902D17F6-FD6B-4E44-BABB-58E788DCD907"
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {}

%s

data "azurerm_client_config" "current" {}

resource "azuread_group" "test" {
  display_name     = "acctestspa-%d"
  security_enabled = true
}

resource "azurerm_linux_function_app" "test" {
  name                = "acctest-LFA-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {}

  app_settings = {
    "%[3]s" = "%[4]s"
  }

  sticky_settings {
    app_setting_names = ["%[3]s"]
  }

  auth_settings_v2 {
    auth_enabled           = true
    unauthenticated_action = "Return401"
    active_directory_v2 {
      client_id                  = data.azurerm_client_config.current.client_id
      client_secret_setting_name = "%[3]s"
      tenant_auth_endpoint       = "https://sts.windows.net/%[5]s/v2.0"
      allowed_groups             = [azuread_group.test.object_id]
      allowed_applications       = ["WhoopsMissedThisOne"]
    }

    login {
      token_store_enabled = true
    }
  }
}
`, r.template(data, planSku), data.RandomInteger, secretSettingName, secretSettingValue, data.Client().TenantID)
}

func (r LinuxFunctionAppResource) authV2AzureActiveDirectoryNoSecretName(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {}

%s

data "azurerm_client_config" "current" {}

resource "azurerm_linux_function_app" "test" {
  name                = "acctest-LFA-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {}


  auth_settings_v2 {
    auth_enabled           = true
    unauthenticated_action = "Return401"
    active_directory_v2 {
      client_id            = data.azurerm_client_config.current.client_id
      tenant_auth_endpoint = "https://sts.windows.net/%[3]s/v2.0"
    }

    login {
      token_store_enabled = true
    }
  }
}
`, r.template(data, planSku), data.RandomInteger, data.Client().TenantID)
}

func (r LinuxFunctionAppResource) authV2Apple(data acceptance.TestData, planSku string) string {
	secretSettingName := "APPLE_PROVIDER_AUTHENTICATION_SECRET"
	secretSettingValue := "902D17F6-FD6B-4E44-BABB-58E788DCD907"
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_client_config" "current" {}

resource "azurerm_linux_function_app" "test" {
  name                = "acctest-LFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

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
`, r.template(data, planSku), data.RandomInteger, secretSettingName, secretSettingValue)
}

func (r LinuxFunctionAppResource) authV2AppleCustomSetting(data acceptance.TestData, planSku string) string {
	secretSettingName := "TEST_PROVIDER_AUTHENTICATION_SECRET"
	secretSettingValue := "902D17F6-FD6B-4E44-BABB-58E788DCD907"
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_client_config" "current" {}

resource "azurerm_linux_function_app" "test" {
  name                = "acctest-LFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

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
`, r.template(data, planSku), data.RandomInteger, secretSettingName, secretSettingValue)
}

// static web app? - Need to add test when Static Web Apps are deployable from TF.

func (r LinuxFunctionAppResource) authV2CustomOIDC(data acceptance.TestData, planSku string) string {
	secretSettingName := "TESTCUSTOM_PROVIDER_AUTHENTICATION_SECRET"
	secretSettingValue := "902D17F6-FD6B-4E44-BABB-58E788DCD907"
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_client_config" "current" {}

resource "azurerm_linux_function_app" "test" {
  name                = "acctest-LFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {}

  app_settings = {
    "%[3]s" = "%[4]s"
  }

  sticky_settings {
    app_setting_names = ["%[3]s"]
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
`, r.template(data, planSku), data.RandomInteger, secretSettingName, secretSettingValue)
}

func (r LinuxFunctionAppResource) authV2Facebook(data acceptance.TestData, planSku string) string {
	secretSettingName := "FACEBOOK_PROVIDER_AUTHENTICATION_SECRET"
	secretSettingValue := "902D17F6-FD6B-4E44-BABB-58E788DCD907"
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_client_config" "current" {}

resource "azurerm_linux_function_app" "test" {
  name                = "acctest-LFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {}

  app_settings = {
    "%[3]s" = "%[4]s"
  }

  sticky_settings {
    app_setting_names = ["%[3]s"]
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
`, r.template(data, planSku), data.RandomInteger, secretSettingName, secretSettingValue)
}

func (r LinuxFunctionAppResource) authV2Github(data acceptance.TestData, planSku string) string {
	secretSettingName := "GITHUB_PROVIDER_AUTHENTICATION_SECRET"
	secretSettingValue := "902D17F6-FD6B-4E44-BABB-58E788DCD907"
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_client_config" "current" {}

resource "azurerm_linux_function_app" "test" {
  name                = "acctest-LFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {}

  app_settings = {
    "%[3]s" = "%[4]s"
  }

  sticky_settings {
    app_setting_names = ["%[3]s"]
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
`, r.template(data, planSku), data.RandomInteger, secretSettingName, secretSettingValue)
}

func (r LinuxFunctionAppResource) authV2Google(data acceptance.TestData, planSku string) string {
	secretSettingName := "GOOGLE_PROVIDER_AUTHENTICATION_SECRET"
	secretSettingValue := "902D17F6-FD6B-4E44-BABB-58E788DCD907"
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_client_config" "current" {}

resource "azurerm_linux_function_app" "test" {
  name                = "acctest-LFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {}

  app_settings = {
    "%[3]s" = "%[4]s"
  }

  sticky_settings {
    app_setting_names = ["%[3]s"]
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
`, r.template(data, planSku), data.RandomInteger, secretSettingName, secretSettingValue)
}

func (r LinuxFunctionAppResource) authV2Microsoft(data acceptance.TestData, planSku string) string {
	secretSettingName := "MICROSOFT_PROVIDER_AUTHENTICATION_SECRET"
	secretSettingValue := "902D17F6-FD6B-4E44-BABB-58E788DCD907"
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_client_config" "current" {}

resource "azurerm_linux_function_app" "test" {
  name                = "acctest-LFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {}

  app_settings = {
    "%[3]s" = "%[4]s"
  }

  sticky_settings {
    app_setting_names = ["%[3]s"]
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
`, r.template(data, planSku), data.RandomInteger, secretSettingName, secretSettingValue)
}

func (r LinuxFunctionAppResource) authV2Twitter(data acceptance.TestData, planSku string) string {
	secretSettingName := "TWITTER_PROVIDER_AUTHENTICATION_SECRET"
	secretSettingValue := "902D17F6-FD6B-4E44-BABB-58E788DCD907"
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_client_config" "current" {}

resource "azurerm_linux_function_app" "test" {
  name                = "acctest-LFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {}

  app_settings = {
    "%[3]s" = "%[4]s"
  }

  sticky_settings {
    app_setting_names = ["%[3]s"]
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
`, r.template(data, planSku), data.RandomInteger, secretSettingName, secretSettingValue)
}

func (r LinuxFunctionAppResource) authV2Multi(data acceptance.TestData, planSku string) string {
	secretSettingValue := "902D17F6-FD6B-4E44-BABB-58E788DCD907"
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_client_config" "current" {}

resource "azurerm_linux_function_app" "test" {
  name                = "acctest-LFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {}

  app_settings = {
    "APPLE_PROVIDER_AUTHENTICATION_SECRET"     = "%[3]s"
    "FACEBOOK_PROVIDER_AUTHENTICATION_SECRET"  = "%[3]s"
    "GITHUB_PROVIDER_AUTHENTICATION_SECRET"    = "%[3]s"
    "GOOGLE_PROVIDER_AUTHENTICATION_SECRET"    = "%[3]s"
    "MICROSOFT_PROVIDER_AUTHENTICATION_SECRET" = "%[3]s"
    "TWITTER_PROVIDER_AUTHENTICATION_SECRET"   = "%[3]s"
  }

  sticky_settings {
    app_setting_names = [
      "APPLE_PROVIDER_AUTHENTICATION_SECRET",
      "FACEBOOK_PROVIDER_AUTHENTICATION_SECRET",
      "GITHUB_PROVIDER_AUTHENTICATION_SECRET",
      "GOOGLE_PROVIDER_AUTHENTICATION_SECRET",
      "MICROSOFT_PROVIDER_AUTHENTICATION_SECRET",
      "TWITTER_PROVIDER_AUTHENTICATION_SECRET",
    ]
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
`, r.template(data, planSku), data.RandomInteger, secretSettingValue)
}

func (r LinuxFunctionAppResource) completeAuthV2(data acceptance.TestData, planSku string) string {
	secretSettingValue := "902D17F6-FD6B-4E44-BABB-58E788DCD907"
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acct-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}


resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_linux_function_app" "test" {
  name                = "acctest-LFA-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  app_settings = {
    "foo"                                      = "bar"
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


  backup {
    name                = "acctest"
    storage_account_url = "https://${azurerm_storage_account.test.name}.blob.core.windows.net/${azurerm_storage_container.test.name}${data.azurerm_storage_account_sas.test.sas}&sr=b"
    schedule {
      frequency_interval = 7
      frequency_unit     = "Day"
    }
  }

  builtin_logging_enabled            = false
  client_certificate_enabled         = true
  client_certificate_mode            = "OptionalInteractiveUser"
  client_certificate_exclusion_paths = "/foo;/bar;/hello;/world"

  connection_string {
    name  = "First"
    value = "some-postgresql-connection-string"
    type  = "PostgreSQL"
  }

  enabled                     = false
  functions_extension_version = "~3"
  https_only                  = true

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  site_config {
    always_on                              = true
    app_command_line                       = "whoami"
    api_definition_url                     = "https://example.com/azure_function_app_def.json"
    application_insights_connection_string = azurerm_application_insights.test.connection_string

    application_stack {
      python_version = "3.8"
    }

    container_registry_use_managed_identity       = true
    container_registry_managed_identity_client_id = azurerm_user_assigned_identity.test.client_id

    default_documents = [
      "first.html",
      "second.jsp",
      "third.aspx",
      "hostingstart.html",
    ]

    http2_enabled = true

    ip_restriction {
      ip_address = "10.10.10.10/32"
      name       = "test-restriction"
      priority   = 123
      action     = "Allow"
      headers {
        x_azure_fdid      = ["55ce4ed1-4b06-4bf1-b40e-4638452104da"]
        x_fd_health_probe = ["1"]
        x_forwarded_for   = ["9.9.9.9/32", "2002::1234:abcd:ffff:c0a8:101/64"]
        x_forwarded_host  = ["example.com"]
      }
    }

    load_balancing_mode       = "LeastResponseTime"
    pre_warmed_instance_count = 2
    remote_debugging_enabled  = true
    remote_debugging_version  = "VS2022"

    scm_ip_restriction {
      ip_address = "10.20.20.20/32"
      name       = "test-scm-restriction"
      priority   = 123
      action     = "Allow"
      headers {
        x_azure_fdid      = ["55ce4ed1-4b06-4bf1-b40e-4638452104da"]
        x_fd_health_probe = ["1"]
        x_forwarded_for   = ["9.9.9.9/32", "2002::1234:abcd:ffff:c0a8:101/64"]
        x_forwarded_host  = ["example.com"]
      }
    }

    scm_ip_restriction {
      ip_address = "fd80::/64"
      name       = "test-scm-restriction-v6"
      priority   = 124
      action     = "Allow"
      headers {
        x_azure_fdid      = ["55ce4ed1-4b06-4bf1-b40e-4638452104da"]
        x_fd_health_probe = ["1"]
        x_forwarded_for   = ["9.9.9.9/32", "2002::1234:abcd:ffff:c0a8:101/64"]
        x_forwarded_host  = ["example.com"]
      }
    }

    use_32_bit_worker                 = true
    websockets_enabled                = true
    ftps_state                        = "FtpsOnly"
    health_check_path                 = "/health-check"
    health_check_eviction_time_in_min = 5
    worker_count                      = 3

    minimum_tls_version     = "1.1"
    scm_minimum_tls_version = "1.1"

    cors {
      allowed_origins = [
        "https://www.contoso.com",
        "www.contoso.com",
      ]

      support_credentials = true
    }

    vnet_route_all_enabled = true
  }

  sticky_settings {
    app_setting_names       = ["foo", "secret"]
    connection_string_names = ["First"]
  }

  tags = {
    terraform = "true"
    Env       = "AccTest"
  }
}
`, r.storageContainerTemplate(data, planSku), data.RandomInteger, secretSettingValue)
}

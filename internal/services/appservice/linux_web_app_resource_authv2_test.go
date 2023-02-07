package appservice_test

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"os"
	"testing"
)

func TestAccLinuxWebApp_authV2AzureActiveDirectory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authV2AzureActiveDirectory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("app,linux"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_authV2Apple(t *testing.T) {
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

func TestAccLinuxWebApp_authV2Facebook(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authV2Facebook(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("app,linux"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_authV2Github(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authV2Github(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("app,linux"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_authV2Google(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authV2Google(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("app,linux"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_authV2Microsoft(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authV2Microsoft(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("app,linux"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_authV2Twitter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authV2Twitter(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("app,linux"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLinuxWebApp_authV2MultipleAuths(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app", "test")
	r := LinuxWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authV2Multi(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("app,linux"),
			),
		},
		data.ImportStep(),
	})
}

func (r LinuxWebAppResource) authV2AzureActiveDirectory(data acceptance.TestData) string {
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

  sticky_settings {
    app_setting_names = ["%[3]s"]
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

// static web app?

// custom OIDC?

// FACEBOOK_PROVIDER_AUTHENTICATION_SECRET
func (r LinuxWebAppResource) authV2Facebook(data acceptance.TestData) string {
	secretSettingName := "FACEBOOK_PROVIDER_AUTHENTICATION_SECRET"
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
    auth_enabled           = true
    unauthenticated_action = "RedirectToLoginPage"
    
    facebook {
      app_id                  = "testFacebookID"
      app_secret_setting_name = "%[3]s"
    }

    login {}
  }
}
`, r.baseTemplate(data), data.RandomInteger, secretSettingName, secretSettingValue)
}

// GITHUB_PROVIDER_AUTHENTICATION_SECRET
func (r LinuxWebAppResource) authV2Github(data acceptance.TestData) string {
	secretSettingName := "GITHUB_PROVIDER_AUTHENTICATION_SECRET"
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
    auth_enabled           = true
    unauthenticated_action = "RedirectToLoginPage"
    
    github {
      client_id                  = "testGithubID"
      client_secret_setting_name = "%[3]s"
    }

    login {}
  }
}
`, r.baseTemplate(data), data.RandomInteger, secretSettingName, secretSettingValue)
}

// GOOGLE_PROVIDER_AUTHENTICATION_SECRET
func (r LinuxWebAppResource) authV2Google(data acceptance.TestData) string {
	secretSettingName := "GOOGLE_PROVIDER_AUTHENTICATION_SECRET"
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
    auth_enabled           = true
    unauthenticated_action = "RedirectToLoginPage"
    
    google {
      client_id                  = "testGoogleID"
      client_secret_setting_name = "%[3]s"
    }

    login {}
  }
}
`, r.baseTemplate(data), data.RandomInteger, secretSettingName, secretSettingValue)
}

// MICROSOFT_PROVIDER_AUTHENTICATION_SECRET
func (r LinuxWebAppResource) authV2Microsoft(data acceptance.TestData) string {
	secretSettingName := "MICROSOFT_PROVIDER_AUTHENTICATION_SECRET"
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
    auth_enabled           = true
    unauthenticated_action = "RedirectToLoginPage"
    
    microsoft {
      client_id                  = "testMSFTID"
      client_secret_setting_name = "%[3]s"
    }

    login {}
  }
}
`, r.baseTemplate(data), data.RandomInteger, secretSettingName, secretSettingValue)
}

// TWITTER_PROVIDER_AUTHENTICATION_SECRET
func (r LinuxWebAppResource) authV2Twitter(data acceptance.TestData) string {
	secretSettingName := "TWITTER_PROVIDER_AUTHENTICATION_SECRET"
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
    auth_enabled           = true
    unauthenticated_action = "RedirectToLoginPage"
    
    twitter {
      consumer_key                 = "testTwitterKey"
      consumer_secret_setting_name = "%[3]s"
    }

    login {}
  }
}
`, r.baseTemplate(data), data.RandomInteger, secretSettingName, secretSettingValue)
}

// multi / Complete
func (r LinuxWebAppResource) authV2Multi(data acceptance.TestData) string {
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

  auth_v2_settings {
    auth_enabled           = true
    unauthenticated_action = "RedirectToLoginPage"

    apple {
      client_id                  = "testAppleID"
      client_secret_setting_name = "APPLE_PROVIDER_AUTHENTICATION_SECRET"
    }

    facebook {
      app_id                  = "testFacebookID"
      app_secret_setting_name = "FACEBOOK_PROVIDER_AUTHENTICATION_SECRET"
    }

    github {
      client_id                  = "testGithubID"
      client_secret_setting_name = "GITHUB_PROVIDER_AUTHENTICATION_SECRET"
    }

    google {
      client_id                  = "testGoogleID"
      client_secret_setting_name = "GOOGLE_PROVIDER_AUTHENTICATION_SECRET"
    }

    microsoft {
      client_id                  = "testMSFTID"
      client_secret_setting_name = "MICROSOFT_PROVIDER_AUTHENTICATION_SECRET"
    }

    twitter {
      consumer_key                 = "testTwitterKey"
      consumer_secret_setting_name = "TWITTER_PROVIDER_AUTHENTICATION_SECRET"
    }

    login {}
  }
}
`, r.baseTemplate(data), data.RandomInteger, secretSettingValue)
}

// update

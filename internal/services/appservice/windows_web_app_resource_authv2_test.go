// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccWindowsWebApp_authV2AzureActiveDirectory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

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

func TestAccWindowsWebApp_authV2AzureActiveDirectoryNoSecretName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

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

func TestAccWindowsWebApp_authV2Apple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

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

func TestAccWindowsWebApp_authV2CustomOIDC(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authV2CustomOIDC(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auth_settings_v2.0.custom_oidc_v2.0.client_secret_setting_name").HasValue("TESTCUSTOM_PROVIDER_AUTHENTICATION_SECRET"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsWebApp_authV2Facebook(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

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

func TestAccWindowsWebApp_authV2Github(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

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

func TestAccWindowsWebApp_authV2Google(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

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

func TestAccWindowsWebApp_authV2Microsoft(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

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

func TestAccWindowsWebApp_authV2Twitter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

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

func TestAccWindowsWebApp_authV2MultipleAuths(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

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

func TestAccWindowsWebApp_authV2Update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authV2Apple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.authV2Facebook(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsWebApp_authV2UpgradeFromV1(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_web_app", "test")
	r := WindowsWebAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.completeAuthV2(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func (r WindowsWebAppResource) authV2AzureActiveDirectory(data acceptance.TestData) string {
	secretSettingName := "MICROSOFT_PROVIDER_AUTHENTICATION_SECRET"
	secretSettingValue := "902D17F6-FD6B-4E44-BABB-58E788DCD907"

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_client_config" "current" {}

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%[2]d"
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

func (r WindowsWebAppResource) authV2AzureActiveDirectoryNoSecretName(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_client_config" "current" {}

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

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

func (r WindowsWebAppResource) authV2Apple(data acceptance.TestData) string {
	secretSettingName := "APPLE_PROVIDER_AUTHENTICATION_SECRET"
	secretSettingValue := "902D17F6-FD6B-4E44-BABB-58E788DCD907"

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%[2]d"
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
`, r.baseTemplate(data), data.RandomInteger, secretSettingName, secretSettingValue)
}

func (r WindowsWebAppResource) authV2CustomOIDC(data acceptance.TestData) string {
	secretSettingName := "TESTCUSTOM_PROVIDER_AUTHENTICATION_SECRET"
	secretSettingValue := "902D17F6-FD6B-4E44-BABB-58E788DCD907"

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%[2]d"
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
`, r.baseTemplate(data), data.RandomInteger, secretSettingName, secretSettingValue)
}

func (r WindowsWebAppResource) authV2Facebook(data acceptance.TestData) string {
	secretSettingName := "FACEBOOK_PROVIDER_AUTHENTICATION_SECRET"
	secretSettingValue := "902D17F6-FD6B-4E44-BABB-58E788DCD907"

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%[2]d"
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
`, r.baseTemplate(data), data.RandomInteger, secretSettingName, secretSettingValue)
}

func (r WindowsWebAppResource) authV2Github(data acceptance.TestData) string {
	secretSettingName := "GITHUB_PROVIDER_AUTHENTICATION_SECRET"
	secretSettingValue := "902D17F6-FD6B-4E44-BABB-58E788DCD907"

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%[2]d"
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
`, r.baseTemplate(data), data.RandomInteger, secretSettingName, secretSettingValue)
}

func (r WindowsWebAppResource) authV2Google(data acceptance.TestData) string {
	secretSettingName := "GOOGLE_PROVIDER_AUTHENTICATION_SECRET"
	secretSettingValue := "902D17F6-FD6B-4E44-BABB-58E788DCD907"

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%[2]d"
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
`, r.baseTemplate(data), data.RandomInteger, secretSettingName, secretSettingValue)
}

func (r WindowsWebAppResource) authV2Microsoft(data acceptance.TestData) string {
	secretSettingName := "MICROSOFT_PROVIDER_AUTHENTICATION_SECRET"
	secretSettingValue := "902D17F6-FD6B-4E44-BABB-58E788DCD907"

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%[2]d"
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
`, r.baseTemplate(data), data.RandomInteger, secretSettingName, secretSettingValue)
}

func (r WindowsWebAppResource) authV2Twitter(data acceptance.TestData) string {
	secretSettingName := "TWITTER_PROVIDER_AUTHENTICATION_SECRET"
	secretSettingValue := "902D17F6-FD6B-4E44-BABB-58E788DCD907"

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%[2]d"
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
`, r.baseTemplate(data), data.RandomInteger, secretSettingName, secretSettingValue)
}

func (r WindowsWebAppResource) authV2Multi(data acceptance.TestData) string {
	secretSettingValue := "902D17F6-FD6B-4E44-BABB-58E788DCD907"

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%[2]d"
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

func (r WindowsWebAppResource) completeAuthV2(data acceptance.TestData) string {
	secretSettingName := "APPLE_PROVIDER_AUTHENTICATION_SECRET"
	secretSettingValue := "902D17F6-FD6B-4E44-BABB-58E788DCD907"

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_web_app" "test" {
  name                = "acctestWA-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

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
      frequency_interval = 1
      frequency_unit     = "Day"
    }
  }

  logs {
    application_logs {
      file_system_level = "Warning"
      azure_blob_storage {
        level             = "Information"
        sas_url           = "http://x.com/"
        retention_in_days = 2
      }
    }

    http_logs {
      azure_blob_storage {
        sas_url           = "https://${azurerm_storage_account.test.name}.blob.core.windows.net/${azurerm_storage_container.test.name}${data.azurerm_storage_account_sas.test.sas}&sr=b"
        retention_in_days = 3
      }
    }
  }

  client_affinity_enabled            = true
  client_certificate_enabled         = true
  client_certificate_exclusion_paths = "/foo;/bar;/hello;/world"

  connection_string {
    name  = "First"
    value = "first-connection-string"
    type  = "Custom"
  }

  connection_string {
    name  = "Second"
    value = "some-postgresql-connection-string"
    type  = "PostgreSQL"
  }

  enabled    = false
  https_only = true

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  site_config {
    always_on        = true
    app_command_line = "/sbin/myserver -b 0.0.0.0"
    default_documents = [
      "first.html",
      "second.jsp",
      "third.aspx",
      "hostingstart.html",
    ]
    http2_enabled                     = true
    scm_use_main_ip_restriction       = true
    local_mysql_enabled               = true
    managed_pipeline_mode             = "Integrated"
    remote_debugging_enabled          = true
    remote_debugging_version          = "VS2022"
    use_32_bit_worker                 = false
    websockets_enabled                = true
    ftps_state                        = "FtpsOnly"
    health_check_path                 = "/health"
    health_check_eviction_time_in_min = 7
    worker_count                      = 1
    minimum_tls_version               = "1.1"
    scm_minimum_tls_version           = "1.1"
    cors {
      allowed_origins = [
        "http://www.contoso.com",
        "www.contoso.com",
      ]

      support_credentials = true
    }

    container_registry_use_managed_identity       = true
    container_registry_managed_identity_client_id = azurerm_user_assigned_identity.test.client_id

    virtual_application {
      virtual_path  = "/"
      physical_path = "site\\wwwroot"
      preload       = true

      virtual_directory {
        virtual_path  = "/stuff"
        physical_path = "site\\stuff"
      }
    }

    auto_heal_setting {
      trigger {
        status_code {
          status_code_range = "500"
          interval          = "00:01:00"
          count             = 10
        }
      }

      action {
        action_type                    = "Recycle"
        minimum_process_execution_time = "00:05:00"
      }
    }
  }

  sticky_settings {
    connection_string_names = ["First", "Third"]

    app_setting_names = [
      "foo",
      "APPLE_PROVIDER_AUTHENTICATION_SECRET",
      "FACEBOOK_PROVIDER_AUTHENTICATION_SECRET",
      "GITHUB_PROVIDER_AUTHENTICATION_SECRET",
      "GOOGLE_PROVIDER_AUTHENTICATION_SECRET",
      "MICROSOFT_PROVIDER_AUTHENTICATION_SECRET",
      "TWITTER_PROVIDER_AUTHENTICATION_SECRET",
    ]
  }

  storage_account {
    name         = "files"
    type         = "AzureFiles"
    account_name = azurerm_storage_account.test.name
    share_name   = azurerm_storage_share.test.name
    access_key   = azurerm_storage_account.test.primary_access_key
    mount_path   = "/mounts/files"
  }

  tags = {
    Environment = "AccTest"
    foo         = "bar"
  }
}
`, r.templateWithStorageAccount(data), data.RandomInteger, secretSettingName, secretSettingValue)
}

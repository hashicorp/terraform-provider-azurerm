package appservice_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccTestAccFunctionAppFlexConsumption_authV2Update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authV2Apple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
				check.That(data.ResourceName).Key("auth_settings_v2.0.auth_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("auth_settings_v2.0.apple_v2.#").HasValue("1"),
				check.That(data.ResourceName).Key("auth_settings_v2.0.apple_v2.0.client_id").IsNotEmpty(),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.authV2Facebook(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
				check.That(data.ResourceName).Key("auth_settings_v2.0.auth_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("auth_settings_v2.0.apple_v2.#").HasValue("0"),
				check.That(data.ResourceName).Key("auth_settings_v2.0.facebook_v2.#").HasValue("1"),
				check.That(data.ResourceName).Key("auth_settings_v2.0.facebook_v2.0.app_id").IsNotEmpty(),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.authV2Removed(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
				check.That(data.ResourceName).Key("auth_settings_v2.0.auth_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("auth_settings_v2.0.apple_v2.#").HasValue("0"),
				check.That(data.ResourceName).Key("auth_settings_v2.0.facebook_v2.#").HasValue("0"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.authV2Apple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
				check.That(data.ResourceName).Key("auth_settings_v2.0.auth_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("auth_settings_v2.0.apple_v2.#").HasValue("1"),
				check.That(data.ResourceName).Key("auth_settings_v2.0.apple_v2.0.client_id").IsNotEmpty(),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func (r FunctionAppFlexConsumptionResource) authV2Apple(data acceptance.TestData) string {
	secretSettingName := "APPLE_PROVIDER_AUTHENTICATION_SECRET"
	secretSettingValue := "902D17F6-FD6B-4E44-BABB-58E788DCD907"
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_container_type      = "blobContainer"
  storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type = "StorageAccountConnectionString"
  storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                = "node"
  runtime_version             = "20"
  maximum_instance_count      = 100
  instance_memory_in_mb       = 2048

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
`, r.template(data), data.RandomInteger, secretSettingName, secretSettingValue)
}

func (r FunctionAppFlexConsumptionResource) authV2Facebook(data acceptance.TestData) string {
	secretSettingName := "FACEBOOK_PROVIDER_AUTHENTICATION_SECRET"
	secretSettingValue := "902D17F6-FD6B-4E44-BABB-58E788DCD907"
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_container_type      = "blobContainer"
  storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type = "StorageAccountConnectionString"
  storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                = "node"
  runtime_version             = "20"
  maximum_instance_count      = 100
  instance_memory_in_mb       = 2048

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
`, r.template(data), data.RandomInteger, secretSettingName, secretSettingValue)
}

func (r FunctionAppFlexConsumptionResource) authV2Removed(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_container_type      = "blobContainer"
  storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type = "StorageAccountConnectionString"
  storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                = "node"
  runtime_version             = "20"
  maximum_instance_count      = 100
  instance_memory_in_mb       = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger)
}

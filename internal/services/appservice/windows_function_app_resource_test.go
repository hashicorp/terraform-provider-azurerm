package appservice_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WindowsFunctionAppResource struct{}

// Plan types
func TestAccWindowsFunctionApp_basicBasicPlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, SkuBasicPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_basicConsumptionPlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, SkuConsumptionPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_basicElasticPremiumPlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, SkuElasticPremiumPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_basicPremiumAppServicePlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, SkuPremiumPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_basicStandardPlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
	})
}

// App Settings by Plan Type

func TestAccWindowsFunctionApp_withAppSettingsBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appSettings(data, SkuBasicPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
				check.That(data.ResourceName).Key("app_settings.%").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_withAppSettingsConsumption(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appSettings(data, SkuConsumptionPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
				check.That(data.ResourceName).Key("app_settings.%").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_withAppSettingsElasticPremiumPlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appSettings(data, SkuElasticPremiumPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
				check.That(data.ResourceName).Key("app_settings.%").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_withCustomContentShareElasticPremiumPlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appSettingsCustomContentShare(data, SkuElasticPremiumPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
				check.That(data.ResourceName).Key("app_settings.%").HasValue("3"),
				check.That(data.ResourceName).Key("app_settings.WEBSITE_CONTENTSHARE").HasValue("test-acc-custom-content-share"),
			),
		},
		data.ImportStep("app_settings.WEBSITE_CONTENTSHARE", "app_settings.%"),
	})
}

func TestAccWindowsFunctionApp_withAppSettingsPremiumPlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appSettings(data, SkuPremiumPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
				check.That(data.ResourceName).Key("app_settings.%").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_withAppSettingsStandardPlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appSettings(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
				check.That(data.ResourceName).Key("app_settings.%").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_withAppSettingsUserSettingUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appSettingsUserSettings(data, SkuConsumptionPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
				check.That(data.ResourceName).Key("app_settings.%").DoesNotExist(),
			),
		},
		data.ImportStep(),
		{
			Config: r.appSettingsUserSettingsUpdate(data, SkuConsumptionPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
				check.That(data.ResourceName).Key("app_settings.%").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_addAppSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appStackNode(data, SkuConsumptionPlan, "~14"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
				check.That(data.ResourceName).Key("app_settings.%").DoesNotExist(),
			),
		},
		data.ImportStep(),
		{
			Config: r.appSettingsAdded(data, SkuConsumptionPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
				check.That(data.ResourceName).Key("app_settings.%").HasValue("2"),
			),
		},
		data.ImportStep(),
		{
			Config: r.appStackNode(data, SkuConsumptionPlan, "~14"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
				check.That(data.ResourceName).Key("app_settings.%").DoesNotExist(),
			),
		},
		data.ImportStep(),
	})
}

// Sticky Settings

func TestAccWindowsFunctionApp_stickySettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.stickySettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("app_settings.foo").HasValue("bar"),
				check.That(data.ResourceName).Key("sticky_settings.0.app_setting_names.#").HasValue("3"),
				check.That(data.ResourceName).Key("sticky_settings.0.app_setting_names.0").HasValue("foo"),
				check.That(data.ResourceName).Key("sticky_settings.0.connection_string_names.#").HasValue("3"),
				check.That(data.ResourceName).Key("sticky_settings.0.connection_string_names.0").HasValue("First"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_stickySettingsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("app_settings").DoesNotExist(),
				check.That(data.ResourceName).Key("sticky_settings").DoesNotExist(),
			),
		},
		data.ImportStep(),
		{
			Config: r.stickySettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("app_settings.foo").HasValue("bar"),
				check.That(data.ResourceName).Key("sticky_settings.0.app_setting_names.#").HasValue("3"),
				check.That(data.ResourceName).Key("sticky_settings.0.app_setting_names.0").HasValue("foo"),
				check.That(data.ResourceName).Key("sticky_settings.0.connection_string_names.#").HasValue("3"),
				check.That(data.ResourceName).Key("sticky_settings.0.connection_string_names.0").HasValue("First"),
			),
		},
		data.ImportStep(),
		{
			Config: r.stickySettingsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("app_settings.foo").HasValue("bar"),
				check.That(data.ResourceName).Key("sticky_settings.0.app_setting_names.#").HasValue("3"),
				check.That(data.ResourceName).Key("sticky_settings.0.app_setting_names.0").HasValue("foo"),
				check.That(data.ResourceName).Key("sticky_settings.0.connection_string_names.#").HasValue("3"),
				check.That(data.ResourceName).Key("sticky_settings.0.connection_string_names.0").HasValue("First"),
			),
		},
		data.ImportStep(),
		{
			Config: r.stickySettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("app_settings.foo").HasValue("bar"),
				check.That(data.ResourceName).Key("sticky_settings.0.app_setting_names.#").HasValue("3"),
				check.That(data.ResourceName).Key("sticky_settings.0.app_setting_names.0").HasValue("foo"),
				check.That(data.ResourceName).Key("sticky_settings.0.connection_string_names.#").HasValue("3"),
				check.That(data.ResourceName).Key("sticky_settings.0.connection_string_names.0").HasValue("First"),
			),
		},
		data.ImportStep(),
		{
			Config: r.stickySettingsRemoved(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("app_settings.foo").HasValue("bar"),
				check.That(data.ResourceName).Key("sticky_settings").DoesNotExist(),
			),
		},
		data.ImportStep(),
	})
}

// backup by plan type

func TestAccWindowsFunctionApp_withBackupElasticPremiumPlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.backup(data, SkuElasticPremiumPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_withBackupPremiumPlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.backup(data, SkuPremiumPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_withBackupStandardPlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.backup(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
	})
}

// Completes by plan type

func TestAccWindowsFunctionApp_consumptionComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.consumptionComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_consumptionCompleteUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, SkuConsumptionPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
		{
			Config: r.consumptionComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, SkuConsumptionPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
	})
}

func TestAccWindowsFunctionApp_elasticPremiumComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.elasticComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.elastic_instance_minimum").HasValue("5"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_standardComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.standardComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

// Individual Settings / Blocks

func TestAccWindowsFunctionApp_withAuthSettingsConsumption(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withAuthSettings(data, SkuConsumptionPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_withAuthSettingsStandard(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withAuthSettings(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_builtInLogging(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.builtInLogging(data, SkuStandardPlan, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_withConnectionStrings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.connectionStrings(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_withConnectionStringsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
		{
			Config: r.connectionStrings(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
		{
			Config: r.connectionStringsUpdate(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_dailyTimeQuotaConsumptionPlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dailyTimeLimitQuota(data, SkuConsumptionPlan, 1000),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_dailyTimeQuotaElasticPremiumPlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dailyTimeLimitQuota(data, SkuElasticPremiumPlan, 2000),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_healthCheckPath(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.healthCheckPath(data, "S1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_healthCheckPathWithEviction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.healthCheckPathWithEviction(data, "S1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_healthCheckPathWithEvictionUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "S1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.healthCheckPathWithEviction(data, "S1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, "S1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_appServiceLogging(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appServiceLogs(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_appServiceLoggingUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
		{
			Config: r.appServiceLogs(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
		{
			Config: r.appServiceLogsWithRetention(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.identitySystemAssigned(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.identityUserAssigned(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.identitySystemAssignedUserAssigned(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_identityKeyVaultIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.identityUserAssignedKeyVaultIdentity(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

// App Stacks

func TestAccWindowsFunctionApp_appStackDotNet31(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appStackDotNet(data, SkuBasicPlan, "3.1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_appStackDotNet6(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appStackDotNet(data, SkuBasicPlan, "6"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_appStackDotNet6Isolated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appStackDotNetIsolated(data, SkuBasicPlan, "6"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
				check.That(data.ResourceName).Key("site_config.0.windows_fx_version").HasValue("DotNet-Isolated|6"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_appStackNode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appStackNode(data, SkuConsumptionPlan, "~14"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_appStackNodeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appStackNode(data, SkuBasicPlan, "~12"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
		{
			Config: r.appStackNode(data, SkuBasicPlan, "~16"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
		{
			Config: r.appStackNode(data, SkuBasicPlan, "~14"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_appStackJava(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appStackJava(data, SkuBasicPlan, "11"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_appStackJavaUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appStackJava(data, SkuBasicPlan, "8"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
		{
			Config: r.appStackJava(data, SkuBasicPlan, "11"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_appStackPowerShellCore(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appStackPowerShellCore(data, SkuBasicPlan, "7"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
	})
}

// Others

func TestAccWindowsFunctionApp_updateServicePlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
		{
			Config: r.servicePlanUpdate(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_updateStorageAccount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateStorageAccount(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_msiStorageAccount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.msiStorageAccount(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_msiStorageAccountUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.msiStorageAccount(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
		{
			Config: r.msiStorageAccountUpdate(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_storageAccountKeyVaultSecret(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageAccountKVSecret(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionApp_storageAccountKeyVaultSecretVersionless(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app", "test")
	r := WindowsFunctionAppResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageAccountKVSecretVersionless(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
	})
}

// Exists

func (r WindowsFunctionAppResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.FunctionAppID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.AppService.WebAppsClient.Get(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Windows Web App %s: %+v", id, err)
	}
	if utils.ResponseWasNotFound(resp.Response) {
		return utils.Bool(false), nil
	}
	return utils.Bool(true), nil
}

// Configs

func (r WindowsFunctionAppResource) basic(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {}
}
`, r.template(data, planSku), data.RandomInteger)
}

func (r WindowsFunctionAppResource) appSettings(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  app_settings = {
    foo    = "bar"
    secret = "sauce"
  }

  site_config {}
}
`, r.template(data, planSku), data.RandomInteger)
}

func (r WindowsFunctionAppResource) appSettingsUserSettings(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  app_settings = {}

  site_config {}
}
`, r.template(data, planSku), data.RandomInteger)
}

func (r WindowsFunctionAppResource) appSettingsUserSettingsUpdate(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  app_settings = {
    foo    = "bar"
    secret = "sauce"
  }

  site_config {}
}
`, r.template(data, planSku), data.RandomInteger)
}

func (r WindowsFunctionAppResource) appSettingsAdded(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  app_settings = {
    foo    = "bar"
    secret = "sauce"
  }

  site_config {
    application_stack {
      node_version = "~14"
    }
  }
}
`, r.template(data, planSku), data.RandomInteger)
}

func (r WindowsFunctionAppResource) appSettingsCustomContentShare(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  app_settings = {
    foo                  = "bar"
    secret               = "sauce"
    WEBSITE_CONTENTSHARE = "test-acc-custom-content-share"
  }

  site_config {}
}
`, r.template(data, planSku), data.RandomInteger)
}

func (r WindowsFunctionAppResource) stickySettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
  site_config {}

  app_settings = {
    foo                                     = "bar"
    secret                                  = "sauce"
    third                                   = "degree"
    "Special chars: !@#$%%^&*()_+-=' \";/?" = "Supported by the Azure portal"
  }

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

  connection_string {
    name  = "Third"
    value = "some-postgresql-connection-string"
    type  = "PostgreSQL"
  }

  connection_string {
    name  = "Special chars: !@#$%%^&*()_+-=' \";/?"
    value = "characters-supported-by-the-Azure-portal"
    type  = "Custom"
  }

  sticky_settings {
    app_setting_names       = ["foo", "secret", "Special chars: !@#$%%^&*()_+-=' \";/?"]
    connection_string_names = ["First", "Third", "Special chars: !@#$%%^&*()_+-=' \";/?"]
  }
}
`, r.template(data, SkuStandardPlan), data.RandomInteger)
}

func (r WindowsFunctionAppResource) stickySettingsRemoved(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {}

  app_settings = {
    foo    = "bar"
    secret = "sauce"
    third  = "degree"
  }

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

  connection_string {
    name  = "Third"
    value = "some-postgresql-connection-string"
    type  = "PostgreSQL"
  }
}
`, r.template(data, SkuStandardPlan), data.RandomInteger)
}

func (r WindowsFunctionAppResource) stickySettingsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {}

  app_settings = {
    foo    = "bar"
    secret = "sauce"
    third  = "degree"
  }

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

  connection_string {
    name  = "Third"
    value = "some-postgresql-connection-string"
    type  = "PostgreSQL"
  }

  sticky_settings {
    app_setting_names       = ["foo", "secret", "third"]
    connection_string_names = ["First", "Second", "Third"]
  }
}
`, r.template(data, SkuStandardPlan), data.RandomInteger)
}

func (r WindowsFunctionAppResource) backup(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  backup {
    name                = "acctest"
    storage_account_url = "https://${azurerm_storage_account.test.name}.blob.core.windows.net/${azurerm_storage_container.test.name}${data.azurerm_storage_account_sas.test.sas}&sr=b"
    schedule {
      frequency_interval = 7
      frequency_unit     = "Day"
    }
  }

  site_config {}
}
`, r.storageContainerTemplate(data, planSku), data.RandomInteger)
}

func (r WindowsFunctionAppResource) consumptionComplete(data acceptance.TestData) string {
	planSku := "Y1"
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

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  app_settings = {
    foo    = "bar"
    secret = "sauce"
  }

  auth_settings {
    enabled = true
    issuer  = "https://sts.windows.net/%[3]s"

    additional_login_parameters = {
      test_key = "test_value"
    }

    active_directory {
      client_id     = "aadclientid"
      client_secret = "aadsecret"

      allowed_audiences = [
        "activedirectorytokenaudiences",
      ]
    }

    facebook {
      app_id     = "facebookappid"
      app_secret = "facebookappsecret"

      oauth_scopes = [
        "facebookscope",
      ]
    }
  }

  builtin_logging_enabled    = false
  client_certificate_enabled = true
  client_certificate_mode    = "Required"

  connection_string {
    name  = "Second"
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
    app_command_line   = "whoami"
    api_definition_url = "https://example.com/azure_function_app_def.json"
    app_scale_limit    = 3
    // api_management_api_id = ""  // TODO
    application_insights_key               = azurerm_application_insights.test.instrumentation_key
    application_insights_connection_string = azurerm_application_insights.test.connection_string

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

    ip_restriction {
      service_tag = "ActionGroup"
      name        = "test-servicetag-restriction"
      priority    = 125
      action      = "Allow"
    }

    load_balancing_mode      = "LeastResponseTime"
    remote_debugging_enabled = true
    remote_debugging_version = "VS2019"

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

    use_32_bit_worker  = true
    websockets_enabled = true
    ftps_state         = "FtpsOnly"
    health_check_path  = "/health-check"

    application_stack {
      powershell_core_version = "7"
    }

    minimum_tls_version     = "1.1"
    scm_minimum_tls_version = "1.1"

    cors {
      allowed_origins = [
        "https://www.contoso.com",
        "www.contoso.com",
      ]

      support_credentials = true
    }
  }

  tags = {
    terraform = "true"
    Env       = "AccTest"
  }
}
`, r.template(data, planSku), data.RandomInteger, data.Client().TenantID)
}

func (r WindowsFunctionAppResource) standardComplete(data acceptance.TestData) string {
	planSku := "S1"
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

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  app_settings = {
    foo    = "bar"
    secret = "sauce"
  }

  auth_settings {
    enabled = true
    issuer  = "https://sts.windows.net/%[3]s"

    additional_login_parameters = {
      test_key = "test_value"
    }

    active_directory {
      client_id     = "aadclientid"
      client_secret = "aadsecret"

      allowed_audiences = [
        "activedirectorytokenaudiences",
      ]
    }

    facebook {
      app_id     = "facebookappid"
      app_secret = "facebookappsecret"

      oauth_scopes = [
        "facebookscope",
      ]
    }
  }

  backup {
    name                = "acctest"
    storage_account_url = "https://${azurerm_storage_account.test.name}.blob.core.windows.net/${azurerm_storage_container.test.name}${data.azurerm_storage_account_sas.test.sas}&sr=b"
    schedule {
      frequency_interval = 7
      frequency_unit     = "Day"
    }
  }

  builtin_logging_enabled    = false
  client_certificate_enabled = true
  client_certificate_mode    = "OptionalInteractiveUser"

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
    always_on          = true
    app_command_line   = "whoami"
    api_definition_url = "https://example.com/azure_function_app_def.json"
    // api_management_api_id = ""  // TODO
    application_insights_key               = azurerm_application_insights.test.instrumentation_key
    application_insights_connection_string = azurerm_application_insights.test.connection_string

    application_stack {
      powershell_core_version = "7"
    }

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
    remote_debugging_version  = "VS2017"

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

    use_32_bit_worker  = true
    websockets_enabled = true
    ftps_state         = "FtpsOnly"
    health_check_path  = "/health-check"
    worker_count       = 3

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
`, r.storageContainerTemplate(data, planSku), data.RandomInteger, data.Client().TenantID)
}

func (r WindowsFunctionAppResource) elasticComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}


resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  app_settings = {
    foo    = "bar"
    secret = "sauce"
  }

  backup {
    name                = "acctest"
    storage_account_url = "https://${azurerm_storage_account.test.name}.blob.core.windows.net/${azurerm_storage_container.test.name}${data.azurerm_storage_account_sas.test.sas}&sr=b"
    schedule {
      frequency_interval = 7
      frequency_unit     = "Day"
    }
  }

  connection_string {
    name  = "Example"
    value = "some-postgresql-connection-string"
    type  = "PostgreSQL"
  }

  site_config {
    app_command_line                       = "whoami"
    api_definition_url                     = "https://example.com/azure_function_app_def.json"
    application_insights_key               = azurerm_application_insights.test.instrumentation_key
    application_insights_connection_string = azurerm_application_insights.test.connection_string

    application_stack {
      powershell_core_version = "7"
    }

    elastic_instance_minimum = 5

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
    remote_debugging_version  = "VS2017"

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

    use_32_bit_worker  = true
    websockets_enabled = true
    ftps_state         = "FtpsOnly"
    health_check_path  = "/health-check"
    worker_count       = 3

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
}
`, r.storageContainerTemplate(data, SkuElasticPremiumPlan), data.RandomInteger)
}

func (r WindowsFunctionAppResource) withAuthSettings(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  auth_settings {
    enabled                       = true
    issuer                        = "https://sts.windows.net/%s"
    runtime_version               = "1.0"
    unauthenticated_client_action = "RedirectToLoginPage"
    token_refresh_extension_hours = 75
    token_store_enabled           = true

    additional_login_parameters = {
      test_key = "test_value"
    }

    allowed_external_redirect_urls = [
      "https://terra.form",
    ]

    active_directory {
      client_id     = "aadclientid"
      client_secret = "aadsecret"

      allowed_audiences = [
        "activedirectorytokenaudiences",
      ]
    }
  }

  site_config {}
}
`, r.template(data, planSku), data.RandomInteger, data.RandomString)
}

func (r WindowsFunctionAppResource) builtInLogging(data acceptance.TestData, planSku string, builtInLogging bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  builtin_logging_enabled = %t

  site_config {}
}
`, r.template(data, planSku), data.RandomInteger, builtInLogging)
}

func (r WindowsFunctionAppResource) connectionStrings(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  connection_string {
    name  = "Example"
    value = "some-postgresql-connection-string"
    type  = "PostgreSQL"
  }

  site_config {}
}
`, r.template(data, planSku), data.RandomInteger)
}

func (r WindowsFunctionAppResource) connectionStringsUpdate(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  connection_string {
    name  = "Example"
    value = "some-postgresql-connection-string"
    type  = "PostgreSQL"
  }

  connection_string {
    name  = "AnotherExample"
    value = "some-other-connection-string"
    type  = "Custom"
  }

  site_config {}
}
`, r.template(data, planSku), data.RandomInteger)
}

func (r WindowsFunctionAppResource) dailyTimeLimitQuota(data acceptance.TestData, planSku string, quota int) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  daily_memory_time_quota = %d

  site_config {}
}
`, r.template(data, planSku), data.RandomInteger, quota)
}

func (r WindowsFunctionAppResource) healthCheckPath(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    health_check_path = "/health"
  }
}
`, r.template(data, planSku), data.RandomInteger)
}

func (r WindowsFunctionAppResource) healthCheckPathWithEviction(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    health_check_path                 = "/health"
    health_check_eviction_time_in_min = 3
  }
}
`, r.template(data, planSku), data.RandomInteger)
}

func (r WindowsFunctionAppResource) appServiceLogs(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    app_service_logs {
      disk_quota_mb = 25
    }
  }
}
`, r.template(data, planSku), data.RandomInteger)
}

func (r WindowsFunctionAppResource) appServiceLogsWithRetention(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    app_service_logs {
      disk_quota_mb         = 25
      retention_period_days = 7
    }
  }
}
`, r.template(data, planSku), data.RandomInteger)
}

func (r WindowsFunctionAppResource) appStackDotNet(data acceptance.TestData, planSku string, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    application_stack {
      dotnet_version = "%s"
    }
  }
}
`, r.template(data, planSku), data.RandomInteger, version)
}

func (r WindowsFunctionAppResource) appStackDotNetIsolated(data acceptance.TestData, planSku string, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    application_stack {
      dotnet_version              = "%s"
      use_dotnet_isolated_runtime = true
    }
  }
}
`, r.template(data, planSku), data.RandomInteger, version)
}

func (r WindowsFunctionAppResource) appStackNode(data acceptance.TestData, planSku string, nodeVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    application_stack {
      node_version = "%s"
    }
  }
}
`, r.template(data, planSku), data.RandomInteger, nodeVersion)
}

func (r WindowsFunctionAppResource) appStackJava(data acceptance.TestData, planSku string, javaVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    application_stack {
      java_version = "%s"
    }
  }
}
`, r.template(data, planSku), data.RandomInteger, javaVersion)
}

func (r WindowsFunctionAppResource) appStackPowerShellCore(data acceptance.TestData, planSku string, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    application_stack {
      powershell_core_version = "%s"
    }
  }
}
`, r.template(data, planSku), data.RandomInteger, version)
}

func (r WindowsFunctionAppResource) servicePlanUpdate(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.update.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {}

  depends_on = [azurerm_service_plan.update]
}
`, r.templateServicePlanUpdate(data, planSku), data.RandomInteger)
}

func (r WindowsFunctionAppResource) updateStorageAccount(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-LFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.update.name
  storage_account_access_key = azurerm_storage_account.update.primary_access_key

  site_config {}
}
`, r.templateExtraStorageAccount(data, planSku), data.RandomInteger)
}

func (r WindowsFunctionAppResource) identitySystemAssigned(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {}

  identity {
    type = "SystemAssigned"
  }
}
`, r.identityTemplate(data, planSku), data.RandomInteger)
}

func (r WindowsFunctionAppResource) identitySystemAssignedUserAssigned(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {}

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, r.identityTemplate(data, planSku), data.RandomInteger)
}

func (r WindowsFunctionAppResource) identityUserAssigned(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {}

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, r.identityTemplate(data, planSku), data.RandomInteger)
}

func (r WindowsFunctionAppResource) identityUserAssignedKeyVaultIdentity(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_user_assigned_identity" "kv" {
  name                = "acctest-kv-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {}

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id, azurerm_user_assigned_identity.kv.id]
  }

  key_vault_reference_identity_id = azurerm_user_assigned_identity.kv.id
}
`, r.identityTemplate(data, planSku), data.RandomInteger)
}

func (r WindowsFunctionAppResource) msiStorageAccount(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s


resource "azurerm_role_assignment" "func_app_access_to_storage" {
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Storage Blob Data Owner"
  principal_id         = azurerm_windows_function_app.test.identity[0].principal_id
}

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name          = azurerm_storage_account.test.name
  storage_uses_managed_identity = true

  identity {
    type = "SystemAssigned"
  }

  site_config {}
}
`, r.template(data, planSku), data.RandomInteger)
}

func (r WindowsFunctionAppResource) msiStorageAccountUpdate(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s


resource "azurerm_role_assignment" "func_app_access_to_storage" {
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Storage Blob Data Owner"
  principal_id         = azurerm_windows_function_app.test.identity[0].principal_id
}

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name          = azurerm_storage_account.test.name
  storage_uses_managed_identity = true

  identity {
    type = "SystemAssigned"
  }

  site_config {
    always_on = true
  }
}
`, r.template(data, planSku), data.RandomInteger)
}

func (r WindowsFunctionAppResource) storageAccountKVSecret(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy    = true
      recover_soft_deleted_key_vaults = true
    }
  }
}

%[1]s

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                       = "acctestkv-%[2]s"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Get",
    ]

    secret_permissions = [
      "Get",
      "Delete",
      "List",
      "Purge",
      "Recover",
      "Set",
    ]
  }

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = azurerm_user_assigned_identity.test.principal_id

    secret_permissions = [
      "Get",
      "List",
    ]
  }

  tags = {
    environment = "AccTest"
  }
}

resource "azurerm_key_vault_secret" "test" {
  name         = "secret-%[2]s"
  value        = "DefaultEndpointsProtocol=https;AccountName=${azurerm_storage_account.test.name};AccountKey=${azurerm_storage_account.test.primary_access_key};EndpointSuffix=core.windows.net"
  key_vault_id = azurerm_key_vault.test.id
}

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%[3]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  key_vault_reference_identity_id = azurerm_user_assigned_identity.test.id
  storage_key_vault_secret_id     = azurerm_key_vault_secret.test.id

  site_config {}

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, r.identityTemplate(data, planSku), data.RandomString, data.RandomInteger)
}

func (r WindowsFunctionAppResource) storageAccountKVSecretVersionless(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy    = true
      recover_soft_deleted_key_vaults = true
    }
  }
}

%[1]s

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                       = "acctestkv-%[2]s"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Get",
    ]

    secret_permissions = [
      "Get",
      "Delete",
      "List",
      "Purge",
      "Recover",
      "Set",
    ]
  }

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = azurerm_user_assigned_identity.test.principal_id

    secret_permissions = [
      "Get",
      "List",
    ]
  }

  tags = {
    environment = "AccTest"
  }
}

resource "azurerm_key_vault_secret" "test" {
  name         = "secret-%[2]s"
  value        = "DefaultEndpointsProtocol=https;AccountName=${azurerm_storage_account.test.name};AccountKey=${azurerm_storage_account.test.primary_access_key};EndpointSuffix=core.windows.net"
  key_vault_id = azurerm_key_vault.test.id
}

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%[3]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  key_vault_reference_identity_id = azurerm_user_assigned_identity.test.id
  storage_key_vault_secret_id     = azurerm_key_vault_secret.test.versionless_id

  site_config {}

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, r.identityTemplate(data, planSku), data.RandomString, data.RandomInteger)
}

// Config Templates

func (WindowsFunctionAppResource) template(data acceptance.TestData, planSku string) string {
	var additionalConfig string
	if strings.EqualFold(planSku, "EP1") {
		additionalConfig = "maximum_elastic_worker_count = 5"
	}
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-LFA-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  os_type             = "Windows"
  sku_name            = "%s"
  %s
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, planSku, additionalConfig)
}

func (r WindowsFunctionAppResource) storageContainerTemplate(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "test"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

data "azurerm_storage_account_sas" "test" {
  connection_string = azurerm_storage_account.test.primary_connection_string
  https_only        = true

  resource_types {
    service   = false
    container = false
    object    = true
  }

  services {
    blob  = true
    queue = false
    table = false
    file  = false
  }

  start  = "2021-04-01"
  expiry = "2024-03-30"

  permissions {
    read    = false
    write   = true
    delete  = false
    list    = false
    add     = false
    create  = false
    update  = false
    process = false
    tag     = false
    filter  = false
  }
}

`, r.template(data, planSku))
}

func (r WindowsFunctionAppResource) identityTemplate(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acct-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

`, r.template(data, planSku), data.RandomInteger)
}

func (r WindowsFunctionAppResource) templateServicePlanUpdate(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_service_plan" "update" {
  name                = "acctestASP2-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  os_type             = "Windows"
  sku_name            = "%s"
}
`, r.template(data, planSku), data.RandomInteger, planSku)
}

func (WindowsFunctionAppResource) templateExtraStorageAccount(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-WFA-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_account" "update" {
  name                     = "acctestsa2%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  os_type             = "Windows"
  sku_name            = "%[4]s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, planSku)
}

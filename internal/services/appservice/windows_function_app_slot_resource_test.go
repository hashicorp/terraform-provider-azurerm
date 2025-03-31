// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WindowsFunctionAppSlotResource struct{}

// Plan types

func TestAccWindowsFunctionAppSlot_basicConsumptionPlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, SkuConsumptionPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_basicElasticPremiumPlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, SkuElasticPremiumPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_basicPremiumAppServicePlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, SkuPremiumPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_basicStandardPlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_withAppSettingsConsumption(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appSettings(data, SkuConsumptionPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
				check.That(data.ResourceName).Key("app_settings.%").HasValue("2"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

// App Settings by Plan Type

func TestAccWindowsFunctionAppSlot_withAppSettingsElasticPremiumPlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appSettings(data, SkuElasticPremiumPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
				check.That(data.ResourceName).Key("app_settings.%").HasValue("2"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_withAppSettingsPremiumPlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appSettings(data, SkuPremiumPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
				check.That(data.ResourceName).Key("app_settings.%").HasValue("2"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_withAppSettingsStandardPlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appSettings(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
				check.That(data.ResourceName).Key("app_settings.%").HasValue("2"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_withCustomContentShareElasticPremiumPlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appSettingsCustomContentShare(data, SkuElasticPremiumPlan, "test-acc-custom-content-share"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
				check.That(data.ResourceName).Key("app_settings.%").HasValue("3"),
				check.That(data.ResourceName).Key("app_settings.WEBSITE_CONTENTSHARE").HasValue("test-acc-custom-content-share"),
			),
		},
		data.ImportStep("app_settings.WEBSITE_CONTENTSHARE", "app_settings.%", "site_credential.0.password"),
		{
			Config: r.appSettingsCustomContentShare(data, SkuElasticPremiumPlan, "test-acc-custom-content-updated"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
				check.That(data.ResourceName).Key("app_settings.%").HasValue("3"),
				check.That(data.ResourceName).Key("app_settings.WEBSITE_CONTENTSHARE").HasValue("test-acc-custom-content-updated"),
			),
		},
		data.ImportStep("app_settings.WEBSITE_CONTENTSHARE", "app_settings.%", "site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_separateStandardPlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.separatePlan(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_separateStandardPlanUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.separatePlan(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.separatePlanUpdate(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_withAppSettingsUserSettingUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appSettingsUserSettings(data, SkuConsumptionPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
				check.That(data.ResourceName).Key("app_settings.%").DoesNotExist(),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.appSettingsUserSettingsUpdate(data, SkuConsumptionPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
				check.That(data.ResourceName).Key("app_settings.%").HasValue("2"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

// backup by plan type

func TestAccWindowsFunctionAppSlot_withBackupElasticPremiumPlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.backup(data, SkuElasticPremiumPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_withBackupPremiumPlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.backup(data, SkuPremiumPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_withBackupStandardPlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.backup(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_withBackupVnetIntegration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.backupVnetIntegration(data, SkuStandardPlan, "true"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.backupVnetIntegration(data, SkuStandardPlan, "false"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

// Completes by plan type

func TestAccWindowsFunctionAppSlot_consumptionComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.consumptionComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_consumptionCompleteUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, SkuConsumptionPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.consumptionComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.basic(data, SkuConsumptionPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
	})
}

func TestAccWindowsFunctionAppSlot_elasticPremiumComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.elasticComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.elastic_instance_minimum").HasValue("2"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_standardComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.standardComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

// Individual Settings / Blocks

func TestAccWindowsFunctionAppSlot_withAuthSettingsConsumption(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withAuthSettings(data, SkuConsumptionPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_withAuthSettingsStandard(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withAuthSettings(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_builtInLogging(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.builtInLogging(data, SkuStandardPlan, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_withConnectionStrings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.connectionStrings(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_withConnectionStringsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.connectionStrings(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.connectionStringsUpdate(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.basic(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_dailyTimeQuotaConsumptionPlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dailyTimeLimitQuota(data, SkuConsumptionPlan, 1000),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_dailyTimeQuotaElasticPremiumPlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dailyTimeLimitQuota(data, SkuElasticPremiumPlan, 2000),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_healthCheckPath(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.healthCheckPath(data, "S1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_basicRuntimeCheck(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.runtimeScaleCheck(data, SkuElasticPremiumPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.runtime_scale_monitoring_enabled").HasValue("true"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_healthCheckPathWithEviction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.healthCheckPathWithEviction(data, "S1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_healthCheckPathWithEvictionUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, "S1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.healthCheckPathWithEviction(data, "S1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.basic(data, "S1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_appServiceLogging(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appServiceLogs(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_appServiceLoggingUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.appServiceLogs(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.basic(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_withIPRestrictions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withIPRestrictions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_withIPRestrictionsDescription(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withIPRestrictions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.withIPRestrictionsDescription(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.withIPRestrictions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_withIPRestrictionsDefaultAction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withIPRestrictions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_withIPRestrictionsDefaultActionUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withIPRestrictions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.withIPRestrictionsDefaultActionDeny(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.withIPRestrictions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

// App Stacks

func TestAccWindowsFunctionAppSlot_appStackDotNet31(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appStackDotNet(data, SkuStandardPlan, "v3.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_appStackDotNet6(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appStackDotNet(data, SkuStandardPlan, "v6.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_appStackDotNet6Isolated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appStackDotNetIsolated(data, SkuStandardPlan, "v6.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_appStackDotNet9Isolated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appStackDotNetIsolated(data, SkuStandardPlan, "v9.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_appStackNode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appStackNode(data, SkuStandardPlan, "~18"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_appStackNode20(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appStackNode(data, SkuStandardPlan, "~20"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_appStackNodeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appStackNode(data, SkuStandardPlan, "~16"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.appStackNode(data, SkuStandardPlan, "~18"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.appStackNode(data, SkuStandardPlan, "~20"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.appStackNode(data, SkuStandardPlan, "~22"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWindowsFunctionAppSlot_appStackJava(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appStackJava(data, SkuStandardPlan, "11"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_appStackJavaUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appStackJava(data, SkuStandardPlan, "1.8"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.appStackJava(data, SkuStandardPlan, "17"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_appStackPowerShellCore7(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appStackPowerShellCore(data, SkuStandardPlan, "7"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_appStackPowerShellCore72(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appStackPowerShellCore(data, SkuStandardPlan, "7.2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

// Others

func TestAccWindowsFunctionAppSlot_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.identitySystemAssigned(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.identityUserAssigned(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.identitySystemAssignedUserAssigned(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.basic(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_identityKeyVaultIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.identityUserAssignedKeyVaultIdentity(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_updateServicePlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.servicePlanUpdate(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_updateStorageAccount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.updateStorageAccount(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_msiStorageAccountConsumption(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.msiStorageAccount(data, SkuConsumptionPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_msiStorageAccountElastic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.msiStorageAccount(data, SkuElasticPremiumPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_msiStorageAccountStandard(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.msiStorageAccount(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_storageAccountKeyVaultSecret(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageAccountKVSecret(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_storageAccountKeyVaultSecretVersionless(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageAccountKVSecretVersionless(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_vNetIntegration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.vNetIntegration_subnet1WithVnetProperties(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("virtual_network_subnet_id").MatchesOtherKey(
					check.That("azurerm_subnet.test1").Key("id"),
				),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_vNetIntegrationUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.vNetIntegration_basic(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.vNetIntegration_subnet1(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("virtual_network_subnet_id").MatchesOtherKey(
					check.That("azurerm_subnet.test1").Key("id"),
				),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.vNetIntegration_subnet2(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("virtual_network_subnet_id").MatchesOtherKey(
					check.That("azurerm_subnet.test2").Key("id"),
				),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.vNetIntegration_basic(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlotASEv3_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withASEV3(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_withStorageAccountBlock(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withStorageAccountSingle(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_withStorageAccountBlocks(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withStorageAccountMultiple(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_withStorageAccountBlocksUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.withStorageAccountSingle(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.withStorageAccountMultiple(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.withStorageAccountSingle(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.basic(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_publicNetworkAccessDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.publicNetworkAccessDisabled(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("false"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_publicNetworkAccessUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("true"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.publicNetworkAccessDisabled(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("false"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.basic(data, SkuStandardPlan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("true"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccWindowsFunctionAppSlot_basicWithTlsOnePointThree(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_function_app_slot", "test")
	r := WindowsFunctionAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withTlsVersion(data, SkuConsumptionPlan, "1.3"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

// Exists

func (r WindowsFunctionAppSlotResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := webapps.ParseSlotID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.AppService.WebAppsClient.GetSlot(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Windows %s: %+v", id, err)
	}
	if response.WasNotFound(resp.HttpResponse) {
		return utils.Bool(false), nil
	}
	return utils.Bool(true), nil
}

// Configs

func (r WindowsFunctionAppSlotResource) basic(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%d"
  function_app_id            = azurerm_windows_function_app.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {}
}
`, r.template(data, planSku), data.RandomInteger)
}

func (r WindowsFunctionAppSlotResource) withIPRestrictions(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%d"
  function_app_id            = azurerm_windows_function_app.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    ip_restriction {
      ip_address = "13.107.6.152/31,13.107.128.0/22"
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
  }
}
`, r.template(data, SkuStandardPlan), data.RandomInteger)
}

func (r WindowsFunctionAppSlotResource) withIPRestrictionsDescription(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%d"
  function_app_id            = azurerm_windows_function_app.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    ip_restriction {
      ip_address = "13.107.6.152/31,13.107.128.0/22"
      name       = "test-restriction"
      priority   = 123
      action     = "Allow"
      headers {
        x_azure_fdid      = ["55ce4ed1-4b06-4bf1-b40e-4638452104da"]
        x_fd_health_probe = ["1"]
        x_forwarded_for   = ["9.9.9.9/32", "2002::1234:abcd:ffff:c0a8:101/64"]
        x_forwarded_host  = ["example.com"]
      }
      description = "Allow ip address windows function app slot"
    }
  }
}
`, r.template(data, SkuStandardPlan), data.RandomInteger)
}

func (r WindowsFunctionAppSlotResource) withIPRestrictionsDefaultActionDeny(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%d"
  function_app_id            = azurerm_windows_function_app.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    ip_restriction_default_action = "Deny"

    ip_restriction {
      ip_address = "13.107.6.152/31,13.107.128.0/22"
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
  }
}
`, r.template(data, SkuStandardPlan), data.RandomInteger)
}

func (r WindowsFunctionAppSlotResource) appSettings(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%d"
  function_app_id            = azurerm_windows_function_app.test.id
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

func (r WindowsFunctionAppSlotResource) appSettingsUserSettings(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%d"
  function_app_id            = azurerm_windows_function_app.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  app_settings = {}

  site_config {}
}
`, r.template(data, planSku), data.RandomInteger)
}

func (r WindowsFunctionAppSlotResource) appSettingsUserSettingsUpdate(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%d"
  function_app_id            = azurerm_windows_function_app.test.id
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

func (r WindowsFunctionAppSlotResource) appSettingsCustomContentShare(data acceptance.TestData, planSku string, customShareName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%d"
  function_app_id            = azurerm_windows_function_app.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  app_settings = {
    foo                  = "bar"
    secret               = "sauce"
    WEBSITE_CONTENTSHARE = "%s"
  }

  site_config {}
}
`, r.template(data, planSku), data.RandomInteger, customShareName)
}

func (r WindowsFunctionAppSlotResource) backup(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%d"
  function_app_id            = azurerm_windows_function_app.test.id
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

func (r WindowsFunctionAppSlotResource) backupVnetIntegration(data acceptance.TestData, planSku string, enabled string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%d"
  function_app_id            = azurerm_windows_function_app.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  virtual_network_subnet_id              = azurerm_subnet.test.id
  virtual_network_backup_restore_enabled = %s

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
`, r.storageWithVnetIntegrationTemplate(data, planSku), data.RandomInteger, enabled)
}

func (r WindowsFunctionAppSlotResource) consumptionComplete(data acceptance.TestData) string {
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

resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%[2]d"
  function_app_id            = azurerm_windows_function_app.test.id
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

  builtin_logging_enabled            = false
  client_certificate_enabled         = true
  client_certificate_mode            = "Required"
  client_certificate_exclusion_paths = "/foo;/bar;/hello;/world"

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
    load_balancing_mode      = "LeastResponseTime"
    remote_debugging_enabled = true
    remote_debugging_version = "VS2022"

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
    health_check_eviction_time_in_min = 3

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

  ftp_publish_basic_authentication_enabled       = false
  webdeploy_publish_basic_authentication_enabled = false

  tags = {
    terraform = "true"
    Env       = "AccTest"
  }
}
`, r.template(data, planSku), data.RandomInteger, data.Client().TenantID)
}

func (r WindowsFunctionAppSlotResource) standardComplete(data acceptance.TestData) string {
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

resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%[2]d"
  function_app_id            = azurerm_windows_function_app.test.id
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
    health_check_eviction_time_in_min = 3
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

  tags = {
    terraform = "true"
    Env       = "AccTest"
  }
}
`, r.storageContainerTemplate(data, planSku), data.RandomInteger, data.Client().TenantID)
}

func (r WindowsFunctionAppSlotResource) elasticComplete(data acceptance.TestData) string {
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


resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%[2]d"
  function_app_id            = azurerm_windows_function_app.test.id
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

    default_documents = [
      "first.html",
      "second.jsp",
      "third.aspx",
      "hostingstart.html",
    ]

    elastic_instance_minimum = 2

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
    health_check_eviction_time_in_min = 3
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
}
`, r.storageContainerTemplate(data, SkuElasticPremiumPlan), data.RandomInteger)
}

func (r WindowsFunctionAppSlotResource) withAuthSettings(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%d"
  function_app_id            = azurerm_windows_function_app.test.id
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

func (r WindowsFunctionAppSlotResource) builtInLogging(data acceptance.TestData, planSku string, builtInLogging bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%d"
  function_app_id            = azurerm_windows_function_app.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  builtin_logging_enabled = %t

  site_config {}
}
`, r.template(data, planSku), data.RandomInteger, builtInLogging)
}

func (r WindowsFunctionAppSlotResource) connectionStrings(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%d"
  function_app_id            = azurerm_windows_function_app.test.id
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

func (r WindowsFunctionAppSlotResource) connectionStringsUpdate(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%d"
  function_app_id            = azurerm_windows_function_app.test.id
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

func (r WindowsFunctionAppSlotResource) dailyTimeLimitQuota(data acceptance.TestData, planSku string, quota int) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%d"
  function_app_id            = azurerm_windows_function_app.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  daily_memory_time_quota = %d

  site_config {}
}
`, r.template(data, planSku), data.RandomInteger, quota)
}

func (r WindowsFunctionAppSlotResource) healthCheckPath(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%d"
  function_app_id            = azurerm_windows_function_app.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    health_check_path                 = "/health"
    health_check_eviction_time_in_min = 3
  }
}
`, r.template(data, planSku), data.RandomInteger)
}

func (r WindowsFunctionAppSlotResource) runtimeScaleCheck(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%d"
  function_app_id            = azurerm_windows_function_app.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    runtime_scale_monitoring_enabled = true
  }
}
`, r.template(data, planSku), data.RandomInteger)
}

func (r WindowsFunctionAppSlotResource) healthCheckPathWithEviction(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%d"
  function_app_id            = azurerm_windows_function_app.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    health_check_path                 = "/health"
    health_check_eviction_time_in_min = 3
  }
}
`, r.template(data, planSku), data.RandomInteger)
}

func (r WindowsFunctionAppSlotResource) appServiceLogs(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%d"
  function_app_id            = azurerm_windows_function_app.test.id
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

func (r WindowsFunctionAppSlotResource) appStackDotNet(data acceptance.TestData, planSku string, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%d"
  function_app_id            = azurerm_windows_function_app.test.id
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

func (r WindowsFunctionAppSlotResource) appStackDotNetIsolated(data acceptance.TestData, planSku string, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%d"
  function_app_id            = azurerm_windows_function_app.test.id
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

// nolint: unparam
func (r WindowsFunctionAppSlotResource) appStackNode(data acceptance.TestData, planSku string, nodeVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%d"
  function_app_id            = azurerm_windows_function_app.test.id
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

// nolint: unparam
func (r WindowsFunctionAppSlotResource) appStackJava(data acceptance.TestData, planSku string, javaVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%d"
  function_app_id            = azurerm_windows_function_app.test.id
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

func (r WindowsFunctionAppSlotResource) appStackPowerShellCore(data acceptance.TestData, planSku string, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%d"
  function_app_id            = azurerm_windows_function_app.test.id
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

func (r WindowsFunctionAppSlotResource) servicePlanUpdate(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%d"
  function_app_id            = azurerm_windows_function_app.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {}

  depends_on = [azurerm_service_plan.update]
}
`, r.templateServicePlanUpdate(data, planSku), data.RandomInteger)
}

func (r WindowsFunctionAppSlotResource) updateStorageAccount(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%d"
  function_app_id            = azurerm_windows_function_app.test.id
  storage_account_name       = azurerm_storage_account.update.name
  storage_account_access_key = azurerm_storage_account.update.primary_access_key

  site_config {}
}
`, r.templateExtraStorageAccount(data, planSku), data.RandomInteger)
}

func (r WindowsFunctionAppSlotResource) identitySystemAssigned(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%d"
  function_app_id            = azurerm_windows_function_app.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {}

  identity {
    type = "SystemAssigned"
  }
}
`, r.identityTemplate(data, planSku), data.RandomInteger)
}

func (r WindowsFunctionAppSlotResource) identitySystemAssignedUserAssigned(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%d"
  function_app_id            = azurerm_windows_function_app.test.id
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

func (r WindowsFunctionAppSlotResource) identityUserAssigned(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%d"
  function_app_id            = azurerm_windows_function_app.test.id
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

func (r WindowsFunctionAppSlotResource) identityUserAssignedKeyVaultIdentity(data acceptance.TestData, planSku string) string {
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

resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%[2]d"
  function_app_id            = azurerm_windows_function_app.test.id
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

func (r WindowsFunctionAppSlotResource) msiStorageAccount(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s


resource "azurerm_role_assignment" "func_app_access_to_storage" {
  scope                = azurerm_storage_account.test.id
  role_definition_name = "Storage Blob Data Owner"
  principal_id         = azurerm_windows_function_app_slot.test.identity[0].principal_id
}

resource "azurerm_windows_function_app_slot" "test" {
  name            = "acctest-WFAS-%d"
  function_app_id = azurerm_windows_function_app.test.id

  storage_account_name          = azurerm_storage_account.test.name
  storage_uses_managed_identity = true

  identity {
    type = "SystemAssigned"
  }

  site_config {
    application_stack {
      dotnet_version = "v6.0"
    }
  }
}
`, r.template(data, planSku), data.RandomInteger)
}

func (r WindowsFunctionAppSlotResource) storageAccountKVSecret(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
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

resource "azurerm_windows_function_app_slot" "test" {
  name            = "acctest-WFAS-%[3]d"
  function_app_id = azurerm_windows_function_app.test.id

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

func (r WindowsFunctionAppSlotResource) storageAccountKVSecretVersionless(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
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

resource "azurerm_windows_function_app_slot" "test" {
  name            = "acctest-WFAS-%[3]d"
  function_app_id = azurerm_windows_function_app.test.id

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

func (r WindowsFunctionAppSlotResource) separatePlan(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_service_plan" "test2" {
  name                = "acctestASP2-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  os_type             = "Windows"
  sku_name            = "%[3]s"
}

resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%[2]d"
  function_app_id            = azurerm_windows_function_app.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  service_plan_id = azurerm_service_plan.test2.id

  site_config {}
}
`, r.template(data, planSku), data.RandomInteger, SkuStandardPlan)
}

func (r WindowsFunctionAppSlotResource) separatePlanUpdate(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_service_plan" "test2" {
  name                = "acctestASP2-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  os_type             = "Windows"
  sku_name            = "%[3]s"
}

resource "azurerm_service_plan" "test3" {
  name                = "acctestASP3-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  os_type             = "Windows"
  sku_name            = "%[4]s"
}

resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%[2]d"
  function_app_id            = azurerm_windows_function_app.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  service_plan_id = azurerm_service_plan.test3.id

  site_config {}
}
`, r.template(data, planSku), data.RandomInteger, SkuStandardPlan, SkuPremiumPlan)
}

func (r WindowsFunctionAppSlotResource) publicNetworkAccessDisabled(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%d"
  function_app_id            = azurerm_windows_function_app.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  public_network_access_enabled = false

  site_config {}
}
`, r.template(data, planSku), data.RandomInteger)
}

// Config Templates

func (WindowsFunctionAppSlotResource) template(data acceptance.TestData, planSku string) string {
	var additionalConfig string
	if strings.EqualFold(planSku, "EP1") {
		additionalConfig = "maximum_elastic_worker_count = 5"
	}
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

resource "azurerm_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  os_type             = "Windows"
  sku_name            = "%[4]s"
  %[5]s
}

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {}
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, planSku, additionalConfig)
}

func (r WindowsFunctionAppSlotResource) storageContainerTemplate(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_container" "test" {
  name                  = "test"
  storage_account_id    = azurerm_storage_account.test.id
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

func (r WindowsFunctionAppSlotResource) storageWithVnetIntegrationTemplate(data acceptance.TestData, planSku string) string {
	timeFormat := "2006-01-02"
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network" "test" {
  name                = "vnet-%[2]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "subnet-%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "delegation"

    service_delegation {
      name    = "Microsoft.Web/serverFarms"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }

  service_endpoints = [
    "Microsoft.Storage"
  ]
}

resource "azurerm_storage_account_network_rules" "test" {
  storage_account_id = azurerm_storage_account.test.id

  default_action             = "Deny"
  virtual_network_subnet_ids = [azurerm_subnet.test.id]
}

resource "azurerm_storage_container" "test" {
  name                  = "test"
  storage_account_id    = azurerm_storage_account.test.id
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

  start  = "%[3]s"
  expiry = "%[4]s"

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
`, r.template(data, planSku), data.RandomInteger, time.Now().Format(timeFormat), time.Now().AddDate(0, 0, 1).Format(timeFormat))
}

func (r WindowsFunctionAppSlotResource) identityTemplate(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acct-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

`, r.template(data, planSku), data.RandomInteger)
}

func (r WindowsFunctionAppSlotResource) templateServicePlanUpdate(data acceptance.TestData, planSku string) string {
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

func (WindowsFunctionAppSlotResource) templateExtraStorageAccount(data acceptance.TestData, planSku string) string {
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

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {}
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, planSku)
}

func (r WindowsFunctionAppSlotResource) vNetIntegration_basic(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
%s
resource "azurerm_virtual_network" "test" {
  name                = "vnet-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
resource "azurerm_subnet" "test1" {
  name                 = "subnet1"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
  delegation {
    name = "delegation"
    service_delegation {
      name    = "Microsoft.Web/serverFarms"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}
resource "azurerm_subnet" "test2" {
  name                 = "subnet2"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
  delegation {
    name = "delegation"
    service_delegation {
      name    = "Microsoft.Web/serverFarms"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}
resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%d"
  function_app_id            = azurerm_windows_function_app.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
  site_config {}
}
`, r.template(data, planSku), data.RandomInteger, data.RandomInteger)
}

func (r WindowsFunctionAppSlotResource) vNetIntegration_subnet1(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
%s
resource "azurerm_virtual_network" "test" {
  name                = "vnet-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
resource "azurerm_subnet" "test1" {
  name                 = "subnet1"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
  delegation {
    name = "delegation"
    service_delegation {
      name    = "Microsoft.Web/serverFarms"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}
resource "azurerm_subnet" "test2" {
  name                 = "subnet2"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
  delegation {
    name = "delegation"
    service_delegation {
      name    = "Microsoft.Web/serverFarms"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}
resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%d"
  function_app_id            = azurerm_windows_function_app.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
  virtual_network_subnet_id  = azurerm_subnet.test1.id
  site_config {}
}
`, r.template(data, planSku), data.RandomInteger, data.RandomInteger)
}

func (r WindowsFunctionAppSlotResource) vNetIntegration_subnet2(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
%s
resource "azurerm_virtual_network" "test" {
  name                = "vnet-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

}
resource "azurerm_subnet" "test1" {
  name                 = "subnet1"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
  delegation {
    name = "delegation"
    service_delegation {
      name    = "Microsoft.Web/serverFarms"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}
resource "azurerm_subnet" "test2" {
  name                 = "subnet2"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
  delegation {
    name = "delegation"
    service_delegation {
      name    = "Microsoft.Web/serverFarms"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}
resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%d"
  function_app_id            = azurerm_windows_function_app.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
  virtual_network_subnet_id  = azurerm_subnet.test2.id
  site_config {}
}
`, r.template(data, planSku), data.RandomInteger, data.RandomInteger)
}

func (r WindowsFunctionAppSlotResource) vNetIntegration_subnet1WithVnetProperties(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
%s
resource "azurerm_virtual_network" "test" {
  name                = "vnet-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
resource "azurerm_subnet" "test1" {
  name                 = "subnet1"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
  delegation {
    name = "delegation"
    service_delegation {
      name    = "Microsoft.Web/serverFarms"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}
resource "azurerm_subnet" "test2" {
  name                 = "subnet2"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
  delegation {
    name = "delegation"
    service_delegation {
      name    = "Microsoft.Web/serverFarms"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}
resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%d"
  function_app_id            = azurerm_windows_function_app.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
  virtual_network_subnet_id  = azurerm_subnet.test1.id
  vnet_image_pull_enabled    = true

  site_config {}
}
`, r.template(data, planSku), data.RandomInteger, data.RandomInteger)
}

func (r WindowsFunctionAppSlotResource) withASEV3(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[2]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_windows_function_app" "test" {
  name                = "acctest-WFA-%[3]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  vnet_image_pull_enabled = true
  site_config {
    vnet_route_all_enabled = true
  }
}

resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%[3]d"
  function_app_id            = azurerm_windows_function_app.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  vnet_image_pull_enabled = true
  site_config {
    vnet_route_all_enabled = true
  }
}

`, ServicePlanResource{}.aseV3(data), data.RandomString, data.RandomInteger)
}

func (r WindowsFunctionAppSlotResource) withStorageAccountSingle(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%d"
  function_app_id            = azurerm_windows_function_app.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  storage_account {
    name         = "files"
    type         = "AzureFiles"
    account_name = azurerm_storage_account.test.name
    share_name   = azurerm_storage_share.test.name
    access_key   = azurerm_storage_account.test.primary_access_key
    mount_path   = "\\mounts\\files"
  }

  site_config {}
}
`, r.templateWithStorageAccountExtras(data, planSku), data.RandomInteger)
}

func (r WindowsFunctionAppSlotResource) withStorageAccountMultiple(data acceptance.TestData, planSku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%d"
  function_app_id            = azurerm_windows_function_app.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  storage_account {
    name         = "files"
    type         = "AzureFiles"
    account_name = azurerm_storage_account.test.name
    share_name   = azurerm_storage_share.test.name
    access_key   = azurerm_storage_account.test.primary_access_key
    mount_path   = "\\mounts\\files"
  }

  storage_account {
    name         = "morefiles"
    type         = "AzureFiles"
    account_name = azurerm_storage_account.test.name
    share_name   = azurerm_storage_share.test2.name
    access_key   = azurerm_storage_account.test.primary_access_key
    mount_path   = "\\mounts\\morefiles"
  }

  site_config {}
}
`, r.templateWithStorageAccountExtras(data, planSku), data.RandomInteger)
}

func (r WindowsFunctionAppSlotResource) templateWithStorageAccountExtras(data acceptance.TestData, planSKU string) string {
	return fmt.Sprintf(`


%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acct-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_storage_container" "test" {
  name                  = "test"
  storage_account_id    = azurerm_storage_account.test.id
  container_access_type = "private"
}

resource "azurerm_storage_share" "test" {
  name               = "test"
  storage_account_id = azurerm_storage_account.test.id
  quota              = 1
}

resource "azurerm_storage_container" "test2" {
  name                  = "test2"
  storage_account_id    = azurerm_storage_account.test.id
  container_access_type = "private"
}

resource "azurerm_storage_share" "test2" {
  name               = "test2"
  storage_account_id = azurerm_storage_account.test.id
  quota              = 1
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
`, r.template(data, planSKU), data.RandomInteger)
}

func (r WindowsFunctionAppSlotResource) withTlsVersion(data acceptance.TestData, planSku string, tlsVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_windows_function_app_slot" "test" {
  name                       = "acctest-WFAS-%d"
  function_app_id            = azurerm_windows_function_app.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  site_config {
    minimum_tls_version = "%s"
  }
}
`, r.template(data, planSku), data.RandomInteger, tlsVersion)
}

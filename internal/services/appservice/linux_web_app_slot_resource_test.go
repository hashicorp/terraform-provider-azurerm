// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LinuxWebAppSlotResource struct{}

func TestAccLinuxWebAppSlot_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_updateTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.basicWithTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("app_settings.WEBSITE_HEALTHCHECK_MAXPINGFAILURES").DoesNotExist(),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccLinuxWebAppSlot_autoSwap(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoSwap(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_separateStandardPlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.separatePlan(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_separateStandardPlanUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.separatePlan(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.separatePlanUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_parentStandardPlanError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.parentStandardPlanError(data),
			ExpectError: regexp.MustCompile("`service_plan_id` should only be specified when it differs from the `service_plan_id` of the associated Web App"),
		},
		{
			Config: r.separatePlan(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config:      r.parentStandardPlanError(data),
			ExpectError: regexp.MustCompile("`service_plan_id` should only be specified when it differs from the `service_plan_id` of the associated Web App"),
		},
	})
}

// Complete

func TestAccLinuxWebAppSlot_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_completeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithStorage(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.always_on").HasValue("true"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.basicWithStorage(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

// Block Tests

func TestAccLinuxWebAppSlot_backup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withBackup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_backupUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithStorage(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.withBackup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.basicWithStorage(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_scmIpRestrictionSubnet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.scmIpRestrictionSubnet(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_scmIpRestrictionHeaders(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.scmIpRestrictionHeaders(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withConnectionStrings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withConnectionStrings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("app,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withConnectionStringsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("app,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.withConnectionStrings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("app,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.withConnectionStringsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("app,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.withConnectionStrings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("app,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("app,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withIPRestrictions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withIPRestrictions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withIPRestrictionsDescription(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withIPRestrictions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.withIPRestrictionsDescription(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.withIPRestrictions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withIPRestrictionsDefaultAction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withIPRestrictionsDefaultActionDeny(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withIPRestrictionsDefaultActionUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withIPRestrictions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.withIPRestrictionsDefaultActionDeny(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.withIPRestrictions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withIPRangeRestrictions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withIPRangeRestrictions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withIPRestrictionsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withIPRestrictions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.withIPRestrictionsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.withIPRestrictions(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withAuthSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withAuthSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withAuthSettingsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.withAuthSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.withAuthSettingsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.withAuthSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withAutoHealRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoHealRules(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccLinuxWebAppSlot_withAutoHealSlowRequestWithPath(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoHealRulesSlowRequestWithPath(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccLinuxWebAppSlot_withAutoHealRulesUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.autoHealRules(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.autoHealRulesUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_appSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("app_settings.foo").HasValue("bar"),
				check.That(data.ResourceName).Key("app_settings.secret").HasValue("sauce"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.identitySystemAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.identityUserAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.identitySystemAssignedUserAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_identityKeyVaultIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.identityUserAssignedKeyVaultIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

// Atrtibutes

func TestAccLinuxWebAppSlot_loadBalancing(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.loadBalancing(data, "WeightedRoundRobin"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("app,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_detailedLogging(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withDetailedLogging(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("app,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withLoggingUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.withDetailedLogging(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.detailed_error_logging_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("logs.0.detailed_error_messages").HasValue("true"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.withLogsHttpBlob(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.detailed_error_logging_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("logs.0.detailed_error_messages").HasValue("false"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.withLoggingComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.detailed_error_logging_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("logs.0.detailed_error_messages").HasValue("true"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

// App Stacks

func TestAccLinuxWebAppSlot_withDotNet31(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dotNet(data, "3.1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withDotNet50(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dotNet(data, "5.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withDotNet60(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dotNet(data, "6.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withDotNet70(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dotNet(data, "7.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withDotNet80(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dotNet(data, "8.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withPhp74(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.php(data, "7.4"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withPhp80(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.php(data, "8.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withPhp81(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.php(data, "8.1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withPhp82(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.php(data, "8.2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withPhp83(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.php(data, "8.3"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withPython37(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.python(data, "3.7"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withPython38(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.python(data, "3.8"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withPython39(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.python(data, "3.9"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withPython310(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.python(data, "3.10"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withPython311(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.python(data, "3.11"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withPython312(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.python(data, "3.12"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withRuby26(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.ruby(data, "2.6"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withRuby27(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.ruby(data, "2.7"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withNode12LTS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.node(data, "12-lts"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withNode14LTS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.node(data, "14-lts"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withNode16LTS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.node(data, "16-lts"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withNode18LTS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.node(data, "18-lts"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withNode20LTS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.node(data, "20-lts"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withJre8Java(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.java(data, "8", "JAVA", "8"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withJre11Java(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.java(data, "11", "JAVA", "11"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.linux_fx_version").HasValue("JAVA|11-java11"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withJava1109(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.java(data, "11", "JAVA", "11.0.9"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.linux_fx_version").HasValue("JAVA|11.0.9"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withJava8u242(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.java(data, "8", "JAVA", "8u242"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.linux_fx_version").HasValue("JAVA|8u242"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withJava11Tomcat9(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.java(data, "11", "TOMCAT", "9.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.linux_fx_version").HasValue("TOMCAT|9.0-java11"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withJava11Tomcat8561(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.java(data, "11", "TOMCAT", "8.5.61"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.linux_fx_version").HasValue("TOMCAT|8.5.61-java11"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withJava8JBOSSEAP73(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.javaPremiumV3Plan(data, "8", "JBOSSEAP", "7.3"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.linux_fx_version").HasValue("JBOSSEAP|7.3-java8"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withDocker(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skipf("Skippped as deprecated property removed in 4.0")
	}
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.docker(data, "mcr.microsoft.com/appsvc/staticsite", "latest"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.linux_fx_version").HasValue("DOCKER|mcr.microsoft.com/appsvc/staticsite:latest"),
			),
		},
		data.ImportStep("app_settings.%",
			"app_settings.DOCKER_REGISTRY_SERVER_PASSWORD",
			"app_settings.DOCKER_REGISTRY_SERVER_URL",
			"app_settings.DOCKER_REGISTRY_SERVER_USERNAME",
			"site_config.0.application_stack.0.docker_image",
			"site_config.0.application_stack.0.docker_image_name",
			"site_config.0.application_stack.0.docker_image_tag",
			"site_config.0.application_stack.0.docker_registry_url",
			"site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withDockerHub(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skipf("Skippped as deprecated property removed in 4.0")
	}
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dockerHub(data, "nginx", "latest"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.linux_fx_version").HasValue("DOCKER|nginx:latest"),
			),
		},
		data.ImportStep("app_settings.%",
			"app_settings.DOCKER_REGISTRY_SERVER_PASSWORD",
			"app_settings.DOCKER_REGISTRY_SERVER_URL",
			"app_settings.DOCKER_REGISTRY_SERVER_USERNAME",
			"site_config.0.application_stack.0.docker_image",
			"site_config.0.application_stack.0.docker_image_name",
			"site_config.0.application_stack.0.docker_image_tag",
			"site_config.0.application_stack.0.docker_registry_url",
			"site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withDockerDeprecatedUpgrade(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skipf("Skippped as deprecated property removed in 4.0")
	}
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dockerHub(data, "nginx", "latest"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.linux_fx_version").HasValue("DOCKER|nginx:latest"),
			),
		},
		data.ImportStep("app_settings.%",
			"app_settings.DOCKER_REGISTRY_SERVER_PASSWORD",
			"app_settings.DOCKER_REGISTRY_SERVER_URL",
			"app_settings.DOCKER_REGISTRY_SERVER_USERNAME",
			"site_config.0.application_stack.0.docker_image",
			"site_config.0.application_stack.0.docker_image_name",
			"site_config.0.application_stack.0.docker_image_tag",
			"site_config.0.application_stack.0.docker_registry_url",
			"site_credential.0.password"),
		{
			Config: r.dockerImageName(data, "https://index.docker.io", "nginx:latest"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.linux_fx_version").HasValue("DOCKER|index.docker.io/nginx:latest"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withDockerImageMCR(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dockerImageName(data, "https://mcr.microsoft.com", "appsvc/staticsite:latest"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.linux_fx_version").HasValue("DOCKER|mcr.microsoft.com/appsvc/staticsite:latest"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withDockerImageDockerHub(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dockerImageName(data, "https://index.docker.io", "nginx:latest"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.linux_fx_version").HasValue("DOCKER|index.docker.io/nginx:latest"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_withDockerImageDockerHubUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dockerImageName(data, "https://index.docker.io", "nginx:latest"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.linux_fx_version").HasValue("DOCKER|index.docker.io/nginx:latest"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.dockerImageName(data, "https://index.docker.io", "nginx:stable"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("site_config.0.linux_fx_version").HasValue("DOCKER|index.docker.io/nginx:stable"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

// Deployments

func TestAccLinuxWebAppSlot_zipDeploy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.zipDeploy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("zip_deploy_file",
			"site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_disableDeployBasicAuth(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.deployBasicAuthDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_disableDeployBasicAuthUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.deployBasicAuthDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

// Network tests

func TestAccLinuxWebAppSlot_vNetIntegration(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.vNetIntegrationWebApp_subnet1(data),
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

func TestAccLinuxWebAppSlot_vNetIntegrationUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.vNetIntegrationWebApp_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.vNetIntegrationWebApp_subnet1(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("virtual_network_subnet_id").MatchesOtherKey(
					check.That("azurerm_subnet.test1").Key("id"),
				),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.vNetIntegrationWebApp_subnet2(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("virtual_network_subnet_id").MatchesOtherKey(
					check.That("azurerm_subnet.test2").Key("id"),
				),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.vNetIntegrationWebApp_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_publicNetworkAccessDisabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.publicNetworkAccessDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("false"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccLinuxWebAppSlot_publicNetworkAccessUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_slot", "test")
	r := LinuxWebAppSlotResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("true"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.publicNetworkAccessDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("false"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("true"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

// Exists

func (r LinuxWebAppSlotResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := webapps.ParseSlotID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.AppService.WebAppsClient.GetSlot(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Linux %s: %+v", id, err)
	}

	if response.WasNotFound(resp.HttpResponse) {
		return utils.Bool(false), nil
	}

	return utils.Bool(true), nil
}

// Configs

func (r LinuxWebAppSlotResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_linux_web_app.test.id

  site_config {}
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r LinuxWebAppSlotResource) basicWithTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_linux_web_app.test.id

  site_config {}

  tags = {
    foo = "bar"
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r LinuxWebAppSlotResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`

%s

resource "azurerm_linux_web_app_slot" "import" {
  name           = azurerm_linux_web_app_slot.test.name
  app_service_id = azurerm_linux_web_app_slot.test.app_service_id

  site_config {}
}
`, r.basic(data))
}

func (r LinuxWebAppSlotResource) appSettings(data acceptance.TestData) string {
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
    foo    = "bar"
    secret = "sauce"
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r LinuxWebAppSlotResource) autoHealRules(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_linux_web_app.test.id

  site_config {
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
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r LinuxWebAppSlotResource) autoHealRulesSlowRequestWithPath(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
%s
resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_linux_web_app.test.id

  site_config {
    auto_heal_setting {
      trigger {
        slow_request {
          count      = "10"
          interval   = "00:10:00"
          time_taken = "00:00:10"
        }
        slow_request_with_path {
          count      = "11"
          interval   = "00:11:00"
          time_taken = "00:00:11"
          path       = "/tftest1"
        }
        slow_request_with_path {
          count      = "12"
          interval   = "00:12:00"
          time_taken = "00:00:12"
          path       = "/tftest2"
        }
      }
      action {
        action_type                    = "Recycle"
        minimum_process_execution_time = "00:05:00"
      }
    }
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r LinuxWebAppSlotResource) autoHealRulesUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_linux_web_app.test.id

  site_config {
    auto_heal_setting {
      trigger {
        status_code {
          status_code_range = "500"
          interval          = "00:01:00"
          count             = 10
        }
        status_code {
          status_code_range = "400-404"
          interval          = "00:10:00"
          count             = 10
        }
      }

      action {
        action_type                    = "Recycle"
        minimum_process_execution_time = "00:10:00"
      }
    }
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r LinuxWebAppSlotResource) autoSwap(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_linux_web_app.test.id

  site_config {
    auto_swap_slot_name = "production"
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r LinuxWebAppSlotResource) basicWithStorage(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_linux_web_app.test.id

  site_config {}
}
`, r.templateWithStorageAccount(data), data.RandomInteger)
}

func (r LinuxWebAppSlotResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_linux_web_app.test.id

  app_settings = {
    foo = "bar"
  }

  auth_settings {
    enabled = true
    issuer  = "https://sts.windows.net/%s"

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
  client_certificate_mode            = "Optional"
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
    use_32_bit_worker                 = true
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

    auto_swap_slot_name = "Production"

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

  storage_account {
    name         = "files"
    type         = "AzureFiles"
    account_name = azurerm_storage_account.test.name
    share_name   = azurerm_storage_share.test.name
    access_key   = azurerm_storage_account.test.primary_access_key
    mount_path   = "/storage/files"
  }

  tags = {
    Environment = "AccTest"
    foo         = "bar"
  }
}
`, r.templateWithStorageAccount(data), data.RandomInteger, data.Client().TenantID)
}

func (r LinuxWebAppSlotResource) withDetailedLogging(data acceptance.TestData, detailedErrorLogging bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_linux_web_app.test.id

  site_config {}

  logs {
    detailed_error_messages = %t
  }
}
`, r.baseTemplate(data), data.RandomInteger, detailedErrorLogging)
}

func (r LinuxWebAppSlotResource) loadBalancing(data acceptance.TestData, loadBalancingMode string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_linux_web_app.test.id

  site_config {
    load_balancing_mode = "%s"
  }
}
`, r.baseTemplate(data), data.RandomInteger, loadBalancingMode)
}

func (r LinuxWebAppSlotResource) withAuthSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_linux_web_app.test.id

  site_config {}

  auth_settings {
    enabled = true
    issuer  = "https://sts.windows.net/%s"

    additional_login_parameters = {
      test_key = "test_value"
    }

    allowed_external_redirect_urls = ["https://example.com"]

    active_directory {
      client_id     = "aadclientid"
      client_secret = "aadsecret"

      allowed_audiences = [
        "activedirectorytokenaudiences",
      ]
    }
  }
}
`, r.baseTemplate(data), data.RandomInteger, data.Client().TenantID)
}

func (r LinuxWebAppSlotResource) withAuthSettingsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_linux_web_app.test.id

  site_config {}

  auth_settings {
    enabled = true
    issuer  = "https://sts.windows.net/%s"

    additional_login_parameters = {
      test_key = "test_value"
    }

    allowed_external_redirect_urls = ["https://example.com"]

    default_provider              = "AzureActiveDirectory"
    token_refresh_extension_hours = 24
    token_store_enabled           = true
    unauthenticated_client_action = "RedirectToLoginPage"

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

    google {
      client_id     = "googleclientid"
      client_secret = "googleclientsecret"

      oauth_scopes = [
        "googlescope",
      ]
    }

    microsoft {
      client_id     = "microsoftclientid"
      client_secret = "microsoftclientsecret"

      oauth_scopes = [
        "microsoftscope",
      ]
    }

    twitter {
      consumer_key    = "twitterconsumerkey"
      consumer_secret = "twitterconsumersecret"
    }
  }
}
`, r.baseTemplate(data), data.RandomInteger, data.Client().TenantID)
}

func (r LinuxWebAppSlotResource) withBackup(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_linux_web_app.test.id

  site_config {}

  backup {
    name                = "acctest"
    storage_account_url = "https://${azurerm_storage_account.test.name}.blob.core.windows.net/${azurerm_storage_container.test.name}${data.azurerm_storage_account_sas.test.sas}&sr=b"
    schedule {
      frequency_interval = 7
      frequency_unit     = "Day"
    }
  }
}
`, r.templateWithStorageAccount(data), data.RandomInteger)
}

func (r LinuxWebAppSlotResource) withConnectionStrings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_linux_web_app.test.id

  connection_string {
    name  = "First"
    value = "first-connection-string"
    type  = "Custom"
  }

  site_config {}
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r LinuxWebAppSlotResource) scmIpRestrictionSubnet(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%[2]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%[2]d"
  app_service_id = azurerm_linux_web_app.test.id

  site_config {
    scm_ip_restriction {
      virtual_network_subnet_id = azurerm_subnet.test.id
    }
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r LinuxWebAppSlotResource) scmIpRestrictionHeaders(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_linux_web_app.test.id

  site_config {
    scm_ip_restriction {
      action     = "Allow"
      ip_address = "10.10.10.10/32"
      name       = "test-restriction"
      priority   = 123
      headers {
        x_azure_fdid      = ["55ce4ed1-4b06-4bf1-b40e-4638452104da"]
        x_fd_health_probe = ["1"]
        x_forwarded_for   = ["9.9.9.9/32", "2002::1234:abcd:ffff:c0a8:101/64"]
        x_forwarded_host  = ["example.com"]
      }
    }
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r LinuxWebAppSlotResource) withConnectionStringsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_linux_web_app.test.id

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

  site_config {}
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r LinuxWebAppSlotResource) withIPRestrictions(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_linux_web_app.test.id

  site_config {
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
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r LinuxWebAppSlotResource) withIPRestrictionsDescription(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_linux_web_app.test.id

  site_config {
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
      description = "Allow ip address linux function app"
    }
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r LinuxWebAppSlotResource) withIPRestrictionsDefaultActionDeny(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_linux_web_app.test.id

  site_config {
    ip_restriction_default_action = "Deny"

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
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r LinuxWebAppSlotResource) withIPRangeRestrictions(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_linux_web_app.test.id

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
      description = "Allow ip address linux web app slot"
    }
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r LinuxWebAppSlotResource) withIPRestrictionsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_linux_web_app.test.id

  site_config {
    ip_restriction {
      ip_address = "10.10.10.10/32"
      name       = "test-restriction"
      priority   = 123
      action     = "Allow"
      headers {
        x_azure_fdid      = ["55ce4ed1-4b06-4bf1-b40e-4638452104da", "6bde7211-57bc-4476-866a-c9676e22b9d7"]
        x_fd_health_probe = ["1"]
        x_forwarded_for   = ["9.9.9.9/32", "2002::1234:abcd:ffff:c0a8:101/64", "9.9.9.8/32"]
        x_forwarded_host  = ["example.com", "anotherexample.com"]
      }
    }
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r LinuxWebAppSlotResource) withLogsHttpBlob(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_linux_web_app.test.id

  site_config {}

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
}

`, r.templateWithStorageAccount(data), data.RandomInteger)
}

func (r LinuxWebAppSlotResource) withLoggingComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_linux_web_app.test.id

  site_config {}

  logs {
    detailed_error_messages = true
    failed_request_tracing  = true

    application_logs {
      file_system_level = "Warning"

      azure_blob_storage {
        level             = "Information"
        sas_url           = "http://x.com/"
        retention_in_days = 7
      }
    }

    http_logs {
      file_system {
        retention_in_days = 4
        retention_in_mb   = 25
      }
    }
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r LinuxWebAppSlotResource) dotNet(data acceptance.TestData, dotNetVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_linux_web_app.test.id

  site_config {
    application_stack {
      dotnet_version = "%s"
    }
  }
}

`, r.baseTemplate(data), data.RandomInteger, dotNetVersion)
}

func (r LinuxWebAppSlotResource) php(data acceptance.TestData, phpVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_linux_web_app.test.id

  site_config {
    application_stack {
      php_version = "%s"
    }
  }
}




`, r.baseTemplate(data), data.RandomInteger, phpVersion)
}

func (r LinuxWebAppSlotResource) python(data acceptance.TestData, pythonVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_linux_web_app.test.id

  site_config {
    application_stack {
      python_version = "%s"
    }
  }
}

`, r.baseTemplate(data), data.RandomInteger, pythonVersion)
}

func (r LinuxWebAppSlotResource) ruby(data acceptance.TestData, rubyVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_linux_web_app.test.id

  site_config {
    application_stack {
      ruby_version = "%s"
    }
  }
}

`, r.baseTemplate(data), data.RandomInteger, rubyVersion)
}

func (r LinuxWebAppSlotResource) node(data acceptance.TestData, nodeVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_linux_web_app.test.id

  site_config {
    application_stack {
      node_version = "%s"
    }
  }
}

`, r.baseTemplate(data), data.RandomInteger, nodeVersion)
}

func (r LinuxWebAppSlotResource) java(data acceptance.TestData, javaVersion, javaServer, javaServerVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_linux_web_app.test.id

  site_config {
    application_stack {
      java_version        = "%s"
      java_server         = "%s"
      java_server_version = "%s"
    }
  }
}

`, r.baseTemplate(data), data.RandomInteger, javaVersion, javaServer, javaServerVersion)
}

func (r LinuxWebAppSlotResource) javaPremiumV3Plan(data acceptance.TestData, javaVersion, javaServer, javaServerVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_linux_web_app.test.id

  site_config {
    application_stack {
      java_version        = "%s"
      java_server         = "%s"
      java_server_version = "%s"
    }
  }
}

`, r.premiumV3PlanTemplate(data), data.RandomInteger, javaVersion, javaServer, javaServerVersion)
}

func (r LinuxWebAppSlotResource) docker(data acceptance.TestData, containerImage, containerTag string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_linux_web_app.test.id

  app_settings = {
    "DOCKER_REGISTRY_SERVER_URL"          = "https://mcr.microsoft.com"
    "DOCKER_REGISTRY_SERVER_USERNAME"     = ""
    "DOCKER_REGISTRY_SERVER_PASSWORD"     = ""
    "WEBSITES_ENABLE_APP_SERVICE_STORAGE" = "false"

  }

  site_config {
    application_stack {
      docker_image     = "%s"
      docker_image_tag = "%s"
    }
  }
}

`, r.baseTemplate(data), data.RandomInteger, containerImage, containerTag)
}

func (r LinuxWebAppSlotResource) dockerHub(data acceptance.TestData, containerImage, containerTag string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_linux_web_app.test.id

  app_settings = {
    "DOCKER_REGISTRY_SERVER_URL"          = "https://index.docker.io"
    "DOCKER_REGISTRY_SERVER_USERNAME"     = ""
    "DOCKER_REGISTRY_SERVER_PASSWORD"     = ""
    "WEBSITES_ENABLE_APP_SERVICE_STORAGE" = "false"
  }

  site_config {
    application_stack {
      docker_image     = "%s"
      docker_image_tag = "%s"
    }
  }
}

`, r.baseTemplate(data), data.RandomInteger, containerImage, containerTag)
}

func (r LinuxWebAppSlotResource) dockerImageName(data acceptance.TestData, registryUrl, containerImage string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_linux_web_app.test.id

  app_settings = {
    "WEBSITES_ENABLE_APP_SERVICE_STORAGE" = "false"
  }

  site_config {
    application_stack {
      docker_image_name   = "%s"
      docker_registry_url = "%s"
    }
  }
}

`, r.baseTemplate(data), data.RandomInteger, containerImage, registryUrl)
}

func (r LinuxWebAppSlotResource) identitySystemAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_linux_web_app.test.id

  site_config {}

  identity {
    type = "SystemAssigned"
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r LinuxWebAppSlotResource) identitySystemAssignedUserAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%[2]d"
  app_service_id = azurerm_linux_web_app.test.id

  site_config {}

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r LinuxWebAppSlotResource) identityUserAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%[2]d"
  app_service_id = azurerm_linux_web_app.test.id

  site_config {}

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r LinuxWebAppSlotResource) identityUserAssignedKeyVaultIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_user_assigned_identity" "kv" {
  name                = "acctestUAI-kv-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}


resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%[2]d"
  app_service_id = azurerm_linux_web_app.test.id

  site_config {}

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id, azurerm_user_assigned_identity.kv.id]
  }

  key_vault_reference_identity_id = azurerm_user_assigned_identity.kv.id
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r LinuxWebAppSlotResource) zipDeploy(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_linux_web_app.test.id

  app_settings = {
    WEBSITE_RUN_FROM_PACKAGE       = "1"
    SCM_DO_BUILD_DURING_DEPLOYMENT = "true"
  }

  site_config {
    application_stack {
      python_version = "3.9"
    }
  }

  zip_deploy_file = "./testdata/msdocs-python-flask-webapp-quickstart-main.zip"
}

`, r.baseTemplate(data), data.RandomInteger)
}

func (r LinuxWebAppSlotResource) separatePlan(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_service_plan" "test2" {
  name                = "acctestASP2-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  os_type             = "Linux"
  sku_name            = "%[3]s"
}

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%[2]d"
  app_service_id = azurerm_linux_web_app.test.id

  service_plan_id = azurerm_service_plan.test2.id

  site_config {}
}
`, r.baseTemplate(data), data.RandomInteger, SkuStandardPlan)
}

func (r LinuxWebAppSlotResource) separatePlanUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_service_plan" "test2" {
  name                = "acctestASP2-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  os_type             = "Linux"
  sku_name            = "%[3]s"
}

resource "azurerm_service_plan" "test3" {
  name                = "acctestASP3-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  os_type             = "Linux"
  sku_name            = "%[4]s"
}

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%[2]d"
  app_service_id = azurerm_linux_web_app.test.id

  service_plan_id = azurerm_service_plan.test3.id

  site_config {}
}
`, r.baseTemplate(data), data.RandomInteger, SkuStandardPlan, SkuPremiumPlan)
}

func (r LinuxWebAppSlotResource) parentStandardPlanError(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_service_plan" "test2" {
  name                = "acctestASP2-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  os_type             = "Linux"
  sku_name            = "%[3]s"
}

resource "azurerm_linux_web_app_slot" "test" {
  name            = "acctestWAS-%[2]d"
  app_service_id  = azurerm_linux_web_app.test.id
  service_plan_id = azurerm_service_plan.test.id

  site_config {}
}
`, r.baseTemplate(data), data.RandomInteger, SkuStandardPlan)
}

func (r LinuxWebAppSlotResource) publicNetworkAccessDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_linux_web_app.test.id

  public_network_access_enabled = false

  site_config {}
}
`, r.baseTemplate(data), data.RandomInteger)
}

func (r LinuxWebAppSlotResource) deployBasicAuthDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%d"
  app_service_id = azurerm_linux_web_app.test.id

  ftp_publish_basic_authentication_enabled       = false
  webdeploy_publish_basic_authentication_enabled = false

  site_config {}
}
`, r.baseTemplate(data), data.RandomInteger)
}

// Templates

func (LinuxWebAppSlotResource) baseTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_service_plan" "test" {
  name                = "acctestASP-WAS-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  os_type             = "Linux"
  sku_name            = "S1"
}

resource "azurerm_linux_web_app" "test" {
  name                = "acctestWA-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {}
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r LinuxWebAppSlotResource) templateWithStorageAccount(data acceptance.TestData) string {
	return fmt.Sprintf(`

%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acct-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "test"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_storage_share" "test" {
  name                 = "test"
  storage_account_name = azurerm_storage_account.test.name
  quota                = 1
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
`, r.baseTemplate(data), data.RandomInteger, data.RandomString)
}

func (LinuxWebAppSlotResource) premiumV3PlanTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  os_type             = "Linux"
  sku_name            = "P1v3"
}

resource "azurerm_linux_web_app" "test" {
  name                = "acctestWA-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {}
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r LinuxWebAppSlotResource) vNetIntegrationWebApp_basic(data acceptance.TestData) string {
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

resource "azurerm_linux_web_app_slot" "test" {
  name           = "acctestWAS-%[2]d"
  app_service_id = azurerm_linux_web_app.test.id

  site_config {}
}
`, r.baseTemplate(data), data.RandomInteger, data.RandomInteger)
}

func (r LinuxWebAppSlotResource) vNetIntegrationWebApp_subnet1(data acceptance.TestData) string {
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

resource "azurerm_linux_web_app_slot" "test" {
  name                      = "acctestWAS-%[2]d"
  app_service_id            = azurerm_linux_web_app.test.id
  virtual_network_subnet_id = azurerm_subnet.test1.id

  site_config {}
}
`, r.baseTemplate(data), data.RandomInteger, data.RandomInteger)
}

func (r LinuxWebAppSlotResource) vNetIntegrationWebApp_subnet2(data acceptance.TestData) string {
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

resource "azurerm_linux_web_app_slot" "test" {
  name                      = "acctestWAS-%[2]d"
  app_service_id            = azurerm_linux_web_app.test.id
  virtual_network_subnet_id = azurerm_subnet.test2.id

  site_config {}
}
`, r.baseTemplate(data), data.RandomInteger, data.RandomInteger)
}

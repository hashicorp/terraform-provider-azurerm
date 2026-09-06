// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package securitycenter_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

// TestAccSecurityCenterAutomation_V0ToV1_501 tests the state migration path from an `id` with lowercased
// static segments to their canonicalized format. It uses v5.0.1 as the setup version because it is
// the last release where the `id` could have been stored in state with lowercased static segments via an import using a non-canonical ID.
func TestAccSecurityCenterAutomation_V0ToV1_501(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_automation", "test")
	r := SecurityCenterAutomationResource{}

	importedResourceName := data.ResourceName + "-import"

	data.ResourceRegressionAdditionalStepsTest(t, r, []acceptance.TestStep{
		{
			Config: r.logicApp(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourceGroups/acctestRG-%[2]d/providers/Microsoft.Security/automations/acctestautomation-%[2]d", data.Subscriptions.Primary, data.RandomInteger)),
			),
		},
		{
			Config: r.logicAppV0(data),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					// Ensure import is a no-op to prevent the resource from normalizing the ID due to a combined CreateUpdate func
					plancheck.ExpectResourceAction(importedResourceName, plancheck.ResourceActionNoop),
				},
			},
			Check: acceptance.ComposeTestCheckFunc(
				check.That(importedResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourcegroups/acctestRG-%[2]d/providers/microsoft.security/automations/acctestautomation-%[2]d", data.Subscriptions.Primary, data.RandomInteger)),
			),
		},
		{
			Config: r.logicAppImported(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(importedResourceName).ExistsInAzure(r),
				check.That(importedResourceName).Key("id").HasValue(fmt.Sprintf("/subscriptions/%[1]s/resourceGroups/acctestRG-%[2]d/providers/Microsoft.Security/automations/acctestautomation-%[2]d", data.Subscriptions.Primary, data.RandomInteger)),
			),
		},
	}, "5.0.1")
}

func (SecurityCenterAutomationResource) v0MigrationTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_logic_app_workflow" "test" {
  name                = "acctestlogicapp-%[1]d"
  location            = "%[2]s"
  resource_group_name = azurerm_resource_group.test.name
}

data "azurerm_client_config" "current" {
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r SecurityCenterAutomationResource) logicAppV0(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_security_center_automation" "test-import" {
  name                = "acctestautomation-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  scopes = [
    "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  ]

  action {
    type        = "LogicApp"
    resource_id = azurerm_logic_app_workflow.test.id
    trigger_url = "https://example.net/this_is_never_validated_by_azure"
  }

  source {
    event_source = "Alerts"
  }

  tags = {
    Env2 = "Test2"
  }

  lifecycle {
    # ignore trigger_url, it's not returned on imports
    ignore_changes = [action[0].trigger_url]
  }
}

removed {
  from = azurerm_security_center_automation.test
  lifecycle {
    destroy = false
  }
}

import {
  to = azurerm_security_center_automation.test-import
  id = "/subscriptions/%[3]s/resourcegroups/acctestRG-%[2]d/providers/microsoft.security/automations/acctestautomation-%[2]d"
}
`, r.v0MigrationTemplate(data), data.RandomInteger, data.Subscriptions.Primary)
}

func (r SecurityCenterAutomationResource) logicAppImported(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_security_center_automation" "test-import" {
  name                = "acctestautomation-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  scopes = [
    "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  ]

  action {
    type        = "LogicApp"
    resource_id = azurerm_logic_app_workflow.test.id
    trigger_url = "https://example.net/this_is_never_validated_by_azure"
  }

  source {
    event_source = "Alerts"
  }

  tags = {
    Env2 = "Test2"
  }
}
`, r.v0MigrationTemplate(data), data.RandomInteger)
}

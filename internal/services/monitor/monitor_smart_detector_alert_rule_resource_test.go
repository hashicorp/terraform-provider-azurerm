// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2019-06-01/smartdetectoralertrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type MonitorSmartDetectorAlertRuleResource struct{}

func TestAccMonitorSmartDetectorAlertRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_smart_detector_alert_rule", "test")
	r := MonitorSmartDetectorAlertRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorSmartDetectorAlertRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_smart_detector_alert_rule", "test")
	r := MonitorSmartDetectorAlertRuleResource{}
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

func TestAccMonitorSmartDetectorAlertRule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_smart_detector_alert_rule", "test")
	r := MonitorSmartDetectorAlertRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorSmartDetectorAlertRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_smart_detector_alert_rule", "test")
	r := MonitorSmartDetectorAlertRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t MonitorSmartDetectorAlertRuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := smartdetectoralertrules.ParseSmartDetectorAlertRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Monitor.SmartDetectorAlertRulesClient.Get(ctx, *id, smartdetectoralertrules.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r MonitorSmartDetectorAlertRuleResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_smart_detector_alert_rule" "test" {
  name                = "acctestSDAR-%d"
  resource_group_name = azurerm_resource_group.test.name
  severity            = "Sev0"
  scope_resource_ids  = [azurerm_application_insights.test.id]
  frequency           = "PT1M"
  detector_type       = "FailureAnomaliesDetector"

  action_group {
    ids = [azurerm_monitor_action_group.test.id]
  }
}
`, r.template(data), data.RandomInteger)
}

func (r MonitorSmartDetectorAlertRuleResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_smart_detector_alert_rule" "import" {
  name                = azurerm_monitor_smart_detector_alert_rule.test.name
  resource_group_name = azurerm_monitor_smart_detector_alert_rule.test.resource_group_name
  severity            = azurerm_monitor_smart_detector_alert_rule.test.severity
  scope_resource_ids  = azurerm_monitor_smart_detector_alert_rule.test.scope_resource_ids
  frequency           = azurerm_monitor_smart_detector_alert_rule.test.frequency
  detector_type       = azurerm_monitor_smart_detector_alert_rule.test.detector_type

  action_group {
    ids = [azurerm_monitor_action_group.test.id]
  }
}
`, r.basic(data))
}

func (r MonitorSmartDetectorAlertRuleResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_smart_detector_alert_rule" "test" {
  name                = "acctestSDAR-%d"
  resource_group_name = azurerm_resource_group.test.name
  severity            = "Sev0"
  scope_resource_ids  = [azurerm_application_insights.test.id]
  frequency           = "PT1M"
  detector_type       = "FailureAnomaliesDetector"

  description = "acctest"
  enabled     = false

  action_group {
    ids             = [azurerm_monitor_action_group.test.id]
    email_subject   = "acctest email subject"
    webhook_payload = <<BODY
{
    "msg": "Acctest payload body"
}
BODY
  }

  throttling_duration = "PT20M"

  tags = {
    Env = "Test"
  }
}
`, r.template(data), data.RandomInteger)
}

func (MonitorSmartDetectorAlertRuleResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-monitor-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

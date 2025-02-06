// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2021-08-08/alertprocessingrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MonitorAlertProcessingRuleSuppressionResource struct{}

func TestAccMonitorAlertProcessingRuleSuppression_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_alert_processing_rule_suppression", "test")
	r := MonitorAlertProcessingRuleSuppressionResource{}

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

func TestAccMonitorAlertProcessingRuleSuppression_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_alert_processing_rule_suppression", "test")
	r := MonitorAlertProcessingRuleSuppressionResource{}
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

func TestAccMonitorAlertProcessingRuleSuppression_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_alert_processing_rule_suppression", "test")
	r := MonitorAlertProcessingRuleSuppressionResource{}

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

func TestAccMonitorAlertProcessingRuleSuppression_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_alert_processing_rule_suppression", "test")
	r := MonitorAlertProcessingRuleSuppressionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
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

func (r MonitorAlertProcessingRuleSuppressionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := alertprocessingrules.ParseActionRuleIDInsensitively(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Monitor.AlertProcessingRulesClient.GetByName(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving (%s): %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r MonitorAlertProcessingRuleSuppressionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_alert_processing_rule_suppression" "test" {
  name                = "acctest-moniter-%d"
  resource_group_name = azurerm_resource_group.test.name
  scopes              = [azurerm_resource_group.test.id]
}
`, r.template(data), data.RandomInteger)
}

func (r MonitorAlertProcessingRuleSuppressionResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_alert_processing_rule_suppression" "import" {
  name                = azurerm_monitor_alert_processing_rule_suppression.test.name
  resource_group_name = azurerm_monitor_alert_processing_rule_suppression.test.resource_group_name
  scopes              = azurerm_monitor_alert_processing_rule_suppression.test.scopes
}
`, r.basic(data))
}

func (r MonitorAlertProcessingRuleSuppressionResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_alert_processing_rule_suppression" "test" {
  name                = "acctest-moniter-%d"
  resource_group_name = azurerm_resource_group.test.name
  scopes              = [azurerm_resource_group.test.id]
  enabled             = false

  condition {
    signal_type {
      operator = "NotEquals"
      values   = ["Health"]
    }
  }

  schedule {
    recurrence {
      weekly {
        days_of_week = ["Monday"]
      }
    }
  }

  tags = {
    ENV = "Test"
  }
}




`, r.template(data), data.RandomInteger)
}

func (r MonitorAlertProcessingRuleSuppressionResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_alert_processing_rule_suppression" "test" {
  name                = "acctest-moniter-%d"
  resource_group_name = azurerm_resource_group.test.name
  description         = "alertprocessingrule-test"
  scopes              = [azurerm_resource_group.test.id]
  enabled             = false

  condition {
    alert_context {
      operator = "Contains"
      values   = ["context1", "context2"]
    }

    alert_rule_id {
      operator = "Contains"
      values   = ["ruleId1", "ruleId2"]
    }

    alert_rule_name {
      operator = "DoesNotContain"
      values   = ["ruleName1", "ruleName2"]
    }

    description {
      operator = "DoesNotContain"
      values   = ["description1", "description2"]
    }

    monitor_condition {
      operator = "NotEquals"
      values   = ["Fired"]
    }

    monitor_service {
      operator = "Equals"
      values   = ["Data Box Gateway", "Resource Health", "Prometheus"]
    }

    severity {
      operator = "Equals"
      values   = ["Sev0", "Sev1", "Sev2"]
    }

    signal_type {
      operator = "Equals"
      values   = ["Metric", "Log"]
    }

    target_resource {
      operator = "Contains"
      values   = ["resourseId1", "resourceId2"]
    }

    target_resource_group {
      operator = "DoesNotContain"
      values   = ["rg1", "rg2"]
    }

    target_resource_type {
      operator = "Equals"
      values   = ["Microsoft.Compute/VirtualMachines", "microsoft.batch/batchaccounts"]
    }
  }

  schedule {
    effective_from  = "2022-01-01T01:02:03"
    effective_until = "2022-02-02T01:02:03"
    time_zone       = "Pacific Standard Time"
    recurrence {
      daily {
        start_time = "17:00:00"
        end_time   = "09:00:00"
      }
      weekly {
        days_of_week = ["Sunday", "Saturday"]
      }
      weekly {
        start_time   = "17:00:00"
        end_time     = "18:00:00"
        days_of_week = ["Monday"]
      }
      monthly {
        start_time    = "09:00:00"
        end_time      = "12:00:00"
        days_of_month = [1, 15]
      }
    }
  }

  tags = {
    ENV  = "Test"
    ENV2 = "Test2"
  }
}
`, r.template(data), data.RandomInteger)
}

func (MonitorAlertProcessingRuleSuppressionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-monitor-maprs-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

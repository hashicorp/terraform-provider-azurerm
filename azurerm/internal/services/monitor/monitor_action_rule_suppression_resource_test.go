package monitor_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/monitor/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type MonitorActionRuleSuppressionResource struct {
}

func TestAccMonitorActionRuleSuppression_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_rule_suppression", "test")
	r := MonitorActionRuleSuppressionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorActionRuleSuppression_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_rule_suppression", "test")
	r := MonitorActionRuleSuppressionResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccMonitorActionRuleSuppression_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_rule_suppression", "test")
	r := MonitorActionRuleSuppressionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorActionRuleSuppression_updateSuppressionConfig(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_action_rule_suppression", "test")
	r := MonitorActionRuleSuppressionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.dailyRecurrence(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.monthlyRecurrence(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t MonitorActionRuleSuppressionResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.ActionRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Monitor.ActionRulesClient.GetByName(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading action rule (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r MonitorActionRuleSuppressionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_action_rule_suppression" "test" {
  name                = "acctest-moniter-%d"
  resource_group_name = azurerm_resource_group.test.name

  suppression {
    recurrence_type = "Always"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r MonitorActionRuleSuppressionResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_action_rule_suppression" "import" {
  name                = azurerm_monitor_action_rule_suppression.test.name
  resource_group_name = azurerm_monitor_action_rule_suppression.test.resource_group_name

  suppression {
    recurrence_type = azurerm_monitor_action_rule_suppression.test.suppression.0.recurrence_type
  }
}
`, r.basic(data))
}

func (r MonitorActionRuleSuppressionResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_action_rule_suppression" "test" {
  name                = "acctest-moniter-%d"
  resource_group_name = azurerm_resource_group.test.name
  enabled             = false
  description         = "actionRule-test"

  scope {
    type         = "ResourceGroup"
    resource_ids = [azurerm_resource_group.test.id]
  }

  suppression {
    recurrence_type = "Weekly"

    schedule {
      start_date_utc = "2019-01-01T01:02:03Z"
      end_date_utc   = "2019-01-03T15:02:07Z"

      recurrence_weekly = ["Sunday", "Monday", "Friday", "Saturday"]
    }
  }

  condition {
    alert_context {
      operator = "Contains"
      values   = ["context1", "context2"]
    }

    alert_rule_id {
      operator = "Contains"
      values   = ["ruleId1", "ruleId2"]
    }

    description {
      operator = "DoesNotContain"
      values   = ["description1", "description2"]
    }

    monitor {
      operator = "NotEquals"
      values   = ["Fired"]
    }

    monitor_service {
      operator = "Equals"
      values   = ["Data Box Edge", "Data Box Gateway", "Resource Health"]
    }

    severity {
      operator = "Equals"
      values   = ["Sev0", "Sev1", "Sev2"]
    }

    target_resource_type {
      operator = "Equals"
      values   = ["Microsoft.Compute/VirtualMachines", "microsoft.batch/batchaccounts"]
    }
  }

  tags = {
    ENV = "Test"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r MonitorActionRuleSuppressionResource) dailyRecurrence(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_action_rule_suppression" "test" {
  name                = "acctest-moniter-%d"
  resource_group_name = azurerm_resource_group.test.name

  scope {
    type         = "ResourceGroup"
    resource_ids = [azurerm_resource_group.test.id]
  }

  suppression {
    recurrence_type = "Daily"

    schedule {
      start_date_utc = "2019-01-01T01:02:03Z"
      end_date_utc   = "2019-01-03T15:02:07Z"
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r MonitorActionRuleSuppressionResource) monthlyRecurrence(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_action_rule_suppression" "test" {
  name                = "acctest-moniter-%d"
  resource_group_name = azurerm_resource_group.test.name

  scope {
    type         = "ResourceGroup"
    resource_ids = [azurerm_resource_group.test.id]
  }

  suppression {
    recurrence_type = "Monthly"

    schedule {
      start_date_utc     = "2019-01-01T01:02:03Z"
      end_date_utc       = "2019-01-03T15:02:07Z"
      recurrence_monthly = [1, 2, 15, 30, 31]
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (MonitorActionRuleSuppressionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-monitor-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

package sentinel_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/preview/securityinsight/mgmt/2019-01-01-preview/securityinsight"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/sentinel/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type SentinelAlertRuleScheduledResource struct {
}

func TestAccSentinelAlertRuleScheduled_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_scheduled", "test")
	r := SentinelAlertRuleScheduledResource{}

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

func TestAccSentinelAlertRuleScheduled_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_scheduled", "test")
	r := SentinelAlertRuleScheduledResource{}

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

func TestAccSentinelAlertRuleScheduled_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_scheduled", "test")
	r := SentinelAlertRuleScheduledResource{}

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

func TestAccSentinelAlertRuleScheduled_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_scheduled", "test")
	r := SentinelAlertRuleScheduledResource{}

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

func TestAccSentinelAlertRuleScheduled_withAlertRuleTemplateGuid(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_scheduled", "test")
	r := SentinelAlertRuleScheduledResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.alertRuleTemplateGuid(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSentinelAlertRuleScheduled_updateEventGroupingSetting(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_scheduled", "test")
	r := SentinelAlertRuleScheduledResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.eventGroupingSetting(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateEventGroupingSetting(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t SentinelAlertRuleScheduledResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.AlertRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Sentinel.AlertRulesClient.Get(ctx, id.ResourceGroup, "Microsoft.OperationalInsights", id.WorkspaceName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading Sentinel Alert Rule Scheduled %q: %v", id, err)
	}

	rule, ok := resp.Value.(securityinsight.ScheduledAlertRule)
	if !ok {
		return nil, fmt.Errorf("the Alert Rule %q is not a Scheduled Alert Rule", id)
	}

	return utils.Bool(rule.ID != nil), nil
}

func (r SentinelAlertRuleScheduledResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_alert_rule_scheduled" "test" {
  name                       = "acctest-SentinelAlertRule-Sche-%d"
  log_analytics_workspace_id = azurerm_log_analytics_solution.test.workspace_resource_id
  display_name               = "Some Rule"
  severity                   = "High"
  query                      = <<QUERY
AzureActivity |
  where OperationName == "Create or Update Virtual Machine" or OperationName =="Create Deployment" |
  where ActivityStatus == "Succeeded" |
  make-series dcount(ResourceId) default=0 on EventSubmissionTimestamp in range(ago(7d), now(), 1d) by Caller
QUERY
}
`, r.template(data), data.RandomInteger)
}

func (r SentinelAlertRuleScheduledResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_alert_rule_scheduled" "test" {
  name                       = "acctest-SentinelAlertRule-Sche-%d"
  log_analytics_workspace_id = azurerm_log_analytics_solution.test.workspace_resource_id
  display_name               = "Updated Rule"
  description                = "Some Description"
  tactics                    = ["Collection", "CommandAndControl"]
  severity                   = "Low"
  enabled                    = false
  incident_configuration {
    create_incident = true
    grouping {
      enabled                 = true
      lookback_duration       = "P7D"
      reopen_closed_incidents = true
      entity_matching_method  = "Custom"
      group_by                = ["Account", "Host"]
    }
  }
  query                = <<QUERY
AzureActivity |
  where OperationName == "Create or Update Virtual Machine" or OperationName =="Create Deployment" |
  where ActivityStatus == "Succeeded" |
  make-series dcount(ResourceId) default=0 on EventSubmissionTimestamp in range(ago(3d), now(), 1d) by Caller
QUERY
  query_frequency      = "PT20M"
  query_period         = "PT40M"
  trigger_operator     = "Equal"
  trigger_threshold    = 5
  suppression_enabled  = true
  suppression_duration = "PT40M"
}
`, r.template(data), data.RandomInteger)
}

func (r SentinelAlertRuleScheduledResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_alert_rule_scheduled" "import" {
  name                       = azurerm_sentinel_alert_rule_scheduled.test.name
  log_analytics_workspace_id = azurerm_sentinel_alert_rule_scheduled.test.log_analytics_workspace_id
  display_name               = azurerm_sentinel_alert_rule_scheduled.test.display_name
  severity                   = azurerm_sentinel_alert_rule_scheduled.test.severity
  query                      = azurerm_sentinel_alert_rule_scheduled.test.query
}
`, r.basic(data))
}

func (r SentinelAlertRuleScheduledResource) alertRuleTemplateGuid(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_alert_rule_scheduled" "test" {
  name                       = "acctest-SentinelAlertRule-Sche-%d"
  log_analytics_workspace_id = azurerm_log_analytics_solution.test.workspace_resource_id
  display_name               = "Some Rule"
  severity                   = "Low"
  alert_rule_template_guid   = "65360bb0-8986-4ade-a89d-af3cf44d28aa"
  query                      = <<QUERY
AzureActivity |
  where OperationName == "Create or Update Virtual Machine" or OperationName =="Create Deployment" |
  where ActivityStatus == "Succeeded" |
  make-series dcount(ResourceId) default=0 on EventSubmissionTimestamp in range(ago(7d), now(), 1d) by Caller
QUERY
}
`, r.template(data), data.RandomInteger)
}

func (r SentinelAlertRuleScheduledResource) eventGroupingSetting(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_alert_rule_scheduled" "test" {
  name                       = "acctest-SentinelAlertRule-Sche-%d"
  log_analytics_workspace_id = azurerm_log_analytics_solution.test.workspace_resource_id
  display_name               = "Some Rule"
  severity                   = "Low"
  alert_rule_template_guid   = "65360bb0-8986-4ade-a89d-af3cf44d28aa"
  query                      = <<QUERY
AzureActivity |
  where OperationName == "Create or Update Virtual Machine" or OperationName =="Create Deployment" |
  where ActivityStatus == "Succeeded" |
  make-series dcount(ResourceId) default=0 on EventSubmissionTimestamp in range(ago(7d), now(), 1d) by Caller
QUERY

  event_grouping {
    aggregation_method = "SingleAlert"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r SentinelAlertRuleScheduledResource) updateEventGroupingSetting(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_alert_rule_scheduled" "test" {
  name                       = "acctest-SentinelAlertRule-Sche-%d"
  log_analytics_workspace_id = azurerm_log_analytics_solution.test.workspace_resource_id
  display_name               = "Some Rule"
  severity                   = "Low"
  alert_rule_template_guid   = "65360bb0-8986-4ade-a89d-af3cf44d28aa"
  query                      = <<QUERY
AzureActivity |
  where OperationName == "Create or Update Virtual Machine" or OperationName =="Create Deployment" |
  where ActivityStatus == "Succeeded" |
  make-series dcount(ResourceId) default=0 on EventSubmissionTimestamp in range(ago(7d), now(), 1d) by Caller
QUERY

  event_grouping {
    aggregation_method = "AlertPerResult"
  }
}
`, r.template(data), data.RandomInteger)
}

func (SentinelAlertRuleScheduledResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sentinel-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}

resource "azurerm_log_analytics_solution" "test" {
  solution_name         = "SecurityInsights"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  workspace_resource_id = azurerm_log_analytics_workspace.test.id
  workspace_name        = azurerm_log_analytics_workspace.test.name

  plan {
    publisher = "Microsoft"
    product   = "OMSGallery/SecurityInsights"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

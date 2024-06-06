// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-15-preview/scheduledqueryrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MonitorScheduledQueryRulesAlertV2Resource struct{}

func TestAccMonitorScheduledQueryRulesAlertV2_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_scheduled_query_rules_alert_v2", "test")
	r := MonitorScheduledQueryRulesAlertV2Resource{}
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

func TestAccMonitorScheduledQueryRulesAlertV2_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_scheduled_query_rules_alert_v2", "test")
	r := MonitorScheduledQueryRulesAlertV2Resource{}
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

func TestAccMonitorScheduledQueryRulesAlertV2_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_scheduled_query_rules_alert_v2", "test")
	r := MonitorScheduledQueryRulesAlertV2Resource{}
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

func TestAccMonitorScheduledQueryRulesAlertV2_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_scheduled_query_rules_alert_v2", "test")
	r := MonitorScheduledQueryRulesAlertV2Resource{}
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
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorScheduledQueryRulesAlertV2_identitySystemAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_scheduled_query_rules_alert_v2", "test")
	r := MonitorScheduledQueryRulesAlertV2Resource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.identitySystemAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").IsUUID(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorScheduledQueryRulesAlertV2_identityUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_scheduled_query_rules_alert_v2", "test")
	r := MonitorScheduledQueryRulesAlertV2Resource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.#").HasValue("0"),
			),
		},
		data.ImportStep(),
		{
			Config: r.identitySystemAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").IsUUID(),
			),
		},
		data.ImportStep(),
		{
			Config: r.identityUserAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("UserAssigned"),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsEmpty(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").IsUUID(),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.#").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func (r MonitorScheduledQueryRulesAlertV2Resource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := scheduledqueryrules.ParseScheduledQueryRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Monitor.ScheduledQueryRulesV2Client
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r MonitorScheduledQueryRulesAlertV2Resource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctest-ai-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctest-mag-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "test mag"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestUAI-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_application_insights.test.id
  role_definition_name = "Reader"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r MonitorScheduledQueryRulesAlertV2Resource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_scheduled_query_rules_alert_v2" "test" {
  name                 = "acctest-isqr-%d"
  resource_group_name  = azurerm_resource_group.test.name
  location             = "%s"
  evaluation_frequency = "PT5M"
  window_duration      = "PT5M"
  scopes               = [azurerm_application_insights.test.id]
  severity             = 3
  criteria {
    query                   = <<-QUERY
      requests
	    | summarize CountByCountry=count() by client_CountryOrRegion
	  QUERY
    time_aggregation_method = "Count"
    threshold               = 5.0
    operator                = "Equal"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r MonitorScheduledQueryRulesAlertV2Resource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_scheduled_query_rules_alert_v2" "import" {
  name                 = azurerm_monitor_scheduled_query_rules_alert_v2.test.name
  resource_group_name  = azurerm_resource_group.test.name
  location             = "%s"
  evaluation_frequency = azurerm_monitor_scheduled_query_rules_alert_v2.test.evaluation_frequency
  window_duration      = azurerm_monitor_scheduled_query_rules_alert_v2.test.window_duration
  scopes               = azurerm_monitor_scheduled_query_rules_alert_v2.test.scopes
  severity             = azurerm_monitor_scheduled_query_rules_alert_v2.test.severity
  criteria {
    query                   = azurerm_monitor_scheduled_query_rules_alert_v2.test.criteria.0.query
    time_aggregation_method = azurerm_monitor_scheduled_query_rules_alert_v2.test.criteria.0.time_aggregation_method
    threshold               = azurerm_monitor_scheduled_query_rules_alert_v2.test.criteria.0.threshold
    operator                = azurerm_monitor_scheduled_query_rules_alert_v2.test.criteria.0.operator
  }
}
`, config, data.Locations.Primary)
}

func (r MonitorScheduledQueryRulesAlertV2Resource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_scheduled_query_rules_alert_v2" "test" {
  name                = "acctest-isqr-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"

  evaluation_frequency = "PT5M"
  window_duration      = "PT5M"
  scopes               = [azurerm_application_insights.test.id]
  severity             = 3
  criteria {
    query                   = <<-QUERY
      requests
	    | summarize CountByCountry=count() by client_CountryOrRegion
	  QUERY
    time_aggregation_method = "Count"
    threshold               = 5.0
    operator                = "GreaterThan"

    resource_id_column = "client_CountryOrRegion"
    dimension {
      name     = "client_CountryOrRegion"
      operator = "Include"
      values   = ["*"]
    }
    failing_periods {
      minimum_failing_periods_to_trigger_alert = 1
      number_of_evaluation_periods             = 1
    }
  }

  auto_mitigation_enabled           = false
  workspace_alerts_storage_enabled  = false
  description                       = "test sqr"
  display_name                      = "acctest-sqr"
  enabled                           = false
  mute_actions_after_alert_duration = "PT10M"
  query_time_range_override         = "PT10M"
  skip_query_validation             = false
  target_resource_types             = ["microsoft.insights/components"]
  action {
    action_groups = [azurerm_monitor_action_group.test.id]
    custom_properties = {
      key = "value"
    }
  }

  tags = {
    key = "value"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r MonitorScheduledQueryRulesAlertV2Resource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_scheduled_query_rules_alert_v2" "test" {
  name                = "acctest-isqr-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"

  evaluation_frequency = "PT10M"
  window_duration      = "PT10M"
  scopes               = [azurerm_application_insights.test.id]
  severity             = 4
  criteria {
    query                   = <<-QUERY
      requests
        | summarize CountByCountry=count() by client_CountryOrRegion
      QUERY
    time_aggregation_method = "Maximum"
    threshold               = 17.5
    operator                = "LessThan"

    resource_id_column    = "client_CountryOrRegion"
    metric_measure_column = "CountByCountry"
    dimension {
      name     = "client_CountryOrRegion"
      operator = "Exclude"
      values   = ["123"]
    }
    failing_periods {
      minimum_failing_periods_to_trigger_alert = 1
      number_of_evaluation_periods             = 1
    }
  }

  auto_mitigation_enabled          = true
  workspace_alerts_storage_enabled = false
  description                      = "test sqr"
  display_name                     = "acctest-sqr"
  enabled                          = true
  query_time_range_override        = "PT1H"
  skip_query_validation            = true
  action {
    action_groups = [azurerm_monitor_action_group.test.id]
    custom_properties = {
      key  = "value"
      key2 = "value2"
    }
  }

  tags = {
    key  = "value"
    key2 = "value2"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r MonitorScheduledQueryRulesAlertV2Resource) identitySystemAssigned(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_scheduled_query_rules_alert_v2" "test" {
  name                 = "acctest-isqr-%d"
  resource_group_name  = azurerm_resource_group.test.name
  location             = "%s"
  evaluation_frequency = "PT5M"
  window_duration      = "PT5M"
  scopes               = [azurerm_application_insights.test.id]
  severity             = 3
  criteria {
    query                   = <<-QUERY
      requests
	    | summarize CountByCountry=count() by client_CountryOrRegion
	  QUERY
    time_aggregation_method = "Count"
    threshold               = 5.0
    operator                = "Equal"
  }
  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r MonitorScheduledQueryRulesAlertV2Resource) identityUserAssigned(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_scheduled_query_rules_alert_v2" "test" {
  name                 = "acctest-isqr-%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  location             = "%s"
  evaluation_frequency = "PT5M"
  window_duration      = "PT5M"
  scopes               = [azurerm_application_insights.test.id]
  severity             = 3
  criteria {
    query                   = <<-QUERY
      requests
	    | summarize CountByCountry=count() by client_CountryOrRegion
	  QUERY
    time_aggregation_method = "Count"
    threshold               = 5.0
    operator                = "Equal"
  }
  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }
  depends_on = [azurerm_role_assignment.test]
}
`, template, data.RandomInteger, data.Locations.Primary)
}

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loganalytics_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/tables"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LogAnalyticsWorkspaceTableCustomLogResource struct{}

func TestAccLogAnalyticsWorkspaceTableCustomLog_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace_table_custom_log", "test")
	r := LogAnalyticsWorkspaceTableCustomLogResource{}

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

func TestAccLogAnalyticsWorkspaceTableCustomLog_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace_table_custom_log", "test")
	r := LogAnalyticsWorkspaceTableCustomLogResource{}

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

func TestAccLogAnalyticsWorkspaceTableCustomLog_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace_table_custom_log", "test")
	r := LogAnalyticsWorkspaceTableCustomLogResource{}

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

func TestAccLogAnalyticsWorkspaceTableCustomLog_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace_table_custom_log", "test")
	r := LogAnalyticsWorkspaceTableCustomLogResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicUpdate(data),
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

func TestAccLogAnalyticsWorkspaceTableCustomLog_planBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace_table_custom_log", "test")
	r := LogAnalyticsWorkspaceTableCustomLogResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicPlan(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config:      r.basicPlanWithRetentionInDays(data),
			PlanOnly:    true,
			ExpectError: regexp.MustCompile("`retention_in_days` cannot be set when `plan` is set to `Basic`"),
		},
	})
}

func (r LogAnalyticsWorkspaceTableCustomLogResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := tables.ParseTableID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.LogAnalytics.TablesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r LogAnalyticsWorkspaceTableCustomLogResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace_table_custom_log" "test" {
  name         = "acctestlawdcr%d_CL"
  workspace_id = azurerm_log_analytics_workspace.test.id

  column {
    name = "TimeGenerated"
    type = "dateTime"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r LogAnalyticsWorkspaceTableCustomLogResource) basicUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace_table_custom_log" "test" {
  name         = "acctestlawdcr%d_CL"
  workspace_id = azurerm_log_analytics_workspace.test.id

  display_name = "Basic Update"

  column {
    name = "TimeGenerated"
    type = "dateTime"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r LogAnalyticsWorkspaceTableCustomLogResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace_table_custom_log" "import" {
  name         = azurerm_log_analytics_workspace_table_custom_log.test.name
  workspace_id = azurerm_log_analytics_workspace_table_custom_log.test.workspace_id

  column {
    name = azurerm_log_analytics_workspace_table_custom_log.test.column.0.name
    type = azurerm_log_analytics_workspace_table_custom_log.test.column.0.type
  }
}
`, r.basic(data))
}

func (r LogAnalyticsWorkspaceTableCustomLogResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace_table_custom_log" "test" {
  name                    = "acctestlawdcr%d_CL"
  workspace_id            = azurerm_log_analytics_workspace.test.id
  display_name            = "Test Custom Log Table"
  description             = "This is a test custom log table"
  plan                    = "Analytics"
  total_retention_in_days = 60
  retention_in_days       = 20

  column {
    display_name = "TimeGenerated"
    description  = "The timestamp when the log was generated"
    name         = "TimeGenerated"
    type         = "dateTime"
  }

  column {
    display_name = "ApplicationName"
    description  = "The name of the application"
    name         = "Application"
    type         = "string"
  }

  column {
    display_name = "RawLogData"
    description  = "The raw log data content"
    name         = "RawData"
    type         = "string"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r LogAnalyticsWorkspaceTableCustomLogResource) basicPlan(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace_table_custom_log" "test" {
  name         = "acctestlawdcr%d_CL"
  workspace_id = azurerm_log_analytics_workspace.test.id

  plan = "Basic"

  column {
    name = "TimeGenerated"
    type = "dateTime"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r LogAnalyticsWorkspaceTableCustomLogResource) basicPlanWithRetentionInDays(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace_table_custom_log" "test" {
  name         = "acctestlawdcr%d_CL"
  workspace_id = azurerm_log_analytics_workspace.test.id

  plan              = "Basic"
  retention_in_days = 60

  column {
    name = "TimeGenerated"
    type = "dateTime"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r LogAnalyticsWorkspaceTableCustomLogResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  retention_in_days   = 30
}

resource "azurerm_monitor_data_collection_endpoint" "test" {
  name                = "acctestdce-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_monitor_data_collection_rule" "test" {
  name                        = "acctestdcr-%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  data_collection_endpoint_id = azurerm_monitor_data_collection_endpoint.test.id

  data_flow {
    destinations  = [replace(azurerm_log_analytics_workspace.test.workspace_id, "-", "")]
    output_stream = "Custom-${azurerm_log_analytics_workspace_table_custom_log.test.name}"
    streams       = ["Custom-${azurerm_log_analytics_workspace_table_custom_log.test.name}"]
    transform_kql = "source"
  }

  stream_declaration {
    stream_name = "Custom-${azurerm_log_analytics_workspace_table_custom_log.test.name}"
    column {
      name = "TimeGenerated"
      type = "datetime"
    }

    column {
      name = "Application"
      type = "string"
    }

    column {
      name = "RawData"
      type = "string"
    }
  }

  destinations {
    log_analytics {
      name                  = replace(azurerm_log_analytics_workspace.test.workspace_id, "-", "")
      workspace_resource_id = azurerm_log_analytics_workspace.test.id
    }
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

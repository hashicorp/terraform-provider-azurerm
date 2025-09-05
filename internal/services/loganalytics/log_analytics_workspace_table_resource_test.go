// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loganalytics_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/tables"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LogAnalyticsWorkspaceTableResource struct{}

func TestAccLogAnalyticsWorkspaceTable_updateTableRetention(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")
	r := LogAnalyticsWorkspaceTableResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.base(data),
		},
		{
			Config:      r.updateRetentionImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_log_analytics_workspace_table"),
		},
		data.ImportStep(),
		{
			Config: r.updateRetentionUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_log_analytics_workspace_table.test").ExistsInAzure(r),
				check.That("azurerm_log_analytics_workspace_table.test").Key("id").Exists(),
				check.That("azurerm_log_analytics_workspace_table.test").Key("name").HasValue("AppEvents"),
				check.That("azurerm_log_analytics_workspace_table.test").Key("retention_in_days").HasValue("7"),
				check.That("azurerm_log_analytics_workspace_table.test").Key("total_retention_in_days").HasValue("32"),
			),
		},
		{
			Config: r.removeRetentionDays(data),
			// The `retention_in_days` is an optional property, when it's not specified, the service will set a default value for it.
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccLogAnalyticsWorkspaceTable_plan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace", "test")
	r := LogAnalyticsWorkspaceTableResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.base(data),
		},
		{
			Config:      r.planImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_log_analytics_workspace_table"),
		},
		data.ImportStep(),
		{
			Config: r.planUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_log_analytics_workspace_table.test").ExistsInAzure(r),
			),
		},
	})
}

func TestAccLogAnalyticsWorkspaceTable_customDcr(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace_table", "test")
	r := LogAnalyticsWorkspaceTableResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.custom(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("name").HasValue("CustomTable_CL"),
				check.That(data.ResourceName).Key("column.0.name").HasValue("CompanyName"),
				check.That(data.ResourceName).Key("retention_in_days").HasValue("7"),
				check.That(data.ResourceName).Key("total_retention_in_days").HasValue("32"),
			),
		},
		data.ImportStep(),
		{
			Config: r.customUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("name").HasValue("CustomTable_CL"),
				check.That(data.ResourceName).Key("column.0.name").HasValue("LogName"),
				check.That(data.ResourceName).Key("retention_in_days").HasValue("0"),
				check.That(data.ResourceName).Key("total_retention_in_days").HasValue("30"),
			),
		},
		data.ImportStep(),
	})
}

func (t LogAnalyticsWorkspaceTableResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := tables.ParseTableID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.LogAnalytics.TablesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading Log Analytics Workspace Table (%s): %+v", id.ID(), err)
	}

	return utils.Bool(resp.Model.Id != nil), nil
}

func (t LogAnalyticsWorkspaceTableResource) updateRetentionImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
import {
  id = "${azurerm_log_analytics_workspace.test.id}/tables/AppEvents"
  to = azurerm_log_analytics_workspace_table.test
}
resource "azurerm_log_analytics_workspace_table" "test" {
  name                    = "AppEvents"
  type                    = "Microsoft"
  sub_type                = "DataCollectionRuleBased"
  workspace_id            = azurerm_log_analytics_workspace.test.id
  total_retention_in_days = 90
}
`, t.base(data))
}

func (t LogAnalyticsWorkspaceTableResource) updateRetentionUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_log_analytics_workspace_table" "test" {
  name                    = "AppEvents"
  type                    = "Microsoft"
  sub_type                = "DataCollectionRuleBased"
  workspace_id            = azurerm_log_analytics_workspace.test.id
  retention_in_days       = 7
  total_retention_in_days = 32
}
`, t.base(data))
}

func (t LogAnalyticsWorkspaceTableResource) planImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
import {
  id = "${azurerm_log_analytics_workspace.test.id}/tables/AppTraces"
  to = azurerm_log_analytics_workspace_table.test
}
resource "azurerm_log_analytics_workspace_table" "test" {
  name                    = "AppTraces"
  type                    = "Microsoft"
  sub_type                = "DataCollectionRuleBased"
  workspace_id            = azurerm_log_analytics_workspace.test.id
  plan                    = "Analytics"
  total_retention_in_days = 90
}
`, t.base(data))
}

func (t LogAnalyticsWorkspaceTableResource) planUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_log_analytics_workspace_table" "test" {
  name                    = "AppTraces"
  type                    = "Microsoft"
  sub_type                = "DataCollectionRuleBased"
  workspace_id            = azurerm_log_analytics_workspace.test.id
  plan                    = "Basic"
  total_retention_in_days = 90
}
`, t.base(data))
}

func (t LogAnalyticsWorkspaceTableResource) custom(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_log_analytics_workspace_table" "test" {
  name                    = "CustomTable_CL"
  type                    = "CustomLog"
  sub_type                = "DataCollectionRuleBased"
  workspace_id            = azurerm_log_analytics_workspace.test.id
  retention_in_days       = 7
  total_retention_in_days = 32

  column {
    name = "CompanyName"
    type = "string"
  }
  column {
    name = "TimeGenerated"
    type = "dateTime"
  }
}
`, t.base(data))
}

func (t LogAnalyticsWorkspaceTableResource) customUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_log_analytics_workspace_table" "test" {
  name                    = "CustomTable_CL"
  type                    = "CustomLog"
  sub_type                = "DataCollectionRuleBased"
  workspace_id            = azurerm_log_analytics_workspace.test.id
  total_retention_in_days = 30

  column {
    name = "LogName"
    type = "string"
  }
  column {
    name = "TimeGenerated"
    type = "dateTime"
  }
}
`, t.base(data))
}

func (t LogAnalyticsWorkspaceTableResource) base(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}
resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  retention_in_days   = 30
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (LogAnalyticsWorkspaceTableResource) removeRetentionDays(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}
resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  retention_in_days   = 30
}
resource "azurerm_log_analytics_workspace_table" "test" {
  name                    = "AppEvents"
  workspace_id            = azurerm_log_analytics_workspace.test.id
  total_retention_in_days = 180
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

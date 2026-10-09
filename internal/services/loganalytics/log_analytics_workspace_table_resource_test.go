// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package loganalytics_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/tables"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LogAnalyticsWorkspaceTableResource struct{}

func TestAccLogAnalyticsWorkspaceTable_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace_table", "test")
	r := LogAnalyticsWorkspaceTableResource{}

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

func TestAccLogAnalyticsWorkspaceTable_updateTableRetention(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace_table", "test")
	r := LogAnalyticsWorkspaceTableResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.retention(data, 10),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.retention(data, 7),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.totalRetention(data, 180),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.retention(data, 10),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsWorkspaceTable_updateTableTotalRetention(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace_table", "test")
	r := LogAnalyticsWorkspaceTableResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.totalRetention(data, 180),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.totalRetention(data, 120),
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
		{
			Config: r.totalRetention(data, 180),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsWorkspaceTable_plan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace_table", "test")
	r := LogAnalyticsWorkspaceTableResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.plan(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.planWithTotalRetention(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.plan(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
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

	return pointer.To(resp.Model != nil), nil
}

func (r LogAnalyticsWorkspaceTableResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace_table" "test" {
  name         = "AppEvents"
  workspace_id = azurerm_log_analytics_workspace.test.id
}
`, r.template(data))
}

func (r LogAnalyticsWorkspaceTableResource) retention(data acceptance.TestData, retentionInDays int) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace_table" "test" {
  name                    = "AppEvents"
  workspace_id            = azurerm_log_analytics_workspace.test.id
  retention_in_days       = %d
  total_retention_in_days = 180
}
`, r.template(data), retentionInDays)
}

func (r LogAnalyticsWorkspaceTableResource) totalRetention(data acceptance.TestData, totalRetentionInDays int) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace_table" "test" {
  name                    = "AppEvents"
  workspace_id            = azurerm_log_analytics_workspace.test.id
  total_retention_in_days = %d
}
`, r.template(data), totalRetentionInDays)
}

func (r LogAnalyticsWorkspaceTableResource) plan(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace_table" "test" {
  name         = "AppTraces"
  workspace_id = azurerm_log_analytics_workspace.test.id
  plan         = "Basic"
}
`, r.template(data))
}

func (r LogAnalyticsWorkspaceTableResource) planWithTotalRetention(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace_table" "test" {
  name                    = "AppTraces"
  workspace_id            = azurerm_log_analytics_workspace.test.id
  plan                    = "Basic"
  total_retention_in_days = 45
}
`, r.template(data))
}

func (LogAnalyticsWorkspaceTableResource) template(data acceptance.TestData) string {
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

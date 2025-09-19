// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loganalytics_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/tables"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LogAnalyticsWorkspaceTableMicrosoftResource struct{}

func TestAccLogAnalyticsWorkspaceTableMicrosoft_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace_table_microsoft", "test")
	r := LogAnalyticsWorkspaceTableMicrosoftResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.base(data),
		},
		{
			Config:      r.updateRetentionImport(data),
			ExpectError: acceptance.RequiresImportError(data.ResourceType),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsWorkspaceTableMicrosoft_updateTableRetention(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace_table_microsoft", "test")
	r := LogAnalyticsWorkspaceTableMicrosoftResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.base(data),
		},
		{
			Config:      r.updateRetentionImport(data),
			ExpectError: acceptance.RequiresImportError(data.ResourceType),
		},
		data.ImportStep(),
		{
			Config: r.updateRetentionUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_log_analytics_workspace_table_microsoft.test").ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.removeRetentionDays(data),
			// The `retention_in_days` is an optional property, when it's not specified, the service will set a default value for it.
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccLogAnalyticsWorkspaceTableMicrosoft_plan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace_table_microsoft", "test")
	r := LogAnalyticsWorkspaceTableMicrosoftResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.base(data),
			SkipFunc: func() (bool, error) {
				fmt.Println(os.WriteFile("/Users/wyatt.fry/tests/TestAccLogAnalyticsWorkspaceTableMicrosoft_plan.tf", []byte(r.base(data)), 0666))
				return true, nil
			},
		},
		{
			Config:      r.planImport(data),
			ExpectError: acceptance.RequiresImportError(data.ResourceType),
		},
		data.ImportStep(),
		{
			Config: r.planUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_log_analytics_workspace_table_microsoft.test").ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t LogAnalyticsWorkspaceTableMicrosoftResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := tables.ParseTableID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.LogAnalytics.TablesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading Log Analytics Workspace Table (%s): %+v", id.ID(), err)
	}

	return pointer.To(resp.Model.Id != nil), nil
}

func (t LogAnalyticsWorkspaceTableMicrosoftResource) updateRetentionImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
import {
  id = "${azurerm_log_analytics_workspace.test.id}/tables/AppEvents"
  to = azurerm_log_analytics_workspace_table_microsoft.test
}
resource "azurerm_log_analytics_workspace_table_microsoft" "test" {
  name                    = "AppEvents"
  sub_type = "Any"
  workspace_id            = azurerm_log_analytics_workspace.test.id
  total_retention_in_days = 90
}
`, t.base(data))
}

func (t LogAnalyticsWorkspaceTableMicrosoftResource) updateRetentionUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_log_analytics_workspace_table_microsoft" "test" {
  name                    = "AppEvents"
  sub_type = "Any"
  workspace_id            = azurerm_log_analytics_workspace.test.id
  retention_in_days       = 7
  total_retention_in_days = 32
}
`, t.base(data))
}

func (t LogAnalyticsWorkspaceTableMicrosoftResource) planImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
import {
  id = "${azurerm_log_analytics_workspace.test.id}/tables/Alert"
  to = azurerm_log_analytics_workspace_table_microsoft.test
}
resource "azurerm_log_analytics_workspace_table_microsoft" "test" {
  name                    = "Alert"
  display_name            = "Alert"
  sub_type                = "Any"
  workspace_id            = azurerm_log_analytics_workspace.test.id
  plan                    = "Analytics"
  total_retention_in_days = 30
}
`, t.base(data))
}

func (t LogAnalyticsWorkspaceTableMicrosoftResource) planUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
import {
  id = "${azurerm_log_analytics_workspace.test.id}/tables/AppTraces"
  to = azurerm_log_analytics_workspace_table_microsoft.test
}
resource "azurerm_log_analytics_workspace_table_microsoft" "test" {
  name                    = "AppTraces"
  sub_type = "Any"
  workspace_id            = azurerm_log_analytics_workspace.test.id
  plan                    = "Basic"
  total_retention_in_days = 90
}
`, t.base(data))
}

func (t LogAnalyticsWorkspaceTableMicrosoftResource) base(data acceptance.TestData) string {
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

func (t LogAnalyticsWorkspaceTableMicrosoftResource) removeRetentionDays(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_log_analytics_workspace_table_microsoft" "test" {
  name                    = "AppEvents"
  sub_type = "Any"
  workspace_id            = azurerm_log_analytics_workspace.test.id
  # retention_in_days       = 7
  total_retention_in_days = 90
}
`, t.base(data))
}

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
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace_table", "test")
	r := LogAnalyticsWorkspaceTableResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.updateRetention(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("name").HasValue("AppEvents"),
				check.That(data.ResourceName).Key("retention_in_days").HasValue("7"),
				check.That(data.ResourceName).Key("total_retention_in_days").HasValue("32"),
			),
		},
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

func (LogAnalyticsWorkspaceTableResource) updateRetention(data acceptance.TestData) string {
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
  retention_in_days       = 7
  total_retention_in_days = 32
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (LogAnalyticsWorkspaceTableResource) plan(data acceptance.TestData) string {
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
  name                    = "AppTraces"
  workspace_id            = azurerm_log_analytics_workspace.test.id
  plan                    = "Basic"
  total_retention_in_days = 32
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

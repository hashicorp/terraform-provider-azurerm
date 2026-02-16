// Copyright (c) HashiCorp, Inc.
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

type LogAnalyticsWorkspaceTableMicrosoftResource struct{}

func TestAccLogAnalyticsWorkspaceTableMicrosoft_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace_table_microsoft", "test")
	r := LogAnalyticsWorkspaceTableMicrosoftResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_log_analytics_workspace_table_microsoft.test").ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsWorkspaceTableMicrosoft_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace_table_microsoft", "test")
	r := LogAnalyticsWorkspaceTableMicrosoftResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_log_analytics_workspace_table_microsoft.test").ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsWorkspaceTableMicrosoft_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace_table_microsoft", "test")
	r := LogAnalyticsWorkspaceTableMicrosoftResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_log_analytics_workspace_table_microsoft.test").ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_log_analytics_workspace_table_microsoft.test").ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_log_analytics_workspace_table_microsoft.test").ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateColumns(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_log_analytics_workspace_table_microsoft.test").ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
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
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model.Id != nil), nil
}

func (t LogAnalyticsWorkspaceTableMicrosoftResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace_table_microsoft" "test" {
  name         = "AppCenterError"
  workspace_id = azurerm_log_analytics_workspace.test.id
}
`, t.template(data))
}

func (t LogAnalyticsWorkspaceTableMicrosoftResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace_table_microsoft" "test" {
  name                    = "AppCenterError"
  display_name            = "AppCenterError"
  workspace_id            = azurerm_log_analytics_workspace.test.id
  total_retention_in_days = 30
  retention_in_days       = 30
  description             = "description"
  labels                  = ["label1", "label2"]
  column {
    name               = "column1_CF"
    display_name       = "Column1"
    description        = "description"
    type               = "string"
    display_by_default = false
    hidden             = true
  }
}
`, t.template(data))
}

func (t LogAnalyticsWorkspaceTableMicrosoftResource) updated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace_table_microsoft" "test" {
  name                    = "AppCenterError"
  display_name            = "AppCenterError"
  workspace_id            = azurerm_log_analytics_workspace.test.id
  total_retention_in_days = 25
  retention_in_days       = 25
  description             = "description updated"
  labels                  = ["label1", "label2", "label3"]
  column {
    name               = "column1updated_CF"
    display_name       = "Column1Updated"
    description        = "description updated"
    type               = "string"
    display_by_default = true
    hidden             = false
  }
  column {
    name               = "column2_CF"
    description        = "description 2"
    type               = "string"
    display_by_default = false
    hidden             = false
  }
}
`, t.template(data))
}

func (t LogAnalyticsWorkspaceTableMicrosoftResource) updateColumns(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace_table_microsoft" "test" {
  name                    = "AppCenterError"
  display_name            = "AppCenterError"
  workspace_id            = azurerm_log_analytics_workspace.test.id
  total_retention_in_days = 25
  retention_in_days       = 25
  description             = "description updated"
  labels                  = ["label1", "label2", "label3"]
  column {
    name         = "column1updated_CF"
    display_name = "Column1Updated"
    type         = "dynamic"
  }
}
`, t.template(data))
}

func (t LogAnalyticsWorkspaceTableMicrosoftResource) template(data acceptance.TestData) string {
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

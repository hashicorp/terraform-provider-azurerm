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
			Config: r.template(data),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_log_analytics_workspace_table_microsoft.test").ExistsInAzure(r),
			),
		},
		data.ImportStep("categories"),
	})
}

func TestAccLogAnalyticsWorkspaceTableMicrosoft_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace_table_microsoft", "test")
	r := LogAnalyticsWorkspaceTableMicrosoftResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.template(data),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_log_analytics_workspace_table_microsoft"),
		},
	})
}

func TestAccLogAnalyticsWorkspaceTableMicrosoft_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace_table_microsoft", "test")
	r := LogAnalyticsWorkspaceTableMicrosoftResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.template(data),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_log_analytics_workspace_table_microsoft.test").ExistsInAzure(r),
			),
		},
		data.ImportStep("categories"),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_log_analytics_workspace_table_microsoft.test").ExistsInAzure(r),
			),
		},
		data.ImportStep("categories"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_log_analytics_workspace_table_microsoft.test").ExistsInAzure(r),
			),
		},
		data.ImportStep("categories"),
	})
}

func TestAccLogAnalyticsWorkspaceTableMicrosoft_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_workspace_table_microsoft", "test")
	r := LogAnalyticsWorkspaceTableMicrosoftResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.template(data),
		},
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_log_analytics_workspace_table_microsoft.test").ExistsInAzure(r),
			),
		},
		data.ImportStep("categories"),
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

func (t LogAnalyticsWorkspaceTableMicrosoftResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
import {
  id = "${azurerm_log_analytics_workspace.test.id}/tables/AppCenterError"
  to = azurerm_log_analytics_workspace_table_microsoft.test
}
resource "azurerm_log_analytics_workspace_table_microsoft" "test" {
  name         = "AppCenterError"
  display_name = "AppCenterError"
  sub_type     = "Any"
  workspace_id = azurerm_log_analytics_workspace.test.id
}
`, t.template(data))
}

func (r LogAnalyticsWorkspaceTableMicrosoftResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_workspace_table_microsoft" "import" {
  name         = azurerm_log_analytics_workspace_table_microsoft.test.name
  display_name = azurerm_log_analytics_workspace_table_microsoft.test.display_name
  sub_type     = azurerm_log_analytics_workspace_table_microsoft.test.sub_type
  workspace_id = azurerm_log_analytics_workspace_table_microsoft.test.workspace_id
}
`, r.basic(data))
}

func (t LogAnalyticsWorkspaceTableMicrosoftResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
import {
  id = "${azurerm_log_analytics_workspace.test.id}/tables/AppCenterError"
  to = azurerm_log_analytics_workspace_table_microsoft.test
}
resource "azurerm_log_analytics_workspace_table_microsoft" "test" {
  name                    = "AppCenterError"
  display_name            = "AppCenterError"
  sub_type                = "Any"
  workspace_id            = azurerm_log_analytics_workspace.test.id
  total_retention_in_days = 30
  retention_in_days       = 30
  description             = "This is a description"
  plan                    = "Analytics"
  categories              = ["AppCenter"] # NOT RETURNED BY API
  labels                  = ["label1", "label2"]
  column {
    name               = "mycustom_CF"
    description        = "test"
    type               = "string"
    display_by_default = false
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

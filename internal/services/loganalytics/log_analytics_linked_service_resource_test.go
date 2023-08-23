// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loganalytics_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/linkedservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LogAnalyticsLinkedServiceResource struct{}

func TestAccLogAnalyticsLinkedService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_linked_service", "test")
	r := LogAnalyticsLinkedServiceResource{}

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

func TestAccLogAnalyticsLinkedService_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_linked_service", "test")
	r := LogAnalyticsLinkedServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_log_analytics_linked_service"),
		},
	})
}

func TestAccLogAnalyticsLinkedService_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_linked_service", "test")
	r := LogAnalyticsLinkedServiceResource{}

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

func TestAccLogAnalyticsLinkedService_withWriteAccessResourceId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_linked_service", "test")
	r := LogAnalyticsLinkedServiceResource{}

	if os.Getenv("ARM_RUN_TEST_LOG_ANALYTICS_CLUSTERS") == "" {
		t.Skip("Skipping as ARM_RUN_TEST_LOG_ANALYTICS_CLUSTERS is not specified")
		return
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withWriteAccessResourceId(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r LogAnalyticsLinkedServiceResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := linkedservices.ParseLinkedServiceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.LogAnalytics.LinkedServicesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r LogAnalyticsLinkedServiceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_linked_service" "test" {
  resource_group_name = azurerm_resource_group.test.name
  workspace_id        = azurerm_log_analytics_workspace.test.id
  read_access_id      = azurerm_automation_account.test.id
}
`, r.template(data))
}

func (r LogAnalyticsLinkedServiceResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_linked_service" "import" {
  resource_group_name = azurerm_log_analytics_linked_service.test.resource_group_name
  workspace_id        = azurerm_log_analytics_linked_service.test.workspace_id
  read_access_id      = azurerm_log_analytics_linked_service.test.read_access_id
}
`, r.basic(data))
}

func (r LogAnalyticsLinkedServiceResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_linked_service" "test" {
  resource_group_name = azurerm_resource_group.test.name
  workspace_id        = azurerm_log_analytics_workspace.test.id
  read_access_id      = azurerm_automation_account.test.id
}
`, r.template(data))
}

func (LogAnalyticsLinkedServiceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-la-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctestAutomation-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku_name = "Basic"

  tags = {
    Environment = "Test"
  }
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r LogAnalyticsLinkedServiceResource) withWriteAccessResourceId(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_cluster" "test" {
  name                = "acctest-LA-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_log_analytics_linked_service" "test" {
  resource_group_name = azurerm_resource_group.test.name
  workspace_id        = azurerm_log_analytics_workspace.test.id
  write_access_id     = azurerm_log_analytics_cluster.test.id
}
`, r.template(data), data.RandomInteger)
}

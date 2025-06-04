// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loganalytics_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/storageinsights"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LogAnalyticsStorageInsightsResource struct{}

func TestAccLogAnalyticsStorageInsights_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_storage_insights", "test")
	r := LogAnalyticsStorageInsightsResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("storage_account_key"), // key is not returned by the API
	})
}

func TestAccLogAnalyticsStorageInsights_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_storage_insights", "test")
	r := LogAnalyticsStorageInsightsResource{}
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

func TestAccLogAnalyticsStorageInsights_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_storage_insights", "test")
	r := LogAnalyticsStorageInsightsResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("storage_account_key"), // key is not returned by the API
	})
}

func TestAccLogAnalyticsStorageInsights_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_storage_insights", "test")
	r := LogAnalyticsStorageInsightsResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("storage_account_key"),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("storage_account_key"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("storage_account_key"), // key is not returned by the API
	})
}

func TestAccLogAnalyticsStorageInsights_updateStorageAccount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_storage_insights", "test")
	r := LogAnalyticsStorageInsightsResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("storage_account_key"),
		{
			Config: r.updateStorageAccount(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("storage_account_key"), // key is not returned by the API
	})
}

func (t LogAnalyticsStorageInsightsResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := storageinsights.ParseStorageInsightConfigID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.LogAnalytics.StorageInsightsClient.StorageInsightConfigsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving Log Analytics Storage Insights (%s): %+v", id.String(), err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (LogAnalyticsStorageInsightsResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-la-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_storage_account" "test" {
  name                = "acctestsads%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString)
}

func (r LogAnalyticsStorageInsightsResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_storage_insights" "test" {
  name                = "acctest-la-%d"
  resource_group_name = azurerm_resource_group.test.name
  workspace_id        = azurerm_log_analytics_workspace.test.id

  storage_account_id  = azurerm_storage_account.test.id
  storage_account_key = azurerm_storage_account.test.primary_access_key
}
`, r.template(data), data.RandomInteger)
}

func (r LogAnalyticsStorageInsightsResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_storage_insights" "import" {
  name                = azurerm_log_analytics_storage_insights.test.name
  resource_group_name = azurerm_log_analytics_storage_insights.test.resource_group_name
  workspace_id        = azurerm_log_analytics_storage_insights.test.workspace_id

  storage_account_id  = azurerm_storage_account.test.id
  storage_account_key = azurerm_storage_account.test.primary_access_key
}
`, r.basic(data))
}

func (r LogAnalyticsStorageInsightsResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_storage_insights" "test" {
  name                = "acctest-la-%d"
  resource_group_name = azurerm_resource_group.test.name
  workspace_id        = azurerm_log_analytics_workspace.test.id

  blob_container_names = ["wad-iis-logfiles"]
  table_names          = ["WADWindowsEventLogsTable", "LinuxSyslogVer2v0"]

  storage_account_id  = azurerm_storage_account.test.id
  storage_account_key = azurerm_storage_account.test.primary_access_key
}
`, r.template(data), data.RandomInteger)
}

func (r LogAnalyticsStorageInsightsResource) updateStorageAccount(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test2" {
  name                = "acctestsads%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_log_analytics_storage_insights" "test" {
  name                = "acctest-la-%d"
  resource_group_name = azurerm_resource_group.test.name
  workspace_id        = azurerm_log_analytics_workspace.test.id

  blob_container_names = ["wad-iis-logfiles"]
  table_names          = ["WADWindowsEventLogsTable", "LinuxSyslogVer2v0"]

  storage_account_id  = azurerm_storage_account.test2.id
  storage_account_key = azurerm_storage_account.test2.primary_access_key
}
`, r.template(data), data.RandomStringOfLength(6), data.RandomInteger)
}

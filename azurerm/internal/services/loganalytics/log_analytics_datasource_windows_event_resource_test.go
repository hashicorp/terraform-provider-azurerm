package loganalytics_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type LogAnalyticsDataSourceWindowsEventResource struct {
}

func TestAccLogAnalyticsDataSourceWindowsEvent_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_datasource_windows_event", "test")
	r := LogAnalyticsDataSourceWindowsEventResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsDataSourceWindowsEvent_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_datasource_windows_event", "test")
	r := LogAnalyticsDataSourceWindowsEventResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsDataSourceWindowsEvent_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_datasource_windows_event", "test")
	r := LogAnalyticsDataSourceWindowsEventResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsDataSourceWindowsEvent_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_datasource_windows_event", "test")
	r := LogAnalyticsDataSourceWindowsEventResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (LogAnalyticsDataSourceWindowsEventResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.LogAnalyticsDataSourceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.LogAnalytics.DataSourcesClient.Get(ctx, id.ResourceGroup, id.Workspace, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Log Analytics Data Source Windows Event %s (resource group: %s): %v", id.Name, id.ResourceGroup, err)
	}

	return utils.Bool(resp.Properties != nil), nil
}

func (r LogAnalyticsDataSourceWindowsEventResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_datasource_windows_event" "test" {
  name                = "acctestLADS-WE-%d"
  resource_group_name = azurerm_resource_group.test.name
  workspace_name      = azurerm_log_analytics_workspace.test.name
  event_log_name      = "Application"
  event_types         = ["error"]
}
`, r.template(data), data.RandomInteger)
}

func (r LogAnalyticsDataSourceWindowsEventResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_datasource_windows_event" "test" {
  name                = "acctestLADS-WE-%d"
  resource_group_name = azurerm_resource_group.test.name
  workspace_name      = azurerm_log_analytics_workspace.test.name
  event_log_name      = "Application"
  event_types         = ["InforMation", "warning", "Error"]
}
`, r.template(data), data.RandomInteger)
}

func (r LogAnalyticsDataSourceWindowsEventResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_datasource_windows_event" "import" {
  name                = azurerm_log_analytics_datasource_windows_event.test.name
  resource_group_name = azurerm_log_analytics_datasource_windows_event.test.resource_group_name
  workspace_name      = azurerm_log_analytics_datasource_windows_event.test.workspace_name
  event_log_name      = azurerm_log_analytics_datasource_windows_event.test.event_log_name
  event_types         = azurerm_log_analytics_datasource_windows_event.test.event_types
}
`, r.basic(data))
}

func (LogAnalyticsDataSourceWindowsEventResource) template(data acceptance.TestData) string {
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
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

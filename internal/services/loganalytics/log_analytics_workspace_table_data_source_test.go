package loganalytics_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type LogAnalyticsWorkspaceTableDataSource struct{}

func TestLogAnalyticsWorkspaceTableDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_log_analytics_workspace_table", "this")
	r := LogAnalyticsWorkspaceTableDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basicWithDataSource(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("retention_in_days").HasValue("30"),
				check.That(data.ResourceName).Key("total_retention_in_days").HasValue("30"),
				check.That(data.ResourceName).Key("plan").HasValue("Analytics"),
			),
		},
	})
}

func (d LogAnalyticsWorkspaceTableDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
		features {}
}

resource "azurerm_resource_group" "this" {
  name     = "test-resource-group"
	location = "%s"
}

resource "azurerm_log_analytics_workspace" "this" {
	name                = "test"
  resource_group_name = "test-resource-group"
  location            = "%s"
  sku                 = "PerGB2018"
  retention_in_days   = 30

  tags = {
  	env = "test"
 	}
}
`, data.Locations.Primary, data.Locations.Primary)
}

func (d LogAnalyticsWorkspaceTableDataSource) basicWithDataSource(data acceptance.TestData) string {
	config := d.basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_log_analytics_workspace_table" "this" {
	name                = "InsightsMetrics"
	workspace_id        = azurerm_log_analytics_workspace.this.id
	resource_group_name = "test-resource-group"
}
`, config)
}

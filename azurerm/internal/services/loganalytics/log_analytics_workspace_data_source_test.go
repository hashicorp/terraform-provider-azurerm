package loganalytics_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type LogAnalyticsWorkspaceDataSource struct {
}

func TestAccDataSourceLogAnalyticsWorkspace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_log_analytics_workspace", "test")
	r := LogAnalyticsWorkspaceDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basicWithDataSource(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("sku").HasValue("pergb2018"),
				check.That(data.ResourceName).Key("retention_in_days").HasValue("30"),
				check.That(data.ResourceName).Key("daily_quota_gb").HasValue("-1"),
			),
		},
	})
}

func TestAccDataSourceLogAnalyticsWorkspace_volumeCapWithDataSource(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_log_analytics_workspace", "test")
	r := LogAnalyticsWorkspaceDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.volumeCapWithDataSource(data, 4.5),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("sku").HasValue("pergb2018"),
				check.That(data.ResourceName).Key("retention_in_days").HasValue("30"),
				check.That(data.ResourceName).Key("daily_quota_gb").HasValue("4.5"),
			),
		},
	})
}

func (LogAnalyticsWorkspaceDataSource) basicWithDataSource(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_log_analytics_workspace" "test" {
  name                = azurerm_log_analytics_workspace.test.name
  resource_group_name = azurerm_log_analytics_workspace.test.resource_group_name
}
`, LogAnalyticsWorkspaceResource{}.complete(data))
}

func (LogAnalyticsWorkspaceDataSource) volumeCapWithDataSource(data acceptance.TestData, volumeCapGb float64) string {
	return fmt.Sprintf(`
%s

data "azurerm_log_analytics_workspace" "test" {
  name                = azurerm_log_analytics_workspace.test.name
  resource_group_name = azurerm_log_analytics_workspace.test.resource_group_name
}
`, LogAnalyticsWorkspaceResource{}.withVolumeCap(data, volumeCapGb))
}

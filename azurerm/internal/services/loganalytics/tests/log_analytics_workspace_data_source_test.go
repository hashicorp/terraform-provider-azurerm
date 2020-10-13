package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMLogAnalyticsWorkspace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_log_analytics_workspace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMLogAnalyticsWorkspace_basicWithDataSource(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "sku", "pergb2018"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_in_days", "30"),
					resource.TestCheckResourceAttr(data.ResourceName, "daily_quota_gb", "-1"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMLogAnalyticsWorkspace_volumeCapWithDataSource(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_log_analytics_workspace", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLogAnalyticsWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMLogAnalyticsWorkspace_volumeCapWithDataSource(data, 4.5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "sku", "pergb2018"),
					resource.TestCheckResourceAttr(data.ResourceName, "retention_in_days", "30"),
					resource.TestCheckResourceAttr(data.ResourceName, "daily_quota_gb", "4.5"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMLogAnalyticsWorkspace_basicWithDataSource(data acceptance.TestData) string {
	config := testAccAzureRMLogAnalyticsWorkspace_complete(data)
	return fmt.Sprintf(`
%s

data "azurerm_log_analytics_workspace" "test" {
  name                = azurerm_log_analytics_workspace.test.name
  resource_group_name = azurerm_log_analytics_workspace.test.resource_group_name
}
`, config)
}

func testAccDataSourceAzureRMLogAnalyticsWorkspace_volumeCapWithDataSource(data acceptance.TestData, volumeCapGb float64) string {
	config := testAccAzureRMLogAnalyticsWorkspace_withVolumeCap(data, volumeCapGb)
	return fmt.Sprintf(`
%s

data "azurerm_log_analytics_workspace" "test" {
  name                = azurerm_log_analytics_workspace.test.name
  resource_group_name = azurerm_log_analytics_workspace.test.resource_group_name
}
`, config)
}

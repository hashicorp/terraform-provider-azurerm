package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMLogAnalyticsWorkspace_basic(t *testing.T) {
	dataSourceName := "data.azurerm_log_analytics_workspace.test"
	ri := tf.AccRandTimeInt()
	config := testAccDataSourceAzureRMLogAnalyticsWorkspace_basicWithDataSource(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "sku", "pergb2018"),
					resource.TestCheckResourceAttr(dataSourceName, "retention_in_days", "30"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMLogAnalyticsWorkspace_basicWithDataSource(rInt int, location string) string {
	config := testAccAzureRMLogAnalyticsWorkspace_complete(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_log_analytics_workspace" "test" {
  name                = "${azurerm_log_analytics_workspace.test.name}"
  resource_group_name = "${azurerm_log_analytics_workspace.test.resource_group_name}"
}
`, config)
}

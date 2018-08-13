package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMLogAnalyticsWorkspac_basic(t *testing.T) {
	dataSourceName := "data.azurerm_log_analytics_workspace.test"
	ri := acctest.RandInt()
	rs := acctest.RandString(4)
	location := testLocation()
	preConfig := testAccDataSourceAzureRMLogAnalyticsWorkspace_basic(ri, rs, location)
	config := testAccDataSourceAzureRMLogAnalyticsWorkspace_basicWithDataSource(ri, rs, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "sku", "Free"),
					resource.TestCheckResourceAttr(dataSourceName, "retention_in_days", "30"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.environment", "production"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMLogAnalyticsWorkspace_basicWithDataSource(rInt int, rString string, location string) string {
	config := testAccAzureRMLogAnalyticsWorkspace_requiredOnly(rInt, rString, location)
	return fmt.Sprintf(`
%s

data "azurerm_log_analytics_workspace" "test" {
  name                = "${azurerm_log_analytics_workspace.test.name}"
  resource_group_name = "${azurerm_log_analytics_workspace.test.resource_group_name}"
}
`, config)
}

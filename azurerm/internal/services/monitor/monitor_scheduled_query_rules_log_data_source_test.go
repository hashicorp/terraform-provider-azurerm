package monitor_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMMonitorScheduledQueryRules_LogToMetricAction(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_monitor_scheduled_query_rules_log", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMMonitorScheduledQueryRules_LogToMetricActionConfig(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "id"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMMonitorScheduledQueryRules_LogToMetricActionConfig(data acceptance.TestData) string {
	template := testAccAzureRMMonitorScheduledQueryRules_LogToMetricActionConfigBasic(data)
	return fmt.Sprintf(`
%s

data "azurerm_monitor_scheduled_query_rules_log" "test" {
  name                = basename(azurerm_monitor_scheduled_query_rules_log.test.id)
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, template)
}

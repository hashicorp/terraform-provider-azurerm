package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMMonitorScheduledQueryRulesAlertingAction(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_monitor_scheduled_query_rules_alert", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMMonitorScheduledQueryRules_alertingActionConfig(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "id"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMMonitorScheduledQueryRulesAlertingActionCrossResource(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_monitor_scheduled_query_rules_alert", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMMonitorScheduledQueryRules_alertingActionCrossResourceConfig(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "id"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMMonitorScheduledQueryRules_alertingActionConfig(data acceptance.TestData) string {
	ts := time.Now().Format(time.RFC3339)
	template := testAccAzureRMMonitorScheduledQueryRulesAlertingActionConfig_basic(data, ts)

	return fmt.Sprintf(`
%s

data "azurerm_monitor_scheduled_query_rules_alert" "test" {
  name                = basename(azurerm_monitor_scheduled_query_rules_alert.test.id)
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, template)
}

func testAccDataSourceAzureRMMonitorScheduledQueryRules_alertingActionCrossResourceConfig(data acceptance.TestData) string {
	template := testAccAzureRMMonitorScheduledQueryRules_alertingActionCrossResourceConfig(data)
	return fmt.Sprintf(`
%s

data "azurerm_monitor_scheduled_query_rules_alert" "test" {
  name                = basename(azurerm_monitor_scheduled_query_rules_alert.test.id)
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, template)
}

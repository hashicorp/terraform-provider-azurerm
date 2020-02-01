package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMMonitorScheduledQueryRules_LogToMetricAction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_scheduled_query_rules_log", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorScheduledQueryRulesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorScheduledQueryRules_logToMetricAction(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorScheduledQueryRulesExists(data.ResourceName),
				),
			},
			{
				ResourceName:      data.ResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAzureRMMonitorScheduledQueryRules_logToMetricAction(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-monitor-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestAppInsights-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  application_type    = "web"
}

resource "azurerm_monitor_action_group" "test" {
	name                = "acctestActionGroup-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  short_name          = "acctestag"
}

resource "azurerm_monitor_scheduled_query_rules_log" "test" {
  name                = "acctestsqr-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
	description         = "test log to metric action"
	enabled             = true

	data_source_id = "${azurerm_application_insights.test.id}"

	criteria {
		metric_name        = "Average_%% Idle Time"
		dimension {
			name             = "InstanceName"
			operator         = "Include"
			values           = ["50"]
		}
	}
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testCheckAzureRMMonitorScheduledQueryRulesLogDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Monitor.ScheduledQueryRulesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_monitor_scheduled_query_rules_log" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Activity log alert still exists:\n%#v", resp)
		}
	}

	return nil
}

func testCheckAzureRMMonitorScheduledQueryRulesExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Activity Log Alert Instance: %s", name)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Monitor.ScheduledQueryRulesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on monitorScheduledQueryRulesClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Activity Log Alert Instance %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

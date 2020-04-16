package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceSentinelAlertRuleMsSecurityIncident_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_sentinel_alert_rule_ms_security_incident", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSentinelAlertRuleMsSecurityIncident_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "id"),
				),
			},
		},
	})
}

func testAccDataSourceSentinelAlertRuleMsSecurityIncident_basic(data acceptance.TestData) string {
	config := testAccAzureRMSentinelAlertRuleMsSecurityIncident_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_sentinel_alert_rule_ms_security_incident" "test" {
  name                       = azurerm_sentinel_alert_rule_ms_security_incident.test.name
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id
}
`, config)
}

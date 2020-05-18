package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceSentinelAlertRuleTemplate_fusionAlertRuleTemplate(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_sentinel_alert_rule_template", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSentinelAlertRuleTemplate_byDisplayName(data, "Advanced Multistage Attack Detection"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "display_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckNoResourceAttr(data.ResourceName, "scheduled_template"),
					resource.TestCheckNoResourceAttr(data.ResourceName, "ms_security_incident_template"),
				),
			},
		},
	})
}

func TestAccDataSourceSentinelAlertRuleTemplate_msSecurityIncidentAlertRuleTemplate(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_sentinel_alert_rule_template", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSentinelAlertRuleTemplate_byDisplayName(data, "Create incidents based on Azure Security Center for IoT alerts"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "display_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ms_security_incident_template.0.description"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ms_security_incident_template.0.product_filter"),
					resource.TestCheckNoResourceAttr(data.ResourceName, "scheduled_template"),
				),
			},
		},
	})
}

func TestAccDataSourceSentinelAlertRuleTemplate_scheduledAlertRuleTemplate(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_sentinel_alert_rule_template", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSentinelAlertRuleTemplate_byDisplayName(data, "Malware in the recycle bin"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "display_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "scheduled_template.0.description"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "scheduled_template.0.tactics.0"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "scheduled_template.0.severity"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "scheduled_template.0.query"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "scheduled_template.0.query_frequency"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "scheduled_template.0.query_period"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "scheduled_template.0.trigger_operator"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "scheduled_template.0.trigger_threshold"),
					resource.TestCheckNoResourceAttr(data.ResourceName, "ms_security_incident_template"),
				),
			},
		},
	})
}

func testAccDataSourceSentinelAlertRuleTemplate_byDisplayName(data acceptance.TestData, displayName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sentinel-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}


data "azurerm_sentinel_alert_rule_template" "test" {
  display_name               = "%s"
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, displayName)
}

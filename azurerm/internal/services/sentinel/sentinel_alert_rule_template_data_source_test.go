package sentinel_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

type SentinelAlertRuleTemplateDataSource struct{}

func TestAccSentinelAlertRuleTemplateDataSource_fusion(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_sentinel_alert_rule_template", "test")
	r := SentinelAlertRuleTemplateDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.byDisplayName(data, "Advanced Multistage Attack Detection"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("display_name").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("scheduled_template").DoesNotExist(),
				check.That(data.ResourceName).Key("ms_security_incident_template").DoesNotExist(),
			),
		},
	})
}

func TestAccSentinelAlertRuleTemplateDataSource_msSecurityIncident(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_sentinel_alert_rule_template", "test")
	r := SentinelAlertRuleTemplateDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.byDisplayName(data, "Create incidents based on Azure Security Center for IoT alerts"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("display_name").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("ms_security_incident_template.0.description").Exists(),
				check.That(data.ResourceName).Key("ms_security_incident_template.0.product_filter").Exists(),
				check.That(data.ResourceName).Key("scheduled_template").DoesNotExist(),
			),
		},
	})
}

func TestAccSentinelAlertRuleTemplateDataSource_scheduled(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_sentinel_alert_rule_template", "test")
	r := SentinelAlertRuleTemplateDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.byDisplayName(data, "Malware in the recycle bin"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("display_name").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("scheduled_template.0.description").Exists(),
				check.That(data.ResourceName).Key("scheduled_template.0.tactics.0").Exists(),
				check.That(data.ResourceName).Key("scheduled_template.0.severity").Exists(),
				check.That(data.ResourceName).Key("scheduled_template.0.query").Exists(),
				check.That(data.ResourceName).Key("scheduled_template.0.query_frequency").Exists(),
				check.That(data.ResourceName).Key("scheduled_template.0.query_period").Exists(),
				check.That(data.ResourceName).Key("scheduled_template.0.trigger_operator").Exists(),
				check.That(data.ResourceName).Key("scheduled_template.0.trigger_threshold").Exists(),
				check.That(data.ResourceName).Key("ms_security_incident_template").DoesNotExist(),
			),
		},
	})
}

func (SentinelAlertRuleTemplateDataSource) byDisplayName(data acceptance.TestData, displayName string) string {
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

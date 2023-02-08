package sentinel_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type SentinelAlertRuleAnomalyDataSource struct{}

func TestAccSentinelAlertRuleAnomalyDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_sentinel_alert_rule_anomaly", "test")
	r := SentinelAlertRuleAnomalyDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("display_name").Exists(),
				check.That(data.ResourceName).Key("anomaly_version").Exists(),
				check.That(data.ResourceName).Key("anomaly_settings_version").Exists(),
				check.That(data.ResourceName).Key("customizable_observations").Exists(),
				check.That(data.ResourceName).Key("description").Exists(),
				check.That(data.ResourceName).Key("enabled").Exists(),
				check.That(data.ResourceName).Key("frequency").Exists(),
				check.That(data.ResourceName).Key("is_default_settings").Exists(),
				check.That(data.ResourceName).Key("required_data_connector.#").HasValue("2"),
				check.That(data.ResourceName).Key("mode").Exists(),
				check.That(data.ResourceName).Key("settings_definition_id").Exists(),
				check.That(data.ResourceName).Key("tactics.#").HasValue("1"),
				check.That(data.ResourceName).Key("techniques.#").HasValue("1"),
			),
		},
	})
}

func (SentinelAlertRuleAnomalyDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_sentinel_alert_rule_anomaly" "test" {
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id
  display_name               = "UEBA Anomalous Sign In"
  depends_on                 = [azurerm_sentinel_log_analytics_workspace_onboarding.test]
}
`, SecurityInsightsSentinelOnboardingStateResource{}.basic(data))
}

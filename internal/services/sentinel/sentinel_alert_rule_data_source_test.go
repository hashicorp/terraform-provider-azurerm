// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sentinel_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type SentinelAlertRuleDataSource struct{}

func TestAccSentinelAlertRuleDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_sentinel_alert_rule", "test")
	r := SentinelAlertRuleDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
			),
		},
	})
}

func (SentinelAlertRuleDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_sentinel_alert_rule" "test" {
  name                       = azurerm_sentinel_alert_rule_ms_security_incident.test.name
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id
}
`, SentinelAlertRuleMsSecurityIncidentResource{}.basic(data))
}

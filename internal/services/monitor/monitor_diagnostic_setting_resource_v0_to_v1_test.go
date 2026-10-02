// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

// TestAccMonitorDiagnosticSetting_V0ToV1_560 tests the state migration from the legacy pipe-separated ID format
// to the canonical ARM Scoped Diagnostic Setting ID.
// It uses v5.6.0 as the setup version because it is the last release where the pipe ID format was used in state.
func TestAccMonitorDiagnosticSetting_V0ToV1_560(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_diagnostic_setting", "test")
	r := MonitorDiagnosticSettingResource{}

	data.ResourceRegressionTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").HasValue(fmt.Sprintf("/providers/Microsoft.Management/managementGroups/acctestMG%[1]d|acctestMDS%[1]d", data.RandomInteger)),
			),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("id").HasValue(fmt.Sprintf("/providers/Microsoft.Management/managementGroups/acctestMG%[1]d/providers/Microsoft.Insights/diagnosticSettings/acctestMDS%[1]d", data.RandomInteger)),
			),
		},
	}, "5.6.0")
}

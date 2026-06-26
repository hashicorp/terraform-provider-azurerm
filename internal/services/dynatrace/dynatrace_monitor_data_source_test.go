// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dynatrace_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type MonitorsDataSource struct{}

func TestAccDynatraceMonitorsDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_dynatrace_monitor", "test")
	d := MonitorsDataSource{}
	r := NewDynatraceMonitorResource()
	r.preCheck(t)

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic(data, r),
			Check:  acceptance.ComposeTestCheckFunc(),
		},
	})
}

func (d MonitorsDataSource) basic(data acceptance.TestData, resource MonitorsResource) string {
	return fmt.Sprintf(`
%s

data "azurerm_dynatrace_monitor" "test" {
  name                = azurerm_dynatrace_monitor.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, resource.basic(data))
}

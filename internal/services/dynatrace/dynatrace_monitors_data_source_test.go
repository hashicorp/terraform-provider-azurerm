// Copyright (c) HashiCorp, Inc.
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
	r := MonitorsDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check:  acceptance.ComposeTestCheckFunc(),
		},
	})
}

func (d MonitorsDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_dynatrace_monitor" "test" {
  name                = azurerm_dynatrace_monitor.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, MonitorsResource{}.basic(data))
}

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type MonitorWorkspaceDataSource struct{}

func TestAccMonitorWorkspaceDataSourceDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_monitor_workspace", "test")
	d := MonitorWorkspaceDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
			),
		},
	})
}

func (d MonitorWorkspaceDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_monitor_workspace" "test" {
  name                = azurerm_monitor_workspace.test.name
  resource_group_name = azurerm_monitor_workspace.test.resource_group_name
}
`, WorkspaceTestResource{}.complete(data))
}

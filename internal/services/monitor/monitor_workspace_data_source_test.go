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

func TestAccMonitorWorkspaceDataSourceDataSource_privateEndpointConnection(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_monitor_workspace", "test")
	d := MonitorWorkspaceDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.privateEndpointConnection(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("private_endpoint_connections.#").HasValue("1"),
				check.That(data.ResourceName).Key("private_endpoint_connections.0.name").IsNotEmpty(),
				check.That(data.ResourceName).Key("private_endpoint_connections.0.id").IsNotEmpty(),
				check.That(data.ResourceName).Key("private_endpoint_connections.0.group_ids.#").HasValue("1"),
			),
		},
	})
}

func (d MonitorWorkspaceDataSource) privateEndpointConnection(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_monitor_workspace" "test" {
  name                = azurerm_monitor_workspace.test.name
  resource_group_name = azurerm_monitor_workspace.test.resource_group_name
}
`, WorkspaceTestResource{}.privateEndpointConnection(data))
}

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type VirtualHubConnectionDataSource struct{}

func TestAccDataSourceAzureRMVirtualHubConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_virtual_hub_connection", "test")
	r := VirtualHubConnectionDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("virtual_hub_id").Exists(),
				check.That(data.ResourceName).Key("remote_virtual_network_id").Exists(),
			),
		},
	})
}

func (VirtualHubConnectionDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_virtual_hub_connection" "test" {
  name                = azurerm_virtual_hub_connection.test.name
  resource_group_name = azurerm_virtual_hub.test.resource_group_name
  virtual_hub_name    = azurerm_virtual_hub.test.name
}
`, VirtualHubRouteTableResource{}.basic(data))
}

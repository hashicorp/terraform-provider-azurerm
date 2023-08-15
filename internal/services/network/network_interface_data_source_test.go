// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type NetworkInterfaceDataSource struct{}

func TestAccDataSourceArmNetworkInterface_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_network_interface", "test")
	r := NetworkInterfaceDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: NetworkInterfaceResource{}.static(data),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("private_ip_address").HasValue("10.0.2.15"),
			),
		},
	})
}

func (NetworkInterfaceDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_network_interface" "test" {
  name                = azurerm_network_interface.test.name
  resource_group_name = azurerm_network_interface.test.resource_group_name
}
`, NetworkInterfaceResource{}.static(data))
}

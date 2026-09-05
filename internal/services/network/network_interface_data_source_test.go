// Copyright IBM Corp. 2014, 2025
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
				check.That(data.ResourceName).Key("internal_domain_name_suffix").IsSet(),
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

func TestAccDataSourceArmNetworkInterface_auxiliary(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_network_interface", "test")
	r := NetworkInterfaceDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.auxiliary(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("auxiliary_mode").HasValue("AcceleratedConnections"),
				check.That(data.ResourceName).Key("auxiliary_sku").HasValue("A2"),
			),
		},
	})
}

func (NetworkInterfaceDataSource) auxiliary(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_network_interface" "test" {
  name                = azurerm_network_interface.test.name
  resource_group_name = azurerm_network_interface.test.resource_group_name
}
`, NetworkInterfaceResource{}.auxiliaryAcceleratedConnections(data))
}

func TestAccDataSourceArmNetworkInterface_edgeZone(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_network_interface", "test")
	r := NetworkInterfaceDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.edgeZone(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("edge_zone").IsSet(),
			),
		},
	})
}

func (NetworkInterfaceDataSource) edgeZone(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_network_interface" "test" {
  name                = azurerm_network_interface.test.name
  resource_group_name = azurerm_network_interface.test.resource_group_name
}
`, NetworkInterfaceResource{}.edgeZone(data))
}

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mobilenetwork_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type MobileNetworkPacketCoreDataPlaneDataSource struct{}

func TestAccMobileNetworkPacketCoreDataPlaneDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_packet_core_data_plane", "test")
	d := MobileNetworkPacketCoreDataPlaneDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key(`location`).Exists(),
				check.That(data.ResourceName).Key(`user_plane_access_name`).HasValue("default-interface"),
				check.That(data.ResourceName).Key(`user_plane_access_ipv4_address`).HasValue("192.168.1.199"),
				check.That(data.ResourceName).Key(`user_plane_access_ipv4_gateway`).HasValue("192.168.1.1"),
				check.That(data.ResourceName).Key(`user_plane_access_ipv4_subnet`).HasValue("192.168.1.0/25"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
			),
		},
	})
}

func (r MobileNetworkPacketCoreDataPlaneDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

data "azurerm_mobile_network_packet_core_data_plane" "test" {
  name                                        = azurerm_mobile_network_packet_core_data_plane.test.name
  mobile_network_packet_core_control_plane_id = azurerm_mobile_network_packet_core_data_plane.test.mobile_network_packet_core_control_plane_id
}
`, MobileNetworkPacketCoreDataPlaneResource{}.complete(data))
}

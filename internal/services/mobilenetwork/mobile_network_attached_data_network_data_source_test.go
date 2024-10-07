// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mobilenetwork_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type MobileNetworkAttachedDataNetworkDataSource struct{}

func TestAccMobileNetworkAttachedDataNetworkDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_attached_data_network", "test")
	d := MobileNetworkAttachedDataNetworkDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key(`location`).Exists(),
				check.That(data.ResourceName).Key(`dns_addresses.0`).HasValue("1.1.1.1"),
				check.That(data.ResourceName).Key(`user_equipment_address_pool_prefixes.0`).HasValue("2.4.1.0/24"),
				check.That(data.ResourceName).Key(`user_equipment_static_address_pool_prefixes.0`).HasValue("2.4.2.0/24"),
				check.That(data.ResourceName).Key(`network_address_port_translation.0.pinhole_maximum_number`).HasValue("65536"),
				check.That(data.ResourceName).Key(`network_address_port_translation.0.icmp_pinhole_timeout_in_seconds`).HasValue("30"),
				check.That(data.ResourceName).Key(`network_address_port_translation.0.tcp_pinhole_timeout_in_seconds`).HasValue("100"),
				check.That(data.ResourceName).Key(`network_address_port_translation.0.udp_pinhole_timeout_in_seconds`).HasValue("39"),
				check.That(data.ResourceName).Key(`network_address_port_translation.0.port_range.0.maximum`).HasValue("49999"),
				check.That(data.ResourceName).Key(`network_address_port_translation.0.port_range.0.minimum`).HasValue("1024"),
				check.That(data.ResourceName).Key(`network_address_port_translation.0.tcp_port_reuse_minimum_hold_time_in_seconds`).HasValue("120"),
				check.That(data.ResourceName).Key(`network_address_port_translation.0.udp_port_reuse_minimum_hold_time_in_seconds`).HasValue("60"),
				check.That(data.ResourceName).Key(`user_plane_access_name`).HasValue("test"),
				check.That(data.ResourceName).Key(`user_plane_access_ipv4_address`).HasValue("10.204.141.4"),
				check.That(data.ResourceName).Key(`user_plane_access_ipv4_gateway`).HasValue("10.204.141.1"),
				check.That(data.ResourceName).Key(`user_plane_access_ipv4_subnet`).HasValue("10.204.141.0/24"),
				check.That(data.ResourceName).Key(`tags.%`).HasValue("1"),
			),
		},
	})
}

func (r MobileNetworkAttachedDataNetworkDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

data "azurerm_mobile_network_attached_data_network" "test" {
  mobile_network_data_network_name         = azurerm_mobile_network_attached_data_network.test.mobile_network_data_network_name
  mobile_network_packet_core_data_plane_id = azurerm_mobile_network_attached_data_network.test.mobile_network_packet_core_data_plane_id
}
`, MobileNetworkAttachedDataNetworkResource{}.complete(data))
}

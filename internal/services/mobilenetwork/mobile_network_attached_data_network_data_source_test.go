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
				check.That(data.ResourceName).Key(`location`).HasValue(data.Locations.Primary),
				check.That(data.ResourceName).Key(`dns_addresses.0`).HasValue("1.1.1.1"),
				check.That(data.ResourceName).Key(`user_equipment_address_pool_prefixes.0`).HasValue("2.4.0.0/16"),
				check.That(data.ResourceName).Key(`user_equipment_static_address_pool_prefixes.0`).HasValue("2.4.0.0/16"),
				check.That(data.ResourceName).Key(`network_address_port_translation_configuration.0.enabled`).HasValue("true"),
				check.That(data.ResourceName).Key(`network_address_port_translation_configuration.0.pinhole_maximum_number`).HasValue("65536"),
				check.That(data.ResourceName).Key(`network_address_port_translation_configuration.0.pinhole_timeouts_in_seconds.0.icmp`).HasValue("30"),
				check.That(data.ResourceName).Key(`network_address_port_translation_configuration.0.pinhole_timeouts_in_seconds.0.tcp`).HasValue("100"),
				check.That(data.ResourceName).Key(`network_address_port_translation_configuration.0.pinhole_timeouts_in_seconds.0.udp`).HasValue("39"),
				check.That(data.ResourceName).Key(`network_address_port_translation_configuration.0.port_range.0.max_port`).HasValue("49999"),
				check.That(data.ResourceName).Key(`network_address_port_translation_configuration.0.port_range.0.min_port`).HasValue("1024"),
				check.That(data.ResourceName).Key(`network_address_port_translation_configuration.0.port_reuse_minimum_hold_time_in_seconds.0.tcp`).HasValue("120"),
				check.That(data.ResourceName).Key(`network_address_port_translation_configuration.0.port_reuse_minimum_hold_time_in_seconds.0.udp`).HasValue("60"),
				check.That(data.ResourceName).Key(`user_plane_data_interface.0.name`).HasValue("test"),
				check.That(data.ResourceName).Key(`user_plane_data_interface.0.ipv4_address`).HasValue("10.204.141.4"),
				check.That(data.ResourceName).Key(`user_plane_data_interface.0.ipv4_gateway`).HasValue("10.204.141.1"),
				check.That(data.ResourceName).Key(`user_plane_data_interface.0.ipv4_subnet`).HasValue("10.204.141.0/24"),
				check.That(data.ResourceName).Key(`tags.%`).HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func (r MobileNetworkAttachedDataNetworkDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

data "azurerm_mobile_network_attached_data_network" "test" {
  name                                     = azurerm_mobile_network_attached_data_network.test.name
  mobile_network_packet_core_data_plane_id = azurerm_mobile_network_attached_data_network.test.mobile_network_packet_core_data_plane_id
}
`, MobileNetworkAttachedDataNetworkResource{}.complete(data))
}

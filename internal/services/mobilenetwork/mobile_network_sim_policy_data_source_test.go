package mobilenetwork_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type MobileNetworkSimPolicyDataSource struct{}

func TestAccMobileNetworkSimPolicyDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_sim_policy", "test")
	d := MobileNetworkSimPolicyDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key(`location`).Exists(),
				check.That(data.ResourceName).Key(`default_slice_id`).Exists(),
				check.That(data.ResourceName).Key(`registration_timer_in_seconds`).HasValue("3240"),
				check.That(data.ResourceName).Key(`rfsp_index`).HasValue("1"),
				check.That(data.ResourceName).Key(`slice_configurations.0.default_data_network_id`).Exists(),
				check.That(data.ResourceName).Key(`slice_configurations.0.slice_id`).Exists(),
				check.That(data.ResourceName).Key(`slice_configurations.0.data_network_configurations.0.allocation_and_retention_priority_level`).HasValue("9"),
				check.That(data.ResourceName).Key(`slice_configurations.0.data_network_configurations.0.default_session_type`).HasValue("IPv4"),
				check.That(data.ResourceName).Key(`slice_configurations.0.data_network_configurations.0.qos_indicator`).HasValue("9"),
				check.That(data.ResourceName).Key(`slice_configurations.0.data_network_configurations.0.preemption_capability`).HasValue("NotPreempt"),
				check.That(data.ResourceName).Key(`slice_configurations.0.data_network_configurations.0.preemption_vulnerability`).HasValue("Preemptable"),
				check.That(data.ResourceName).Key(`slice_configurations.0.data_network_configurations.0.allowed_services_ids.#`).HasValue("1"),
				check.That(data.ResourceName).Key(`slice_configurations.0.data_network_configurations.0.data_network_id`).Exists(),
				check.That(data.ResourceName).Key(`slice_configurations.0.data_network_configurations.0.max_buffered_packets`).HasValue("200"),
				check.That(data.ResourceName).Key(`slice_configurations.0.data_network_configurations.0.session_aggregate_maximum_bit_rate.0.downlink`).HasValue("1 Gbps"),
				check.That(data.ResourceName).Key(`slice_configurations.0.data_network_configurations.0.session_aggregate_maximum_bit_rate.0.uplink`).HasValue("500 Mbps"),
				check.That(data.ResourceName).Key(`user_equipment_aggregate_maximum_bit_rate.0.downlink`).HasValue("1 Gbps"),
				check.That(data.ResourceName).Key(`user_equipment_aggregate_maximum_bit_rate.0.uplink`).HasValue("500 Mbps"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func (r MobileNetworkSimPolicyDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

data "azurerm_mobile_network_sim_policy" "test" {
  name              = azurerm_mobile_network_sim_policy.test.name
  mobile_network_id = azurerm_mobile_network_sim_policy.test.mobile_network_id
}
`, MobileNetworkSimPolicyResource{}.complete(data))
}

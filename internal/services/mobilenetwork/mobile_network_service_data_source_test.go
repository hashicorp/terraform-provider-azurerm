// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mobilenetwork_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type MobileNetworkServiceDataSource struct{}

func TestAccMobileNetworkServiceDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_service", "test")
	d := MobileNetworkServiceDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key(`location`).Exists(),
				check.That(data.ResourceName).Key(`service_precedence`).HasValue("0"),
				check.That(data.ResourceName).Key(`pcc_rule.0.name`).HasValue("default-rule"),
				check.That(data.ResourceName).Key(`pcc_rule.0.precedence`).HasValue("1"),
				check.That(data.ResourceName).Key(`pcc_rule.0.traffic_control_enabled`).HasValue("true"),
				check.That(data.ResourceName).Key(`pcc_rule.0.qos_policy.0.allocation_and_retention_priority_level`).HasValue("9"),
				check.That(data.ResourceName).Key(`pcc_rule.0.qos_policy.0.qos_indicator`).HasValue("9"),
				check.That(data.ResourceName).Key(`pcc_rule.0.qos_policy.0.preemption_capability`).HasValue("NotPreempt"),
				check.That(data.ResourceName).Key(`pcc_rule.0.qos_policy.0.preemption_vulnerability`).HasValue("Preemptable"),
				check.That(data.ResourceName).Key(`pcc_rule.0.qos_policy.0.guaranteed_bit_rate.0.downlink`).HasValue("100 Mbps"),
				check.That(data.ResourceName).Key(`pcc_rule.0.qos_policy.0.guaranteed_bit_rate.0.uplink`).HasValue("10 Mbps"),
				check.That(data.ResourceName).Key(`pcc_rule.0.qos_policy.0.maximum_bit_rate.0.downlink`).HasValue("1 Gbps"),
				check.That(data.ResourceName).Key(`pcc_rule.0.qos_policy.0.maximum_bit_rate.0.uplink`).HasValue("100 Mbps"),
				check.That(data.ResourceName).Key(`pcc_rule.0.service_data_flow_template.0.direction`).HasValue("Uplink"),
				check.That(data.ResourceName).Key(`pcc_rule.0.service_data_flow_template.0.name`).HasValue("IP-to-server"),
				check.That(data.ResourceName).Key(`pcc_rule.0.service_data_flow_template.0.protocol.0`).HasValue("ip"),
				check.That(data.ResourceName).Key(`pcc_rule.0.service_data_flow_template.0.remote_ip_list.0`).HasValue("10.3.4.0/24"),
				check.That(data.ResourceName).Key(`service_qos_policy.0.allocation_and_retention_priority_level`).HasValue("9"),
				check.That(data.ResourceName).Key(`service_qos_policy.0.qos_indicator`).HasValue("9"),
				check.That(data.ResourceName).Key(`service_qos_policy.0.preemption_capability`).HasValue("NotPreempt"),
				check.That(data.ResourceName).Key(`service_qos_policy.0.preemption_vulnerability`).HasValue("Preemptable"),
				check.That(data.ResourceName).Key(`service_qos_policy.0.maximum_bit_rate.0.downlink`).HasValue("1 Gbps"),
				check.That(data.ResourceName).Key(`service_qos_policy.0.maximum_bit_rate.0.uplink`).HasValue("100 Mbps"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
			),
		},
	})
}

func (r MobileNetworkServiceDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

data "azurerm_mobile_network_service" "test" {
  name              = azurerm_mobile_network_service.test.name
  mobile_network_id = azurerm_mobile_network_service.test.mobile_network_id
}
`, MobileNetworkServiceResource{}.complete(data))
}

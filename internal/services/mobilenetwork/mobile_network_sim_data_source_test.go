package mobilenetwork_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type MobileNetworkSimDataSource struct{}

func TestAccMobileNetworkSimDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_sim", "test")
	d := MobileNetworkSimDataSource{}
	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key(`integrated_circuit_card_identifier`).HasValue("8900000000000000000"),
				check.That(data.ResourceName).Key(`international_mobile_subscriber_identity`).HasValue("000000000000000"),
				check.That(data.ResourceName).Key(`sim_policy_id`).Exists(),
				check.That(data.ResourceName).Key(`static_ip_configuration.0.attached_data_network_id`).Exists(),
				check.That(data.ResourceName).Key(`static_ip_configuration.0.slice_id`).Exists(),
				check.That(data.ResourceName).Key(`static_ip_configuration.0.static_ipv4_address`).HasValue("2.4.0.1"),
			),
		},
	})
}

func (r MobileNetworkSimDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

data "azurerm_mobile_network_sim" "test" {
  name                        = azurerm_mobile_network_sim.test.name
  mobile_network_sim_group_id = azurerm_mobile_network_sim.test.mobile_network_sim_group_id
}
`, MobileNetworkSimResource{}.complete(data))
}

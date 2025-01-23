package azurestackhci_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type StackHCILogicalNetworkDataSource struct{}

func TestAccStackHCILogicalNetworkDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_stack_hci_logical_network", "test")
	d := StackHCILogicalNetworkDataSource{}

	data.DataSourceTestInSequence(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").IsNotEmpty(),
				check.That(data.ResourceName).Key("custom_location_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("virtual_switch_name").HasValue("ConvergedSwitch(managementcompute)"),
				check.That(data.ResourceName).Key("dns_servers.#").HasValue("2"),
				check.That(data.ResourceName).Key("subnet.0.ip_allocation_method").HasValue("Static"),
				check.That(data.ResourceName).Key("subnet.0.address_prefix").Exists(),
				check.That(data.ResourceName).Key("subnet.0.vlan_id").HasValue("123"),
				check.That(data.ResourceName).Key("subnet.0.ip_pool.0.start").Exists(),
				check.That(data.ResourceName).Key("subnet.0.ip_pool.1.end").Exists(),
				check.That(data.ResourceName).Key("subnet.0.route.0.address_prefix").Exists(),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
			),
		},
	})
}

func (d StackHCILogicalNetworkDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_stack_hci_logical_network" "test" {
  name                = azurerm_stack_hci_logical_network.test.name
  resource_group_name = azurerm_stack_hci_logical_network.test.resource_group_name
}
`, StackHCILogicalNetworkResource{}.complete(data))
}

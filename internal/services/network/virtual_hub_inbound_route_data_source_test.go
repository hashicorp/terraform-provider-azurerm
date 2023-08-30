package network_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type VirtualHubInboundRouteDataSource struct{}

func TestAccDataSourceVirtualHubInboundRoute_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_virtual_hub_inbound_route", "test")
	r := VirtualHubInboundRouteDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check:  acceptance.ComposeTestCheckFunc(),
		},
	})
}

func (VirtualHubInboundRouteDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_virtual_hub_inbound_route" "test" {
  virtual_hub_id     = azurerm_virtual_hub.test.id
  target_resource_id = azurerm_virtual_hub_connection.test.id
  connection_type    = "HubVirtualNetworkConnection"
}
`, VirtualHubConnectionResource{}.basic(data))
}

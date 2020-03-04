package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func testAccDataSourceAzureRMExpressRoute_basicMetered(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_express_route_circuit", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMExpressRoute_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "service_provider_properties.0.service_provider_name", "Equinix"),
					resource.TestCheckResourceAttr(data.ResourceName, "service_provider_properties.0.peering_location", "Silicon Valley"),
					resource.TestCheckResourceAttr(data.ResourceName, "service_provider_properties.0.bandwidth_in_mbps", "50"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.tier", "Standard"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.family", "MeteredData"),
					resource.TestCheckResourceAttr(data.ResourceName, "service_provider_provisioning_state", "NotProvisioned"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMExpressRoute_basic(data acceptance.TestData) string {
	config := testAccAzureRMExpressRouteCircuit_basicMeteredConfig(data)

	return fmt.Sprintf(`
%s

data "azurerm_express_route_circuit" "test" {
  resource_group_name = azurerm_resource_group.test.name
  name                = azurerm_express_route_circuit.test.name
}
`, config)
}

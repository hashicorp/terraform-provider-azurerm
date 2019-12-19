package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func testAccDataSourceAzureRMExpressRoute_basicMetered(t *testing.T) {
	dataSourceName := "data.azurerm_express_route_circuit.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMExpressRouteCircuitDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMExpressRoute_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "service_provider_properties.0.service_provider_name", "Equinix"),
					resource.TestCheckResourceAttr(dataSourceName, "service_provider_properties.0.peering_location", "Silicon Valley"),
					resource.TestCheckResourceAttr(dataSourceName, "service_provider_properties.0.bandwidth_in_mbps", "50"),
					resource.TestCheckResourceAttr(dataSourceName, "sku.0.tier", "Standard"),
					resource.TestCheckResourceAttr(dataSourceName, "sku.0.family", "MeteredData"),
					resource.TestCheckResourceAttr(dataSourceName, "service_provider_provisioning_state", "NotProvisioned"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMExpressRoute_basic(rInt int, location string) string {
	config := testAccAzureRMExpressRouteCircuit_basicMeteredConfig(rInt, location)

	return fmt.Sprintf(`
%s

data "azurerm_express_route_circuit" test {
  resource_group_name = "${azurerm_resource_group.test.name}"
  name                = "${azurerm_express_route_circuit.test.name}"
}
`, config)
}

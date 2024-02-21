// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ExpressRouteCircuitPeeringDataSource struct{}

func testAccDataSourceExpressRouteCircuitPeering_privatePeering(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_peering", "test")
	r := ExpressRouteCircuitPeeringResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.privatePeering(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("peering_type").HasValue("AzurePrivatePeering"),
				check.That(data.ResourceName).Key("microsoft_peering_config.#").HasValue("0"),
			),
		},
		data.ImportStep("shared_key"), // is not returned by the API
	})
}

func (r ExpressRouteCircuitPeeringDataSource) privatePeering(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_express_route_circuit_peering" "data" {
  peering_type                  = azurerm_express_route_circuit_peering.test.peering_type
  express_route_circuit_name    = azurerm_express_route_circuit_peering.test.express_route_circuit_name
  resource_group_name           = azurerm_express_route_circuit_peering.test.resource_group_name
}
`, r.privatePeering(data))
}

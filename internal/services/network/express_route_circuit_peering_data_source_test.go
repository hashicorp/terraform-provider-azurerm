// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type ExpressRouteCircuitPeeringDataSource struct{}

func testAccDataSourceExpressRouteCircuitPeering_privatePeering(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_express_route_circuit_peering", "test")
	d := ExpressRouteCircuitPeeringDataSource{}

	data.DataSourceTestInSequence(t, []acceptance.TestStep{
		{
			Config: d.privatePeering(data),
		},
		data.ImportStep("shared_key"), // is not returned by the API
	})
}

func (d ExpressRouteCircuitPeeringDataSource) privatePeering(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_express_route_circuit_peering" "data" {
  peering_type               = azurerm_express_route_circuit_peering.test.peering_type
  express_route_circuit_name = azurerm_express_route_circuit_peering.test.express_route_circuit_name
  resource_group_name        = azurerm_express_route_circuit_peering.test.resource_group_name
}
`, ExpressRouteCircuitPeeringResource{}.privatePeering(data))
}

package network_test

import (
    "fmt"
    "testing"

    "github.com/hashicorp/terraform-plugin-sdk/helper/resource"
    "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

type ExpressRouteCircuitConnectionDataSource struct{}

var config  = `
provider "azurerm" {
  features {}
}

data "azurerm_express_route_circuit_connection" "example" {
  name = "xz-test-conn"
  resource_group_name = "xz3-test"
  circuit_name = "tf-er"
  peering_name = "AzurePrivatePeering"
}
`
func TestAccExpressRouteCircuitConnectionDataSource_ddd(t *testing.T) {
    data := acceptance.BuildTestData(t, "data.azurerm_express_route_circuit_connection", "test")

    data.DataSourceTest(t, []resource.TestStep{
        {
            Config: config,
            Check: resource.ComposeTestCheckFunc(

            ),
        },
    })
}

func TestAccExpressRouteCircuitConnectionDataSource_basic(t *testing.T) {
    data := acceptance.BuildTestData(t, "data.azurerm_express_route_circuit_connection", "test")
    r := ExpressRouteCircuitConnectionDataSource{}
    data.DataSourceTest(t, []resource.TestStep{
        {
            Config: r.basic(data),
            Check: resource.ComposeTestCheckFunc(

            ),
        },
    })
}
func(ExpressRouteCircuitConnectionDataSource) basic(data acceptance.TestData) string {
    return fmt.Sprintf(`
%s

data "azurerm_express_route_circuit_connection" "test" {
  name = azurerm_express_route_circuit_connection.test.name
  resource_group_name = azurerm_express_route_circuit_connection.test.resource_group_name
  circuit_name = azurerm_express_route_circuit_connection.test.circuit_name
  peering_name = azurerm_express_route_circuit_connection.test.peering_name
}
`, ExpressRouteCircuitConnectionResource{}.basic(data))
}

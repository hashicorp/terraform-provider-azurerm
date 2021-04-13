package network_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

type ExpressRouteCircuitConnectionDataSource struct{}

func TestAccExpressRouteCircuitConnectionDataSource_basic(t *testing.T) {
	if os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP") == "" || os.Getenv("ARM_TEST_CIRCUIT_NAME_FIRST") == "" ||
		os.Getenv("ARM_TEST_CIRCUIT_NAME_SECOND") == "" {
		t.Skip("Skipping as ARM_TEST_DATA_RESOURCE_GROUP and/or ARM_TEST_CIRCUIT_NAME_FIRST and/or ARM_TEST_CIRCUIT_NAME_SECOND are not specified")
		return
	}

	data := acceptance.BuildTestData(t, "data.azurerm_express_route_circuit_connection", "test")
	r := ExpressRouteCircuitConnectionDataSource{}
	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}
func (ExpressRouteCircuitConnectionDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_express_route_circuit_connection" "test" {
  name = azurerm_express_route_circuit_connection.test.name
  resource_group_name = azurerm_express_route_circuit_connection.test.resource_group_name
  circuit_name = azurerm_express_route_circuit_connection.test.circuit_name
}
`, ExpressRouteCircuitConnectionResource{}.basic(data))
}

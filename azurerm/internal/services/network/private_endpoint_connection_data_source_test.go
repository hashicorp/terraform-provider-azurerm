package network_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type PrivateEndpointConnectionDataSource struct {
}

func TestAccDataSourcePrivateEndpointConnection_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_private_endpoint_connection", "test")
	r := PrivateEndpointConnectionDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("private_service_connection.0.status").HasValue("Approved"),
			),
		},
	})
}

func (PrivateEndpointConnectionDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_private_endpoint_connection" "test" {
  name                = azurerm_private_endpoint.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, PrivateEndpointResource{}.basic(data))
}

package network_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type PrivateLinkServiceEndpointConnectionDataSource struct {
}

func TestAccDataSourcePrivateLinkServiceEndpointConnections_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_private_link_service_endpoint_connections", "test")
	r := PrivateLinkServiceEndpointConnectionDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("private_endpoint_connections.0.action_required").HasValue("None"),
				check.That(data.ResourceName).Key("private_endpoint_connections.0.description").HasValue("Approved"),
				check.That(data.ResourceName).Key("private_endpoint_connections.0.status").HasValue("Approved"),
				check.That(data.ResourceName).Key("private_endpoint_connections.0.connection_id").Exists(),
				check.That(data.ResourceName).Key("private_endpoint_connections.0.connection_name").Exists(),
				check.That(data.ResourceName).Key("private_endpoint_connections.0.private_endpoint_id").Exists(),
				check.That(data.ResourceName).Key("private_endpoint_connections.0.private_endpoint_name").Exists(),
			),
		},
	})
}

func (PrivateLinkServiceEndpointConnectionDataSource) complete(data acceptance.TestData) string {
	// azurerm_private_link_service_endpoint_connections depends on azurerm_private_endpoint, we deliberately introduce
	// this dependency here via reference, rather than using `depends_on` since `depends_on` on data source will make
	// it never converge.
	return fmt.Sprintf(`
%s

data "azurerm_private_link_service_endpoint_connections" "test" {
  service_id          = azurerm_private_endpoint.test.private_service_connection.0.private_connection_resource_id
  resource_group_name = azurerm_resource_group.test.name
}
`, PrivateLinkServiceResource{}.basic(data))
}

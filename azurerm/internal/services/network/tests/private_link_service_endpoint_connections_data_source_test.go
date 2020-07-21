package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourcePrivateLinkServiceEndpointConnections_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_private_link_service_endpoint_connections", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePrivateLinkServiceEndpointConnections_complete(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "private_endpoint_connections.0.action_required", "None"),
					resource.TestCheckResourceAttr(data.ResourceName, "private_endpoint_connections.0.description", "Approved"),
					resource.TestCheckResourceAttr(data.ResourceName, "private_endpoint_connections.0.status", "Approved"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_endpoint_connections.0.connection_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_endpoint_connections.0.connection_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_endpoint_connections.0.private_endpoint_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "private_endpoint_connections.0.private_endpoint_name"),
				),
			},
		},
	})
}

func testAccDataSourcePrivateLinkServiceEndpointConnections_complete(data acceptance.TestData) string {
	// azurerm_private_link_service_endpoint_connections depends on azurerm_private_endpoint, we deliberately introduce
	// this dependency here via reference, rather than using `depends_on` since `depends_on` on data source will make
	// it never converge.
	return fmt.Sprintf(`
%s

data "azurerm_private_link_service_endpoint_connections" "test" {
  service_id          = azurerm_private_endpoint.test.private_service_connection.0.private_connection_resource_id
  resource_group_name = azurerm_resource_group.test.name
}
`, testAccAzureRMPrivateEndpoint_basic(data))
}

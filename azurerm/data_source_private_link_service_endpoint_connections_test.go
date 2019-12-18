package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourcePrivateLinkServiceEndpointConnections_complete(t *testing.T) {
	dataSourceName := "data.azurerm_private_link_service_endpoint_connections.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePrivateLinkServiceEndpointConnections_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "private_endpoint_connections.0.action_required", "None"),
					resource.TestCheckResourceAttr(dataSourceName, "private_endpoint_connections.0.description", "Approved"),
					resource.TestCheckResourceAttr(dataSourceName, "private_endpoint_connections.0.status", "Approved"),
					resource.TestCheckResourceAttrSet(dataSourceName, "private_endpoint_connections.0.connection_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "private_endpoint_connections.0.connection_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "private_endpoint_connections.0.private_endpoint_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "private_endpoint_connections.0.private_endpoint_name"),
				),
			},
		},
	})
}

func testAccDataSourcePrivateLinkServiceEndpointConnections_complete(rInt int, location string) string {
	return fmt.Sprintf(`
%s

data "azurerm_private_link_service_endpoint_connections" "test" {
  service_id          = azurerm_private_link_service.test.id
	resource_group_name = azurerm_resource_group.test.name
	depends_on          = [azurerm_private_link_endpoint.test,]
}
`, testAccAzureRMPrivateEndpoint_basic(rInt, location))
}

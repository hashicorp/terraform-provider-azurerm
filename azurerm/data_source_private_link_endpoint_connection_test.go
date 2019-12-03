package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMPrivateLinkEndpointConnection_complete(t *testing.T) {
	dataSourceName := "data.azurerm_private_link_endpoint_connection.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePrivateLinkEndpointConnection_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "private_service_connection.0.status", "Approved"),
				),
			},
		},
	})
}

func testAccDataSourcePrivateLinkEndpointConnection_complete(rInt int, location string) string {
	return fmt.Sprintf(`
%s

data "azurerm_private_link_endpoint_connection" "test" {
  name                = azurerm_private_link_endpoint.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, testAccAzureRMPrivateEndpoint_basic(rInt, location))
}

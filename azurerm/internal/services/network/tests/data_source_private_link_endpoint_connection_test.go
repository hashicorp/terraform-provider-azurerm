package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMPrivateLinkEndpointConnection_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_private_link_endpoint_connection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePrivateLinkEndpointConnection_complete(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "private_service_connection.0.status", "Approved"),
				),
			},
		},
	})
}

func testAccDataSourcePrivateLinkEndpointConnection_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_private_link_endpoint_connection" "test" {
  name                = azurerm_private_link_endpoint.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, testAccAzureRMPrivateLinkEndpoint_basic(data))
}

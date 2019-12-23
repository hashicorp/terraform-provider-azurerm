package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMAppServiceVirtualNetworkConnectionGateway_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_service_virtual_network_connection_gateway", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppServiceVirtualNetworkConnectionGateway_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "certificate_thumbprint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "certificate_blob"),
				),
			},
		},
	})
}

func testAccDataSourceAppServiceVirtualNetworkConnectionGateway_basic(data acceptance.TestData) string {
	config := testAccAzureRMAppServiceVirtualNetworkConnectionGateway_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_app_service_virtual_network_connection_gateway" "test" {
	app_service_name      = "${azurerm_app_service_virtual_network_connection_gateway.test.app_service_name}"
	resource_group_name   = "${azurerm_app_service_virtual_network_connection_gateway.test.resource_group_name}"
	virtual_network_name  = "${azurerm_virtual_network.test.name}"
}
`, config)
}

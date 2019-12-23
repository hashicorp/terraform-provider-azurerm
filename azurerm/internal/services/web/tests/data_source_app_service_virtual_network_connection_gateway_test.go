package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMAppServiceVirtualNetworkConnectionGateway_basic(t *testing.T) {
	dataSourceName := "data.azurerm_app_service_virtual_network_connection_gateway.example"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppServiceVirtualNetworkConnectionGateway_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "certificate_thumbprint"),
					resource.TestCheckResourceAttrSet(dataSourceName, "certificate_blob"),
				),
			},
		},
	})
}

func testAccDataSourceAppServiceVirtualNetworkConnectionGateway_basic(rInt int, location string) string {
	config := testAccAzureRMAppServiceVirtualNetworkConnectionGateway_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_app_service_virtual_network_connection_gateway" "example" {
	app_service_name      = "${azurerm_app_service_virtual_network_connection_gateway.example.app_service_name}"
	resource_group_name   = "${azurerm_app_service_virtual_network_connection_gateway.example.resource_group_name}"
	virtual_network_name  = "${azurerm_virtual_network.example.name}"
}
`, config)
}

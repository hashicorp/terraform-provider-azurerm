package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMSignalRService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_signalr_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSignalRServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMSignalRService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "hostname"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_port"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "server_port"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_connection_string"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMSignalRService_basic(data acceptance.TestData) string {
	template := testAccAzureRMSignalRService_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_signalr_service" "test" {
  name                = azurerm_signalr_service.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}

package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMAppServiceEnvironment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_service_environment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppServiceEnvironment_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "front_end_scale_factor"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "pricing_tier"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "internal_ip_address"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "service_ip_address"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "outbound_ip_addresses"),
				),
			},
		},
	})
}

func testAccDataSourceAppServiceEnvironment_basic(data acceptance.TestData) string {
	config := testAccAzureRMAppServiceEnvironment_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_app_service_environment" "test" {
  name                = azurerm_app_service_environment.test.name
  resource_group_name = azurerm_app_service_environment.test.resource_group_name
}
`, config)
}

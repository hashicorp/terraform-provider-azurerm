package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMSpringCloudService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_spring_cloud_service", "test")

	//lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSpringCloudService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "id"),
				),
			},
		},
	})
}

func testAccDataSourceSpringCloudService_basic(data acceptance.TestData) string {
	config := testAccAzureRMSpringCloudService_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_spring_cloud_service" "test" {
  name                = azurerm_spring_cloud_service.test.name
  resource_group_name = azurerm_spring_cloud_service.test.resource_group_name
}
`, config)
}

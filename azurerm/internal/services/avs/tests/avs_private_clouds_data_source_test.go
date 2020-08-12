package tests

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"testing"
)

func TestAccDataSourceAzureRMavsPrivateCloud_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_avs_private_cloud", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMavsPrivateCloudDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceavsPrivateCloud_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMavsPrivateCloudExists(data.ResourceName),
				),
			},
		},
	})
}

func testAccDataSourceavsPrivateCloud_basic(data acceptance.TestData) string {
	config := testAccAzureRMavsPrivateCloud_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_avs_private_cloud" "test" {
  name = azurerm_avs_private_cloud.test.name
  resource_group_name = azurerm_avs_private_cloud.test.resource_group_name
}
`, config)
}

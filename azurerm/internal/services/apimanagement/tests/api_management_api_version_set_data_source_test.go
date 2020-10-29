package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMApiManagementApiVersionSet_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_api_management_api_version_set", "test")

	//lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceApiManagementApiVersionSet_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "api_management_name"),
				),
			},
		},
	})
}

func testAccDataSourceApiManagementApiVersionSet_basic(data acceptance.TestData) string {
	config := testAccAzureRMApiManagementApiVersionSet_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_api_management_api_version_set" "test" {
  name                = azurerm_api_management_api_version_set.test.name
  resource_group_name = azurerm_api_management_api_version_set.test.resource_group_name
  api_management_name = azurerm_api_management_api_version_set.test.api_management_name
}
`, config)
}

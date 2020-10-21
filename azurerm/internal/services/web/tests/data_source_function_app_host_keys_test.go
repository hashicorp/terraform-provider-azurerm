package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMFunctionAppHostKeys_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_function_app_host_keys", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMFunctionAppHostKeys_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(data.ResourceName),
					testCheckAzureRMFunctionAppHasNoContentShare(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "default_function_key"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMFunctionAppHostKeys_basic(data acceptance.TestData) string {
	template := testAccAzureRMFunctionApp_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_function_app_host_keys" "test" {
  name                = azurerm_function_app.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}

package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMApiManagementApi_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_api_management_api", "test")

	//lintignore:AT001
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceApiManagementApi_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "display_name", "api1"),
					resource.TestCheckResourceAttr(data.ResourceName, "path", "api1"),
					resource.TestCheckResourceAttr(data.ResourceName, "protocols.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "protocols.0", "https"),
					resource.TestCheckResourceAttr(data.ResourceName, "soap_pass_through", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "subscription_required", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "is_current", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "is_online", "false"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMApiManagementApi_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_api_management_api", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceApiManagementApi_complete(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "display_name", "Butter Parser"),
					resource.TestCheckResourceAttr(data.ResourceName, "path", "butter-parser"),
					resource.TestCheckResourceAttr(data.ResourceName, "protocols.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "What is my purpose? You parse butter."),
					resource.TestCheckResourceAttr(data.ResourceName, "service_url", "https://example.com/foo/bar"),
					resource.TestCheckResourceAttr(data.ResourceName, "soap_pass_through", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "subscription_key_parameter_names.0.header", "X-Butter-Robot-API-Key"),
					resource.TestCheckResourceAttr(data.ResourceName, "subscription_key_parameter_names.0.query", "location"),
					resource.TestCheckResourceAttr(data.ResourceName, "is_current", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "is_online", "false"),
				),
			},
		},
	})
}

func testAccDataSourceApiManagementApi_basic(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementApi_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_api_management_api" "test" {
  name                = azurerm_api_management_api.test.name
  api_management_name = azurerm_api_management_api.test.api_management_name
  resource_group_name = azurerm_api_management_api.test.resource_group_name
  revision            = azurerm_api_management_api.test.revision
}
`, template)
}

func testAccDataSourceApiManagementApi_complete(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementApi_complete(data)
	return fmt.Sprintf(`
%s

data "azurerm_api_management_api" "test" {
  name                = azurerm_api_management_api.test.name
  api_management_name = azurerm_api_management_api.test.api_management_name
  resource_group_name = azurerm_api_management_api.test.resource_group_name
  revision            = azurerm_api_management_api.test.revision
}
`, template)
}

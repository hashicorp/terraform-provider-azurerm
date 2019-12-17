package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMApiManagementApi_basic(t *testing.T) {
	dataSourceName := "data.azurerm_api_management_api.test"
	rInt := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceApiManagementApi_basic(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "display_name", "api1"),
					resource.TestCheckResourceAttr(dataSourceName, "path", "api1"),
					resource.TestCheckResourceAttr(dataSourceName, "protocols.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "protocols.0", "https"),
					resource.TestCheckResourceAttr(dataSourceName, "soap_pass_through", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "is_current", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "is_online", "false"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMApiManagementApi_complete(t *testing.T) {
	dataSourceName := "data.azurerm_api_management_api.test"
	rInt := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceApiManagementApi_complete(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "display_name", "Butter Parser"),
					resource.TestCheckResourceAttr(dataSourceName, "path", "butter-parser"),
					resource.TestCheckResourceAttr(dataSourceName, "protocols.#", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "description", "What is my purpose? You parse butter."),
					resource.TestCheckResourceAttr(dataSourceName, "service_url", "https://example.com/foo/bar"),
					resource.TestCheckResourceAttr(dataSourceName, "soap_pass_through", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "subscription_key_parameter_names.0.header", "X-Butter-Robot-API-Key"),
					resource.TestCheckResourceAttr(dataSourceName, "subscription_key_parameter_names.0.query", "location"),
					resource.TestCheckResourceAttr(dataSourceName, "is_current", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "is_online", "false"),
				),
			},
		},
	})
}

func testAccDataSourceApiManagementApi_basic(rInt int, location string) string {
	template := testAccAzureRMApiManagementApi_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_api_management_api" "test" {
  name                = "${azurerm_api_management_api.test.name}"
  api_management_name = "${azurerm_api_management_api.test.api_management_name}"
  resource_group_name = "${azurerm_api_management_api.test.resource_group_name}"
  revision            = "${azurerm_api_management_api.test.revision}"
}
`, template)
}

func testAccDataSourceApiManagementApi_complete(rInt int, location string) string {
	template := testAccAzureRMApiManagementApi_complete(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_api_management_api" "test" {
  name                = "${azurerm_api_management_api.test.name}"
  api_management_name = "${azurerm_api_management_api.test.api_management_name}"
  resource_group_name = "${azurerm_api_management_api.test.resource_group_name}"
  revision            = "${azurerm_api_management_api.test.revision}"
}
`, template)
}

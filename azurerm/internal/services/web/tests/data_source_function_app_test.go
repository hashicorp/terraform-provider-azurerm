package tests

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceAzureRMFunctionApp_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_function_app", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMFunctionApp_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(data.ResourceName),
					testCheckAzureRMFunctionAppHasNoContentShare(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "outbound_ip_addresses"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "possible_outbound_ip_addresses"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMFunctionApp_appSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_function_app", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMFunctionApp_appSettings(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "app_settings.hello", "world"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMFunctionApp_connectionStrings(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_function_app", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMFunctionApp_connectionStrings(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.0.name", "Example"),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.0.value", "some-postgresql-connection-string"),
					resource.TestCheckResourceAttr(data.ResourceName, "connection_string.0.type", "PostgreSQL"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMFunctionApp_basic(data acceptance.TestData) string {
	template := testAccAzureRMFunctionApp_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_function_app" "test" {
  name                = "${azurerm_function_app.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, template)
}

func testAccDataSourceAzureRMFunctionApp_connectionStrings(data acceptance.TestData) string {
	template := testAccAzureRMFunctionApp_connectionStrings(data)
	return fmt.Sprintf(`
%s

data "azurerm_function_app" "test" {
  name                = "${azurerm_function_app.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, template)
}

func testAccDataSourceAzureRMFunctionApp_appSettings(data acceptance.TestData) string {
	template := testAccAzureRMFunctionApp_appSettings(data)
	return fmt.Sprintf(`
%s

data "azurerm_function_app" "test" {
  name                = "${azurerm_function_app.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, template)
}

package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
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
					resource.TestCheckResourceAttrSet(data.ResourceName, "custom_domain_verification_id"),
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

func TestAccDataSourceAzureRMFunctionApp_withSourceControl(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_function_app", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMFunctionApp_withSourceControl(data, "main"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "source_control.0.branch", "main"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMFunctionApp_siteConfig(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFunctionAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMFunctionApp_withSiteConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFunctionAppExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.0.ip_address", "10.10.10.10/32"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.1.ip_address", "20.20.20.0/24"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.2.ip_address", "30.30.0.0/16"),
					resource.TestCheckResourceAttr(data.ResourceName, "site_config.0.ip_restriction.3.ip_address", "192.168.1.2/24"),
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
  name                = azurerm_function_app.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}

func testAccDataSourceAzureRMFunctionApp_connectionStrings(data acceptance.TestData) string {
	template := testAccAzureRMFunctionApp_connectionStrings(data)
	return fmt.Sprintf(`
%s

data "azurerm_function_app" "test" {
  name                = azurerm_function_app.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}

func testAccDataSourceAzureRMFunctionApp_appSettings(data acceptance.TestData) string {
	template := testAccAzureRMFunctionApp_appSettings(data)
	return fmt.Sprintf(`
%s

data "azurerm_function_app" "test" {
  name                = azurerm_function_app.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}

func testAccDataSourceAzureRMFunctionApp_withSourceControl(data acceptance.TestData, branch string) string {
	config := testAccAzureRMFunctionApp_withSourceControl(data, branch)
	return fmt.Sprintf(`
%s

data "azurerm_function_app" "test" {
  name                = azurerm_function_app.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, config)
}

func testAccDataSourceAzureRMFunctionApp_withSiteConfig(data acceptance.TestData) string {
	config := testAccAzureRMFunctionApp_manyIpRestrictions(data)
	return fmt.Sprintf(`
%s

data "azurerm_function_app" "test" {
  name                = azurerm_function_app.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, config)
}

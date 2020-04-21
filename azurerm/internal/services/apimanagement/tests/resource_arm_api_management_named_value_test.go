package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAPIManagementNamedValue_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_named_value", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAPIManagementNamedValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAPIManagementNamedValue_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementNamedValueExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "display_name", fmt.Sprintf("TestProperty%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "value", "Test Value"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.0", "tag1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.1", "tag2"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAPIManagementNamedValue_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_named_value", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAPIManagementNamedValueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAPIManagementNamedValue_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementNamedValueExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "display_name", fmt.Sprintf("TestProperty%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "value", "Test Value"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.0", "tag1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.1", "tag2"),
				),
			},
			{
				Config: testAccAzureRMAPIManagementNamedValue_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementNamedValueExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "display_name", fmt.Sprintf("TestProperty2%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "value", "Test Value2"),
					resource.TestCheckResourceAttr(data.ResourceName, "secret", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.0", "tag3"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.1", "tag4"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMAPIManagementNamedValueDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.NamedValueClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_named_value" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		resp, err := client.Get(ctx, resourceGroup, serviceName, name)

		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}

		return nil
	}
	return nil
}

func testCheckAzureRMAPIManagementNamedValueExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.NamedValueClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		resp, err := client.Get(ctx, resourceGroup, serviceName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: API Management Property %q (Resource Group %q / API Management Service %q) does not exist", name, resourceGroup, serviceName)
			}
			return fmt.Errorf("Bad: Get on apiManagement.NamedValueClient: %+v", err)
		}

		return nil
	}
}

/*

 */

func testAccAzureRMAPIManagementNamedValue_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"
}

resource "azurerm_api_management_named_value" "test" {
  name                = "acctestAMProperty-%d"
  resource_group_name = azurerm_api_management.test.resource_group_name
  api_management_name = azurerm_api_management.test.name
  display_name        = "TestProperty%d"
  value               = "Test Value"
  tags                = ["tag1", "tag2"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAPIManagementNamedValue_update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"
}

resource "azurerm_api_management_named_value" "test" {
  name                = "acctestAMProperty-%d"
  resource_group_name = azurerm_api_management.test.resource_group_name
  api_management_name = azurerm_api_management.test.name
  display_name        = "TestProperty2%d"
  value               = "Test Value2"
  secret              = true
  tags                = ["tag3", "tag4"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

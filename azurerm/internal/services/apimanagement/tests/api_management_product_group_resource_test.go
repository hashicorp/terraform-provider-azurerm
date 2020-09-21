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

func TestAccAzureRMApiManagementProductGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_product_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAPIManagementProductGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementProductGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementProductGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagementProductGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_product_group", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAPIManagementProductGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementProductGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementProductGroupExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMApiManagementProductGroup_requiresImport),
		},
	})
}

func testCheckAzureRMAPIManagementProductGroupDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.ProductGroupsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_product_group" {
			continue
		}

		productId := rs.Primary.Attributes["product_id"]
		groupName := rs.Primary.Attributes["group_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		resp, err := client.CheckEntityExists(ctx, resourceGroup, serviceName, productId, groupName)
		if err != nil {
			if !utils.ResponseWasNotFound(resp) {
				return err
			}
		}

		return nil
	}
	return nil
}

func testCheckAzureRMAPIManagementProductGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.ProductGroupsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		productId := rs.Primary.Attributes["product_id"]
		groupName := rs.Primary.Attributes["group_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		resp, err := client.CheckEntityExists(ctx, resourceGroup, serviceName, productId, groupName)
		if err != nil {
			if utils.ResponseWasNotFound(resp) {
				return fmt.Errorf("Bad: Product %q / Group %q (API Management Service %q / Resource Group %q) does not exist", productId, groupName, serviceName, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on apiManagement.ProductGroupsClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMApiManagementProductGroup_basic(data acceptance.TestData) string {
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

resource "azurerm_api_management_product" "test" {
  product_id            = "test-product"
  api_management_name   = azurerm_api_management.test.name
  resource_group_name   = azurerm_resource_group.test.name
  display_name          = "Test Product"
  subscription_required = true
  approval_required     = false
  published             = true
}

resource "azurerm_api_management_group" "test" {
  name                = "acctestAMGroup-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "Test Group"
}

resource "azurerm_api_management_product_group" "test" {
  product_id          = azurerm_api_management_product.test.product_id
  group_name          = azurerm_api_management_group.test.name
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApiManagementProductGroup_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementProductGroup_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_product_group" "import" {
  product_id          = azurerm_api_management_product_group.test.product_id
  group_name          = azurerm_api_management_product_group.test.group_name
  api_management_name = azurerm_api_management_product_group.test.api_management_name
  resource_group_name = azurerm_api_management_product_group.test.resource_group_name
}
`, template)
}

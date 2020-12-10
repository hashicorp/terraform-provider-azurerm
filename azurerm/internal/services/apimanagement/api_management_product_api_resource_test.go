package apimanagement_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMApiManagementProductApi_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_product_api", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAPIManagementProductApiDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementProductApi_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementProductApiExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagementProductApi_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_product_api", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAPIManagementProductApiDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementProductApi_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementProductApiExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMApiManagementProductApi_requiresImport),
		},
	})
}

func testCheckAzureRMAPIManagementProductApiDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.ProductApisClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_product_api" {
			continue
		}

		apiName := rs.Primary.Attributes["api_name"]
		productId := rs.Primary.Attributes["product_id"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		resp, err := client.CheckEntityExists(ctx, resourceGroup, serviceName, productId, apiName)
		if err != nil {
			if !utils.ResponseWasNotFound(resp) {
				return err
			}
		}

		return nil
	}
	return nil
}

func testCheckAzureRMAPIManagementProductApiExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.ProductApisClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		apiName := rs.Primary.Attributes["api_name"]
		productId := rs.Primary.Attributes["product_id"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		resp, err := client.CheckEntityExists(ctx, resourceGroup, serviceName, productId, apiName)
		if err != nil {
			if utils.ResponseWasNotFound(resp) {
				return fmt.Errorf("Bad: API %q / Product %q (API Management Service %q / Resource Group %q) does not exist", apiName, productId, serviceName, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on apiManagement.ProductApisClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMApiManagementProductApi_basic(data acceptance.TestData) string {
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

resource "azurerm_api_management_api" "test" {
  name                = "acctestapi-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "api1"
  path                = "api1"
  protocols           = ["https"]
  revision            = "1"
}

resource "azurerm_api_management_product_api" "test" {
  product_id          = azurerm_api_management_product.test.product_id
  api_name            = azurerm_api_management_api.test.name
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApiManagementProductApi_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementProductApi_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_product_api" "import" {
  api_name            = azurerm_api_management_product_api.test.api_name
  product_id          = azurerm_api_management_product_api.test.product_id
  api_management_name = azurerm_api_management_product_api.test.api_management_name
  resource_group_name = azurerm_api_management_product_api.test.resource_group_name
}
`, template)
}

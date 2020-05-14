package tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMApiManagementProduct_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_product", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementProductDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementProduct_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementProductExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "approval_required", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "display_name", "Test Product"),
					resource.TestCheckResourceAttr(data.ResourceName, "product_id", "test-product"),
					resource.TestCheckResourceAttr(data.ResourceName, "published", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "subscription_required", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "terms", ""),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagementProduct_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_product", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementProductDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementProduct_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementProductExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMApiManagementProduct_requiresImport),
		},
	})
}

func testCheckAzureRMApiManagementProductDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.ProductsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_product" {
			continue
		}

		productId := rs.Primary.Attributes["product_id"]
		serviceName := rs.Primary.Attributes["api_management_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		resp, err := conn.Get(ctx, resourceGroup, serviceName, productId)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return nil
	}

	return nil
}

func TestAccAzureRMApiManagementProduct_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_product", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementProductDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementProduct_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementProductExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "approval_required", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "display_name", "Test Product"),
					resource.TestCheckResourceAttr(data.ResourceName, "product_id", "test-product"),
					resource.TestCheckResourceAttr(data.ResourceName, "published", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "subscription_required", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "terms", ""),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMApiManagementProduct_updated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementProductExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "approval_required", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "display_name", "Test Updated Product"),
					resource.TestCheckResourceAttr(data.ResourceName, "product_id", "test-product"),
					resource.TestCheckResourceAttr(data.ResourceName, "published", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "subscription_required", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "terms", ""),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMApiManagementProduct_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementProductExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "description", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "display_name", "Test Product"),
					resource.TestCheckResourceAttr(data.ResourceName, "product_id", "test-product"),
					resource.TestCheckResourceAttr(data.ResourceName, "published", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "subscription_required", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "terms", ""),
				),
			},
		},
	})
}

func TestAccAzureRMApiManagementProduct_subscriptionsLimit(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_product", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementProductDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementProduct_subscriptionLimits(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementProductExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "approval_required", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "subscription_required", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "subscriptions_limit", "2"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagementProduct_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_product", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementProductDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementProduct_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementProductExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "approval_required", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "This is an example description"),
					resource.TestCheckResourceAttr(data.ResourceName, "display_name", "Test Product"),
					resource.TestCheckResourceAttr(data.ResourceName, "product_id", "test-product"),
					resource.TestCheckResourceAttr(data.ResourceName, "published", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "subscriptions_limit", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "subscription_required", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "terms", "These are some example terms and conditions"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagementProduct_approvalRequiredError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_product", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementProductDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementProduct_approvalRequiredError(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementProductExists(data.ResourceName)),
				ExpectError: regexp.MustCompile("`subscription_required` must be true and `subscriptions_limit` must be greater than 0 to use `approval_required`"),
			},
		},
	})
}

func testCheckAzureRMApiManagementProductExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.ProductsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		productId := rs.Primary.Attributes["product_id"]
		serviceName := rs.Primary.Attributes["api_management_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, serviceName, productId)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Product %q (API Management Service %q / Resource Group %q) does not exist", productId, serviceName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on apiManagement.ProductsClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMApiManagementProduct_basic(data acceptance.TestData) string {
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
  subscription_required = false
  published             = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMApiManagementProduct_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementProduct_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_product" "import" {
  product_id            = azurerm_api_management_product.test.product_id
  api_management_name   = azurerm_api_management_product.test.api_management_name
  resource_group_name   = azurerm_api_management_product.test.resource_group_name
  display_name          = azurerm_api_management_product.test.display_name
  subscription_required = azurerm_api_management_product.test.subscription_required
  approval_required     = azurerm_api_management_product.test.approval_required
  published             = azurerm_api_management_product.test.published
}
`, template)
}

func testAccAzureRMApiManagementProduct_updated(data acceptance.TestData) string {
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
  display_name          = "Test Updated Product"
  subscription_required = true
  approval_required     = true
  subscriptions_limit   = 1
  published             = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMApiManagementProduct_subscriptionLimits(data acceptance.TestData) string {
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
  approval_required     = true
  subscriptions_limit   = 2
  published             = false
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMApiManagementProduct_complete(data acceptance.TestData) string {
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
  approval_required     = true
  published             = true
  subscriptions_limit   = 2
  description           = "This is an example description"
  terms                 = "These are some example terms and conditions"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMApiManagementProduct_approvalRequiredError(data acceptance.TestData) string {
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
  sku_name            = "Developer_1"
}

resource "azurerm_api_management_product" "test" {
  product_id            = "test-product"
  api_management_name   = azurerm_api_management.test.name
  resource_group_name   = azurerm_resource_group.test.name
  display_name          = "Test Product"
  approval_required     = true
  subscription_required = false
  published             = true
  description           = "This is an example description"
  terms                 = "These are some example terms and conditions"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

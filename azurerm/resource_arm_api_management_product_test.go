package azurerm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMApiManagementProduct_basic(t *testing.T) {
	resourceName := "azurerm_api_management_product.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementProductDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementProduct_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementProductExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "approval_required", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Test Product"),
					resource.TestCheckResourceAttr(resourceName, "product_id", "test-product"),
					resource.TestCheckResourceAttr(resourceName, "published", "false"),
					resource.TestCheckResourceAttr(resourceName, "subscription_required", "false"),
					resource.TestCheckResourceAttr(resourceName, "terms", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMApiManagementProduct_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_api_management_product.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementProductDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementProduct_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementProductExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMApiManagementProduct_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_api_management_product"),
			},
		},
	})
}

func testCheckAzureRMApiManagementProductDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).apiManagement.ProductsClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_product" {
			continue
		}

		productId := rs.Primary.Attributes["product_id"]
		serviceName := rs.Primary.Attributes["api_management_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
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
	resourceName := "azurerm_api_management_product.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementProductDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementProduct_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementProductExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "approval_required", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Test Product"),
					resource.TestCheckResourceAttr(resourceName, "product_id", "test-product"),
					resource.TestCheckResourceAttr(resourceName, "published", "false"),
					resource.TestCheckResourceAttr(resourceName, "subscription_required", "false"),
					resource.TestCheckResourceAttr(resourceName, "terms", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMApiManagementProduct_updated(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementProductExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "approval_required", "true"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Test Updated Product"),
					resource.TestCheckResourceAttr(resourceName, "product_id", "test-product"),
					resource.TestCheckResourceAttr(resourceName, "published", "true"),
					resource.TestCheckResourceAttr(resourceName, "subscription_required", "true"),
					resource.TestCheckResourceAttr(resourceName, "terms", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMApiManagementProduct_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementProductExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Test Product"),
					resource.TestCheckResourceAttr(resourceName, "product_id", "test-product"),
					resource.TestCheckResourceAttr(resourceName, "published", "false"),
					resource.TestCheckResourceAttr(resourceName, "subscription_required", "false"),
					resource.TestCheckResourceAttr(resourceName, "terms", ""),
				),
			},
		},
	})
}

func TestAccAzureRMApiManagementProduct_subscriptionsLimit(t *testing.T) {
	resourceName := "azurerm_api_management_product.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementProductDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementProduct_subscriptionLimits(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementProductExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "approval_required", "true"),
					resource.TestCheckResourceAttr(resourceName, "subscription_required", "true"),
					resource.TestCheckResourceAttr(resourceName, "subscriptions_limit", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMApiManagementProduct_complete(t *testing.T) {
	resourceName := "azurerm_api_management_product.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementProductDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementProduct_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementProductExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "approval_required", "true"),
					resource.TestCheckResourceAttr(resourceName, "description", "This is an example description"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Test Product"),
					resource.TestCheckResourceAttr(resourceName, "product_id", "test-product"),
					resource.TestCheckResourceAttr(resourceName, "published", "true"),
					resource.TestCheckResourceAttr(resourceName, "subscriptions_limit", "2"),
					resource.TestCheckResourceAttr(resourceName, "subscription_required", "true"),
					resource.TestCheckResourceAttr(resourceName, "terms", "These are some example terms and conditions"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMApiManagementProduct_approvalRequiredError(t *testing.T) {
	resourceName := "azurerm_api_management_product.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementProductDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementProduct_approvalRequiredError(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementProductExists(resourceName)),
				ExpectError: regexp.MustCompile("`subscription_required` must be true and `subscriptions_limit` must be greater than 0 to use `approval_required`"),
			},
		},
	})
}

func testCheckAzureRMApiManagementProductExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		productId := rs.Primary.Attributes["product_id"]
		serviceName := rs.Primary.Attributes["api_management_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		conn := testAccProvider.Meta().(*ArmClient).apiManagement.ProductsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
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

func testAccAzureRMApiManagementProduct_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku {
    name     = "Developer"
    capacity = 1
  }
}

resource "azurerm_api_management_product" "test" {
  product_id            = "test-product"
  api_management_name   = "${azurerm_api_management.test.name}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  display_name          = "Test Product"
  subscription_required = false
  published             = false
}
`, rInt, location, rInt)
}

func testAccAzureRMApiManagementProduct_requiresImport(rInt int, location string) string {
	template := testAccAzureRMApiManagementProduct_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_product" "import" {
  product_id            = "${azurerm_api_management_product.test.product_id}"
  api_management_name   = "${azurerm_api_management_product.test.api_management_name}"
  resource_group_name   = "${azurerm_api_management_product.test.resource_group_name}"
  display_name          = "${azurerm_api_management_product.test.display_name}"
  subscription_required = "${azurerm_api_management_product.test.subscription_required}"
  approval_required     = "${azurerm_api_management_product.test.approval_required}"
  published             = "${azurerm_api_management_product.test.published}"
}
`, template)
}

func testAccAzureRMApiManagementProduct_updated(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku {
    name     = "Developer"
    capacity = 1
  }
}

resource "azurerm_api_management_product" "test" {
  product_id            = "test-product"
  api_management_name   = "${azurerm_api_management.test.name}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  display_name          = "Test Updated Product"
  subscription_required = true
  approval_required     = true
  subscriptions_limit   = 1
  published             = true
}
`, rInt, location, rInt)
}

func testAccAzureRMApiManagementProduct_subscriptionLimits(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku {
    name     = "Developer"
    capacity = 1
  }
}

resource "azurerm_api_management_product" "test" {
  product_id            = "test-product"
  api_management_name   = "${azurerm_api_management.test.name}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  display_name          = "Test Product"
  subscription_required = true
  approval_required     = true
  subscriptions_limit   = 2
  published             = false
}
`, rInt, location, rInt)
}

func testAccAzureRMApiManagementProduct_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku {
    name     = "Developer"
    capacity = 1
  }
}

resource "azurerm_api_management_product" "test" {
  product_id            = "test-product"
  api_management_name   = "${azurerm_api_management.test.name}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  display_name          = "Test Product"
  subscription_required = true
  approval_required     = true
  published             = true
  subscriptions_limit   = 2
  description           = "This is an example description"
  terms                 = "These are some example terms and conditions"
}
`, rInt, location, rInt)
}

func testAccAzureRMApiManagementProduct_approvalRequiredError(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku {
    name     = "Developer"
    capacity = 1
  }
}

resource "azurerm_api_management_product" "test" {
  product_id            = "test-product"
  api_management_name   = "${azurerm_api_management.test.name}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  display_name          = "Test Product"
  approval_required     = true
  subscription_required = false
  published             = true
  description           = "This is an example description"
  terms                 = "These are some example terms and conditions"
}
`, rInt, location, rInt)
}

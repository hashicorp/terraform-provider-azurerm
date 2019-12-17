package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMApiManagementProductPolicy_basic(t *testing.T) {
	resourceName := "azurerm_api_management_product_policy.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementProductPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementProductPolicy_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementProductPolicyExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"xml_link"},
			},
		},
	})
}

func TestAccAzureRMApiManagementProductPolicy_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_api_management_product_policy.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementProductPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementProductPolicy_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementProductPolicyExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMApiManagementProductPolicy_requiresImport(ri, location),
				ExpectError: acceptance.RequiresImportError("azurerm_api_management_product"),
			},
		},
	})
}

func TestAccAzureRMApiManagementProductPolicy_update(t *testing.T) {
	resourceName := "azurerm_api_management_product_policy.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementProductPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementProductPolicy_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementProductPolicyExists(resourceName),
				),
			},
			{
				Config: testAccAzureRMApiManagementProductPolicy_updated(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementProductPolicyExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"xml_link"},
			},
		},
	})
}

func testCheckAzureRMApiManagementProductPolicyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		productID := rs.Primary.Attributes["product_id"]
		serviceName := rs.Primary.Attributes["api_management_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		conn := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.ProductPoliciesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		resp, err := conn.Get(ctx, resourceGroup, serviceName, productID)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Product Policy (API Management Service %q / Product %q/  Resource Group %q) does not exist", serviceName, productID, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on apiManagement.ProductPoliciesClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMApiManagementProductPolicyDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.ProductPoliciesClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_product_policy" {
			continue
		}

		productID := rs.Primary.Attributes["product_id"]
		serviceName := rs.Primary.Attributes["api_management_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		resp, err := conn.Get(ctx, resourceGroup, serviceName, productID)
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

func testAccAzureRMApiManagementProductPolicy_basic(rInt int, location string) string {
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

resource "azurerm_api_management_product_policy" "test" {
  product_id          = "${azurerm_api_management_product.test.product_id}"
  api_management_name = "${azurerm_api_management.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  xml_link            = "https://gist.githubusercontent.com/riordanp/ca22f8113afae0eb38cc12d718fd048d/raw/d6ac89a2f35a6881a7729f8cb4883179dc88eea1/example.xml"
}
`, rInt, location, rInt)
}

func testAccAzureRMApiManagementProductPolicy_requiresImport(rInt int, location string) string {
	template := testAccAzureRMApiManagementProductPolicy_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_product_policy" "import" {
  product_id          = "${azurerm_api_management_product_policy.test.product_id}"
  api_management_name = "${azurerm_api_management_product_policy.test.api_management_name}"
  resource_group_name = "${azurerm_api_management_product_policy.test.resource_group_name}"
  xml_link            = "${azurerm_api_management_product_policy.test.xml_link}"
}
`, template)
}

func testAccAzureRMApiManagementProductPolicy_updated(rInt int, location string) string {
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

resource "azurerm_api_management_product_policy" "test" {
  product_id          = "${azurerm_api_management_product.test.product_id}"
  api_management_name = "${azurerm_api_management.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  xml_content = <<XML
<policies>
  <inbound>
    <set-variable name="abc" value="@(context.Request.Headers.GetValueOrDefault("X-Header-Name", ""))" />
    <find-and-replace from="xyz" to="abc" />
  </inbound>
</policies>
XML
}
`, rInt, location, rInt)
}

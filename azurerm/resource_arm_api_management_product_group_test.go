package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAPIManagementProductGroup_basic(t *testing.T) {
	resourceName := "azurerm_api_management_product_group.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAPIManagementProductGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAPIManagementProductGroup_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementProductGroupExists(resourceName),
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

func TestAccAzureRMAPIManagementProductGroup_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_api_management_product_group.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAPIManagementProductGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAPIManagementProductGroup_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementProductGroupExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMAPIManagementProductGroup_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_api_management_product_group"),
			},
		},
	})
}

func testCheckAzureRMAPIManagementProductGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).apiManagement.ProductGroupsClient
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_product_group" {
			continue
		}

		productId := rs.Primary.Attributes["product_id"]
		groupName := rs.Primary.Attributes["group_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		ctx := testAccProvider.Meta().(*ArmClient).StopContext
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
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		productId := rs.Primary.Attributes["product_id"]
		groupName := rs.Primary.Attributes["group_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		client := testAccProvider.Meta().(*ArmClient).apiManagement.ProductGroupsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
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

func testAccAzureRMAPIManagementProductGroup_basic(rInt int, location string) string {
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
  approval_required     = false
  published             = true
}

resource "azurerm_api_management_group" "test" {
  name                = "acctestAMGroup-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  api_management_name = "${azurerm_api_management.test.name}"
  display_name        = "Test Group"
}

resource "azurerm_api_management_product_group" "test" {
  product_id          = "${azurerm_api_management_product.test.product_id}"
  group_name          = "${azurerm_api_management_group.test.name}"
  api_management_name = "${azurerm_api_management.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAPIManagementProductGroup_requiresImport(rInt int, location string) string {
	template := testAccAzureRMAPIManagementProductGroup_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_product_group" "import" {
  product_id          = "${azurerm_api_management_product_group.test.product_id}"
  group_name          = "${azurerm_api_management_product_group.test.group_name}"
  api_management_name = "${azurerm_api_management_product_group.test.api_management_name}"
  resource_group_name = "${azurerm_api_management_product_group.test.resource_group_name}"
}
`, template)
}

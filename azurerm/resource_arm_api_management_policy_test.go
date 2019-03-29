package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMApiManagementPolicy_basic(t *testing.T) {
	resourceName := "azurerm_api_management_policy.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementPolicy_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "xml_content", "<policies></policies>"),
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

func TestAccAzureRMApiManagementPolicy_findReplace(t *testing.T) {
	resourceName := "azurerm_api_management_policy.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementPolicy_findReplace(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "xml_content", "<policies>\r\n\t<inbound>\r\n\t\t<find-and-replace from=\"xyz\" to=\"abc\" />\r\n\t</inbound>\r\n</policies>"),
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

func TestAccAzureRMApiManagementPolicy_update(t *testing.T) {
	resourceName := "azurerm_api_management_policy.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementPolicy_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "xml_content", "<policies></policies>"),
				),
			},
			{
				Config: testAccAzureRMApiManagementPolicy_findReplace(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "xml_content", "<policies>\r\n\t<inbound>\r\n\t\t<find-and-replace from=\"xyz\" to=\"abc\" />\r\n\t</inbound>\r\n</policies>"),
				),
			},
			{
				Config: testAccAzureRMApiManagementPolicy_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "xml_content", "<policies></policies>"),
				),
			},
		},
	})
}

func testCheckAzureRMApiManagementPolicyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Api Management Policy not found: %s", resourceName)
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		client := testAccProvider.Meta().(*ArmClient).apiManagementPolicyClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		if resp, err := client.Get(ctx, resourceGroup, serviceName); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Global Policy (API Management Service %q / Resource Group %q) does not exist", serviceName, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on apiManagementPolicyClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMApiManagementPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).apiManagementPolicyClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_policy" {
			continue
		}

		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		if resp, err := client.Get(ctx, resourceGroup, serviceName); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on apiManagementPolicyClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMApiManagementPolicy_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku = {
    name     = "Developer"
    capacity = 1
  }
}

resource "azurerm_api_management_policy" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  api_management_name = "${azurerm_api_management.test.name}"
  xml_content         = "<policies></policies>"
}
`, rInt, location, rInt)
}

func testAccAzureRMApiManagementPolicy_findReplace(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku = {
    name     = "Developer"
    capacity = 1
  }
}

resource "azurerm_api_management_policy" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  api_management_name = "${azurerm_api_management.test.name}"
  xml_content         = "<policies><inbound><find-and-replace from='xyz' to='abc' /></inbound></policies>"
}
`, rInt, location, rInt)
}

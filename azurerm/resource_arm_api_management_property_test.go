package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAPIManagementProperty_basic(t *testing.T) {
	resourceName := "azurerm_api_management_property.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAPIManagementProperty_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAPIManagementPropertyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementPropertyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", fmt.Sprintf("TestProperty%d", ri)),
					resource.TestCheckResourceAttr(resourceName, "value", "Test Value"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "tag1"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "tag2"),
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

func TestAccAzureRMAPIManagementProperty_update(t *testing.T) {
	resourceName := "azurerm_api_management_property.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMAPIManagementProperty_basic(ri, testLocation())
	config2 := testAccAzureRMAPIManagementProperty_update(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAPIManagementPropertyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementPropertyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", fmt.Sprintf("TestProperty%d", ri)),
					resource.TestCheckResourceAttr(resourceName, "value", "Test Value"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "tag1"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "tag2"),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementPropertyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", fmt.Sprintf("TestProperty2%d", ri)),
					resource.TestCheckResourceAttr(resourceName, "value", "Test Value2"),
					resource.TestCheckResourceAttr(resourceName, "secret", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "tag3"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "tag4"),
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

func testCheckAzureRMAPIManagementPropertyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).apiManagement.PropertyClient
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_property" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		ctx := testAccProvider.Meta().(*ArmClient).StopContext
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

func testCheckAzureRMAPIManagementPropertyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		client := testAccProvider.Meta().(*ArmClient).apiManagement.PropertyClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, resourceGroup, serviceName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: API Management Property %q (Resource Group %q / API Management Service %q) does not exist", name, resourceGroup, serviceName)
			}
			return fmt.Errorf("Bad: Get on apiManagement.PropertyClient: %+v", err)
		}

		return nil
	}
}

/*

 */

func testAccAzureRMAPIManagementProperty_basic(rInt int, location string) string {
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

resource "azurerm_api_management_property" "test" {
  name                = "acctestAMProperty-%d"
  resource_group_name = "${azurerm_api_management.test.resource_group_name}"
  api_management_name = "${azurerm_api_management.test.name}"
  display_name        = "TestProperty%d"
  value               = "Test Value"
  tags                = ["tag1", "tag2"]
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMAPIManagementProperty_update(rInt int, location string) string {
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

resource "azurerm_api_management_property" "test" {
  name                = "acctestAMProperty-%d"
  resource_group_name = "${azurerm_api_management.test.resource_group_name}"
  api_management_name = "${azurerm_api_management.test.name}"
  display_name        = "TestProperty2%d"
  value               = "Test Value2"
  secret              = true
  tags                = ["tag3", "tag4"]
}
`, rInt, location, rInt, rInt, rInt)
}

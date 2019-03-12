package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMApiManagementApiVersionSet_basic(t *testing.T) {
	resourceName := "azurerm_api_management_api_version_set.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMApiManagementApiVersionSet_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiVersionSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiVersionSetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", fmt.Sprintf("TestDescription1")),
					resource.TestCheckResourceAttr(resourceName, "display_name", fmt.Sprintf("TestApiVersionSet1%d", ri)),
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

func TestAccAzureRMApiManagementApiVersionSet_update(t *testing.T) {
	resourceName := "azurerm_api_management_api_version_set.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMApiManagementApiVersionSet_basic(ri, testLocation())
	config2 := testAccAzureRMApiManagementApiVersionSet_update(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementApiVersionSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiVersionSetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", fmt.Sprintf("TestDescription1")),
					resource.TestCheckResourceAttr(resourceName, "display_name", fmt.Sprintf("TestApiVersionSet1%d", ri)),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementApiVersionSetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", fmt.Sprintf("TestDescription2")),
					resource.TestCheckResourceAttr(resourceName, "display_name", fmt.Sprintf("TestApiVersionSet2%d", ri)),
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

func testCheckAzureRMApiManagementApiVersionSetDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).apiManagementApiVersionSetClient
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_api_version_set" {
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

func testCheckAzureRMApiManagementApiVersionSetExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		client := testAccProvider.Meta().(*ArmClient).apiManagementApiVersionSetClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, resourceGroup, serviceName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Api Management Api Version Set %q (Resource Group %q / Api Management Service %q) does not exist", name, resourceGroup, serviceName)
			}
			return fmt.Errorf("Bad: Get on apiManagementApiVersionSetClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMApiManagementApiVersionSet_basic(rInt int, location string) string {
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

resource "azurerm_api_management_api_version_set" "test" {
  name                = "acctestAMAVS-%d"
  resource_group_name = "${azurerm_api_management.test.resource_group_name}"
  api_management_name = "${azurerm_api_management.test.name}"
  description         = "TestDescription1"
  display_name        = "TestApiVersionSet1%d"
  versioning_schema   = "Segment"
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMApiManagementApiVersionSet_update(rInt int, location string) string {
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

resource "azurerm_api_management_api_version_set" "test" {
  name                = "acctestAMAVS-%d"
  resource_group_name = "${azurerm_api_management.test.resource_group_name}"
  api_management_name = "${azurerm_api_management.test.name}"
  description         = "TestDescription2"
  display_name        = "TestApiVersionSet2%d"
  versioning_schema   = "Segment"
}
`, rInt, location, rInt, rInt, rInt)
}

package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSpringCloud_basic(t *testing.T) {
	resourceName := "azurerm_spring_cloud.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSpringCloudDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSpringCloud_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSpringCloudExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "service_id"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.env", "test"),
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

func TestAccAzureRMSpringCloud_update(t *testing.T) {
	resourceName := "azurerm_spring_cloud.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSpringCloudDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSpringCloud_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSpringCloudExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.env", "test"),
				),
			},
			{
				Config: testAccAzureRMSpringCloud_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSpringCloudExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.env", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags.version", "1"),
				),
			},
			{
				Config: testAccAzureRMSpringCloud_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSpringCloudExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.env", "test"),
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

func TestAccAzureRMSpringCloud_complete(t *testing.T) {
	resourceName := "azurerm_spring_cloud.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSpringCloudDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSpringCloud_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSpringCloudExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "service_id"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.env", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags.version", "1"),
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

func testCheckAzureRMSpringCloudExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Spring Cloud not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := testAccProvider.Meta().(*ArmClient).AppPlatform.ServicesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		if resp, err := client.Get(ctx, resourceGroup, name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Spring Cloud Service %q (Resource Group %q) does not exist", name, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on AppPlatform.ServicesClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMSpringCloudDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).AppPlatform.ServicesClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_spring_cloud" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on AppPlatform.ServicesClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMSpringCloud_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_spring_cloud" "test" {
  name                     = "acctestsc-%d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name

  tags = {
    env = "test"
  }
}

`, rInt, location, rInt)
}

func testAccAzureRMSpringCloud_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_spring_cloud" "test" {
  name                     = "acctestsc-%d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name

  tags = {
    env = "test"
    version = "1"
  }
}
`, rInt, location, rInt)
}

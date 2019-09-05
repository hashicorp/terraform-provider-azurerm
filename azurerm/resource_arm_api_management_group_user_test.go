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

func TestAccAzureRMAPIManagementGroupUser_basic(t *testing.T) {
	resourceName := "azurerm_api_management_group_user.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAPIManagementGroupUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAPIManagementGroupUser_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementGroupUserExists(resourceName),
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

func TestAccAzureRMAPIManagementGroupUser_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_api_management_group_user.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAPIManagementGroupUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAPIManagementGroupUser_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementGroupUserExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMAPIManagementGroupUser_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_api_management_group_user"),
			},
		},
	})
}

func testCheckAzureRMAPIManagementGroupUserDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).apiManagement.GroupUsersClient
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_group_user" {
			continue
		}

		userId := rs.Primary.Attributes["user_id"]
		groupName := rs.Primary.Attributes["group_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.CheckEntityExists(ctx, resourceGroup, serviceName, groupName, userId)
		if err != nil {
			if !utils.ResponseWasNotFound(resp) {
				return err
			}
		}

		return nil
	}
	return nil
}

func testCheckAzureRMAPIManagementGroupUserExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		userId := rs.Primary.Attributes["user_id"]
		groupName := rs.Primary.Attributes["group_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		client := testAccProvider.Meta().(*ArmClient).apiManagement.GroupUsersClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.CheckEntityExists(ctx, resourceGroup, serviceName, groupName, userId)
		if err != nil {
			if utils.ResponseWasNotFound(resp) {
				return fmt.Errorf("Bad: User %q / Group %q (API Management Service %q / Resource Group %q) does not exist", userId, groupName, serviceName, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on apiManagement.GroupUsersClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMAPIManagementGroupUser_basic(rInt int, location string) string {
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

resource "azurerm_api_management_group" "test" {
  name                = "acctestAMGroup-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  api_management_name = "${azurerm_api_management.test.name}"
  display_name        = "Test Group"
}

resource "azurerm_api_management_user" "test" {
  user_id             = "acctestuser%d"
  api_management_name = "${azurerm_api_management.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  first_name          = "Acceptance"
  last_name           = "Test"
  email               = "azure-acctest%d@example.com"
}

resource "azurerm_api_management_group_user" "test" {
  user_id             = "${azurerm_api_management_user.test.user_id}"
  group_name          = "${azurerm_api_management_group.test.name}"
  api_management_name = "${azurerm_api_management.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt, rInt, rInt, rInt)
}

func testAccAzureRMAPIManagementGroupUser_requiresImport(rInt int, location string) string {
	template := testAccAzureRMAPIManagementGroupUser_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_group_user" "import" {
  user_id             = "${azurerm_api_management_group_user.test.user_id}"
  group_name          = "${azurerm_api_management_group_user.test.group_name}"
  api_management_name = "${azurerm_api_management_group_user.test.api_management_name}"
  resource_group_name = "${azurerm_api_management_group_user.test.resource_group_name}"
}
`, template)
}

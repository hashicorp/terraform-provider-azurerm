package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMApiManagementGroupUser_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_group_user", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAPIManagementGroupUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementGroupUser_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementGroupUserExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagementGroupUser_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_group_user", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAPIManagementGroupUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementGroupUser_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementGroupUserExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMApiManagementGroupUser_requiresImport),
		},
	})
}

func testCheckAzureRMAPIManagementGroupUserDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.GroupUsersClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_group_user" {
			continue
		}

		userId := rs.Primary.Attributes["user_id"]
		groupName := rs.Primary.Attributes["group_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

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
		client := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.GroupUsersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		userId := rs.Primary.Attributes["user_id"]
		groupName := rs.Primary.Attributes["group_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

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

func testAccAzureRMApiManagementGroupUser_basic(data acceptance.TestData) string {
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

resource "azurerm_api_management_group" "test" {
  name                = "acctestAMGroup-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "Test Group"
}

resource "azurerm_api_management_user" "test" {
  user_id             = "acctestuser%d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  first_name          = "Acceptance"
  last_name           = "Test"
  email               = "azure-acctest%d@example.com"
}

resource "azurerm_api_management_group_user" "test" {
  user_id             = azurerm_api_management_user.test.user_id
  group_name          = azurerm_api_management_group.test.name
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApiManagementGroupUser_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementGroupUser_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_group_user" "import" {
  user_id             = azurerm_api_management_group_user.test.user_id
  group_name          = azurerm_api_management_group_user.test.group_name
  api_management_name = azurerm_api_management_group_user.test.api_management_name
  resource_group_name = azurerm_api_management_group_user.test.resource_group_name
}
`, template)
}

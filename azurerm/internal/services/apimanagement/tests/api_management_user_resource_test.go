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

func TestAccAzureRMApiManagementUser_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_user", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementUser_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementUserExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "first_name", "Acceptance"),
					resource.TestCheckResourceAttr(data.ResourceName, "last_name", "Test"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagementUser_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_user", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementUser_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementUserExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMApiManagementUser_requiresImport),
		},
	})
}

func TestAccAzureRMApiManagementUser_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_user", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementUser_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementUserExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "first_name", "Acceptance"),
					resource.TestCheckResourceAttr(data.ResourceName, "last_name", "Test"),
					resource.TestCheckResourceAttr(data.ResourceName, "state", "active"),
				),
			},
			{
				Config: testAccAzureRMApiManagementUser_updatedBlocked(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementUserExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "first_name", "Acceptance Updated"),
					resource.TestCheckResourceAttr(data.ResourceName, "last_name", "Test Updated"),
					resource.TestCheckResourceAttr(data.ResourceName, "state", "blocked"),
				),
			},
			{
				Config: testAccAzureRMApiManagementUser_updatedActive(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementUserExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "first_name", "Acceptance"),
					resource.TestCheckResourceAttr(data.ResourceName, "last_name", "Test"),
					resource.TestCheckResourceAttr(data.ResourceName, "state", "active"),
				),
			},
		},
	})
}

func TestAccAzureRMApiManagementUser_password(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_user", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementUser_password(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementUserExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "first_name", "Acceptance"),
					resource.TestCheckResourceAttr(data.ResourceName, "last_name", "Test"),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccAzureRMApiManagementUser_invite(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_user", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementUser_invited(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementUserExists(data.ResourceName),
				),
			},
			{
				ResourceName:      data.ResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// not returned
					"confirmation",
				},
			},
		},
	})
}

func TestAccAzureRMApiManagementUser_signup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_user", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementUser_signUp(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementUserExists(data.ResourceName),
				),
			},
			{
				ResourceName:      data.ResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// not returned
					"confirmation",
				},
			},
		},
	})
}

func TestAccAzureRMApiManagementUser_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_user", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementUser_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementUserExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "first_name", "Acceptance"),
					resource.TestCheckResourceAttr(data.ResourceName, "last_name", "Test"),
					resource.TestCheckResourceAttr(data.ResourceName, "note", "Used for testing in dimension C-137."),
				),
			},
			{
				ResourceName:      data.ResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// not returned
					"confirmation",
				},
			},
		},
	})
}

func testCheckAzureRMApiManagementUserDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.UsersClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_user" {
			continue
		}

		userId := rs.Primary.Attributes["user_id"]
		serviceName := rs.Primary.Attributes["api_management_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		resp, err := conn.Get(ctx, resourceGroup, serviceName, userId)
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

func testCheckAzureRMApiManagementUserExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.UsersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		userId := rs.Primary.Attributes["user_id"]
		serviceName := rs.Primary.Attributes["api_management_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, serviceName, userId)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: User %q (API Management Service %q / Resource Group %q) does not exist", userId, serviceName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on apiManagement.UsersClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMApiManagementUser_basic(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementUser_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_user" "test" {
  user_id             = "acctestuser%d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  first_name          = "Acceptance"
  last_name           = "Test"
  email               = "azure-acctest%d@example.com"
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApiManagementUser_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementUser_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_user" "import" {
  user_id             = azurerm_api_management_user.test.user_id
  api_management_name = azurerm_api_management_user.test.api_management_name
  resource_group_name = azurerm_api_management_user.test.resource_group_name
  first_name          = azurerm_api_management_user.test.first_name
  last_name           = azurerm_api_management_user.test.last_name
  email               = azurerm_api_management_user.test.email
  state               = azurerm_api_management_user.test.state
}
`, template)
}

func testAccAzureRMApiManagementUser_password(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementUser_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_user" "test" {
  user_id             = "acctestuser%d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  first_name          = "Acceptance"
  last_name           = "Test"
  email               = "azure-acctest%d@example.com"
  state               = "active"
  password            = "3991bb15-282d-4b9b-9de3-3d5fc89eb530"
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApiManagementUser_updatedActive(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementUser_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_user" "test" {
  user_id             = "acctestuser%d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  first_name          = "Acceptance"
  last_name           = "Test"
  email               = "azure-acctest%d@example.com"
  state               = "active"
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApiManagementUser_updatedBlocked(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementUser_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_user" "test" {
  user_id             = "acctestuser%d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  first_name          = "Acceptance Updated"
  last_name           = "Test Updated"
  email               = "azure-acctest%d@example.com"
  state               = "blocked"
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApiManagementUser_invited(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementUser_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_user" "test" {
  user_id             = "acctestuser%d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  first_name          = "Acceptance"
  last_name           = "Test User"
  email               = "azure-acctest%d@example.com"
  state               = "blocked"
  confirmation        = "invite"
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApiManagementUser_signUp(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementUser_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_user" "test" {
  user_id             = "acctestuser%d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  first_name          = "Acceptance"
  last_name           = "Test User"
  email               = "azure-acctest%d@example.com"
  state               = "blocked"
  confirmation        = "signup"
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApiManagementUser_complete(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementUser_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_user" "test" {
  user_id             = "acctestuser%d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  first_name          = "Acceptance"
  last_name           = "Test"
  email               = "azure-acctest%d@example.com"
  state               = "active"
  confirmation        = "signup"
  note                = "Used for testing in dimension C-137."
}
`, template, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMApiManagementUser_template(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

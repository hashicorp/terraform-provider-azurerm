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

func TestAccAzureRMApiManagementUser_basic(t *testing.T) {
	resourceName := "azurerm_api_management_user.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementUser_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementUserExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "first_name", "Acceptance"),
					resource.TestCheckResourceAttr(resourceName, "last_name", "Test"),
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

func TestAccAzureRMApiManagementUser_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_api_management_user.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementUser_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementUserExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMApiManagementUser_requiresImport(ri, location),
				ExpectError: acceptance.RequiresImportError("azurerm_api_management_user"),
			},
		},
	})
}

func TestAccAzureRMApiManagementUser_update(t *testing.T) {
	resourceName := "azurerm_api_management_user.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementUser_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementUserExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "first_name", "Acceptance"),
					resource.TestCheckResourceAttr(resourceName, "last_name", "Test"),
					resource.TestCheckResourceAttr(resourceName, "state", "active"),
				),
			},
			{
				Config: testAccAzureRMApiManagementUser_updatedBlocked(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementUserExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "first_name", "Acceptance Updated"),
					resource.TestCheckResourceAttr(resourceName, "last_name", "Test Updated"),
					resource.TestCheckResourceAttr(resourceName, "state", "blocked"),
				),
			},
			{
				Config: testAccAzureRMApiManagementUser_updatedActive(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementUserExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "first_name", "Acceptance"),
					resource.TestCheckResourceAttr(resourceName, "last_name", "Test"),
					resource.TestCheckResourceAttr(resourceName, "state", "active"),
				),
			},
		},
	})
}

func TestAccAzureRMApiManagementUser_password(t *testing.T) {
	resourceName := "azurerm_api_management_user.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementUser_password(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementUserExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "first_name", "Acceptance"),
					resource.TestCheckResourceAttr(resourceName, "last_name", "Test"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccAzureRMApiManagementUser_invite(t *testing.T) {
	resourceName := "azurerm_api_management_user.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementUser_invited(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementUserExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
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
	resourceName := "azurerm_api_management_user.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementUser_signUp(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementUserExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
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
	resourceName := "azurerm_api_management_user.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementUser_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementUserExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "first_name", "Acceptance"),
					resource.TestCheckResourceAttr(resourceName, "last_name", "Test"),
					resource.TestCheckResourceAttr(resourceName, "note", "Used for testing in dimension C-137."),
				),
			},
			{
				ResourceName:      resourceName,
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

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_user" {
			continue
		}

		userId := rs.Primary.Attributes["user_id"]
		serviceName := rs.Primary.Attributes["api_management_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
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
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		userId := rs.Primary.Attributes["user_id"]
		serviceName := rs.Primary.Attributes["api_management_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		conn := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.UsersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
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

func testAccAzureRMApiManagementUser_basic(rInt int, location string) string {
	template := testAccAzureRMApiManagementUser_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_user" "test" {
  user_id             = "acctestuser%d"
  api_management_name = "${azurerm_api_management.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  first_name          = "Acceptance"
  last_name           = "Test"
  email               = "azure-acctest%d@example.com"
}
`, template, rInt, rInt)
}

func testAccAzureRMApiManagementUser_requiresImport(rInt int, location string) string {
	template := testAccAzureRMApiManagementUser_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_user" "import" {
  user_id             = "${azurerm_api_management_user.test.user_id}"
  api_management_name = "${azurerm_api_management_user.test.api_management_name}"
  resource_group_name = "${azurerm_api_management_user.test.resource_group_name}"
  first_name          = "${azurerm_api_management_user.test.first_name}"
  last_name           = "${azurerm_api_management_user.test.last_name}"
  email               = "${azurerm_api_management_user.test.email}"
  state               = "${azurerm_api_management_user.test.state}"
}
`, template)
}

func testAccAzureRMApiManagementUser_password(rInt int, location string) string {
	template := testAccAzureRMApiManagementUser_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_user" "test" {
  user_id             = "acctestuser%d"
  api_management_name = "${azurerm_api_management.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  first_name          = "Acceptance"
  last_name           = "Test"
  email               = "azure-acctest%d@example.com"
  state               = "active"
  password            = "3991bb15-282d-4b9b-9de3-3d5fc89eb530"
}
`, template, rInt, rInt)
}

func testAccAzureRMApiManagementUser_updatedActive(rInt int, location string) string {
	template := testAccAzureRMApiManagementUser_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_user" "test" {
  user_id             = "acctestuser%d"
  api_management_name = "${azurerm_api_management.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  first_name          = "Acceptance"
  last_name           = "Test"
  email               = "azure-acctest%d@example.com"
  state               = "active"
}
`, template, rInt, rInt)
}

func testAccAzureRMApiManagementUser_updatedBlocked(rInt int, location string) string {
	template := testAccAzureRMApiManagementUser_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_user" "test" {
  user_id             = "acctestuser%d"
  api_management_name = "${azurerm_api_management.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  first_name          = "Acceptance Updated"
  last_name           = "Test Updated"
  email               = "azure-acctest%d@example.com"
  state               = "blocked"
}
`, template, rInt, rInt)
}

func testAccAzureRMApiManagementUser_invited(rInt int, location string) string {
	template := testAccAzureRMApiManagementUser_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_user" "test" {
  user_id             = "acctestuser%d"
  api_management_name = "${azurerm_api_management.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  first_name          = "Acceptance"
  last_name           = "Test User"
  email               = "azure-acctest%d@example.com"
  state               = "blocked"
  confirmation        = "invite"
}
`, template, rInt, rInt)
}

func testAccAzureRMApiManagementUser_signUp(rInt int, location string) string {
	template := testAccAzureRMApiManagementUser_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_user" "test" {
  user_id             = "acctestuser%d"
  api_management_name = "${azurerm_api_management.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  first_name          = "Acceptance"
  last_name           = "Test User"
  email               = "azure-acctest%d@example.com"
  state               = "blocked"
  confirmation        = "signup"
}
`, template, rInt, rInt)
}

func testAccAzureRMApiManagementUser_complete(rInt int, location string) string {
	template := testAccAzureRMApiManagementUser_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_user" "test" {
  user_id             = "acctestuser%d"
  api_management_name = "${azurerm_api_management.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  first_name          = "Acceptance"
  last_name           = "Test"
  email               = "azure-acctest%d@example.com"
  state               = "active"
  confirmation        = "signup"
  note                = "Used for testing in dimension C-137."
}
`, template, rInt, rInt)
}

func testAccAzureRMApiManagementUser_template(rInt int, location string) string {
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

  sku_name = "Developer_1"
}
`, rInt, location, rInt)
}

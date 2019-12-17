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

func TestAccAzureRMAPIManagementAuthorizationServer_basic(t *testing.T) {
	resourceName := "azurerm_api_management_authorization_server.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAPIManagementAuthorizationServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAPIManagementAuthorizationServer_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementAuthorizationServerExists(resourceName),
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

func TestAccAzureRMAPIManagementAuthorizationServer_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_api_management_authorization_server.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAPIManagementAuthorizationServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAPIManagementAuthorizationServer_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementAuthorizationServerExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMAPIManagementAuthorizationServer_requiresImport(ri, location),
				ExpectError: acceptance.RequiresImportError("azurerm_api_management_authorization_server"),
			},
		},
	})
}

func TestAccAzureRMAPIManagementAuthorizationServer_complete(t *testing.T) {
	resourceName := "azurerm_api_management_authorization_server.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAPIManagementAuthorizationServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAPIManagementAuthorizationServer_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementAuthorizationServerExists(resourceName),
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

func testCheckAzureRMAPIManagementAuthorizationServerDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.AuthorizationServersClient
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_authorization_server" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
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

func testCheckAzureRMAPIManagementAuthorizationServerExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		client := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.AuthorizationServersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		resp, err := client.Get(ctx, resourceGroup, serviceName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Authorization Server %q (API Management Service %q / Resource Group %q) does not exist", name, serviceName, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on apiManagementAuthorizationServersClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMAPIManagementAuthorizationServer_basic(rInt int, location string) string {
	template := testAccAzureRMAPIManagementAuthorizationServer_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_authorization_server" "test" {
  name                         = "acctestauthsrv-%d"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  api_management_name          = "${azurerm_api_management.test.name}"
  display_name                 = "Test Group"
  authorization_endpoint       = "https://azacctest.hashicorptest.com/client/authorize"
  client_id                    = "42424242-4242-4242-4242-424242424242"
  client_registration_endpoint = "https://azacctest.hashicorptest.com/client/register"

  grant_types = [
    "implicit",
  ]

  authorization_methods = [
    "GET",
  ]
}
`, template, rInt)
}

func testAccAzureRMAPIManagementAuthorizationServer_requiresImport(rInt int, location string) string {
	template := testAccAzureRMAPIManagementAuthorizationServer_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_authorization_server" "import" {
  name                         = "${azurerm_api_management_authorization_server.test.name}"
  resource_group_name          = "${azurerm_api_management_authorization_server.test.resource_group_name}"
  api_management_name          = "${azurerm_api_management_authorization_server.test.api_management_name}"
  display_name                 = "${azurerm_api_management_authorization_server.test.display_name}"
  authorization_endpoint       = "${azurerm_api_management_authorization_server.test.authorization_endpoint}"
  client_id                    = "${azurerm_api_management_authorization_server.test.client_id}"
  client_registration_endpoint = "${azurerm_api_management_authorization_server.test.client_registration_endpoint}"
  grant_types                  = "${azurerm_api_management_authorization_server.test.grant_types}"

  authorization_methods = [
    "GET",
  ]
}
`, template)
}

func testAccAzureRMAPIManagementAuthorizationServer_complete(rInt int, location string) string {
	template := testAccAzureRMAPIManagementAuthorizationServer_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_authorization_server" "test" {
  name                         = "acctestauthsrv-%d"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  api_management_name          = "${azurerm_api_management.test.name}"
  display_name                 = "Test Group"
  authorization_endpoint       = "https://azacctest.hashicorptest.com/client/authorize"
  client_id                    = "42424242-4242-4242-4242-424242424242"
  client_registration_endpoint = "https://azacctest.hashicorptest.com/client/register"

  grant_types = [
    "authorizationCode",
  ]

  authorization_methods = [
    "GET",
    "POST",
  ]

  bearer_token_sending_methods = [
    "authorizationHeader",
  ]

  client_secret           = "n1n3-m0re-s3a5on5-m0r1y"
  default_scope           = "read write"
  token_endpoint          = "https://azacctest.hashicorptest.com/client/token"
  resource_owner_username = "rick"
  resource_owner_password = "C-193P"
  support_state           = true
}
`, template, rInt)
}

func testAccAzureRMAPIManagementAuthorizationServer_template(rInt int, location string) string {
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
`, rInt, location, rInt)
}

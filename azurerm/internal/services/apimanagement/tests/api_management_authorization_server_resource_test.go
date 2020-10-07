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

func TestAccAzureRMApiManagementAuthorizationServer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_authorization_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAPIManagementAuthorizationServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementAuthorizationServer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementAuthorizationServerExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagementAuthorizationServer_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_authorization_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAPIManagementAuthorizationServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementAuthorizationServer_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementAuthorizationServerExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMApiManagementAuthorizationServer_requiresImport),
		},
	})
}

func TestAccAzureRMApiManagementAuthorizationServer_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_authorization_server", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAPIManagementAuthorizationServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementAuthorizationServer_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementAuthorizationServerExists(data.ResourceName),
				),
			},
			data.ImportStep("client_secret"),
		},
	})
}

func testCheckAzureRMAPIManagementAuthorizationServerDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.AuthorizationServersClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_authorization_server" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

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
		client := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.AuthorizationServersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

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

func testAccAzureRMApiManagementAuthorizationServer_basic(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementAuthorizationServer_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_authorization_server" "test" {
  name                         = "acctestauthsrv-%d"
  resource_group_name          = azurerm_resource_group.test.name
  api_management_name          = azurerm_api_management.test.name
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
`, template, data.RandomInteger)
}

func testAccAzureRMApiManagementAuthorizationServer_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementAuthorizationServer_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_authorization_server" "import" {
  name                         = azurerm_api_management_authorization_server.test.name
  resource_group_name          = azurerm_api_management_authorization_server.test.resource_group_name
  api_management_name          = azurerm_api_management_authorization_server.test.api_management_name
  display_name                 = azurerm_api_management_authorization_server.test.display_name
  authorization_endpoint       = azurerm_api_management_authorization_server.test.authorization_endpoint
  client_id                    = azurerm_api_management_authorization_server.test.client_id
  client_registration_endpoint = azurerm_api_management_authorization_server.test.client_registration_endpoint
  grant_types                  = azurerm_api_management_authorization_server.test.grant_types

  authorization_methods = [
    "GET",
  ]
}
`, template)
}

func testAccAzureRMApiManagementAuthorizationServer_complete(data acceptance.TestData) string {
	template := testAccAzureRMApiManagementAuthorizationServer_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_authorization_server" "test" {
  name                         = "acctestauthsrv-%d"
  resource_group_name          = azurerm_resource_group.test.name
  api_management_name          = azurerm_api_management.test.name
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
`, template, data.RandomInteger)
}

func testAccAzureRMApiManagementAuthorizationServer_template(data acceptance.TestData) string {
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
  sku_name            = "Developer_1"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

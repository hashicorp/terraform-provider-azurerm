package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ApiManagementAuthorizationServerResource struct {
}

func TestAccApiManagementAuthorizationServer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_authorization_server", "test")
	r := ApiManagementAuthorizationServerResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementAuthorizationServer_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_authorization_server", "test")
	r := ApiManagementAuthorizationServerResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccApiManagementAuthorizationServer_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_authorization_server", "test")
	r := ApiManagementAuthorizationServerResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("client_secret"),
	})
}

func (t ApiManagementAuthorizationServerResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}

	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	name := id.Path["authorizationServers"]

	resp, err := clients.ApiManagement.AuthorizationServersClient.Get(ctx, resourceGroup, serviceName, name)
	if err != nil {
		return nil, fmt.Errorf("reading ApiManagement Authorization Server (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r ApiManagementAuthorizationServerResource) basic(data acceptance.TestData) string {
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
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementAuthorizationServerResource) requiresImport(data acceptance.TestData) string {
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
`, r.basic(data))
}

func (r ApiManagementAuthorizationServerResource) complete(data acceptance.TestData) string {
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
`, r.template(data), data.RandomInteger)
}

func (ApiManagementAuthorizationServerResource) template(data acceptance.TestData) string {
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

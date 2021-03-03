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

type ApiManagementOpenIDConnectProviderResource struct {
}

func TestAccApiManagementOpenIDConnectProvider_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_openid_connect_provider", "test")
	r := ApiManagementOpenIDConnectProviderResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("client_secret"),
	})
}

func TestAccApiManagementOpenIDConnectProvider_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_openid_connect_provider", "test")
	r := ApiManagementOpenIDConnectProviderResource{}

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

func TestAccApiManagementOpenIDConnectProvider_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_openid_connect_provider", "test")
	r := ApiManagementOpenIDConnectProviderResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("client_secret"),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("client_secret"),
	})
}

func (ApiManagementOpenIDConnectProviderResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	name := id.Path["openidConnectProviders"]

	resp, err := clients.ApiManagement.OpenIdConnectClient.Get(ctx, resourceGroup, serviceName, name)
	if err != nil {
		return nil, fmt.Errorf("reading ApiManagement Open ID Connect (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r ApiManagementOpenIDConnectProviderResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_openid_connect_provider" "test" {
  name                = "acctest-%d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  client_id           = "00001111-2222-3333-%d"
  client_secret       = "%d-cwdavsxbacsaxZX-%d"
  display_name        = "Initial Name"
  metadata_endpoint   = "https://azacctest.hashicorptest.com/example/foo"
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementOpenIDConnectProviderResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_openid_connect_provider" "import" {
  name                = azurerm_api_management_openid_connect_provider.test.name
  api_management_name = azurerm_api_management_openid_connect_provider.test.api_management_name
  resource_group_name = azurerm_api_management_openid_connect_provider.test.resource_group_name
  client_id           = azurerm_api_management_openid_connect_provider.test.client_id
  client_secret       = azurerm_api_management_openid_connect_provider.test.client_secret
  display_name        = azurerm_api_management_openid_connect_provider.test.display_name
  metadata_endpoint   = azurerm_api_management_openid_connect_provider.test.metadata_endpoint
}
`, r.basic(data))
}

func (r ApiManagementOpenIDConnectProviderResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_openid_connect_provider" "test" {
  name                = "acctest-%d"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  client_id           = "00001111-3333-2222-%d"
  client_secret       = "%d-423egvwdcsjx-%d"
  display_name        = "Updated Name"
  description         = "Example description"
  metadata_endpoint   = "https://azacctest.hashicorptest.com/example/updated"
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (ApiManagementOpenIDConnectProviderResource) template(data acceptance.TestData) string {
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

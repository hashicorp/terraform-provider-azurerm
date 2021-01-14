package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2019-12-01/apimanagement"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ApiManagementIdentityProviderAADResource struct {
}

func TestAccApiManagementIdentityProviderAAD_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_identity_provider_aad", "test")
	r := ApiManagementIdentityProviderAADResource{}

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

func TestAccApiManagementIdentityProviderAAD_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_identity_provider_aad", "test")
	r := ApiManagementIdentityProviderAADResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("client_id").HasValue("00000000-0000-0000-0000-000000000000"),
				check.That(data.ResourceName).Key("client_secret").HasValue("00000000000000000000000000000000"),
				check.That(data.ResourceName).Key("allowed_tenants.#").HasValue("1"),
				check.That(data.ResourceName).Key("allowed_tenants.0").HasValue(data.Client().TenantID),
			),
		},
		{
			Config: r.update(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("client_id").HasValue("11111111-1111-1111-1111-111111111111"),
				check.That(data.ResourceName).Key("client_secret").HasValue("11111111111111111111111111111111"),
				check.That(data.ResourceName).Key("allowed_tenants.#").HasValue("2"),
				check.That(data.ResourceName).Key("allowed_tenants.0").HasValue(data.Client().TenantID),
				check.That(data.ResourceName).Key("allowed_tenants.1").HasValue(data.Client().TenantID),
			),
		},
		data.ImportStep("client_secret"),
	})
}

func TestAccApiManagementIdentityProviderAAD_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_identity_provider_aad", "test")
	r := ApiManagementIdentityProviderAADResource{}

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

func (t ApiManagementIdentityProviderAADResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	identityProviderName := id.Path["identityProviders"]

	resp, err := clients.ApiManagement.IdentityProviderClient.Get(ctx, resourceGroup, serviceName, apimanagement.IdentityProviderType(identityProviderName))
	if err != nil {
		return nil, fmt.Errorf("reading ApiManagement Identity Provider AAD (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (ApiManagementIdentityProviderAADResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-api-%d"
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

resource "azurerm_api_management_identity_provider_aad" "test" {
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  client_id           = "00000000-0000-0000-0000-000000000000"
  client_secret       = "00000000000000000000000000000000"
  signin_tenant       = "00000000-0000-0000-0000-000000000000"
  allowed_tenants     = ["%s"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Client().TenantID)
}

func (ApiManagementIdentityProviderAADResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-api-%d"
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

resource "azurerm_api_management_identity_provider_aad" "test" {
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  client_id           = "11111111-1111-1111-1111-111111111111"
  client_secret       = "11111111111111111111111111111111"
  allowed_tenants     = ["%s", "%s"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Client().TenantID, data.Client().TenantID)
}

func (r ApiManagementIdentityProviderAADResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_identity_provider_aad" "import" {
  resource_group_name = azurerm_api_management_identity_provider_aad.test.resource_group_name
  api_management_name = azurerm_api_management_identity_provider_aad.test.api_management_name
  client_id           = azurerm_api_management_identity_provider_aad.test.client_id
  client_secret       = azurerm_api_management_identity_provider_aad.test.client_secret
  allowed_tenants     = azurerm_api_management_identity_provider_aad.test.allowed_tenants
}
`, r.basic(data))
}

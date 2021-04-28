package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2019-12-01/apimanagement"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ApiManagementIdentityProviderMicrosoftResource struct {
}

func TestAccApiManagementIdentityProviderMicrosoft_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_identity_provider_microsoft", "test")
	r := ApiManagementIdentityProviderMicrosoftResource{}

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

func TestAccApiManagementIdentityProviderMicrosoft_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_identity_provider_microsoft", "test")
	r := ApiManagementIdentityProviderMicrosoftResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("client_id").HasValue("00000000-0000-0000-0000-000000000000"),
			),
		},
		data.ImportStep("client_secret"),
		{
			Config: r.update(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("client_id").HasValue("11111111-1111-1111-1111-111111111111"),
			),
		},
		data.ImportStep("client_secret"),
	})
}

func TestAccApiManagementIdentityProviderMicrosoft_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_identity_provider_microsoft", "test")
	r := ApiManagementIdentityProviderMicrosoftResource{}

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

func (ApiManagementIdentityProviderMicrosoftResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.IdentityProviderID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.IdentityProviderClient.Get(ctx, id.ResourceGroup, id.ServiceName, apimanagement.IdentityProviderType(id.Name))
	if err != nil {
		return nil, fmt.Errorf("reading ApiManagement Identity Provider Microsoft (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (ApiManagementIdentityProviderMicrosoftResource) basic(data acceptance.TestData) string {
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

resource "azurerm_api_management_identity_provider_microsoft" "test" {
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  client_id           = "00000000-0000-0000-0000-000000000000"
  client_secret       = "00000000000000000000000000000000"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ApiManagementIdentityProviderMicrosoftResource) update(data acceptance.TestData) string {
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

resource "azurerm_api_management_identity_provider_microsoft" "test" {
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  client_id           = "11111111-1111-1111-1111-111111111111"
  client_secret       = "11111111111111111111111111111111"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ApiManagementIdentityProviderMicrosoftResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_identity_provider_microsoft" "import" {
  resource_group_name = azurerm_api_management_identity_provider_microsoft.test.resource_group_name
  api_management_name = azurerm_api_management_identity_provider_microsoft.test.api_management_name
  client_id           = azurerm_api_management_identity_provider_microsoft.test.client_id
  client_secret       = azurerm_api_management_identity_provider_microsoft.test.client_secret
}
`, r.basic(data))
}

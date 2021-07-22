package apimanagement_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ApiManagementGatewayAPIResource struct {
}

func TestAccApiManagementGatewayApi_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_gateway_api", "test")
	r := ApiManagementGatewayAPIResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementGatewayApi_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_gateway_api", "test")
	r := ApiManagementGatewayAPIResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (ApiManagementGatewayAPIResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.GatewayApiID(state.ID)
	if err != nil {
		return nil, err
	}
	if resp, err := clients.ApiManagement.GatewayApisClient.GetEntityTag(ctx, id.ResourceGroup, id.ServiceName, id.ID(), id.ApiName); err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil, fmt.Errorf("reading ApiManagement Gateway (%s): %+v", id, err)
		}

		if !utils.ResponseWasStatusCode(resp, http.StatusNoContent) {
			return nil, fmt.Errorf("reading ApiManagement Gateway (%s): %+v", id, err)
		}
	}

	return utils.Bool(true), nil
}

func (ApiManagementGatewayAPIResource) basic(data acceptance.TestData) string {
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

resource "azurerm_api_management_gateway" "test" {
  name              = "acctestAMGateway-%d"
  api_management_id = azurerm_api_management.test.id
  description       = "this is a test gateway"

  location_data {
    name     = "old world"
    city     = "test city"
    district = "test district"
    region   = "test region"
  }
}

resource "azurerm_api_management_api" "test" {
  name                = "acctestapi-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "api1"
  path                = "api1"
  protocols           = ["https"]
  revision            = "1"
}

resource "azurerm_api_management_gateway_api" "test" {
  gateway_id = azurerm_api_management_gateway.test.id
  api_id     = azurerm_api_management_api.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementGatewayAPIResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_gateway_api" "import" {
  gateway_id = azurerm_api_management_gateway_api.test.id
  api_id     = azurerm_api_management_gateway_api.test.id
}
`, r.basic(data))
}

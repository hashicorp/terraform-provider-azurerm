package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
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
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	gatewayId := id.Path["gateways"]
	apiName := id.Path["apis"]

	if _, err = clients.ApiManagement.GatewayApisClient.GetEntityTag(ctx, resourceGroup, serviceName, gatewayId, apiName); err != nil {
		return nil, fmt.Errorf("reading ApiManagement Gateway (%s): %+v", id, err)
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
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  gateway_id          = "TestGateway"
  location            = "old world updated"
  description         = "this is a test gateway updated"
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
  gateway_id          = azurerm_api_management_gateway.test.gateway_id
  api_name            = azurerm_api_management_api.test.name
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementGatewayAPIResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_gateway_api" "import" {
  api_name            = azurerm_api_management_gateway_api.test.api_name
  gateway_id          = azurerm_api_management_gateway_api.test.gateway_id
  api_management_name = azurerm_api_management_gateway_api.test.api_management_name
  resource_group_name = azurerm_api_management_gateway_api.test.resource_group_name
}
`, r.basic(data))
}

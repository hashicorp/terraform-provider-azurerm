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

type ApiManagementGatewayResource struct {
}

func TestAccApiManagementGateway_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_gateway", "test")
	r := ApiManagementGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("location").HasValue("old world"),
				check.That(data.ResourceName).Key("description").HasValue("this is a test gateway"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementGateway_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_gateway", "test")
	r := ApiManagementGatewayResource{}

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

func TestAccApiManagementGateway_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_gateway", "test")
	r := ApiManagementGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("location").HasValue("old world"),
				check.That(data.ResourceName).Key("description").HasValue("this is a test gateway"),
			),
		},
		data.ImportStep(),
		{
			Config: r.updated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("location").HasValue("old world updated"),
				check.That(data.ResourceName).Key("description").HasValue("this is a test gateway updated"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("location").HasValue("old world"),
				check.That(data.ResourceName).Key("description").HasValue("this is a test gateway"),
			),
		},
	})
}

func (ApiManagementGatewayResource) updated(data acceptance.TestData) string {
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
	api_management_name   = azurerm_api_management.test.name
	resource_group_name   = azurerm_resource_group.test.name
	gateway_id          	= "TestGateway"
	location 				= "old world updated"
	description     		= "this is a test gateway updated"
  }
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ApiManagementGatewayResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	gatewayId := id.Path["gateways"]

	resp, err := clients.ApiManagement.GatewayClient.Get(ctx, resourceGroup, serviceName, gatewayId)
	if err != nil {
		return nil, fmt.Errorf("reading ApiManagement Gateway (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (ApiManagementGatewayResource) basic(data acceptance.TestData) string {
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
	api_management_name   = azurerm_api_management.test.name
	resource_group_name   = azurerm_resource_group.test.name
	gateway_id          	= "TestGateway"
	location 				= "old world"
	description     		= "this is a test gateway"
  }
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ApiManagementGatewayResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_gateway" "import" {
	api_management_name   = azurerm_api_management_gateway.test.api_management_name
	resource_group_name   = azurerm_api_management_gateway.test.resource_group_name
	name          		  = azurerm_api_management_gateway.test.name
	location 			  = azurerm_api_management_gateway.test.location
	description     	  = azurerm_api_management_gateway.test.description
}
`, r.basic(data))
}

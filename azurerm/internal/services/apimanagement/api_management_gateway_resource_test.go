package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
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
				check.That(data.ResourceName).Key("description").HasValue(""),
				check.That(data.ResourceName).Key("location_data.0.name").HasValue("test"),
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

func TestAccApiManagementGateway_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_gateway", "test")
	r := ApiManagementGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, "test description", "test location"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("description").HasValue("test description"),
				check.That(data.ResourceName).Key("location_data.0.name").HasValue("test location"),
				check.That(data.ResourceName).Key("location_data.0.city").HasValue("test city"),
				check.That(data.ResourceName).Key("location_data.0.district").HasValue("test district"),
				check.That(data.ResourceName).Key("location_data.0.region").HasValue("test region"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementGateway_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_gateway", "test")
	r := ApiManagementGatewayResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, "original description", "original location"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("description").HasValue("original description"),
				check.That(data.ResourceName).Key("location_data.#").HasValue("1"),
				check.That(data.ResourceName).Key("location_data.0.name").HasValue("original location"),
				check.That(data.ResourceName).Key("location_data.0.city").HasValue("test city"),
				check.That(data.ResourceName).Key("location_data.0.district").HasValue("test district"),
				check.That(data.ResourceName).Key("location_data.0.region").HasValue("test region"),
			),
		},
		{
			Config: r.update(data, "updated description", "updated location"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("description").HasValue("updated description"),
				check.That(data.ResourceName).Key("location_data.0.name").HasValue("updated location"),
				check.That(data.ResourceName).Key("location_data.0.city").HasValue(""),
				check.That(data.ResourceName).Key("location_data.0.district").HasValue(""),
				check.That(data.ResourceName).Key("location_data.0.region").HasValue(""),
			),
		},
		{
			Config: r.complete(data, "original description", "original location"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("description").HasValue("original description"),
				check.That(data.ResourceName).Key("location_data.0.name").HasValue("original location"),
				check.That(data.ResourceName).Key("location_data.0.city").HasValue("test city"),
				check.That(data.ResourceName).Key("location_data.0.district").HasValue("test district"),
				check.That(data.ResourceName).Key("location_data.0.region").HasValue("test region"),
			),
		},
	})
}

func (ApiManagementGatewayResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.GatewayID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.GatewayClient.Get(ctx, id.ResourceGroup, id.ServiceName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading ApiManagementGateway (%s): %+v", id, err)
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
  name              = "acctestAMGateway-%d"
  api_management_id = azurerm_api_management.test.id

  location_data {
    name = "test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementGatewayResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_gateway" "import" {
  name              = azurerm_api_management_gateway.test.name
  api_management_id = azurerm_api_management_gateway.test.api_management_id

  location_data {
    name = "test"
  }
}
`, r.basic(data))
}

func (ApiManagementGatewayResource) update(data acceptance.TestData, description string, locationName string) string {
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
  description       = "%s"

  location_data {
    name = "%s"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, description, locationName)
}

func (ApiManagementGatewayResource) complete(data acceptance.TestData, description string, locationName string) string {
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
  description       = "%s"

  location_data {
    name     = "%s"
    city     = "test city"
    district = "test district"
    region   = "test region"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, description, locationName)
}

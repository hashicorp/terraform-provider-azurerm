package apimanagement_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ApiManagementGatewayHostnameConfigurationDataSource struct {
}

func TestAccDataSourceApiManagementGatewayHostnameConfiguration_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_api_management_gateway_hostname_configuration", "test")
	r := ApiManagementGatewayHostnameConfigurationDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("hostname").HasValue("example.apim.net"),
				check.That(data.ResourceName).Key("name").HasValue("example-hostname"),
			),
		},
	})
}

func (ApiManagementGatewayHostnameConfigurationDataSource) basic(data acceptance.TestData) string {
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

resource "azurerm_api_management_gateway" "test" {
  name              = "acctestAMGateway-%d"
  api_management_id = azurerm_api_management.test.id

  location_data {
    name = "test"
  }
}

resource "azurerm_api_management_certificate" "test" {
  name                = "example-cert"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  data                = filebase64("testdata/keyvaultcert.pfx")
}

resource "azurerm_api_management_gateway_hostname_configuration" "test" {
  name = "example-hostname"
  api_management_gateway_id = data.azurerm_api_management_gateway.test.id
  hostname = "example.apim.net"
  certificate_id = azurerm_api_management_certificate.test.id
}

data "azurerm_api_management_gateway_hostname_configuration" "main" {
  depends_on = [
    azurerm_api_management_gateway_hostname_configuration.test,
  ]
  name = "example-hostname"
  api_management_gateway_id = azurerm_api_management_gateway.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

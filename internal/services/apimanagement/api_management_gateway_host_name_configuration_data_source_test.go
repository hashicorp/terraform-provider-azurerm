// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	data := acceptance.BuildTestData(t, "data.azurerm_api_management_gateway_host_name_configuration", "test")
	r := ApiManagementGatewayHostnameConfigurationDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("gateway_name").HasValue(fmt.Sprintf("acctestAMGateway-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("host_name").HasValue(fmt.Sprintf("host-name-%s", data.RandomString)),
				check.That(data.ResourceName).Key("request_client_certificate_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("http2_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("tls10_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("tls11_enabled").HasValue("true"),
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
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Developer_1"
}

resource "azurerm_api_management_gateway" "test" {
  name              = "acctestAMGateway-%[1]d"
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

resource "azurerm_api_management_gateway_host_name_configuration" "test" {
  name              = "acctestAMGatewayHostNameConfiguration-%[1]d"
  api_management_id = azurerm_api_management.test.id
  gateway_name      = azurerm_api_management_gateway.test.name

  certificate_id                     = azurerm_api_management_certificate.test.id
  host_name                          = "host-name-%[3]s"
  request_client_certificate_enabled = false
  http2_enabled                      = false
  tls10_enabled                      = false
  tls11_enabled                      = true
}

data "azurerm_api_management_gateway_host_name_configuration" "test" {
  name              = azurerm_api_management_gateway_host_name_configuration.test.name
  api_management_id = azurerm_api_management.test.id
  gateway_name      = azurerm_api_management_gateway.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

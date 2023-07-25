// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type AppGatewayDataSource struct{}

func TestAccDataSourceAppGateway_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_application_gateway", "test")
	r := AppGatewayDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
			),
		},
	})
}

func TestAccDataSourceAppGateway_userAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_application_gateway", "test")
	r := AppGatewayDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.userAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("1"),
			),
		},
	})
}
func TestAccDataSourceAppGateway_backendAddressPool(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_application_gateway", "test")
	r := AppGatewayDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.backendAddressPool(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("backend_address_pool.0.id").Exists(),
				check.That(data.ResourceName).Key("backend_address_pool.0.name").Exists(),
			),
		},
	})
}

func (AppGatewayDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_application_gateway" "test" {
  resource_group_name = azurerm_application_gateway.test.resource_group_name
  name                = azurerm_application_gateway.test.name
}
`, ApplicationGatewayResource{}.basic(data))
}

func (AppGatewayDataSource) userAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_application_gateway" "test" {
  resource_group_name = azurerm_application_gateway.test.resource_group_name
  name                = azurerm_application_gateway.test.name
}
`, ApplicationGatewayResource{}.UserDefinedIdentity(data))
}
func (AppGatewayDataSource) backendAddressPool(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_application_gateway" "test" {
  resource_group_name = azurerm_application_gateway.test.resource_group_name
  name                = azurerm_application_gateway.test.name
}
`, ApplicationGatewayResource{}.backendAddressPoolEmptyIpList(data))
}

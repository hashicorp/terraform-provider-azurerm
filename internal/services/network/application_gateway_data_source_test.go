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

func TestAccDataSourceAppGateway_sslProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_application_gateway", "test")
	r := AppGatewayDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.sslProfile(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("sku.0.capacity").Exists(),
				check.That(data.ResourceName).Key("gateway_ip_configuration.0.subnet_id").Exists(),
				check.That(data.ResourceName).Key("frontend_port.0.port").Exists(),
				check.That(data.ResourceName).Key("frontend_ip_configuration.0.public_ip_address_id").Exists(),
				check.That(data.ResourceName).Key("backend_address_pool.0.name").Exists(),
				check.That(data.ResourceName).Key("backend_http_settings.0.cookie_based_affinity").Exists(),
				check.That(data.ResourceName).Key("http_listener.0.frontend_port_name").Exists(),
				check.That(data.ResourceName).Key("http_listener.0.protocol").Exists(),
				check.That(data.ResourceName).Key("request_routing_rule.0.priority").Exists(),
				check.That(data.ResourceName).Key("request_routing_rule.0.rule_type").Exists(),
				check.That(data.ResourceName).Key("ssl_profile.0.verify_client_certificate_revocation").Exists(),
				check.That(data.ResourceName).Key("ssl_profile.0.ssl_policy.0.policy_type").Exists(),
				check.That(data.ResourceName).Key("ssl_profile.0.ssl_policy.0.policy_name").Exists(),
				check.That(data.ResourceName).Key("ssl_certificate.0.name").Exists(),
			),
		},
	})
}

func TestAccDataSourceAppGateway_sslProfileWithClientCertificateVerification(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_application_gateway", "test")
	r := AppGatewayDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.sslProfileWithClientCertificateVerification(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("trusted_client_certificate.0.name").IsNotEmpty(),
				check.That(data.ResourceName).Key("trusted_client_certificate.0.data").IsNotEmpty(),
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

func (AppGatewayDataSource) sslProfile(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_application_gateway" "test" {
  resource_group_name = azurerm_application_gateway.test.resource_group_name
  name                = azurerm_application_gateway.test.name
}
`, ApplicationGatewayResource{}.sslProfileUpdateOne(data))
}

func (AppGatewayDataSource) sslProfileWithClientCertificateVerification(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_application_gateway" "test" {
  resource_group_name = azurerm_application_gateway.test.resource_group_name
  name                = azurerm_application_gateway.test.name
}
`, ApplicationGatewayResource{}.sslProfileWithClientCertificateVerification(data))
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

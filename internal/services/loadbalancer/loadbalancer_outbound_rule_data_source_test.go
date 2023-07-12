// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loadbalancer_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccAzureRMDataSourceLoadBalancerOutboundRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_lb_outbound_rule", "test")
	r := LoadBalancerOutboundRule{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basicDataSource(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("frontend_ip_configuration.0.name").Exists(),
				check.That(data.ResourceName).Key("protocol").Exists(),
			),
		},
	})
}

func TestAccAzureRMDataSourceLoadBalancerOutboundRule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_lb_outbound_rule", "test")
	r := LoadBalancerOutboundRule{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.completeDataSource(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("frontend_ip_configuration.0.name").Exists(),
				check.That(data.ResourceName).Key("protocol").Exists(),
				check.That(data.ResourceName).Key("backend_address_pool_id").Exists(),
				check.That(data.ResourceName).Key("idle_timeout_in_minutes").Exists(),
				check.That(data.ResourceName).Key("tcp_reset_enabled").Exists(),
			),
		},
	})
}

func (r LoadBalancerOutboundRule) basicDataSource(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_lb_outbound_rule" "test" {
  name            = azurerm_lb_outbound_rule.test.name
  loadbalancer_id = azurerm_lb_outbound_rule.test.loadbalancer_id
}
`, template)
}

func (r LoadBalancerOutboundRule) completeDataSource(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test1" {
  name                = "test-ip-1-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_public_ip" "test2" {
  name                = "test-ip-2-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_public_ip" "test3" {
  name                = "test-ip-3-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_lb" "test" {
  name                = "arm-test-loadbalancer-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "fe1-%d"
    public_ip_address_id = azurerm_public_ip.test1.id
  }

  frontend_ip_configuration {
    name                 = "fe2-%d"
    public_ip_address_id = azurerm_public_ip.test2.id
  }

  frontend_ip_configuration {
    name                 = "fe3-%d"
    public_ip_address_id = azurerm_public_ip.test3.id
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  loadbalancer_id = azurerm_lb.test.id
  name            = "be-%d"
}

resource "azurerm_lb_outbound_rule" "test" {
  loadbalancer_id         = azurerm_lb.test.id
  name                    = "OutboundRule1-%d"
  protocol                = "All"
  backend_address_pool_id = azurerm_lb_backend_address_pool.test.id

  frontend_ip_configuration {
    name = "fe1-%d"
  }
}

resource "azurerm_lb_outbound_rule" "test2" {
  loadbalancer_id         = azurerm_lb.test.id
  name                    = "OutboundRule2-%d"
  protocol                = "Tcp"
  backend_address_pool_id = azurerm_lb_backend_address_pool.test.id
  enable_tcp_reset        = true
  idle_timeout_in_minutes = 5

  frontend_ip_configuration {
    name = "fe2-%d"
  }
}

resource "azurerm_lb_outbound_rule" "test3" {
  loadbalancer_id         = azurerm_lb.test.id
  name                    = "OutboundRule3-%d"
  protocol                = "Udp"
  backend_address_pool_id = azurerm_lb_backend_address_pool.test.id

  frontend_ip_configuration {
    name = "fe3-%d"
  }
}

data "azurerm_lb_outbound_rule" "test" {
  name            = azurerm_lb_outbound_rule.test2.name
  loadbalancer_id = azurerm_lb_outbound_rule.test.loadbalancer_id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

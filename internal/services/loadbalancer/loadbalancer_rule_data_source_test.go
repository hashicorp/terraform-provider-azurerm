// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loadbalancer_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
)

func TestAccAzureRMDataSourceLoadBalancerRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_lb_rule", "test")
	r := LoadBalancerRule{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basicDataSource(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("frontend_ip_configuration_name").Exists(),
				check.That(data.ResourceName).Key("protocol").Exists(),
				check.That(data.ResourceName).Key("frontend_port").Exists(),
				check.That(data.ResourceName).Key("backend_port").Exists(),
			),
		},
	})
}

func TestAccAzureRMDataSourceLoadBalancerRule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_lb_rule", "test")
	r := LoadBalancerRule{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.completeDataSource(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("frontend_ip_configuration_name").Exists(),
				check.That(data.ResourceName).Key("protocol").Exists(),
				check.That(data.ResourceName).Key("frontend_port").Exists(),
				check.That(data.ResourceName).Key("backend_port").Exists(),
				check.That(data.ResourceName).Key("backend_address_pool_id").Exists(),
				check.That(data.ResourceName).Key("probe_id").Exists(),
				check.That(data.ResourceName).Key("floating_ip_enabled").Exists(),
				check.That(data.ResourceName).Key("tcp_reset_enabled").Exists(),
				check.That(data.ResourceName).Key("disable_outbound_snat").Exists(),
				check.That(data.ResourceName).Key("idle_timeout_in_minutes").Exists(),
				check.That(data.ResourceName).Key("load_distribution").Exists(),
			),
		},
	})
}

func (r LoadBalancerRule) basicDataSource(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_lb_rule" "test" {
  name            = azurerm_lb_rule.test.name
  loadbalancer_id = azurerm_lb_rule.test.loadbalancer_id
}
`, template)
}

func (r LoadBalancerRule) completeDataSource(data acceptance.TestData) string {
	if !features.FivePointOh() {
		return fmt.Sprintf(`
%s
resource "azurerm_lb_backend_address_pool" "test" {
  name            = "acctest-pool-%s"
  loadbalancer_id = azurerm_lb.test.id
}

resource "azurerm_lb_probe" "test" {
  name            = "acctest-probe-%s"
  loadbalancer_id = azurerm_lb.test.id
  protocol        = "Tcp"
  port            = 443
}

resource "azurerm_lb_rule" "test" {
  name            = "acctest-lb-rule-%s"
  loadbalancer_id = azurerm_lb.test.id

  protocol      = "Tcp"
  frontend_port = 3389
  backend_port  = 3389

  disable_outbound_snat   = true
  enable_floating_ip      = true
  enable_tcp_reset        = true
  idle_timeout_in_minutes = 10

  backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
  probe_id                 = azurerm_lb_probe.test.id

  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
}

data "azurerm_lb_rule" "test" {
  name            = azurerm_lb_rule.test.name
  loadbalancer_id = azurerm_lb_rule.test.loadbalancer_id
}
`, r.template(data, "Standard"), data.RandomStringOfLength(8), data.RandomStringOfLength(8), data.RandomStringOfLength(8))
	}

	return fmt.Sprintf(`
%s
resource "azurerm_lb_backend_address_pool" "test" {
  name            = "acctest-pool-%s"
  loadbalancer_id = azurerm_lb.test.id
}

resource "azurerm_lb_probe" "test" {
  name            = "acctest-probe-%s"
  loadbalancer_id = azurerm_lb.test.id
  protocol        = "Tcp"
  port            = 443
}

resource "azurerm_lb_rule" "test" {
  name            = "acctest-lb-rule-%s"
  loadbalancer_id = azurerm_lb.test.id

  protocol      = "Tcp"
  frontend_port = 3389
  backend_port  = 3389

  disable_outbound_snat   = true
  floating_ip_enabled     = true
  tcp_reset_enabled       = true
  idle_timeout_in_minutes = 10

  backend_address_pool_ids = [azurerm_lb_backend_address_pool.test.id]
  probe_id                 = azurerm_lb_probe.test.id

  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
}

data "azurerm_lb_rule" "test" {
  name            = azurerm_lb_rule.test.name
  loadbalancer_id = azurerm_lb_rule.test.loadbalancer_id
}
`, r.template(data, "Standard"), data.RandomStringOfLength(8), data.RandomStringOfLength(8), data.RandomStringOfLength(8))
}

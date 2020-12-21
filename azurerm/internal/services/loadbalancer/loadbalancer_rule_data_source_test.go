package loadbalancer_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

func TestAccAzureRMDataSourceLoadBalancerRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_lb_rule", "test")
	r := LoadBalancerRule{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basicDataSource(data),
			Check: resource.ComposeTestCheckFunc(
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

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.completeDataSource(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("frontend_ip_configuration_name").Exists(),
				check.That(data.ResourceName).Key("protocol").Exists(),
				check.That(data.ResourceName).Key("frontend_port").Exists(),
				check.That(data.ResourceName).Key("backend_port").Exists(),
				check.That(data.ResourceName).Key("backend_address_pool_id").Exists(),
				check.That(data.ResourceName).Key("probe_id").Exists(),
				check.That(data.ResourceName).Key("enable_floating_ip").Exists(),
				check.That(data.ResourceName).Key("enable_tcp_reset").Exists(),
				check.That(data.ResourceName).Key("disable_outbound_snat").Exists(),
				check.That(data.ResourceName).Key("idle_timeout_in_minutes").Exists(),
				check.That(data.ResourceName).Key("load_distribution").Exists(),
			),
		},
	})
}

func (r LoadBalancerRule) basicDataSource(data acceptance.TestData) string {
	template := r.basic(data, "Basic")
	return fmt.Sprintf(`
%s

data "azurerm_lb_rule" "test" {
  name                = azurerm_lb_rule.test.name
  resource_group_name = azurerm_lb_rule.test.resource_group_name
  loadbalancer_id     = azurerm_lb_rule.test.loadbalancer_id
}
`, template)
}

func (r LoadBalancerRule) completeDataSource(data acceptance.TestData) string {
	template := r.template(data, "Standard")
	return fmt.Sprintf(`
%s
resource "azurerm_lb_backend_address_pool" "test" {
  name                = "LbPool-%s"
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
}

resource "azurerm_lb_probe" "test" {
  name                = "LbProbe-%s"
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
  protocol            = "Tcp"
  port                = 443
}

resource "azurerm_lb_rule" "test" {
  name                = "LbRule-%s"
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id

  protocol      = "Tcp"
  frontend_port = 3389
  backend_port  = 3389

  disable_outbound_snat   = true
  enable_floating_ip      = true
  enable_tcp_reset        = true
  idle_timeout_in_minutes = 10

  backend_address_pool_id = azurerm_lb_backend_address_pool.test.id
  probe_id                = azurerm_lb_probe.test.id

  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
}

data "azurerm_lb_rule" "test" {
  name                = azurerm_lb_rule.test.name
  resource_group_name = azurerm_lb_rule.test.resource_group_name
  loadbalancer_id     = azurerm_lb_rule.test.loadbalancer_id
}
`, template, data.RandomStringOfLength(8), data.RandomStringOfLength(8), data.RandomStringOfLength(8))
}

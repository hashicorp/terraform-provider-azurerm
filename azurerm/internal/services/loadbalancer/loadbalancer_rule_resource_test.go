package loadbalancer_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loadbalancer/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type LoadBalancerRule struct {
}

func TestAccAzureRMLoadBalancerRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_rule", "test")
	r := LoadBalancerRule{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, "Basic"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMLoadBalancerRule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_rule", "test")
	r := LoadBalancerRule{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMLoadBalancerRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_rule", "test")
	r := LoadBalancerRule{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, "Basic"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				// TODO - More attributes?
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, "Basic"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMLoadBalancerRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_rule", "test")
	r := LoadBalancerRule{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, "Basic"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccAzureRMLoadBalancerRule_removal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_rule", "test")
	r := LoadBalancerRule{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, "Basic"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.template(data, "Basic"),
			Check: resource.ComposeTestCheckFunc(
				r.IsMissing("azurerm_lb.test", fmt.Sprintf("LbRule-%s", data.RandomStringOfLength(8))),
			),
		},
	})
}

// https://github.com/hashicorp/terraform/issues/9424
func TestAccAzureRMLoadBalancerRule_inconsistentReads(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_rule", "test")
	r := LoadBalancerRule{}
	p := LoadBalancerProbe{}
	b := LoadBalancerBackendAddressPool{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.inconsistentRead(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That("azurerm_lb_probe.test").ExistsInAzure(p),
				check.That("azurerm_lb_backend_address_pool.test").ExistsInAzure(b),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMLoadBalancerRule_updateMultipleRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_rule", "test")
	data2 := acceptance.BuildTestData(t, "azurerm_lb_rule", "test2")
	r := LoadBalancerRule{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.multipleRules(data, data2),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data2.ResourceName).ExistsInAzure(r),
				check.That(data2.ResourceName).Key("frontend_port").HasValue("3390"),
				check.That(data2.ResourceName).Key("backend_port").HasValue("3390"),
			),
		},
		data.ImportStep(),
		data2.ImportStep(),
		{
			Config: r.multipleRulesUpdate(data, data2),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data2.ResourceName).ExistsInAzure(r),
				check.That(data2.ResourceName).Key("frontend_port").HasValue("3391"),
				check.That(data2.ResourceName).Key("backend_port").HasValue("3391"),
			),
		},
		data.ImportStep(),
		data2.ImportStep(),
	})
}

func (r LoadBalancerRule) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.LoadBalancingRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	lb, err := client.LoadBalancers.LoadBalancersClient.Get(ctx, id.ResourceGroup, id.LoadBalancerName, "")
	if err != nil {
		if utils.ResponseWasNotFound(lb.Response) {
			return nil, fmt.Errorf("Load Balancer %q (resource group %q) not found for Load Balancing Rule %q", id.LoadBalancerName, id.ResourceGroup, id.Name)
		}
		return nil, fmt.Errorf("failed reading Load Balancer %q (resource group %q) for Load Balancing Rule %q", id.LoadBalancerName, id.ResourceGroup, id.Name)
	}
	props := lb.LoadBalancerPropertiesFormat
	if props == nil || props.LoadBalancingRules == nil || len(*props.LoadBalancingRules) == 0 {
		return nil, fmt.Errorf("Load Balancing Rule %q not found in Load Balancer %q (resource group %q)", id.Name, id.LoadBalancerName, id.ResourceGroup)
	}

	found := false
	for _, v := range *props.LoadBalancingRules {
		if v.Name != nil && *v.Name == id.Name {
			found = true
		}
	}
	if !found {
		return nil, fmt.Errorf("Load Balancing Rule %q not found in Load Balancer %q (resource group %q)", id.Name, id.LoadBalancerName, id.ResourceGroup)
	}
	return utils.Bool(found), nil
}

func (r LoadBalancerRule) IsMissing(loadBalancerName string, loadBalancingRuleName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).LoadBalancers.LoadBalancersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[loadBalancerName]
		if !ok {
			return fmt.Errorf("not found: %q", loadBalancerName)
		}

		id, err := parse.LoadBalancerID(rs.Primary.ID)
		if err != nil {
			return err
		}

		lb, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
		if err != nil {
			if utils.ResponseWasNotFound(lb.Response) {
				return fmt.Errorf("Load Balancer %q (resource group %q) not found while checking for Load Balancing Rule removal", id.Name, id.ResourceGroup)
			}
			return fmt.Errorf("failed reading Load Balancer %q (resource group %q) for Load Balancing Rule removal", id.Name, id.ResourceGroup)
		}
		props := lb.LoadBalancerPropertiesFormat
		if props == nil || props.LoadBalancingRules == nil {
			return fmt.Errorf("Load Balancing Rule %q not found in Load Balancer %q (resource group %q)", loadBalancingRuleName, id.Name, id.ResourceGroup)
		}

		found := false
		for _, v := range *props.LoadBalancingRules {
			if v.Name != nil && *v.Name == loadBalancingRuleName {
				found = true
			}
		}
		if found {
			return fmt.Errorf("Outbound Rule %q not removed from Load Balancer %q (resource group %q)", loadBalancingRuleName, id.Name, id.ResourceGroup)
		}
		return nil
	}
}

func (r LoadBalancerRule) template(data acceptance.TestData, sku string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-lb-%[1]d"
  location = "%[2]s"
}

resource "azurerm_public_ip" "test" {
  name                = "test-ip-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "%[3]s"
}

resource "azurerm_lb" "test" {
  name                = "arm-test-loadbalancer-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "%[3]s"

  frontend_ip_configuration {
    name                 = "one-%[1]d"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, sku)
}

// nolint: unparam
func (r LoadBalancerRule) basic(data acceptance.TestData, sku string) string {
	template := r.template(data, sku)
	return fmt.Sprintf(`
%s

resource "azurerm_lb_rule" "test" {
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  name                           = "LbRule-%s"
  protocol                       = "Tcp"
  frontend_port                  = 3389
  backend_port                   = 3389
  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
}
`, template, data.RandomStringOfLength(8))
}

func (r LoadBalancerRule) complete(data acceptance.TestData) string {
	template := r.template(data, "Standard")
	return fmt.Sprintf(`
%s

resource "azurerm_lb_rule" "test" {
  name                = "LbRule-%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.test.id}"

  protocol      = "Tcp"
  frontend_port = 3389
  backend_port  = 3389

  disable_outbound_snat   = true
  enable_floating_ip      = true
  enable_tcp_reset        = true
  idle_timeout_in_minutes = 10

  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
}
`, template, data.RandomStringOfLength(8))
}

func (r LoadBalancerRule) requiresImport(data acceptance.TestData) string {
	template := r.basic(data, "Basic")
	return fmt.Sprintf(`
%s

resource "azurerm_lb_rule" "import" {
  name                           = azurerm_lb_rule.test.name
  resource_group_name            = azurerm_lb_rule.test.resource_group_name
  loadbalancer_id                = azurerm_lb_rule.test.loadbalancer_id
  frontend_ip_configuration_name = azurerm_lb_rule.test.frontend_ip_configuration_name
  protocol                       = "Tcp"
  frontend_port                  = 3389
  backend_port                   = 3389
}
`, template)
}

// https://github.com/hashicorp/terraform/issues/9424
func (r LoadBalancerRule) inconsistentRead(data acceptance.TestData) string {
	template := r.template(data, "Basic")
	return fmt.Sprintf(`
%s

resource "azurerm_lb_backend_address_pool" "test" {
  name                = "%d-address-pool"
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.test.id}"
}

resource "azurerm_lb_probe" "test" {
  name                = "probe-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.test.id}"
  protocol            = "Tcp"
  port                = 443
}

resource "azurerm_lb_rule" "test" {
  name                           = "LbRule-%s"
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  protocol                       = "Tcp"
  frontend_port                  = 3389
  backend_port                   = 3389
  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomStringOfLength(8))
}

func (r LoadBalancerRule) multipleRules(data, data2 acceptance.TestData) string {
	template := r.template(data, "Basic")
	return fmt.Sprintf(`
%s

resource "azurerm_lb_rule" "test" {
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  name                           = "LbRule-%s"
  protocol                       = "Udp"
  frontend_port                  = 3389
  backend_port                   = 3389
  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
}

resource "azurerm_lb_rule" "test2" {
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  name                           = "LbRule-%s"
  protocol                       = "Udp"
  frontend_port                  = 3390
  backend_port                   = 3390
  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
}
`, template, data.RandomStringOfLength(8), data2.RandomStringOfLength(8))
}

func (r LoadBalancerRule) multipleRulesUpdate(data, data2 acceptance.TestData) string {
	template := r.template(data, "Basic")
	return fmt.Sprintf(`
%s

resource "azurerm_lb_rule" "test" {
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  name                           = "LbRule-%s"
  protocol                       = "Udp"
  frontend_port                  = 3389
  backend_port                   = 3389
  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
}

resource "azurerm_lb_rule" "test2" {
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  name                           = "LbRule-%s"
  protocol                       = "Udp"
  frontend_port                  = 3391
  backend_port                   = 3391
  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
}
`, template, data.RandomStringOfLength(8), data2.RandomStringOfLength(8))
}

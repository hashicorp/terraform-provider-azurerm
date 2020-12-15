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

type LoadBalancerNatRule struct {
}

func TestAccAzureRMLoadBalancerNatRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_nat_rule", "test")
	r := LoadBalancerNatRule{}

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

func TestAccAzureRMLoadBalancerNatRule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_nat_rule", "test")
	r := LoadBalancerNatRule{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data, "Standard"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMLoadBalancerNatRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_nat_rule", "test")
	r := LoadBalancerNatRule{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, "Standard"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data, "Standard"),
			Check: resource.ComposeTestCheckFunc(
				// TODO - More attributes here?
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, "Standard"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMLoadBalancerNatRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_nat_rule", "test")
	r := LoadBalancerNatRule{}

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

func TestAccAzureRMLoadBalancerNatRule_removal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_nat_rule", "test")
	r := LoadBalancerNatRule{}

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
				r.IsMissing("azurerm_lb.test", fmt.Sprintf("NatRule-%d", data.RandomInteger)),
			),
		},
	})
}

func TestAccAzureRMLoadBalancerNatRule_updateMultipleRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_nat_rule", "test")
	data2 := acceptance.BuildTestData(t, "azurerm_lb_nat_rule", "test2")

	r := LoadBalancerNatRule{}

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
	})
}

func (r LoadBalancerNatRule) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.LoadBalancerInboundNatRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	lb, err := client.LoadBalancers.LoadBalancersClient.Get(ctx, id.ResourceGroup, id.LoadBalancerName, "")
	if err != nil {
		if utils.ResponseWasNotFound(lb.Response) {
			return nil, fmt.Errorf("Load Balancer %q (resource group %q) not found for Nat Rule %q", id.LoadBalancerName, id.ResourceGroup, id.InboundNatRuleName)
		}
		return nil, fmt.Errorf("failed reading Load Balancer %q (resource group %q) for Nat Rule %q", id.LoadBalancerName, id.ResourceGroup, id.InboundNatRuleName)
	}
	props := lb.LoadBalancerPropertiesFormat
	if props == nil || props.InboundNatRules == nil || len(*props.InboundNatRules) == 0 {
		return nil, fmt.Errorf("Nat Rule %q not found in Load Balancer %q (resource group %q)", id.InboundNatRuleName, id.LoadBalancerName, id.ResourceGroup)
	}

	found := false
	for _, v := range *props.InboundNatRules {
		if v.Name != nil && *v.Name == id.InboundNatRuleName {
			found = true
		}
	}
	if !found {
		return nil, fmt.Errorf("Nat Rule %q not found in Load Balancer %q (resource group %q)", id.InboundNatRuleName, id.LoadBalancerName, id.ResourceGroup)
	}
	return utils.Bool(found), nil
}

func (r LoadBalancerNatRule) IsMissing(loadBalancerName string, natRuleName string) resource.TestCheckFunc {
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
				return fmt.Errorf("Load Balancer %q (resource group %q) not found while checking for Nat Rule removal", id.Name, id.ResourceGroup)
			}
			return fmt.Errorf("failed reading Load Balancer %q (resource group %q) for Nat Rule removal", id.Name, id.ResourceGroup)
		}
		props := lb.LoadBalancerPropertiesFormat
		if props == nil || props.InboundNatRules == nil {
			return fmt.Errorf("Nat Rule %q not found in Load Balancer %q (resource group %q)", natRuleName, id.Name, id.ResourceGroup)
		}

		found := false
		for _, v := range *props.InboundNatRules {
			if v.Name != nil && *v.Name == natRuleName {
				found = true
			}
		}
		if found {
			return fmt.Errorf("Nat Rule %q not removed from Load Balancer %q (resource group %q)", natRuleName, id.Name, id.ResourceGroup)
		}
		return nil
	}
}

func (r LoadBalancerNatRule) template(data acceptance.TestData, sku string) string {
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

func (r LoadBalancerNatRule) basic(data acceptance.TestData, sku string) string {
	template := r.template(data, sku)
	return fmt.Sprintf(`
%s

resource "azurerm_lb_nat_rule" "test" {
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  name                           = "NatRule-%d"
  protocol                       = "Tcp"
  frontend_port                  = 3389
  backend_port                   = 3389
  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
}
`, template, data.RandomInteger)
}

func (r LoadBalancerNatRule) complete(data acceptance.TestData, sku string) string {
	template := r.template(data, sku)
	return fmt.Sprintf(`
%s

resource "azurerm_lb_nat_rule" "test" {
  name                = "NatRule-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.test.id}"

  protocol      = "Tcp"
  frontend_port = 3389
  backend_port  = 3389

  enable_floating_ip      = true
  enable_tcp_reset        = true
  idle_timeout_in_minutes = 10

  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
}
`, template, data.RandomInteger)
}

func (r LoadBalancerNatRule) requiresImport(data acceptance.TestData) string {
	template := r.basic(data, "Basic")
	return fmt.Sprintf(`
%s

resource "azurerm_lb_nat_rule" "import" {
  name                           = azurerm_lb_nat_rule.test.name
  loadbalancer_id                = azurerm_lb_nat_rule.test.loadbalancer_id
  resource_group_name            = azurerm_lb_nat_rule.test.resource_group_name
  frontend_ip_configuration_name = azurerm_lb_nat_rule.test.frontend_ip_configuration_name
  protocol                       = "Tcp"
  frontend_port                  = 3389
  backend_port                   = 3389
}
`, template)
}

func (r LoadBalancerNatRule) multipleRules(data, data2 acceptance.TestData) string {
	template := r.template(data, "Basic")
	return fmt.Sprintf(`
%s

resource "azurerm_lb_nat_rule" "test" {
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  name                           = "NatRule-%d"
  protocol                       = "Tcp"
  frontend_port                  = 3389
  backend_port                   = 3389
  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
}

resource "azurerm_lb_nat_rule" "test2" {
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  name                           = "NatRule-%d"
  protocol                       = "Tcp"
  frontend_port                  = 3390
  backend_port                   = 3390
  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
}
`, template, data.RandomInteger, data2.RandomInteger)
}

func (r LoadBalancerNatRule) multipleRulesUpdate(data, data2 acceptance.TestData) string {
	template := r.template(data, "Basic")
	return fmt.Sprintf(`
%s
resource "azurerm_lb_nat_rule" "test" {
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  name                           = "NatRule-%d"
  protocol                       = "Tcp"
  frontend_port                  = 3389
  backend_port                   = 3389
  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
}

resource "azurerm_lb_nat_rule" "test2" {
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  name                           = "NatRule-%d"
  protocol                       = "Tcp"
  frontend_port                  = 3391
  backend_port                   = 3391
  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
}
`, template, data.RandomInteger, data2.RandomInteger)
}

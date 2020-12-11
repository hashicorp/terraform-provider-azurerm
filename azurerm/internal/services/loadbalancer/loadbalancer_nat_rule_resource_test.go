package loadbalancer_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loadbalancer/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMLoadBalancerNatRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_nat_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerNatRule_basic(data, "Basic"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerNatRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLoadBalancerNatRule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_nat_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerNatRule_complete(data, "Standard"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerNatRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLoadBalancerNatRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_nat_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerNatRule_basic(data, "Standard"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerNatRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMLoadBalancerNatRule_complete(data, "Standard"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerNatRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMLoadBalancerNatRule_basic(data, "Standard"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerNatRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLoadBalancerNatRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_nat_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerNatRule_basic(data, "Basic"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerNatRuleExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMLoadBalancerNatRule_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_lb_nat_rule"),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerNatRule_removal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_nat_rule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerNatRule_basic(data, "Basic"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerNatRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMLoadBalancerNatRule_template(data, "Basic"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerNatRuleIsMissing("azurerm_lb.test", fmt.Sprintf("NatRule-%d", data.RandomInteger)),
				),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerNatRule_updateMultipleRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_nat_rule", "test")
	data2 := acceptance.BuildTestData(t, "azurerm_lb_nat_rule", "test2")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerNatRule_multipleRules(data, data2),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerNatRuleExists(data.ResourceName),
					testCheckAzureRMLoadBalancerNatRuleExists(data2.ResourceName),
					resource.TestCheckResourceAttr(data2.ResourceName, "frontend_port", "3390"),
					resource.TestCheckResourceAttr(data2.ResourceName, "backend_port", "3390"),
				),
			},
			{
				Config: testAccAzureRMLoadBalancerNatRule_multipleRulesUpdate(data, data2),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerNatRuleExists(data.ResourceName),
					testCheckAzureRMLoadBalancerNatRuleExists(data2.ResourceName),
					resource.TestCheckResourceAttr(data2.ResourceName, "frontend_port", "3391"),
					resource.TestCheckResourceAttr(data2.ResourceName, "backend_port", "3391"),
				),
			},
		},
	})
}

func testCheckAzureRMLoadBalancerNatRuleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).LoadBalancers.LoadBalancersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %q", resourceName)
		}

		id, err := parse.LoadBalancerInboundNatRuleID(rs.Primary.ID)
		if err != nil {
			return err
		}

		lb, err := client.Get(ctx, id.ResourceGroup, id.LoadBalancerName, "")
		if err != nil {
			if utils.ResponseWasNotFound(lb.Response) {
				return fmt.Errorf("Load Balancer %q (resource group %q) not found for Nat Rule %q", id.LoadBalancerName, id.ResourceGroup, id.InboundNatRuleName)
			}
			return fmt.Errorf("failed reading Load Balancer %q (resource group %q) for Nat Rule %q", id.LoadBalancerName, id.ResourceGroup, id.InboundNatRuleName)
		}
		props := lb.LoadBalancerPropertiesFormat
		if props == nil || props.InboundNatRules == nil || len(*props.InboundNatRules) == 0 {
			return fmt.Errorf("Backend Pool %q not found in Load Balancer %q (resource group %q)", id.InboundNatRuleName, id.LoadBalancerName, id.ResourceGroup)
		}

		found := false
		for _, v := range *props.InboundNatRules {
			if v.Name != nil && *v.Name == id.InboundNatRuleName {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("Backend Pool %q not found in Load Balancer %q (resource group %q)", id.InboundNatRuleName, id.LoadBalancerName, id.ResourceGroup)
		}
		return nil

	}
}

func testCheckAzureRMLoadBalancerNatRuleIsMissing(loadBalancerName string, natRuleName string) resource.TestCheckFunc {
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

func testAccAzureRMLoadBalancerNatRule_template(data acceptance.TestData, sku string) string {
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

func testAccAzureRMLoadBalancerNatRule_basic(data acceptance.TestData, sku string) string {
	template := testAccAzureRMLoadBalancerNatRule_template(data, sku)
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

func testAccAzureRMLoadBalancerNatRule_complete(data acceptance.TestData, sku string) string {
	template := testAccAzureRMLoadBalancerNatRule_template(data, sku)
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

func testAccAzureRMLoadBalancerNatRule_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMLoadBalancerNatRule_basic(data, "Basic")
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

func testAccAzureRMLoadBalancerNatRule_multipleRules(data, data2 acceptance.TestData) string {
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
`, testAccAzureRMLoadBalancerNatRule_template(data, "Basic"), data.RandomInteger, data2.RandomInteger)
}

func testAccAzureRMLoadBalancerNatRule_multipleRulesUpdate(data, data2 acceptance.TestData) string {
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
`, testAccAzureRMLoadBalancerNatRule_template(data, "Basic"), data.RandomInteger, data2.RandomInteger)
}

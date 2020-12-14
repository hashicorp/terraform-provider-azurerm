package loadbalancer_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-03-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loadbalancer"
)

func TestAccAzureRMLoadBalancerNatRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_nat_rule", "test")

	var lb network.LoadBalancer
	natRuleName := fmt.Sprintf("NatRule-%d", data.RandomInteger)

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	natRuleId := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/acctestRG-lb-%d/providers/Microsoft.Network/loadBalancers/arm-test-loadbalancer-%d/inboundNatRules/%s",
		subscriptionID, data.RandomInteger, data.RandomInteger, natRuleName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerNatRule_basic(data, natRuleName, "Basic"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerNatRuleExists(natRuleName, &lb),
					resource.TestCheckResourceAttr(data.ResourceName, "id", natRuleId),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLoadBalancerNatRule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_nat_rule", "test")

	var lb network.LoadBalancer
	natRuleName := fmt.Sprintf("NatRule-%d", data.RandomInteger)

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	natRuleId := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/acctestRG-lb-%d/providers/Microsoft.Network/loadBalancers/arm-test-loadbalancer-%d/inboundNatRules/%s",
		subscriptionID, data.RandomInteger, data.RandomInteger, natRuleName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerNatRule_complete(data, natRuleName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerNatRuleExists(natRuleName, &lb),
					resource.TestCheckResourceAttr(data.ResourceName, "id", natRuleId),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLoadBalancerNatRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_nat_rule", "test")

	var lb network.LoadBalancer
	natRuleName := fmt.Sprintf("NatRule-%d", data.RandomInteger)

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	natRuleId := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/acctestRG-lb-%d/providers/Microsoft.Network/loadBalancers/arm-test-loadbalancer-%d/inboundNatRules/%s",
		subscriptionID, data.RandomInteger, data.RandomInteger, natRuleName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerNatRule_basic(data, natRuleName, "Standard"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerNatRuleExists(natRuleName, &lb),
					resource.TestCheckResourceAttr(data.ResourceName, "id", natRuleId),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMLoadBalancerNatRule_complete(data, natRuleName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerNatRuleExists(natRuleName, &lb),
					resource.TestCheckResourceAttr(
						"azurerm_lb_nat_rule.test", "id", natRuleId),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMLoadBalancerNatRule_basic(data, natRuleName, "Standard"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerNatRuleExists(natRuleName, &lb),
					resource.TestCheckResourceAttr(data.ResourceName, "id", natRuleId),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLoadBalancerNatRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_nat_rule", "test")

	var lb network.LoadBalancer
	natRuleName := fmt.Sprintf("NatRule-%d", data.RandomInteger)

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	natRuleId := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/acctestRG-lb-%d/providers/Microsoft.Network/loadBalancers/arm-test-loadbalancer-%d/inboundNatRules/%s",
		subscriptionID, data.RandomInteger, data.RandomInteger, natRuleName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerNatRule_basic(data, natRuleName, "Basic"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerNatRuleExists(natRuleName, &lb),
					resource.TestCheckResourceAttr(data.ResourceName, "id", natRuleId),
				),
			},
			{
				Config:      testAccAzureRMLoadBalancerNatRule_requiresImport(data, natRuleName),
				ExpectError: acceptance.RequiresImportError("azurerm_lb_nat_rule"),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerNatRule_removal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_nat_rule", "test")
	var lb network.LoadBalancer
	natRuleName := fmt.Sprintf("NatRule-%d", data.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerNatRule_basic(data, natRuleName, "Basic"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerNatRuleExists(natRuleName, &lb),
				),
			},
			{
				Config: testAccAzureRMLoadBalancerNatRule_template(data, "Basic"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerNatRuleNotExists(natRuleName, &lb),
				),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerNatRule_updateMultipleRules(t *testing.T) {
	data1 := acceptance.BuildTestData(t, "azurerm_lb_nat_rule", "test")
	data2 := acceptance.BuildTestData(t, "azurerm_lb_nat_rule", "test2")
	var lb network.LoadBalancer
	natRuleName := fmt.Sprintf("NatRule-%d", data1.RandomInteger)
	natRule2Name := fmt.Sprintf("NatRule-%d", data2.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerNatRule_multipleRules(data1, natRuleName, natRule2Name),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerNatRuleExists(natRuleName, &lb),
					testCheckAzureRMLoadBalancerNatRuleExists(natRule2Name, &lb),
					resource.TestCheckResourceAttr(data2.ResourceName, "frontend_port", "3390"),
					resource.TestCheckResourceAttr(data2.ResourceName, "backend_port", "3390"),
				),
			},
			{
				Config: testAccAzureRMLoadBalancerNatRule_multipleRulesUpdate(data1, natRuleName, natRule2Name),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerNatRuleExists(natRuleName, &lb),
					testCheckAzureRMLoadBalancerNatRuleExists(natRule2Name, &lb),
					resource.TestCheckResourceAttr(data2.ResourceName, "frontend_port", "3391"),
					resource.TestCheckResourceAttr(data2.ResourceName, "backend_port", "3391"),
				),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerNatRule_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_nat_rule", "test")

	var lb network.LoadBalancer
	natRuleName := fmt.Sprintf("NatRule-%d", data.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerNatRule_basic(data, natRuleName, "Basic"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerNatRuleExists(natRuleName, &lb),
					testCheckAzureRMLoadBalancerNatRuleDisappears(natRuleName, &lb),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMLoadBalancerNatRuleExists(natRuleName string, lb *network.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, _, exists := loadbalancer.FindLoadBalancerNatRuleByName(lb, natRuleName)
		if !exists {
			return fmt.Errorf("A NAT Rule with name %q cannot be found.", natRuleName)
		}

		return nil
	}
}

func testCheckAzureRMLoadBalancerNatRuleNotExists(natRuleName string, lb *network.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, _, exists := loadbalancer.FindLoadBalancerNatRuleByName(lb, natRuleName)
		if exists {
			return fmt.Errorf("A NAT Rule with name %q has been found.", natRuleName)
		}

		return nil
	}
}

func testCheckAzureRMLoadBalancerNatRuleDisappears(natRuleName string, lb *network.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).LoadBalancers.LoadBalancersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		_, i, exists := loadbalancer.FindLoadBalancerNatRuleByName(lb, natRuleName)
		if !exists {
			return fmt.Errorf("A Nat Rule with name %q cannot be found.", natRuleName)
		}

		currentRules := *lb.LoadBalancerPropertiesFormat.InboundNatRules
		rules := append(currentRules[:i], currentRules[i+1:]...)
		lb.LoadBalancerPropertiesFormat.InboundNatRules = &rules

		id, err := azure.ParseAzureResourceID(*lb.ID)
		if err != nil {
			return err
		}

		future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, *lb.Name, *lb)
		if err != nil {
			return fmt.Errorf("Error Creating/Updating Load Balancer %+v", err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for the completion of Load Balancer %q (Resource Group %q): %+v", *lb.Name, id.ResourceGroup, err)
		}

		_, err = client.Get(ctx, id.ResourceGroup, *lb.Name, "")
		return err
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

func testAccAzureRMLoadBalancerNatRule_basic(data acceptance.TestData, natRuleName string, sku string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_lb_nat_rule" "test" {
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  name                           = "%s"
  protocol                       = "Tcp"
  frontend_port                  = 3389
  backend_port                   = 3389
  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
}
`, testAccAzureRMLoadBalancerNatRule_template(data, sku), natRuleName)
}

func testAccAzureRMLoadBalancerNatRule_complete(data acceptance.TestData, natRuleName string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_lb_nat_rule" "test" {
  name                = "%s"
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
`, testAccAzureRMLoadBalancerNatRule_template(data, "Standard"), natRuleName)
}

func testAccAzureRMLoadBalancerNatRule_requiresImport(data acceptance.TestData, name string) string {
	template := testAccAzureRMLoadBalancerNatRule_basic(data, name, "Basic")
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

func testAccAzureRMLoadBalancerNatRule_multipleRules(data acceptance.TestData, natRuleName, natRule2Name string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_lb_nat_rule" "test" {
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  name                           = "%s"
  protocol                       = "Tcp"
  frontend_port                  = 3389
  backend_port                   = 3389
  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
}

resource "azurerm_lb_nat_rule" "test2" {
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  name                           = "%s"
  protocol                       = "Tcp"
  frontend_port                  = 3390
  backend_port                   = 3390
  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
}
`, testAccAzureRMLoadBalancerNatRule_template(data, "Basic"), natRuleName, natRule2Name)
}

func testAccAzureRMLoadBalancerNatRule_multipleRulesUpdate(data acceptance.TestData, natRuleName, natRule2Name string) string {
	return fmt.Sprintf(`
%s
resource "azurerm_lb_nat_rule" "test" {
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  name                           = "%s"
  protocol                       = "Tcp"
  frontend_port                  = 3389
  backend_port                   = 3389
  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
}

resource "azurerm_lb_nat_rule" "test2" {
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  name                           = "%s"
  protocol                       = "Tcp"
  frontend_port                  = 3391
  backend_port                   = 3391
  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
}
`, testAccAzureRMLoadBalancerNatRule_template(data, "Basic"), natRuleName, natRule2Name)
}

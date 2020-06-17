package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-03-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	nw "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network"
)

func TestResourceAzureRMLoadBalancerRuleNameLabel_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "-word",
			ErrCount: 1,
		},
		{
			Value:    "testing-",
			ErrCount: 1,
		},
		{
			Value:    "test#test",
			ErrCount: 1,
		},
		{
			Value:    acctest.RandStringFromCharSet(81, "abcdedfed"),
			ErrCount: 1,
		},
		{
			Value:    "test.rule",
			ErrCount: 0,
		},
		{
			Value:    "test_rule",
			ErrCount: 0,
		},
		{
			Value:    "test-rule",
			ErrCount: 0,
		},
		{
			Value:    "TestRule",
			ErrCount: 0,
		},
		{
			Value:    "Test123Rule",
			ErrCount: 0,
		},
		{
			Value:    "TestRule",
			ErrCount: 0,
		},
	}

	for _, tc := range cases {
		_, errors := nw.ValidateArmLoadBalancerRuleName(tc.Value, "azurerm_lb_rule")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Azure RM Load Balancer Rule Name Label to trigger a validation error")
		}
	}
}

func TestAccAzureRMLoadBalancerRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_rule", "test")
	var lb network.LoadBalancer
	lbRuleName := fmt.Sprintf("LbRule-%s", acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	lbRule_id := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/acctestRG-lb-%d/providers/Microsoft.Network/loadBalancers/arm-test-loadbalancer-%d/loadBalancingRules/%s",
		subscriptionID, data.RandomInteger, data.RandomInteger, lbRuleName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerRule_basic(data, lbRuleName, "Basic"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerRuleExists(lbRuleName, &lb),
					resource.TestCheckResourceAttr(
						"azurerm_lb_rule.test", "id", lbRule_id),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLoadBalancerRule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_rule", "test")
	var lb network.LoadBalancer
	lbRuleName := fmt.Sprintf("LbRule-%s", acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	lbRule_id := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/acctestRG-lb-%d/providers/Microsoft.Network/loadBalancers/arm-test-loadbalancer-%d/loadBalancingRules/%s",
		subscriptionID, data.RandomInteger, data.RandomInteger, lbRuleName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerRule_complete(data, lbRuleName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerRuleExists(lbRuleName, &lb),
					resource.TestCheckResourceAttr(
						"azurerm_lb_rule.test", "id", lbRule_id),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLoadBalancerRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_rule", "test")
	var lb network.LoadBalancer
	lbRuleName := fmt.Sprintf("LbRule-%s", acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	lbRule_id := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/acctestRG-lb-%d/providers/Microsoft.Network/loadBalancers/arm-test-loadbalancer-%d/loadBalancingRules/%s",
		subscriptionID, data.RandomInteger, data.RandomInteger, lbRuleName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerRule_basic(data, lbRuleName, "Basic"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerRuleExists(lbRuleName, &lb),
					resource.TestCheckResourceAttr(
						"azurerm_lb_rule.test", "id", lbRule_id),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMLoadBalancerRule_complete(data, lbRuleName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerRuleExists(lbRuleName, &lb),
					resource.TestCheckResourceAttr(
						"azurerm_lb_rule.test", "id", lbRule_id),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMLoadBalancerRule_basic(data, lbRuleName, "Basic"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerRuleExists(lbRuleName, &lb),
					resource.TestCheckResourceAttr(
						"azurerm_lb_rule.test", "id", lbRule_id),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLoadBalancerRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb", "test")

	var lb network.LoadBalancer
	lbRuleName := fmt.Sprintf("LbRule-%s", acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	lbRule_id := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/acctestRG-lb-%d/providers/Microsoft.Network/loadBalancers/arm-test-loadbalancer-%d/loadBalancingRules/%s",
		subscriptionID, data.RandomInteger, data.RandomInteger, lbRuleName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerRule_basic(data, lbRuleName, "Basic"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerRuleExists(lbRuleName, &lb),
					resource.TestCheckResourceAttr(
						"azurerm_lb_rule.test", "id", lbRule_id),
				),
			},
			{
				Config:      testAccAzureRMLoadBalancerRule_requiresImport(data, lbRuleName),
				ExpectError: acceptance.RequiresImportError("azurerm_lb_rule"),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerRule_removal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_rule", "test")
	var lb network.LoadBalancer
	lbRuleName := fmt.Sprintf("LbRule-%s", acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerRule_basic(data, lbRuleName, "Basic"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerRuleExists(lbRuleName, &lb),
				),
			},
			{
				Config: testAccAzureRMLoadBalancerRule_template(data, "Basic"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerRuleNotExists(lbRuleName, &lb),
				),
			},
		},
	})
}

// https://github.com/hashicorp/terraform/issues/9424
func TestAccAzureRMLoadBalancerRule_inconsistentReads(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_rule", "test")
	var lb network.LoadBalancer
	backendPoolName := fmt.Sprintf("LbPool-%s", acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	lbRuleName := fmt.Sprintf("LbRule-%s", acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	probeName := fmt.Sprintf("LbProbe-%s", acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerRule_inconsistentRead(data, backendPoolName, probeName, lbRuleName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerBackEndAddressPoolExists(backendPoolName, &lb),
					testCheckAzureRMLoadBalancerRuleExists(lbRuleName, &lb),
					testCheckAzureRMLoadBalancerProbeExists(probeName, &lb),
				),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerRule_updateMultipleRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_rule", "test")

	var lb network.LoadBalancer
	lbRuleName := fmt.Sprintf("LbRule-%s", acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	lbRule2Name := fmt.Sprintf("LbRule-%s", acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	lbRuleID := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/acctestRG-lb-%d/providers/Microsoft.Network/loadBalancers/arm-test-loadbalancer-%d/loadBalancingRules/%s",
		subscriptionID, data.RandomInteger, data.RandomInteger, lbRuleName)

	lbRule2ID := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/acctestRG-lb-%d/providers/Microsoft.Network/loadBalancers/arm-test-loadbalancer-%d/loadBalancingRules/%s",
		subscriptionID, data.RandomInteger, data.RandomInteger, lbRule2Name)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerRule_multipleRules(data, lbRuleName, lbRule2Name),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerRuleExists(lbRuleName, &lb),
					testCheckAzureRMLoadBalancerRuleExists(lbRule2Name, &lb),
					resource.TestCheckResourceAttr("azurerm_lb_rule.test", "id", lbRuleID),
					resource.TestCheckResourceAttr("azurerm_lb_rule.test2", "id", lbRule2ID),
					resource.TestCheckResourceAttr("azurerm_lb_rule.test2", "frontend_port", "3390"),
					resource.TestCheckResourceAttr("azurerm_lb_rule.test2", "backend_port", "3390"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMLoadBalancerRule_multipleRulesUpdate(data, lbRuleName, lbRule2Name),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerRuleExists(lbRuleName, &lb),
					testCheckAzureRMLoadBalancerRuleExists(lbRule2Name, &lb),
					resource.TestCheckResourceAttr("azurerm_lb_rule.test", "id", lbRuleID),
					resource.TestCheckResourceAttr("azurerm_lb_rule.test2", "id", lbRule2ID),
					resource.TestCheckResourceAttr("azurerm_lb_rule.test2", "frontend_port", "3391"),
					resource.TestCheckResourceAttr("azurerm_lb_rule.test2", "backend_port", "3391"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMLoadBalancerRule_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb", "test")

	var lb network.LoadBalancer
	lbRuleName := fmt.Sprintf("LbRule-%s", acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerRule_basic(data, lbRuleName, "Basic"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerRuleExists(lbRuleName, &lb),
					testCheckAzureRMLoadBalancerRuleDisappears(lbRuleName, &lb),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMLoadBalancerRuleExists(lbRuleName string, lb *network.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, _, exists := nw.FindLoadBalancerRuleByName(lb, lbRuleName)
		if !exists {
			return fmt.Errorf("A Load Balancer Rule with name %q cannot be found.", lbRuleName)
		}

		return nil
	}
}

func testCheckAzureRMLoadBalancerRuleNotExists(lbRuleName string, lb *network.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, _, exists := nw.FindLoadBalancerRuleByName(lb, lbRuleName)
		if exists {
			return fmt.Errorf("A Load Balancer Rule with name %q has been found.", lbRuleName)
		}

		return nil
	}
}

func testCheckAzureRMLoadBalancerRuleDisappears(ruleName string, lb *network.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.LoadBalancersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		_, i, exists := nw.FindLoadBalancerRuleByName(lb, ruleName)
		if !exists {
			return fmt.Errorf("A Rule with name %q cannot be found.", ruleName)
		}

		currentRules := *lb.LoadBalancerPropertiesFormat.LoadBalancingRules
		rules := append(currentRules[:i], currentRules[i+1:]...)
		lb.LoadBalancerPropertiesFormat.LoadBalancingRules = &rules

		id, err := azure.ParseAzureResourceID(*lb.ID)
		if err != nil {
			return err
		}

		future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, *lb.Name, *lb)
		if err != nil {
			return fmt.Errorf("Error Creating/Updating Load Balancer %q (Resource Group %q): %+v", *lb.Name, id.ResourceGroup, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for completion of Load Balancer %q (Resource Group %q): %+v", *lb.Name, id.ResourceGroup, err)
		}

		_, err = client.Get(ctx, id.ResourceGroup, *lb.Name, "")
		return err
	}
}

func testAccAzureRMLoadBalancerRule_template(data acceptance.TestData, sku string) string {
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
func testAccAzureRMLoadBalancerRule_basic(data acceptance.TestData, lbRuleName, sku string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_lb_rule" "test" {
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  name                           = "%s"
  protocol                       = "Tcp"
  frontend_port                  = 3389
  backend_port                   = 3389
  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
}
`, testAccAzureRMLoadBalancerRule_template(data, sku), lbRuleName)
}

func testAccAzureRMLoadBalancerRule_complete(data acceptance.TestData, lbRuleName string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_lb_rule" "test" {
  name                = "%s"
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
`, testAccAzureRMLoadBalancerRule_template(data, "Standard"), lbRuleName)
}

func testAccAzureRMLoadBalancerRule_requiresImport(data acceptance.TestData, name string) string {
	template := testAccAzureRMLoadBalancerRule_basic(data, name, "Basic")
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
func testAccAzureRMLoadBalancerRule_inconsistentRead(data acceptance.TestData, backendPoolName, probeName, lbRuleName string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_lb_backend_address_pool" "teset" {
  name                = "%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.test.id}"
}

resource "azurerm_lb_probe" "test" {
  name                = "%s"
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.test.id}"
  protocol            = "Tcp"
  port                = 443
}

resource "azurerm_lb_rule" "test" {
  name                           = "%s"
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  protocol                       = "Tcp"
  frontend_port                  = 3389
  backend_port                   = 3389
  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
}
`, testAccAzureRMLoadBalancerRule_template(data, "Basic"), backendPoolName, probeName, lbRuleName)
}

func testAccAzureRMLoadBalancerRule_multipleRules(data acceptance.TestData, lbRuleName, lbRule2Name string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_lb_rule" "test" {
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  name                           = "%s"
  protocol                       = "Udp"
  frontend_port                  = 3389
  backend_port                   = 3389
  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
}

resource "azurerm_lb_rule" "test2" {
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  name                           = "%s"
  protocol                       = "Udp"
  frontend_port                  = 3390
  backend_port                   = 3390
  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
}
`, testAccAzureRMLoadBalancerRule_template(data, "Basic"), lbRuleName, lbRule2Name)
}

func testAccAzureRMLoadBalancerRule_multipleRulesUpdate(data acceptance.TestData, lbRuleName, lbRule2Name string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_lb_rule" "test" {
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  name                           = "%s"
  protocol                       = "Udp"
  frontend_port                  = 3389
  backend_port                   = 3389
  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
}

resource "azurerm_lb_rule" "test2" {
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  name                           = "%s"
  protocol                       = "Udp"
  frontend_port                  = 3391
  backend_port                   = 3391
  frontend_ip_configuration_name = azurerm_lb.test.frontend_ip_configuration.0.name
}
`, testAccAzureRMLoadBalancerRule_template(data, "Basic"), lbRuleName, lbRule2Name)
}

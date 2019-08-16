package azurerm

import (
	"fmt"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-12-01/network"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccAzureRMLoadBalancerNatRule_basic(t *testing.T) {
	var lb network.LoadBalancer
	ri := tf.AccRandTimeInt()
	natRuleName := fmt.Sprintf("NatRule-%d", ri)

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	natRuleId := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/acctestRG-%d/providers/Microsoft.Network/loadBalancers/arm-test-loadbalancer-%d/inboundNatRules/%s",
		subscriptionID, ri, ri, natRuleName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerNatRule_basic(ri, natRuleName, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerNatRuleExists(natRuleName, &lb),
					resource.TestCheckResourceAttr(
						"azurerm_lb_nat_rule.test", "id", natRuleId),
				),
			},
			{
				ResourceName:      "azurerm_lb.test",
				ImportState:       true,
				ImportStateVerify: true,
				// location is deprecated and was never actually used
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func TestAccAzureRMLoadBalancerNatRule_requiresImport(t *testing.T) {
	if !requireResourcesToBeImported {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	var lb network.LoadBalancer
	ri := tf.AccRandTimeInt()
	natRuleName := fmt.Sprintf("NatRule-%d", ri)
	location := testLocation()

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	natRuleId := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/acctestRG-%d/providers/Microsoft.Network/loadBalancers/arm-test-loadbalancer-%d/inboundNatRules/%s",
		subscriptionID, ri, ri, natRuleName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerNatRule_basic(ri, natRuleName, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerNatRuleExists(natRuleName, &lb),
					resource.TestCheckResourceAttr(
						"azurerm_lb_nat_rule.test", "id", natRuleId),
				),
			},
			{
				Config:      testAccAzureRMLoadBalancerNatRule_requiresImport(ri, natRuleName, location),
				ExpectError: testRequiresImportError("azurerm_lb_nat_rule"),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerNatRule_removal(t *testing.T) {
	var lb network.LoadBalancer
	ri := tf.AccRandTimeInt()
	natRuleName := fmt.Sprintf("NatRule-%d", ri)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerNatRule_basic(ri, natRuleName, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerNatRuleExists(natRuleName, &lb),
				),
			},
			{
				Config: testAccAzureRMLoadBalancerNatRule_removal(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerNatRuleNotExists(natRuleName, &lb),
				),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerNatRule_update(t *testing.T) {
	var lb network.LoadBalancer
	ri := tf.AccRandTimeInt()
	natRuleName := fmt.Sprintf("NatRule-%d", ri)
	natRule2Name := fmt.Sprintf("NatRule-%d", tf.AccRandTimeInt())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerNatRule_multipleRules(ri, natRuleName, natRule2Name, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerNatRuleExists(natRuleName, &lb),
					testCheckAzureRMLoadBalancerNatRuleExists(natRule2Name, &lb),
					resource.TestCheckResourceAttr("azurerm_lb_nat_rule.test2", "frontend_port", "3390"),
					resource.TestCheckResourceAttr("azurerm_lb_nat_rule.test2", "backend_port", "3390"),
				),
			},
			{
				Config: testAccAzureRMLoadBalancerNatRule_multipleRulesUpdate(ri, natRuleName, natRule2Name, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerNatRuleExists(natRuleName, &lb),
					testCheckAzureRMLoadBalancerNatRuleExists(natRule2Name, &lb),
					resource.TestCheckResourceAttr("azurerm_lb_nat_rule.test2", "frontend_port", "3391"),
					resource.TestCheckResourceAttr("azurerm_lb_nat_rule.test2", "backend_port", "3391"),
				),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerNatRule_disappears(t *testing.T) {
	var lb network.LoadBalancer
	ri := tf.AccRandTimeInt()
	natRuleName := fmt.Sprintf("NatRule-%d", ri)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerNatRule_basic(ri, natRuleName, testLocation()),
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

func TestAccAzureRMLoadBalancerNatRule_enableFloatingIP(t *testing.T) {
	var lb network.LoadBalancer
	ri := tf.AccRandTimeInt()
	natRuleName := fmt.Sprintf("NatRule-%d", ri)
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerNatRule_enableFloatingIP(ri, natRuleName, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerNatRuleExists(natRuleName, &lb),
				),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerNatRule_disableFloatingIP(t *testing.T) {
	var lb network.LoadBalancer
	ri := tf.AccRandTimeInt()
	natRuleName := fmt.Sprintf("NatRule-%d", ri)
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerNatRule_basic(ri, natRuleName, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerNatRuleExists(natRuleName, &lb),
				),
			},
			{
				Config: testAccAzureRMLoadBalancerNatRule_enableFloatingIP(ri, natRuleName, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerNatRuleExists(natRuleName, &lb),
				),
			},
			{
				Config: testAccAzureRMLoadBalancerNatRule_basic(ri, natRuleName, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerNatRuleExists(natRuleName, &lb),
				),
			},
		},
	})
}

func testCheckAzureRMLoadBalancerNatRuleExists(natRuleName string, lb *network.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, _, exists := findLoadBalancerNatRuleByName(lb, natRuleName)
		if !exists {
			return fmt.Errorf("A NAT Rule with name %q cannot be found.", natRuleName)
		}

		return nil
	}
}

func testCheckAzureRMLoadBalancerNatRuleNotExists(natRuleName string, lb *network.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, _, exists := findLoadBalancerNatRuleByName(lb, natRuleName)
		if exists {
			return fmt.Errorf("A NAT Rule with name %q has been found.", natRuleName)
		}

		return nil
	}
}

func testCheckAzureRMLoadBalancerNatRuleDisappears(natRuleName string, lb *network.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ArmClient).network.LoadBalancersClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		_, i, exists := findLoadBalancerNatRuleByName(lb, natRuleName)
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

func testAccAzureRMLoadBalancerNatRule_basic(rInt int, natRuleName string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "test-ip-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "arm-test-loadbalancer-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  frontend_ip_configuration {
    name                 = "one-%d"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }
}

resource "azurerm_lb_nat_rule" "test" {
  location                       = "${azurerm_resource_group.test.location}"
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  name                           = "%s"
  protocol                       = "Tcp"
  frontend_port                  = 3389
  backend_port                   = 3389
  frontend_ip_configuration_name = "one-%d"
}
`, rInt, location, rInt, rInt, rInt, natRuleName, rInt)
}

func testAccAzureRMLoadBalancerNatRule_requiresImport(rInt int, name string, location string) string {
	template := testAccAzureRMLoadBalancerNatRule_basic(rInt, name, location)
	return fmt.Sprintf(`
%s

resource "azurerm_lb_nat_rule" "import" {
  name                           = "${azurerm_lb_nat_rule.test.name}"
  loadbalancer_id                = "${azurerm_lb_nat_rule.test.loadbalancer_id}"
  location                       = "${azurerm_lb_nat_rule.test.location}"
  resource_group_name            = "${azurerm_lb_nat_rule.test.resource_group_name}"
  frontend_ip_configuration_name = "${azurerm_lb_nat_rule.test.frontend_ip_configuration_name}"
  protocol                       = "Tcp"
  frontend_port                  = 3389
  backend_port                   = 3389
}
`, template)
}

func testAccAzureRMLoadBalancerNatRule_removal(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "test-ip-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "arm-test-loadbalancer-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  frontend_ip_configuration {
    name                 = "one-%d"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMLoadBalancerNatRule_multipleRules(rInt int, natRuleName, natRule2Name string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "test-ip-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "arm-test-loadbalancer-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  frontend_ip_configuration {
    name                 = "one-%d"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }
}

resource "azurerm_lb_nat_rule" "test" {
  location                       = "${azurerm_resource_group.test.location}"
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  name                           = "%s"
  protocol                       = "Tcp"
  frontend_port                  = 3389
  backend_port                   = 3389
  frontend_ip_configuration_name = "one-%d"
}

resource "azurerm_lb_nat_rule" "test2" {
  location                       = "${azurerm_resource_group.test.location}"
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  name                           = "%s"
  protocol                       = "Tcp"
  frontend_port                  = 3390
  backend_port                   = 3390
  frontend_ip_configuration_name = "one-%d"
}
`, rInt, location, rInt, rInt, rInt, natRuleName, rInt, natRule2Name, rInt)
}

func testAccAzureRMLoadBalancerNatRule_multipleRulesUpdate(rInt int, natRuleName, natRule2Name string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "test-ip-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "arm-test-loadbalancer-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  frontend_ip_configuration {
    name                 = "one-%d"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }
}

resource "azurerm_lb_nat_rule" "test" {
  location                       = "${azurerm_resource_group.test.location}"
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  name                           = "%s"
  protocol                       = "Tcp"
  frontend_port                  = 3389
  backend_port                   = 3389
  frontend_ip_configuration_name = "one-%d"
}

resource "azurerm_lb_nat_rule" "test2" {
  location                       = "${azurerm_resource_group.test.location}"
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  name                           = "%s"
  protocol                       = "Tcp"
  frontend_port                  = 3391
  backend_port                   = 3391
  frontend_ip_configuration_name = "one-%d"
}
`, rInt, location, rInt, rInt, rInt, natRuleName, rInt, natRule2Name, rInt)
}

func testAccAzureRMLoadBalancerNatRule_enableFloatingIP(rInt int, natRuleName string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "test-ip-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
}

resource "azurerm_lb" "test" {
  name                = "arm-test-loadbalancer-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  frontend_ip_configuration {
    name                 = "one-%d"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }
}

resource "azurerm_lb_nat_rule" "test" {
  location                       = "${azurerm_resource_group.test.location}"
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  name                           = "%s"
  protocol                       = "Tcp"
  frontend_port                  = 3389
  backend_port                   = 3389
  frontend_ip_configuration_name = "one-%d"
}
`, rInt, location, rInt, rInt, rInt, natRuleName, rInt)
}

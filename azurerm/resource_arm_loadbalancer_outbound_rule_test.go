package azurerm

import (
	"fmt"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMLoadBalancerOutboundRule_basic(t *testing.T) {
	var lb network.LoadBalancer
	ri := tf.AccRandTimeInt()
	outboundRuleName := fmt.Sprintf("OutboundRule-%d", ri)

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	outboundRuleId := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/acctestRG-%d/providers/Microsoft.Network/loadBalancers/arm-test-loadbalancer-%d/outboundRules/%s",
		subscriptionID, ri, ri, outboundRuleName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerOutboundRule_basic(ri, outboundRuleName, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerOutboundRuleExists(outboundRuleName, &lb),
					resource.TestCheckResourceAttr(
						"azurerm_lb_outbound_rule.test", "id", outboundRuleId),
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

func TestAccAzureRMLoadBalancerOutboundRule_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	var lb network.LoadBalancer
	ri := tf.AccRandTimeInt()
	outboundRuleName := fmt.Sprintf("OutboundRule-%d", ri)
	location := acceptance.Location()

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	outboundRuleId := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/acctestRG-%d/providers/Microsoft.Network/loadBalancers/arm-test-loadbalancer-%d/outboundRules/%s",
		subscriptionID, ri, ri, outboundRuleName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerOutboundRule_basic(ri, outboundRuleName, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerOutboundRuleExists(outboundRuleName, &lb),
					resource.TestCheckResourceAttr(
						"azurerm_lb_outbound_rule.test", "id", outboundRuleId),
				),
			},
			{
				Config:      testAccAzureRMLoadBalancerOutboundRule_requiresImport(ri, outboundRuleName, location),
				ExpectError: acceptance.RequiresImportError("azurerm_lb_outbound_rule"),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerOutboundRule_removal(t *testing.T) {
	var lb network.LoadBalancer
	ri := tf.AccRandTimeInt()
	outboundRuleName := fmt.Sprintf("OutboundRule-%d", ri)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerOutboundRule_basic(ri, outboundRuleName, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerOutboundRuleExists(outboundRuleName, &lb),
				),
			},
			{
				Config: testAccAzureRMLoadBalancerOutboundRule_removal(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerOutboundRuleNotExists(outboundRuleName, &lb),
				),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerOutboundRule_update(t *testing.T) {
	var lb network.LoadBalancer
	ri := tf.AccRandTimeInt()
	outboundRuleName := fmt.Sprintf("OutboundRule-%d", ri)
	outboundRule2Name := fmt.Sprintf("OutboundRule-%d", tf.AccRandTimeInt())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerOutboundRule_multipleRules(ri, outboundRuleName, outboundRule2Name, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerOutboundRuleExists(outboundRuleName, &lb),
					testCheckAzureRMLoadBalancerOutboundRuleExists(outboundRule2Name, &lb),
					resource.TestCheckResourceAttr("azurerm_lb_outbound_rule.test2", "protocol", "Udp"),
				),
			},
			{
				Config: testAccAzureRMLoadBalancerOutboundRule_multipleRulesUpdate(ri, outboundRuleName, outboundRule2Name, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerOutboundRuleExists(outboundRuleName, &lb),
					testCheckAzureRMLoadBalancerOutboundRuleExists(outboundRule2Name, &lb),
					resource.TestCheckResourceAttr("azurerm_lb_outbound_rule.test2", "protocol", "All"),
				),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerOutboundRule_withPublicIPPrefix(t *testing.T) {
	var lb network.LoadBalancer
	ri := tf.AccRandTimeInt()
	outboundRuleName := fmt.Sprintf("OutboundRule-%d", ri)

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	outboundRuleId := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/acctestRG-%d/providers/Microsoft.Network/loadBalancers/arm-test-loadbalancer-%d/outboundRules/%s",
		subscriptionID, ri, ri, outboundRuleName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerOutboundRule_withPublicIPPrefix(ri, outboundRuleName, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerOutboundRuleExists(outboundRuleName, &lb),
					resource.TestCheckResourceAttr(
						"azurerm_lb_outbound_rule.test", "id", outboundRuleId),
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

func TestAccAzureRMLoadBalancerOutboundRule_disappears(t *testing.T) {
	var lb network.LoadBalancer
	ri := tf.AccRandTimeInt()
	outboundRuleName := fmt.Sprintf("OutboundRule-%d", ri)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerOutboundRule_basic(ri, outboundRuleName, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerOutboundRuleExists(outboundRuleName, &lb),
					testCheckAzureRMLoadBalancerOutboundRuleDisappears(outboundRuleName, &lb),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMLoadBalancerOutboundRuleExists(outboundRuleName string, lb *network.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if _, _, exists := findLoadBalancerOutboundRuleByName(lb, outboundRuleName); !exists {
			return fmt.Errorf("A Load Balancer Outbound Rule with name %q cannot be found.", outboundRuleName)
		}

		return nil
	}
}

func testCheckAzureRMLoadBalancerOutboundRuleNotExists(outboundRuleName string, lb *network.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if _, _, exists := findLoadBalancerOutboundRuleByName(lb, outboundRuleName); exists {
			return fmt.Errorf("A Load Balancer Outbound Rule with name %q has been found.", outboundRuleName)
		}

		return nil
	}
}

func testCheckAzureRMLoadBalancerOutboundRuleDisappears(ruleName string, lb *network.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.LoadBalancersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		_, i, exists := findLoadBalancerOutboundRuleByName(lb, ruleName)
		if !exists {
			return fmt.Errorf("A Outbound Rule with name %q cannot be found.", ruleName)
		}

		currentRules := *lb.LoadBalancerPropertiesFormat.OutboundRules
		rules := append(currentRules[:i], currentRules[i+1:]...)
		lb.LoadBalancerPropertiesFormat.OutboundRules = &rules

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

func testAccAzureRMLoadBalancerOutboundRule_basic(rInt int, outboundRuleName string, location string) string {
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
  sku                 = "Standard"
}

resource "azurerm_lb" "test" {
  name                = "arm-test-loadbalancer-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "one-%d"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.test.id}"
  name                = "be-%d"
}

resource "azurerm_lb_outbound_rule" "test" {
  resource_group_name     = "${azurerm_resource_group.test.name}"
  loadbalancer_id         = "${azurerm_lb.test.id}"
  name                    = "%s"
  backend_address_pool_id = "${azurerm_lb_backend_address_pool.test.id}"
  protocol                = "All"

  frontend_ip_configuration {
    name = "one-%d"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt, outboundRuleName, rInt)
}

func testAccAzureRMLoadBalancerOutboundRule_requiresImport(rInt int, name string, location string) string {
	template := testAccAzureRMLoadBalancerOutboundRule_basic(rInt, name, location)
	return fmt.Sprintf(`
%s

resource "azurerm_lb_outbound_rule" "import" {
  name                    = "${azurerm_lb_outbound_rule.test.name}"
  resource_group_name     = "${azurerm_lb_outbound_rule.test.resource_group_name}"
  loadbalancer_id         = "${azurerm_lb_outbound_rule.test.loadbalancer_id}"
  backend_address_pool_id = "${azurerm_lb_backend_address_pool.test.id}"
  protocol                = "All"

  frontend_ip_configuration {
    name = "${azurerm_lb_outbound_rule.test.frontend_ip_configuration.0.name}"
  }
}
`, template)
}

func testAccAzureRMLoadBalancerOutboundRule_removal(rInt int, location string) string {
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
  sku                 = "Standard"
}

resource "azurerm_lb_backend_address_pool" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.test.id}"
  name                = "be-%d"
}

resource "azurerm_lb" "test" {
  name                = "arm-test-loadbalancer-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "one-%d"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt)
}

func testAccAzureRMLoadBalancerOutboundRule_multipleRules(rInt int, outboundRuleName, outboundRule2Name string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test1" {
  name                = "test-ip-1-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_public_ip" "test2" {
  name                = "test-ip-2-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_lb" "test" {
  name                = "arm-test-loadbalancer-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "fe1-%d"
    public_ip_address_id = "${azurerm_public_ip.test1.id}"
  }

  frontend_ip_configuration {
    name                 = "fe2-%d"
    public_ip_address_id = "${azurerm_public_ip.test2.id}"
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.test.id}"
  name                = "be-%d"
}

resource "azurerm_lb_outbound_rule" "test" {
  resource_group_name     = "${azurerm_resource_group.test.name}"
  loadbalancer_id         = "${azurerm_lb.test.id}"
  name                    = "%s"
  protocol                = "Tcp"
  backend_address_pool_id = "${azurerm_lb_backend_address_pool.test.id}"

  frontend_ip_configuration {
    name = "fe1-%d"
  }
}

resource "azurerm_lb_outbound_rule" "test2" {
  resource_group_name     = "${azurerm_resource_group.test.name}"
  loadbalancer_id         = "${azurerm_lb.test.id}"
  name                    = "%s"
  protocol                = "Udp"
  backend_address_pool_id = "${azurerm_lb_backend_address_pool.test.id}"

  frontend_ip_configuration {
    name = "fe2-%d"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt, outboundRuleName, rInt, outboundRule2Name, rInt)
}

func testAccAzureRMLoadBalancerOutboundRule_multipleRulesUpdate(rInt int, outboundRuleName, outboundRule2Name string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test1" {
  name                = "test-ip-1-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_public_ip" "test2" {
  name                = "test-ip-2-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_lb" "test" {
  name                = "arm-test-loadbalancer-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "fe1-%d"
    public_ip_address_id = "${azurerm_public_ip.test1.id}"
  }

  frontend_ip_configuration {
    name                 = "fe2-%d"
    public_ip_address_id = "${azurerm_public_ip.test2.id}"
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.test.id}"
  name                = "be-%d"
}

resource "azurerm_lb_outbound_rule" "test" {
  resource_group_name     = "${azurerm_resource_group.test.name}"
  loadbalancer_id         = "${azurerm_lb.test.id}"
  name                    = "%s"
  protocol                = "All"
  backend_address_pool_id = "${azurerm_lb_backend_address_pool.test.id}"

  frontend_ip_configuration {
    name = "fe1-%d"
  }
}

resource "azurerm_lb_outbound_rule" "test2" {
  resource_group_name     = "${azurerm_resource_group.test.name}"
  loadbalancer_id         = "${azurerm_lb.test.id}"
  name                    = "%s"
  protocol                = "All"
  backend_address_pool_id = "${azurerm_lb_backend_address_pool.test.id}"

  frontend_ip_configuration {
    name = "fe2-%d"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt, rInt, rInt, outboundRuleName, rInt, outboundRule2Name, rInt)
}

func testAccAzureRMLoadBalancerOutboundRule_withPublicIPPrefix(rInt int, outboundRuleName string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip_prefix" "test" {
  name                = "test-ip-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  prefix_length       = 31
}

resource "azurerm_lb" "test" {
  name                = "arm-test-loadbalancer-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Standard"

  frontend_ip_configuration {
    name                = "one-%d"
    public_ip_prefix_id = "${azurerm_public_ip_prefix.test.id}"
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.test.id}"
  name                = "be-%d"
}

resource "azurerm_lb_outbound_rule" "test" {
  resource_group_name     = "${azurerm_resource_group.test.name}"
  loadbalancer_id         = "${azurerm_lb.test.id}"
  name                    = "%s"
  backend_address_pool_id = "${azurerm_lb_backend_address_pool.test.id}"
  protocol                = "All"

  frontend_ip_configuration {
    name = "one-%d"
  }
}
`, rInt, location, rInt, rInt, rInt, rInt, outboundRuleName, rInt)
}

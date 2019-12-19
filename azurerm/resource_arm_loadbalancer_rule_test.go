package azurerm

import (
	"fmt"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
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
		_, errors := validateArmLoadBalancerRuleName(tc.Value, "azurerm_lb_rule")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Azure RM Load Balancer Rule Name Label to trigger a validation error")
		}
	}
}

func TestAccAzureRMLoadBalancerRule_basic(t *testing.T) {
	var lb network.LoadBalancer
	ri := tf.AccRandTimeInt()
	lbRuleName := fmt.Sprintf("LbRule-%s", acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	lbRule_id := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/acctestRG-%d/providers/Microsoft.Network/loadBalancers/arm-test-loadbalancer-%d/loadBalancingRules/%s",
		subscriptionID, ri, ri, lbRuleName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerRule_basic(ri, lbRuleName, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerRuleExists(lbRuleName, &lb),
					resource.TestCheckResourceAttr(
						"azurerm_lb_rule.test", "id", lbRule_id),
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

func TestAccAzureRMLoadBalancerRule_disableoutboundsnat(t *testing.T) {
	var lb network.LoadBalancer
	ri := tf.AccRandTimeInt()
	lbRuleName := fmt.Sprintf("LbRule-%s", acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	lbRule_id := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/acctestRG-%d/providers/Microsoft.Network/loadBalancers/arm-test-loadbalancer-%d/loadBalancingRules/%s",
		subscriptionID, ri, ri, lbRuleName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerRule_disableoutboundsnat(ri, lbRuleName, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerRuleExists(lbRuleName, &lb),
					resource.TestCheckResourceAttr("azurerm_lb_rule.test", "id", lbRule_id),
					resource.TestCheckResourceAttr("azurerm_lb_rule.test", "disable_outbound_snat", "true"),
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

func TestAccAzureRMLoadBalancerRule_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	var lb network.LoadBalancer
	ri := tf.AccRandTimeInt()
	lbRuleName := fmt.Sprintf("LbRule-%s", acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	location := acceptance.Location()

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	lbRule_id := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/acctestRG-%d/providers/Microsoft.Network/loadBalancers/arm-test-loadbalancer-%d/loadBalancingRules/%s",
		subscriptionID, ri, ri, lbRuleName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerRule_basic(ri, lbRuleName, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerRuleExists(lbRuleName, &lb),
					resource.TestCheckResourceAttr(
						"azurerm_lb_rule.test", "id", lbRule_id),
				),
			},
			{
				Config:      testAccAzureRMLoadBalancerRule_requiresImport(ri, lbRuleName, location),
				ExpectError: acceptance.RequiresImportError("azurerm_lb_rule"),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerRule_removal(t *testing.T) {
	var lb network.LoadBalancer
	ri := tf.AccRandTimeInt()
	lbRuleName := fmt.Sprintf("LbRule-%s", acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerRule_basic(ri, lbRuleName, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerRuleExists(lbRuleName, &lb),
				),
			},
			{
				Config: testAccAzureRMLoadBalancerRule_removal(ri, acceptance.Location()),
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
	var lb network.LoadBalancer
	ri := tf.AccRandTimeInt()
	backendPoolName := fmt.Sprintf("LbPool-%s", acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	lbRuleName := fmt.Sprintf("LbRule-%s", acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	probeName := fmt.Sprintf("LbProbe-%s", acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerRule_inconsistentRead(ri, backendPoolName, probeName, lbRuleName, acceptance.Location()),
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

func TestAccAzureRMLoadBalancerRule_update(t *testing.T) {
	var lb network.LoadBalancer
	ri := tf.AccRandTimeInt()
	lbRuleName := fmt.Sprintf("LbRule-%s", acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))
	lbRule2Name := fmt.Sprintf("LbRule-%s", acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	lbRuleID := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/acctestRG-%d/providers/Microsoft.Network/loadBalancers/arm-test-loadbalancer-%d/loadBalancingRules/%s",
		subscriptionID, ri, ri, lbRuleName)

	lbRule2ID := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/acctestRG-%d/providers/Microsoft.Network/loadBalancers/arm-test-loadbalancer-%d/loadBalancingRules/%s",
		subscriptionID, ri, ri, lbRule2Name)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerRule_multipleRules(ri, lbRuleName, lbRule2Name, acceptance.Location()),
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
			{
				Config: testAccAzureRMLoadBalancerRule_multipleRulesUpdate(ri, lbRuleName, lbRule2Name, acceptance.Location()),
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
		},
	})
}

func TestAccAzureRMLoadBalancerRule_disappears(t *testing.T) {
	var lb network.LoadBalancer
	ri := tf.AccRandTimeInt()
	lbRuleName := fmt.Sprintf("LbRule-%s", acctest.RandStringFromCharSet(8, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerRule_basic(ri, lbRuleName, acceptance.Location()),
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
		_, _, exists := findLoadBalancerRuleByName(lb, lbRuleName)
		if !exists {
			return fmt.Errorf("A Load Balancer Rule with name %q cannot be found.", lbRuleName)
		}

		return nil
	}
}

func testCheckAzureRMLoadBalancerRuleNotExists(lbRuleName string, lb *network.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, _, exists := findLoadBalancerRuleByName(lb, lbRuleName)
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

		_, i, exists := findLoadBalancerRuleByName(lb, ruleName)
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

func testAccAzureRMLoadBalancerRule_basic(rInt int, lbRuleName string, location string) string {
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

resource "azurerm_lb_rule" "test" {
  location                       = "${azurerm_resource_group.test.location}"
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  name                           = "%s"
  protocol                       = "Tcp"
  frontend_port                  = 3389
  backend_port                   = 3389
  frontend_ip_configuration_name = "one-%d"
}
`, rInt, location, rInt, rInt, rInt, lbRuleName, rInt)
}

func testAccAzureRMLoadBalancerRule_disableoutboundsnat(rInt int, lbRuleName string, location string) string {
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

resource "azurerm_lb_rule" "test" {
  location                       = "${azurerm_resource_group.test.location}"
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  name                           = "%s"
  protocol                       = "Tcp"
  frontend_port                  = 3389
  backend_port                   = 3389
  frontend_ip_configuration_name = "one-%d"
  disable_outbound_snat          = true
}
`, rInt, location, rInt, rInt, rInt, lbRuleName, rInt)
}

func testAccAzureRMLoadBalancerRule_requiresImport(rInt int, name string, location string) string {
	template := testAccAzureRMLoadBalancerRule_basic(rInt, name, location)
	return fmt.Sprintf(`
%s

resource "azurerm_lb_rule" "import" {
  name                           = "${azurerm_lb_rule.test.name}"
  location                       = "${azurerm_lb_rule.test.location}"
  resource_group_name            = "${azurerm_lb_rule.test.resource_group_name}"
  loadbalancer_id                = "${azurerm_lb_rule.test.loadbalancer_id}"
  frontend_ip_configuration_name = "${azurerm_lb_rule.test.frontend_ip_configuration_name}"
  protocol                       = "Tcp"
  frontend_port                  = 3389
  backend_port                   = 3389
}
`, template)
}

func testAccAzureRMLoadBalancerRule_removal(rInt int, location string) string {
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

// https://github.com/hashicorp/terraform/issues/9424
func testAccAzureRMLoadBalancerRule_inconsistentRead(rInt int, backendPoolName, probeName, lbRuleName string, location string) string {
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

resource "azurerm_lb_backend_address_pool" "teset" {
  name                = "%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.test.id}"
}

resource "azurerm_lb_probe" "test" {
  name                = "%s"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id     = "${azurerm_lb.test.id}"
  protocol            = "Tcp"
  port                = 443
}

resource "azurerm_lb_rule" "test" {
  name                           = "%s"
  location                       = "${azurerm_resource_group.test.location}"
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  protocol                       = "Tcp"
  frontend_port                  = 3389
  backend_port                   = 3389
  frontend_ip_configuration_name = "one-%d"
}
`, rInt, location, rInt, rInt, rInt, backendPoolName, probeName, lbRuleName, rInt)
}

func testAccAzureRMLoadBalancerRule_multipleRules(rInt int, lbRuleName, lbRule2Name string, location string) string {
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

resource "azurerm_lb_rule" "test" {
  location                       = "${azurerm_resource_group.test.location}"
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  name                           = "%s"
  protocol                       = "Udp"
  frontend_port                  = 3389
  backend_port                   = 3389
  frontend_ip_configuration_name = "one-%d"
}

resource "azurerm_lb_rule" "test2" {
  location                       = "${azurerm_resource_group.test.location}"
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  name                           = "%s"
  protocol                       = "Udp"
  frontend_port                  = 3390
  backend_port                   = 3390
  frontend_ip_configuration_name = "one-%d"
}
`, rInt, location, rInt, rInt, rInt, lbRuleName, rInt, lbRule2Name, rInt)
}

func testAccAzureRMLoadBalancerRule_multipleRulesUpdate(rInt int, lbRuleName, lbRule2Name string, location string) string {
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

resource "azurerm_lb_rule" "test" {
  location                       = "${azurerm_resource_group.test.location}"
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  name                           = "%s"
  protocol                       = "Udp"
  frontend_port                  = 3389
  backend_port                   = 3389
  frontend_ip_configuration_name = "one-%d"
}

resource "azurerm_lb_rule" "test2" {
  location                       = "${azurerm_resource_group.test.location}"
  resource_group_name            = "${azurerm_resource_group.test.name}"
  loadbalancer_id                = "${azurerm_lb.test.id}"
  name                           = "%s"
  protocol                       = "Udp"
  frontend_port                  = 3391
  backend_port                   = 3391
  frontend_ip_configuration_name = "one-%d"
}
`, rInt, location, rInt, rInt, rInt, lbRuleName, rInt, lbRule2Name, rInt)
}

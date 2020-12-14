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

func TestAccAzureRMLoadBalancerOutboundRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_nat_rule", "test")

	var lb network.LoadBalancer
	outboundRuleName := fmt.Sprintf("OutboundRule-%d", data.RandomInteger)

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	outboundRuleId := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/acctestRG-%d/providers/Microsoft.Network/loadBalancers/arm-test-loadbalancer-%d/outboundRules/%s",
		subscriptionID, data.RandomInteger, data.RandomInteger, outboundRuleName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerOutboundRule_basic(data, outboundRuleName),
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
	data := acceptance.BuildTestData(t, "azurerm_lb_nat_rule", "test")

	var lb network.LoadBalancer
	outboundRuleName := fmt.Sprintf("OutboundRule-%d", data.RandomInteger)

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	outboundRuleId := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/acctestRG-%d/providers/Microsoft.Network/loadBalancers/arm-test-loadbalancer-%d/outboundRules/%s",
		subscriptionID, data.RandomInteger, data.RandomInteger, outboundRuleName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerOutboundRule_basic(data, outboundRuleName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerOutboundRuleExists(outboundRuleName, &lb),
					resource.TestCheckResourceAttr(
						"azurerm_lb_outbound_rule.test", "id", outboundRuleId),
				),
			},
			{
				Config:      testAccAzureRMLoadBalancerOutboundRule_requiresImport(data, outboundRuleName),
				ExpectError: acceptance.RequiresImportError("azurerm_lb_outbound_rule"),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerOutboundRule_removal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_nat_rule", "test")

	var lb network.LoadBalancer
	outboundRuleName := fmt.Sprintf("OutboundRule-%d", data.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerOutboundRule_basic(data, outboundRuleName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerOutboundRuleExists(outboundRuleName, &lb),
				),
			},
			{
				Config: testAccAzureRMLoadBalancerOutboundRule_removal(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerOutboundRuleNotExists(outboundRuleName, &lb),
				),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerOutboundRule_update(t *testing.T) {
	data1 := acceptance.BuildTestData(t, "azurerm_lb_outbound_rule", "test")
	data2 := acceptance.BuildTestData(t, "azurerm_lb_outbound_rule", "test2")

	var lb network.LoadBalancer
	outboundRuleName := fmt.Sprintf("OutboundRule-%d", data1.RandomInteger)
	outboundRule2Name := fmt.Sprintf("OutboundRule-%d", data2.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerOutboundRule_multipleRules(data1, outboundRuleName, outboundRule2Name),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerOutboundRuleExists(outboundRuleName, &lb),
					testCheckAzureRMLoadBalancerOutboundRuleExists(outboundRule2Name, &lb),
				),
			},
			data1.ImportStep(),
			data2.ImportStep(),
			{
				Config: testAccAzureRMLoadBalancerOutboundRule_multipleRulesUpdate(data1, outboundRuleName, outboundRule2Name),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerOutboundRuleExists(outboundRuleName, &lb),
					testCheckAzureRMLoadBalancerOutboundRuleExists(outboundRule2Name, &lb),
				),
			},
			data1.ImportStep(),
			data2.ImportStep(),
		},
	})
}

func TestAccAzureRMLoadBalancerOutboundRule_withPublicIPPrefix(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lb_nat_rule", "test")

	var lb network.LoadBalancer
	outboundRuleName := fmt.Sprintf("OutboundRule-%d", data.RandomInteger)

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	outboundRuleId := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/acctestRG-%d/providers/Microsoft.Network/loadBalancers/arm-test-loadbalancer-%d/outboundRules/%s",
		subscriptionID, data.RandomInteger, data.RandomInteger, outboundRuleName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerOutboundRule_withPublicIPPrefix(data, outboundRuleName),
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
	data := acceptance.BuildTestData(t, "azurerm_lb_nat_rule", "test")

	var lb network.LoadBalancer
	outboundRuleName := fmt.Sprintf("OutboundRule-%d", data.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerOutboundRule_basic(data, outboundRuleName),
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
		if _, _, exists := loadbalancer.FindLoadBalancerOutboundRuleByName(lb, outboundRuleName); !exists {
			return fmt.Errorf("A Load Balancer Outbound Rule with name %q cannot be found.", outboundRuleName)
		}

		return nil
	}
}

func testCheckAzureRMLoadBalancerOutboundRuleNotExists(outboundRuleName string, lb *network.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if _, _, exists := loadbalancer.FindLoadBalancerOutboundRuleByName(lb, outboundRuleName); exists {
			return fmt.Errorf("A Load Balancer Outbound Rule with name %q has been found.", outboundRuleName)
		}

		return nil
	}
}

func testCheckAzureRMLoadBalancerOutboundRuleDisappears(ruleName string, lb *network.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).LoadBalancers.LoadBalancersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		_, i, exists := loadbalancer.FindLoadBalancerOutboundRuleByName(lb, ruleName)
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

func testAccAzureRMLoadBalancerOutboundRule_basic(data acceptance.TestData, outboundRuleName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "test-ip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_lb" "test" {
  name                = "arm-test-loadbalancer-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "one-%d"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
  name                = "be-%d"
}

resource "azurerm_lb_outbound_rule" "test" {
  resource_group_name     = azurerm_resource_group.test.name
  loadbalancer_id         = azurerm_lb.test.id
  name                    = "%s"
  backend_address_pool_id = azurerm_lb_backend_address_pool.test.id
  protocol                = "All"

  frontend_ip_configuration {
    name = "one-%d"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, outboundRuleName, data.RandomInteger)
}

func testAccAzureRMLoadBalancerOutboundRule_requiresImport(data acceptance.TestData, name string) string {
	template := testAccAzureRMLoadBalancerOutboundRule_basic(data, name)
	return fmt.Sprintf(`
%s

resource "azurerm_lb_outbound_rule" "import" {
  name                    = azurerm_lb_outbound_rule.test.name
  resource_group_name     = azurerm_lb_outbound_rule.test.resource_group_name
  loadbalancer_id         = azurerm_lb_outbound_rule.test.loadbalancer_id
  backend_address_pool_id = azurerm_lb_backend_address_pool.test.id
  protocol                = "All"

  frontend_ip_configuration {
    name = azurerm_lb_outbound_rule.test.frontend_ip_configuration[0].name
  }
}
`, template)
}

func testAccAzureRMLoadBalancerOutboundRule_removal(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "test-ip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_lb_backend_address_pool" "test" {
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
  name                = "be-%d"
}

resource "azurerm_lb" "test" {
  name                = "arm-test-loadbalancer-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "one-%d"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMLoadBalancerOutboundRule_multipleRules(data acceptance.TestData, outboundRuleName, outboundRule2Name string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test1" {
  name                = "test-ip-1-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_public_ip" "test2" {
  name                = "test-ip-2-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_lb" "test" {
  name                = "arm-test-loadbalancer-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "fe1-%d"
    public_ip_address_id = azurerm_public_ip.test1.id
  }

  frontend_ip_configuration {
    name                 = "fe2-%d"
    public_ip_address_id = azurerm_public_ip.test2.id
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
  name                = "be-%d"
}

resource "azurerm_lb_outbound_rule" "test" {
  resource_group_name     = azurerm_resource_group.test.name
  loadbalancer_id         = azurerm_lb.test.id
  name                    = "%s"
  protocol                = "Tcp"
  backend_address_pool_id = azurerm_lb_backend_address_pool.test.id

  frontend_ip_configuration {
    name = "fe1-%d"
  }
}

resource "azurerm_lb_outbound_rule" "test2" {
  resource_group_name     = azurerm_resource_group.test.name
  loadbalancer_id         = azurerm_lb.test.id
  name                    = "%s"
  protocol                = "Udp"
  backend_address_pool_id = azurerm_lb_backend_address_pool.test.id

  frontend_ip_configuration {
    name = "fe2-%d"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, outboundRuleName, data.RandomInteger, outboundRule2Name, data.RandomInteger)
}

func testAccAzureRMLoadBalancerOutboundRule_multipleRulesUpdate(data acceptance.TestData, outboundRuleName, outboundRule2Name string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test1" {
  name                = "test-ip-1-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_public_ip" "test2" {
  name                = "test-ip-2-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_lb" "test" {
  name                = "arm-test-loadbalancer-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "fe1-%d"
    public_ip_address_id = azurerm_public_ip.test1.id
  }

  frontend_ip_configuration {
    name                 = "fe2-%d"
    public_ip_address_id = azurerm_public_ip.test2.id
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
  name                = "be-%d"
}

resource "azurerm_lb_outbound_rule" "test" {
  resource_group_name     = azurerm_resource_group.test.name
  loadbalancer_id         = azurerm_lb.test.id
  name                    = "%s"
  protocol                = "All"
  backend_address_pool_id = azurerm_lb_backend_address_pool.test.id

  frontend_ip_configuration {
    name = "fe1-%d"
  }
}

resource "azurerm_lb_outbound_rule" "test2" {
  resource_group_name     = azurerm_resource_group.test.name
  loadbalancer_id         = azurerm_lb.test.id
  name                    = "%s"
  protocol                = "All"
  backend_address_pool_id = azurerm_lb_backend_address_pool.test.id

  frontend_ip_configuration {
    name = "fe2-%d"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, outboundRuleName, data.RandomInteger, outboundRule2Name, data.RandomInteger)
}

func testAccAzureRMLoadBalancerOutboundRule_withPublicIPPrefix(data acceptance.TestData, outboundRuleName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip_prefix" "test" {
  name                = "test-ip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  prefix_length       = 31
}

resource "azurerm_lb" "test" {
  name                = "arm-test-loadbalancer-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  frontend_ip_configuration {
    name                = "one-%d"
    public_ip_prefix_id = azurerm_public_ip_prefix.test.id
  }
}

resource "azurerm_lb_backend_address_pool" "test" {
  resource_group_name = azurerm_resource_group.test.name
  loadbalancer_id     = azurerm_lb.test.id
  name                = "be-%d"
}

resource "azurerm_lb_outbound_rule" "test" {
  resource_group_name     = azurerm_resource_group.test.name
  loadbalancer_id         = azurerm_lb.test.id
  name                    = "%s"
  backend_address_pool_id = azurerm_lb_backend_address_pool.test.id
  protocol                = "All"

  frontend_ip_configuration {
    name = "one-%d"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, outboundRuleName, data.RandomInteger)
}
